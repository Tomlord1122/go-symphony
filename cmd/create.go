/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Tomlord1122/go-symphony/cmd/flags"
	"github.com/Tomlord1122/go-symphony/cmd/program"
	"github.com/Tomlord1122/go-symphony/cmd/steps"
	"github.com/Tomlord1122/go-symphony/cmd/ui/multiInput"
	"github.com/Tomlord1122/go-symphony/cmd/ui/multiSelect"
	"github.com/Tomlord1122/go-symphony/cmd/ui/spinner"
	"github.com/Tomlord1122/go-symphony/cmd/ui/textinput"
	"github.com/Tomlord1122/go-symphony/cmd/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const logo = `

   _____                       __                     
  / ___/__  ______ ___  ____  / /_  ____  ____  __  __
  \__ \/ / / / __ \__ \/ __ \/ __ \/ __ \/ __ \/ / / /
 ___/ / /_/ / / / / / / /_/ / / / / /_/ / / / / /_/ / 
/____/\__, /_/ /_/ /_/ .___/_/ /_/\____/_/ /_/\__, /  
     /____/         /_/                      /____/                

`

// These are the styles for the logo, tip message, and ending message.
// Using lipgloss to style the text.
// TODO: Modify to go-recipe style
var (
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	tipMsgStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#B74545")).Italic(true)
	tipCommand     = lipgloss.NewStyle().Foreground(lipgloss.Color("#F3F300")).Italic(true)
	endingMsgStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D9D9D9")).Bold(true)
)

func init() {
	var flagFramework flags.Framework
	var flagDBDriver flags.Database
	var advancedFeatures flags.AdvancedFeatures
	var flagGit flags.Git
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringP("name", "n", "", "Name of project to create")
	createCmd.Flags().VarP(&flagFramework, "framework", "f", fmt.Sprintf("Framework to use. Allowed values: %s", strings.Join(flags.AllowedProjectTypes, ", ")))
	createCmd.Flags().VarP(&flagDBDriver, "driver", "d", fmt.Sprintf("Database drivers to use. Allowed values: %s", strings.Join(flags.AllowedDBDrivers, ", ")))
	createCmd.Flags().BoolP("advanced", "a", false, "Get prompts for advanced features")
	createCmd.Flags().Var(&advancedFeatures, "feature", fmt.Sprintf("Advanced feature to use. Allowed values: %s", strings.Join(flags.AllowedAdvancedFeatures, ", ")))
	createCmd.Flags().VarP(&flagGit, "git", "g", fmt.Sprintf("Git to use. Allowed values: %s", strings.Join(flags.AllowedGitsOptions, ", ")))

	utils.RegisterStaticCompletions(createCmd, "framework", flags.AllowedProjectTypes)
	utils.RegisterStaticCompletions(createCmd, "driver", flags.AllowedDBDrivers)
	utils.RegisterStaticCompletions(createCmd, "feature", flags.AllowedAdvancedFeatures)
	utils.RegisterStaticCompletions(createCmd, "git", flags.AllowedGitsOptions)
}

type Options struct {
	ProjectName *textinput.Output
	ProjectType *multiInput.Selection
	DBDriver    *multiInput.Selection
	Advanced    *multiSelect.Selection
	Workflow    *multiInput.Selection
	Git         *multiInput.Selection
}

// createCmd defines the "create" command for the CLI
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Go project and don't worry about the structure",
	Long:  "Go Symphony is a CLI tool that allows you to focus on the actual Go code, and not the project structure. Perfect for someone new to the Go language",

	Run: func(cmd *cobra.Command, args []string) {
		var tprogram *tea.Program
		var err error

		isInteractive := false
		flagName := cmd.Flag("name").Value.String()

		if flagName != "" && !utils.ValidateModuleName(flagName) {
			err = fmt.Errorf("'%s' is not a valid module name. Please choose a different name", flagName)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		rootDirName := utils.GetRootDir(flagName)
		if rootDirName != "" && doesDirectoryExistAndIsNotEmpty(rootDirName) {
			err = fmt.Errorf("directory '%s' already exists and is not empty. Please choose a different name", rootDirName)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		// VarP already validates the contents of the framework flag.
		// If this flag is filled, it is always valid
		flagFramework := flags.Framework(cmd.Flag("framework").Value.String())
		flagDBDriver := flags.Database(cmd.Flag("driver").Value.String())
		flagGit := flags.Git(cmd.Flag("git").Value.String())

		options := Options{
			ProjectName: &textinput.Output{},
			ProjectType: &multiInput.Selection{},
			DBDriver:    &multiInput.Selection{},
			Advanced: &multiSelect.Selection{
				Choices: make(map[string]bool),
			},
			Git: &multiInput.Selection{},
		}

		project := &program.Project{
			ProjectName:     flagName,
			ProjectType:     flagFramework,
			DBDriver:        flagDBDriver,
			FrameworkMap:    make(map[flags.Framework]program.Framework),
			DBDriverMap:     make(map[flags.Database]program.Driver),
			AdvancedOptions: make(map[string]bool),
			GitOptions:      flagGit,
		}

		steps := steps.InitSteps(flagFramework, flagDBDriver)
		fmt.Printf("%s\n", logoStyle.Render(logo))

		// Advanced option steps:
		flagAdvanced, err := cmd.Flags().GetBool("advanced")
		if err != nil {
			log.Fatal("failed to retrieve advanced flag")
		}

		if flagAdvanced {
			fmt.Println(tipMsgStyle.Render("Advanced mode\n"))
		}

		if project.ProjectName == "" {
			isInteractive = true
			tprogram := tea.NewProgram(textinput.InitialTextInputModel(options.ProjectName, "What is the name of your project?", project))
			if _, err := tprogram.Run(); err != nil {
				log.Printf("Name of project contains an error: %v", err)
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}

			if options.ProjectName.Output != "" && !utils.ValidateModuleName(options.ProjectName.Output) {
				err = fmt.Errorf("'%s' is not a valid module name. Please choose a different name", options.ProjectName.Output)
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}

			rootDirName = utils.GetRootDir(options.ProjectName.Output)
			if doesDirectoryExistAndIsNotEmpty(rootDirName) {
				err = fmt.Errorf("directory '%s' already exists and is not empty. Please choose a different name", rootDirName)
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}
			project.ExitCLI(tprogram)

			project.ProjectName = options.ProjectName.Output
			err := cmd.Flag("name").Value.Set(project.ProjectName)
			if err != nil {
				log.Fatal("failed to set the name flag value", err)
			}
		}

		if project.ProjectType == "" {
			isInteractive = true
			step := steps.Steps["framework"]
			tprogram = tea.NewProgram(multiInput.InitialModelMulti(step.Options, options.ProjectType, step.Headers, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}
			project.ExitCLI(tprogram)

			step.Field = options.ProjectType.Choice

			// this type casting is always safe since the user interface can
			// only pass strings that can be cast to a flags.Framework instance
			project.ProjectType = flags.Framework(strings.ToLower(options.ProjectType.Choice))
			err := cmd.Flag("framework").Value.Set(project.ProjectType.String())
			if err != nil {
				log.Fatal("failed to set the framework flag value", err)
			}
		}

		if project.DBDriver == "" {
			isInteractive = true
			step := steps.Steps["driver"]
			tprogram = tea.NewProgram(multiInput.InitialModelMulti(step.Options, options.DBDriver, step.Headers, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}
			project.ExitCLI(tprogram)

			// this type casting is always safe since the user interface can
			// only pass strings that can be cast to a flags.Database instance
			project.DBDriver = flags.Database(strings.ToLower(options.DBDriver.Choice))
			err := cmd.Flag("driver").Value.Set(project.DBDriver.String())
			if err != nil {
				log.Fatal("failed to set the driver flag value", err)
			}
		}

		if flagAdvanced {

			featureFlags := cmd.Flag("feature").Value.String()

			if featureFlags != "" {
				featuresFlagValues := strings.Split(featureFlags, ",")
				for _, key := range featuresFlagValues {
					project.AdvancedOptions[key] = true
				}
			} else {
				isInteractive = true
				step := steps.Steps["advanced"]
				tprogram = tea.NewProgram((multiSelect.InitialModelMultiSelect(step.Options, options.Advanced, step.Headers, project)))
				if _, err := tprogram.Run(); err != nil {
					cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
				}
				project.ExitCLI(tprogram)
				for key, opt := range options.Advanced.Choices {
					project.AdvancedOptions[strings.ToLower(key)] = opt
					err := cmd.Flag("feature").Value.Set(strings.ToLower(key))
					if err != nil {
						log.Fatal("failed to set the feature flag value", err)
					}
				}
				if err != nil {
					log.Fatal("failed to set the htmx option", err)
				}
			}

		}

		if project.GitOptions == "" {
			isInteractive = true
			step := steps.Steps["git"]
			tprogram = tea.NewProgram(multiInput.InitialModelMulti(step.Options, options.Git, step.Headers, project))
			if _, err := tprogram.Run(); err != nil {
				cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
			}
			project.ExitCLI(tprogram)

			project.GitOptions = flags.Git(strings.ToLower(options.Git.Choice))
			err := cmd.Flag("git").Value.Set(project.GitOptions.String())
			if err != nil {
				log.Fatal("failed to set the git flag value", err)
			}
		}

		currentWorkingDir, err := os.Getwd()
		if err != nil {
			log.Printf("could not get current working directory: %v", err)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}
		project.AbsolutePath = currentWorkingDir

		spinner := tea.NewProgram(spinner.InitialModelNew())

		// add synchronization to wait for spinner to finish
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := spinner.Run(); err != nil {
				cobra.CheckErr(err)
			}
		}()

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("The program encountered an unexpected issue and had to exit. The error was:", r)
				fmt.Println("If you continue to experience this issue, please post a message on our GitHub page or join our Discord server for support.")
				if releaseErr := spinner.ReleaseTerminal(); releaseErr != nil {
					log.Printf("Problem releasing terminal: %v", releaseErr)
				}
			}
		}()

		err = project.CreateMainFile()

		releaseErr := spinner.ReleaseTerminal()
		if releaseErr != nil {
			log.Printf("Problem releasing terminal: %v", releaseErr)
		}

		if err != nil {
			log.Printf("Problem creating files for project.")
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		} else {
			fmt.Println(endingMsgStyle.Render("Project created successfully!"))
			fmt.Println(tipCommand.Render(fmt.Sprintf("cd into the newly created project with: `cd %s`\n", utils.GetRootDir(project.ProjectName))))

			if isInteractive {
				nonInteractiveCommand := utils.NonInteractiveCommand(cmd.Use, cmd.Flags())
				fmt.Println(tipMsgStyle.Render("Tip: Repeat the equivalent Symphony with the following non-interactive command:"))
				fmt.Println(tipMsgStyle.Italic(false).Render(fmt.Sprintf("• %s\n", nonInteractiveCommand)))
			}
		}
	},
}

// doesDirectoryExistAndIsNotEmpty checks if the directory exists and is not empty
func doesDirectoryExistAndIsNotEmpty(name string) bool {
	if _, err := os.Stat(name); err == nil {
		dirEntries, err := os.ReadDir(name)
		if err != nil {
			log.Printf("could not read directory: %v", err)
			cobra.CheckErr(textinput.CreateErrorInputModel(err))
		}
		if len(dirEntries) > 0 {
			return true
		}
	}
	return false
}
