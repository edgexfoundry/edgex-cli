// Copyright © 2019 Dell Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pager

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Writer is used to write text to a terminal using the terminal's pager (less).
// It implements the WriteCloser interface and should be closed when writing is finished.
type Writer struct {
	w io.WriteCloser
	c chan struct{}
	o sync.Once
}

// getPagerCommand inspects the PAGER environment variable for a pager command,
// otherwise returns a reasonable default ("less -FRX")
func getPagerCommand() (string, []string) {
	fromEnv := os.Getenv("PAGER")
	if fromEnv == "" {
		return "less", []string{"-FRX"}
	}

	split := strings.Split(fromEnv, " ")
	return split[0], split[1:]
}

// NewWriter returns a new pager.Writer to be used to write
// paged text to the terminal
func NewWriter() (*Writer, error) {
	name, args := getPagerCommand()
	pager := exec.Command(name, args...)

	// Create an os Pipe to allow clients to write into our pager
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	pager.Stdin = r
	pager.Stdout = os.Stdout
	pager.Stderr = os.Stderr

	c := make(chan struct{})

	writer := Writer{
		w: w,
		c: c,
		o: sync.Once{},
	}

	// run pager in goroutine
	go func() {
		defer func() {
			close(c)
			// Defer a close on the writer so potentially blocked clients will
			// become unblocked.  w has potentially already been closed by the
			// Writer's Close method, but without this Close call (or a ReadAll call,
			// alternatively) writers can be blocked forever.
			_ = writer.close()
		}()
		_ = pager.Run()
	}()

	return &writer, nil
}

// Write writes a slice of bytes to the Writer
func (w *Writer) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

// close ensures the Close() method of our pipe is called only once, regardless of whether
// the client calls Close first or the pager terminates first.
func (w *Writer) close() error {
	var err error
	w.o.Do(func() {
		err = w.w.Close()
	})

	return err
}

// Close is potentially a blocking call, as it will wait
// until the pager command terminates before returning.
func (w *Writer) Close() error {
	err := w.close()
	<-w.c
	return err
}
