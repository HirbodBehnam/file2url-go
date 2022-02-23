package bot

import (
	"context"
	"errors"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"sync"
)

// sessionPool contains some sessions for client in order to do not create new sessions
var sessionPool = sync.Pool{
	New: func() interface{} {
		return new(session.StorageMemory)
	},
}

// GetNewApi creates a new Telegram instance API
func GetNewApi(ctx context.Context) (*tg.Client, error) {
	clientChannel := make(chan *tg.Client) // Send the client here
	go func() {
		_ = telegram.BotFromEnvironment(ctx, telegram.Options{SessionStorage: sessionPool.Get().(*session.StorageMemory)},
			func(ctx context.Context, client *telegram.Client) error { return nil },
			func(ctx context.Context, client *telegram.Client) error {
				clientChannel <- client.API()
				<-ctx.Done()
				return nil
			})
		close(clientChannel) // If we reach here without sending anything in clientChannel, we close it, so we can detect it
	}()
	client, ok := <-clientChannel
	if !ok { // channel closed; error
		return nil, errors.New("cannot start client")
	}
	return client, nil
}
