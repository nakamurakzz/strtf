package main

import (
	"encoding/base64"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there is been an error: %v", err)
		os.Exit(1)
	}
}

type State int

const (
	StateSelect State = iota
	StateInput
	StateResult
)

type model struct {
	choices       []string
	cursor        int
	selected      int
	s             string
	state         State // 0: select window, 1: input window 2: result window
	input         string
	translateFunc map[int]func(string) string
}

func InitialModel() model {
	return model{
		choices: []string{
			"1. Base64 Encoding",
			"2. Base64 Decoding",
			"3. Byte Encoding",
			"4. Byte Decoding",
		},
		translateFunc: map[int]func(string) string{
			0: func(s string) string {
				//	base64 Encoding
				return base64.URLEncoding.EncodeToString([]byte(s))
			},
			1: func(s string) string {
				//	base64 Decoding
				decoded, err := base64.URLEncoding.DecodeString(s)
				if err != nil {
					return "Error: " + err.Error()
				}
				return string(decoded)
			},
			2: func(s string) string {
				//	byte Encoding
				return fmt.Sprintf("%v", []byte(s))
			},
			3: func(s string) string {
				//	byte Decoding
				return string([]byte(s))
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctr+c":
			return m, tea.Quit
		case "q":
			if m.state == StateSelect {
				return m, tea.Quit
			} else if m.state == StateInput {
				m.input += "q"
				if len(m.input) > 4096 {
					return m, tea.Quit
				}
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "m":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.cursor
			if m.state == StateSelect {
				m.state = StateInput
				m.s = m.choices[m.selected]
			} else if m.state == StateInput {
				m.state = StateResult
				m.s = m.choices[m.selected] + "\n\n"

				m.s += m.translateFunc[m.selected](m.input)
			} else {
				m.state = StateSelect
				m.s = ""
				m.input = ""
			}
		default:
			m.input += msg.String()
			if len(m.input) > 4096 {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case StateSelect:
		s := m.s
		for i, choice := range m.choices {

			cursor := " " // no cursor
			if m.cursor == i {
				cursor = ">" // cursor!
			}

			// Render the row
			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		// The footer
		s += "\nPress ctr+c to quit.\n"

		// Send the UI for rendering
		return s

	case StateInput:
		s := m.s + "\n\n"
		s += "Please Enter the string: " + m.input
		return s
	}
	return m.s
}
