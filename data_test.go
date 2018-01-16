package ultraViolet

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestStreamPrint(t *testing.T) {
	s0 := Snapshot{
		Time: time.Date(2017, time.December, 31, 15, 0, 0, 0, time.UTC),
		Windows: []*Window{
			&Window{ID: 0, Desktop: 0, Name: "foo - bar"},
		},
		Active:  0,
		Visible: []int{0},
	}
	tests := []struct {
		st       Stream
		expected string
	}{
		{
			Stream{
				Snapshots: []*Snapshot{
					&s0,
					&s0,
					&s0,
				},
			},
			"Sun Dec 31 15:00:00 +0000 UTC 2017\n\tActive: [bar||foo]\nSun Dec 31 15:00:00 +0000 UTC 2017\n\tActive: [bar||foo]\nSun Dec 31 15:00:00 +0000 UTC 2017\n\tActive: [bar||foo]\n",
		},
	}
	for i, tt := range tests {
		if tt.st.Print() != tt.expected {
			t.Errorf("case%d\nexpected:\n`%v`\nactual:\n`%v`", i, tt.expected, tt.st.Print())
		}
	}

}

func TestSnapshotPrint(t *testing.T) {
	tests := []struct {
		ss       Snapshot
		expected string
	}{
		{
			Snapshot{
				Time:    time.Time{},
				Windows: nil,
				Active:  0,
				Visible: []int{0},
			},
			"Mon Jan 1 00:00:00 +0000 UTC 0001\n",
		},

		{
			Snapshot{
				Time: time.Date(2017, time.December, 31, 15, 0, 0, 0, time.UTC),
				Windows: []*Window{
					&Window{ID: 0, Desktop: 0, Name: "foo - bar"},
				},
				Active:  0,
				Visible: []int{0},
			},
			"Sun Dec 31 15:00:00 +0000 UTC 2017\n\tActive: [bar||foo]\n",
		},
		{
			Snapshot{
				Time:    time.Date(2017, time.December, 31, 15, 0, 0, 0, time.UTC),
				Windows: nil,
				Active:  0,
				Visible: []int{0},
			},
			"Sun Dec 31 15:00:00 +0000 UTC 2017\n",
		},
		{

			Snapshot{
				Time: time.Date(2017, time.December, 31, 15, 0, 0, 0, time.UTC),
				Windows: []*Window{
					&Window{ID: 0, Desktop: 0, Name: "foo - bar"},
				},
				Active:  0,
				Visible: nil,
			},
			"Sun Dec 31 15:00:00 +0000 UTC 2017\n\tActive: [bar||foo]\n",
		},
		{
			Snapshot{
				Time: time.Date(2017, time.December, 31, 15, 0, 0, 0, time.UTC),
				Windows: []*Window{
					&Window{ID: -1, Desktop: -1, Name: "a-1 - b-1 - c-1"},
					&Window{ID: 0, Desktop: 0, Name: "a0 - b0 - c0"},
					&Window{ID: 1, Desktop: 0, Name: "a1 - b1 - c1"},
					&Window{ID: 10, Desktop: 1, Name: "a10 - b10 - c10"},
				},
				Active:  0,
				Visible: []int{0},
			},
			"Sun Dec 31 15:00:00 +0000 UTC 2017\n\tActive: [c0||a0 - b0]\n\tOther: [c-1||a-1 - b-1], [c1||a1 - b1], [c10||a10 - b10], \n",
		},
	}

	for i, tt := range tests {
		actual := tt.ss.Print()
		if actual != tt.expected {
			t.Errorf("case%d\nexpected:\n%s\nactual:\n%s", i, tt.expected, actual)
		}
	}
}

func TestWriteWindows(t *testing.T) {
	tests := []struct {
		windows  []*Window
		name     string
		expected string
	}{
		{
			nil,
			"",
			"",
		},
		{
			make([]*Window, 0),
			"",
			"",
		},
		{
			[]*Window{&Window{ID: 0, Desktop: 0, Name: "title0 - app0"},
				&Window{ID: 1, Desktop: 1, Name: "title1 - app1"},
				&Window{ID: 10, Desktop: 10, Name: "title10 - app10"},
				&Window{ID: 11, Desktop: 11, Name: "title11 - app11"},
			},
			"case",
			"\tcase: [app0||title0], [app1||title1], [app10||title10], [app11||title11], \n",
		},
	}

	for i, tt := range tests {
		var b bytes.Buffer
		writeWindows(&b, tt.windows, tt.name)
		actual := string(b.Bytes())
		if actual != tt.expected {
			t.Errorf("case%d:\nexpected:\n%s\nactual:\n%s\n", i, tt.expected, actual)
		}
	}
}

