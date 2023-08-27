package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"tapi/config"
	"tapi/internal/broker"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	tlgBot *tgbotapi.BotAPI

	producer broker.Producer[Command]
	consumer broker.Consumer[Command]
}

func New(producer broker.Producer[Command], consumer broker.Consumer[Command]) *Bot {
	return &Bot{
		producer: producer,
		consumer: consumer,
	}
}

func (b *Bot) Init() error {
	token := config.Env.Telegram.BotToken

	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	b.tlgBot = botApi
	b.tlgBot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := b.tlgBot.GetUpdatesChan(u)

	go b.listenerUpdates(updates)

	msgs, err := b.consumer.Consume(context.TODO())
	if err != nil {
		return err
	}

	go b.handleBrokerMessages(msgs)

	return err
}

func (b *Bot) listenerUpdates(ups tgbotapi.UpdatesChannel) {
	for up := range ups {
		cmd, err := NewCommand(
			up.Message.Command(),
			up.Message.CommandArguments(),
			up.Message.From.ID,
			up.Message.Chat.ID,
			up.Message.From.UserName,
		)
		if err != nil {
			b.handleErrorMessage(up.Message.Chat.ID, up.Message.Command(), err)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = b.producer.Publish(ctx, cmd)
		if err != nil {
			b.handleErrorMessage(up.Message.Chat.ID, up.Message.Command(), err)
			continue
		}
	}
}

func (b *Bot) handleErrorMessage(chatID int64, cmd string, err error) {
	text := fmt.Sprintf("Invalid command %s, description: %s", cmd, err.Error())
	b.sendMessage(chatID, text)
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)

	if _, err := b.tlgBot.Send(msg); err != nil {
		log.Println("Telegram Bot: sendMessage error: ", err)
	}
}

func (b *Bot) handleBrokerMessages(msgs <-chan broker.Message[Command]) {
	for msg := range msgs {
		b.handleBrokerMessage(msg)
	}
}

func (b *Bot) handleBrokerMessage(msg broker.Message[Command]) {
	b.sendMessage(msg.Data.TelegramChatID, msg.Data.Text)
}
