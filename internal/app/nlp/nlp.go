/*
Package nlp provides natural language processing capabilities for text classification.
It includes a Processor for text tokenization and lemmatization, and a Classifier
that implements a Multinomial Naive Bayes algorithm to categorize text.

The Processor takes raw text, cleans it by removing symbols and digits, splits it into tokens,
removes common stopwords for English and Russian, and reduces words to their base or root form (lemmatization).

The Classifier is trained on a dataset of labeled documents (notes). Once trained, it can
predict the category of a new piece of text. It uses logarithmic probabilities for numerical stability.

Usage:

	// Create a new processor.
	processor, err := nlp.NewProcessor()
	if err != nil {
		// handle error
	}

	// Create a new classifier.
	classifier := nlp.NewClassifier(processor)

	// Train the classifier on a dataset of notes.
	classifier.Train(trainingData)

	// Predict the category of a new text.
	_, category := classifier.Predict("This is a new text to classify.")
*/
package nlp
