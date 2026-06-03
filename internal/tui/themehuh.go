package tui

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"neocut/internal/theme"
)

func buildHuhTheme(themeName, colorMode string) *huh.Theme {
	rc := theme.ResolveColors(themeName, colorMode)
	t := huh.ThemeBase()

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color(rc.Accent))
	t.Focused.Card = t.Focused.Base
	t.Focused.Title = t.Focused.Title.Foreground(lipgloss.Color(rc.Accent))
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(lipgloss.Color(rc.Accent))
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.Color(rc.Primary))
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(lipgloss.Color(rc.Error))
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(lipgloss.Color(rc.Error))
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(lipgloss.Color(rc.Warning))
	t.Focused.SelectedPrefix = t.Focused.SelectedPrefix.Foreground(lipgloss.Color(rc.Success))
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(lipgloss.Color(rc.Primary)).Background(lipgloss.Color(rc.Accent))
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(lipgloss.Color(rc.Primary)).Background(lipgloss.Color(rc.Success))

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(lipgloss.Color(rc.Accent))
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lipgloss.Color(rc.Primary))
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(lipgloss.Color(rc.Warning))

	t.Blurred = t.Focused
	t.Blurred.Base = t.Blurred.Base.BorderStyle(lipgloss.HiddenBorder())
	t.Blurred.Card = t.Blurred.Base
	t.Blurred.NextIndicator = lipgloss.NewStyle()
	t.Blurred.PrevIndicator = lipgloss.NewStyle()

	t.Group.Title = t.Focused.Title
	t.Group.Description = t.Focused.Description

	return t
}
