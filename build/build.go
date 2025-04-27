package build

import (
	"encoding/base64"
	"runtime"
	"runtime/debug"
	"strings"
)

var base64SourceDiff string

const infoUnknownMsg = "unknown"

// GoVersion used to build the app.
func GoVersion() string {
	return runtime.Version()
}

// CommitHash used to build the app.
func CommitHash() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return infoUnknownMsg
	}
	for _, setting := range info.Settings {
		if setting.Key != "vcs.revision" {
			continue
		}
		return setting.Value
	}
	return infoUnknownMsg
}

// CommitDate of the hash used to build the app.
func CommitDate() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return infoUnknownMsg
	}
	for _, setting := range info.Settings {
		if setting.Key != "vcs.time" {
			continue
		}
		return setting.Value
	}
	return infoUnknownMsg
}

// SourceModified since the last commit.
func SourceModified() bool {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return false
	}
	for _, setting := range info.Settings {
		if setting.Key != "vcs.modified" {
			continue
		}
		return strings.Contains(strings.ToLower(setting.Value), "true")
	}
	return false
}

// SourceDiff the current state against the last commit.
func SourceDiff() string {
	data, err := base64.StdEncoding.DecodeString(base64SourceDiff)
	if err != nil {
		return infoUnknownMsg
	}
	return string(data)
}
