package controller

import (
	"fmt"
	"net/http"
	"strings"

	"forum/domain/models"
	"forum/domain/usecase"
	"forum/pkg/utils"
	"forum/pkg/validator"
)

type SignupController struct {
	SignupUsecase usecase.SignupUsecase
	*Controller
}

func (sc *SignupController) Signup(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != sc.Data.Endpoints.SignupEndpoint {
		sc.logger.PrintError(fmt.Errorf("Controller: signup: not found"))
		sc.ResponseNotFound(w)
		return
	}

	switch r.Method {
	case http.MethodGet:
		err := sc.tmpl.ExecuteTemplate(w, "signup.html", nil)
		if err != nil {
			sc.logger.PrintError(fmt.Errorf("Controller: signup: domain server error"))
			sc.ResponseServerError(w)
			return
		}
	case http.MethodPost:
		user := &models.User{}
		user.Name = strings.TrimSpace(r.FormValue("name"))
		user.Email = strings.TrimSpace(r.FormValue("email"))
		user.Password.Plaintext = strings.TrimSpace(r.FormValue("password"))

		signupErrors := validator.SignupValidation(user)
		if len(signupErrors) != 0 {
			sc.Data.Errors = signupErrors
			w.WriteHeader(http.StatusBadRequest)
			sc.render(w, "signup.html", sc.Data)
			return
		}

		err := sc.SignupUsecase.Signup(user)
		if err == utils.ErrNameIsTaken {
			errMap := make(map[string]string)
			errMap["name"] = "name is already in use"
			sc.Data.Errors = errMap
			sc.logger.PrintError(err)
			w.WriteHeader(http.StatusBadRequest)
			sc.render(w, "signup.html", sc.Data)
			return
		}

		if err == utils.ErrEmailIsTaken {
			errMap := make(map[string]string)
			errMap["email"] = "email is already in use"
			sc.Data.Errors = errMap
			sc.logger.PrintError(err)
			w.WriteHeader(http.StatusBadRequest)
			sc.render(w, "signup.html", sc.Data)
			return
		}

		if err == utils.ErrInternalServer {
			sc.logger.PrintError(err)
			sc.ResponseServerError(w)
			return
		}

		http.Redirect(w, r, sc.Data.Endpoints.LoginEndpoint, http.StatusSeeOther)
	default:
		sc.logger.PrintError(fmt.Errorf("Controller: signup: method not allowed"))
		sc.ResponseMethodNotAllowed(w)
		return
	}
}
