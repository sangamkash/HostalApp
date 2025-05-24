package LogColor

const (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorOrange  = "\033[38;5;208m" // 256-color orange
	colorReset   = "\033[0m"
)

func Red(s string) string    { return colorRed + s + colorReset }
func Green(s string) string  { return colorGreen + s + colorReset }
func Yellow(s string) string { return colorYellow + s + colorReset }
func Blue(s string) string   { return colorBlue + s + colorReset }
func Pink(s string) string   { return colorMagenta + s + colorReset }
func Orange(s string) string { return colorOrange + s + colorReset }
