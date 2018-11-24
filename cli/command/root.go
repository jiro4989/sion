package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "appName",
	Short: "short",
	Long:  "long",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root")
	},
}

func init() {
	cobra.OnInitialize()
}
