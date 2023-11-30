package middlewares

import (
	"database/sql"
	"errors"

	"go.uber.org/zap"

	"github.com/k-orolevsk-y/resale-bot/internal/bot/entities"
	"github.com/k-orolevsk-y/resale-bot/pkg/bot"
)

func (s *service) Auth(ctx *bot.Context) {
	user := ctx.From()
	if user == nil {
		return
	} else if user.IsBot {
		ctx.Abort()
		return
	}

	u, err := s.rep.GetUserByTgID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u = &entities.User{
				ID:        user.ID,
				Tag:       user.UserName,
				IsManager: false,
			}

			if err = s.rep.CreateUser(ctx, u); err != nil {
				s.logger.Error("error register user in repository", zap.Error(err))
			}
		} else {
			s.logger.Error("error get user for register in repository", zap.Error(err))
		}
	} else {
		if u.Tag != user.UserName {
			u.Tag = user.UserName
			if err = s.rep.EditUser(ctx, u); err != nil {
				s.logger.Error("error change user tag in repository", zap.Error(err))
			}
		}
	}

	ctx.Set("user", u)
}
