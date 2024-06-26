package services

import (
	"fmt"
	"time"

	"github.com/alogsDiu/Oqiga/internal/repos"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Repo *repos.User
}

func (u *User) Authorize(username *string, pass *string) (string, int) {

	storedPass, status := u.Repo.Get_user_pass_from_username(username)

	if status == 0 {
		return "", 0
	}

	fmt.Println()
	fmt.Println("stored password :" + storedPass)
	fmt.Println()

	if !checkPasswordHash(*pass, storedPass) {
		return "", 1
	} else {
		token, err := createToken(*pass)
		if err != nil {
			fmt.Println("unable to create a token ", err)
			return "", 1
		}
		valid_until := time.Now().Add(3 * time.Hour).Format(time.ANSIC)
		u.Repo.Add_authorized_user(username, &token, &valid_until)

		return token, 2
	}
}
func (u *User) Authenticate(token string) bool {
	return u.Repo.Does_user_have_a_session(&token)
}
func (u *User) Authenticate_and_give_name(token *string) (string, bool) {
	return u.Repo.Get_name_and_status(token)
}

func (u *User) Register(uname *string, password *string, email *string) error {

	hashedPass, err := hashPassword(*password)
	fmt.Println("Registration moment :" + hashedPass)
	if err != nil {
		return err
	}

	return u.Repo.Create_new_user(uname, &hashedPass, email)
}
func (u *User) ChangePassword(uname string, pass1 string, pass2 string) {
	println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
}
func (u *User) GetUserProfileInfo(token *string, user_name *string, email *string, number *string, about *string, rating *float32, coments *[]string) {
	u.Repo.Get_user_info_from_token(token, user_name, email, number, about, rating, coments)
}
func hashPassword(password string) (string, error) {

	fmt.Println("What does hasher hash:", password, len(password))

	passtoStore, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passtoStore), err
}
func createToken(password string) (string, error) {
	token, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(token), err
}
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
