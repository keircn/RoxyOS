package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	primaryColor = lipgloss.Color("#1a5276")
	accentColor  = lipgloss.Color("#5dade2")
	bgColor      = lipgloss.Color("#0a1628")
	textColor    = lipgloss.Color("#e8f4fc")
	mutedColor   = lipgloss.Color("#a8c8d8")
	successColor = lipgloss.Color("#4a6b5a")
	errorColor   = lipgloss.Color("#8b4a2b")
	warningColor = lipgloss.Color("#c9a227")

	titleStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginBottom(2)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(1, 2)

	successStyle = lipgloss.NewStyle().Foreground(successColor)
	errorStyle   = lipgloss.NewStyle().Foreground(errorColor)
	infoStyle    = lipgloss.NewStyle().Foreground(accentColor)
	mutedStyle   = lipgloss.NewStyle().Foreground(mutedColor)
	accentStyle  = lipgloss.NewStyle().Foreground(accentColor)
)

type step int

const (
	stepWelcome step = iota
	stepSelectDM
	stepSelectComponents
	stepBackup
	stepInstall
	stepPlymouth
	stepComplete
)

type component struct {
	name        string
	description string
	enabled     bool
}

type model struct {
	step          step
	width, height int
	spinner       spinner.Model
	dmList        list.Model
	components    []component
	cursor        int
	selectedDM    string
	backupPath    string
	logs          []string
	err           error
	installing    bool
	done          bool
}

type dmItem struct {
	name, desc string
}

func (i dmItem) Title() string       { return i.name }
func (i dmItem) Description() string { return i.desc }
func (i dmItem) FilterValue() string { return i.name }

type installMsg struct{ log string }
type doneMsg struct{ err error }

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(accentColor)

	items := []list.Item{
		dmItem{"sddm", "Simple Desktop Display Manager with RoxyOS theme"},
		dmItem{"ly", "TUI display manager for console lovers"},
		dmItem{"none", "Skip display manager configuration"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(accentColor)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(mutedColor)

	dmList := list.New(items, delegate, 60, 10)
	dmList.Title = "Select Display Manager"
	dmList.SetShowStatusBar(false)
	dmList.SetFilteringEnabled(false)
	dmList.Styles.Title = titleStyle

	return model{
		step:    stepWelcome,
		spinner: s,
		dmList:  dmList,
		components: []component{
			{"niri", "Niri window manager configuration", true},
			{"waybar", "Status bar configuration", true},
			{"hyprlock", "Lock screen configuration", true},
			{"kitty", "Terminal emulator configuration", true},
			{"rofi", "Application launcher configuration", true},
			{"fish", "Fish shell and Starship prompt", true},
			{"mako", "Notification daemon configuration", true},
			{"plymouth", "Boot splash theme", true},
		},
	}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.dmList.SetSize(min(60, m.width-4), min(12, m.height-10))
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.step != stepInstall {
				return m, tea.Quit
			}
		case "enter":
			return m.handleEnter()
		case "up", "k":
			if m.step == stepSelectComponents && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.step == stepSelectComponents && m.cursor < len(m.components)-1 {
				m.cursor++
			}
		case " ":
			if m.step == stepSelectComponents {
				m.components[m.cursor].enabled = !m.components[m.cursor].enabled
			}
		case "a":
			if m.step == stepSelectComponents {
				allEnabled := true
				for _, c := range m.components {
					if !c.enabled {
						allEnabled = false
						break
					}
				}
				for i := range m.components {
					m.components[i].enabled = !allEnabled
				}
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case installMsg:
		m.logs = append(m.logs, msg.log)
		if len(m.logs) > 10 {
			m.logs = m.logs[1:]
		}
		return m, nil

	case doneMsg:
		m.installing = false
		m.done = true
		m.err = msg.err
		if msg.err == nil {
			m.step = stepComplete
		}
		return m, nil
	}

	if m.step == stepSelectDM {
		var cmd tea.Cmd
		m.dmList, cmd = m.dmList.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case stepWelcome:
		m.step = stepSelectDM
	case stepSelectDM:
		if i, ok := m.dmList.SelectedItem().(dmItem); ok {
			m.selectedDM = i.name
		}
		m.step = stepSelectComponents
	case stepSelectComponents:
		m.step = stepBackup
	case stepBackup:
		m.step = stepInstall
		m.installing = true
		return m, m.runInstall()
	case stepComplete:
		return m, tea.Quit
	}
	return m, nil
}

func (m model) runInstall() tea.Cmd {
	return func() tea.Msg {
		home, err := os.UserHomeDir()
		if err != nil {
			return doneMsg{err}
		}

		configDir := filepath.Join(home, ".config")
		dataDir := filepath.Join(home, ".local", "share")
		backupDir := filepath.Join(dataDir, "roxyos", "backups")
		roxyosDir := "/usr/share/roxyos"

		os.MkdirAll(backupDir, 0755)

		installed := 0
		for _, comp := range m.components {
			if !comp.enabled {
				continue
			}

			src := filepath.Join(roxyosDir, "configs", comp.name)
			dst := filepath.Join(configDir, comp.name)

			if comp.name == "hyprlock" {
				src = filepath.Join(roxyosDir, "configs", "hyprlock.conf")
				dst = filepath.Join(configDir, "hypr", "hyprlock.conf")
				os.MkdirAll(filepath.Join(configDir, "hypr"), 0755)
			}

			if comp.name == "starship" {
				src = filepath.Join(roxyosDir, "configs", "starship")
				dst = filepath.Join(configDir, "starship")
			}

			if _, err := os.Stat(src); err == nil {
				if err := copyDir(src, dst); err != nil {
					return doneMsg{fmt.Errorf("failed to copy %s: %w", comp.name, err)}
				}
				installed++
			}
		}

		if installed == 0 {
			return doneMsg{fmt.Errorf("no configurations were installed")}
		}

		return doneMsg{nil}
	}
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}

func (m model) View() string {
	var content string

	switch m.step {
	case stepWelcome:
		content = m.viewWelcome()
	case stepSelectDM:
		content = m.viewSelectDM()
	case stepSelectComponents:
		content = m.viewSelectComponents()
	case stepBackup:
		content = m.viewBackup()
	case stepInstall:
		content = m.viewInstall()
	case stepComplete:
		content = m.viewComplete()
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

func (m model) viewWelcome() string {
	banner := `
    ██████╗  ██████╗ ██╗  ██╗██╗   ██╗ ██████╗ ███████╗
    ██╔══██╗██╔═══██╗╚██╗██╔╝╚██╗ ██╔╝██╔═══██╗██╔════╝
    ██████╔╝██║   ██║ ╚███╔╝  ╚████╔╝ ██║   ██║███████╗
    ██╔══██╗██║   ██║ ██╔██╗   ╚██╔╝  ██║   ██║╚════██║
    ██║  ██║╚██████╔╝██╔╝ ██╗   ██║   ╚██████╔╝███████║
    ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚══════╝`

	bannerStyle := lipgloss.NewStyle().Foreground(accentColor)

	return lipgloss.JoinVertical(lipgloss.Center,
		bannerStyle.Render(banner),
		"",
		subtitleStyle.Render("Roxy Migurdia themed Niri configuration for Arch Linux"),
		"",
		infoStyle.Render("Press Enter to begin setup"),
		mutedStyle.Render("Press q to quit"),
	)
}

func (m model) viewSelectDM() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Display Manager"),
		"",
		m.dmList.View(),
		"",
		mutedStyle.Render("↑/↓ navigate • enter select"),
	)
}

func (m model) viewSelectComponents() string {
	var items []string

	for i, comp := range m.components {
		cursor := "  "
		if i == m.cursor {
			cursor = "▸ "
		}

		check := "○"
		if comp.enabled {
			check = successStyle.Render("●")
		}

		name := comp.name
		if i == m.cursor {
			name = accentStyle.Render(comp.name)
		}

		items = append(items, fmt.Sprintf("%s%s %s  %s", cursor, check, name, mutedStyle.Render(comp.description)))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Select Components"),
		subtitleStyle.Render("Choose which configurations to install"),
		"",
		strings.Join(items, "\n"),
		"",
		mutedStyle.Render("↑/↓ navigate • space toggle • a toggle all • enter continue"),
	)
}

func (m model) viewBackup() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("Ready to Install"),
		"",
		"The following will be configured:",
		"",
		m.getSelectedSummary(),
		"",
		fmt.Sprintf("Display Manager: %s", accentStyle.Render(m.selectedDM)),
		"",
		boxStyle.Render("Existing configurations will be backed up to\n~/.local/share/roxyos/backups/"),
		"",
		infoStyle.Render("Press Enter to install"),
		mutedStyle.Render("Press q to cancel"),
	)
}

