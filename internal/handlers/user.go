package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/alogsDiu/Oqiga/internal/renderer"
	"github.com/alogsDiu/Oqiga/internal/services"
)

type User struct {
	Service *services.User
}

type login_problems struct {
	Val int8
}

func (u *User) AboutPage(w http.ResponseWriter, r *http.Request) {
	renderer.Renderer.Render(w, "about.html", r, nil)
}

func (u *User) Log_in(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")
	if err == nil {
		if u.Service.Authenticate(cookie.Value) {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
			return
		}
	}
	if r.Method == http.MethodGet {
		renderer.Renderer.Render(w, "log_in.html", r, nil)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			http.Error(w, "Unable to parse Login form", http.StatusBadRequest)
			return
		}

		uname := r.FormValue("username")
		pass := r.FormValue("password")

		pass = strings.Trim(pass, " ")

		token, problem_type := u.Service.Authorize(&uname, &pass) //problem_type => int, 0 => problem with username, 1 => problem with pass, 2 => OK!

		if problem_type < 2 {
			if problem_type == 0 {
				renderer.Renderer.Render(w, "login", r, login_problems{Val: 0})
			} else {
				renderer.Renderer.Render(w, "login", r, login_problems{Val: 1})
			}
			return
		}

		cookie := &http.Cookie{
			Name:     "authToken",
			Value:    token,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(4 * time.Hour),
		}

		http.SetCookie(w, cookie)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (u *User) Sign_up(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("authToken")
	if err == nil {
		if u.Service.Authenticate(cookie.Value) {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}

	if r.Method == http.MethodGet {
		renderer.Renderer.Render(w, "sign_up.html", r, nil)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			http.Error(w, "Unable to proces ur request", http.StatusBadRequest)
			return
		}

		uname := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		password = strings.Trim(password, " ")

		err = u.Service.Register(&uname, &password, &email)

		if err != nil {
			http.Redirect(w, r, "/sign_up", http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/confirm", http.StatusSeeOther)
	}
}

func (u *User) Password_restoration(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")
	if err == nil {
		if u.Service.Authenticate(cookie.Value) {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
	}

	if r.Method == http.MethodGet {
		renderer.Renderer.Render(w, "forgot_password.html", r, nil)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Unable to proces ur request", http.StatusBadRequest)
			return
		}
		uname := r.FormValue("username")
		pass1 := r.FormValue("new_password1")
		pass2 := r.FormValue("new_password2")

		u.Service.ChangePassword(uname, pass1, pass2)

		http.Redirect(w, r, "/confirm", http.StatusSeeOther)
	}
}

func (u *User) Confirm(w http.ResponseWriter, r *http.Request) {
	renderer.Renderer.Render(w, "confirm.html", r, nil)
}

type Userprofile struct {
	Name    string
	Number  string
	Email   string
	About   string
	Rating  float32
	Coments []string
}

func (u *User) Dashboard(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")

	if err != nil {
		http.Redirect(w, r, "/log_in", http.StatusMovedPermanently)
		return
	}
	token := cookie.Value

	if !u.Service.Authenticate(token) {
		http.Redirect(w, r, "/log_in", http.StatusMovedPermanently)
	}

	var user Userprofile = Userprofile{}

	u.Service.GetUserProfileInfo(&token, &user.Name, &user.Email, &user.Number, &user.About, &user.Rating, &user.Coments)

	renderer.Renderer.Render(w, "main.html", r, user)
}

func (u *User) Profile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("authToken")

	if err != nil {
		http.Redirect(w, r, "/log_in", http.StatusMovedPermanently)
		return
	}
	token := cookie.Value

	if !u.Service.Authenticate(token) {
		http.Redirect(w, r, "/log_in", http.StatusMovedPermanently)
	}

	var user Userprofile = Userprofile{}

	u.Service.GetUserProfileInfo(&token, &user.Name, &user.Email, &user.Number, &user.About, &user.Rating, &user.Coments)

	renderer.Renderer.Render(w, "profile", r, user)
}
