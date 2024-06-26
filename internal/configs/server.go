package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alogsDiu/Oqiga/internal/handlers"
	"github.com/alogsDiu/Oqiga/internal/renderer"
	"github.com/alogsDiu/Oqiga/internal/repos"
	"github.com/alogsDiu/Oqiga/internal/services"
	_ "github.com/mattn/go-sqlite3"
)

func Run() {

	adr := ":8080"

	renderer.CreateRenderer("views/templates")

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("views/css"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("views/img"))))

	user_db, err := sql.Open("sqlite3", "db/user.db")

	if err != nil {
		log.Fatal("Yo dude not connecting to USER DATABASE :", err)
	}

	defer user_db.Close()

	userRepo := &repos.User{DB: user_db}
	userService := &services.User{Repo: userRepo}
	userHandler := &handlers.User{Service: userService}

	eventRepo := &repos.Event{DB: user_db}
	eventService := &services.Event{Repo: eventRepo}
	eventHandler := &handlers.Event{
		Service: eventService,
		Auth:    userHandler.Service.Authenticate_and_give_name,
	}

	http.HandleFunc("/", userHandler.AboutPage)

	http.HandleFunc("/log_in", userHandler.Log_in)
	http.HandleFunc("/sign_up", userHandler.Sign_up)
	http.HandleFunc("/forgot_pass", userHandler.Password_restoration)
	http.HandleFunc("/confirm", userHandler.Confirm)

	http.HandleFunc("/dashboard", userHandler.Dashboard)
	http.HandleFunc("/profile", userHandler.Profile)
	http.HandleFunc("/recomendations", eventHandler.Gen_recomendations)
	http.HandleFunc("/all_recomendations", eventHandler.All_recomendations)
	http.HandleFunc("/confirmed_recomendations", eventHandler.Confirmed_recomendations)
	http.HandleFunc("/my_history", eventHandler.Have_been)
	http.HandleFunc("/chat", eventHandler.Chat)
	http.HandleFunc("/without_chat", eventHandler.Single_party)
	http.HandleFunc("/send_message", eventHandler.Send_message)
	http.HandleFunc("/my_parties", eventHandler.My_parties)
	http.HandleFunc("/history", eventHandler.History)

	println("Server started at http://localhost" + adr)

	log.Fatal(http.ListenAndServe(adr, nil))

}
