package stats_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/status-im/msgstat/stats"
)

// TestAggregationReadToJSON validates the generation of aggregation data into
// json and it's output format.
func TestAggregationReadToJSON(t *testing.T) {
	reader, err := os.Open("../fixtures/status.log")
	if err != nil {
		t.Fatalf("Should have successfully opened log file")
	}
	t.Logf("Should have successfully opened log file")

	var buf bytes.Buffer
	if err := stats.ReadAggregates(reader, &writeCloser{&buf}, "json"); err != nil {
		t.Fatalf("Should have successfully processed log data")
	}
	t.Logf("Should have successfully processed log data")

	expected, err := ioutil.ReadFile("../fixtures/expected.json")
	if err != nil {
		t.Fatalf("Should have successfully opened expected output file")
	}
	t.Logf("Should have successfully opened expected output file")

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Logf("Expected: %+q\n", expected)
		t.Logf("Received: %+q\n", buf.Bytes())
		t.Fatalf("Should have successfully matched aggregated data with expected data")
	}
	t.Logf("Should have successfully matched aggregated data with expected data")
}

// TestAggregationReadToTOML validates the generation of aggregation data into
// toml and it's output format.
func TestAggregationReadToTOML(t *testing.T) {
	reader, err := os.Open("../fixtures/status.log")
	if err != nil {
		t.Fatalf("Should have successfully opened log file")
	}
	t.Logf("Should have successfully opened log file")

	var buf bytes.Buffer
	if err := stats.ReadAggregates(reader, &writeCloser{&buf}, "toml"); err != nil {
		t.Fatalf("Should have successfully processed log data")
	}
	t.Logf("Should have successfully processed log data")

	expected, err := ioutil.ReadFile("../fixtures/expected.toml")
	if err != nil {
		t.Fatalf("Should have successfully opened expected output file")
	}
	t.Logf("Should have successfully opened expected output file")

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Logf("Expected: %+q\n", expected)
		t.Logf("Received: %+q\n", buf.Bytes())
		t.Fatalf("Should have successfully matched aggregated data with expected data")
	}
	t.Logf("Should have successfully matched aggregated data with expected data")
}

// TestAggregationReadToYAML validates the generation of aggregation data into
// yaml and it's output format.
func TestAggregationReadToYAML(t *testing.T) {
	reader, err := os.Open("../fixtures/status.log")
	if err != nil {
		t.Fatalf("Should have successfully opened log file")
	}
	t.Logf("Should have successfully opened log file")

	var buf bytes.Buffer
	if err := stats.ReadAggregates(reader, &writeCloser{&buf}, "yaml"); err != nil {
		t.Fatalf("Should have successfully processed log data")
	}
	t.Logf("Should have successfully processed log data")

	expected, err := ioutil.ReadFile("../fixtures/expected.yaml")
	if err != nil {
		t.Fatalf("Should have successfully opened expected output file")
	}
	t.Logf("Should have successfully opened expected output file")

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Logf("Expected: %+q\n", expected)
		t.Logf("Received: %+q\n", buf.Bytes())
		t.Fatalf("Should have successfully matched aggregated data with expected data")
	}
	t.Logf("Should have successfully matched aggregated data with expected data")
}

// writeCloser provides a decorator which adds a `Close` method to a composed io.Writer.
// It implements the io.WriteCloser interface.
type writeCloser struct {
	io.Writer
}

// Close does nothing.
func (writeCloser) Close() error {
	return nil
}
