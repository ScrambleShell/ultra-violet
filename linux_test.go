package ultraViolet

import (
	"testing"
)

func Test_collectWindows(t *testing.T) {

	var outWmctrl = []string{
		`0x01e00002  0 desktop XdndCollectionWindowImp
0x01e00005  0 desktop unity-launcher
0x01e00008  0 desktop unity-panel
0x01e0000b  0 desktop unity-dash
0x01e0000c  0 desktop Hud
0x0160000a -1 desktop Desktop
0x0340000a  0 desktop Terminal
0x03400232  0 desktop Terminal`,
		`x b a f`,
		``,
		`0x01eeeeee 0 desktop foo`,
		`0x02000000 0 desktop bar`,
	}

	var expectedWindows = [][]*Window{
		[]*Window{
			&Window{ID: 54525962, Desktop: 0, Name: "Terminal"},
			&Window{ID: 54526514, Desktop: 0, Name: "Terminal"},
		},
		nil,
		make([]*Window, 0, 128),
		make([]*Window, 0, 128),
		[]*Window{&Window{ID: 33554432, Desktop: 0, Name: "bar"}},
	}

	var isExpectedWindowsError = []bool{
		false,
		true,
		false,
		false,
		false,
	}

	for i, out := range outWmctrl {
		windows, err := _collectWindows(out)
		for j, w := range windows {
			if w == nil || expectedWindows[i][j] == nil {
				if w == expectedWindows[i][j] {
					t.Error("window: nil")
				}
				continue
			} else if !(w.ID == expectedWindows[i][j].ID && w.Desktop == expectedWindows[i][j].Desktop && w.Name == expectedWindows[i][j].Name) {
				t.Errorf("len: %d, window: %v, case: %d, row: %d", len(windows), w, i, j)
			}
		}
		if (err == nil) == isExpectedWindowsError[i] {
			t.Errorf("error, %d", i)
		}
	}
}

func Test_findCurrentDesktop(t *testing.T) {

	var outDesktop = []string{
		`0  * DG: 3840x2160  VP: 1920,1080  WA: 65,24 1855x1056  N/A`,
		`1 * 2`,
		`0 ** 9 0 9 `,
		`j   i`,
		` 1 i dk k
	3 * jl`,
	}

	var expectedCurrentDesktop = []int{
		0,
		1,
		0,
		0,
		3,
	}

	var isExpectedDesktopError = []bool{
		false,
		false,
		true,
		true,
		false,
	}

	for i, out := range outDesktop {
		cd, err := _findCurrentDesktop(out)
		if !(cd == expectedCurrentDesktop[i] && (err != nil) == isExpectedDesktopError[i]) {
			t.Error("cuurent desktop")
		}
	}
}
