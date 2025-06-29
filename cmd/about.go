package cmd

import (
	"fmt"
	"math/rand"

	"github.com/spf13/cobra"
)

var aboutCmd = &cobra.Command{
	Use: "about",
	Run: func(cmd *cobra.Command, args []string) {
		moomId := 1 + rand.Intn(10000)

		isReply := false
		bubbleText := `Based on MoomBirdsAtHome, an NFT project by NFTtank: https://x.com/moombirbsathome
moomsay is a project by https://farcaster.xyz/vrypan
		`

		Say(moomId, bubbleText, isReply)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}
