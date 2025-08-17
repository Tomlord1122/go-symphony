// Package singleSelection provides functions that
// help define and draw a multi-input step
package singleSelection

import (
	"fmt"

	"github.com/Tomlord1122/go-symphony/cmd/program"
	"github.com/Tomlord1122/go-symphony/cmd/steps"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Claude CLI color scheme
var (
	focusedStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("#FACC15")).Bold(true)                                                        // Yellow for focused items
	titleStyle            = lipgloss.NewStyle().Background(lipgloss.Color("#FACC15")).Foreground(lipgloss.Color("#000000")).Bold(true).Padding(0, 1, 0) // Yellow background
	selectedItemStyle     = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#34D399")).Bold(true)                                         // Cyan green for selected
	selectedItemDescStyle = lipgloss.NewStyle().PaddingLeft(1).Foreground(lipgloss.Color("#34D399"))                                                    // Cyan green for selected desc
	descriptionStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))                                                                   // Light gray for descriptions
)

// A Selection represents a choice made in a singleSelection step
type Selection struct {
	Choice string
}

// Update changes the value of a Selection's Choice
func (s *Selection) Update(value string) {
	s.Choice = value
}

// A singleSelection.model contains the data for the singleSelection step.
//
// It has the required methods that make it a bubbletea.Model
type model struct {
	cursor   int
	choices  []steps.Item
	selected map[int]struct{}
	choice   *Selection
	header   string
	exit     *bool
}

func (m model) Init() tea.Cmd {
	return nil
}

// InitialModelMulti initializes a singleSelection step with
// the given data
func InitialModelMulti(choices []steps.Item, selection *Selection, header string, program *program.Project) model {
	return model{
		choices:  choices,
		selected: make(map[int]struct{}),
		choice:   selection,
		header:   titleStyle.Render(header),
		exit:     &program.Exit,
	}
}

// Update is called when "things happen", it checks for
// important keystrokes to signal when to quit, change selection,
// and confirm the selection.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			*m.exit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.selected) == 1 {
				// If exactly one item is selected, confirm the selection
				for selectedKey := range m.selected {
					m.choice.Update(m.choices[selectedKey].Title)
					m.cursor = selectedKey
				}
				return m, tea.Quit
			} else if len(m.selected) == 0 {
				// If no item is selected, select the current cursor position and confirm
				m.choice.Update(m.choices[m.cursor].Title)
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// View is called to draw the singleSelection step
func (m model) View() string {
	s := m.header + "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = focusedStyle.Render(">")
			choice.Title = selectedItemStyle.Render(choice.Title)
			choice.Desc = selectedItemDescStyle.Render(choice.Desc)
		}

		checked := " "
		// In single-select mode, show the current cursor position as selected
		// OR if the item is explicitly selected
		_, isSelected := m.selected[i]
		if m.cursor == i || isSelected {
			checked = focusedStyle.Render("x")
		}

		title := focusedStyle.Render(choice.Title)
		description := descriptionStyle.Render(choice.Desc)

		s += fmt.Sprintf("%s [%s] %s\n%s\n\n", cursor, checked, title, description)
	}

	s += fmt.Sprintf("Press %s to confirm choice.\n\n", focusedStyle.Render("enter"))
	return s
}
