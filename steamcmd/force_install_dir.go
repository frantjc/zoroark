package steamcmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ForceInstallDirCommand string

func (c ForceInstallDirCommand) String() string {
	return string(c)
}

func (ForceInstallDirCommand) Check(flags *PromptFlags) error {
	if flags.LoggedIn {
		return fmt.Errorf("force_install_dir after login")
	}

	return nil
}

func (c ForceInstallDirCommand) Args() ([]string, error) {
	if c == "" {
		return nil, fmt.Errorf("empty force_install_dir")
	}

	a, err := filepath.Abs(c.String())
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(a, 0644); err != nil {
		return nil, err
	}

	return []string{"force_install_dir", a}, nil
}

func (c ForceInstallDirCommand) ReadOutput(ctx context.Context, r io.Reader) error {
	return new(baseCommand).ReadOutput(ctx, r)
}

func (c ForceInstallDirCommand) Modify(flags *PromptFlags) error {
	return new(baseCommand).Modify(flags)
}
