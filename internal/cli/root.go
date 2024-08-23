package cli

import (
	"email-checker/internal/pkg"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type model struct {
	email    string
	loading  bool
	form     *huh.Form
	spinner  spinner.Model
	response pkg.EmailData
	table    table.Model
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.loading {
		// TODO: spinner not spinnin
		time.Sleep(2 * time.Second)
		m.email = m.form.GetString("email")
		emailData, err := pkg.CheckDomain(m.email)
		if err != nil {
			// TODO: show error
			fmt.Println("Error: ", err)
		}
		m.response = emailData
		m.loading = false
		return m, tea.Quit
	}

	if m.form.State == huh.StateCompleted {
		m.loading = true
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := ""
	switch m.form.State {
	case huh.StateCompleted:
		if m.loading {
			s += fmt.Sprintf("\n\n   %s Loading...\n\n", m.spinner.View())
			return s
		}
		if m.response != (pkg.EmailData{}) {
			m.table.SetRows([]table.Row{
				{strconv.FormatBool(m.response.HasMX), m.response.SpfRecord, m.response.DmarcRecord},
			})
			s += "\n" + m.table.View()

			return s
		}
		return s
	default:
		header := "Email Checker"
		form := strings.TrimSuffix(m.form.View(), "\n\n")
		footer := "\nPress q to quit.\n"
		s = header + "\n" + form + "\n\n" + footer
		return s
	}
}

func initModel() model {
	m := model{}

	s := spinner.New()
	s.Spinner = spinner.Dot
	m.spinner = s

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter email:").
				Value(&m.email).
				Key("email").
				Validate(func(str string) error {
					// TODO: email regex with error messages
					return nil
				}),
		).WithShowHelp(false),
	)

	columns := []table.Column{
		{Title: "MX", Width: 4},
		{Title: "SPF", Width: 30},
		{Title: "DMARC", Width: 30},
	}

	t := table.New(table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(7))
	m.table = t

	return m
}

var rootCmd = &cobra.Command{
	Use: "monstera",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
