package cli

import (
	"fmt"

	"github.com/cuongtl1992/vibe-skills/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vibe-skills %s\n", version.GetFullVersion())
	},
}
