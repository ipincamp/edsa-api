package constant

// ColorName adalah tipe untuk nama warna ANSI
type ColorName string

const (
	ColorReset  ColorName = "reset"
	ColorBold   ColorName = "bold"
	ColorRed    ColorName = "red"
	ColorGreen  ColorName = "green"
	ColorYellow ColorName = "yellow"
	ColorBlue   ColorName = "blue"
	ColorCyan   ColorName = "cyan"
	ColorGray   ColorName = "gray"
)

// Color mengembalikan ANSI escape code string untuk nama warna/style tertentu.
// Jika nama tidak dikenali, return string kosong.
func Color(name ColorName) string {
	switch name {
	case ColorReset:
		return "\033[0m"
	case ColorBold:
		return "\033[1m"
	case ColorRed:
		return "\033[31m"
	case ColorGreen:
		return "\033[32m"
	case ColorYellow:
		return "\033[33m"
	case ColorBlue:
		return "\033[34m"
	case ColorCyan:
		return "\033[36m"
	case ColorGray:
		return "\033[90m"
	default:
		return ""
	}
}
