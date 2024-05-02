package app

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kjkondratuk/sylgar/api"
	"log/slog"
	"strings"
	"time"
)

type (
	// TODO : try using bubbles timer instead?
	tickMsg struct{}
)

type MainModel struct {
	user          string
	viewport      viewport.Model
	capi          api.ChatApi
	refreshTimer  tea.Cmd
	chats         []string
	activeChannel string
	activeMessage string
	ready         bool
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
	if !m.ready {
		return "\n  Loading..."
	}
	return fmt.Sprintf("%s\n%s", m.viewport.View())
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		//headerHeight := lipgloss.Height(m.headerView())
		//footerHeight := lipgloss.Height(m.footerView())
		//verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height)
			m.viewport.YPosition = 1
			//m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(strings.Join(m.chats, "\n"))
			//if msg.Width != 0 && msg.Height != 0 {
			m.ready = true
			//}

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			//m.viewport.YPosition = 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}

		//cmds = append(cmds, viewport.Sync(m.viewport))

		// TODO : don't use the UI to schedule this
	case tickMsg:
		chats, err := m.capi.FetchChatsForUser(m.user, api.WithBefore(time.Now().UTC()))
		if err != nil {
			return nil, tea.Printf("ERROR: %s", err.Error())
		}

		m.chats = make([]string, 0)
		for _, c := range chats {
			if c.IsJoin {
				m.chats = append(m.chats, fmt.Sprintf("%f : %s ----> %s", c.T, c.Channel, c.User))
			} else if c.IsLeave {
				m.chats = append(m.chats, fmt.Sprintf("%f : %s <---- %s", c.T, c.Channel, c.User))
			} else {
				m.chats = append(m.chats, fmt.Sprintf("%f : %s|%s - %s", c.T, c.Channel, c.User, c.Msg))
			}
		}
		content := strings.Join(m.chats, "\n")
		m.viewport.SetContent(content)
		cmds = append(cmds, tick())
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case msg.Type == tea.KeyEnter:
			slog.Info("sending message", "msg", m.activeMessage)
		}
	}

	// handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append([]tea.Cmd{cmd}, cmds...)

	return m, tea.Batch(cmds...)
}

func tick() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}
