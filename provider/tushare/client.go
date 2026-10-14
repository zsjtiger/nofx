// Package tushare reads A-share market data from the Tushare Pro HTTP API.
// Ts codes use the same form as NOFX symbols ("600519.SH"), so no mapping is needed.
package tushare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultBaseURL is the public Tushare Pro endpoint.
const DefaultBaseURL = "https://api.tushare.pro"

// Client calls the Tushare Pro API with a user token.
type Client struct {
	token   string
	baseURL string
	http    *http.Client
}

// NewClient returns a client for the public endpoint.
func NewClient(token string) *Client {
	return &Client{
		token:   token,
		baseURL: DefaultBaseURL,
		http:    &http.Client{Timeout: 15 * time.Second},
	}
}

// Table is the tabular payload Tushare returns: named columns and row arrays.
type Table struct {
	Fields []string
	Items  [][]any
}

// Rows converts the table into one map per row, keyed by column name.
func (t *Table) Rows() []map[string]any {
	rows := make([]map[string]any, 0, len(t.Items))
	for _, item := range t.Items {
		row := make(map[string]any, len(t.Fields))
		for i, f := range t.Fields {
			if i < len(item) {
				row[f] = item[i]
			}
		}
		rows = append(rows, row)
	}
	return rows
}

type apiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *struct {
		Fields []string `json:"fields"`
		Items  [][]any  `json:"items"`
	} `json:"data"`
}

// Call invokes one Tushare API by name and returns its table.
func (c *Client) Call(ctx context.Context, apiName string, params map[string]string, fields string) (*Table, error) {
	if c.token == "" {
		return nil, fmt.Errorf("tushare token is not configured")
	}
	body, err := json.Marshal(map[string]any{
		"api_name": apiName,
		"token":    c.token,
		"params":   params,
		"fields":   fields,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tushare %s: %w", apiName, err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tushare %s: HTTP %d", apiName, resp.StatusCode)
	}

	var parsed apiResponse
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, fmt.Errorf("tushare %s: invalid response: %w", apiName, err)
	}
	if parsed.Code != 0 {
		return nil, fmt.Errorf("tushare %s: code %d: %s", apiName, parsed.Code, parsed.Msg)
	}
	if parsed.Data == nil {
		return &Table{}, nil
	}
	return &Table{Fields: parsed.Data.Fields, Items: parsed.Data.Items}, nil
}
