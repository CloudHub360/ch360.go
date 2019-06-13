package commands

import (
	"context"
	"fmt"
	"github.com/CloudHub360/ch360.go/ch360"
	"io"
)

//go:generate mockery -name "ModuleGetter"

const ListModulesCommandString = "list modules"

type ModuleGetter interface {
	GetAll(ctx context.Context) (ch360.ModuleList, error)
}

type ListModules struct {
	client ModuleGetter
	writer io.Writer
}

func NewListModules(client ModuleGetter, out io.Writer) *ListModules {
	return &ListModules{
		client: client,
		writer: out,
	}
}

func (cmd *ListModules) Execute(ctx context.Context) error {
	Modules, err := cmd.client.GetAll(ctx)
	if err != nil {
		fmt.Fprintln(cmd.writer, "[FAILED]")
		return err
	}

	if len(Modules) == 0 {
		fmt.Fprintln(cmd.writer, "No modules found.")
	}

	for _, Module := range Modules {
		fmt.Fprintln(cmd.writer, Module.ID)
	}

	return nil
}

func (cmd ListModules) Usage() string {
	return ListModulesCommandString
}
