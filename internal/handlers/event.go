package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alogsDiu/Oqiga/internal/models"
	"github.com/alogsDiu/Oqiga/internal/renderer"
	"github.com/alogsDiu/Oqiga/internal/services"
)

type Event struct {
	Service *services.Event
	Auth    func(token *string) (string, bool)
}

type Eventprofile struct {
	Name         string
	Date         time.Time
	Hour         int8
	Description  string
	Price        float32
	Requirements string
	Organizer    string
}

type Events struct {
	Parties []models.Party
}

func (e *Event) Gen_recomendations(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")

	if err != nil {
		fmt.Println("Yo ", err)
		return
	}
	_, status := e.Auth(&cookie.Value)

	if status {

		parties := []models.Party{}

		err := e.Service.Get_parties(&parties)

		if err != nil {
			log.Fatalf("not able to get parties: %e", err)
			return
		}

		events := Events{
			Parties: parties,
		}

		renderer.Renderer.Render(w, "recomendations", r, events)
	}
}

func (e *Event) All_recomendations(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")

	if err != nil {
		fmt.Println("Yo ", err)
		return
	}
	_, status := e.Auth(&cookie.Value)

	if status {

		parties := []models.Party{}

		err := e.Service.Get_parties(&parties)

		if err != nil {
			log.Fatalf("not able to get parties: %e", err)
			return
		}

		events := Events{
			Parties: parties,
		}

		renderer.Renderer.Render(w, "parties", r, events)
	}
}

func (e *Event) Confirmed_recomendations(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")

	if err != nil {
		http.Redirect(w, r, "/log_in", http.StatusMovedPermanently)
		return
	}

	user_name, status := e.Auth(&cookie.Value)

	if status {

		parties := []models.Party{}

		e.Service.Get_confirmed_parties(&parties, &user_name)

		events := Events{
			Parties: parties,
		}

		renderer.Renderer.Render(w, "parties", r, events)
	}
}

func (e *Event) Have_been(w http.ResponseWriter, r *http.Request) {
}

func (e *Event) Chat(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")

	if err != nil {
		http.Redirect(w, r, "/log_in", http.StatusMovedPermanently)
		return
	}

	user_name, status := e.Auth(&cookie.Value)

	if status {
		id := r.URL.Query().Get("id")

		var name string
		var about string
		var city string
		var organizer string
		var date string
		var messages = []models.Message{}

		e.Service.Get_single_party(&id, &name, &about, &city, &organizer, &date)
		e.Service.Get_messages(&id, &user_name, &messages)

		party := models.Party{
			Name:      name,
			City:      city,
			Organizer: organizer,
			Date:      date,
			About:     about,
			Id:        id,
			Messages:  messages,
		}

		renderer.Renderer.Render(w, "with_chat", r, party)
	}
}

func (e *Event) Single_party(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")

	if err != nil {
		http.Redirect(w, r, "/log_in", http.StatusMovedPermanently)
		return
	}

	_, status := e.Auth(&cookie.Value)

	if status {
		id := r.URL.Query().Get("id")

		var name string
		var about string
		var city string
		var organizer string
		var date string

		e.Service.Get_single_party(&id, &name, &about, &city, &organizer, &date)

		party := models.Party{
			Name:      name,
			City:      city,
			Organizer: organizer,
			Date:      date,
			About:     about,
			Id:        id,
		}

		renderer.Renderer.Render(w, "without_chat", r, party)
	}
}

func (e *Event) Send_message(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")

	if err != nil {
		http.Redirect(w, r, "/log_in", http.StatusMovedPermanently)
		return
	}

	user_name, status := e.Auth(&cookie.Value)

	if status {
		id := r.URL.Query().Get("id")
		whose := r.URL.Query().Get("whose")

		if err := r.ParseForm(); err != nil {
			fmt.Println("Unparsable")
			return
		}

		body := r.FormValue("ci" + id)
		if body == "" {
			fmt.Println("ci" + id)
			fmt.Println("empty lol")
			return
		}

		e.Service.Send_message(&user_name, &id, &whose, &body)

		message := models.Message{
			Whose: whose,
			Text:  body,
		}

		renderer.Renderer.Render(w, "message", r, message)
	}
}

func (e *Event) My_parties(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")

	if err != nil {
		http.Redirect(w, r, "/log_in", http.StatusMovedPermanently)
		return
	}

	user_name, status := e.Auth(&cookie.Value)

	if status {
		parties := []models.Party{}
		e.Service.Get_parties_organized_by_user(&user_name, &parties)

		renderer.Renderer.Render(w, "my_parties", r, parties)
	}
}

func (e *Event) History(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("authToken")
	if err != nil {
		http.Redirect(w, r, "/log_in", http.StatusMovedPermanently)
		return
	}

	renderer.Renderer.Render(w, "history", r, nil)
}
