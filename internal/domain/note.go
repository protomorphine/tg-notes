// Package domain contains domain types.
package domain

type Category string

// Note struct represent a note.
type Note struct {
	Title    string
	Content  string
	Category Category
}

