package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type themeNoBg struct{}

var _ fyne.Theme = (*themeNoBg)(nil)

func (t themeNoBg) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameButton:
		return color.Transparent
	}

	return theme.DefaultTheme().Color(name, variant)
}

func (m themeNoBg) Icon(name fyne.ThemeIconName) fyne.Resource {

	return theme.DefaultTheme().Icon(name)
}

func (m themeNoBg) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m themeNoBg) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
