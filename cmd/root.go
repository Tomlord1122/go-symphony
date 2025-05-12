/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-symphony",
	Short: "A CLI tool to quickly create a Go project with Chi and GORM",
	Long:  `Go Symphony is a CLI tool that can build a Go web project with some popluar framework! (Personally, I like to use Chi and GORM)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
