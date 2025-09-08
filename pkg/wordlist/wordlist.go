package wordlist

import "strings"

// Wordlist manages target-specific words
type Wordlist struct {
    companies     []string
    abbreviations []string
    ssids         []string
}

// New creates a new Wordlist instance
func New() *Wordlist {
    return &Wordlist{
        companies:     make([]string, 0),
        abbreviations: make([]string, 0),
        ssids:         make([]string, 0),
    }
}

// AddCompany adds a company name to the wordlist
func (w *Wordlist) AddCompany(company string) {
    if company = strings.TrimSpace(company); company != "" {
        w.companies = append(w.companies, company)
    }
}

// AddAbbreviation adds an abbreviation to the wordlist
func (w *Wordlist) AddAbbreviation(abbr string) {
    if abbr = strings.TrimSpace(abbr); abbr != "" {
        w.abbreviations = append(w.abbreviations, abbr)
    }
}

// AddSSID adds an SSID to the wordlist
func (w *Wordlist) AddSSID(ssid string) {
    if ssid = strings.TrimSpace(ssid); ssid != "" {
        w.ssids = append(w.ssids, ssid)
    }
}

// GetAllWords returns all words combined
func (w *Wordlist) GetAllWords() []string {
    var allWords []string
    allWords = append(allWords, w.companies...)
    allWords = append(allWords, w.abbreviations...)
    allWords = append(allWords, w.ssids...)
    return allWords
}

// GetCompanies returns company names
func (w *Wordlist) GetCompanies() []string {
    return w.companies
}

// GetAbbreviations returns abbreviations
func (w *Wordlist) GetAbbreviations() []string {
    return w.abbreviations
}

// GetSSIDs returns SSIDs
func (w *Wordlist) GetSSIDs() []string {
    return w.ssids
}