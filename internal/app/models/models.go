// Package models contains an application level models.
package models

import "protomorphine/tg-notes/internal/domain"

type SaveResult struct {
	Title    string
	Category domain.Category
}
