package log

import (
	"io"
	"os"
	"sync"

	"github.com/muesli/termenv"
)

var (
	target io.Writer = os.Stderr
	mu     sync.Mutex
)

// Writer is a writer that forward slog output to target writer.
var Writer = forwarder{}

func SetWriter(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()
	target = w
}

func ColorProfile() termenv.Profile {
	mu.Lock()
	defer mu.Unlock()
	return termenv.NewOutput(target).Profile
}

type forwarder struct{}

var _ io.Writer = (*forwarder)(nil)

func (f forwarder) Write(p []byte) (n int, err error) {
	mu.Lock()
	defer mu.Unlock()
	return target.Write(p)
}
