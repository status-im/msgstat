package stats

import (
	"bufio"
	"fmt"
	"io"
)

// AggregateLogs processes incoming data from the reader and writes appropriate
// data with respect to format required.
func AggregateLogs(r io.ReadCloser, w io.WriteCloser, format string) error {
	defer w.Close()
	defer r.Close()

	bufReader := bufio.NewReader(r)

	for {
		line, _, err := bufReader.ReadLine()

		// if we have reach end of file, just return.
		if err != nil {
			if err == io.EOF {
				return nil
			}

			return err
		}

		fmt.Printf("Currentline: %+q", line)
	}
	return nil
}
