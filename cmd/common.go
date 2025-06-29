package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/vrypan/moomsay/pkg/assets"
	"golang.org/x/term"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// extractLines splits content into lines, appending a space to each.
func extractLines(content string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		lines = append(lines, scanner.Text()+" ")
	}
	return lines
}

// buildSpeechBubble creates a "speech bubble" of text, wrapping and padding as needed.
func buildSpeechBubble(text string, maxWidth int, rightAlign bool) []string {
	paragraphs := strings.Split(text, "\n")
	var lines []string

	lines = append(lines, "──")

	for _, paragraph := range paragraphs {
		words := strings.Fields(paragraph)
		var currentLine string

		for _, word := range words {
			if utf8.RuneCountInString(currentLine+" "+word) > maxWidth {
				lines = append(lines, strings.TrimSpace(currentLine))
				currentLine = word + " "
			} else {
				currentLine += word + " "
			}
		}

		if currentLine != "" {
			lines = append(lines, strings.TrimSpace(currentLine))
		}

		// Add an empty line to preserve explicit line breaks
		lines = append(lines, "")
	}

	lines = append(lines, "──")

	bubble := make([]string, 0, len(lines))

	for i, line := range lines {
		padding := strings.Repeat(" ", maxWidth-utf8.RuneCountInString(line))
		var textPart string

		if rightAlign {
			textPart = padding + line
			switch i {
			case 0:
				bubble = append(bubble, fmt.Sprintf(" %s┐ ", textPart))
			case len(lines) - 1:
				bubble = append(bubble, fmt.Sprintf(" %s┘ ", textPart))
			case len(lines) - 5:
				bubble = append(bubble, fmt.Sprintf("%s ├─", textPart))
			default:
				bubble = append(bubble, fmt.Sprintf("%s │ ", textPart))
			}
		} else {
			textPart = line + padding
			switch i {
			case 0:
				bubble = append(bubble, fmt.Sprintf(" ┌%s", textPart))
			case len(lines) - 1:
				bubble = append(bubble, fmt.Sprintf(" └%s", textPart))
			case len(lines) - 5:
				bubble = append(bubble, fmt.Sprintf("─┤ %s", textPart))
			default:
				bubble = append(bubble, fmt.Sprintf(" │ %s", textPart))
			}
		}
	}

	return bubble
}

// printSideBySide prints image block and speech bubble with correct alignment.
func printSideBySide(block, bubble []string, bubblePosition string) {
	imageWidth := getImageWidth(block)
	bubbleWidth := getImageWidth(bubble)

	blockHeight := len(block)
	bubbleHeight := len(bubble)
	totalLines := blockHeight
	if bubbleHeight > totalLines {
		totalLines = bubbleHeight
	}
	blockPadTop := totalLines - blockHeight + 1
	bubblePadTop := totalLines - bubbleHeight

	for i := 0; i < totalLines; i++ {
		var blockPart, bubblePart string

		if bubblePosition == "right" {
			if i < blockPadTop {
				blockPart = strings.Repeat(" ", imageWidth)
			} else {
				blockPart = block[i-blockPadTop]
			}
			if i < bubblePadTop {
				bubblePart = ""
			} else {
				bubblePart = bubble[i-bubblePadTop]
			}

			padding := imageWidth - utf8.RuneCountInString(ansiRegexp.ReplaceAllString(blockPart, ""))
			if padding < 0 {
				padding = 0
			}

			fmt.Printf("%s%s  %s\n", blockPart, strings.Repeat(" ", padding), bubblePart)

		} else {
			if i < bubblePadTop {
				bubblePart = strings.Repeat(" ", bubbleWidth)
			} else {
				bubblePart = bubble[i-bubblePadTop]
			}
			if i < blockPadTop {
				blockPart = ""
			} else {
				blockPart = block[i-blockPadTop]
			}

			padding := bubbleWidth - utf8.RuneCountInString(ansiRegexp.ReplaceAllString(bubblePart, ""))
			if padding < 0 {
				padding = 0
			}

			fmt.Printf("%s%s  %s\n", bubblePart, strings.Repeat(" ", padding), blockPart)
		}
	}
}

// mirrorImage reverses the given image lines, preserving ANSI blocks.
func mirrorImage(lines []string) []string {
	var mirrored []string
	ansiUnitRegexp := regexp.MustCompile(`(\x1b\[[0-9;]*m ?)+ ?`)

	for _, line := range lines {
		var units []string
		matches := ansiUnitRegexp.FindAllStringIndex(line, -1)
		cursor := 0
		for _, match := range matches {
			if match[0] > cursor {
				units = append(units, line[cursor:match[0]])
			}
			unit := line[match[0]:match[1]]
			if !strings.Contains(unit, "\x1b[0m") {
				unit += "\x1b[0m"
			}
			units = append(units, unit)
			cursor = match[1]
		}
		if cursor < len(line) {
			units = append(units, line[cursor:])
		}
		for i, j := 0, len(units)-1; i < j; i, j = i+1, j-1 {
			units[i], units[j] = units[j], units[i]
		}
		mirrored = append(mirrored, strings.Join(units, ""))
	}
	return mirrored
}

// getImageWidth computes the max visible width of block lines (excludes ANSI escapes).
func getImageWidth(block []string) int {
	maxWidth := 0
	for _, line := range block {
		plain := ansiRegexp.ReplaceAllString(line, "")
		width := utf8.RuneCountInString(plain)
		if width > maxWidth {
			maxWidth = width
		}
	}
	return maxWidth
}

func Say(id int, text string, isReply bool) {
	filename := fmt.Sprintf("assets/%d.txt", id)
	content, err := assets.Files.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading %s: %v\n", filename, err)
		os.Exit(1)
	}

	blockLines := extractLines(string(content))

	termWidth, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Could not get terminal size, defaulting width to 80")
		termWidth = 80
	}

	imageWidth := getImageWidth(blockLines)
	availableWidth := termWidth - imageWidth - 4

	bubbleWidth := availableWidth
	if bubbleWidth > 80 {
		bubbleWidth = 80
	}

	bubbleLines := buildSpeechBubble(text, bubbleWidth-2, !isReply)

	if isReply {
		printSideBySide(mirrorImage(blockLines), bubbleLines, "right")
	} else {
		printSideBySide(blockLines, bubbleLines, "left")
	}
}
