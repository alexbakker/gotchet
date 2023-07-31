package tui

import (
	"fmt"
	"io"

	"github.com/alexbakker/gotchet/internal/format"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var (
	styleDoc = lipgloss.NewStyle().Padding(1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4"))
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(1)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(1).PaddingBottom(1)
	quitTextStyle   = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type Model struct {
	c     *format.TestCapture
	ct    *format.Test
	tests []*format.Test

	prevIndex     *int
	width, height int
	list          list.Model
	viewPort      viewport.Model
}

type Item struct {
	Text   string
	Status string
}

func (i Item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(Item)
	if !ok {
		return
	}

	text := item.Text
	if index == m.Index() {
		text = selectedItemStyle.Render(text)
	}

	fmt.Fprint(w, fmt.Sprintf("%s %s", item.Status, text))
}

func New(c *format.TestCapture) *Model {
	m := Model{
		c: c,
	}
	if c.Test != nil {
		// Go deep until there's more than one test
		test := c.Test
		for test.Tests != nil && len(test.Tests) == 1 {
			test = maps.Values(test.Tests)[0]
		}
		m.setCurrentTest(test)
	}
	return &m
}

func (m *Model) setContent(s string) {
	content := lipgloss.NewStyle().Width(m.viewPort.Width).Render(s)
	m.viewPort.SetContent(content)
	m.viewPort.GotoBottom()
}

func (m *Model) setCurrentTest(t *format.Test) {
	m.ct = t
	m.tests = maps.Values(m.ct.Tests)
	slices.SortFunc(m.tests, func(a *format.Test, b *format.Test) int {
		return a.Index - b.Index
	})

	width, height := m.getSize()
	if width == 0 || height == 0 {
		return
	}

	var items []list.Item
	for _, test := range m.tests {
		items = append(items, Item{
			Text:   test.Name(),
			Status: renderTestStatus(test),
		})
	}

	l := list.New(items, itemDelegate{}, width/2, height)
	if t.Parent == nil {
		l.Title = "Tests"
	} else {
		l.Title = t.FullName
	}
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowTitle(true)
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.KeyMap.Quit.Unbind()

	m.list = l

	viewPort := viewport.New(width/2, height)
	m.viewPort = viewPort
	if len(m.tests) > 0 {
		m.setContent(m.tests[0].FullOutput().String())
	} else {
		m.setContent("")
	}
	m.viewPort.KeyMap = viewport.KeyMap{}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) getSize() (int, int) {
	if m.width == 0 || m.height == 0 {
		return 0, 0
	}
	top, right, bottom, left := styleDoc.GetPadding()
	return m.width - left - right, m.height - top - bottom
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	curIndex := m.list.Index()

	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.setCurrentTest(m.ct)
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.list.CursorUp()
		case tea.MouseWheelDown:
			m.list.CursorDown()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			test := m.tests[curIndex]
			if len(test.Tests) > 0 {
				m.prevIndex = &curIndex
				m.setCurrentTest(test)
			}
		case "esc":
			if m.ct.Parent != nil {
				m.setCurrentTest(m.ct.Parent)
				if m.prevIndex != nil {
					m.list.Select(*m.prevIndex)
				}
			}
		}
	}

	m.list, cmd = m.list.Update(msg)
	newIndex := m.list.Index()
	if curIndex != newIndex {
		m.setContent(m.tests[newIndex].FullOutput().String())
	}

	return m, cmd
}

func renderTestStatus(t *format.Test) string {
	if t.Done {
		if t.Passed {
			return lipgloss.NewStyle().SetString("✓").
				Foreground(colorCheck).
				String()
		} else {
			return lipgloss.NewStyle().SetString("x").
				Foreground(colorCross).
				String()
		}
	} else {
		return lipgloss.NewStyle().SetString("?").
			Foreground(colorQuestion).
			String()
	}
}

var (
	colorCheck    = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	colorCross    = lipgloss.AdaptiveColor{Light: "#FF0000", Dark: "#FF0000"}
	colorQuestion = lipgloss.AdaptiveColor{Light: "#FFA500", Dark: "#FFA500"}
)

func (m *Model) View() string {
	m.viewPort.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1)

	doc := lipgloss.JoinHorizontal(lipgloss.Center,
		lipgloss.NewStyle().Width(m.list.Width()).Render(m.list.View()),
		m.viewPort.View(),
	)

	return styleDoc.Render(doc)
}
