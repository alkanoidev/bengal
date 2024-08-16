package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type model struct {
	email string
	form  *huh.Form
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	var cmds []tea.Cmd

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	s := ""
	switch m.form.State {
	case huh.StateCompleted:
		m.email = m.form.GetString("email")
		s += m.email
		return s
	default:
		header := "Tabby"
		form := strings.TrimSuffix(m.form.View(), "\n\n")
		footer := "\nPress q to quit.\n"
		s = header + "\n" + form + "\n\n" + footer
		return s
	}
}

func initModel() model {
	m := model{}
	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What’s your name?").
				Value(&m.email).
				Key("email").
				Validate(func(str string) error {
					if str == "Frank" {
						return errors.New("Sorry, we don’t serve customers named Frank.")
					}
					return nil
				}),
		).WithShowHelp(false),
	)
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
