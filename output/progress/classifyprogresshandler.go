package progress

import (
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360/types"
	"github.com/CloudHub360/ch360.go/output/resultsWriters"
	"github.com/gosuri/uiprogress"
	"os"
)

type ClassifyProgressHandler struct {
	resultsWriter resultsWriters.ResultsWriter
	showProgress  bool
	progress      *uiprogress.Progress
	progressBar   *uiprogress.Bar
	out           *os.File
}

func NewClassifyProgressHandler(resultsWriter resultsWriters.ResultsWriter, showProgress bool, progressOut *os.File) *ClassifyProgressHandler {
	progress := uiprogress.New()
	progress.SetOut(progressOut)
	return &ClassifyProgressHandler{
		resultsWriter: resultsWriter,
		showProgress:  showProgress,
		progress:      progress,
		out:           progressOut,
	}
}

func (c *ClassifyProgressHandler) handleClassifyComplete() {
	if c.showProgress {
		c.progressBar.Incr()
	}
}

func (c *ClassifyProgressHandler) Notify(filename string, result *types.ClassificationResult) error {
	c.handleClassifyComplete()
	return c.resultsWriter.WriteResult(filename, result)
}

func (c *ClassifyProgressHandler) NotifyErr(filename string, err error) {
	c.handleClassifyComplete()
	fmt.Fprintln(c.out, err)
}

func (c *ClassifyProgressHandler) initProgressBar(total int) {
	c.progress.Start()
	c.progressBar = c.progress.AddBar(total).PrependFunc(func(bar *uiprogress.Bar) string {
		return fmt.Sprintf("Classifying file [%d/%d]", bar.Current(), bar.Total)
	})
}

func (c *ClassifyProgressHandler) NotifyStart(total int) error {
	if c.showProgress {
		c.initProgressBar(total)
	}
	return c.resultsWriter.Start()
}

func (c *ClassifyProgressHandler) NotifyFinish() error {
	if c.showProgress {
		c.progress.Stop()
	}
	return c.resultsWriter.Finish()
}
