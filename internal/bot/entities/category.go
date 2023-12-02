package entities

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Category struct {
	ID   uuid.UUID `db:"id"`
	Name string    `db:"name"`
	Type int       `db:"c_type"`
}

func (c *Category) StringForBot(botURL string) string {
	var texts []string

	texts = append(texts, fmt.Sprintf("Категория #%s:", strings.Split(c.ID.String(), "-")[0]))
	texts = append(texts, fmt.Sprintf("\t\tНазвание: <b>%s</b>", c.Name))

	if c.Type == 0 {
		texts = append(texts, "\t\tТип: <b>Новые</b>")
	} else {
		texts = append(texts, "\t\tТип: <b>Б/У</b>")
	}

	texts = append(texts, fmt.Sprintf("<a href=\"%smc-%s\">Редактировать</a>", botURL, c.ID.String()))

	return strings.Join(texts, "\n")
}

func (c *Category) StringForManager() string {
	var texts []string

	texts = append(texts, fmt.Sprintf("Категория #%s:", strings.Split(c.ID.String(), "-")[0]))
	texts = append(texts, fmt.Sprintf("\t\tНазвание: <b>%s</b>", c.Name))

	if c.Type == 0 {
		texts = append(texts, "\t\tТип: <b>Новые</b>")
	} else {
		texts = append(texts, "\t\tТип: <b>Б/У</b>")
	}

	return strings.Join(texts, "\n")
}
