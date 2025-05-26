package LogHelper

const (
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorOrange  = "\033[38;5;208m" // 256-color orange
	colorReset   = "\033[0m"
)

func LogServiceStarting(s string) string {
	return colorYellow + "=====starting ::" + s + " ::starting======" + colorReset
}
func LogServiceStarted(s string) string {
	return colorGreen + "=====started ::" + s + " ::successfully======" + colorReset
}
func LogServiceFailToStarted(s string) string {
	return colorRed + "=====fail to start ::" + s + " :: fail!!!======" + colorReset
}

func LogPanic(s string) string {
	return colorRed + "=====!!!!!Panic ::" + s + " :: Panic!!!!!======" + colorReset
}

func LogValidator(s string) string {
	return colorYellow + "Validator::" + s + colorReset
}