func TestWindowIsSticky(t *testing.T) {
	tests := []struct {
		window   *Window
		isSticky bool
	}{
		{&Window{ID: 0, Desktop: 0, Name: ""}, false},
		{&Window{ID: 1, Desktop: 1, Name: "Desktop"}, false},
		{&Window{ID: -1, Desktop: -1, Name: "Google"}, true},
	}

	for i, tt := range tests {
		if tt.window.IsSticky() != tt.isSticky {
			t.Errorf("case%d", i)
		}
	}
}

func TestWindowIsOnDesktop(t *testing.T) {
	tests := []struct {
		window    *Window
		onDesktop bool
	}{
		{&Window{ID: 0, Desktop: 0, Name: ""}, true},
		{&Window{ID: 1, Desktop: 1, Name: "Desktop"}, false},
		{&Window{ID: -1, Desktop: -1, Name: "Google"}, true},
	}
	for i, tt := range tests {
		if tt.window.IsOnDesktop(0) != tt.onDesktop {
			t.Errorf("case:%d", i)
		}
	}
}

type testTracker0 int

func (_ testTracker0) Snap() (*Snapshot, error) { return new(Snapshot), nil }
func (_ testTracker0) Deps() string             { return "" }

type testTracker1 string

func (_ testTracker1) Snap() (*Snapshot, error) { return new(Snapshot), errors.New("error!") }
func (_ testTracker1) Deps() string             { return "deps" }

type testTracker2 struct{}

func (_ testTracker2) Snap() (*Snapshot, error) {
	s := &Snapshot{
		Time:    time.Now(),
		Windows: nil,
		Active:  0,
		Visible: []int{0, 1, 2, 3},
	}
	return s, nil
}
func (_ testTracker2) Deps() string { return "depends on Go" }

func TestRegisterTracker(t *testing.T) {
	trackers = make(map[string]func() Tracker)
	tests := []struct {
		f          func() Tracker
		tracker    Tracker
		name       string
		duplicated bool
	}{
		// first cases
		{
			func() Tracker { return testTracker0(3) },
			testTracker0(3),
			"case0",
			false,
		},
		{
			func() Tracker { return testTracker1("case1") },
			testTracker1("case1"),
			"case1",
			false,
		},
		{
			func() Tracker { return testTracker2{} },
			testTracker2{},
			"duplicate case",
			false,
		},
		// duplicated cases
		{
			func() Tracker { return testTracker2{} },
			testTracker2{},
			"duplicate case",
			true,
		},
	}

	for i, tt := range tests {
		err := RegisterTracker(tt.name, tt.f)
		if (err == nil) == tt.duplicated {
			t.Errorf("case%d", i)
		}
		if !tt.duplicated && !(reflect.DeepEqual(trackers[tt.name](), tt.tracker)) {
			t.Errorf("case%d", i)
		}
	}
}

func TestNewTracker(t *testing.T) {
	trackers = make(map[string]func() Tracker)
	testsNotError := []struct {
		f       func() Tracker
		tracker Tracker
		name    string
	}{
		// cases are unique
		{
			func() Tracker { return testTracker0(3) },
			testTracker0(3),
			"case0",
		},
		{
			func() Tracker { return testTracker1("case1") },
			testTracker1("case1"),
			"case1",
		},
		{
			func() Tracker { return testTracker2{} },
			testTracker2{},
			"duplicate case",
		},
	}
	testsError := []struct {
		name string
	}{
		{"foo"},
	}
	// define trackers
	for _, tt := range testsNotError {
		trackers[tt.name] = tt.f
	}

	// test
	for i, tt := range testsNotError {
		tracker, err := NewTracker(tt.name)
		if err != nil {
			t.Errorf("case%d", i)
		}
		if !reflect.DeepEqual(tracker, tt.tracker) {
			t.Errorf("case%d", i)
		}
	}
	for i, tt := range testsError {
		_, err := NewTracker(tt.name)
		if err == nil {
			t.Errorf("case%d", i)
		}
	}
}

var testWindowsSlice = [][]*Window{
	{
		&Window{
			ID:      0,
			Desktop: 0,
			Name:    "",
		},
		&Window{
			ID:      1,
			Desktop: 1,
			Name:    "a",
		},
		&Window{
			ID:      -1,
			Desktop: -1,
			Name:    "foo",
		},
	},
	{nil},
}

