package controller

import (
	"html/template"
	"os"

	"forum/pkg/logger"
	"forum/pkg/validator"
)

type Controller struct {
	tmpl      *template.Template
	logger    *logger.Logger
	validator *validator.Validator
}


// TODO: 
type Data struct {

}

func NewController() *Controller {
	return &Controller{
		tmpl:      template.Must(template.ParseGlob("./templates/*")),
		validator: validator.NewValidator(),
		logger:    logger.NewLogger(os.Stdout, logger.LevelInfo),
	}
}
