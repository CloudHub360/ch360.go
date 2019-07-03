package ch360

import (
	"context"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/net"
	"strings"
)

type ModulesClient struct {
	baseUrl       string
	requestSender net.HttpDoer
}

func NewModulesClient(baseUrl string, requestSender net.HttpDoer) *ModulesClient {
	return &ModulesClient{
		baseUrl:       baseUrl,
		requestSender: requestSender,
	}
}

type Module struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
	Fields      []struct {
		Name        string      `json:"name"`
		Description interface{} `json:"description"`
	} `json:"fields"`
	Parameters []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		Required    bool   `json:"required"`
	} `json:"parameters"`
}

type ModuleList []Module

func (m *ModulesClient) GetAll(ctx context.Context) (ModuleList, error) {

	response, err := newRequest(ctx, "GET", m.baseUrl+"/modules", nil).
		issue(m.requestSender)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var modulesResponse struct {
		Modules ModuleList
	}
	err = json.NewDecoder(response.Body).Decode(&modulesResponse)

	if err != nil {
		return nil, err
	}

	return modulesResponse.Modules, nil
}

// Find performs a case-insensitive search for the provided ID
// string, and returns the Module if found and nil if not.
func (l ModuleList) Find(id string) *Module {
	for _, existingModule := range l {
		if strings.ToLower(id) == strings.ToLower(existingModule.ID) {
			return &existingModule
		}
	}

	return nil
}
