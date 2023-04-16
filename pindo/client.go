package pindo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"net/http"

	"github.com/rugwirobaker/hermes/build"
)

var (
	userAgent = build.Info().ServiceName + "/" + build.Info().Version
	baseURL   = "https://api.pindo.io/v1"
)

type Client struct {
	userAgent string
	baseURL   string
	apiKey    string
	client    *http.Client
}

// SetBaseURL sets an alternative base URL for the API.
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// SetUserAgent sets an alternative user agent for the API.
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

func New(apiKey string) *Client {
	return NewWithClient(apiKey, http.DefaultClient)
}

func NewWithClient(apiKey string, client *http.Client) *Client {
	return &Client{
		apiKey:    apiKey,
		client:    client,
		baseURL:   baseURL,
		userAgent: userAgent,
	}
}

func (c *Client) Do(ctx context.Context, method, endpoint string, in, out interface{}, headers map[string][]string) error {
	req, err := c.NewRequest(ctx, method, endpoint, in, headers)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	if resp.StatusCode > 299 {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response body: %w", err)
		}
		return handleAPIError(resp.StatusCode, buf.Bytes())
	}

	if out != nil {
		err = json.NewDecoder(resp.Body).Decode(out)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) NewRequest(ctx context.Context, method, path string, in interface{}, headers map[string][]string) (*http.Request, error) {
	var (
		body io.Reader
		url  = c.baseURL + path
	)

	if headers == nil {
		headers = make(map[string][]string)
	}

	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return nil, err
		}
		headers["Content-Type"] = []string{"application/json"}
		body = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)

	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header[k] = v
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	return req, nil
}

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("code: %d: %s", e.Status, e.Message)
}

func handleAPIError(statusCode int, responseBody []byte) error {
	switch statusCode / 100 {
	case 1, 3:
		return fmt.Errorf("API returned unexpected status, %d", statusCode)
	case 4:
		var apiErr APIError
		err := json.Unmarshal(responseBody, &apiErr)
		if err != nil {
			return fmt.Errorf("API returned unexpected status, %d", statusCode)
		}
		return &apiErr
	case 5:
		return fmt.Errorf("API returned unexpected status, %d", statusCode)
	default:
		return errors.New("something went terribly wrong")
	}
}
