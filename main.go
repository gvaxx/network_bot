package main

import (
	"log"
	"fmt"
	"time"
	"os"
	"strings"
	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/joho/godotenv"
	models "./models"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func saveContact(contact *models.Contact) {
	allContacts[(*contact).FirstName + (*contact).LastName] = contact
}

func viewContacts(contacts []models.Contact) string {
	var message string
	for contact := range contacts {
		str := fmt.Sprintf("%#v", contact)
		fmt.Println(str)
		message += str
	}
	return message
}

func inputDescription(b *tb.Bot, m *tb.Message, selector *tb.ReplyMarkup) {
	b.Send(m.Sender, "Please input description(where your first met, for example)", selector)
}

var allContacts = make(map[string] *models.Contact)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token: os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	contacts := make(map[int] *models.Contact)
	selector := &tb.ReplyMarkup{}

	btnCancel := selector.Data("cancel", "cancel")
	selector.Inline(
		selector.Row(btnCancel),
	)

	b.Handle(tb.OnContact, func(m *tb.Message) {
		if _, isset := contacts[m.Sender.ID]; isset {
			b.Send(m.Sender, "Please, add description for previous contact or cancel")
			inputDescription(b, m, selector, )
			return
		}

		contacts[m.Sender.ID] = &models.Contact{
			FirstName: m.Contact.FirstName,
			LastName: m.Contact.LastName,
			Telephone: m.Contact.PhoneNumber,
			UserID: m.Contact.UserID,
		}			
		inputDescription(b, m, selector)
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		if _, isset := contacts[m.Sender.ID]; isset {
			(*contacts[m.Sender.ID]).Description = m.Text
			saveContact(contacts[m.Sender.ID])
			delete(contacts, m.Sender.ID);
			b.Send(m.Sender, "Description succesuful save!")
		} else {
			contactViews := []models.Contact{}
			for name, contact := range allContacts {
				// fmt.Println(allContacts)
				fmt.Println(name)
				fmt.Println(m.Text)
				if (strings.Contains(name, m.Text)) {
					contactViews = append(contactViews, *contact)
				}
			}
			b.Send(m.Sender, viewContacts(contactViews))
		}
	})

	b.Handle(&btnCancel, func(c *tb.Callback) {
		if _, isset := contacts[c.Sender.ID]; isset {
			delete(contacts, c.Sender.ID)
			b.Send(c.Sender, "Contact is not saved!")
		}
		b.Respond(c)
	})

	b.Handle("/search", func(m *tb.Message) {
		b.Send(m.Sender, "Please input lastname or its part")
	})

	b.Start()
}