package stat

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/alijared/nr-log-parser/internal/context"
	"github.com/alijared/nr-log-parser/internal/progressbar"
	"github.com/alijared/nr-log-parser/internal/scanner"
	"github.com/alijared/nr-log-parser/pkg/errors"
)

var filename string

var statCMD = &cobra.Command{
	Use:   "stat",
	Short: "Get stats for log file",
	RunE:  stat,
}

func StatCMD() *cobra.Command {
	statCMD.Flags().StringVarP(
		&filename,
		"file",
		"f",
		"",
		"log filepath to get stats from",
	)
	_ = statCMD.MarkFlagRequired("file")

	return statCMD
}

func stat(_ *cobra.Command, _ []string) error {
	f, err := os.Open(filename)
	if err != nil {
		return errors.NewExecutionError("unable to open log file: %s", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Unable to close log file: %s", err)
		}
	}()

	ls := &logStat{}
	if err := progressbar.Wrap(f, func(bar *pb.ProgressBar) error {
		scnr := scanner.New(f)
		for scnr.Scan() {
			line := scnr.Bytes()
			line = append(line, '\n')

			ls.addLevel(scanner.GetLogLevel(line))
			t, err := scanner.ParseTime(line)
			if err != nil {
				log.Printf("unable to parse log time: %s", err)
				bar.Add(len(line))
				continue
			}
			ls.addTime(t)
			bar.Add(len(line))
		}
		if err := scnr.Err(); err != nil {
			return errors.NewExecutionError("error scanning file: %s", err)
		}
		return nil
	}); err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Error count", "Warning Count", "Info Count", "Debug Count", "Min Time", "Max Time"})
	table.Append([]string{
		getIntStr(ls.errCount),
		getIntStr(ls.warnCount),
		getIntStr(ls.infoCount),
		getIntStr(ls.debugCount),
		getTimeStr(ls.minTime),
		getTimeStr(ls.maxTime),
	})

	fmt.Println("\nStats for log file")
	table.Render()
	return nil
}

func getIntStr(i int) string {
	return fmt.Sprintf("%d", i)
}

func getTimeStr(t time.Time) string {
	return fmt.Sprintf("%s (UTC)", t.Format(context.DateFormat()))
}
