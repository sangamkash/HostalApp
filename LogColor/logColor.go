package LogColor

const (
	colorReset = "\033[0m"

	// Basic 8-color ANSI
	colorBlack   = "\033[30m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"

	// Bright versions
	colorBrightBlack   = "\033[90m"
	colorBrightRed     = "\033[91m"
	colorBrightGreen   = "\033[92m"
	colorBrightYellow  = "\033[93m"
	colorBrightBlue    = "\033[94m"
	colorBrightMagenta = "\033[95m"
	colorBrightCyan    = "\033[96m"
	colorBrightWhite   = "\033[97m"

	// 256-color support (foreground)
	colorOrange     = "\033[38;5;208m"
	colorPink       = "\033[38;5;205m"
	colorViolet     = "\033[38;5;177m"
	colorLightGray  = "\033[38;5;250m"
	colorDarkGray   = "\033[38;5;240m"
	colorLightGreen = "\033[38;5;120m"
	colorTeal       = "\033[38;5;37m"
	colorGold       = "\033[38;5;178m"
	colorSkyBlue    = "\033[38;5;117m"
	colorPeach      = "\033[38;5;216m"
)

func Black(s string) string   { return colorBlack + s + colorReset }
func Red(s string) string     { return colorRed + s + colorReset }
func Green(s string) string   { return colorGreen + s + colorReset }
func Yellow(s string) string  { return colorYellow + s + colorReset }
func Blue(s string) string    { return colorBlue + s + colorReset }
func Magenta(s string) string { return colorMagenta + s + colorReset }
func Cyan(s string) string    { return colorCyan + s + colorReset }
func White(s string) string   { return colorWhite + s + colorReset }

func BrightBlack(s string) string   { return colorBrightBlack + s + colorReset }
func BrightRed(s string) string     { return colorBrightRed + s + colorReset }
func BrightGreen(s string) string   { return colorBrightGreen + s + colorReset }
func BrightYellow(s string) string  { return colorBrightYellow + s + colorReset }
func BrightBlue(s string) string    { return colorBrightBlue + s + colorReset }
func BrightMagenta(s string) string { return colorBrightMagenta + s + colorReset }
func BrightCyan(s string) string    { return colorBrightCyan + s + colorReset }
func BrightWhite(s string) string   { return colorBrightWhite + s + colorReset }

func Orange(s string) string     { return colorOrange + s + colorReset }
func Pink(s string) string       { return colorPink + s + colorReset }
func Violet(s string) string     { return colorViolet + s + colorReset }
func LightGray(s string) string  { return colorLightGray + s + colorReset }
func DarkGray(s string) string   { return colorDarkGray + s + colorReset }
func LightGreen(s string) string { return colorLightGreen + s + colorReset }
func Teal(s string) string       { return colorTeal + s + colorReset }
func Gold(s string) string       { return colorGold + s + colorReset }
func SkyBlue(s string) string    { return colorSkyBlue + s + colorReset }
func Peach(s string) string      { return colorPeach + s + colorReset }
