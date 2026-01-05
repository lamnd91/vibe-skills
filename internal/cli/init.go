package cli

import (
	"fmt"
	"os"

	"github.com/cuongtl/vibe-skills/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new .vibe-skills.yaml config file",
	Long:  `Creates a new .vibe-skills.yaml configuration file in the current directory with default skills.`,
	RunE:  runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if config.Exists(cwd) {
		return fmt.Errorf("config file already exists: %s", config.ConfigFileName)
	}

	cfg := config.GetDefaultConfig()
	if err := config.Save(cwd, cfg); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	fmt.Printf("Created %s with default skills:\n", config.ConfigFileName)
	for _, skill := range cfg.Skills {
		fmt.Printf("  - %s\n", skill)
	}
	fmt.Println("\nRun 'vibe-skills install' to install these skills.")

	return nil
}
