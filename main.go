package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/status-im/msgstat/stats"
)

var (
	gitCommit    = ""      // rely on linker: -ldflags -X main.gitCommit"
	buildStamp   = ""      // rely on linker: -ldflags -X main.buildStamp"
	versionStamp = "0.0.1" // rely on linker -ldflags -X main.versionStamp"
)

var (
	sourceFile = flag.String("out", "", "output file to save processed data into")
	targetFile = flag.String("file", "", "input file containing status-go logs")
	format     = flag.String("format", "json", "format of output data. -format=json|yaml|csv")
	version    = flag.Bool("v", false, "Print version")
)

type wopCloser struct {
	io.Writer
}

// Close does nothing.
func (wopCloser) Close() error {
	return nil
}

func main() {
	flag.Usage = printUsage
	flag.Parse()

	// if we are to print version.
	if *version {
		printVersion()
		return
	}

	if *sourceFile == "" && !hasIncomingData() {
		printUsage()
		return
	}

	var wc io.WriteCloser

	if *targetFile != "" {
		targetOutputFile, err := os.Open(*targetFile)
		if err != nil {
			log.Fatalf("Failed to open output file %+q: %v", *targetFile, err)
			return
		}

		wc = targetOutputFile
	} else {
		wc = wopCloser{Writer: os.Stdout}
	}

	var rc io.ReadCloser

	if *sourceFile != "" {
		targetInputFile, err := os.OpenFile(*sourceFile, os.O_RDONLY, 0700)
		if err != nil {
			log.Fatalf("Failed to open input file %+q: %v", *sourceFile, err)
			return
		}

		rc = targetInputFile
	} else {
		rc = ioutil.NopCloser(os.Stdin)
	}

	if err := stats.AggregateLogs(rc, wc, *format); err != nil {
		log.Fatalf("Failed to process logs successfully: %v", err)
		return
	}

	log.Println("Aggregation Complete!")

	if *targetFile != "" {
		log.Printf("See output in %+q\n", *targetFile)
	}
}

func printVersion() {
	var vers []string
	vers = append(vers, versionStamp)

	if buildStamp != "" {
		vers = append(vers, fmt.Sprintf("build#%s", buildStamp))
	}

	if gitCommit != "" {
		vers = append(vers, fmt.Sprintf("git#%s", gitCommit))
	}

	fmt.Fprint(os.Stdout, strings.Join(vers, " "))
}

func printUsage() {
	fmt.Fprintf(os.Stdout, `Usage: msgstat [options]
Msgstat processes status-go lines to generate useful message delivery facts.

EXAMPLES:

  cat status.log | msgstat                                  # read from stdin
  msgstat -file=status.log                                  # read from file
  msgstat -file=status.log -out=proc.json                   # read from file and write to output file
  msgstat -format=yaml -file=status.log -out=proc.json      # read from file and write in yaml format

FLAGS:

  -v          Print version.
  -file       Path to status-go log file for processing.
  -out        Path to store processed data from logs.
  -format     Format to use to store processed data (json, yaml, csv).

`)

	// flag.PrintDefaults()
}

// hasIncomingData returns true/false if data is pending in stdin file.
func hasIncomingData() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	if stat.Size() == 0 {
		return false
	}

	return true
}
