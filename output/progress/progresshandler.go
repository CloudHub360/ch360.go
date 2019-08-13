package progress

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/output/resultsWriters"
	"github.com/gosuri/uiprogress"
	"github.com/pkg/errors"
	"io"
)

type ProgressHandler struct {
	resultsWriter resultsWriters.ResultsWriter
	showProgress  bool
	progress      *uiprogress.Progress
	progressBar   *uiprogress.Bar
	out           io.Writer
	started       bool
}

//func NewProgressHandlerFor(params *config.GlobalFlags, progressOut io.Writer) (*ProgressHandler, error) {
//	resultsWriter, err := resultsWriters.NewResultsWriterFor(params)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return NewProgressHandler(resultsWriter, params.ShowProgress, progressOut), nil
//}

func NewProgressHandler(resultsWriter resultsWriters.ResultsWriter, showProgress bool, progressOut io.Writer) *ProgressHandler {
	progress := uiprogress.New()
	progress.SetOut(progressOut)

	progress.Start()

	return &ProgressHandler{
		resultsWriter: resultsWriter,
		showProgress:  showProgress,
		progress:      progress,
		out:           progressOut,
	}
}

func (c *ProgressHandler) updateProgressBar() {
	if c.showProgress {
		c.progressBar.Incr()
	}
}

func (c *ProgressHandler) Notify(filename string, result interface{}) error {
	if !c.started {
		return errors.New("NotifyStart must be called before Notify")
	}
	c.updateProgressBar()
	return c.resultsWriter.WriteResult(filename, result)
}

func (c *ProgressHandler) NotifyErr(filename string, err error) error {
	if !c.started {
		return errors.New("NotifyStart must be called before NotifyErr")
	}
	c.updateProgressBar()
	return nil
}

func (c *ProgressHandler) initProgressBar(total int) {
	c.progressBar = c.progress.AddBar(total).PrependFunc(func(bar *uiprogress.Bar) string {
		return fmt.Sprintf("Processing file [%d/%d]", bar.Current(), bar.Total)
	})
}

func (c *ProgressHandler) NotifyStart(total int) error {
	if c.showProgress {
		c.initProgressBar(total)
	}
	c.started = true
	return c.resultsWriter.Start()
}

func (c *ProgressHandler) NotifyFinish() error {
	if !c.started {
		return errors.New("NotifyStart must be called before NotifyFinish")
	}
	if c.showProgress {
		c.progress.Stop()
	}
	return c.resultsWriter.Finish()
}
