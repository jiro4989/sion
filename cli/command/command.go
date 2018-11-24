package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var CommandCommand = &cobra.Command{
	Use:   "command",
	Short: "cmd",
	Long:  "command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("command")
	},
}

func init() {
	RootCommand.AddCommand(CommandCommand)
}
