package entities

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Repair struct {
	ID           uuid.UUID      `db:"id"`
	ProducerName string         `db:"producer_name"`
	ModelName    string         `db:"model_name"`
	Name         string         `db:"name"`
	Description  sql.NullString `db:"description"`
	Price        float64        `db:"price"`
}

func (r *Repair) String() string {
	var texts []string

	texts = append(texts, fmt.Sprintf("<b>%s %s</b>", r.ProducerName, r.ModelName))
	texts = append(texts, r.Name)

	if r.Description.Valid {
		texts = append(texts, fmt.Sprintf("\n%s\n", r.Description.String))
	}

	texts = append(texts, fmt.Sprintf("Цена: <b>%.2f ₽</b>", r.Price))

	return strings.Join(texts, "\n")
}

func (r *Repair) StringWithoutDescription() string {
	var texts []string

	texts = append(texts, fmt.Sprintf("<b>%s %s</b>", r.ProducerName, r.ModelName))
	texts = append(texts, r.Name)
	texts = append(texts, fmt.Sprintf("Цена: <b>%.2f ₽</b>", r.Price))

	return strings.Join(texts, "\n")
}

func (r *Repair) StringForBot(botURL string) string {
	var texts []string

	texts = append(texts, fmt.Sprintf("Ремонт #%s", strings.Split(r.ID.String(), "-")[0]))
	texts = append(texts, fmt.Sprintf("\t\tПроизводитель: <b>%s</b>", r.ProducerName))
	texts = append(texts, fmt.Sprintf("\t\tМодель: <b>%s</b>", r.ModelName))
	texts = append(texts, fmt.Sprintf("\t\tНазвание ремонта: <b>%s</b>", r.Name))
	texts = append(texts, fmt.Sprintf("\t\tЦена: <b>%.2f ₽</b>", r.Price))

	if r.Description.Valid {
		splitDescription := strings.Split(r.Description.String, "\n")
		description := strings.Join(splitDescription, "\n\t\t")

		texts = append(texts, fmt.Sprintf("\n\t\t%s\n", description))
	}

	if botURL != "" {
		texts = append(texts, fmt.Sprintf("<a href=\"%smrp-%s\">Редактировать</a>", botURL, r.ID.String()))
	}

	return strings.Join(texts, "\n")
}
