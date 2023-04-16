package cmd

import (
	"github.com/spf13/cobra"

	"github.com/code-game-project/codegame-cli/doctor"
)

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check for missing dependencies and misconfigurations",
	Run: func(_ *cobra.Command, _ []string) {
		doctor.Doctor()
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
