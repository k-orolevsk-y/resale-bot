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

func (p *Product) StringForBot(botURL string) string {
	var texts []string

	texts = append(texts, fmt.Sprintf("\t\tПроизводитель: <b>%s</b>", p.Producer))
	texts = append(texts, fmt.Sprintf("\t\tМодель: <b>%s</b>", p.Model))
	texts = append(texts, fmt.Sprintf("\t\tАтрибуты: <b>%s</b>", p.Additional))
	texts = append(texts, fmt.Sprintf("\t\tКатегория: <a href=\"%smc-%s\">#%s</a>", botURL, p.CategoryID, strings.Split(p.CategoryID.String(), "-")[0]))

	if p.Price != 0 {
		if p.OldPrice == 0 {
			texts = append(texts, fmt.Sprintf("\n\t\tЦена: <b>%.2f ₽</b>", p.Price))
		} else {
			texts = append(texts, fmt.Sprintf("\n\t\tЦена: <s>%.2f ₽</s> <b>%.2f ₽</b>", p.OldPrice, p.Price))
		}
	}

	if p.IsSale {
		texts = append(texts, "\t\tТовар по акции: <b>да</b>")
	} else {
		texts = append(texts, "\t\tТовар по акции: <b>нет</b>")
	}

	if p.Description != "" {
		splitDescription := strings.Split(p.Description, "\n")
		description := strings.Join(splitDescription, "\n\t\t")

		texts = append(texts, fmt.Sprintf("\n\t\t%s", description))
	}

	texts = append(texts, fmt.Sprintf("<a href=\"%smp-%s\">Редактировать</a>", botURL, p.ID))
	return strings.Join(texts, "\n")
}

func (p *Product) StringForManager(botURL string) string {
	var texts []string

	texts = append(texts, fmt.Sprintf("\t\tПроизводитель: <b>%s</b>", p.Producer))
	texts = append(texts, fmt.Sprintf("\t\tМодель: <b>%s</b>", p.Model))
	texts = append(texts, fmt.Sprintf("\t\tАтрибуты: <b>%s</b>", p.Additional))
	texts = append(texts, fmt.Sprintf("\t\tКатегория: <a href=\"%smc-%s\">#%s</a>", botURL, p.CategoryID, strings.Split(p.CategoryID.String(), "-")[0]))

	if p.Price != 0 {
		if p.OldPrice == 0 {
			texts = append(texts, fmt.Sprintf("\n\t\tЦена: <b>%.2f ₽</b>", p.Price))
		} else {
			texts = append(texts, fmt.Sprintf("\n\t\tЦена: <s>%.2f ₽</s> <b>%.2f ₽</b>", p.OldPrice, p.Price))
		}
	}

	if p.IsSale {
		texts = append(texts, "\t\tТовар по акции: <b>да</b>")
	} else {
		texts = append(texts, "\t\tТовар по акции: <b>нет</b>")
	}

	if p.Description != "" {
		splitDescription := strings.Split(p.Description, "\n")
		description := strings.Join(splitDescription, "\n\t\t")

		texts = append(texts, fmt.Sprintf("\n\t\t%s", description))
	}

	return strings.Join(texts, "\n")
}
