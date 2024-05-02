package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kjkondratuk/sylgar/api"
	"github.com/kjkondratuk/sylgar/internal/app"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		slog.Error("invalid number of arguments: must specify username followed by chat password")
		os.Exit(-1)
	}
	user := os.Args[1]
	capi := api.New(http.DefaultClient, os.Args[2])
	// start the application on the main thread
	if _, err := tea.NewProgram(app.New(user, capi)).Run(); err != nil {
		slog.Error("", "error", err)
		os.Exit(-1)
	}
}
