package controller

import (
	"errors"
	"fmt"
	"net/http"

	"forum/domain/usecase"
)

type LogoutConrtroller struct {
	LogoutUsecase usecase.LogoutUsecase
	Controller
}

func (lc *LogoutConrtroller) Logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		lc.logger.PrintError(fmt.Errorf("Controller: logout: not found"))
		lc.ResponseNotFound(w)
		return
	}

	c, err := r.Cookie("access_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			lc.logger.PrintError(fmt.Errorf("Controller: logout: unauthorized"))
			lc.errorPage(w, http.StatusUnauthorized)
			return
		}
		fmt.Println("Controller: logout: " + err.Error())
		lc.ResponseServerError(w)
		return
	}

	err = lc.LogoutUsecase.DeleteToken(c.Value)
	if err != nil {
		lc.logger.PrintError(fmt.Errorf("Controller: logout: couldn't delete token"))
		lc.ResponseServerError(w)
		return
	}

	cookie := &http.Cookie{
		Name:   "access_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
