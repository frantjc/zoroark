package steamcmd

import (
	"context"
	"io"
)

type QuitCommand byte

func (c QuitCommand) Check(flags *PromptFlags) error {
	return new(baseCommand).Check(flags)
}

func (c *QuitCommand) Args() ([]string, error) {
	return []string{"quit"}, nil
}

func (c *QuitCommand) ReadOutput(context.Context, io.Reader) error {
	return nil
}

func (c QuitCommand) Modify(flags *PromptFlags) error {
	return new(baseCommand).Check(flags)
}
