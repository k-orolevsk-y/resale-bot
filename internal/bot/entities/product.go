package entities

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Product struct {
	ID              uuid.UUID      `db:"id"`
	CategoryID      uuid.UUID      `db:"category_id"`
	Producer        string         `db:"producer"`
	Model           string         `db:"model"`
	Additional      string         `db:"additional"`
	OperatingSystem int            `db:"operating_system"`
	Description     string         `db:"description"`
	Photo           sql.NullString `db:"photo"`
	Price           float64        `db:"price"`
	OldPrice        float64        `db:"old_price"`
	IsSale          bool           `db:"is_sale"`
}

func (p *Product) String() string {
	var texts []string

	texts = append(texts, fmt.Sprintf("<b>%s — %s</b>\n%s", p.Producer, p.Model, p.Additional))

	if p.Price != 0 {
		if p.OldPrice == 0 {
			texts = append(texts, fmt.Sprintf("Цена: <b>%.2f ₽</b>", p.Price))
		} else {
			texts = append(texts, fmt.Sprintf("Цена: <s>%.2f ₽</s> <b>%.2f ₽</b>", p.OldPrice, p.Price))
		}
	}

	if p.Description != "" {
		texts = append(texts, p.Description)
	}

	return strings.Join(texts, "\n\n")
}

func (p *Product) StringWithoutDescription() string {
	var texts []string

	texts = append(texts, fmt.Sprintf("<b>%s — %s</b>\n%s", p.Producer, p.Model, p.Additional))

	if p.Price != 0 {
		if p.OldPrice == 0 {
			texts = append(texts, fmt.Sprintf("Цена: <b>%.2f ₽</b>", p.Price))
		} else {
			texts = append(texts, fmt.Sprintf("Цена: <s>%.2f ₽</s> <b>%.2f ₽</b>", p.OldPrice, p.Price))
		}
	}

	return strings.Join(texts, "\n\n")
}
