package theme

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type Role int

const (
	RolePrimary Role = iota
	RoleSuccess
	RoleWarning
	RoleError
	RoleAccent
)

type Theme struct {
	Name   string
	Label  string
	Colors []string
}

func (t Theme) Hex(role Role) string {
	idx := int(role)
	if idx >= len(t.Colors) {
		idx %= len(t.Colors)
	}
	return t.Colors[idx]
}

const defaultThemeName = "sunny_beach_day"

var defaultTheme = Theme{
	Name: defaultThemeName, Label: "Sunny Beach Day (Default)",
	Colors: []string{"#264653", "#2a9d8f", "#e9c46a", "#f4a261", "#e76f51"},
}

var Themes = []Theme{
	{
		Name: "dark", Label: "Dark Theme",
		Colors: []string{"#0f172a", "#111827", "#1e293b", "#334155", "#64748b", "#e2e8f0"},
	},
	{
		Name: "light", Label: "Light Theme",
		Colors: []string{"#ffffff", "#f8fafc", "#e2e8f0", "#cbd5e1", "#475569", "#0f172a"},
	},
	defaultTheme,
	{
		Name: "olive_garden_feast", Label: "Olive Garden Feast",
		Colors: []string{"#606c38", "#283618", "#fefae0", "#dda15e", "#bc6c25"},
	},
	{
		Name: "summer_ocean_breeze", Label: "Summer Ocean Breeze",
		Colors: []string{"#e63946", "#f1faee", "#a8dadc", "#457b9d", "#1d3557"},
	},
	{
		Name: "refreshing_summer_fun", Label: "Refreshing Summer Fun",
		Colors: []string{"#8ecae6", "#219ebc", "#023047", "#ffb703", "#fb8500"},
	},
	{
		Name: "black_gold_elegance", Label: "Black & Gold Elegance",
		Colors: []string{"#000000", "#14213d", "#fca311", "#e5e5e5", "#ffffff"},
	},
	{
		Name: "vibrant_color_fiesta", Label: "Vibrant Color Fiesta",
		Colors: []string{"#ffbe0b", "#fb5607", "#ff006e", "#8338ec", "#3a86ff"},
	},
	{
		Name: "light_steel", Label: "Light Steel",
		Colors: []string{"#f8f9fa", "#e9ecef", "#dee2e6", "#ced4da", "#adb5bd", "#6c757d", "#495057", "#343a40", "#212529"},
	},
	{
		Name: "golden_twilight", Label: "Golden Twilight",
		Colors: []string{"#000814", "#001d3d", "#003566", "#ffc300", "#ffd60a"},
	},
	{
		Name: "deep_sea", Label: "Deep Sea",
		Colors: []string{"#0d1b2a", "#1b263b", "#415a77", "#778da9", "#e0e1dd"},
	},
	{
		Name: "bright_green", Label: "Bright Green",
		Colors: []string{"#004b23", "#006400", "#007200", "#008000", "#38b000", "#70e000", "#9ef01a", "#ccff33"},
	},
	{
		Name: "vivid_nightfall", Label: "Vivid Nightfall",
		Colors: []string{"#10002b", "#240046", "#3c096c", "#5a189a", "#7b2cbf", "#9d4edd", "#c77dff", "#e0aaff"},
	},
}

var activeTheme Theme
var activeMu sync.RWMutex

func SetActive(name, colorMode string) {
	activeMu.Lock()
	defer activeMu.Unlock()
	activeTheme = Resolve(name, colorMode)
}

func Active() Theme {
	activeMu.RLock()
	defer activeMu.RUnlock()
	return activeTheme
}

func RoleColor(role Role) string {
	return Active().Hex(role)
}

func Find(name string) (Theme, bool) {
	for _, t := range Themes {
		if strings.EqualFold(t.Name, name) {
			return t, true
		}
	}
	return Theme{}, false
}

func Names() []string {
	names := make([]string, len(Themes))
	for i, t := range Themes {
		names[i] = t.Name
	}
	return names
}

func Labels() []string {
	labels := make([]string, len(Themes))
	for i, t := range Themes {
		labels[i] = t.Label
	}
	return labels
}

func DefaultTheme() Theme {
	return defaultTheme
}

type RoleColors struct {
	Primary string
	Success string
	Warning string
	Error   string
	Accent  string
}

func Resolve(name, colorMode string) Theme {
	switch colorMode {
	case "dark":
		if t, ok := Find("dark"); ok {
			return t
		}
	case "light":
		if t, ok := Find("light"); ok {
			return t
		}
	}
	t, ok := Find(name)
	if !ok {
		return defaultTheme
	}
	return t
}

func ResolveColors(name, colorMode string) RoleColors {
	t := Resolve(name, colorMode)
	return RoleColors{
		Primary: t.Hex(RolePrimary),
		Success: t.Hex(RoleSuccess),
		Warning: t.Hex(RoleWarning),
		Error:   t.Hex(RoleError),
		Accent:  t.Hex(RoleAccent),
	}
}

func (t Theme) IsDark() bool {
	r, g, b := parseHex(t.Hex(RolePrimary))
	lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
	return lum < 128
}

func parseHex(hex string) (uint8, uint8, uint8) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 255, 255, 255
	}
	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)
	return uint8(r), uint8(g), uint8(b)
}

func Sprintf(text, hex string) string {
	r, g, b := parseHex(hex)
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}

func SprintfBold(text, hex string) string {
	r, g, b := parseHex(hex)
	return fmt.Sprintf("\x1b[1;38;2;%d;%d;%dm%s\x1b[0m", r, g, b, text)
}
