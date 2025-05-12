/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

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
	logoStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)
	tipMsgStyle    = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("190")).Italic(true)
	endingMsgStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("170")).Bold(true)
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Go project and selected your favorite framework",
	Long: `Create a new Go project and selected your favorite framework.
	You can choose the project name, and the project will be created in the current directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
