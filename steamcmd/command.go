package steamcmd

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/frantjc/go-fn"
)

type Command interface {
	Check(*PromptFlags) error
	Args() ([]string, error)
	ReadOutput(context.Context, io.Reader) error
	Modify(*PromptFlags) error
}

func (p *Prompt) Run(ctx context.Context, cmd Command) error {
	if err := cmd.Check(p.flags); err != nil {
		return err
	}

	args, err := cmd.Args()
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(p.stdin, fn.Map(args, func(arg string, _ int) any {
		return arg
	})...); err != nil {
		return err
	}

	if err := cmd.ReadOutput(ctx, p.stdout); err != nil {
		return err
	}

	return cmd.Modify(p.flags)
}

type baseCommand byte

func (*baseCommand) Check(flags *PromptFlags) error {
	return nil
}

func (*baseCommand) Args() ([]string, error) {
	return make([]string, 0), nil
}

func (*baseCommand) ReadOutput(ctx context.Context, r io.Reader) error {
	var (
		errC     = make(chan error, 1)
		buf      = new(bytes.Buffer)
		successB = []byte("Steam>")
		errB     = []byte("ERROR! ")
	)

	go func() {
		defer close(errC)

		for {
			var b [4096]byte

			if _, err := r.Read(b[:]); err != nil {
				errC <- err
				return
			}

			if _, err := buf.Write(b[:]); err != nil {
				errC <- err
				return
			}

			p := buf.Bytes()
			if _, msgB, found := bytes.Cut(p, errB); found {
				msgB, _, _ = bytes.Cut(msgB, []byte("\n"))
				errC <- fmt.Errorf(string(msgB))
				return
			} else if bytes.Contains(p, successB) {
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errC:
		return err
	}
}

func (*baseCommand) Modify(_ *PromptFlags) error {
	return nil
}
