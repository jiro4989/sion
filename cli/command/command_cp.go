package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cpCommand = &cobra.Command{
	Use:   "cp",
	Short: "copy",
	Long:  "copy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("copy")
	},
}

func init() {
	fmt.Println("cp init")
	CommandCommand.AddCommand(cpCommand)
}
