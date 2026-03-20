package nlp

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"unicode"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/aaaton/golem/v4/dicts/ru"
)

//go:embed resources/stopwords.txt
var stopwordsData []byte

// Processor handles tokenization and lemmatization of text.
type Processor struct {
	ruLemmatizer *golem.Lemmatizer
	enLemmatizer *golem.Lemmatizer
	stopwords    map[string]struct{}
}

// NewProcessor creates a new Processor.
func NewProcessor() (*Processor, error) {
	ruLemmatizer, err := golem.New(ru.New())
	if err != nil {
		return nil, fmt.Errorf("failed to create ru lemmatizer: %w", err)
	}

	enLemmatizer, err := golem.New(en.New())
	if err != nil {
		return nil, fmt.Errorf("failed to create en lemmatizer: %w", err)
	}

	stopwords, err := loadStopwords(stopwordsData)
	if err != nil {
		return nil, fmt.Errorf("failed to load stopwords: %w", err)
	}

	return &Processor{
		ruLemmatizer: ruLemmatizer,
		enLemmatizer: enLemmatizer,
		stopwords:    stopwords,
	}, nil
}

func loadStopwords(data []byte) (map[string]struct{}, error) {
	stopwords := make(map[string]struct{})
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())

		if word != "" {
			stopwords[word] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan stopwords: %w", err)
	}

	return stopwords, nil
}

// Process tokenizes and lemmatizes a document.
func (p *Processor) Process(doc string) []string {
	clearText := strings.Map(func(r rune) rune {
		if unicode.IsSymbol(r) || unicode.IsDigit(r) || unicode.IsPunct(r) {
			return -1
		}
		return r
	}, doc)

	fields := strings.Fields(strings.ToLower(clearText))
	tokens := make([]string, 0, len(fields))

	for _, token := range fields {
		if _, ok := p.stopwords[token]; ok {
			continue
		}

		if p.ruLemmatizer.InDict(token) {
			tokens = append(tokens, p.ruLemmatizer.Lemma(token))
		} else if p.enLemmatizer.InDict(token) {
			tokens = append(tokens, p.enLemmatizer.Lemma(token))
		} else {
			tokens = append(tokens, token)
		}
	}

	return tokens
}
