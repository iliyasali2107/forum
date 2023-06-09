package controller

import (
	"fmt"
	"forum/domain/usecase"
	"net/http"
)

type PostController struct {
	PostDetailsUsecase     usecase.PostDetailsUsecase
	CommentReactionUsecase usecase.CommentReactionUsecase
	*Controller
}

func (pc *PostController) PostController(w http.ResponseWriter, r *http.Request) {
	user := pc.contextGetUser(r)
	if user != nil {
		pc.Data.IsAuthorized = true
	} else {
		pc.Data.IsAuthorized = false
	}

	id, err := GetIdFromShortURL(r.URL.Path)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: not found: %w", err))
		pc.ResponseNotFound(w)
		return
	}

	post, err := pc.PostDetailsUsecase.GetPost(id)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: GetPost: %w", err))
		pc.ResponseNotFound(w)
		return
	}

	postLikes, err := pc.PostDetailsUsecase.GetPostLikes(id)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: GetPostLikes: %w", err))
		pc.ResponseServerError(w)
		return
	}

	postDislikes, err := pc.PostDetailsUsecase.GetPostDislikes(id)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: GetPostDislikes: %w", err))
		pc.ResponseServerError(w)
		return
	}

	comments, err := pc.PostDetailsUsecase.GetCommentsByPostId(id)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: GetPostDislikes: %w", err))
		pc.ResponseServerError(w)
		return
	}

	post.Likes = postLikes
	post.Dislikes = postDislikes
	post.Comments = comments

	for i, comment := range post.Comments {

		likes, err := pc.CommentReactionUsecase.CommentLikeCount(comment.ID)
		if err != nil {
			pc.logger.PrintError(err)
			pc.ResponseServerError(w)
			return
		}

		post.Comments[i].Likes = likes

		dislikes, err := pc.CommentReactionUsecase.CommentDislikeCount(comment.ID)
		if err != nil {
			pc.logger.PrintError(err)
			pc.ResponseServerError(w)
			return
		}

		post.Comments[i].Dislikes = dislikes
	}

	pc.Data.Post = post
	err = pc.tmpl.ExecuteTemplate(w, "post.html", pc.Data)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: ExecuteTemplate post.html: %w", err))
		pc.ResponseServerError(w)
		return
	}
}
