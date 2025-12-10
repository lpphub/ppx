package cmd

import (
	"fmt"
	"os"

	"ppx/generator"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go web project",
	Long:  `Create a new Go web project with predefined directory structure including cmd, infra, logic, web, and config folders.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		if err := generator.CreateProject(projectName); err != nil {
			fmt.Printf("Error creating project: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Project '%s' created successfully!\n", projectName)
	},
}
