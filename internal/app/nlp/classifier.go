package nlp

import (
	"math"

	"protomorphine/tg-notes/internal/domain"
)

// Classifier implements a Multinomial Naive Bayes classifier for text documents.
type Classifier struct {
	nlpProcessor   *Processor
	vocabSize      int
	wordCountByCat map[domain.Category]int
	catProbs       map[domain.Category]float64
	freqByCat      map[domain.Category]map[string]int
}

// NewClassifier creates a new Classifier.
func NewClassifier(processor *Processor, dataset []domain.Note) *Classifier {
	c := &Classifier{
		wordCountByCat: make(map[domain.Category]int),
		freqByCat:      make(map[domain.Category]map[string]int),
		catProbs:       make(map[domain.Category]float64),
		nlpProcessor:   processor,
	}

	c.train(dataset)
	return c
}

func (c *Classifier) train(dataset []domain.Note) {
	vocab := make(map[string]struct{})
	docsInCat := make(map[domain.Category]int)
	totalDocs := len(dataset)

	for _, note := range dataset {
		category := note.Category
		if _, ok := c.freqByCat[category]; !ok {
			c.freqByCat[category] = make(map[string]int)
		}

		docsInCat[category]++

		tokens := c.nlpProcessor.Process(note.Content)
		for _, token := range tokens {
			vocab[token] = struct{}{}
			c.wordCountByCat[category]++
			c.freqByCat[category][token]++
		}
	}

	for category, count := range docsInCat {
		c.catProbs[category] = math.Log(float64(count) / float64(totalDocs))
	}

	c.vocabSize = len(vocab)
}

func (c *Classifier) Predict(text string) (map[domain.Category]float64, domain.Category) {
	logPredictions := make(map[domain.Category]float64)
	tokens := c.nlpProcessor.Process(text)

	for category, freqs := range c.freqByCat {
		logProb := c.catProbs[category]

		for _, token := range tokens {
			// P(token|category) = (count(token, category) + 1) / (total words in category + vocab size)
			numerator := float64(1 + freqs[token])
			denominator := float64(c.vocabSize + c.wordCountByCat[category])
			logProb += math.Log(numerator / denominator)
		}

		logPredictions[category] = logProb
	}

	// Convert log probabilities to linear scale and normalize
	predictions := make(map[domain.Category]float64)
	maxLogProb := -math.MaxFloat64
	for _, logProb := range logPredictions {
		if logProb > maxLogProb {
			maxLogProb = logProb
		}
	}

	var sumExp float64
	for category, logProb := range logPredictions {
		// Subtract maxLogProb for numerical stability (softmax trick)
		p := math.Exp(logProb - maxLogProb)
		predictions[category] = p
		sumExp += p
	}

	if sumExp > 0 {
		for category := range predictions {
			predictions[category] /= sumExp
		}
	}

	var bestCategory domain.Category
	maxProb := -1.0

	for cat, prob := range predictions {
		if prob > maxProb {
			maxProb = prob
			bestCategory = cat
		}
	}

	return predictions, bestCategory
}
