package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	repoURL  = "https://repo.roxyproxy.de/roxyos"
	repoName = "roxyos"
)

var (
	accentColor  = lipgloss.Color("#5dade2")
	mutedColor   = lipgloss.Color("#a8c8d8")
	successColor = lipgloss.Color("#4a6b5a")
	errorColor   = lipgloss.Color("#8b4a2b")

	titleStyle    = lipgloss.NewStyle().Foreground(accentColor).Bold(true).MarginBottom(1)
	subtitleStyle = lipgloss.NewStyle().Foreground(mutedColor).MarginBottom(2)
	boxStyle      = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(accentColor).Padding(1, 2)
	successStyle  = lipgloss.NewStyle().Foreground(successColor)
	errorStyle    = lipgloss.NewStyle().Foreground(errorColor)
	infoStyle     = lipgloss.NewStyle().Foreground(accentColor)
	mutedStyle    = lipgloss.NewStyle().Foreground(mutedColor)
	accentStyle   = lipgloss.NewStyle().Foreground(accentColor)
)

type step int

const (
	stepWelcome step = iota
	stepSelectDM
	stepSelectComponents
	stepConfirm
	stepInstall
	stepApplyConfigs
	stepComplete
)

type component struct {
	name        string
	pkg         string
	description string
	enabled     bool
	isCore      bool
}

type model struct {
	step          step
	width, height int
	spinner       spinner.Model
	dmList        list.Model
	components    []component
	cursor        int
	selectedDM    string
	logs          []string
	err           error
	installStatus string
}

type dmItem struct{ name, desc string }

func (i dmItem) Title() string       { return i.name }
func (i dmItem) Description() string { return i.desc }
func (i dmItem) FilterValue() string { return i.name }

