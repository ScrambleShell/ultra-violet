package ultraViolet

import (
	"reflect"
	"sort"
	"testing"
)

func Test_snapAll(t *testing.T) {
	casesAllProcWins := []map[process][]*Window{
		nil,
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{ID: -1, Name: "1. zsh - iTerm2"},
			},
			{"Terminal", 370000}:           nil,
			{"Keynote", 420000}:            nil,
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{Name: "1. zsh - iTerm2"},
				&Window{Name: "2. bash - iTerm2"},
			},
			{"Terminal", 370000}: nil,
			{"Keynote", 420000}: []*Window{
				&Window{Name: "foo.key - Keynote"},
				&Window{Name: "bar.key - Keynote"},
			},
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
	}
	expectedAllWindows := [][]*Window{
		make([]*Window, 0, 0),
		[]*Window{
			&Window{ID: -1, Name: "1. zsh - iTerm2"},
			&Window{ID: 69000, Name: "Finder"},
			&Window{ID: 100000, Name: "Google Chrome"},
			&Window{ID: 370000, Name: "Terminal"},
			&Window{ID: 420000, Name: "Keynote"},
			&Window{ID: 640000, Name: "System Preferences"},
		},
		[]*Window{
			&Window{ID: 0, Name: "1. zsh - iTerm2"},
			&Window{ID: 0, Name: "2. bash - iTerm2"},
			&Window{ID: 0, Name: "bar.key - Keynote"},
			&Window{ID: 0, Name: "foo.key - Keynote"},
			&Window{ID: 69000, Name: "Finder"},
			&Window{ID: 100000, Name: "Google Chrome"},
			&Window{ID: 370000, Name: "Terminal"},
			&Window{ID: 640000, Name: "System Preferences"},
		},
	}
	for i, allProcWins := range casesAllProcWins {
		actualAllWindows := _snapAll(allProcWins)
		tmp := windowsSlice(actualAllWindows)
		sort.Sort(tmp)
		actualAllWindows = []*Window(tmp)
		if !reflect.DeepEqual(actualAllWindows, expectedAllWindows[i]) {
			t.Errorf("case%d: %v", i, actualAllWindows)
		}
	}
}

type windowsSlice []*Window

func (w windowsSlice) Len() int { return len(w) }

func (w windowsSlice) Less(i, j int) bool {
	if w[i].ID == w[j].ID {
		return w[i].Name < w[j].Name
	}
	return w[i].ID < w[j].ID
}

func (w windowsSlice) Swap(i, j int) {
	w[j], w[i] = w[i], w[j]
}

func Test_snapAcitive(t *testing.T) {
	casesProcWins := []map[process][]*Window{
		nil,
		make(map[process][]*Window),
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
		},
		map[process][]*Window{
			{"iTerm2", 190000}: []*Window{
				&Window{ID: -1, Name: "1. zsh - iTerm2"},
			},
		},
		map[process][]*Window{
			{"iterm2", 190000}: []*Window{
				&Window{Name: "1. zsh - iterm2"},
				&Window{Name: "2. bash - iterm2"},
			},
		},
		map[process][]*Window{
			{"", 0}:  []*Window{},
			{"", 40}: []*Window{},
		},
	}
	expectedActive := []int{0, 0, 100000, -1, 0, 0}
	expectedIsError := []bool{
		false,
		false,
		false,
		false,
		true,
		true,
	}
	for i, procWins := range casesProcWins {
		actualActive, actualErr := _snapActive(procWins)
		if (actualErr != nil) != expectedIsError[i] {
			t.Errorf("case%d: %v", i, actualErr)
		}
		if actualActive != expectedActive[i] {
			t.Errorf("case%d: %d", i, actualActive)
		}
	}
}

