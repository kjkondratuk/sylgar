package app

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kjkondratuk/sylgar/api"
	"log/slog"
	"time"
)

type (
	tickMsg struct{}
)

type MainModel struct {
	user          string
	capi          api.ChatApi
	refreshTimer  tea.Cmd
	chats         []string
	activeChannel string
	activeMessage string
}

func New(user string, capi api.ChatApi) MainModel {
	return MainModel{
		user: user,
		capi: capi,
	}
}

func (m MainModel) Init() tea.Cmd {
	return tick()
}

func (m MainModel) View() string {
	dsp := ""
	for _, i := range m.chats {
		dsp += i + "\n"
	}
	return dsp
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		chats, err := m.capi.FetchChatsForUser(m.user, api.WithBefore(time.Now().UTC()))
		if err != nil {
			return nil, tea.Printf("ERROR: %s", err.Error())
		}

		m.chats = make([]string, 0)
		for _, c := range chats {
			if len(m.chats) < 20 {
				m.chats = append(m.chats, fmt.Sprintf("%f : %s|%s - %s", c.T, c.Channel, c.User, c.Msg))
			} else {
				break
			}
		}

		return m, tick()
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case msg.Type == tea.KeyEnter:
			slog.Info("sending message", "msg", m.activeMessage)
		}
	}

	return m, nil
}

func tick() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
