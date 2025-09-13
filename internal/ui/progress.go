package ui

import "fmt"

func (p *Printer) PrintProgress(count int) {
	fmt.Printf("\rGenerated %s%s%s unique passwords...", p.colors.Bold, p.humanizeNumber(count), p.colors.Reset)
}

func (p *Printer) PrintFinalCount(count int) {
	fmt.Printf("\n\n%sGenerated %s%s%s%s total unique passwords%s\n",
		p.colors.Green, p.colors.Bold, p.humanizeNumber(count), p.colors.Reset, p.colors.Green, p.colors.Reset)
}

func (p *Printer) PrintOutputFile(path string) {
	p.Success(fmt.Sprintf("Output saved to: %s%s%s\n", p.colors.Bold, path, p.colors.Reset))
}

func (p *Printer) PrintApproximateCount(count int) {
	fmt.Printf("\n%sApproximately %s%s%s%s unique passwords will be generated.%s\n",
		p.colors.Cyan, p.colors.Bold, p.humanizeNumber(count), p.colors.Reset, p.colors.Cyan, p.colors.Reset)
}