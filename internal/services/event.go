package services

import (
	"github.com/alogsDiu/Oqiga/internal/models"
	"github.com/alogsDiu/Oqiga/internal/repos"
)

type Event struct {
	Repo *repos.Event
}

func (e *Event) Get_parties(parites *[]models.Party) error {
	return e.Repo.Get_all_parties(parites)
}
func (e *Event) Get_single_party(id *string, name *string, about *string, city *string, organizer *string, date *string) {
	e.Repo.Get_single_party(id, name, about, city, organizer, date)
}

func (e *Event) Get_confirmed_parties(parites *[]models.Party, user_name *string) {
	e.Repo.Get_confirmed_parties(parites, user_name)
}
func (e *Event) Send_message(user_name, party_id, whose, body *string) {
	e.Repo.Send_message(user_name, party_id, whose, body)
}
func (e *Event) Get_messages(id, user_name *string, messages *[]models.Message) {
	e.Repo.Get_messages(id, user_name, messages)
}
func (e *Event) Get_parties_organized_by_user(user_name *string, parties *[]models.Party) {
	e.Repo.Get_parties_organized_by_user(user_name, parties)
}
