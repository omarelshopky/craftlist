package ui

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/omarelshopky/craftlist/internal/interfaces"
	"golang.org/x/text/language"
    "golang.org/x/text/message"
)

type Printer struct {
	colors Colors
	humanizer *message.Printer
}

func NewPrinter() *Printer {
	return &Printer{
		colors: DefaultColors,
		humanizer: message.NewPrinter(language.English),
	}
}

func (p *Printer) Info(message string) {
	fmt.Printf("%s%s%s\n", p.colors.Cyan, message, p.colors.Reset)
}

func (p *Printer) Success(message string) {
	fmt.Printf("%s%s%s\n", p.colors.Green, message, p.colors.Reset)
}

func (p *Printer) Error(message string) {
	fmt.Printf("%s%s%s\n", p.colors.Red, message, p.colors.Reset)
}

func (p *Printer) Warning(message string) {
	fmt.Printf("%s%s%s\n", p.colors.Yellow, message, p.colors.Reset)
}

func (p *Printer) Bold(message string) {
	fmt.Printf("%s%s%s\n", p.colors.Bold, message, p.colors.Reset)
}

func (p *Printer) PrintIntro(version string) {
	fmt.Printf(`
                 __ _   _ _     _   
                / _| | | (_)   | |  
  ___ _ __ __ _| |_| |_| |_ ___| |_ 
 / __| '__/ _' |  _| __| | / __| __|
| (__| | | (_| | | | |_| | \__ \ |_ 
 \___|_|  \__,_|_|  \__|_|_|___/\__|
                                      
v%s             By Omar Elshopky

`, version)
}

func (p *Printer) PrintPlaceholders(placeholders interfaces.PlaceholdersConfig) {
	fmt.Printf("%sAvailable Placeholders:%s\n\n", p.colors.Bold, p.colors.Reset)
	fmt.Printf("%s%-15s %s%s\n", p.colors.Green, "PLACEHOLDER", "DESCRIPTION", p.colors.Reset)
	fmt.Printf("%s%-15s %s%s\n", p.colors.Green, strings.Repeat("-", 15), strings.Repeat("-", 50), p.colors.Reset)

	values := reflect.ValueOf(placeholders)
	for idx := 0; idx < values.NumField(); idx++ {
		if placeholder, ok := values.Field(idx).Interface().(interfaces.Placeholder); ok {
			fmt.Printf("%s%-15s %s%s\n", p.colors.Yellow, placeholder.Format, p.colors.Reset, placeholder.Description)
		}
	}
}

func (p *Printer) PrintLoadedWords(category string, count int) {
	fmt.Printf("%sLoaded %s%d%s%s words for %s%s\n",
		p.colors.Cyan, p.colors.Bold, count, p.colors.Reset, p.colors.Cyan, category, p.colors.Reset)
}