package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/aimof/ultra-violet"
	"github.com/jessevdk/go-flags"
)

var CLI = flags.NewNamedParser("uv", flags.PrintErrors|flags.PassDoubleDash)

func init() {
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stdout)
	CLI.Usage = `
UV - automatically track which applications you use and for how long.

UV is a simple time tracker that tracks active window names and collects
statistics over active, open, and visible windows. Statistics are collected
into a local JSON file, which is used to generate a pretty HTML report.

UV is a local CLI tool and does not send any data over the network.

Example usage:

  uv dep
  uv track -o <file>
  uv show  -i <file> -w stats > viz.html

`

	if _, err := CLI.AddCommand("track", "record current windows", "Record current window metadata as JSON printed to stdout or a file. If a filename is specified and the file already exists, Thyme will append the new snapshot data to the existing data.", &trackCmd); err != nil {
		log.Fatal(err)
	}
	if _, err := CLI.AddCommand("show", "visualize data", "Generate an HTML page visualizing the data from a file written to by `uv track`.", &showCmd); err != nil {
		log.Fatal(err)
	}
	if _, err := CLI.AddCommand("dep", "dep install instructions", "Show installation instructions for required external dependencies (which vary depending on your OS and windowing system).", &depCmd); err != nil {
		log.Fatal(err)
	}
	if _, err := CLI.AddCommand("watch", "'Friend Copmputer' is watching you", "Allow 'Friend Computer' to watch your activities. Happiness is Mandatory.", &watchCmd); err != nil {
		log.Fatal(err)
	}
}

// TrackCmd is the subcommand that tracks application usage.
type TrackCmd struct {
	Out string `long:"out" short:"o" description:"output file"`
}

var trackCmd TrackCmd

func (c *TrackCmd) Execute(args []string) error {
	err := track(c.Out)
	return err
}

func track(outFile string) error {
	t, err := getTracker()
	if err != nil {
		return err
	}
	snap, err := t.Snap()
	if err != nil {
		return err
	}

	if outFile == "" {
		out, err := json.MarshalIndent(snap, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(out))
	} else {
		f, err := os.OpenFile(outFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		b, err := json.Marshal(snap)
		if err != nil {
			return err
		}

		if _, err = f.Write(b); err != nil {
			return err
		}

		if _, err = f.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}

// WatchCmd allows your 'Friend Computer' watch your activities.
type WatchCmd struct {
	Dir string `long:"dir" short:"d" description:"data and log directory"`
}

var watchCmd WatchCmd

func (c *WatchCmd) Execute(args []string) error {
	workDir, err := filepath.Abs(c.Dir)
	if err != nil {
		workDir, err = filepath.Abs(".")
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := os.Stat(workDir); err != nil {
		log.Fatalln("Dir does not exist!")
	}

	if err := os.Chdir(workDir); err != nil {
		log.Fatalln(err)
	}

	t := time.Now()
	dataDir := "./data/" + t.Format("2006/01")
	if _, err = os.Stat(dataDir); err != nil {
		if err = os.MkdirAll(dataDir, 0775); err != nil {
			log.Fatalln(err)
		}
	}
	dataFilePath := dataDir + t.Format("/02.json")

	outFilePath := workDir + "/uv.html"

	if err := track(dataFilePath); err != nil {
		return err
	}

	stream, err := readStream(dataFilePath)
	if err != nil {
		return err
	} else {
		f, err := os.OpenFile(outFilePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		if err := ultraViolet.Stats(&stream, f); err != nil {
			return err
		}
		f.Close()
	}

	return nil
}

// ShowCmd is the subcommand that reads the data emitted by the track
// subcommand and displays the data to the user.
type ShowCmd struct {
	In   string `long:"in" short:"i" description:"input file"`
	What string `long:"what" short:"w" description:"what to show {list,stats}" default:"list"`
}

var showCmd ShowCmd

func (c *ShowCmd) Execute(args []string) error {
	if c.In == "" {
		var snap ultraViolet.Snapshot
		if err := json.NewDecoder(os.Stdin).Decode(&snap); err != nil {
			return err
		}
		for _, w := range snap.Windows {
			fmt.Printf("%+v\n", w.Info())
		}
	} else {
		stream, err := readStream(c.In)
		if err != nil {
			return err
		}

		switch c.What {
		case "stats":
			if err := ultraViolet.Stats(&stream, os.Stdout); err != nil {
				return err
			}
		case "list":
			fallthrough
		default:
			ultraViolet.List(&stream)
		}
	}
	return nil
}

func readStream(in string) (ultraViolet.Stream, error) {
	var stream ultraViolet.Stream
	stream.Snapshots = make([]*ultraViolet.Snapshot, 0, 4086)
	f, err := os.Open(in)
	if err != nil {
		return stream, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var s *ultraViolet.Snapshot
		err := json.Unmarshal(scanner.Bytes(), &s)
		if err != nil {
			return stream, err
		}
		stream.Snapshots = append(stream.Snapshots, s)
	}
	return stream, nil
}

type DepCmd struct{}

var depCmd DepCmd

func (c *DepCmd) Execute(args []string) error {
	t, err := getTracker()
	if err != nil {
		return err
	}
	fmt.Println(t.Deps())
	return nil
}

func main() {
	run := func() error {
		_, err := CLI.Parse()
		if err != nil {
			if _, isFlagsErr := err.(*flags.Error); isFlagsErr {
				CLI.WriteHelp(os.Stderr)
				return nil
			} else {
				return err
			}
		}
		return nil
	}

	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func getTracker() (ultraViolet.Tracker, error) {
	switch runtime.GOOS {
	case "linux":
		return ultraViolet.NewTracker("linux")
	case "darwin":
		return ultraViolet.NewTracker("darwin")
	default:
		log.Println("Sorry, your OS is not supported.")
		os.Exit(1)
		return nil, errors.New("OS not supported.")
	}
}
