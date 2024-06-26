package repos

import (
	"database/sql"
	"fmt"

	"github.com/alogsDiu/Oqiga/internal/models"
)

type Event struct {
	DB *sql.DB
}

func (e *Event) Get_all_parties(parites *[]models.Party) error {
	query := "select * from Party"

	rows, err := e.DB.Query(query)

	if err != nil {
		fmt.Println("The query is not being able to be executed", err)
		return err
	}
	for rows.Next() {
		var name string
		var city string
		var organizer string
		var date string
		var about string
		var id string

		err := rows.Scan(&name, &city, &organizer, &date, &about, &id)

		if err != nil {
			fmt.Println("Not able to parce one row")
		}

		party := models.Party{
			Name:      name,
			City:      city,
			Organizer: organizer,
			Date:      date,
			About:     about,
			Id:        id,
		}

		*parites = append(*parites, party)
	}
	return nil
}
func (e *Event) Get_confirmed_parties(parites *[]models.Party, user_name *string) {
	query := "select party_id from confirmed_parties where user_name=?"
	rows, err := e.DB.Query(query, *user_name)
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var name string
		var city string
		var organizer string
		var date string
		var about string
		var id string

		err := rows.Scan(&id)

		if err != nil {
			fmt.Println(err)
		}

		query = "select * from party where id=?"
		row := e.DB.QueryRow(query, id)

		err = row.Scan(&id, &name, &about, &city, &organizer, &date)

		if err != nil {
			fmt.Println(err)
		}

		party := models.Party{
			Name:      name,
			City:      city,
			Organizer: organizer,
			Date:      date,
			About:     about,
			Id:        id,
		}

		*parites = append(*parites, party)
	}
}
func (e *Event) Get_single_party(id *string, name *string, about *string, city *string, organizer *string, date *string) {
	query := "select * from party where id =?"

	row := e.DB.QueryRow(query, *id)

	err := row.Scan(id, name, about, city, organizer, date)

	if err != nil {
		fmt.Println(err)
	}
}
func (e *Event) Send_message(user_name, party_id, whose, body *string) {
	query := "insert into messages(user_name, whose, party_id, body) values (?,?,?,?)"
	_, err := e.DB.Exec(query, *user_name, *whose, *party_id, body)
	if err != nil {
		fmt.Println(err)
	}
}
func (e *Event) Get_messages(id, user_name *string, messages *[]models.Message) {
	query := "select body, whose from messages where party_id=? and user_name=?"
	rows, err := e.DB.Query(query, *id, *user_name)

	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var body string
		var whose string
		err := rows.Scan(&body, &whose)
		if err != nil {
			fmt.Println(err)
		}
		*messages = append(*messages, models.Message{Text: body, Whose: whose})
	}
}
func (e *Event) Get_parties_organized_by_user(user_name *string, parties *[]models.Party) {

	query := "select name, id from Party where organizer=?"

	rows, err := e.DB.Query(query, *user_name)

	if err != nil {
		fmt.Println(err)
		return
	}

	for rows.Next() {
		var name string
		var id string
		rows.Scan(&name, &id)
		*parties = append(*parties, models.Party{Name: name, Id: id})
	}
}
