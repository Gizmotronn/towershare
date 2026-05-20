package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	baseURL := flag.String("url", "http://localhost:8090", "PocketBase base URL")
	from := flag.String("from", "+61400000001", "Sender identifier (phone number)")
	flag.Parse()

	m := newModel(*baseURL, *from)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
