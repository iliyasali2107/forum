package controller

import (
	"database/sql"
	"forum/domain/models"
	"forum/domain/repository"
	"forum/domain/usecase"
	"forum/pkg/logger"
	"html/template"
	"os"
)

type Controller struct {
	tmpl         *template.Template
	logger       *logger.Logger
	Data         Data
	TokenUsecase usecase.TokenUsecase
}

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
	CommentReactionEndpoint    string
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
		CommentReactionEndpoint:    "/comment/reaction/create",
	}

	ur := repository.NewUserRepository(db)

	return &Controller{
		tmpl:         template.Must(template.ParseGlob("./templates/*")),
		logger:       logger.NewLogger(os.Stdout, logger.LevelInfo),
		Data:         Data{Endpoints: endpts},
		TokenUsecase: usecase.NewTokenUsecae(ur),
	}
}
