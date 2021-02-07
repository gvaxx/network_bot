package main

import (
	"fmt"
	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"networkbot-v1/db"
	"networkbot-v1/models"
	"os"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func saveContact(contact *models.Contact) {
	_, err := db.Save(contact)
	if err != nil {
		log.Printf("SAVE ERR: %v", err)
	}
}

//func viewContacts(contacts []models.Contact) string {
//	var message string
//	for contact := range contacts {
//		str := fmt.Sprintf("%#v", contact)
//		fmt.Println(str)
//		message += str
//	}
//	return message
//}

func inputDescription(b *tb.Bot, m *tb.Message, selector *tb.ReplyMarkup) {
	b.Send(m.Sender, "Please input description(where your first met, for example)", selector)
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token: os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
		Verbose: true,
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
			inputDescription(b, m, selector)
			return
		}
		c, err := db.GetByTel(m.Contact.PhoneNumber, m.Sender.ID)

		if err != nil {
			contacts[m.Sender.ID] = &models.Contact{
				Name: m.Contact.FirstName + m.Contact.LastName,
				Telephone: m.Contact.PhoneNumber,
				UserID: m.Sender.ID,
			}
			inputDescription(b, m, selector)
		} else {

			b.Send(m.Sender, fmt.Sprintf("Contact already exists:\nName:%v\nDesc:%v", c.Name, c.Description))
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		if _, isset := contacts[m.Sender.ID]; isset {
			(*contacts[m.Sender.ID]).Description = m.Text
			saveContact(contacts[m.Sender.ID])
			delete(contacts, m.Sender.ID);
			b.Send(m.Sender, "Description succesuful save!")
		}
		b.Send(m.Sender, "Send me a contact!")
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

	b.Handle("/start", func(m *tb.Message) {
		if db.InsertUser(m) != nil {
			b.Send(m.Sender, "ERROR: " + err.Error())
		} else {
			b.Send(m.Sender, "WELCOME!")
		}
	} )

	b.Start()
}