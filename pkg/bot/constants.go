package bot

import (
	"fmt"
	"regexp"
)

var (
	regexForAllStrings = regexp.MustCompile(".*")

	ErrInvalidChat            = fmt.Errorf("invalid chat")
	ErrInvalidBigCallbackData = fmt.Errorf("invalid big callback data")
)

func defaultRecovery(ctx *Context, err interface{}) {
	_ = ctx.Message(fmt.Sprintf("Не удалось обработать сообщение. Сообщите эту информацию разработчикам:\n%s", err))
	ctx.Abort()
}
