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
		bubbleText := `moomsay is a command line tool by https://farcaster.xyz/vrypan

moomsay is also an art project, about the history of computer graphics, combining ANSI art, Pixel art, and NFTs.

Based on MoomBirdsAtHome, an NFT project by NFTtank: https://x.com/moombirbsathome
`

		Say(moomId, bubbleText, isReply)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}
