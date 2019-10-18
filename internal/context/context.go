package context

import (
	"github.com/spf13/cobra"
)

var (
	cmd        *cobra.Command
	dateFormat string
)

func CMD() *cobra.Command {
	return cmd
}

func SetCMD(c *cobra.Command) {
	cmd = c
}

func DateFormat() string {
	return dateFormat
}

func SetDateFormat(fmt string) {
	dateFormat = fmt
}
