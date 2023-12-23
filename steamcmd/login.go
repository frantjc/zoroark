package steamcmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

type LoginCommand struct {
	Username       string
	Password       string
	SteamGuardCode string
}

func (c *LoginCommand) Check(flags *PromptFlags) error {
	return new(baseCommand).Check(flags)
}

func (c *LoginCommand) Args() ([]string, error) {
	if c == nil || c.Username == "" || c.Username == "anonymous" {
		return []string{"login", "anonymous"}, nil
	}

	args := []string{"login", c.Username}

	if c.Password != "" {
		args = append(args, c.Password)

		if c.SteamGuardCode != "" {
			args = append(args, c.SteamGuardCode)
		}
	} else if c.SteamGuardCode != "" {
		return nil, fmt.Errorf("specified Steam guard code without password")
	}

	return args, nil
}

func (c *LoginCommand) ReadOutput(ctx context.Context, r io.Reader) error {
	var (
		errC      = make(chan error, 1)
		buf       = new(bytes.Buffer)
		successB  = []byte("Steam>")
		errB      = []byte("ERROR! ")
		couldHang = false
	)

	// The "password:" prompt doesn't get written to
	// stdout OR stderr, so we're relegated to this.
	if c.Username != "" && c.Username != "anonymous" && (c.Password == "" || c.SteamGuardCode == "") {
		couldHang = true
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*9)
		defer cancel()
	}

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
				errC <- fmt.Errorf(strings.ToLower(string(msgB)))
				return
			} else if bytes.Contains(p, successB) {
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		err := ctx.Err()

		if couldHang && errors.Is(err, context.DeadlineExceeded) {
			return fmt.Errorf("password or Steam guard code not given for nonanonymous, non-logged-in user; assumed login command hung")
		}

		return ctx.Err()
	case err := <-errC:
		return err
	}
}

func (*LoginCommand) Modify(flags *PromptFlags) error {
	if flags == nil {
		flags = &PromptFlags{}
	}

	flags.LoggedIn = true

	return nil
}
