package bot

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

type Context struct {
	engine *Engine
	update tgbotapi.Update

	index    int8
	handlers HandlersChain
}

var (
	abortIndex int8 = math.MaxInt8 >> 1
)

/************************************/
/************ FLOW CONTROL **********/
/************************************/

func (ctx *Context) Next() {
	ctx.index++
	for ctx.index < int8(len(ctx.handlers)) {
		ctx.handlers[ctx.index](ctx)
		ctx.index++
	}
}

func (ctx *Context) IsAborted() bool {
	return ctx.index >= abortIndex
}

func (ctx *Context) Abort() {
	ctx.index = abortIndex
}

// AbortWith panic function
func (ctx *Context) AbortWith(f interface{}, args ...interface{}) {
	v := reflect.ValueOf(f)

	var callArgs []reflect.Value
	for _, arg := range args {
		callArgs = append(callArgs, reflect.ValueOf(arg))
	}

	result := v.Call(callArgs)
	for _, res := range result {
		if err, ok := res.Interface().(error); ok && err != nil { // nolint
			panic(err)
		}
	}

	ctx.Abort()
}

/************************************/
/********** CONTEXT HANDLERS ********/
/************************************/

func (ctx *Context) NameMainHandler() string {
	mainHandler := ctx.handlers[len(ctx.handlers)-1]
	value := reflect.ValueOf(mainHandler)

	return runtime.FuncForPC(value.Pointer()).Name()
}

func (ctx *Context) CurrentNameHandler() string {
	lenHandlers := int8(len(ctx.handlers))
	if ctx.index < 0 || ctx.index >= lenHandlers {
		return "DONE"
	}

	handler := ctx.handlers[ctx.index]
	value := reflect.ValueOf(handler)

	return runtime.FuncForPC(value.Pointer()).Name()
}

/************************************/
/*************** DATA ***************/
/************************************/

func (ctx *Context) From() *tgbotapi.User {
	return ctx.update.SentFrom()
}

func (ctx *Context) Chat() *tgbotapi.Chat {
	return ctx.update.FromChat()
}

func (ctx *Context) Method() string {
	if ctx.update.Message != nil {
		message := ctx.update.Message
		if message.IsCommand() {
			return "command"
		} else {
			return "message"
		}
	} else if ctx.update.EditedMessage != nil {
		if ctx.update.EditedMessage.IsCommand() {
			return "edited_command"
		} else {
			return "edited_message"
		}
	} else if ctx.update.ChannelPost != nil {
		return "channel_post"
	} else if ctx.update.EditedChannelPost != nil {
		return "edited_channel_post"
	} else if ctx.update.InlineQuery != nil {
		return "inline_query"
	} else if ctx.update.ChosenInlineResult != nil {
		return "chosen_inline_result"
	} else if ctx.update.CallbackQuery != nil {
		return "callback"
	} else if ctx.update.ShippingQuery != nil {
		return "shipping_query"
	} else if ctx.update.PreCheckoutQuery != nil {
		return "pre_checkout_query"
	} else if ctx.update.Poll != nil {
		return "poll"
	} else if ctx.update.PollAnswer != nil {
		return "poll_answer"
	} else if ctx.update.MyChatMember != nil {
		return "my_chat_member"
	} else if ctx.update.ChatMember != nil {
		return "chat_member"
	} else if ctx.update.ChatJoinRequest != nil {
		return "chat_join_request"
	}

	return ""
}

func (ctx *Context) Query() string {
	if ctx.update.Message != nil {
		message := ctx.update.Message
		if message.IsCommand() {
			return message.Command()
		} else {
			return message.Text
		}
	} else if ctx.update.EditedMessage != nil {
		if ctx.update.EditedMessage.IsCommand() {
			return ctx.update.EditedMessage.Command()
		} else {
			return ctx.update.EditedMessage.Text
		}
	} else if ctx.update.ChannelPost != nil {
		return ctx.update.ChannelPost.Text
	} else if ctx.update.EditedChannelPost != nil {
		return ctx.update.EditedChannelPost.Text
	} else if ctx.update.InlineQuery != nil {
		return ctx.update.InlineQuery.Query
	} else if ctx.update.ChosenInlineResult != nil {
		return ctx.update.ChosenInlineResult.Query
	} else if ctx.update.CallbackQuery != nil {
		return ctx.GetCallbackTemplate()
	} else if ctx.update.ShippingQuery != nil {
		return ctx.update.ShippingQuery.InvoicePayload
	} else if ctx.update.PreCheckoutQuery != nil {
		return ctx.update.PreCheckoutQuery.InvoicePayload
	} else if ctx.update.Poll != nil {
		return ctx.update.Poll.ID
	} else if ctx.update.PollAnswer != nil {
		return ctx.update.PollAnswer.PollID
	} else if ctx.update.MyChatMember != nil {
		return "*"
	} else if ctx.update.ChatMember != nil {
		return "*"
	} else if ctx.update.ChatJoinRequest != nil {
		return "*"
	}

	return ""
}

