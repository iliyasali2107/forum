package controller

import (
	"fmt"
	"net/http"

	"forum/domain/models"
	"forum/domain/usecase"
)

type CreatePostController struct {
	CreatePostUsecase usecase.CreatePostUsecase

	Controller
}

// PostsController TODO: invalid field messages for each field
func (cpc *CreatePostController) CreatePostController(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/posts/create" {
		cpc.logger.PrintError(fmt.Errorf("Controller: post-create: not found"))
		cpc.ResponseNotFound(w)
		return
	}

	user := cpc.contextGetUser(r)

	switch r.Method {
	case http.MethodGet:
		categories, err := cpc.CreatePostUsecase.GetAllCategories()

		data := struct {
			Errors     map[string]string
			Categories []*models.Category
		}{
			Categories: categories,
		}

		if err != nil {
			cpc.logger.PrintError(fmt.Errorf("Controller: post-create: GetAllCategories error"))
			cpc.logger.PrintError(err)
			cpc.ResponseServerError(w)
			return
		}

		err = cpc.tmpl.ExecuteTemplate(w, "create_post.html", data)
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
		post.Content = r.FormValue("content")
		post.User = user

		err = cpc.CreatePostUsecase.CreatePost(cpc.validator, post)
		if err != nil {
			if err == ErrFormValidation {
				categories, err := cpc.CreatePostUsecase.GetAllCategories()
				if err != nil {
					cpc.logger.PrintError(fmt.Errorf("Controller: post-create: CreatePost ErrFormValidation error"))
					cpc.logger.PrintError(err)
					cpc.ResponseServerError(w)
					return
				}

				data := struct {
					Errors     map[string]string
					Categories []*models.Category
				}{
					Errors:     cpc.validator.Errors,
					Categories: categories,
				}
				cpc.ResponseBadRequest(w)
				cpc.render(w, "create_post.html", data)
				return
			}

			// TODO: ------------------
			cpc.logger.PrintError(err)
			return
		}
		http.Redirect(w, r, "/posts", http.StatusSeeOther)
	default:
		cpc.logger.PrintError(fmt.Errorf("create-post: method not allowed"))
		cpc.ResponseMethodNotAllowed(w)
		return
	}
}
