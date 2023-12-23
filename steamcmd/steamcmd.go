package steamcmd

import (
	"context"
	"fmt"
	"os/exec"
	"sync"
)

func IsOnPath() bool {
	bin, err := exec.LookPath("steamcmd")
	return bin != "" && err == nil
}

type Steamcmd struct {
	Path string
}

func (s *Steamcmd) Start(ctx context.Context) (*Prompt, error) {
	cmd := exec.CommandContext(ctx, s.Path)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	p := &Prompt{&PromptFlags{}, stdin, stdout, sync.Mutex{}, nil}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		p.err = cmd.Wait()
	}()

	return p, p.Run(ctx, new(baseCommand))
}

type CommandLine struct {
	ForceInstallDir   string
	Login             *LoginCommand
	ForcePlatformType PlatformType
	AppUpdate         *AppUpdateCommand
}

func (c *CommandLine) Args() ([]string, error) {
	args := []string{}

	if c.ForceInstallDir != "" {
		args = append(args, "+force_install_dir", c.ForceInstallDir)
	}

	if c.Login == nil || c.Login.Username == "" || c.Login.Username == "anonymous" {
		args = append(args, "+login", "anonymous")
	} else if c.Login.Password != "" {
		args = append(args, "+login", c.Login.Username, c.Login.Password)

		if c.Login.SteamGuardCode != "" {
			args = append(args, "+login", c.Login.SteamGuardCode)
		}
	} else if c.Login.SteamGuardCode != "" {
		return nil, fmt.Errorf("steam guard code requires password")
	}

	if c.ForcePlatformType != "" {
		args = append(args, "+@sSteamCmdForcePlatformType", c.ForcePlatformType.String())
	}

	if c.AppUpdate != nil {
		if c.AppUpdate.AppID == "" {
			return nil, fmt.Errorf("app_update requires app ID")
		}

		args = append(args, "+app_update", c.AppUpdate.AppID)

		if c.AppUpdate.Beta != "" {
			args = append(args, "-beta", c.AppUpdate.Beta)
		} else if c.AppUpdate.BetaPassword != "" {
			return nil, fmt.Errorf("beta requires betapassword")
		}

		if c.AppUpdate.BetaPassword != "" {
			args = append(args, "-betapassword", c.AppUpdate.BetaPassword)
		}

		if c.AppUpdate.Validate {
			args = append(args, "-validate")
		}
	}

	return append(args, "+quit"), nil
}

func (s *Steamcmd) Run(ctx context.Context, c *CommandLine) error {
	if args, err := c.Args(); err != nil {
		return err
	} else {
		return exec.CommandContext(ctx, s.Path, args...).Run()
	}
}

func (s *Steamcmd) init() error {
	if s.Path == "" {
		s.Path = "steamcmd"
	}

	return nil
}
