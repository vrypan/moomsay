package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var sayCmd = &cobra.Command{
	Use:   "say <moom-id>",
	Short: "Display an image and a speech bubble",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		moomId, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Printf("Invalid moom-id: %s\n", args[0])
			os.Exit(1)
		}

		isReply := false
		bubbleText, _ := cmd.Flags().GetString("text")

		// Read bubble text from stdin if not provided
		if bubbleText == "" {
			scanner := bufio.NewScanner(os.Stdin)
			var inputLines []string
			for scanner.Scan() {
				inputLines = append(inputLines, scanner.Text())
			}
			bubbleText = strings.Join(inputLines, " ")
		}

		Say(moomId, bubbleText, isReply)
	},
}

func init() {
	rootCmd.AddCommand(sayCmd)
	sayCmd.Flags().String("text", "", "Text for the speech bubble (if empty, read from stdin)")
}
