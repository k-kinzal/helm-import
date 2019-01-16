package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"net/url"
	"os"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "import",
		Long: "Import Helm Chart where only public code into local repository.",
		Example: `
  $ helm import https://example.com/path/to/file.tgz
  $ helm import https://github.com/user/repo/
  $ helm import https://github.com/user/repo/tree/branch
  $ helm import https://github.com/user/repo/tree/branch/path/to`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u, err := url.Parse(args[0])
			if err != nil {
				return err
			}
			return Import(*u)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	cmd.SetUsageTemplate(`Usage:{{if .Runnable}}
  helm import URL{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}`)

	cmd.AddCommand()

	return cmd
}

func main() {
	cmd := newRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
