package interfaces

type Placeholder struct {
	Format      string `mapstructure:"format" json:"format"`
	Description string `mapstructure:"description" json:"description"`
}

type PlaceholdersConfig struct {
	CustomWord Placeholder `mapstructure:"custom_word" json:"custom_word"`
	CommonWord Placeholder `mapstructure:"common_word" json:"common_word"`
	SSID       Placeholder `mapstructure:"ssid" json:"ssid"`
	Separator  Placeholder `mapstructure:"separator" json:"separator"`
	Year       Placeholder `mapstructure:"year" json:"year"`
	ShortYear  Placeholder `mapstructure:"short_year" json:"short_year"`
	Number     Placeholder `mapstructure:"number" json:"number"`
}

type Printer interface {
	Info(message string)
	Success(message string)
	Error(message string)
	Warning(message string)
	Bold(message string)
	PrintIntro(version string)
	PrintPlaceholders(placeholders PlaceholdersConfig)
	PrintLoadedWords(category string, count int)
	PrintProgress(count int)
	PrintFinalCount(count int)
	PrintOutputFile(path string)
}