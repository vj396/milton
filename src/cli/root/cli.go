package root

import (
	"github.com/spf13/cobra"
)

var (
	root = &cobra.Command{
		Use:   "milton",
		Short: "Milton is a slack bot",
		Long:  "A robust interrupt management and workflow track for slack",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func GetRoot() *cobra.Command {
	return root
}
