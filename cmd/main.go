package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nick96/pixelmatch"
)

var (
	output    = flag.String("output", "", "File to output the diff to")
	threshold = flag.Float64("threshold", 0.1, "Sensitivity of diff [0, 1]")
	includeAA = flag.Bool("include-anti-aliasing", false, "Do anti-aliasing detection")
)

func abort(format string, vals ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, vals...)
	os.Exit(1)
}

func main() {
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...] <expected> <actual>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintf(os.Stderr, "Actual and expected image files must both be provided")
		flag.Usage()
		os.Exit(1)
	}

	expected, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		abort("Failed to read %s: %v", flag.Arg(0), err)
	}
	actual, err := ioutil.ReadFile(flag.Arg(1))
	if err != nil {
		abort("Failed to read %s: %v", flag.Arg(1), err)
	}

	diff, count, err := pixelmatch.PixelMatch(
		expected,
		actual,
		pixelmatch.Threshold(float32(*threshold)),
		pixelmatch.AntiAliasDetection(*includeAA),
	)
	if err != nil {
		abort("Failed to compare %s and %s: %v", flag.Arg(0), flag.Arg(1), err)
	}

	if *output == "" {
		fmt.Printf("There are %d pixels different between %s and %s\n", count, flag.Arg(0), flag.Arg(1))
	}

	if err := ioutil.WriteFile(*output, diff, 0666); err != nil {
		fmt.Printf("Failed to write diff to %s: %v\n", *output, err)
	}
}
