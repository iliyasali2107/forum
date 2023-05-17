package controller

import (
	"fmt"
	"net/http"

	"forum/domain/models"
	"forum/domain/usecase"
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
		user.Name = r.FormValue("name")
		user.Email = r.FormValue("email")
		user.Password.Plaintext = r.FormValue("password")

		err := sc.SignupUsecase.Signup(sc.validator, user)
		if err == ErrUserExists {
			sc.logger.PrintError(err)
			sc.ResponseEditConflict(w)
			return
		}

		if err == ErrInternalServer {
			sc.logger.PrintError(err)
			sc.ResponseServerError(w)
			return
		}

		if err == ErrFormValidation {
			sc.ResponseBadRequest(w)
			sc.render(w, "signup.html", sc.validator)
			sc.validator.Errors = map[string]string{}
			return
		}

		http.Redirect(w, r, sc.Data.Endpoints.LoginEndpoint, http.StatusSeeOther)
	default:
		sc.logger.PrintError(fmt.Errorf("Controller: signup: method not allowed"))
		sc.ResponseMethodNotAllowed(w)
		return
	}
}
