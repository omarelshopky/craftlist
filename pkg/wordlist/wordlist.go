package wordlist

import "strings"

type Wordlist struct {
	words []string
	ssids []string
}

func New() *Wordlist {
	return &Wordlist{
		words: make([]string, 0),
		ssids: make([]string, 0),
	}
}

func (w *Wordlist) AddWord(word string) {
	if word = strings.TrimSpace(word); word != "" {
		w.words = append(w.words, word)
	}
}

func (w *Wordlist) AddSSID(ssid string) {
	if ssid = strings.TrimSpace(ssid); ssid != "" {
		w.ssids = append(w.ssids, ssid)
	}
}

func (w *Wordlist) GetAllWords() []string {
	var allWords []string

	allWords = append(allWords, w.words...)
	allWords = append(allWords, w.ssids...)

	return allWords
}

func (w *Wordlist) GetWords() []string {
	return w.words
}

func (w *Wordlist) GetSSIDs() []string {
	return w.ssids
}
