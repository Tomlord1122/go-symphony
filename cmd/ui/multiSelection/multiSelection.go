// Package multiSelection provides functions that
// help define and draw a multi-select step
package multiSelection

import (
	"fmt"

	"github.com/Tomlord1122/go-symphony/v2/cmd/program"
	"github.com/Tomlord1122/go-symphony/v2/cmd/steps"

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

// A Selection represents a choice made in a multiSelection step
type Selection struct {
	Choices map[string]bool
}

// Update changes the value of a Selection's Choice
func (s *Selection) Update(optionName string, value bool) {
	s.Choices[optionName] = value
}

// A multiSelection.model contains the data for the multiSelection step.
//
// It has the required methods that make it a bubbletea.Model
type model struct {
	cursor   int
	options  []steps.Item
	selected map[int]struct{}
	choices  *Selection
	header   string
	exit     *bool
}

func (m model) Init() tea.Cmd {
	return nil
}

// InitialModelMulti initializes a multiSelection step with
// the given data
func InitialModelMultiSelect(options []steps.Item, selection *Selection, header string, program *program.Project) model {
	return model{
		options:  options,
		selected: make(map[int]struct{}),
		choices:  selection,
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
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			// Only confirm if at least one item is selected
			if len(m.selected) > 0 {
				for selectedKey := range m.selected {
					m.choices.Update(m.options[selectedKey].Flag, true)
				}
				return m, tea.Quit
			}
			// Otherwise, do nothing (stay in selection)
			return m, nil
		}
	}
	return m, nil
}

// View is called to draw the multiSelection step
func (m model) View() string {
	s := m.header + "\n\n"

	for i, option := range m.options {
		cursor := " "
		if m.cursor == i {
			cursor = focusedStyle.Render(">")
			option.Title = selectedItemStyle.Render(option.Title)
			option.Desc = selectedItemDescStyle.Render(option.Desc)
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = focusedStyle.Render("*")
		}

		title := focusedStyle.Render(option.Title)
		description := descriptionStyle.Render(option.Desc)

		s += fmt.Sprintf("%s [%s] %s\n%s\n\n", cursor, checked, title, description)
	}

	s += fmt.Sprintf("Press %s to select/deselect, %s to confirm choice.\n", focusedStyle.Render("space"), focusedStyle.Render("enter"))
	return s
}
