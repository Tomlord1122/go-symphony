package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Tomlord1122/go-symphony/v2/cmd/flags"
	"github.com/Tomlord1122/go-symphony/v2/cmd/program"
	"github.com/Tomlord1122/go-symphony/v2/cmd/steps"
	"github.com/Tomlord1122/go-symphony/v2/cmd/template/nextjs"
	"github.com/Tomlord1122/go-symphony/v2/cmd/template/sveltekit"
	"github.com/Tomlord1122/go-symphony/v2/cmd/ui/multiSelection"
	"github.com/Tomlord1122/go-symphony/v2/cmd/ui/singleSelection"
	"github.com/Tomlord1122/go-symphony/v2/cmd/ui/spinner"
	"github.com/Tomlord1122/go-symphony/v2/cmd/ui/textinput"
	"github.com/Tomlord1122/go-symphony/v2/cmd/utils"
	"github.com/Tomlord1122/go-symphony/v2/internal/scaffold"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// ASNI SHADOW
const logo = `
███████╗██╗   ██╗███╗   ███╗██████╗ ██╗  ██╗ ██████╗ ███╗   ██╗██╗   ██╗
██╔════╝╚██╗ ██╔╝████╗ ████║██╔══██╗██║  ██║██╔═══██╗████╗  ██║╚██╗ ██╔╝
███████╗ ╚████╔╝ ██╔████╔██║██████╔╝███████║██║   ██║██╔██╗ ██║ ╚████╔╝ 
╚════██║  ╚██╔╝  ██║╚██╔╝██║██╔═══╝ ██╔══██║██║   ██║██║╚██╗██║  ╚██╔╝  
███████║   ██║   ██║ ╚═╝ ██║██║     ██║  ██║╚██████╔╝██║ ╚████║   ██║   
╚══════╝   ╚═╝   ╚═╝     ╚═╝╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   
`

var (
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#CE9178")).Bold(true)                  // main text
	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#F99D77")).Italic(true) // yellow highlight
	headerStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#F25717")).Bold(true)                  // yellow title
	successStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#34D399")).Bold(true)                  // green success message
	warningStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#F99D77"))                             // yellow warning
	secondaryStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))                             // light gray secondary text
)

