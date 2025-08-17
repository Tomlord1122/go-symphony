/*
Copyright © 2025 HSIU-CHI LIU aa12359346@gmail.com
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-symphony",
	Short: "Modern Go project generator with Gin + PostgreSQL + SQLC + Supabase",
	Long: `Symphony is a CLI tool for creating modern Go projects with best practices built-in.
Features: Gin web framework, PostgreSQL with pgx driver, SQLC for type-safe queries, Supabase for auth and real-time features, and more.`,
}

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
