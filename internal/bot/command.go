package bot

import (
	"errors"
)

type CommandName string

var (
	AddSchedule  string = "addschedule"
	GetSchedules string = "getschedules"
)

type Command struct {
	Name             CommandName
	Text             string
	TelegramUserID   int64
	TelegramChatID   int64
	TelegramUsername string
}

func NewCommand(
	name string,
	text string,
	telegramUserID int64,
	telegramChatID int64,
	telegramUsername string,
) (Command, error) {
	cmdName, err := newCommandName(name)
	if err != nil {
		return Command{}, err
	}

	return Command{cmdName, text, telegramUserID, telegramChatID, telegramUsername}, nil
}

func newCommandName(name string) (CommandName, error) {
	switch name {
	case AddSchedule, GetSchedules:
		return CommandName(name), nil
	default:
		return CommandName(""), errors.New("invalid command name")
	}
}
