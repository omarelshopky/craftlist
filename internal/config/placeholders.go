package config

import "github.com/omarelshopky/craftlist/internal/interfaces"

type PlaceholdersConfig = interfaces.PlaceholdersConfig
type Placeholder = interfaces.Placeholder

func NewDefaultPlaceholdersConfig() PlaceholdersConfig {
	return PlaceholdersConfig{
		CustomWord: Placeholder{
			Format:      "<CUSTOM>",
			Description: "Inserts custom word variations from the file specified with the --words flag",
		},
		CommonWord: Placeholder{
			Format:      "<COMMON>",
			Description: "Inserts common word variations based on the list defined in your config file",
		},
		SSID: Placeholder{
			Format:      "<SSID>",
			Description: "Inserts SSID variations from the file specified with the --ssids flag",
		},
		Separator: Placeholder{
			Format:      "<SEP>",
			Description: "Inserts separators based on the list defined in your config file",
		},
		Year: Placeholder{
			Format:      "<YEAR>",
			Description: "Inserts full year based on the range defined in flags or config file (e.g., 2025)",
		},
		ShortYear: Placeholder{
			Format:      "<SHORTYEAR>",
			Description: "Inserts two-digit year based on the range defined in flags or config file (e.g., 25)",
		},
		Number: Placeholder{
			Format:      "<NUM>",
			Description: "Inserts numbers based on the list defined in your config file",
		},
	}
}