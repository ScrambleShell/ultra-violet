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
	windows0 := []*Window{
		&Window{ID: 0, Desktop: 0, Name: "title0 - app0"},
		&Window{ID: 1, Desktop: 1, Name: "title1 - app1"},
		&Window{ID: 10, Desktop: 10, Name: "title10 - app10"},
		&Window{ID: 11, Desktop: 11, Name: "title11 - app11"},
	}

	casesWindows := make([][]*Window, 3)

	casesWindows[0] = nil
	casesWindows[1] = make([]*Window, 0)
	casesWindows[2] = windows0

	casesName := []string{
		"",
		"",
		"case",
	}

	expectedStrings := []string{
		"",
		"",
		"\tcase: [app0||title0], [app1||title1], [app10||title10], [app11||title11], \n",
	}

	for i, w := range casesWindows {
		var b bytes.Buffer
		writeWindows(&b, w, casesName[i])
		var actualString string = ""
		actualString = string(b.Bytes())
		if actualString != expectedStrings[i] {
			t.Errorf("case%d: %s", i, actualString)
		}
	}
}

var testWindows = []*Window{
	{ID: 0, Desktop: 0, Name: ""},
	{ID: 1, Desktop: 1, Name: "Desktop"},
	{ID: -1, Desktop: -1, Name: "Google"},
}

var expectedIsSticky = []bool{false, false, true}

func TestWindowIsSticky(t *testing.T) {
	for i, w := range testWindows {
		if w.IsSticky() != expectedIsSticky[i] {
			t.Errorf("%d", i)
		}
	}
}

var expectedIsOnDesktop = []bool{true, false, true}

func TestWindowIsOnDesktop(t *testing.T) {
	for i, w := range testWindows {
		if w.IsOnDesktop(0) != expectedIsOnDesktop[i] {
			t.Errorf("%d", i)
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

var testTrackerFunc = []func() Tracker{
	func() Tracker { return testTracker0(3) },
	func() Tracker { return testTracker1("case1") },
	func() Tracker { return testTracker2{} },
}
var testTrackerName = []string{"case0", "case1", "case2"}

func TestRegisterTracker(t *testing.T) {
	for i, n := range testTrackerName {
		err := RegisterTracker(n, testTrackerFunc[i])
		if err != nil {
			t.Errorf("case%d", i)
		}
		if !(reflect.DeepEqual(trackers[n](), testTrackerFunc[i]())) {
			t.Errorf("case%d", i)
		}
	}
	if err := RegisterTracker(testTrackerName[0], testTrackerFunc[0]); err == nil {
		t.Errorf("case duplicate")
	}
}

func TestNewTracker(t *testing.T) {
	trackers = make(map[string]func() Tracker)
	for i, n := range testTrackerName {
		err := RegisterTracker(n, testTrackerFunc[i])
		if err != nil {
			t.Errorf("case%d", i)
		}
		tracker, err := NewTracker(n)
		if !reflect.DeepEqual(tracker, testTrackerFunc[i]()) {
			t.Errorf("case%d", i)
		}
		if err != nil {
			t.Errorf("case%d", i)
		}
	}
	if _, err := NewTracker("foo"); err == nil {
		t.Error("foo exist")
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
			Time:    time.Now(),
			Windows: testWindows,
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
