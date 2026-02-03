package terminal

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

type theme byte

const (
	DarkTheme theme = iota
	LightTheme
)

type colorStyle byte

const (
	PrimaryColor colorStyle = iota
	SecondaryColor
	ClockColor
	SizeColor
	ErrorColor
	ErrorLabelColor
	WarningColor
	WarningLabelColor
	SuccessColor
	SuccessLabelColor
	InfoLevelColor
	DebugLevelColor
	WarnLevelColor
	ErrorLevelColor
	TableTitleColor
	SpinnerColor
)

var colorMap = map[theme]map[colorStyle]lipgloss.Style{
	DarkTheme: {
		PrimaryColor:   lipgloss.NewStyle().Foreground(lipgloss.Color("#ADB3BF")),
		SecondaryColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6370")),
		ClockColor:     lipgloss.NewStyle().Foreground(lipgloss.Color("#61AFEF")),
		SizeColor:      lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C07B")),
		ErrorColor:     lipgloss.NewStyle().Foreground(lipgloss.Color("#E06C75")),
		ErrorLabelColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#E06C75")),
		WarningColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C07B")),
		WarningLabelColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282C34")).
			Background(lipgloss.Color("#E5C07B")),
		SuccessColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#98C379")),
		SuccessLabelColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282C34")).
			Background(lipgloss.Color("#98C379")),
		InfoLevelColor:  lipgloss.NewStyle().Foreground(lipgloss.Color("#56B6C2")),
		DebugLevelColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#C678DD")),
		WarnLevelColor:  lipgloss.NewStyle().Foreground(lipgloss.Color("#D19A66")),
		ErrorLevelColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#E06C75")),
		TableTitleColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282C34")).
			Background(lipgloss.Color("#61AFEF")).Bold(true),
		SpinnerColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#E5C07B")),
	},
	LightTheme: {
		PrimaryColor:   lipgloss.NewStyle().Foreground(lipgloss.Color("#373941")),
		SecondaryColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A1A7")),
		ClockColor:     lipgloss.NewStyle().Foreground(lipgloss.Color("#4078F2")),
		SizeColor:      lipgloss.NewStyle().Foreground(lipgloss.Color("#C18401")),
		ErrorColor:     lipgloss.NewStyle().Foreground(lipgloss.Color("#E45649")),
		ErrorLabelColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#E45649")),
		WarningColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#C18401")),
		WarningLabelColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#C18401")),
		SuccessColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#50A14F")),
		SuccessLabelColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#50A14F")),
		InfoLevelColor:  lipgloss.NewStyle().Foreground(lipgloss.Color("#0184BC")),
		DebugLevelColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#A626A4")),
		WarnLevelColor:  lipgloss.NewStyle().Foreground(lipgloss.Color("#986801")),
		ErrorLevelColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#E45649")),
		TableTitleColor: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#4078F2")).Bold(true),
		SpinnerColor: lipgloss.NewStyle().Foreground(lipgloss.Color("#C18401")),
	},
}

func getColorStyle(cs colorStyle) lipgloss.Style {
	theme := DarkTheme

	if !lipgloss.HasDarkBackground() {
		theme = LightTheme
	}

	return colorMap[theme][cs]
}

func HasColor() bool {
	return !nocolor && termenv.ColorProfile() != termenv.Ascii
}

func SetNoColor(forceNoColor bool) {
	nocolor = forceNoColor
}
