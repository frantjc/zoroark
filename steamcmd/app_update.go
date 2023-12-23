package steamcmd

import (
	"context"
	"fmt"
	"io"
)

type AppUpdateCommand struct {
	AppID        string
	Beta         string
	BetaPassword string
	Validate     bool
}

func (*AppUpdateCommand) Check(flags *PromptFlags) error {
	if !flags.LoggedIn {
		return fmt.Errorf("app_update before login")
	}

	return nil
}

func (c *AppUpdateCommand) Args() ([]string, error) {
	if c == nil || c.AppID == "" {
		return nil, fmt.Errorf("app_update requires app ID")
	}

	args := []string{"app_update", c.AppID}

	if c.Beta != "" {
		args = append(args, "-beta", c.Beta)
	}

	if c.BetaPassword != "" {
		args = append(args, "-betapassword", c.BetaPassword)
	}

	if c.Validate {
		args = append(args, "validate")
	}

	return args, nil
}

func (c *AppUpdateCommand) ReadOutput(ctx context.Context, r io.Reader) error {
	return new(baseCommand).ReadOutput(ctx, r)
}

func (c *AppUpdateCommand) Modify(flags *PromptFlags) error {
	return new(baseCommand).Modify(flags)
}
