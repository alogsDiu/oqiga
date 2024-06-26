package repos

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	DB *sql.DB
}

func (u *User) Get_user_pass_from_username(user_name *string) (string, int) {
	query := "select pass from user where user_name=?"

	var pass string

	err := u.DB.QueryRow(query, *user_name).Scan(&pass)
	if err != nil {
		return "", 0
	}
	return pass, 2
}

func (u *User) Create_new_user(user_name *string, pass *string, email *string) error {
	query := "INSERT INTO USER(user_name, pass, email, rating, number,about) VALUES (?, ?, ?, ?, ?, ?)"

	_, err := u.DB.Exec(query, *user_name, *pass, *email, 50, "", "")

	if err != nil {
		fmt.Println("THE REGISTRATION IS ROUGH")
		return err
	}

	return nil
}

func (u *User) Does_user_have_a_session(token *string) bool {

	var valid_until string

	query := "SELECT valid_until FROM tokens WHERE token=?"

	err := u.DB.QueryRow(query, *token).Scan(&valid_until)

	if err != nil {
		fmt.Println("Error executing query:", err)
		return false
	}

	stored_time, err := time.Parse(time.ANSIC, valid_until)

	if err != nil {
		fmt.Println("Yo the stored time is not being able to be parsed", err)
	}

	if time.Now().Before(stored_time) {
		return true
	} else {
		query = "delete * from tokens where token=?"
		_, err := u.DB.Exec(query, *token)
		if err != nil {
			fmt.Println("Unable to delete the users session", err)
		}
	}

	return false
}

func (u *User) Add_authorized_user(user_name *string, token *string, valid_util *string) error {
	_, err := u.DB.Exec("delete  from tokens  where user_name=?", *user_name)

	if err != nil {
		return err
	}

	query := "INSERT INTO tokens(user_name, token, valid_until) values (?,?,?)"

	_, err = u.DB.Exec(query, *user_name, *token, valid_util)

	if err != nil {
		fmt.Println("so there are some troubles with registering token to a user ", err)
		return err
	}

	return nil
}

func (u *User) Get_user_info_from_token(token *string, user_name *string, email *string, number *string, about *string, rating *float32, coments *[]string) {
	query_to_get_a_user_name := "select user_name from tokens"

	err := u.DB.QueryRow(query_to_get_a_user_name, *token).Scan(user_name)
	if err != nil {
		fmt.Println("Unable to get a username", err)
		return
	}

	query_to_get_user_info := "select email, number, about, rating from user where user_name = ?"

	err = u.DB.QueryRow(query_to_get_user_info, *user_name).Scan(email, number, about, rating)

	if err != nil {
		fmt.Println("Unable to get a username", err)
		return
	}

	array_of_comments := []string{}

	query_to_get_comments := "select comment from coments where user_name = ?"

	rows, err := u.DB.Query(query_to_get_comments, *user_name)

	if err != nil {
		fmt.Println("Unable to get comments", err)
		return
	}

	for rows.Next() {
		var comment string
		err := rows.Scan(&comment)
		if err != nil {
			fmt.Println("One row got fed up", err)
		}
		array_of_comments = append(array_of_comments, comment)
	}

	*coments = array_of_comments
}

func (u *User) Get_name_and_status(token *string) (string, bool) {

	var valid_until string
	var user_name string

	query := "SELECT valid_until,user_name FROM tokens WHERE token=?"

	err := u.DB.QueryRow(query, *token).Scan(&valid_until, &user_name)

	if err != nil {
		fmt.Println("Error executing query:", err)
		return "", false
	}

	stored_time, err := time.Parse(time.ANSIC, valid_until)

	if err != nil {
		fmt.Println("Yo the stored time is not being able to be parsed", err)
		return "", false
	}

	if time.Now().Before(stored_time) {
		return user_name, true
	} else {
		query = "delete * from tokens where token=?"
		_, err := u.DB.Exec(query, *token)
		if err != nil {
			fmt.Println("Unable to delete the users session", err)
		}
	}

	return "", false
}
