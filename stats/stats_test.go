package stats_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/status-im/msgstat/stats"
)

func TestAggregationReadToJSON(t *testing.T) {
	reader, err := os.Open("../fixtures/status.log")
	if err != nil {
		t.Fatalf("Should have succesfully opened log file")
	}
	t.Logf("Should have succesfully opened log file")

	var buf bytes.Buffer
	if err := stats.ReadAggregates(reader, &wopCloser{&buf}, "json"); err != nil {
		t.Fatalf("Should have succesfully processed log data")
	}
	t.Logf("Should have succesfully processed log data")

	expected, err := ioutil.ReadFile("../fixtures/expected.json")
	if err != nil {
		t.Fatalf("Should have succesfully opened expected output file")
	}
	t.Logf("Should have succesfully opened expected output file")

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Logf("Expected: %+q\n", expected)
		t.Logf("Received: %+q\n", buf.Bytes())
		t.Fatalf("Should have succesfully matched aggregated data with expected data")
	}
	t.Logf("Should have succesfully matched aggregated data with expected data")
}

func TestAggregationReadToTOML(t *testing.T) {
	reader, err := os.Open("../fixtures/status.log")
	if err != nil {
		t.Fatalf("Should have succesfully opened log file")
	}
	t.Logf("Should have succesfully opened log file")

	var buf bytes.Buffer
	if err := stats.ReadAggregates(reader, &wopCloser{&buf}, "toml"); err != nil {
		t.Fatalf("Should have succesfully processed log data")
	}
	t.Logf("Should have succesfully processed log data")

	expected, err := ioutil.ReadFile("../fixtures/expected.toml")
	if err != nil {
		t.Fatalf("Should have succesfully opened expected output file")
	}
	t.Logf("Should have succesfully opened expected output file")

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Logf("Expected: %+q\n", expected)
		t.Logf("Received: %+q\n", buf.Bytes())
		t.Fatalf("Should have succesfully matched aggregated data with expected data")
	}
	t.Logf("Should have succesfully matched aggregated data with expected data")
}

func TestAggregationReadToYAML(t *testing.T) {
	reader, err := os.Open("../fixtures/status.log")
	if err != nil {
		t.Fatalf("Should have succesfully opened log file")
	}
	t.Logf("Should have succesfully opened log file")

	var buf bytes.Buffer
	if err := stats.ReadAggregates(reader, &wopCloser{&buf}, "yaml"); err != nil {
		t.Fatalf("Should have succesfully processed log data")
	}
	t.Logf("Should have succesfully processed log data")

	expected, err := ioutil.ReadFile("../fixtures/expected.yaml")
	if err != nil {
		t.Fatalf("Should have succesfully opened expected output file")
	}
	t.Logf("Should have succesfully opened expected output file")

	if !bytes.Equal(expected, buf.Bytes()) {
		t.Logf("Expected: %+q\n", expected)
		t.Logf("Received: %+q\n", buf.Bytes())
		t.Fatalf("Should have succesfully matched aggregated data with expected data")
	}
	t.Logf("Should have succesfully matched aggregated data with expected data")
}

type wopCloser struct {
	io.Writer
}

// Close does nothing.
func (wopCloser) Close() error {
	return nil
}
