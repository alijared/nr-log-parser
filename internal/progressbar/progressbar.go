package progressbar

import (
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"

	"github.com/alijared/nr-log-parser/pkg/errors"
)

func New(f *os.File) (*pb.ProgressBar, error) {
	filesize, err := fileSize(f)
	if err != nil {
		return nil, err
	}

	bar := pb.New64(filesize)
	bar.SetRefreshRate(time.Millisecond)
	return bar.Start(), nil
}

func Wrap(f *os.File, cb func(*pb.ProgressBar) error) error {
	bar, err := New(f)
	if err != nil {
		return err
	}
	defer bar.Finish()
	return cb(bar)
}

func fileSize(f *os.File) (int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, errors.NewExecutionError("unable to get log file info: %s", err)
	}
	return fi.Size(), nil
}
