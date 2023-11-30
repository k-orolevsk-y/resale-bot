package entities

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Repair struct {
	ID          uuid.UUID      `db:"id"`
	ModelID     uuid.UUID      `db:"model_id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Price       float64        `db:"price"`
}

type RepairWithModelAndCategory struct {
	Repair
	ModelName    string `db:"model_name"`
	CategoryName string `db:"category_name"`
}

func (r *RepairWithModelAndCategory) String() string {
	var texts []string

	texts = append(texts, fmt.Sprintf("<b>%s %s</b>", r.CategoryName, r.ModelName))
	texts = append(texts, r.Name)

	if r.Description.Valid {
		texts = append(texts, fmt.Sprintf("\n%s\n", r.Description.String))
	}

	texts = append(texts, fmt.Sprintf("Цена: <b>%.2f ₽</b>", r.Price))

	return strings.Join(texts, "\n")
}

func (r *RepairWithModelAndCategory) StringWithoutDescription() string {
	var texts []string

	texts = append(texts, fmt.Sprintf("<b>%s %s</b>", r.CategoryName, r.ModelName))
	texts = append(texts, r.Name)
	texts = append(texts, fmt.Sprintf("Цена: <b>%.2f ₽</b>", r.Price))

	return strings.Join(texts, "\n")
}
