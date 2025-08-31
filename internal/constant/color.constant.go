package constant

// Color returns the ANSI escape code string for the given color or style name.
// Supported names are:
//   - "reset"  : Reset all attributes
//   - "bold"   : Bold text
//   - "red"    : Red foreground
//   - "green"  : Green foreground
//   - "yellow" : Yellow foreground
//   - "blue"   : Blue foreground
//   - "cyan"   : Cyan foreground
//   - "gray"   : Gray foreground
//
// If the name is not recognized, an empty string is returned.
func Color(name string) string {
	switch name {
	case "reset":
		return "\033[0m"
	case "bold":
		return "\033[1m"
	case "red":
		return "\033[31m"
	case "green":
		return "\033[32m"
	case "yellow":
		return "\033[33m"
	case "blue":
		return "\033[34m"
	case "cyan":
		return "\033[36m"
	case "gray":
		return "\033[90m"
	default:
		return ""
	}
}
