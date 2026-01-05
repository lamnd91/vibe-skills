package cli

import (
	"fmt"
	"os"

	"github.com/cuongtl/vibe-skills/internal/installer"
	"github.com/cuongtl/vibe-skills/internal/registry"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for skills",
	Long: `Search for skills by name or description.

Examples:
  vibe-skills search database
  vibe-skills search "code review"`,
	Args: cobra.ExactArgs(1),
	RunE: runSearch,
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	reg, err := registry.New()
	if err != nil {
		return fmt.Errorf("failed to load skills registry: %w", err)
	}

	inst := installer.New(reg, cwd)

	results := reg.Search(query)
	if len(results) == 0 {
		fmt.Printf("No skills found matching: %s\n", query)
		return nil
	}

	fmt.Printf("Found %d skill(s) matching '%s':\n\n", len(results), query)
	for _, skill := range results {
		installed := ""
		if inst.IsInstalled(skill.Name) {
			installed = " [installed]"
		}
		fmt.Printf("  %s/%s%s\n", skill.Stack, skill.Name, installed)
		if skill.Description != "" {
			fmt.Printf("    %s\n", skill.Description)
		}
		fmt.Println()
	}

	return nil
}
