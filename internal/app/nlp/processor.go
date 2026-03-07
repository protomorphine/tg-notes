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

//go:embed resources/stopwords_en.txt
var enStopwordsData []byte

//go:embed resources/stopwords_ru.txt
var ruStopwordsData []byte

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

	stopwords, err := loadStopwords(enStopwordsData, ruStopwordsData)
	if err != nil {
		return nil, fmt.Errorf("failed to load stopwords: %w", err)
	}

	return &Processor{
		ruLemmatizer: ruLemmatizer,
		enLemmatizer: enLemmatizer,
		stopwords:    stopwords,
	}, nil
}

func loadStopwords(datas ...[]byte) (map[string]struct{}, error) {
	stopwords := make(map[string]struct{})
	for _, data := range datas {
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
	}

	return stopwords, nil
}

// Process tokenizes and lemmatizes a document.
func (p *Processor) Process(doc string) []string {
	tokens := p.tokenize(doc)
	return p.lemmatize(tokens)
}

// tokenize tokenizes a document and removes stopwords.
func (p *Processor) tokenize(doc string) []string {
	text := strings.Map(func(r rune) rune {
		if unicode.IsSymbol(r) || unicode.IsDigit(r) {
			return -1
		}
		return r
	}, string(doc))

	fields := strings.Fields(strings.ToLower(text))

	tokens := make([]string, 0, len(fields))
	for _, field := range fields {
		if _, ok := p.stopwords[field]; ok {
			continue
		}
		tokens = append(tokens, field)
	}

	return tokens
}

// lemmatize lemmatizes a list of tokens.
func (p *Processor) lemmatize(tokens []string) []string {
	lemmas := make([]string, 0, len(tokens))

	for _, token := range tokens {
		if p.ruLemmatizer.InDict(token) {
			lemmas = append(lemmas, p.ruLemmatizer.Lemma(token))
		} else if p.enLemmatizer.InDict(token) {
			lemmas = append(lemmas, p.enLemmatizer.Lemma(token))
		} else {
			lemmas = append(lemmas, token)
		}
	}
	return lemmas
}
