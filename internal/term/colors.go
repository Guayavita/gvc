package term

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var errorStyle = lipgloss.NewStyle().
	Bold(true).Foreground(
	lipgloss.Color("#B22222"))

var infoStyle = lipgloss.NewStyle().
	Bold(true).Foreground(
	lipgloss.Color("#6464C9"))

func Error(message string, args ...interface{}) {
	fmt.Println(errorStyle.Render(fmt.Sprintf(message, args...)))
}

func Info(message string, args ...interface{}) {
	fmt.Println(infoStyle.Render(fmt.Sprintf(message, args...)))
}
