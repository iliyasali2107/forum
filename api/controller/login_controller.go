package controller

import (
	"fmt"
	"net/http"
	"strings"

	"forum/domain/models"
	"forum/domain/usecase"
	"forum/pkg/utils"
)

type LoginController struct {
	LoginUsecase usecase.LoginUsecase
	*Controller
}

func (lc *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != lc.Data.Endpoints.LoginEndpoint {
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
		user.Email = strings.TrimSpace(r.FormValue("email"))
		user.Password.Plaintext = strings.TrimSpace(r.FormValue("password"))

		errors := make(map[string]string)

		err := lc.LoginUsecase.Login(user)
		if err != nil {
			switch err {
			case utils.ErrUserNotFound:
				errors["email"] = "user not found"
				lc.Data.Errors = errors
				w.WriteHeader(http.StatusBadRequest)
				lc.render(w, "login.html", lc.Data)
				return
			case utils.ErrInvalidPassword:
				errors["password"] = "password is not correct"
				lc.Data.Errors = errors
				w.WriteHeader(http.StatusBadRequest)
				lc.render(w, "login.html", lc.Data)
				return
			default:
				lc.logger.PrintError(fmt.Errorf("handler: login: %w", err))
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

		lc.Data.IsAuthorized = true

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
