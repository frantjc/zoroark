package steamcmd

import (
	"errors"
	"io"
	"sync"
)

type PromptFlags struct {
	LoggedIn bool
}

type Prompt struct {
	flags  *PromptFlags
	stdin  io.Writer
	stdout io.Reader
	mu     sync.Mutex
	err    error
}

func (p *Prompt) Close() error {
	return errors.Join(p.err)
}
