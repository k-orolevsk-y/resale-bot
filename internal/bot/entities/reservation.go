package entities

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	ID        uuid.UUID `db:"id"`
	UserID    int64     `db:"user_id"`
	ProductID uuid.UUID `db:"product_id"`
	CreatedAt time.Time `db:"created_at"`
	Completed int       `db:"completed"`
}

type ReservationWithAdditionalData struct {
	Reservation
	CategoryName    string `db:"category_name"`
	ProductFullName string `db:"product_full_name"`
}

func (r *ReservationWithAdditionalData) StringForBot(botURL string) string {
	var texts []string

	texts = append(texts, fmt.Sprintf("Бронь #%s:", strings.Split(r.ID.String(), "-")[0]))
	texts = append(texts, fmt.Sprintf("\t\tПользователь: <a href=\"%smu-%d\">%d</a>", botURL, r.UserID, r.UserID))
	texts = append(texts, fmt.Sprintf("\t\tКатегория товара: <b>%s</b>", r.CategoryName))
	texts = append(texts, fmt.Sprintf("\t\tТовар: <a href=\"%smp-%s\">%s</a>", botURL, r.ProductID, r.ProductFullName))

	if r.Completed == -1 {
		texts = append(texts, "\t\tСтатус: <b>Отменен</b>")
	} else if r.Completed == 1 {
		texts = append(texts, "\t\tСтатус: <b>Выполнен</b>")
	} else {
		texts = append(texts, "\t\tСтатус: <b>Рассматривается</b>")
	}

	texts = append(texts, fmt.Sprintf("\t\tДата: <b>%s</b>", r.CreatedAt.Format("02.01.2006 15:04:05")))
	texts = append(texts, fmt.Sprintf("<a href=\"%smr-%s\">Редактировать</a>", botURL, r.ID.String()))

	return strings.Join(texts, "\n")
}

func (r *ReservationWithAdditionalData) StringForManager(botURL string) string {
	var texts []string

	texts = append(texts, fmt.Sprintf("Бронь #%s:", strings.Split(r.ID.String(), "-")[0]))
	texts = append(texts, fmt.Sprintf("\t\tПользователь: <a href=\"%smu-%d\">%d</a>", botURL, r.UserID, r.UserID))
	texts = append(texts, fmt.Sprintf("\t\tКатегория товара: <b>%s</b>", r.CategoryName))
	texts = append(texts, fmt.Sprintf("\t\tТовар: <a href=\"%smp-%s\">%s</a>", botURL, r.ProductID, r.ProductFullName))

	if r.Completed == -1 {
		texts = append(texts, "\t\tСтатус: <b>Отменен</b>")
	} else if r.Completed == 1 {
		texts = append(texts, "\t\tСтатус: <b>Выполнен</b>")
	} else {
		texts = append(texts, "\t\tСтатус: <b>Рассматривается</b>")
	}

	texts = append(texts, fmt.Sprintf("\t\tДата: <b>%s</b>", r.CreatedAt.Format("02.01.2006 15:04:05")))

	return strings.Join(texts, "\n")
}
