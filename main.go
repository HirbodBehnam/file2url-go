package main

import (
	"context"
	"file2url/api"
	"file2url/bot"
	"file2url/config"
	"file2url/database"
	"file2url/shared"
	"fmt"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	// Load config
	var err error
	config.LoadConfig()
	shared.Database, err = database.LoadDatabaseFromEnv()
	if err != nil {
		log.Fatalln("cannot open the database: ", err)
	}
	defer func() {
		if err := shared.Database.Close(); err != nil {
			log.Println("cannot close the database: ", err)
		}
	}()
	// Load the bot
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	stopChannel := make(chan struct{})
	go api.StartWebserver(stopChannel)
	fmt.Println("file2url-go version " + config.Version)
	err = telegram.BotFromEnvironment(ctx,
		telegram.Options{SessionStorage: &session.FileStorage{Path: "session"}, UpdateHandler: shared.Dispatcher},
		bot.RunBot,
		telegram.RunUntilCanceled)
	stopChannel <- struct{}{}
	if err != nil {
		log.Fatalln(err)
	}
	<-stopChannel // wait until the server is down
}
