package controller

import (
	"fmt"
	"net/http"

	"forum/domain/usecase"
)

type ListPostsController struct {
	ListPostUsecase usecase.ListPostsUsecase
	*Controller
}

func (lpc *ListPostsController) ListPostsController(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		lpc.logger.PrintError(fmt.Errorf("Controller: listPost: method not allowed"))
		lpc.ResponseMethodNotAllowed(w)
		return
	}

	if r.URL.Path != lpc.Data.Endpoints.PostsAllEndpoint {
		lpc.logger.PrintError(fmt.Errorf("Controller: listPost: not found"))
		lpc.ResponseNotFound(w)
		return
	}

	user := lpc.contextGetUser(r)

	queryMap := r.URL.Query()
	filter, ok := queryMap["filter"]
	if !ok {
		posts, err := lpc.ListPostUsecase.GetAllPosts()
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("Controller: listPost: GetAllPosts error"))
			lpc.logger.PrintError(err)
			lpc.ResponseServerError(w)
			return
		}

		lpc.Data.Posts = posts
		err = lpc.tmpl.ExecuteTemplate(w, "show_posts.html", lpc.Data)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("Controller: listPost: ExecuteTemplate error"))
			lpc.logger.PrintError(err)
			lpc.ResponseServerError(w)
			return
		}
		return
	}

	if user == nil {
		lpc.logger.PrintError(fmt.Errorf("Controller: listPost: unauthorized"))
		lpc.ResponseUnauthorized(w)
		return
	}

	switch filter[0] {
	case "created":
		posts, err := lpc.ListPostUsecase.GetCreatedPosts(user.ID)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("Controller: listPost: created GetCreatedPosts"))
			lpc.logger.PrintError(err)
			lpc.ResponseServerError(w)
			return
		}

		lpc.Data.Posts = posts
		err = lpc.tmpl.ExecuteTemplate(w, "show_posts.html", lpc.Data)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("Controller: listPost: created ExecuteTemplate show_posts.html"))
			lpc.logger.PrintError(err)
			lpc.ResponseServerError(w)
			return
		}
	case "liked":
		posts, err := lpc.ListPostUsecase.GetLikedPosts(user.ID)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("Controller: listPost: liked GetLikedPosts"))
			lpc.logger.PrintError(err)
			lpc.ResponseServerError(w)
			return
		}

		lpc.Data.Posts = posts
		err = lpc.tmpl.ExecuteTemplate(w, "show_posts.html", lpc.Data)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("Controller: listPost: liked ExecuteTemplate show_posts.html"))
			lpc.logger.PrintError(err)
			lpc.ResponseServerError(w)
			return
		}
	case "disliked":
		posts, err := lpc.ListPostUsecase.GetDislikedPosts(user.ID)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("Controller: listPost: disliked GetLikedPosts"))
			lpc.logger.PrintError(err)
			lpc.ResponseServerError(w)
			return
		}

		lpc.Data.Posts = posts
		err = lpc.tmpl.ExecuteTemplate(w, "show_posts.html", lpc.Data)
		if err != nil {
			lpc.logger.PrintError(fmt.Errorf("Controller: listPost: disliked ExecuteTemplate show_posts.html"))
			lpc.logger.PrintError(err)
			lpc.ResponseServerError(w)
			return
		}
	default:
		lpc.logger.PrintError(fmt.Errorf("Controller: listPost: bad request"))
		lpc.ResponseBadRequest(w)
		return
	}
}
