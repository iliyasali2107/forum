package delivery

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"forum/pkg/validator"
)

type envelope map[string]interface{}

//func (h *Handler) readIdParam(r *http.Request) (int64, error) {
//	params := httprouter.ParamsFromContext(r.Context())
//
//	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
//	if err != nil || id < 1 {
//		return 0, errors.New("invalid id parameter")
//	}
//
//	return id, nil
//}

// The readString() helper returns value from the query string, or the provided
// default value if no matching key could be found.
/*func (h *Handler) readString(qs url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the query string. If no key exists
	// this will return the empty string "".
	s := qs.Get(key)

	// If no keys exist (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}

	// Otherwise return the string
	return s
}*/

// The readCSV() helper reads a string value from the query string and then splits it
// into a slice on the comma character. If no matching key could be found, it returns
// the provided default value.
func (h *Handler) readCSV(qs url.Values, key string, v *validator.Validator, defaultValue []string) []string {
	// Extract the value from the query string.
	csv := qs.Get(key)

	// If no key exists (or the value is empty) the return the default value.
	if csv == "" {
		return defaultValue
	}

	// Otherwise parse the value into a []string slice and return it.
	return strings.Split(csv, ",")
}

// The readInt() helper reads a string value from the equery string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value. If the value couldn't be converted to an integer, then we record an
// error message tn the provided Validator instance.
func (h *Handler) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	// Extract the value from the query string.
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	// Try to convert the value to an int. If this fails, add an error message to the
	// validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func (h *Handler) background(fn func()) {
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

func (h *Handler) render(w http.ResponseWriter, name string, td any) {
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
