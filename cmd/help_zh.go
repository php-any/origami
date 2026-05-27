package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.SetHelpCommand(newHelpCommand())
}

func newHelpCommand() *cobra.Command {
	bin := rootCmd.Name()
	return &cobra.Command{
		Use:   "help [command]",
		Short: "查看命令帮助",
		Long: fmt.Sprintf(`查看 %s 各命令的帮助信息。

使用 "%s [command] --help" 查看子命令详细说明。`, bin, bin),
		ValidArgsFunction: func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			var completions []string
			cmd, _, e := c.Root().Find(args)
			if e != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			if cmd == nil {
				cmd = c.Root()
			}
			for _, subCmd := range cmd.Commands() {
				if subCmd.IsAvailableCommand() {
					if strings.HasPrefix(subCmd.Name(), toComplete) {
						completions = append(completions, subCmd.Name())
					}
				}
			}
			return completions, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				_ = c.Root().Usage()
				return
			}
			if cmd.Context() == nil {
				cmd.SetContext(c.Context())
			}
			_ = cmd.Help()
		},
	}
}
