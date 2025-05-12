/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/spf13/cobra"
)

// GoSymphonyVersion will be overridden by the goreleaser in the CI run.
// It will be set to the actual version of the release in github
var GoSymphonyVersion string

func getGoSymphonyVersion() string {
	noAvailable := "No version info available"

	if len(GoSymphonyVersion) != 0 {
		return GoSymphonyVersion
	}
	// TODO: add comment
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return noAvailable
	}

	// If no main version is available, Go default it to (devel)
	if buildInfo.Main.Version != "(devel)" {
		return buildInfo.Main.Version
	}
	// vcs is abbreviation for version control system
	var vcsRevision string
	var vcsTime time.Time
	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			vcsRevision = setting.Value
		case "vcs.time":
			vcsTime, _ = time.Parse(time.RFC3339, setting.Value)
		}
	}

	if vcsRevision != "" {
		return fmt.Sprintf("%s, (%s)", vcsRevision, vcsTime)
	}
	return buildInfo.Main.Version
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of Go Symphony",
	Long:  `Show the current version of Go Symphony. It will be embedded at compile time.`,
	Run: func(cmd *cobra.Command, args []string) {
		version := getGoSymphonyVersion()
		fmt.Println("Go Symphony version: ", version)
	},
}
