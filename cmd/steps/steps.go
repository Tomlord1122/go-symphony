// Package steps provides utility for creating
// each step of the CLI
package steps

import "github.com/Tomlord1122/go-symphony/cmd/flags"

// A StepSchema contains the data that is used
// for an individual step of the CLI
type StepSchema struct {
	StepName string // The name of a given step
	Options  []Item // The slice of each option for a given step
	Headers  string // The title displayed at the top of a given step
	Field    string
}

// Steps contains a slice of steps
type Steps struct {
	Steps map[string]StepSchema
}

// An Item contains the data for each option
// in a StepSchema.Options
type Item struct {
	Flag, Title, Desc string
}

// InitSteps initializes and returns the *Steps to be used in the CLI program
func InitSteps(projectType flags.Framework, databaseType flags.Database) *Steps {
	steps := &Steps{
		map[string]StepSchema{
			"framework": {
				StepName: "Go Project Framework",
				Options: []Item{
					{
						Title: "Standard-library",
						Desc:  "The built-in Go standard library HTTP package",
					},
					{
						Title: "Chi",
						Desc:  "A lightweight, idiomatic and composable router for building Go HTTP services",
					},
				},
				Headers: "What framework do you want to use in your Go project?",
				Field:   projectType.String(),
			},
			"driver": {
				StepName: "Go Project Database Driver",
				Options: []Item{
					{
						Title: "Postgres",
						Desc:  "Go postgres driver for Go's database/sql package"},
					{
						Title: "None",
						Desc:  "Choose this option if you don't wish to install a specific database driver."},
				},
				Headers: "What database driver do you want to use in your Go project?",
				Field:   databaseType.String(),
			},
			"advanced": {
				StepName: "Advanced Features",
				Headers:  "Which advanced features do you want?",
				Options: []Item{
					{
						Flag:  "GitHubAction",
						Title: "Go Project Workflow",
						Desc:  "Workflow templates for testing, cross-compiling and releasing Go projects",
					},
					{
						Flag:  "Websocket",
						Title: "Websocket endpoint",
						Desc:  "Add a websocket endpoint",
					},
					{
						Flag:  "Docker",
						Title: "Docker",
						Desc:  "Dockerfile and docker-compose generic configuration for go project",
					},
				},
			},
			"git": {
				StepName: "Git Repository",
				Headers:  "Which git option would you like to select for your project?",
				Options: []Item{
					{
						Title: "Commit",
						Desc:  "Initialize a new git repository and commit all the changes",
					},
					{
						Title: "Skip",
						Desc:  "Proceed without initializing a git repository",
					},
				},
			},
		},
	}

	return steps
}
