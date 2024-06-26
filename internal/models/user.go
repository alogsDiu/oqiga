package models

type User struct {
	User_name string
	Email     string
	About     string
	Rating    float32
}

type Coment struct {
	User_name string
	Text      string
}
