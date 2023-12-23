package steamcmd

import (
	"context"
	"fmt"
	"io"
)

type PlatformType string

func (t PlatformType) String() string {
	return string(t)
}

var (
	PlatformTypeWindows PlatformType = "windows"
	PlatformTypeLinux   PlatformType = "linux"
	PlatformTypeMacOS   PlatformType = "macos"
)

type ForcePlatformTypeCommand PlatformType

func (c ForcePlatformTypeCommand) String() string {
	return string(c)
}

func (c ForcePlatformTypeCommand) Check(flags *PromptFlags) error {
	return nil
}

func (c ForcePlatformTypeCommand) Args() ([]string, error) {
	if c == "" {
		return nil, fmt.Errorf("empty PlatformType")
	}

	return []string{"@sSteamCmdForcePlatformType", c.String()}, nil
}

func (c ForcePlatformTypeCommand) ReadOutput(ctx context.Context, r io.Reader) error {
	return new(baseCommand).ReadOutput(ctx, r)
}

func (c ForcePlatformTypeCommand) Modify(flags *PromptFlags) error {
	return nil
}
