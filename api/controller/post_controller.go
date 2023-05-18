package controller

import (
	"fmt"
	"net/http"

	"forum/domain/usecase"
)

type PostController struct {
	PostDetailsUsecase usecase.PostDetailsUsecase
	*Controller
}

func (pc *PostController) PostController(w http.ResponseWriter, r *http.Request) {
	id, err := GetIdFromShortURL(r.URL.Path)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: not found"))
		pc.logger.PrintError(err)
		pc.ResponseNotFound(w)
		return
	}

	post, err := pc.PostDetailsUsecase.GetPost(id)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: GetPost"))
		pc.logger.PrintError(err)
		pc.ResponseNotFound(w)
		return
	}

	likes, err := pc.PostDetailsUsecase.GetPostLikes(id)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: GetPostLikes"))
		pc.logger.PrintError(err)
		pc.ResponseServerError(w)
		return
	}

	dislikes, err := pc.PostDetailsUsecase.GetPostDislikes(id)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: GetPostDislikes"))
		pc.logger.PrintError(err)
		pc.ResponseServerError(w)
		return
	}

	comments, err := pc.PostDetailsUsecase.GetCommentsByPostId(id)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: GetPostDislikes"))
		pc.logger.PrintError(err)
		pc.ResponseServerError(w)
		return
	}

	

	post.Likes = likes
	post.Dislikes = dislikes
	post.Comments = comments

	// comments, err := pc.Service.CommentService.GetCommentsByPostId(id)
	// fmt.Println(comments)
	// if err != nil {
	// 	pc.logger.PrintError(fmt.Errorf("Controller: PostController: GetCommentsByPostId"))
	// 	pc.logger.PrintError(err)
	// 	pc.ResponseServerError(w)
	// 	return
	// }

	// post.Comments = comments
	pc.Data.Post = post
	err = pc.tmpl.ExecuteTemplate(w, "post.html", pc.Data)
	if err != nil {
		pc.logger.PrintError(fmt.Errorf("Controller: PostController: ExecuteTemplate post.html"))
		pc.logger.PrintError(err)
		pc.ResponseServerError(w)
		return
	}
}
