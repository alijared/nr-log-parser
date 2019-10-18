package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/alijared/nr-log-parser/internal/cmd/search"
	"github.com/alijared/nr-log-parser/internal/cmd/stat"
	"github.com/alijared/nr-log-parser/internal/context"
	"github.com/alijared/nr-log-parser/pkg/errors"
)

const DEFAULT_DATE_FORMAT = "2006-01-02 15:04:05"

var dateFormat string

var rootCMD = &cobra.Command{
	Use:           "nrlp",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		context.SetCMD(cmd)
		context.SetDateFormat(dateFormat)
	},
}

func Execute() error {
	rootCMD.AddCommand(search.SearchCMD(), stat.StatCMD())
	rootCMD.PersistentFlags().StringVarP(
		&dateFormat,
		"df",
		"",
		DEFAULT_DATE_FORMAT,
		"date format",
	)
	if err := rootCMD.Execute(); err != nil {
		if _, ok := err.(errors.CMDError); !ok {
			cmd, _, _ := rootCMD.Traverse(os.Args[1:])
			context.SetCMD(cmd)
			return errors.NewValidationError(err.Error())
		}
		return err
	}
	return nil
}
