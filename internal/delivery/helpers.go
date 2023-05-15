package delivery

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (h *Controller) background(fn func()) {
	// Increment the WaitGroup counter
	h.wg.Add(1)

	go func() {
		// Use defer to decrement the WaitGroup counter before the goroutine returns.
		defer h.wg.Done()

		// Recover from any panic
		defer func() {
			if err := recover(); err != nil {
				h.logger.PrintError(fmt.Errorf("error: error while recovering"))
			}
		}()

		// Execute the arbitrary function that we passed as the parameter
		fn()
	}()
}

func (h *Controller) render(w http.ResponseWriter, name string, td any) {
	err := h.tmpl.ExecuteTemplate(w, name, td)
	if err != nil {
		h.logger.PrintInfo("render: " + err.Error())
		h.ResponseServerError(w)
		return
	}
}

func GetIdFromURL(path string) (int, error) {
	s := strings.Split(path, "/")

	if len(s) <= 3 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	if len(s[3:]) > 1 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	id, err := strconv.Atoi(s[3])
	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetIdFromShortURL(path string) (int, error) {
	s := strings.Split(path, "/")
	if len(s) <= 2 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	if len(s[2:]) > 1 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	id, err := strconv.Atoi(s[2])
	if err != nil {
		return 0, err
	}

	return id, nil
}

// TODO: replace all GetIdFromURL and GetIdFromShortURL with function below
func GetIdFromURL2(numOfWords int, path string) (int, error) {
	s := strings.Split(path, "/")
	if len(s) <= numOfWords+1 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	if len(s[numOfWords+1:]) > 1 {
		return 0, fmt.Errorf("%s", "invalid url")
	}

	id, err := strconv.Atoi(s[numOfWords+1])
	if err != nil {
		return 0, err
	}

	return id, nil
}
