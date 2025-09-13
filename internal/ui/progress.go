package ui

import "fmt"

func (p *Printer) PrintProgress(count int) {
	fmt.Printf("\rGenerated %s%d%s unique passwords...", p.colors.Bold, count, p.colors.Reset)
}

func (p *Printer) PrintFinalCount(count int) {
	fmt.Printf("\n\n%sGenerated %s%d%s%s total unique passwords%s\n",
		p.colors.Green, p.colors.Bold, count, p.colors.Reset, p.colors.Green, p.colors.Reset)
}

func (p *Printer) PrintOutputFile(path string) {
	p.Success(fmt.Sprintf("Output saved to: %s%s%s\n", p.colors.Bold, path, p.colors.Reset))
}