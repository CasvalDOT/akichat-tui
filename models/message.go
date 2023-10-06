package models

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var colors = map[string]string{
	"gray":    "#808080",
	"silver":  "#c0c0c0",
	"white":   "#ffffff",
	"yellow":  "#f8f801",
	"orange":  "#ffa500",
	"red":     "#ff0000",
	"fuchsia": "#ff00ff",
	"purple":  "#800080",
	"navy":    "#000080",
	"blue":    "#0000ff",
	"aqua":    "#00ffff",
	"teal":    "#008080",
	"green":   "#008000",
	"lime":    "#00ff00",
	"olive":   "#808000",
	"maroon":  "#800000",
	"black":   "#000000",
}

type MessageText string

func (m MessageText) ToString() string {
	return string(m)
}

func (m MessageText) Parse() string {
	wordRegex := regexp.MustCompile(`(?m)(\[.*?\].*?\[.*?\]|\w|\W)`)
	tagRegex := regexp.MustCompile(`(?m)\[(.*?)(=(.*?)|)\](.*?)\[\/.*?\]`)
	cmdRegex := regexp.MustCompile(`^(\/.*?)\s`)

	content := []string{}

	baseContent := m.ToString()

	baseContent = cmdRegex.ReplaceAllString(baseContent, "[color=red]$1[/color] ")

	words := wordRegex.FindAllString(baseContent, -1)

	if len(words) > 0 {
		for _, word := range words {
			attributes := tagRegex.FindAllStringSubmatch(word, -1)

			if attributes == nil {
				content = append(content, word)
				continue
			}

			tagName := attributes[0][1]
			tagValue := attributes[0][3]
			text := attributes[0][4]

			switch tagName {
			case "b":
				content = append(content, lipgloss.NewStyle().Bold(true).Render(text))
			case "i":
				content = append(content, lipgloss.NewStyle().Italic(true).Render(text))
			case "u":
				content = append(content, lipgloss.NewStyle().Underline(true).Render(text))
			case "quote":
				content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render(text))
			case "code":
				content = append(content, lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render(text))
			case "color":
				matchColor, ok := colors[tagValue]
				if !ok {
					matchColor = tagValue
				}
				content = append(
					content,
					lipgloss.NewStyle().Foreground(lipgloss.Color(matchColor)).Render(text),
				)
			default:
				content = append(content, text)
			}
		}
	}
	return strings.Join(content, "")
}
