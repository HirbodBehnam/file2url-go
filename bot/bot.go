package bot

import (
	"context"
	"file2url/config"
	"file2url/database"
	"file2url/shared"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	"log"
	"net/url"
	"time"
)

// RunBot runs the bot to receive the updates
func RunBot(_ context.Context, client *telegram.Client) error {
	sender := message.NewSender(tg.NewClient(client))
	shared.Dispatcher.OnNewMessage(func(ctx context.Context, entities tg.Entities, u *tg.UpdateNewMessage) error {
		m, ok := u.Message.(*tg.Message)
		if !ok || m.Out {
			// Outgoing message, not interesting.
			return nil
		}
		// Check if the user is allowed to use the bot
		if !checkAllowedUser(m.PeerID) {
			return nil
		}
		// Check file
		replyText := "Please send a media to bot"
		if m.Media != nil {
			doc, ok := m.Media.(*tg.MessageMediaDocument)
			if ok {
				if doc, ok := doc.Document.AsNotEmpty(); ok {
					var filename string
					for _, attribute := range doc.Attributes {
						if name, ok := attribute.(*tg.DocumentAttributeFilename); ok {
							filename = name.FileName
							break
						}
					}
					id, err := shared.Database.Store(database.File{
						FileReference: doc.FileReference,
						Name:          filename,
						ID:            doc.ID,
						AccessHash:    doc.AccessHash,
						Size:          int64(doc.Size),
						MimeType:      doc.MimeType,
						AddedTime:     time.Now(),
					})
					if err != nil {
						replyText = "cannot insert data in database"
						log.Println("cannot insert data in database: ", err)
					} else {
						replyText = config.Config.URLPrefix + "/" + id + "/" + url.PathEscape(filename)
					}
				}
			}
		}

		// Send the link or error
		_, err := sender.Reply(entities, u).NoWebpage().Text(ctx, replyText)
		return err
	})
	return nil
}

// checkAllowedUser checks if a user is allowed to use the bot or not
func checkAllowedUser(peer tg.PeerClass) bool {
	p, ok := peer.(*tg.PeerUser)
	return ok && config.IsUserAllowed(p.UserID)
}
