package ultraViolet

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	RegisterTracker("linux", NewLinuxTracker)
}

// LinuxTracker tracks application usage on Linux via a few standard command-line utilities.
type LinuxTracker struct{}

var _ Tracker = (*LinuxTracker)(nil)

func NewLinuxTracker() Tracker {
	return &LinuxTracker{}
}

func (t *LinuxTracker) Deps() string {
	return `
Install the following command-line utilities via your package manager of choice:
* xdpyinfo
* xwininfo
* xdotool
* wmctrl

For example:
* Debian: apt-get install x11-utils xdotool wmctrl

Note: this command prints out this message regardless of whether the dependencies are already installed.
`
}

func (t *LinuxTracker) Snap() (*Snapshot, error) {

	windows, err := collectWindows()
	if err != nil {
		return nil, err
	}
	currentDesktop, err := findCurrentDesktop()
	if err != nil {
		return nil, err
	}
	visible, err := getVisible(windows, currentDesktop)
	if err != nil {
		return nil, err
	}

	active, err := getActiveWindow()
	if err != nil {
		return nil, err
	}

	return &Snapshot{Windows: windows, Active: active, Visible: visible, Time: time.Now()}, nil
}

func getVisible(windows []*Window, currentDesktop int) ([]int, error) {
	var visible = make([]int, 0, len(windows))
	for _, window := range windows {
		cmd := exec.Command("xwininfo", "-id", fmt.Sprintf("%d", window.ID), "-stats")
		cmd.Env = append(cmd.Env, "path=/usr/bin", "DISPLAY=:0")
		out_, err := cmd.Output()
		log.Println(string(out_))
		if err != nil {
			return nil, fmt.Errorf("xwininfo failed with error: %s", err)
		}
		if window.IsOnDesktop(currentDesktop) && vis.Match(out_) {
			visible = append(visible, window.ID)
		}
	}
	return visible, nil
}

var vis = regexp.MustCompile(`Map State:\s+IsViewable`)

func collectWindows() ([]*Window, error) {
	var windows = make([]*Window, 0, 128)
	cmd := exec.Command("wmctrl", "-l")
	cmd.Env = append(cmd.Env, "path=/usr/bin", "DISPLAY=:0")
	out_, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	windows, err = _collectWindows(string(out_))
	return windows, err
}

func _collectWindows(out string) ([]*Window, error) {
	windows := make([]*Window, 0, 128)
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		id_, desktop_, name := fields[0], fields[1], strings.Join(fields[3:], " ")
		id64, err := strconv.ParseInt(id_, 0, 64)
		if err != nil {
			return nil, err
		}
		desktop, err := strconv.Atoi(desktop_)
		if err != nil {
			return nil, err
		}
		w := Window{ID: int(id64), Desktop: desktop, Name: name}
		if w.ID > 33554432 {
			windows = append(windows, &w)
		}
	}
	return windows, nil
}

func findCurrentDesktop() (int, error) {
	cmd := exec.Command("wmctrl", "-d")
	cmd.Env = append(cmd.Env, "PATH=/usr/bin", "DISPLAY=:0")
	out_, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	currentDesktop, err := _findCurrentDesktop(string(out_))
	return currentDesktop, nil
}

func _findCurrentDesktop(out string) (int, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		id_, mode := fields[0], fields[1]
		if mode == "*" {

			id, err := strconv.Atoi(id_)
			if err != nil {
				return 0, err
			}
			return id, nil
		}
	}
	return 0, errors.New("Cannot find current desktop")
}

func getActiveWindow() (int, error) {
	cmd := exec.Command("xdotool", "getactivewindow")
	cmd.Env = append(cmd.Env, "PATH=/usr/bin", "DISPLAY=:0")
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("xdotool failed with error: %s. Try running `xdotool getactivewindow` to diagnose.", err)
	}
	id, err := strconv.Atoi(strings.TrimSpace(string(out)))
	return id, err
}
