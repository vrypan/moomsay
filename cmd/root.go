package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "moomsay",
	Short: "A CLI app with speech bubbles",
}

// Execute is the entry point
func Execute() {
	// Check before parsing if we need to inject the default subcommand
	if len(os.Args) > 1 && !isKnownSubcommand(os.Args[1]) {
		// Inject the default subcommand name here (example: "say")
		os.Args = append([]string{os.Args[0], "say"}, os.Args[1:]...)
	}

	// Now run Cobra as usual
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// Checks if the first argument is a known subcommand
func isKnownSubcommand(arg string) bool {
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == arg || contains(cmd.Aliases, arg) {
			return true
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
