package ultraViolet

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

// trackers is the list of Tracker constructors that are available on this system. Tracker implementations should call
// the RegisterTracker function to make themselves available.
var trackers = make(map[string]func() Tracker)

// RegisterTracker makes a Tracker constructor available to clients of this package.
func RegisterTracker(name string, t func() Tracker) error {
	if _, exists := trackers[name]; exists {
		return errors.New("a tracker already exists with the name" + name)
	}
	trackers[name] = t
	return nil
}

// NewTracker returns a new Tracker instance whose type is `name`.
func NewTracker(name string) (Tracker, error) {
	if _, exists := trackers[name]; !exists {
		log.Println("error")
		return nil, errors.New("no Tracker constructor has been registered with name " + name)
	}
	return trackers[name](), nil
}

// Tracker tracks application usage. An implementation that satisfies
// this interface is required for each OS windowing system Thyme
// supports.
type Tracker interface {
	// Snap returns a Snapshot reflecting the currently in-use windows
	// at the current time.
	Snap() (*Snapshot, error)

	// Deps returns a string listing the dependencies that still need
	// to be installed with instructions for how to install them.
	Deps() string
}

// Stream represents all the sampling data gathered by Thyme.
type Stream struct {
	// Snapshots is a list of window snapshots ordered by time.
	Snapshots []*Snapshot
}

// Print returns a pretty-printed representation of the snapshot.
func (s Stream) Print() string {
	var b bytes.Buffer
	for _, snap := range s.Snapshots {
		fmt.Fprintf(&b, "%s", snap.Print())
	}
	return string(b.Bytes())
}

// Snapshot represents the current state of all in-use application
// windows at a moment in time.
type Snapshot struct {
	Time    time.Time
	Windows []*Window
	Active  int
	Visible []int
}

// Print returns a pretty-printed representation of the snapshot.
func (s Snapshot) Print() string {
	var b bytes.Buffer

	var active *Window
	visible := make([]*Window, 0, len(s.Windows))
	other := make([]*Window, 0, len(s.Windows))
s_Windows:
	for _, w := range s.Windows {
		if w.ID == s.Active {
			active = w
			continue s_Windows
		}
		for _, v := range s.Visible {
			if w.ID == v {
				visible = append(visible, w)
				continue s_Windows
			}
		}
		other = append(other, w)
	}

	fmt.Fprintf(&b, "%s\n", s.Time.Format("Mon Jan 2 15:04:05 -0700 MST 2006"))
	if active != nil {
		fmt.Fprintf(&b, "\tActive: %s\n", active.Info().Print())
	}
	writeWindows(&b, visible, "Visible")
	writeWindows(&b, other, "Other")
	return string(b.Bytes())
}

func writeWindows(b *bytes.Buffer, windows []*Window, name string) {
	if len(windows) > 0 {
		fmt.Fprintf(b, "\t%s: ", name)
		for _, w := range windows {
			fmt.Fprintf(b, "%s, ", w.Info().Print())
		}
		fmt.Fprintf(b, "\n")
	}
}

// Window represents an application window.
type Window struct {
	// ID is the numerical identifier of the window.
	ID int

	// Desktop is the numerical identifier of the desktop the
	// window belongs to.  Equal to -1 if the window is sticky.
	Desktop int

	// Name is the display name of the window (typically what the
	// windowing system shows in the top bar of the window).
	Name string
}

// IsSticky returns true if the window is a sticky window (i.e.
// present on all desktops)
func (w *Window) IsSticky() bool {
	return w.Desktop == -1
}

// IsOnDesktop returns true if the window is present on the specified
// desktop
func (w *Window) IsOnDesktop(desktop int) bool {
	return w.IsSticky() || w.Desktop == desktop
}

const (
	defaultWindowTitleSeparator       = " - "
	microsoftEdgeWindowTitleSeparator = "\u200e- "
)

// Info returns more structured metadata about a window. The metadata
// is extracted using heuristics.
//
// Assumptions:
//     1) Most windows use " - " to separate their window names from their content
//     2) Most windows use the " - " with the application name at the end.
//     3) The few programs that reverse this convention only reverse it.
func (w *Window) Info() *Winfo {
	// Special Cases
	wi, isChrome := chromeInfo(w.Name)
	if isChrome {
		return wi
	}

	// Normal Cases
	if beforeSep := strings.Index(w.Name, defaultWindowTitleSeparator); beforeSep > -1 && beforeSep < len(w.Name) {

		// parameter of slackInfo() must be validated.
		wi, isSlack := slackInfo(w.Name, beforeSep)
		if isSlack {
			return wi
		}

		//parameter of sepDefault() must be validated.
		return sepDefault(w.Name)
	}

	// No Application name separator
	return &Winfo{
		Title: w.Name,
	}
}

// Winfo is structured metadata info about a window.
type Winfo struct {
	// App is the application that controls the window.
	App string

	// SubApp is the sub-application that controls the window. An
	// example is a web app (e.g., Sourcegraph) that runs
	// inside a Chrome tab. In this case, the App field would be
	// "Google Chrome" and the SubApp field would be "Sourcegraph".
	SubApp string

	// Title is the title of the window after the App and SubApp name
	// have been stripped.
	Title string
}

// Print returns a pretty-printed representation of the snapshot.
func (w Winfo) Print() string {
	return fmt.Sprintf("[%s|%s|%s]", w.App, w.SubApp, w.Title)
}

func chromeInfo(wName string) (wi *Winfo, isChrome bool) {
	fields := strings.Split(wName, defaultWindowTitleSeparator)
	if len(fields) > 1 {
		last := strings.TrimSpace(fields[len(fields)-1])
		if last == "Google Chrome" {
			return &Winfo{
				App:    "Google Chrome",
				SubApp: strings.TrimSpace(fields[len(fields)-2]),
				Title:  strings.Join(fields[0:len(fields)-2], defaultWindowTitleSeparator),
			}, true
		}
	}
	return nil, false
}

// beforeSep must be varidated.
func slackInfo(wName string, beforeSep int) (wi *Winfo, isSlack bool) {
	// App Name First
	if wName[:beforeSep] == "Slack" {
		afterSep := beforeSep + len(defaultWindowTitleSeparator)
		return &Winfo{
			App:   strings.TrimSpace(wName[:beforeSep]),
			Title: strings.TrimSpace(wName[afterSep:]),
		}, true
	}
	return nil, false
}

func sepDefault(wName string) (wi *Winfo) {
	// App Name Last
	beforeSep := strings.LastIndex(wName, defaultWindowTitleSeparator)
	afterSep := beforeSep + len(defaultWindowTitleSeparator)
	return &Winfo{
		App:   strings.TrimSpace(wName[afterSep:]),
		Title: strings.TrimSpace(wName[:beforeSep]),
	}
}
