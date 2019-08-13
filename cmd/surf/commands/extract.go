package commands

//
//import (
//	"context"
//	"errors"
//	"fmt"
//	"github.com/CloudHub360/ch360.go/ch360"
//	"github.com/CloudHub360/ch360.go/ch360/results"
//	"github.com/CloudHub360/ch360.go/config"
//	"github.com/CloudHub360/ch360.go/output/progress"
//	"github.com/CloudHub360/ch360.go/pool"
//	"github.com/mattn/go-zglob"
//	"io"
//	"os"
//)
//
//const ExtractFilesCommand = "extract"
//
////go:generate mockery -name "FileExtractor"
//
//type FileExtractor interface {
//	Extract(ctx context.Context, fileContent io.Reader, extractorName string) (*results.ExtractionResult, error)
//}
//
//type Extract struct {
//	fileExtractor   FileExtractor
//	documentGetter  ch360.DocumentGetter
//	parallelWorkers int
//	progressHandler ProgressHandler
//
//	extractorName string
//	filesPattern  string
//}
//
//func NewExtractCommand(progressHandler ProgressHandler,
//	fileExtractor FileExtractor,
//	docGetter ch360.DocumentGetter,
//	parallelism int,
//	filesPattern string,
//	extractorName string) *Extract {
//	return &Extract{
//		progressHandler: progressHandler,
//		fileExtractor:   fileExtractor,
//		documentGetter:  docGetter,
//		parallelWorkers: parallelism,
//		filesPattern:    filesPattern,
//		extractorName:   extractorName,
//	}
//}
//
//func (cmd *Extract) handlerFor(cancel context.CancelFunc, filename string, errs *[]error) pool.HandlerFunc {
//	return func(value interface{}, err error) {
//		if err != nil {
//			err = errors.New(fmt.Sprintf("Error extracting file %s: %v", filename, err))
//			cmd.progressHandler.NotifyErr(filename, err)
//			*errs = append(*errs, err)
//			// Don't process any more if there's an error
//			cancel()
//		} else {
//			if err = cmd.progressHandler.Notify(filename, value.(*results.ExtractionResult)); err != nil {
//				// An error occurred while writing output
//				*errs = append(*errs, err)
//				cancel()
//			}
//		}
//	}
//}
//
//func (cmd *Extract) Execute(ctx context.Context) error {
//	files, err := zglob.Glob(cmd.filesPattern)
//	if err != nil {
//		if os.IsNotExist(err) {
//			// The file pattern is for a specific (single) file that doesn't exist
//			return errors.New(fmt.Sprintf("File %s does not exist", cmd.filesPattern))
//		} else {
//			return err
//		}
//	}
//
//	fileCount := len(files)
//	if fileCount == 0 {
//		return errors.New(fmt.Sprintf("File glob pattern %s does not match any files. Run 'surf -h' for glob pattern examples.", cmd.filesPattern))
//	}
//
//	// Get the current number of documents, so we know how many slots are available
//	docs, err := cmd.documentGetter.GetAll(ctx)
//	if err != nil {
//		return err
//	}
//	// Limit the number of workers to the number of available doc slots
//	cmd.parallelWorkers = min(cmd.parallelWorkers, ch360.TotalDocumentSlots-len(docs))
//
//	ctx, cancel := context.WithCancel(ctx)
//
//	var (
//		processFileJobs []pool.Job
//		errs            []error
//	)
//
//	for _, filename := range files {
//		// The memory of the 'filename' var is reused here, see:
//		// https://golang.org/doc/faq#closures_and_goroutines
//		// The workaround is to copy it:
//		filename := filename // <- copy
//
//		processFileJob := pool.NewJob(
//			func() (interface{}, error) {
//				return cmd.processFile(ctx, filename, cmd.extractorName)
//			},
//			cmd.handlerFor(cancel, filename, &errs))
//
//		processFileJobs = append(processFileJobs, processFileJob)
//	}
//
//	workPool := pool.NewPool(processFileJobs, cmd.parallelWorkers)
//
//	// Print results
//	cmd.progressHandler.NotifyStart(len(processFileJobs))
//	defer cmd.progressHandler.NotifyFinish()
//	workPool.Run(ctx)
//
//	// Just return the first error.
//	if len(errs) > 0 {
//		return errs[0]
//	}
//
//	return nil
//}
//
//func (cmd *Extract) processFile(ctx context.Context, filePath string, extractorName string) (*results.ExtractionResult, error) {
//	file, err := os.Open(filePath)
//	if err != nil {
//		return nil, err
//	}
//
//	return cmd.fileExtractor.Extract(ctx, file, extractorName)
//
//}
//
//func NewExtractFilesCommandFromArgs(params *config.RunParams, client *ch360.ApiClient) (*Extract, error) {
//
//	progressHandler, err := progress.NewProgressHandlerFor(params, os.Stderr)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return NewExtractCommand(progressHandler,
//		ch360.NewFileExtractor(client.Documents, client.Documents, client.Documents),
//		client.Documents,
//		10,
//		params.FilePattern,
//		params.ExtractorName), nil
//}
//
//func (cmd Extract) Usage() string {
//	return ExtractFilesCommand
//}
