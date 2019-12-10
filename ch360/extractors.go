package ch360

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/waives/surf/net"
	"io"
)

type ExtractorsClient struct {
	baseUrl       string
	requestSender net.HttpDoer
}

func NewExtractorsClient(baseUrl string, requestSender net.HttpDoer) *ExtractorsClient {
	return &ExtractorsClient{
		baseUrl:       baseUrl,
		requestSender: requestSender,
	}
}

type Extractor struct {
	Name string
}

type ExtractorList []Extractor

func (client *ExtractorsClient) Create(ctx context.Context, name string, config io.Reader) error {
	_, err := newRequest(ctx, "POST", client.baseUrl+"/extractors/"+name, config).
		issue(client.requestSender)

	return err
}

func (client *ExtractorsClient) CreateFromJson(ctx context.Context, name string, jsonTemplate io.Reader) error {
	template, err := NewModulesTemplateFromJson(jsonTemplate)

	if err != nil {
		return err
	}

	return client.CreateFromModules(ctx, name, *template)
}

type ExtractorTemplate struct {
	Modules []ModuleTemplate `json:"modules"`
}

type ModuleTemplate struct {
	ID           string                 `json:"id"`
	Arguments    map[string]interface{} `json:"arguments,omitempty"`
	FieldAliases []FieldAliasTemplate   `json:"field_aliases,omitempty"`
}

type FieldAliasTemplate struct {
	Field string `json:"field"`
	Alias string `json:"alias"`
}

func NewModulesTemplateFromJson(stream io.Reader) (*ExtractorTemplate, error) {
	template := ExtractorTemplate{}
	err := json.NewDecoder(stream).Decode(&template)

	return &template, err
}

func (client *ExtractorsClient) CreateFromModules(ctx context.Context, name string, modules ExtractorTemplate) error {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	jsonTemplate, err := json.Marshal(modules)

	if err != nil {
		return err
	}

	_, err = newRequest(ctx, "POST", client.baseUrl+"/extractors/"+name, bytes.NewBuffer(jsonTemplate)).
		withHeaders(headers).
		issue(client.requestSender)

	return err
}

func (client *ExtractorsClient) Delete(ctx context.Context, name string) error {
	_, err := newRequest(ctx, "DELETE", client.baseUrl+"/extractors/"+name, nil).
		issue(client.requestSender)

	return err
}

func (client *ExtractorsClient) GetAll(ctx context.Context) (ExtractorList, error) {
	response, err := newRequest(ctx, "GET", client.baseUrl+"/extractors", nil).
		issue(client.requestSender)

	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(response.Body)

	if err != nil {
		return nil, err
	}

	var extractorsResponse struct {
		Extractors []Extractor
	}
	err = json.Unmarshal(buf.Bytes(), &extractorsResponse)

	if err != nil {
		return nil, err
	}

	return extractorsResponse.Extractors, nil
}

func (e ExtractorList) Contains(item string) bool {
	for _, b := range e {
		if b.Name == item {
			return true
		}
	}
	return false
}
