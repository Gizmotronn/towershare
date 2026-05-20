package main

import "github.com/charmbracelet/lipgloss"

var (
	// Colour palette
	colAccent  = lipgloss.Color("#00b4d8")
	colSurface = lipgloss.Color("#1e1e1e")
	colBorder  = lipgloss.Color("#2a2a2a")
	colMuted   = lipgloss.Color("#707070")
	colWhite   = lipgloss.Color("#f0f0f0")
	colGreen   = lipgloss.Color("#34c759")
	colOrange  = lipgloss.Color("#ff9500")

	// Sent bubble (right-aligned, accent background)
	sentStyle = lipgloss.NewStyle().
			Background(colAccent).
			Foreground(colWhite).
			Padding(0, 1).
			MaxWidth(55)

	// Received bubble (left-aligned, dark surface)
	receivedStyle = lipgloss.NewStyle().
			Background(colSurface).
			Foreground(colWhite).
			Padding(0, 1).
			MaxWidth(55)

	// Sender label above received bubble
	senderStyle = lipgloss.NewStyle().
			Foreground(colMuted).
			PaddingLeft(1).
			Italic(true)

	// Action chip — unselected
	chipStyle = lipgloss.NewStyle().
			Foreground(colAccent).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colBorder).
			Padding(0, 1)

	// Action chip — selected
	chipSelectedStyle = lipgloss.NewStyle().
				Foreground(colWhite).
				Background(colAccent).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colAccent).
				Padding(0, 1)

	// Input bar
	inputBarStyle = lipgloss.NewStyle().
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(colBorder).
			Padding(0, 1)

	// Header bar
	headerStyle = lipgloss.NewStyle().
			Background(colSurface).
			Foreground(colAccent).
			Bold(true).
			Padding(0, 2).
			Width(0) // set dynamically

	// Status bar
	statusStyle = lipgloss.NewStyle().
			Foreground(colMuted).
			Padding(0, 2)

	// Channel chip in header
	channelStyle = lipgloss.NewStyle().
			Foreground(colMuted).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colBorder).
			Padding(0, 1)

	channelActiveStyle = lipgloss.NewStyle().
				Foreground(colWhite).
				Background(colAccent).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colAccent).
				Padding(0, 1)

	// Dot separators
	dotStyle = lipgloss.NewStyle().Foreground(colMuted)
)
