package clwrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Wrapper[Req any, Res any] struct {
}

func New[Req any, Res any](req Req, res Res) *Wrapper[Req, Res] {
	return &Wrapper[Req, Res]{}
}

func (w *Wrapper[Req, Res]) MakeRequest(ctx context.Context, urlPath string, req Req) (*Res, error) {
	rawData, err := json.Marshal(&req)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, urlPath, bytes.NewBuffer(rawData))
	if err != nil {
		return nil, fmt.Errorf("prepare request: %w", err)
	}

	httpResponse, err := http.DefaultClient.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", httpResponse.StatusCode)
	}

	var response Res

	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("decode request: %w", err)
	}

	return &response, nil
}