type (
	installLogMsg  struct{ log string }
	installDoneMsg struct{ err error }
)

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(accentColor)

	items := []list.Item{
		dmItem{"ly", "Console-based display manager (lightweight)"},
		dmItem{"sddm", "Graphical display manager with RoxyOS theme"},
		dmItem{"none", "Skip display manager"},
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
			{"Niri", "roxyos-niri", "Scrolling window manager", true, true},
			{"Waybar", "roxyos-waybar", "Status bar", true, true},
			{"Rofi", "roxyos-rofi", "Application launcher", true, true},
			{"Hypr", "roxyos-hypr", "Lock, wallpaper, idle daemon", true, false},
			{"Kitty", "roxyos-kitty", "GPU-accelerated terminal", true, false},
			{"Ghostty", "roxyos-ghostty", "Modern terminal emulator", false, false},
			{"Fish", "roxyos-fish", "Shell with Starship prompt", true, false},
			{"Mako", "roxyos-mako", "Notification daemon", true, false},
			{"Plymouth", "roxyos-plymouth", "Boot splash theme", false, false},
			{"Assets", "roxyos-assets", "Wallpapers and themes", true, false},
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
			if m.step != stepInstall && m.step != stepApplyConfigs {
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
			if m.step == stepSelectComponents && !m.components[m.cursor].isCore {
				m.components[m.cursor].enabled = !m.components[m.cursor].enabled
			}
		case "a":
			if m.step == stepSelectComponents {
				allEnabled := true
				for _, c := range m.components {
					if !c.enabled && !c.isCore {
						allEnabled = false
						break
					}
				}
				for i := range m.components {
					if !m.components[i].isCore {
						m.components[i].enabled = !allEnabled
					}
				}
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case installLogMsg:
		m.installStatus = msg.log
		m.logs = append(m.logs, msg.log)
		if len(m.logs) > 5 {
			m.logs = m.logs[1:]
		}
		return m, nil

	case installDoneMsg:
		if msg.err != nil {
			m.err = msg.err
			m.step = stepComplete
		} else if m.step == stepInstall {
			m.step = stepApplyConfigs
			return m, m.applyConfigs()
		} else {
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
		m.step = stepConfirm
	case stepConfirm:
		m.step = stepInstall
		return m, m.runInstall()
	case stepComplete:
		return m, tea.Quit
	}
	return m, nil
}

func (m model) runInstall() tea.Cmd {
	return func() tea.Msg {
		if err := ensureRepo(); err != nil {
			return installDoneMsg{fmt.Errorf("failed to add repo: %w", err)}
		}

		packages := m.getSelectedPackages()
		if len(packages) == 0 {
			return installDoneMsg{fmt.Errorf("no packages selected")}
		}

		args := append([]string{"-S", "--noconfirm", "--needed"}, packages...)
		cmd := exec.Command("sudo", append([]string{"pacman"}, args...)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if err := cmd.Run(); err != nil {
			return installDoneMsg{fmt.Errorf("pacman failed: %w", err)}
		}

		return installDoneMsg{nil}
	}
}

func (m model) applyConfigs() tea.Cmd {
	return func() tea.Msg {
		home, err := os.UserHomeDir()
		if err != nil {
			return installDoneMsg{err}
		}

		configDir := filepath.Join(home, ".config")
		roxyosDir := "/usr/share/roxyos/configs"

		configMap := map[string]string{
			"roxyos-niri":    "niri",
			"roxyos-waybar":  "waybar",
			"roxyos-rofi":    "rofi",
			"roxyos-kitty":   "kitty",
			"roxyos-ghostty": "ghostty",
			"roxyos-fish":    "fish",
			"roxyos-mako":    "mako",
			"roxyos-hypr":    "hypr",
		}

		for _, comp := range m.components {
			if !comp.enabled {
				continue
			}

			configName, ok := configMap[comp.pkg]
			if !ok {
				continue
			}

			src := filepath.Join(roxyosDir, configName)
			dst := filepath.Join(configDir, configName)

			if _, err := os.Stat(src); err == nil {
				if err := copyDir(src, dst); err != nil {
					return installDoneMsg{fmt.Errorf("failed to copy %s config: %w", configName, err)}
				}
			}
		}

		for _, comp := range m.components {
			if comp.pkg == "roxyos-fish" && comp.enabled {
				src := filepath.Join(roxyosDir, "starship")
				dst := filepath.Join(configDir, "starship")
				if _, err := os.Stat(src); err == nil {
					copyDir(src, dst)
				}
			}
		}

		if m.selectedDM == "sddm" || m.selectedDM == "ly" {
			exec.Command("sudo", "systemctl", "enable", m.selectedDM).Run()
		}

		return installDoneMsg{nil}
	}
}

func ensureRepo() error {
	data, err := os.ReadFile("/etc/pacman.conf")
	if err != nil {
		return err
	}

	if strings.Contains(string(data), "["+repoName+"]") {
		return nil
	}

	repoBlock := fmt.Sprintf("\n[%s]\nSigLevel = Optional TrustAll\nServer = %s\n", repoName, repoURL)

	cmd := exec.Command("sudo", "sh", "-c", fmt.Sprintf("echo '%s' >> /etc/pacman.conf && pacman -Sy", repoBlock))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func (m model) getSelectedPackages() []string {
	var packages []string

	for _, comp := range m.components {
		if comp.enabled {
			packages = append(packages, comp.pkg)
		}
	}

	if m.selectedDM == "sddm" {
		packages = append(packages, "roxyos-sddm")
	} else if m.selectedDM == "ly" {
		packages = append(packages, "roxyos-ly")
	}

	return packages
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

		if strings.HasPrefix(filepath.Base(path), "custom.") {
			if _, err := os.Stat(dstPath); err == nil {
				return nil
			}
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
	case stepConfirm:
		content = m.viewConfirm()
	case stepInstall:
		content = m.viewInstall("Installing packages...")
	case stepApplyConfigs:
		content = m.viewInstall("Applying configurations...")
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
		boxStyle.Render("This wizard will:\n\n  1. Add the RoxyOS repository\n  2. Install selected packages via pacman\n  3. Apply configurations to ~/.config"),
		"",
		infoStyle.Render("Press Enter to begin"),
		mutedStyle.Render("Press q to quit"),
	)
}

func (m model) viewSelectDM() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Display Manager"),
		"",
		m.dmList.View(),
		"",
		mutedStyle.Render("↑/↓ navigate  enter select"),
	)
}

func (m model) viewSelectComponents() string {
	var items []string

	for i, comp := range m.components {
		cursor := "  "
		if i == m.cursor {
			cursor = accentStyle.Render("> ")
		}

		check := "[ ]"
		if comp.enabled {
			check = "[x]"
		}
		checkStyled := check
		if comp.enabled {
			checkStyled = successStyle.Render(check)
		}

		paddedName := fmt.Sprintf("%-12s", comp.name)
		name := paddedName
		if i == m.cursor {
			name = accentStyle.Render(paddedName)
		}

		core := "          "
		if comp.isCore {
			core = mutedStyle.Render("(required)")
		}

		line := fmt.Sprintf("%s%s %s %s  %s", cursor, checkStyled, name, core, mutedStyle.Render(comp.description))
		items = append(items, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render("Select Components"),
		subtitleStyle.Render("Choose which packages to install"),
		"",
		strings.Join(items, "\n"),
		"",
		mutedStyle.Render("j/k navigate | space toggle | a toggle all | enter continue"),
	)
}

func (m model) viewConfirm() string {
	packages := m.getSelectedPackages()

	var pkgList []string
	for _, pkg := range packages {
		pkgList = append(pkgList, successStyle.Render("")+" "+pkg)
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("Confirm Installation"),
		"",
		"The following packages will be installed:",
		"",
		strings.Join(pkgList, "\n"),
		"",
		boxStyle.Render("This will:\n  Add RoxyOS repo to /etc/pacman.conf\n  Install packages via pacman\n  Copy configs to ~/.config"),
		"",
		infoStyle.Render("Press Enter to install"),
		mutedStyle.Render("Press q to cancel"),
	)
}

func (m model) viewInstall(status string) string {
	return lipgloss.JoinVertical(lipgloss.Center,
		titleStyle.Render("Installing RoxyOS"),
		"",
		fmt.Sprintf("%s %s", m.spinner.View(), status),
		"",
		mutedStyle.Render("This may take a moment..."),
	)
}

func (m model) viewComplete() string {
	if m.err != nil {
		return lipgloss.JoinVertical(lipgloss.Center,
			errorStyle.Render(" Installation Failed"),
			"",
			m.err.Error(),
			"",
			mutedStyle.Render("Press Enter to exit"),
		)
	}

	nextSteps := []string{"Next steps:", ""}

	if m.selectedDM != "none" && m.selectedDM != "" {
		nextSteps = append(nextSteps, "  1. Reboot your system")
		nextSteps = append(nextSteps, "  2. Select 'Niri' from the login screen")
	} else {
		nextSteps = append(nextSteps, "  1. Log out and start Niri manually")
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		successStyle.Render(" Installation Complete"),
		"",
		"RoxyOS has been installed and configured!",
		"",
		boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, nextSteps...)),
		"",
		"Keybindings:",
		fmt.Sprintf("  %s  Terminal", accentStyle.Render("Super+T")),
		fmt.Sprintf("  %s  Launcher", accentStyle.Render("Super+D")),
		fmt.Sprintf("  %s  Lock", accentStyle.Render("Super+L")),
		fmt.Sprintf("  %s  Clipboard", accentStyle.Render("Super+P")),
		"",
		mutedStyle.Render("Press Enter to exit"),
	)
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
