package controller

import (
	"forum/internal/service"
	"forum/pkg/logger"
	"forum/pkg/validator"
	"html/template"
	"os"
	"sync"
)

type Controller struct {
	tmpl      *template.Template
	Service   *service.Service
	logger    *logger.Logger
	validator *validator.Validator
	wg        sync.WaitGroup
}

func NewController(service *service.Service) *Controller {
	return &Controller{
		tmpl:      template.Must(template.ParseGlob("./templates/*")),
		validator: validator.NewValidator(),
		logger:    logger.NewLogger(os.Stdout, logger.LevelInfo),
	}
}
