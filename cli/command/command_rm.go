package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rmCommand = &cobra.Command{
	Use:   "rm",
	Short: "copy",
	Long:  "copy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("rm")
	},
}

func init() {
	CommandCommand.AddCommand(rmCommand)
}