func (ctx *Context) GetMessage() *tgbotapi.Message {
	if ctx.update.Message != nil {
		return ctx.update.Message
	} else if ctx.update.EditedMessage != nil {
		return ctx.update.EditedMessage
	} else if ctx.update.CallbackQuery != nil {
		return ctx.update.CallbackQuery.Message
	} else if ctx.update.ChannelPost != nil {
		return ctx.update.ChannelPost
	} else if ctx.update.EditedChannelPost != nil {
		return ctx.update.EditedChannelPost
	}

	return nil
}

func (ctx *Context) GetCallback() *tgbotapi.CallbackQuery {
	if ctx.update.CallbackQuery != nil {
		return ctx.update.CallbackQuery
	}

	return nil
}

func (ctx *Context) GetCallbackID() string {
	data := ctx.GetCallback()
	if data != nil {
		return data.ID
	}

	return ""
}

func (ctx *Context) GetUpdateID() int {
	return ctx.update.UpdateID
}

/************************************/
/************** RESPONSE ************/
/************************************/

func (ctx *Context) Message(text string) error {
	message, err := ctx.buildMessage(text)
	if err != nil {
		return err
	}

	_, err = ctx.engine.botAPI.Send(message)
	return err
}

func (ctx *Context) MessageWithKeyboard(text string, keyboard interface{}) error {
	message, err := ctx.buildMessage(text)
	if err != nil {
		return err
	}
	message.ReplyMarkup = keyboard

	_, err = ctx.engine.botAPI.Send(message)
	return err
}

func (ctx *Context) Answer(text string) error {
	message, err := ctx.buildMessage(text)
	if err != nil {
		return err
	}
	message.ReplyToMessageID = ctx.GetMessage().MessageID

	_, err = ctx.engine.botAPI.Send(message)
	return err
}

func (ctx *Context) AnswerWithKeyboard(text string, keyboard interface{}) error {
	message, err := ctx.buildMessage(text)
	if err != nil {
		return err
	}

	message.ReplyMarkup = keyboard
	message.ReplyToMessageID = ctx.GetMessage().MessageID

	_, err = ctx.engine.botAPI.Send(message)
	return err
}

func (ctx *Context) Callback(showAlert bool, text string) error {
	cb := tgbotapi.NewCallback(ctx.GetCallbackID(), text)
	cb.ShowAlert = showAlert

	_, err := ctx.engine.botAPI.Request(cb)
	return err
}

/************************************/
/************** CALLBACK ************/
/************************************/

func (ctx *Context) CallbackDone() {
	_ = ctx.Callback(false, "")
}

func (ctx *Context) BigCallbackData(template string, data interface{}) string {
	id, err := ctx.engine.callbackStorage.Add(data)
	if err != nil {
		panic(fmt.Errorf("error add in callback storage: %w", err))
	}

	return fmt.Sprintf("bigData:%s:%s", template, id)
}

func (ctx *Context) GetCallbackTemplate() string {
	cbData := ctx.GetCallback().Data
	if strings.HasPrefix(cbData, "bigData") {
		split := strings.Split(cbData, ":")
		if len(split) != 3 {
			return ""
		}

		return split[1]
	}

	return cbData
}

func (ctx *Context) GetCallbackData() (interface{}, error) {
	cbData := ctx.GetCallback().Data
	if strings.HasPrefix(cbData, "bigData") {
		split := strings.Split(cbData, ":")
		if len(split) != 3 {
			return nil, ErrInvalidBigCallbackData
		}

		id, err := uuid.Parse(split[2])
		if err != nil {
			return nil, fmt.Errorf("error parse uuid of data: %w", err)
		}

		return ctx.engine.callbackStorage.Get(id)
	}

	return cbData, nil
}

func (ctx *Context) MustDeleteCallbackData() {
	if err := ctx.DeleteCallbackData(); err != nil {
		panic(err)
	}
}

func (ctx *Context) DeleteCallbackData() error {
	cbData := ctx.GetCallback().Data
	if strings.HasPrefix(cbData, "bigData") {
		split := strings.Split(cbData, ":")
		if len(split) != 3 {
			return ErrInvalidBigCallbackData
		}

		id, err := uuid.Parse(split[2])
		if err != nil {
			return fmt.Errorf("error parse uuid of data: %w", err)
		}

		return ctx.engine.callbackStorage.Delete(id)
	}

	return nil
}

/************************************/
/************ GO CONTEXT ************/
/************************************/

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

func (ctx *Context) Done() <-chan struct{} {
	return nil
}

func (ctx *Context) Err() error {
	return nil
}

func (ctx *Context) Value(_ any) any {
	return nil
}

func (ctx *Context) String() string {
	return "bot.Context"
}

/************************************/
/*************** UTILS **************/
/************************************/

func (ctx *Context) buildMessage(text string) (tgbotapi.MessageConfig, error) {
	chat := ctx.update.FromChat()
	if chat == nil {
		return tgbotapi.MessageConfig{}, ErrInvalidChat
	}

	return tgbotapi.NewMessage(chat.ID, text), nil
}

func (ctx *Context) reset() {
	ctx.index = -1
}
