package main

import (
	"log"
	// "fmt"
	"time"
	"os"
	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func isNonDescriptionContact(m *tb.Message, contacts map[int64]int64) bool {
	_, ok := contacts[m.Chat.ID]
	return ok
}

func saveInDb(m *tb.Message) int64 {
	return 12345
}

func saveDescription(contact_id int64, description string) {
	
}


func inputDescription(b *tb.Bot, m *tb.Message, selector *tb.ReplyMarkup) {
	b.Send(m.Sender, "Please input description(where your first met, for example)", selector)
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		// You can also set custom API URL.
		// If field is empty it equals to "https://api.telegram.org".
		Token:  os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	var contacts = make(map[int64]int64)

	if err != nil {
		log.Fatal(err)
		return
	}
	var (
		// Universal markup builders.
		menu     = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		selector = &tb.ReplyMarkup{}

		// Reply buttons.
		btnHelp     = menu.Text("ℹ Help")
		btnSettings = menu.Text("⚙ Settings")

		// Inline buttons.
		//
		// Pressing it will cause the client to
		// send the bot a callback.
		//
		// Make sure Unique stays unique as per button kind,
		// as it has to be for callback routing to work.
		//
		btnCancel = selector.Data("cancel", "cancel")
	)

	menu.Reply(
		menu.Row(btnHelp),
		menu.Row(btnSettings),
	)
	selector.Inline(
		selector.Row(btnCancel),
	)

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "Hello World!", selector)
	})

	b.Handle(tb.OnContact, func(m *tb.Message) {
		if isNonDescriptionContact(m, contacts) {
			b.Send(m.Sender, "Please, add description for previous contact or cancel")
			inputDescription(b, m, selector)
			return
		}

		contact_id := saveInDb(m)
		if contact_id > 0 {
			contacts[m.Chat.ID] = contact_id
			inputDescription(b, m, selector)
		} else {
			b.Send(m.Sender, "We have some database error! Please try again later")
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {

		if isNonDescriptionContact(m, contacts) {
			saveDescription(contacts[m.Chat.ID], m.Text)
			delete(contacts, m.Chat.ID);
			b.Send(m.Sender, "Description succesuful save!")
		} else {
			b.Send(m.Sender, "Please, add contact or click /search")
		}
		// all the text messages that weren't
		// captured by existing handlers
	})
	b.Handle(&btnCancel, func(c *tb.Callback) {
		if isNonDescriptionContact(c.Message, contacts) {
			delete(contacts, c.Message.Chat.ID)
			b.Send(c.Sender, "Contact is not saved!")
		}
		b.Respond(c)
	})
	b.Start()
}