package context

import (
	"github.com/spf13/cobra"
)

var cmd *cobra.Command

func CMD() *cobra.Command {
	return cmd
}

func SetCMD(c *cobra.Command) {
	cmd = c
}
