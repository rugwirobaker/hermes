package helmes

import "fmt"

// Build information
type Build struct {
	Version string `json:"version,omitempty"`
	Date    string `json:"date,omitempty"`
}

var version, buildDate string = "unset", "unset"

// String version of Build
func String() string {
	return fmt.Sprintf("Build Details:\n\tVersion:\t%s\n\tDate:\t\t%s", version, buildDate)
}

// Data returns build details as a struct
func Data() *Build {
	return &Build{
		Version: version,
		Date:    buildDate,
	}
}
