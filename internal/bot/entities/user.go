package entities

import (
	"fmt"
	"strings"
	"time"
)

type User struct {
	ID           int64     `db:"id"`
	Tag          string    `db:"tag"`
	IsManager    bool      `db:"is_manager"`
	IsBanned     bool      `db:"is_banned"`
	RegisteredAt time.Time `db:"registered_at"`
}

func (u *User) String() string {
	var texts []string

	texts = append(texts, fmt.Sprintf("ID: <b>%d</b>", u.ID))
	texts = append(texts, fmt.Sprintf("Тег: <b>%s</b>", u.Tag))

	if u.IsManager {
		texts = append(texts, "Менеджер: <b>да</b>")
	} else {
		texts = append(texts, "Менеджер: <b>нет</b>")
	}

	if u.IsBanned {
		texts = append(texts, "Заблокирован: <b>да</b>")
	} else {
		texts = append(texts, "Заблокирован: <b>нет</b>")
	}

	texts = append(texts, fmt.Sprintf("Дата регистрации: <b>%s</b>", u.RegisteredAt.Format("02.01.2006 15:04:05")))

	return strings.Join(texts, "\n")
}
