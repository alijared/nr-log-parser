package search

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"

	"github.com/alijared/nr-log-parser/internal/context"
	"github.com/alijared/nr-log-parser/internal/progressbar"
	"github.com/alijared/nr-log-parser/internal/scanner"
	"github.com/alijared/nr-log-parser/pkg/errors"
)

const (
	LOG_FILE_OUTPUT = "nrlp.log"
)

var (
	filename,
	level,
	component,
	customSearch,
	bf,
	af,
	outputFile string
	before,
	after time.Time
	hasTimeFilter bool
	attrs         [][]byte
)

var searchCMD = &cobra.Command{
	Use:     "search",
	Short:   "Search log file for matching attributes",
	PreRunE: validateFlags,
	RunE:    search,
	Example: searchExamples(),
}

func SearchCMD() *cobra.Command {
	searchCMD.Flags().StringVarP(&filename, "file", "f", "", "log filepath to search")
	searchCMD.Flags().StringVarP(&level, "level", "l", "", "match lines with log level")
	searchCMD.Flags().StringVarP(
		&component,
		"component",
		"c",
		"",
		"match lines with component",
	)
	searchCMD.Flags().StringVarP(
		&customSearch,
		"custom",
		"q",
		"",
		"match lines with custom substring",
	)
	searchCMD.Flags().StringVarP(
		&outputFile,
		"output",
		"o",
		LOG_FILE_OUTPUT,
		"filtered log filepath",
	)
	searchCMD.Flags().StringVarP(&bf, "before", "b", "", "match lines before time")
	searchCMD.Flags().StringVarP(&af, "after", "a", "", "match lines after time")
	_ = searchCMD.MarkFlagRequired("file")

	return searchCMD
}

func search(_ *cobra.Command, _ []string) error {
	f, err := os.Open(filename)
	if err != nil {
		return errors.NewExecutionError("unable to open log file: %s", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Unable to close log file: %s", err)
		}
	}()

	scanCount := 0
	matchCount := 0
	var buffer []byte
	if err := progressbar.Wrap(f, func(bar *pb.ProgressBar) error {
		scnr := scanner.New(f)
		for scnr.Scan() {
			line := scnr.Bytes()
			line = append(line, '\n')
			if inLine(line, attrs...) {
				buffer = append(buffer, line...)
				matchCount++
			}
			bar.Add(len(line))
			scanCount++
		}
		if err := scnr.Err(); err != nil {
			return errors.NewExecutionError("error scanning file: %s", err)
		}
		return nil
	}); err != nil {
		return err
	}

	fmt.Printf("Scanned %d lines, matched %d lines\n", scanCount, matchCount)
	if buffer != nil {
		if err := ioutil.WriteFile(outputFile, buffer, 0644); err != nil {
			return errors.NewExecutionError("unable to write output to new file: %s", err)
		}
		fmt.Printf("Matched lines written to %s\n", outputFile)
	}
	return nil
}

func inLine(line []byte, attrs ...[]byte) bool {
	for _, attr := range attrs {
		if bytes.Index(line, attr) == -1 {
			return false
		}
	}
	if hasTimeFilter {
		t, err := scanner.ParseTime(line)
		if err != nil {
			log.Printf("unable to parse log time: %s", err)
			return false
		}
		if !before.IsZero() {
			if before.Before(t) || before.Equal(t) {
				return false
			}
		}
		if !after.IsZero() {
			if after.After(t) || after.Equal(t) {
				return false
			}
		}
	}
	return true
}

func validateFlags(_ *cobra.Command, _ []string) error {
	if level != "" {
		attrs = append(attrs, searchAttribute("level", level))
	}
	if component != "" {
		attrs = append(attrs, searchAttribute("component", component))
	}
	if customSearch != "" {
		split := strings.Split(customSearch, ",")
		for _, attr := range split {
			attrs = append(attrs, []byte(attr))
		}
	}
	if err := validateTimeFilter(); err != nil {
		return errors.NewValidationError("invalid time filter: %s", err)
	}
	if attrs == nil && !hasTimeFilter {
		return errors.NewValidationError("you must add at least a level, component or custom search")
	}

	return nil
}

func validateTimeFilter() error {
	if af != "" {
		if err := setTime(&after, af); err != nil {
			return err
		}
	}
	if bf != "" {
		if err := setTime(&before, bf); err != nil {
			return err
		}
	}
	return nil
}

func setTime(t *time.Time, s string) error {
	ct, err := time.Parse(context.DateFormat(), s)
	if err != nil {
		return err
	}
	*t = ct.UTC()
	hasTimeFilter = true
	return nil
}

func searchAttribute(prefix, value string) []byte {
	return []byte(fmt.Sprintf("%s=%s", prefix, value))
}

func searchExamples() string {
	examples := []string{
		"nrlp search -f mylog.log --level error",
		"nrlp search -f mylog.log --custom key=value",
		"nrlp search -f mylog.log --before \"2019-10-18 15:05:00\" --after \"2019-10-18 15:00:00\"",
	}
	s := ""
	for _, e := range examples {
		s += fmt.Sprintf("  %s\n", e)
	}
	return s[:len(s)-1]
}
