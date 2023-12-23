package command

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/frantjc/zoroark"
	"github.com/frantjc/zoroark/x11"
	"github.com/spf13/cobra"
)

// NewShim returns the root command for
// zoroark's shim which acts as its CLI entrypoint.
// It runs `xauth add` to give access to the docker
// host's x11 server and then runs a subcommand with
// the appropriate DISPLAY environment variable.
func NewShim() *cobra.Command {
	var (
		verbosity int
		display   = &x11.Display{
			Hostname: "host.docker.internal",
		}
		xauthorityEntry = &x11.XAuthorityEntry{
			Display: display,
		}
		cmd = &cobra.Command{
			Use:           "shim cookie",
			Version:       zoroark.GetSemver(),
			SilenceErrors: true,
			SilenceUsage:  true,
			Args:          cobra.MinimumNArgs(1),
			PreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(
					zoroark.WithLogger(cmd.Context(), zoroark.NewLogger().V(2-verbosity)),
				)
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx   = cmd.Context()
					xauth = new(x11.XAuth)
					sub   = exec.CommandContext(ctx, args[0], args[0:]...)
				)
				sub.Env = append(os.Environ(), fmt.Sprintf("%s=%s", x11.EnvVarDisplay, display.String()))
				sub.Stdout = cmd.OutOrStdout()
				sub.Stderr = cmd.ErrOrStderr()

				if err := xauth.Add(ctx, xauthorityEntry); err != nil {
					return err
				}

				return sub.Run()
			},
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.Flags().CountVarP(&verbosity, "verbose", "V", "verbosity for zoroark")

	cmd.Flags().IntVar(&display.Display, "display", 0, "x11 display")
	cmd.Flags().IntVar(&display.Screen, "screen", 0, "x11 screen")

	cmd.Flags().StringVar(&xauthorityEntry.Proto, "proto", ".", "x11 protocol")
	cmd.Flags().StringVar(&xauthorityEntry.Cookie, "cookie", "", "x11 cookie")
	_ = cmd.MarkFlagRequired("cookie")

	return cmd
}
