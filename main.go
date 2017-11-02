package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/status-im/msgstat/stats"
)

var (
	version    = "0.0.1" // rely on linker -ldflags -X main.version"
	gitCommit  = ""      // rely on linker: -ldflags -X main.gitCommit"
	buildStamp = ""      // rely on linker: -ldflags -X main.buildStamp"
)

var (
	sourceFile = flag.String("file", "", "input file containing status-go logs")
	targetFile = flag.String("out", "", "output file to save processed data into")
	format     = flag.String("format", "json", "format of output data. -format=json|yaml|csv")
	getVersion = flag.Bool("v", false, "Print getVersion")
)

func main() {
	flag.Usage = printUsage
	flag.Parse()

	// if we are to print getVersion.
	if *getVersion {
		printVersion()
		return
	}

	if *sourceFile == "" && !hasIncomingData() {
		printUsage()
		return
	}

	rc := io.ReadCloser(os.Stdin)
	wc := io.WriteCloser(os.Stdout)

	if *targetFile != "" {
		targetOutputFile, err := os.Create(*targetFile)
		if err != nil {
			log.Fatalf("Failed to open output file %+q: %v", *targetFile, err)
			return
		}

		wc = targetOutputFile
	}

	if *sourceFile != "" {
		targetInputFile, err := os.OpenFile(*sourceFile, os.O_RDONLY, 0700)
		if err != nil {
			log.Fatalf("Failed to open input file %+q: %v", *sourceFile, err)
			return
		}

		rc = targetInputFile
	}

	if err := stats.ReadAggregates(rc, wc, *format); err != nil {
		log.Fatalf("Failed to process logs successfully: %v", err)
		return
	}

	log.Println("Aggregation Complete!")

	if *targetFile != "" {
		log.Printf("See output in %+q\n", *targetFile)
	}
}

// printVersion prints corresponding build getVersion with associated build stamp and git commit if provided.
func printVersion() {
	var vers []string
	vers = append(vers, version)

	if buildStamp != "" {
		vers = append(vers, fmt.Sprintf("build#%s", buildStamp))
	}

	if gitCommit != "" {
		vers = append(vers, fmt.Sprintf("git#%s", gitCommit))
	}

	fmt.Fprint(os.Stdout, strings.Join(vers, " "))
}

// printUsage prints out usage message for CLI tool.
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
		log.Printf("Error: unable to retrieve `stat` for os.Stdin: %+q", err)
		return false
	}

	return stat.Size() > 0
}