var testSt = Stream{
	Snapshots: []*Snapshot{
		&Snapshot{},
		&Snapshot{
			Time:    time.Now(),
			Windows: []*Window{&Window{}},
			Active:  2,
			Visible: []int{0, 1, 2},
		},
		&Snapshot{
			Time:    time.Now(),
			Windows: nil,
			Active:  0,
			Visible: nil,
		},
		&Snapshot{
			Time: time.Now(),
			Windows: []*Window{
				&Window{ID: 0, Desktop: 0, Name: ""},
				&Window{ID: 1, Desktop: 1, Name: "Desktop"},
				&Window{ID: -1, Desktop: -1, Name: "Google"},
			},
			Active:  -1,
			Visible: []int{-1, 0, 1},
		},
	},
}

/*
	Winfo Test
*/
type typeExpectedWinfoCase struct {
	winfo     *Winfo
	isFulfill bool
}

var caseWinfoWindows = []*Window{
	&Window{Name: ""},
	&Window{Name: "Default"},
	&Window{Name: "Title - Example - Google Chrome"},
	&Window{Name: "Slack - Channel bar"},
	&Window{Name: "foo - bar"},
}

var expectedWinfo = []*Winfo{
	&Winfo{},
	&Winfo{
		Title: "Default",
	},
	&Winfo{
		App:    "Google Chrome",
		SubApp: "Example",
		Title:  "Title",
	},
	&Winfo{
		App:   "Slack",
		Title: "Channel bar",
	},
	&Winfo{
		App:   "bar",
		Title: "foo",
	},
}

func TestWindowInfo(t *testing.T) {
	for i, w := range caseWinfoWindows {
		if !reflect.DeepEqual(w.Info(), expectedWinfo[i]) {
			t.Errorf("case%d", i)
		}
	}
}

var expectedChrome = []typeExpectedWinfoCase{
	{winfo: nil, isFulfill: false},
	{winfo: nil, isFulfill: false},
	{
		winfo: &Winfo{
			App:    "Google Chrome",
			SubApp: "Example",
			Title:  "Title",
		},
		isFulfill: true,
	},
	{winfo: nil, isFulfill: false},
	{winfo: nil, isFulfill: false},
}

func TestChromeInfo(t *testing.T) {
	for i, w := range caseWinfoWindows {
		wi, isChrome := chromeInfo(w.Name)
		if isChrome != expectedChrome[i].isFulfill {
			t.Errorf("case%d", i)
		} else {
			if !reflect.DeepEqual(wi, expectedChrome[i].winfo) {
				t.Errorf("case%d", i)
			}
		}
	}
}

var expectedSlack = []typeExpectedWinfoCase{
	{winfo: nil, isFulfill: false},
	{winfo: nil, isFulfill: false},
	{winfo: nil, isFulfill: false},
	{
		winfo: &Winfo{
			App:   "Slack",
			Title: "Channel bar",
		},
		isFulfill: true,
	},
	{winfo: nil, isFulfill: false},
}

func TestSlackInfo(t *testing.T) {
	for i, w := range caseWinfoWindows {
		n := strings.Index(w.Name, defaultWindowTitleSeparator)
		if n > -1 && n < len(w.Name) {
			wi, isSlack := slackInfo(w.Name, n)
			if isSlack != expectedSlack[i].isFulfill {
				t.Errorf("case%d", i)
			} else {
				if !reflect.DeepEqual(wi, expectedSlack[i].winfo) {
					t.Errorf("case%d", i)
				}
			}
		} else {
			if expectedSlack[i].isFulfill {
				t.Errorf("case%d", i)
			}
		}
	}
}

var expectedSepDefault = []*Winfo{
	nil,
	nil,
	&Winfo{
		App:   "Google Chrome",
		Title: "Title - Example",
	},
	&Winfo{
		App:   "Channel bar",
		Title: "Slack",
	},
	&Winfo{
		App:   "bar",
		Title: "foo",
	},
}

func TestSepDefault(t *testing.T) {
	for i, w := range caseWinfoWindows {
		if sep := strings.Index(w.Name, defaultWindowTitleSeparator); sep > -1 && sep < len(w.Name) {
			if !reflect.DeepEqual(sepDefault(w.Name), expectedSepDefault[i]) {
				t.Errorf("case%d", i)
			}
		}
	}
}
