package main

import (
	"fmt"
	"os"

	"github.com/aga-absolut/Vault-System/internal/client"
	"github.com/aga-absolut/Vault-System/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.InitialModel(client.NewClient(":50051")), tea.WithAltScreen())

	fmt.Println("Запуск GophKeeper TUI...")
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}
}
