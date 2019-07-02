package ch360

import (
	"context"
	"encoding/json"
	"github.com/CloudHub360/ch360.go/net"
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

// Map returns a map of moduleId : Module
func (l ModuleList) Map() map[string]Module {
	modulesMap := map[string]Module{}

	for _, module := range l {
		modulesMap[module.ID] = module
	}

	return modulesMap
}
