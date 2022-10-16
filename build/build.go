package build

import "fmt"

var serviceName = "hermes"

// Build information
type Build struct {
	Date        string `json:"date,omitempty"`
	Version     string `json:"version,omitempty"`
	ServiceName string `json:"service_name"`
}

var version, buildDate string = "unset", "unset"

// String version of Build
func String() string {
	return fmt.Sprintf("Build Details:\n\tVersion:\t%s\n\tDate:\t\t%s", version, buildDate)
}

// Info returns build details as a struct
func Info() *Build {
	return &Build{
		Date:        buildDate,
		Version:     version,
		ServiceName: serviceName,
	}
}
