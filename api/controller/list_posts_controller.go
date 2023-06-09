package controller

import (
	"errors"
	"fmt"
	"forum/domain/usecase"
	"forum/pkg/utils"
	"net/http"
	"strconv"
)

type ListPostsController struct {
	ListPostUsecase usecase.ListPostsUsecase
	*Controller
}

func (lpc *ListPostsController) ListPostsController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		lpc.logger.PrintError(fmt.Errorf("list-post: method not allowed"))
		lpc.ResponseMethodNotAllowed(w)
		return
	}

	if r.URL.Path != lpc.Data.Endpoints.PostsAllEndpoint {
		lpc.logger.PrintError(fmt.Errorf("list-post: not found"))
		lpc.ResponseNotFound(w)
		return
	}

	user := lpc.contextGetUser(r)
	if user != nil {
		lpc.Data.IsAuthorized = true
	} else {
		lpc.Data.IsAuthorized = false
	}

	_, ok1 := r.URL.Query()["filter"]
	_, ok2 := r.URL.Query()["category_filter"]

	if !ok1 && !ok2 && len(r.URL.Query()) != 0 {
		lpc.ResponseBadRequest(w)
		lpc.logger.PrintError(fmt.Errorf("list-post: incorrect filter"))
		return
	}

	queryMap := r.URL.Query()
	filter, ok := queryMap["filter"]
	if !ok {
		categoryMap, ok := queryMap["category_filter"]
		if !ok {
			posts, err := lpc.ListPostUsecase.GetAllPosts()
			if err != nil {
				lpc.logger.PrintError(fmt.Errorf("list-post: %w", err))
				lpc.ResponseServerError(w)
				return
			}

			lpc.Data.Posts = posts

			categories, err := lpc.ListPostUsecase.GetAllCategories()
			if err != nil {
				lpc.logger.PrintError(fmt.Errorf("list-post: %w", err))
				lpc.ResponseServerError(w)
			}

			lpc.Data.Categories = categories

			err = lpc.tmpl.ExecuteTemplate(w, "index.html", lpc.Data)
			if err != nil {
				lpc.logger.PrintError(fmt.Errorf("list-post: ExecuteTemplate error: %w", err))
				lpc.ResponseServerError(w)
				return
			}

			return
		}

		ids := []int{}
		for _, categoryIDStr := range categoryMap {
			categoryID, err := strconv.Atoi(categoryIDStr)
			if err != nil {
				lpc.logger.PrintError(fmt.Errorf("list-post: invalid category: %w", err))
				lpc.ResponseBadRequest(w)
				return
			}

			ids = append(ids, categoryID)
		}

		posts, err := lpc.ListPostUsecase.GetPostsByCategories(ids...)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("list-post: %w", err))
			lpc.ResponseServerError(w)
			return
		}
		lpc.Data.Posts = posts

		categories, err := lpc.ListPostUsecase.GetAllCategories()
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("list-post: %w", err))
			lpc.ResponseServerError(w)
		}

		lpc.Data.Categories = categories

		err = lpc.tmpl.ExecuteTemplate(w, "index.html", lpc.Data)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("list-post: ExecuteTemplate error: %w", err))
			lpc.ResponseServerError(w)
			return
		}

		return

	}

	if user == nil {
		lpc.logger.PrintError(fmt.Errorf("list-post: unauthorized"))
		lpc.ResponseUnauthorized(w)
		return
	}

	switch filter[0] {
	case "created":
		posts, err := lpc.ListPostUsecase.GetCreatedPosts(user.ID)
		if err != nil {
			if errors.Is(err, utils.ErrNoPosts) {
				err = lpc.tmpl.ExecuteTemplate(w, "index.html", lpc.Data)
				if err != nil {
					lpc.logger.PrintError(fmt.Errorf("list-post: created ExecuteTemplate index.html: %w", err))
					lpc.ResponseServerError(w)
					return
				}
			}

			lpc.logger.PrintError(fmt.Errorf("list-post: %w", err))
			lpc.ResponseServerError(w)
			return
		}

		lpc.Data.Posts = posts

		categories, err := lpc.ListPostUsecase.GetAllCategories()
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("list-post: %w", err))
			lpc.ResponseServerError(w)
		}

		lpc.Data.Categories = categories
		err = lpc.tmpl.ExecuteTemplate(w, "index.html", lpc.Data)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("list-post: created ExecuteTemplate index.html: %w", err))
			lpc.ResponseServerError(w)
			return
		}
	case "liked":
		posts, err := lpc.ListPostUsecase.GetLikedPosts(user.ID)
		if err != nil {
			if errors.Is(err, utils.ErrNoPosts) {
				err = lpc.tmpl.ExecuteTemplate(w, "index.html", lpc.Data)
				if err != nil {
					lpc.logger.PrintError(fmt.Errorf("list-post: liked ExecuteTemplate index.html: %w", err))
					lpc.ResponseServerError(w)
					return
				}
			}
			lpc.logger.PrintError(fmt.Errorf("list-post: %w", err))
			lpc.ResponseServerError(w)
			return
		}

		lpc.Data.Posts = posts

		categories, err := lpc.ListPostUsecase.GetAllCategories()
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("list-post: %w", err))
			lpc.ResponseServerError(w)
		}

		lpc.Data.Categories = categories

		err = lpc.tmpl.ExecuteTemplate(w, "index.html", lpc.Data)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("list-post: ExecuteTemplate index.html: %w", err))
			lpc.ResponseServerError(w)
			return
		}

	default:
		posts, err := lpc.ListPostUsecase.GetAllPosts()
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("list-post: %w", err))
			lpc.ResponseServerError(w)
			return
		}

		lpc.Data.Posts = posts

		categories, err := lpc.ListPostUsecase.GetAllCategories()
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("list-post: %w", err))
			lpc.ResponseServerError(w)
		}

		lpc.Data.Categories = categories

		err = lpc.tmpl.ExecuteTemplate(w, "index.html", lpc.Data)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("list-post: ExecuteTemplate error: %w", err))
			lpc.ResponseServerError(w)
			return
		}
		return
	}
}