func init() {
	var flagDBDriver flags.Database
	var advancedFeatures flags.AdvancedFeatures
	var flagGit flags.Git
	var flagSupabaseMode flags.SupabaseMode
	var flagFrontendFramework flags.FrontendFramework
	var flagSvelteKitTemplate flags.SvelteKitTemplate
	var flagSvelteKitTypes flags.SvelteKitTypes
	var flagSvelteKitPackageManager flags.SvelteKitPackageManager

	createCmd.Flags().StringP("name", "n", "", "Name of project to create")
	createCmd.Flags().VarP(&flagDBDriver, "driver", "d", fmt.Sprintf("Database drivers to use. Allowed values: %s", strings.Join(flags.AllowedDBDrivers, ", ")))
	createCmd.Flags().BoolP("advanced", "a", false, "Get prompts for advanced features")
	createCmd.Flags().Bool("dry-run", false, "Print the planned scaffold steps without creating files")
	createCmd.Flags().Bool("no-interactive", false, "Disable interactive prompts and require all needed flags")
	createCmd.Flags().Bool("skip-install", false, "Skip dependency installation and formatting commands")
	createCmd.Flags().String("output", string(scaffold.OutputText), "Output format: text or json")
	createCmd.Flags().Var(&advancedFeatures, "feature", fmt.Sprintf("Advanced feature to use. Allowed values: %s", strings.Join(flags.AllowedAdvancedFeatures, ", ")))
	createCmd.Flags().VarP(&flagGit, "git", "g", fmt.Sprintf("Git to use. Allowed values: %s", strings.Join(flags.AllowedGitsOptions, ", ")))
	createCmd.Flags().Var(&flagSupabaseMode, "supabase-mode", fmt.Sprintf("Supabase mode when using Supabase. Allowed values: %s", strings.Join(flags.AllowedSupabaseModes, ", ")))
	createCmd.Flags().Var(&flagFrontendFramework, "frontend", fmt.Sprintf("Frontend framework to use. Allowed values: %s", strings.Join(flags.AllowedFrontendFrameworks, ", ")))
	createCmd.Flags().Var(&flagSvelteKitTemplate, "sveltekit-template", fmt.Sprintf("SvelteKit template when using SvelteKit. Allowed values: %s", strings.Join(flags.AllowedSvelteKitTemplates, ", ")))
	createCmd.Flags().Var(&flagSvelteKitTypes, "sveltekit-types", fmt.Sprintf("SvelteKit type checking when using SvelteKit. Allowed values: %s", strings.Join(flags.AllowedSvelteKitTypes, ", ")))
	createCmd.Flags().Var(&flagSvelteKitPackageManager, "sveltekit-package-manager", fmt.Sprintf("SvelteKit package manager when using SvelteKit. Allowed values: %s", strings.Join(flags.AllowedSvelteKitPackageManagers, ", ")))

	utils.RegisterStaticCompletions(createCmd, "driver", flags.AllowedDBDrivers)
	utils.RegisterStaticCompletions(createCmd, "feature", flags.AllowedAdvancedFeatures)
	utils.RegisterStaticCompletions(createCmd, "git", flags.AllowedGitsOptions)
	utils.RegisterStaticCompletions(createCmd, "supabase-mode", flags.AllowedSupabaseModes)
	utils.RegisterStaticCompletions(createCmd, "frontend", flags.AllowedFrontendFrameworks)
	utils.RegisterStaticCompletions(createCmd, "output", []string{string(scaffold.OutputText), string(scaffold.OutputJSON)})
	utils.RegisterStaticCompletions(createCmd, "sveltekit-template", flags.AllowedSvelteKitTemplates)
	utils.RegisterStaticCompletions(createCmd, "sveltekit-types", flags.AllowedSvelteKitTypes)
	utils.RegisterStaticCompletions(createCmd, "sveltekit-package-manager", flags.AllowedSvelteKitPackageManagers)
}

type Options struct {
	ProjectName *textinput.Output
	DBDriver    *singleSelection.Selection
	Advanced    *multiSelection.Selection
	Frontend    *singleSelection.Selection
	Workflow    *singleSelection.Selection
	Git         *singleSelection.Selection
}

