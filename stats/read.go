package stats

import "io"

// AggregateLogs processes incoming data from the reader and writes appropriate
// data with respect to format required.
func AggregateLogs(r io.ReadCloser, w io.WriteCloser, format string) error {
	defer w.Close()
	defer r.Close()

	return nil
}
