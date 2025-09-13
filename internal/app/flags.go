package app

import "time"

type Flags struct {
	CfgFile          string
	WordsFile        string
	SSIDsFile        string
	OutputFile       string
	MinLength        int
	MaxLength        int
	MinYear          int
	MaxYear          int
	ListPlaceholders bool
	CountPasswords   bool
}

func NewFlags() *Flags {
	return &Flags{
		OutputFile: "passwords.txt",
		MinLength:  8,
		MaxLength:  64,
		MinYear:    1990,
		MaxYear:    time.Now().Year(),
	}
}