package db

import (
	md "networkbot-v1/models"
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"time"
)

type User struct {
	UserId		int 	`db:"user_id"`
	Username 	string	`db:"username"`
	StartedAt	time.Time	`db:"started_at"`
}

// Returns contact by tel and error
func GetByTel(tel string, userid int) (md.Contact, error) {
	ctx := context.Background()
	conn := connDB()
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, `SELECT name, user_id, tel, description FROM Contacts WHERE tel=$1 AND user_id=$2`, tel, userid)
	var c md.Contact
	err = pgxscan.ScanOne(&c, rows)
	return c, err
}

// Saves md.Contact c
// Returns userID and error
func Save(c *md.Contact) (int, error) {
	_, err := GetByTel(c.Telephone, c.UserID)
	ctx := context.Background()
	conn := connDB()
	defer conn.Close(ctx)

	if err != nil {
		_, err = conn.Exec(ctx, "INSERT INTO Contacts (user_id, tel, description, name) VALUES ($1, $2, $3, $4)", c.UserID, c.Telephone, c.Description, c.Name)
	} else {
		_, err = conn.Exec(ctx, "UPDATE Contacts SET (description, name) = ($1, $2) WHERE tel=$3", c.Description, c.Name, c.Telephone)
	}
	return c.UserID, err
}

func getUser(userId int) (User, error) {
	ctx := context.Background()
	conn := connDB()
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, `SELECT * FROM users WHERE user_id=$1`, userId)
	var user User
	err = pgxscan.ScanOne(&user, rows)
	return user, err
}

func InsertUser(m *tb.Message) error {
	_, err := getUser(m.Sender.ID)

	if err != nil {
		ctx := context.Background()
		conn := connDB()
		defer conn.Close(ctx)

		_, err = conn.Exec(ctx, "INSERT INTO users VALUES ($1, $2, now())", m.Sender.ID, m.Sender.Username)
	}
	return err
}

func connDB() *pgx.Conn {
	//url := "postgres://pahijagrkuygtu:cb5d57616624bb5b6548e6b3ef202376119a3504b0049158424b19f0aff3dc67@ec2-52-205-3-3.compute-1.amazonaws.com:5432/dbb6mku8d18n1v"
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return conn
}