func Test_snapVisible(t *testing.T) {
	casesAllProcWins := []map[process][]*Window{
		nil,
		nil,
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{ID: -1, Name: "1. zsh - iTerm2"},
			},
			{"Terminal", 370000}:           nil,
			{"Keynote", 420000}:            nil,
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{ID: -1, Name: "1. zsh - iTerm2"},
			},
			{"Terminal", 370000}:           nil,
			{"Keynote", 420000}:            nil,
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{ID: -1, Name: "1. zsh - iTerm2"},
				&Window{ID: 8, Name: "2. bash - iTerm2"},
			},
			{"Terminal", 370000}: nil,
			{"Keynote", 420000}: []*Window{
				&Window{ID: 9, Name: "foo.key - Keynote"},
				&Window{ID: 23, Name: "bar.key - Keynote"},
			},
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{ID: -1, Name: "1. zsh - iTerm2"},
				&Window{ID: 8, Name: "2. bash - iTerm2"},
			},
			{"Terminal", 370000}: nil,
			{"Keynote", 420000}: []*Window{
				&Window{ID: 9, Name: "foo.key - Keynote"},
				&Window{ID: 23, Name: "bar.key - Keynote"},
			},
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
	}
	casesProcWins := []map[process][]*Window{
		nil,
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{ID: -1, Name: "1. zsh - iTerm2"},
			},
			{"Terminal", 370000}:           nil,
			{"Keynote", 420000}:            nil,
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
		nil,
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{ID: -1, Name: "1. zsh - iTerm2"},
			},
			{"Terminal", 370000}:           nil,
			{"Keynote", 420000}:            nil,
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{Name: "1. zsh - iTerm2"},
				&Window{Name: "2. bash - iTerm2"},
			},
			{"Terminal", 370000}: nil,
			{"Keynote", 420000}: []*Window{
				&Window{Name: "foo.key - Keynote"},
				&Window{Name: "bar.key - Keynote"},
			},
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{Name: "zsh - iTerm2"},
				&Window{Name: "bash - iTerm2"},
			},
			{"Terminal", 370000}: nil,
			{"Keynote", 420000}: []*Window{
				&Window{Name: ".key - Keynote"},
				&Window{Name: ".key - Keynote"},
			},
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
	}
	expectedVisible := [][]int{
		make([]int, 0, 6),
		[]int{190000},
		make([]int, 0, 6),
		[]int{-1},
		[]int{-1, 8, 9, 23},
		make([]int, 0, 6),
	}
	for i, allProcWins := range casesAllProcWins {
		actualVisible := _snapVisible(allProcWins, casesProcWins[i])
		tmp := sort.IntSlice(actualVisible)
		sort.Sort(tmp)
		actualVisible = []int(tmp)
		if !reflect.DeepEqual(actualVisible, expectedVisible[i]) {
			t.Errorf("case%d: %v", i, actualVisible)
		}

	}

}

func TestParseASOutput(t *testing.T) {
	outs := []string{
		``,
		`PROCESS 69000:Finder
PROCESS 100000:Google Chrome
PROCESS 190000:iTerm2
WINDOW 5001:1. zsh
WINDOW 42000:. nvim
PROCESS 370000:Terminal
PROCESS 420000:Keynote
PROCESS 640000:System Preferences
WINDOW -1:Panel
WINDOW 5000:セキュリティとプライバシー`,
		`PROCESS 69000:Finder
PROCESS 100000:Google Chrome
PROCESS 190000:iTerm2
WINDOW -1:1. zsh
PROCESS 370000:Terminal
PROCESS 420000:Keynote
PROCESS 640000:System Preferences`,
		`WINDOW 5000:1. zsh`,
	}

	expectedProcWins := []map[process][]*Window{
		nil,
		map[process][]*Window{
			{"System Preferences", 640000}: []*Window{
				&Window{ID: -1, Name: "Panel - System Preferences"},
				&Window{ID: 5000, Name: "セキュリティとプライバシー - System Preferences"},
			},
			{"Finder", 69000}:         nil,
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{ID: 5001, Name: "1. zsh - iTerm2"},
				&Window{ID: 42000, Name: ". nvim - iTerm2"},
			},
			{"Terminal", 370000}: nil,
			{"Keynote", 420000}:  nil,
		},
		map[process][]*Window{
			{"Google Chrome", 100000}: nil,
			{"iTerm2", 190000}: []*Window{
				&Window{ID: -1, Name: "1. zsh - iTerm2"},
			},
			{"Terminal", 370000}:           nil,
			{"Keynote", 420000}:            nil,
			{"System Preferences", 640000}: nil,
			{"Finder", 69000}:              nil,
		},
		map[process][]*Window{
			{"", 0}: []*Window{
				&Window{ID: 5000, Name: "1. zsh - "},
			},
		},
	}
	isExpectedErrors := []bool{
		true,
		false,
		false,
		false,
	}
	for i, out := range outs {
		actual, err := parseASOutput(out)
		if (err != nil) != isExpectedErrors[i] {
			t.Errorf("case%d: %v", i, err)
		} else {
			if !(reflect.DeepEqual(actual, expectedProcWins[i])) {
				t.Errorf("case%d, '%v'", i, actual)
			}
		}
	}
}

func TestParseWindowLine(t *testing.T) {
	lines := []string{
		"WINDOW :Сегодня я банан",
		"WINDOW :01234567",
		"0123456789012:00000000000",
		"12345678901:2",
	}
	procIds := []int{
		1234567,
		0,
		-1,
		0,
	}
	expectedWin := []string{
		"Сегодня я банан",
		"01234567",
		"00000000000",
		"2",
	}
	expectedWinID := []int{
		2482953042,
		277576791,
		789012,
		8901,
	}
	for i, line := range lines {
		win, winID := parseWindowLine(line, procIds[i])
		if win != expectedWin[i] || winID != expectedWinID[i] {
			t.Errorf("case: %d, win: %s, winID: %d", i, win, winID)
		}
	}
}

func TestHash(t *testing.T) {
	cases := []string{
		"Сегодня я банан",
		"",
		"0",
	}
	hashes := []int64{
		893410942,
		2166136261,
		890022063,
	}
	for i, s := range cases {
		if h := hash(s); h != hashes[i] {
			t.Errorf("error: case %d hash is %d", i, h)
		}
	}
}