func (m model) viewInstall() string {
	status := "Copying configurations to ~/.config..."

	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("Installing RoxyOS"),
		"",
		fmt.Sprintf("%s %s", m.spinner.View(), status),
		"",
		mutedStyle.Render("Please wait..."),
	)
}

func (m model) viewComplete() string {
	if m.err != nil {
		return lipgloss.JoinVertical(lipgloss.Center,
			errorStyle.Render("Installation Failed"),
			"",
			m.err.Error(),
			"",
			mutedStyle.Render("Press Enter to exit"),
		)
	}

	dmInstructions := ""
	if m.selectedDM == "sddm" {
		dmInstructions = "sudo systemctl enable sddm"
	} else if m.selectedDM == "ly" {
		dmInstructions = "sudo systemctl enable ly"
	}

	nextSteps := []string{
		"Next steps:",
		"",
	}

	if dmInstructions != "" {
		nextSteps = append(nextSteps, fmt.Sprintf("  1. Enable display manager: %s", accentStyle.Render(dmInstructions)))
		nextSteps = append(nextSteps, "  2. Reboot your system")
		nextSteps = append(nextSteps, "  3. Select 'Niri' from your display manager")
	} else {
		nextSteps = append(nextSteps, "  1. Log out of your current session")
		nextSteps = append(nextSteps, "  2. Start Niri manually or via your preferred method")
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		successStyle.Render(" Installation Complete"),
		"",
		"RoxyOS has been configured successfully!",
		"",
		boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, nextSteps...)),
		"",
		"Keybindings:",
		fmt.Sprintf("  %s  Terminal", accentStyle.Render("Super+T")),
		fmt.Sprintf("  %s  Launcher", accentStyle.Render("Super+D")),
		fmt.Sprintf("  %s  Lock", accentStyle.Render("Super+L")),
		"",
		mutedStyle.Render("Press Enter to exit"),
	)
}

func (m model) getSelectedSummary() string {
	var selected []string
	for _, comp := range m.components {
		if comp.enabled {
			selected = append(selected, successStyle.Render("• ")+comp.name)
		}
	}
	return strings.Join(selected, "\n")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if os.Geteuid() == 0 {
		fmt.Println("Please run as a regular user, not root.")
		fmt.Println("The installer will prompt for sudo when needed.")
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
