/*
Copyright © 2026 Julian Easterling

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/dcjulian29/proxmoxctl/internal/color"
	"github.com/dcjulian29/proxmoxctl/internal/settings"
	"github.com/spf13/viper"
)

const KeyInsecureSkipVerify = "tls_insecure"

type Client struct {
	BaseURL    string
	APIToken   string
	Username   string
	httpClient *http.Client
}

type APIResponse struct {
	Data any `json:"data"`
}

func New() (*Client, error) {
	if err := settings.RequireConfig(); err != nil {
		return nil, err
	}

	tls := allowInsecure()

	transport := &http.Transport{
		TLSClientConfig: tls,
	}

	return &Client{
		BaseURL:  viper.GetString(settings.KeyServerURL),
		APIToken: viper.GetString(settings.KeyAPIToken),
		httpClient: &http.Client{
			Timeout:   120 * time.Second, // 2 minutes
			Transport: transport,
		},
	}, nil
}

func (c *Client) DefaultNode() (string, error) {
	nodes, err := c.Nodes()

	if err != nil {
		return "", err
	}

	if len(nodes) == 0 {
		return "", fmt.Errorf("no nodes found in cluster")
	}

	if name, ok := nodes[0]["node"].(string); ok {
		fmt.Fprintln(os.Stderr, color.Info("Using node: "+name))
		return name, nil
	}

	return "", fmt.Errorf("could not parse node name")
}

func (c *Client) Delete(path string) error {
	return c.do(http.MethodDelete, path, nil, nil)
}

func (c *Client) Get(path string, dest any) error {
	return c.do(http.MethodGet, path, nil, dest)
}

func (c *Client) Nodes() ([]map[string]any, error) {
	var resp struct {
		Data []map[string]any `json:"data"`
	}

	if err := c.Get("/nodes", &resp); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c *Client) Post(path string, body any, dest any) error {
	return c.do(http.MethodPost, path, body, dest)
}

func (c *Client) Put(path string, body any, dest any) error {
	return c.do(http.MethodPut, path, body, dest)
}

// Proxmox API tokens are formatted as "PVEAPIToken=USER@REALM!TOKENID=SECRET"
func (c *Client) authHeader() string {
	return fmt.Sprintf("PVEAPIToken=%s", c.APIToken)
}

func (c *Client) do(method, path string, body any, dest any) error {
	url := fmt.Sprintf("%s/api2/json%s", c.BaseURL, path)

	var reqBody io.Reader

	if body != nil {
		data, err := json.Marshal(body)

		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}

		reqBody = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", c.authHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("\n%w", err)
	}

	defer resp.Body.Close() //nolint:errcheck

	respBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if len(respBytes) > 0 {
			return fmt.Errorf("API error %d: %s", resp.StatusCode, string(respBytes))
		} else {
			return fmt.Errorf("API error %d", resp.StatusCode)
		}
	}

	if dest != nil {
		if err := json.Unmarshal(respBytes, dest); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

func allowInsecure() *tls.Config {
	insecure := viper.GetString(KeyInsecureSkipVerify)

	tls := &tls.Config{InsecureSkipVerify: false} //nolint:gosec

	if insecure == "true" {
		// Accept self-signed certs common on Proxmox homelab installs.
		tls.InsecureSkipVerify = true //nolint:gosec
	}

	return tls
}
