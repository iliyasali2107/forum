package controller

import (
	"database/sql"
	"html/template"
	"os"

	"forum/domain/models"
	"forum/domain/repository"
	"forum/domain/usecase"
	"forum/pkg/logger"
	"forum/pkg/validator"
)

type Controller struct {
	tmpl         *template.Template
	logger       *logger.Logger
	validator    *validator.Validator
	Data         Data
	TokenUsecase usecase.TokenUsecase
}

// TODO:
type Data struct {
	Endpoints  Endpoints
	Post       *models.Post
	Posts      []*models.Post
	Comment    *models.Comment
	Comments   []*models.Comment
	Errors     map[string]string
	Categories []*models.Category
}

type Endpoints struct {
	SignupEndpoint             string
	LoginEndpoint              string
	LogoutEndpoint             string
	CreatePostEndpoint         string
	PostDetailsEndpoint        string
	PostsAllEndpoint           string
	CreateCommentEndpoint      string
	CommentDetailsEndpoint     string
	CreatePostReactionEndpoint string
}

func NewController(db *sql.DB) *Controller {
	endpts := Endpoints{
		SignupEndpoint:             "/signup",
		LoginEndpoint:              "/login",
		LogoutEndpoint:             "/logout",
		CreatePostEndpoint:         "/posts/create",
		PostDetailsEndpoint:        "/posts/",
		PostsAllEndpoint:           "/",
		CreateCommentEndpoint:      "/posts/comment/create",
		CommentDetailsEndpoint:     "/posts/comment/",
		CreatePostReactionEndpoint: "/post/reaction/create",
	}

	ur := repository.NewUserRepository(db)

	return &Controller{
		tmpl:         template.Must(template.ParseGlob("./templates/*")),
		validator:    validator.NewValidator(),
		logger:       logger.NewLogger(os.Stdout, logger.LevelInfo),
		Data:         Data{Endpoints: endpts},
		TokenUsecase: usecase.NewTokenUsecae(ur),
	}
}
