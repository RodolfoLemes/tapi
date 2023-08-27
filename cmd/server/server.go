package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tapi/internal/bot"
	"tapi/internal/broker"
	"tapi/internal/container"
	"tapi/internal/server"
)

func main() {
	consumer, err := broker.NewConsumer[bot.Command]()
	if err != nil {
		log.Fatalf("Server: consumer initialization err: %s", err.Error())
	}

	err = consumer.AttachQueue(broker.CommandRequestQueue)
	if err != nil {
		log.Fatalf("Server: consumer attach queue err: %s", err.Error())
	}

	producer, err := broker.NewProducer[bot.Command]()
	if err != nil {
		log.Fatalf("Server: producer initialization err: %s", err.Error())
	}

	err = producer.AttachQueue(broker.CommandResponseQueue)
	if err != nil {
		log.Fatalf("Server: producer attach queue err: %s", err.Error())
	}

	container, err := container.New()
	if err != nil {
		log.Fatalf("Server: container initialization err: %s", err.Error())
	}

	server := server.New(consumer, producer, container.ScheduleService)

	ctx, cancel := context.WithCancel(context.Background())

	err = server.Init(ctx)
	if err != nil {
		log.Fatalf("Server: bot initialization err: %s", err.Error())
	}

	log.Println("server starts...")

	exit := make(
		chan os.Signal,
		1, // we need to reserve to buffer size 1, so the notifier are not blocked
	)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGSTOP)

	select {
	case sig := <-exit:
		cancel()
		// Log the status of this shutdown.
		if sig == syscall.SIGSTOP {
			log.Fatal("Server: explode")
		} else {
			log.Println("Server: graceful shutdown")
		}
	}
}
