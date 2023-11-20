package fly

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/rugwirobaker/hermes"
)

const DefaultTimeout = 2 * time.Second

// check if this implements the Environment interface
var _ hermes.Environment = (*Environment)(nil)

type Environment struct {
	HTTPClient *http.Client

	Timeout time.Duration
}

func NewEnvironment() *Environment {
	return &Environment{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", "/.fly/api")
				},
			},
		},
		Timeout: DefaultTimeout,
	}
}

func (e *Environment) Type() string {
	return "fly"
}

func (e *Environment) GetNodeRole(ctx context.Context) (string, error) {
	appName := AppName()
	if appName == "" {
		return "", fmt.Errorf("cannot get primary status on host environment: %w", fmt.Errorf("app name unavailable"))
	}

	machineID := MachineID()
	if machineID == "" {
		return "", fmt.Errorf("cannot get primary status on host environment: %w", fmt.Errorf("machine id unavailable"))
	}

	u := url.URL{
		Scheme: "http",
		Host:   "localhost",
		Path:   path.Join("/v1", "apps", appName, "machines", machineID, "metadata"),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("cannot get primary status on host environment: %w", err)
	}

	req.Header.Set("Accept", "text/plain")

	resp, err := e.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("cannot get primary status on host environment: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cannot get primary status on host environment: %w", fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	var metadata = make(metadataResponse)

	//unmarshal the response body into the metadata map
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return "", fmt.Errorf("cannot get primary status on host environment: %w", err)
	}

	//check if the metadata has the role key and return the value
	if role, ok := metadata["role"]; ok {
		return role, nil
	}
	return "", fmt.Errorf("cannot get primary status on host environment: %w", fmt.Errorf("role metadata unavailable"))
}

// a map of all the key/value pairs in the metadata
type metadataResponse map[string]string

// Available returns true if currently running in a Fly.io environment.
func Available() bool { return AppName() != "" }

// AppName returns the name of the current Fly.io application.
func AppName() string {
	return os.Getenv("FLY_APP_NAME")
}

// MachineID returns the identifier for the current Fly.io machine.
func MachineID() string {
	return os.Getenv("FLY_MACHINE_ID")
}
