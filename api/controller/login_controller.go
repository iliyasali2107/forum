package controller

import (
	"fmt"
	"net/http"

	"forum/domain/models"
	"forum/domain/usecase"
)

type LoginController struct {
	LoginUsecase usecase.LoginUsecase
	*Controller
}

func (lc *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		lc.logger.PrintError(fmt.Errorf("handler: login: not found"))
		lc.ResponseNotFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := lc.tmpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			lc.logger.PrintError(err)
			lc.ResponseServerError(w)
			return
		}
	case http.MethodPost:
		user := &models.User{}
		user.Email = r.FormValue("email")
		user.Password.Plaintext = r.FormValue("password")

		err := lc.LoginUsecase.Login(user)
		if err != nil {
			switch err {
			case ErrUserNotFound:
				lc.logger.PrintError(fmt.Errorf("handler: login: user not found"))
				lc.ResponseBadRequest(w)
				return
			case ErrInvalidPassword:
				lc.logger.PrintError(fmt.Errorf("handler: login: password is not correct"))
				lc.ResponseBadRequest(w)
				return
			default:
				lc.logger.PrintError(fmt.Errorf("handler: login: password is not correct"))
				lc.ResponseServerError(w)
				return
			}
		}

		cookie := http.Cookie{}
		cookie.Name = "access_token"
		cookie.Value = *user.Token
		cookie.Expires = *user.Expires
		cookie.Path = "/"
		cookie.HttpOnly = true
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