// createCmd defines the "create" command for the CLI
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Go project with Gin framework and modern architecture",
	Long:  "Symphony is a CLI tool for creating modern Go projects with Gin framework, PostgreSQL database, and SQLC for type-safe queries. Focus on building features, not project structure.",

	Run: func(cmd *cobra.Command, args []string) {
		var tprogram *tea.Program
		var err error

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

		flagDBDriver := flags.Database(cmd.Flag("driver").Value.String())
		flagGit := flags.Git(cmd.Flag("git").Value.String())
		flagSupabaseMode := flags.SupabaseMode(cmd.Flag("supabase-mode").Value.String())
		flagFrontendFramework := flags.FrontendFramework(cmd.Flag("frontend").Value.String())
		flagSvelteKitTemplate := flags.SvelteKitTemplate(cmd.Flag("sveltekit-template").Value.String())
		flagSvelteKitTypes := flags.SvelteKitTypes(cmd.Flag("sveltekit-types").Value.String())
		flagSvelteKitPackageManager := flags.SvelteKitPackageManager(cmd.Flag("sveltekit-package-manager").Value.String())
		flagDryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			log.Fatal("failed to retrieve dry-run flag")
		}
		flagNoInteractive, err := cmd.Flags().GetBool("no-interactive")
		if err != nil {
			log.Fatal("failed to retrieve no-interactive flag")
		}
		flagSkipInstall, err := cmd.Flags().GetBool("skip-install")
		if err != nil {
			log.Fatal("failed to retrieve skip-install flag")
		}
		flagOutput := scaffold.OutputFormat(cmd.Flag("output").Value.String())

		featureFlags := cmd.Flag("feature").Value.String()
		featureValues := []string{}
		if featureFlags != "" {
			featureValues = strings.Split(featureFlags, ",")
		}

		spec := scaffold.BuildSpec(
			flagName,
			flagDBDriver,
			featureValues,
			flagGit,
			flagSupabaseMode,
			flagFrontendFramework,
			flagSvelteKitTemplate,
			flagSvelteKitTypes,
			flagSvelteKitPackageManager,
			scaffold.ExecutionOptions{
				DryRun:        flagDryRun,
				NoInteractive: flagNoInteractive,
				SkipInstall:   flagSkipInstall,
				Output:        flagOutput,
			},
		)

		if err := scaffold.ValidateSpec(spec); err != nil {
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		options := Options{
			ProjectName: &textinput.Output{},
			DBDriver:    &singleSelection.Selection{},
			Advanced: &multiSelection.Selection{
				Choices: make(map[string]bool),
			},
			Frontend: &singleSelection.Selection{},
			Git:      &singleSelection.Selection{},
		}

		project := &program.Project{
			ProjectName:     flagName,
			ProjectType:     flags.Gin, // Always use Gin framework
			DBDriver:        flagDBDriver,
			FrameworkMap:    make(map[flags.Framework]program.Framework),
			DBDriverMap:     make(map[flags.Database]program.Driver),
			AdvancedOptions: make(map[string]bool),
			GitOptions:      flagGit,
			SupabaseMode:    flagSupabaseMode,
		}

		steps := steps.InitSteps(flags.Gin, flagDBDriver)
		if flagOutput != scaffold.OutputJSON {
			fmt.Printf("%s\n", logoStyle.Render(logo))
		}

		// Advanced option steps:
		flagAdvanced, err := cmd.Flags().GetBool("advanced")
		if err != nil {
			log.Fatal("failed to retrieve advanced flag")
		}

		if flagAdvanced {
			fmt.Println(headerStyle.Render("*** Advanced Mode Enabled ***\n\n"))
		}

		collectProjectName(cmd, project, options, flagNoInteractive)

		// Skip framework selection - always use Gin
		// project.ProjectType is already set to flags.Gin above

		collectDatabaseDriver(cmd, project, options, steps, flagNoInteractive)

		featureFlags = cmd.Flag("feature").Value.String()
		if featureFlags != "" {
			featuresFlagValues := strings.Split(featureFlags, ",")
			for _, key := range featuresFlagValues {
				normalized := strings.ToLower(strings.TrimSpace(key))
				if normalized != "" {
					project.AdvancedOptions[normalized] = true
				}
			}
		} else if flagAdvanced {
			if flagNoInteractive {
				cobra.CheckErr(textinput.CreateErrorInputModel(fmt.Errorf("advanced features must be provided with --feature when --no-interactive is set")).Err())
			}
			step := steps.Steps["advanced"]
			tprogram = tea.NewProgram((multiSelection.InitialModelMultiSelect(step.Options, options.Advanced, step.Headers, project)))
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
		}

		collectFrontendFramework(cmd, project, options, steps, flagFrontendFramework, flagNoInteractive)

		collectGitMode(cmd, project, options, steps, flagNoInteractive)

		currentWorkingDir, err := os.Getwd()
		if err != nil {
			log.Printf("could not get current working directory: %v", err)
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}
		project.AbsolutePath = currentWorkingDir
		spec = buildScaffoldSpec(
			cmd,
			project,
			options.Advanced.Choices,
			scaffold.ExecutionOptions{
				DryRun:        flagDryRun,
				NoInteractive: flagNoInteractive,
				SkipInstall:   flagSkipInstall,
				Output:        flagOutput,
			},
		)

		if err := scaffold.ValidateSpec(spec); err != nil {
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		if err := maybeHandleDryRun(spec, currentWorkingDir); err != nil {
			cobra.CheckErr(err)
		}
		if spec.Execution.DryRun {
			return
		}

		if spec.Execution.SkipInstall {
			project.SkipInstall = true
		}

		frontendFramework := flags.FrontendFramework(cmd.Flag("frontend").Value.String())
		executeProjectCreation(project, frontendFramework, flagSvelteKitTemplate, flagSvelteKitTypes, flagSvelteKitPackageManager, !flagNoInteractive)
		if spec.Execution.Output != scaffold.OutputJSON {
			printNextSteps(project, frontendFramework, flagSupabaseMode)
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

// handleSvelteKitSetup manages the SvelteKit frontend creation flow
func handleSvelteKitSetup(project *program.Project, template flags.SvelteKitTemplate, types flags.SvelteKitTypes, packageManager flags.SvelteKitPackageManager) error {
	// Check if Node.js and npx are installed
	if err := sveltekit.CheckNodeJS(); err != nil {
		return err
	}

	if err := sveltekit.CheckNPX(); err != nil {
		return err
	}

	// Check if pnpm is installed (user preference)
	if packageManager == flags.PNPM || packageManager == "" {
		if err := sveltekit.CheckPNPM(); err != nil {
			fmt.Println(warningStyle.Render("⚠️ Warning: pnpm not found, will use npm instead"))
			packageManager = flags.NPM
		}
	}

	// Use the parent directory of the Go project for SvelteKit
	parentPath := project.AbsolutePath // This is the parent directory where Go project was created
	frontendName := project.ProjectName + "-frontend"

	// Create SvelteKit manager
	manager := sveltekit.NewSvelteKitManager(parentPath, frontendName)

	// Set options if provided via flags
	manager.SetOptions(template, types, packageManager)

	fmt.Println(headerStyle.Render("\n🚀 Setting up SvelteKit frontend...\n"))

	// Determine which mode to use based on flags
	if template == "" && types == "" && packageManager == "" {
		fmt.Println(tipMsgStyle.Render("💡 Creating SvelteKit project with interactive prompts..."))
		fmt.Println(secondaryStyle.Render("You will be able to answer the prompts to customize your SvelteKit project"))

		// Use interactive mode - npx sv create will handle all the prompts
		if err := manager.CreateInteractive(); err != nil {
			return fmt.Errorf("failed to create SvelteKit project: %w", err)
		}
	} else {
		// Use preset mode with flags
		fmt.Println(tipMsgStyle.Render("💡 Creating SvelteKit project with preset options..."))
		fmt.Printf("Template: %s, Types: %s, Package Manager: %s\n",
			manager.Template, manager.Types, manager.PackageManager)

		if err := manager.Create(); err != nil {
			return fmt.Errorf("failed to create SvelteKit project: %w", err)
		}
	}

	fmt.Println(successStyle.Render("✅ SvelteKit frontend setup completed successfully!"))

	return nil
}

// handleNextJSSetup manages the Next.js frontend creation flow
func handleNextJSSetup(project *program.Project, packageManager flags.SvelteKitPackageManager) error {
	// Check if Node.js and npx are installed
	if err := nextjs.CheckNodeJS(); err != nil {
		return err
	}

	if err := nextjs.CheckNPX(); err != nil {
		return err
	}

	// Check if pnpm is installed (user preference)
	if packageManager == flags.PNPM || packageManager == "" {
		if err := nextjs.CheckPNPM(); err != nil {
			fmt.Println(warningStyle.Render("⚠️ Warning: pnpm not found, will use npm instead"))
			packageManager = flags.NPM
		}
	}

	// Use the parent directory of the Go project for Next.js
	parentPath := project.AbsolutePath // This is the parent directory where Go project was created
	frontendName := project.ProjectName + "-frontend"

	// Create Next.js manager
	manager := nextjs.NewNextJSManager(parentPath, frontendName)

	// Set package manager option
	manager.PackageManager = packageManager

	fmt.Println(headerStyle.Render("\n🚀 Setting up Next.js frontend...\n"))

	// Use interactive mode for now (similar to SvelteKit)
	fmt.Println(tipMsgStyle.Render("💡 Creating Next.js project with interactive prompts..."))
	fmt.Println(secondaryStyle.Render("You will be able to answer the prompts to customize your Next.js project"))

	if err := manager.CreateInteractive(); err != nil {
		return fmt.Errorf("failed to create Next.js project: %w", err)
	}

	fmt.Println(successStyle.Render("✅ Next.js frontend setup completed successfully!"))

	return nil
}

func runOptionalBootstrap(project *program.Project, frontendFramework flags.FrontendFramework, template flags.SvelteKitTemplate, types flags.SvelteKitTypes, packageManager flags.SvelteKitPackageManager) {
	var err error

	switch frontendFramework {
	case flags.SvelteKitFrontend:
		err = handleSvelteKitSetup(project, template, types, packageManager)
		if err != nil {
			fmt.Println(warningStyle.Render("\n⚠️ SvelteKit setup was skipped or failed:"))
			fmt.Println(secondaryStyle.Render(fmt.Sprintf("   %v", err)))
			fmt.Println(tipMsgStyle.Render("💡 You can create the frontend manually later with:"))
			fmt.Println(secondaryStyle.Render(fmt.Sprintf("   npx sv create %s-frontend", project.ProjectName)))
		}
	case flags.NextJSFrontend:
		err = handleNextJSSetup(project, packageManager)
		if err != nil {
			fmt.Println(warningStyle.Render("\n⚠️ Next.js setup was skipped or failed:"))
			fmt.Println(secondaryStyle.Render(fmt.Sprintf("   %v", err)))
			fmt.Println(tipMsgStyle.Render("💡 You can create the frontend manually later with:"))
			fmt.Println(secondaryStyle.Render(fmt.Sprintf("   npx create-next-app@latest %s-frontend", project.ProjectName)))
		}
	}
}

func collectProjectName(cmd *cobra.Command, project *program.Project, options Options, noInteractive bool) {
	if project.ProjectName != "" {
		return
	}
	if noInteractive {
		cobra.CheckErr(textinput.CreateErrorInputModel(fmt.Errorf("project name is required when --no-interactive is set")).Err())
	}
	tprogram := tea.NewProgram(textinput.InitialTextInputModel(options.ProjectName, "What is the name of your project?", project))
	if _, err := tprogram.Run(); err != nil {
		log.Printf("Name of project contains an error: %v", err)
		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
	}
	if options.ProjectName.Output != "" && !utils.ValidateModuleName(options.ProjectName.Output) {
		err := fmt.Errorf("'%s' is not a valid module name. Please choose a different name", options.ProjectName.Output)
		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
	}
	rootDirName := utils.GetRootDir(options.ProjectName.Output)
	if doesDirectoryExistAndIsNotEmpty(rootDirName) {
		err := fmt.Errorf("directory '%s' already exists and is not empty. Please choose a different name", rootDirName)
		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
	}
	project.ExitCLI(tprogram)
	project.ProjectName = options.ProjectName.Output
	if err := cmd.Flag("name").Value.Set(project.ProjectName); err != nil {
		log.Fatal("failed to set the name flag value", err)
	}
}

func collectDatabaseDriver(cmd *cobra.Command, project *program.Project, options Options, steps *steps.Steps, noInteractive bool) {
	if project.DBDriver != "" {
		return
	}
	if noInteractive {
		cobra.CheckErr(textinput.CreateErrorInputModel(fmt.Errorf("database driver is required when --no-interactive is set")).Err())
	}
	step := steps.Steps["driver"]
	tprogram := tea.NewProgram(singleSelection.InitialModelMulti(step.Options, options.DBDriver, step.Headers, project))
	if _, err := tprogram.Run(); err != nil {
		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
	}
	project.ExitCLI(tprogram)
	project.DBDriver = flags.Database(strings.ToLower(options.DBDriver.Choice))
	if err := cmd.Flag("driver").Value.Set(project.DBDriver.String()); err != nil {
		log.Fatal("failed to set the driver flag value", err)
	}
}

func collectFrontendFramework(cmd *cobra.Command, project *program.Project, options Options, steps *steps.Steps, current flags.FrontendFramework, noInteractive bool) {
	if current == "" {
		if noInteractive {
			cobra.CheckErr(textinput.CreateErrorInputModel(fmt.Errorf("frontend framework is required when --no-interactive is set")).Err())
		}
		step := steps.Steps["frontend"]
		tprogram := tea.NewProgram(singleSelection.InitialModelMulti(step.Options, options.Frontend, step.Headers, project))
		if _, err := tprogram.Run(); err != nil {
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}
		project.ExitCLI(tprogram)
		current = frontendChoiceFromSelection(options.Frontend.Choice)
	}
	if err := cmd.Flag("frontend").Value.Set(current.String()); err != nil {
		log.Fatal("failed to set the frontend flag value", err)
	}
}

func collectGitMode(cmd *cobra.Command, project *program.Project, options Options, steps *steps.Steps, noInteractive bool) {
	if project.GitOptions != "" {
		return
	}
	if noInteractive {
		cobra.CheckErr(textinput.CreateErrorInputModel(fmt.Errorf("git mode is required when --no-interactive is set")).Err())
	}
	step := steps.Steps["git"]
	tprogram := tea.NewProgram(singleSelection.InitialModelMulti(step.Options, options.Git, step.Headers, project))
	if _, err := tprogram.Run(); err != nil {
		cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
	}
	project.ExitCLI(tprogram)
	project.GitOptions = flags.Git(strings.ToLower(options.Git.Choice))
	if err := cmd.Flag("git").Value.Set(project.GitOptions.String()); err != nil {
		log.Fatal("failed to set the git flag value", err)
	}
}

func frontendChoiceFromSelection(choice string) flags.FrontendFramework {
	switch strings.ToLower(choice) {
	case "sveltekit":
		return flags.SvelteKitFrontend
	case "next.js":
		return flags.NextJSFrontend
	case "none":
		return flags.NoneFrontend
	default:
		return flags.NoneFrontend
	}
}

func executeProjectCreation(project *program.Project, frontendFramework flags.FrontendFramework, template flags.SvelteKitTemplate, types flags.SvelteKitTypes, packageManager flags.SvelteKitPackageManager, useSpinner bool) {
	if useSpinner {
		spinnerProgram := tea.NewProgram(spinner.InitialModelNew())
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := spinnerProgram.Run(); err != nil {
				cobra.CheckErr(err)
			}
		}()

		defer func() {
			if r := recover(); r != nil {
				fmt.Println("The program encountered an unexpected issue and had to exit. The error was:", r)
				fmt.Println("If you continue to experience this issue, please post a message on our GitHub page or join our Discord server for support.")
				if releaseErr := spinnerProgram.ReleaseTerminal(); releaseErr != nil {
					log.Printf("Problem releasing terminal: %v", releaseErr)
				}
			}
		}()

		err := project.CreateMainFile()
		if err != nil {
			if releaseErr := spinnerProgram.ReleaseTerminal(); releaseErr != nil {
				log.Printf("Problem releasing terminal: %v", releaseErr)
			}
			log.Printf("Problem creating files for project.")
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}

		if err := spinnerProgram.ReleaseTerminal(); err != nil {
			log.Printf("Could not release terminal: %v", err)
		}
	} else {
		if err := project.CreateMainFile(); err != nil {
			log.Printf("Problem creating files for project.")
			cobra.CheckErr(textinput.CreateErrorInputModel(err).Err())
		}
	}

	if !project.SkipInstall {
		runOptionalBootstrap(project, frontendFramework, template, types, packageManager)
	}
}

func printNextSteps(project *program.Project, frontendFramework flags.FrontendFramework, supabaseMode flags.SupabaseMode) {
	for _, line := range nextStepLines(project, frontendFramework, supabaseMode) {
		fmt.Println(line)
	}
	fmt.Println()
}

func nextStepLines(project *program.Project, frontendFramework flags.FrontendFramework, supabaseMode flags.SupabaseMode) []string {
	lines := []string{
		headerStyle.Render("\n🎉 Project created successfully!\n"),
		successStyle.Render("Next steps:"),
		fmt.Sprintf("• %-25s %s", fmt.Sprintf("cd %s", utils.GetRootDir(project.ProjectName)), secondaryStyle.Render("# Change to project directory")),
	}

	if project.DBDriver == flags.Supabase {
		if supabaseMode == flags.LocalDB {
			lines = append(lines, fmt.Sprintf("• %-25s %s", "supabase status", secondaryStyle.Render("# Check Supabase local instance")))
		} else {
			lines = append(lines,
				fmt.Sprintf("• %-25s %s", "supabase link", secondaryStyle.Render("# Link to your Supabase project")),
				fmt.Sprintf("• %-25s %s", "supabase start", secondaryStyle.Render("# Start local development")),
			)
		}
		if project.AdvancedOptions[string(flags.Sqlc)] {
			lines = append(lines, fmt.Sprintf("• %-25s %s", "sqlc generate", secondaryStyle.Render("# Generate type-safe Go code from SQL")))
		}
	} else {
		if project.AdvancedOptions[string(flags.Sqlc)] {
			lines = append(lines, fmt.Sprintf("• %-25s %s", "make sqlc-generate", secondaryStyle.Render("# Generate type-safe Go code from SQL")))
		}
		if project.AdvancedOptions[string(flags.Docker)] && project.DBDriver != "none" {
			lines = append(lines, fmt.Sprintf("• %-25s %s", "make docker-run", secondaryStyle.Render("# Start PostgreSQL database")))
		}
	}

	lines = append(lines, fmt.Sprintf("• %-25s %s", "make run", secondaryStyle.Render("# Start the server")))

	switch frontendFramework {
	case flags.SvelteKitFrontend:
		frontendName := project.ProjectName + "-frontend"
		lines = append(lines,
			fmt.Sprintf("• %-25s %s", fmt.Sprintf("cd %s", frontendName), secondaryStyle.Render("# Switch to frontend directory")),
			fmt.Sprintf("• %-25s %s", "pnpm dev", secondaryStyle.Render("# Start SvelteKit development server")),
		)
	case flags.NextJSFrontend:
		frontendName := project.ProjectName + "-frontend"
		lines = append(lines,
			fmt.Sprintf("• %-25s %s", fmt.Sprintf("cd %s", frontendName), secondaryStyle.Render("# Switch to frontend directory")),
			fmt.Sprintf("• %-25s %s", "pnpm dev", secondaryStyle.Render("# Start Next.js development server")),
		)
	}

	return lines
}

func buildScaffoldSpec(
	cmd *cobra.Command,
	project *program.Project,
	advancedChoices map[string]bool,
	exec scaffold.ExecutionOptions,
) scaffold.CreateSpec {
	return scaffold.BuildSpec(
		project.ProjectName,
		project.DBDriver,
		mapKeys(advancedChoices, project.AdvancedOptions),
		project.GitOptions,
		flags.SupabaseMode(cmd.Flag("supabase-mode").Value.String()),
		flags.FrontendFramework(cmd.Flag("frontend").Value.String()),
		flags.SvelteKitTemplate(cmd.Flag("sveltekit-template").Value.String()),
		flags.SvelteKitTypes(cmd.Flag("sveltekit-types").Value.String()),
		flags.SvelteKitPackageManager(cmd.Flag("sveltekit-package-manager").Value.String()),
		exec,
	)
}

func maybeHandleDryRun(spec scaffold.CreateSpec, baseDir string) error {
	if !spec.Execution.DryRun {
		return nil
	}
	plan := scaffold.BuildPlan(spec, baseDir)
	return scaffold.WritePlan(os.Stdout, plan, spec.Execution.Output)
}

func mapKeys(interactive map[string]bool, existing map[string]bool) []string {
	keys := make([]string, 0, len(interactive)+len(existing))
	seen := map[string]bool{}
	for key, enabled := range existing {
		if enabled && !seen[key] {
			seen[key] = true
			keys = append(keys, key)
		}
	}
	for key, enabled := range interactive {
		if enabled {
			normalized := strings.ToLower(key)
			if !seen[normalized] {
				seen[normalized] = true
				keys = append(keys, normalized)
			}
		}
	}
	return keys
}
