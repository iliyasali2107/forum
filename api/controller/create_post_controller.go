package controller

import (
	"fmt"
	"net/http"
	"strings"

	"forum/domain/models"
	"forum/domain/usecase"
	"forum/pkg/validator"
)

type CreatePostController struct {
	CreatePostUsecase usecase.CreatePostUsecase

	*Controller
}

// PostsController TODO: invalid field messages for each field
func (cpc *CreatePostController) CreatePostController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != cpc.Data.Endpoints.CreatePostEndpoint {
		cpc.logger.PrintError(fmt.Errorf("create-post: not found"))
		cpc.ResponseNotFound(w)
		return
	}

	user := cpc.contextGetUser(r)
	if user != nil {
		cpc.Data.IsAuthorized = true
	} else {
		cpc.Data.IsAuthorized = false
	}

	switch r.Method {
	case http.MethodGet:
		categories, err := cpc.CreatePostUsecase.GetAllCategories()
		if err != nil {
			cpc.logger.PrintError(fmt.Errorf("create-post: %w", err))
			cpc.ResponseServerError(w)
			return
		}

		cpc.Data.Categories = categories

		err = cpc.tmpl.ExecuteTemplate(w, "create_post.html", cpc.Data)
		if err != nil {
			cpc.logger.PrintError(fmt.Errorf("post-create: ExecuteTemplate error: %w", err))
			cpc.ResponseServerError(w)
			return
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			cpc.logger.PrintError(fmt.Errorf("error while parsing form: %w", err))
			cpc.ResponseBadRequest(w)
			return
		}

		post := &models.Post{}
		post.Categories = r.Form["category"]
		post.Title = r.FormValue("title")
		title := strings.TrimSpace(r.FormValue("title"))
		post.Title = title

		content := strings.TrimSpace(r.FormValue("content"))
		post.Content = content

		errMap := validator.CreatePostValidation(post)
		if len(errMap) != 0 {
			cpc.Data.Errors = errMap
			categories, err := cpc.CreatePostUsecase.GetAllCategories()
			if err != nil {
				cpc.logger.PrintError(fmt.Errorf("create-post: %w", err))
				cpc.ResponseServerError(w)
				return
			}
			cpc.Data.Categories = categories
			w.WriteHeader(http.StatusBadRequest)
			cpc.render(w, "create_post.html", cpc.Data)
			return
		}

		post.User = user

		postID, err := cpc.CreatePostUsecase.CreatePost(post)
		if err != nil {
			cpc.logger.PrintError(fmt.Errorf("server error occured while creating a post: %w", err))
			cpc.ResponseServerError(w)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("%s%d", cpc.Data.Endpoints.PostDetailsEndpoint, postID), http.StatusSeeOther)
	default:
		cpc.logger.PrintError(fmt.Errorf("create-post: method not allowed"))
		cpc.ResponseMethodNotAllowed(w)
		return
	}
}
