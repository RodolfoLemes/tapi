package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"tapi/internal/bot"
	"tapi/internal/broker"
)

func main() {
	producer, err := broker.NewProducer[bot.Command]()
	if err != nil {
		log.Fatalf("Tapi: producer initialization err: %s", err.Error())
	}

	err = producer.AttachQueue(broker.CommandRequestQueue)
	if err != nil {
		log.Fatalf("Tapi: producer attach queue err: %s", err.Error())
	}

	consumer, err := broker.NewConsumer[bot.Command]()
	if err != nil {
		log.Fatalf("Tapi: consumer initialization err: %s", err.Error())
	}

	err = consumer.AttachQueue(broker.CommandResponseQueue)
	if err != nil {
		log.Fatalf("Tapi: consumer attach queue err: %s", err.Error())
	}

	bot := bot.New(producer, consumer)

	err = bot.Init()
	if err != nil {
		log.Fatalf("Tapi: bot initialization err: %s", err.Error())
	}

	log.Println("Tapi starts...")

	exit := make(
		chan os.Signal,
		1, // we need to reserve to buffer size 1, so the notifier are not blocked
	)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP)

	select {
	case sig := <-exit:
		// Log the status of this shutdown.
		if sig == syscall.SIGSTOP {
			log.Fatal("Tapi: explode")
		}
	}

	log.Println("Tapi: graceful shutdown")
}
