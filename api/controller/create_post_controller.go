package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"forum/domain/models"
	"forum/domain/usecase"
)

type CreatePostController struct {
	CreatePostUsecase usecase.CreatePostUsecase

	*Controller
}

// PostsController TODO: invalid field messages for each field
func (cpc *CreatePostController) CreatePostController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != cpc.Data.Endpoints.CreatePostEndpoint {
		cpc.logger.PrintError(fmt.Errorf("Controller: post-create: not found"))
		cpc.ResponseNotFound(w)
		return
	}

	user := cpc.contextGetUser(r)

	switch r.Method {
	case http.MethodGet:
		categories, err := cpc.CreatePostUsecase.GetAllCategories()

		if err != nil {
			cpc.logger.PrintError(fmt.Errorf("Controller: post-create: GetAllCategories error"))
			cpc.logger.PrintError(err)
			cpc.ResponseServerError(w)
			return
		}

		cpc.Data.Categories = categories

		err = cpc.tmpl.ExecuteTemplate(w, "create_post.html", cpc.Data)
		if err != nil {
			cpc.logger.PrintError(fmt.Errorf("Controller: post-create: ExecuteTemplate error"))
			cpc.logger.PrintError(err)
			cpc.ResponseServerError(w)
			return
		}
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			cpc.logger.PrintError(fmt.Errorf("Controller: post-create: ParseForm error"))
			cpc.logger.PrintError(err)
			cpc.ResponseBadRequest(w)
			return
		}

		post := &models.Post{}
		post.Categories = r.Form["category"]
		post.Title = r.FormValue("title")
		title := strings.TrimSpace(r.FormValue("title"))
		if title != "" {
			post.Title = title
		}

		content := strings.TrimSpace(r.FormValue("content"))
		if content != "" {
			post.Content = content
		}

		post.User = user

		postID, err := cpc.CreatePostUsecase.CreatePost(cpc.validator, post)
		if err != nil {
			if errors.Is(err, usecase.ErrFormValidation) {
				categories, err := cpc.CreatePostUsecase.GetAllCategories()
				if err != nil {
					cpc.logger.PrintError(fmt.Errorf("Controller: post-create: CreatePost ErrFormValidation error"))
					cpc.logger.PrintError(err)
					cpc.ResponseServerError(w)
					return
				}
				cpc.Data.Categories = categories
				cpc.Data.Errors = cpc.validator.Errors
				cpc.validator.Errors = map[string]string{}

				// cpc.ResponseBadRequest(w)
				cpc.render(w, "create_post.html", cpc.Data)
				return
			}

			cpc.logger.PrintError(err)
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
