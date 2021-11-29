package main

import (
	"context"
	"file2url/api"
	"file2url/bot"
	"file2url/config"
	"file2url/shared"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load config
	if len(os.Args) > 1 {
		config.LoadConfig(os.Args[1])
	} else {
		config.LoadConfig("config.json")
	}
	// Load the bot
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	stopChannel := make(chan struct{})
	go api.StartWebserver(stopChannel)
	err := telegram.BotFromEnvironment(ctx,
		telegram.Options{SessionStorage: &session.FileStorage{Path: "session"}, UpdateHandler: shared.Dispatcher},
		bot.RunBot,
		telegram.RunUntilCanceled)
	stopChannel <- struct{}{}
	if err != nil {
		log.Fatalln(err)
	}
}
