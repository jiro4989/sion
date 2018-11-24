package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCommand = &cobra.Command{
	Use:   "config",
	Short: "c",
	Long:  "config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config")
	},
}

func init() {
	RootCommand.AddCommand(configCommand)
}
