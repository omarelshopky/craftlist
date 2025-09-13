package ui

type Colors struct {
	Reset  string
	Red    string
	Yellow string
	Green  string
	Cyan   string
	Bold   string
}

var DefaultColors = Colors{
	Reset:  "\033[0m",
	Red:    "\033[31m",
	Yellow: "\033[33m",
	Green:  "\033[32m",
	Cyan:   "\033[36m",
	Bold:   "\033[1m",
}