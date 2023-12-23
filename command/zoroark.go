package command

import (
	"runtime"

	"github.com/frantjc/zoroark"
	"github.com/frantjc/zoroark/steamcmd"
	"github.com/frantjc/zoroark/x11"
	"github.com/spf13/cobra"
)

// NewZoroark returns the root command for
// zoroark which acts as its CLI entrypoint.
func NewZoroark() *cobra.Command {
	var (
		verbosity int
		cmd       = &cobra.Command{
			Use:           "zoroark",
			Version:       zoroark.GetSemver(),
			SilenceErrors: true,
			SilenceUsage:  true,
			PreRun: func(cmd *cobra.Command, args []string) {
				cmd.SetContext(
					zoroark.WithLogger(cmd.Context(), zoroark.NewLogger().V(2-verbosity)),
				)
			},
			RunE: func(cmd *cobra.Command, args []string) error {
				var (
					ctx   = cmd.Context()
					xauth = new(x11.XAuth)
				)

				_, err := xauth.List(ctx)
				if err != nil {
					return err
				}

				// if len(xauthorityEntries) == 0 {
				// 	return fmt.Errorf("found no .XAuthority cookies")
				// }

				// use xauthorityEntries[0] with &x11.Display{ Hostname: "host.docker.internal" } inside of a
				// container spawned like `docker run -e DISPLAY=host.docker.internal:0 ...`

				prompt, err := (&steamcmd.Steamcmd{Path: "/Users/fran/Library/Application Support/Steam/steamcmd"}).Start(ctx)
				if err != nil {
					return err
				}
				defer prompt.Close()

				if err := prompt.Run(ctx, steamcmd.ForceInstallDirCommand("/tmp")); err != nil {
					return err
				}

				if err := prompt.Run(ctx, steamcmd.ForcePlatformTypeCommand(steamcmd.PlatformTypeLinux)); err != nil {
					return err
				}

				appInfoPrint := steamcmd.AppInfoPrintCommand("2357570") // Overwatch 2
				// "740") // CSGO
				// "470") // Not found
				// "892970") // Valheim

				if err := prompt.Run(ctx, appInfoPrint); err != nil {
					return err
				}

				// appInfoPrint.AppInfo()

				return prompt.Run(ctx, new(steamcmd.QuitCommand))
			},
		}
	)

	cmd.SetVersionTemplate("{{ .Name }}{{ .Version }} " + runtime.Version() + "\n")
	cmd.Flags().CountVarP(&verbosity, "verbose", "V", "verbosity for zoroark")

	return cmd
}
