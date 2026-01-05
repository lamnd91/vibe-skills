package cli

import (
	"fmt"

	"github.com/cuongtl/vibe-skills/internal/updater"
	"github.com/cuongtl/vibe-skills/internal/version"
	"github.com/spf13/cobra"
)

var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update vibe-skills to the latest version",
	Long:  `Downloads and installs the latest version of vibe-skills from GitHub releases.`,
	RunE:  runSelfUpdate,
}

func runSelfUpdate(cmd *cobra.Command, args []string) error {
	fmt.Printf("Current version: %s\n", version.GetVersion())
	fmt.Println("Checking for updates...")

	latestVersion, hasUpdate, err := updater.CheckForUpdate()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !hasUpdate {
		fmt.Println("You are already running the latest version.")
		return nil
	}

	fmt.Printf("New version available: %s\n", latestVersion)
	fmt.Println("Downloading update...")

	if err := updater.SelfUpdate(); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	fmt.Printf("Successfully updated to version %s\n", latestVersion)
	return nil
}
