package config

import "time"

func NewDefaultGeneratorConfig() GeneratorConfig {
	return GeneratorConfig{
		MinYear:        1990,
		MaxYear:        time.Now().Year(),
		MinPasswordLen: 8,
		MaxPasswordLen: 64,
		CommonWords:    getDefaultCommonWords(),
		Separators:     getDefaultSeparators(),
		Substitutions:  getDefaultSubstitutions(),
		NumberPatterns: getDefaultNumberPatterns(),
		Patterns:       getDefaultPatterns(),
	}
}

func NewDefaultOutputConfig() OutputConfig {
	return OutputConfig{
		Filename: "passwords.txt",
	}
}

func getDefaultCommonWords() []string {
	return []string{
		"password", "admin", "guest", "wifi", "wireless", "IT", "tech",
		"pass", "login", "user", "root", "default", "access",
		"network", "internet", "secure", "temp", "test",
	}
}

func getDefaultSeparators() []string {
	return []string{"", "@", "_", "-", ".", "#", "!", "*", "+", "~", "%", "&", "^"}
}

func getDefaultSubstitutions() map[string][]string {
	return map[string][]string{
		"a": {"4", "@", "^"},
		"A": {"4", "@", "^"},
		"b": {"6"},
		"B": {"8"},
		"c": {"<", "("},
		"C": {"<", "("},
		"D": {")"},
		"e": {"3"},
		"E": {"3"},
		"g": {"9", "6", "&"},
		"G": {"9", "6", "&"},
		"h": {"#"},
		"H": {"#"},
		"i": {"1", "!", "|"},
		"I": {"1", "!", "|"},
		"l": {"1", "|", "7", "2"},
		"L": {"1", "|", "7", "2"},
		"o": {"0"},
		"O": {"0"},
		"p": {"9"},
		"P": {"9"},
		"q": {"9", "2", "&"},
		"Q": {"9", "2", "&"},
		"s": {"5", "$"},
		"S": {"5", "$"},
		"t": {"7", "+"},
		"T": {"7", "+"},
		"z": {"2"},
		"Z": {"2"},
	}
}

func getDefaultNumberPatterns() []string {
	return []string{
		"d", "dd", "ddd", "dddd", "12345", "123456",
	}
}

func getDefaultPatterns() []string {
	return []string{
		"<CUSTOM>",
		"<COMMON>",
		"<SSID>",
		"<CUSTOM><SEP><YEAR>",
		"<CUSTOM><SEP><SHORTYEAR>",
		"<CUSTOM><SEP><NUM>",
		"<CUSTOM><SEP><COMMON>",
		"<COMMON><SEP><YEAR>",
		"<COMMON><SEP><SHORTYEAR>",
		"<COMMON><SEP><NUM>",
		"<COMMON><SEP><CUSTOM>",
		"<SSID><SEP><YEAR>",
		"<SSID><SEP><SHORTYEAR>",
		"<SSID><SEP><NUM>",
		"<SSID><SEP><CUSTOM>",
		"<YEAR><SEP><CUSTOM>",
		"<YEAR><SEP><COMMON>",
		"<YEAR><SEP><SSID>",
		"<SHORTYEAR><SEP><CUSTOM>",
		"<SHORTYEAR><SEP><COMMON>",
		"<SHORTYEAR><SEP><SSID>",
		"<NUM><SEP><CUSTOM>",
		"<NUM><SEP><COMMON>",
		"<NUM><SEP><SSID>",
		"<SEP><CUSTOM><SEP><SSID><SEP><YEAR>",
		"<COMMON><SEP><CUSTOM><YEAR>",
	}
}