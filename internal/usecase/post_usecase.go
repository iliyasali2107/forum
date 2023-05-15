package usecase

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type postGetUsecase struct {
	postRepository     repository.PostRepository
	userRepository     repository.UserRepository
	categoryRepository repository.CategoryRepository
	reactionRepository repository.ReactionRepository
	commentRepository  repository.CommentRepository
}

type PostGetUsecase interface {
	GetPost(postID int) (*models.Post, error)
	GetPostLikes(postID int) (int, error)
	GetPostDislikes(postID int) (int, error)
}

func NewPostGetUsecase(postRepository repository.PostRepository, userRepository repository.UserRepository, categoryRepository repository.CategoryRepository) PostGetUsecase {
	return &postGetUsecase{
		postRepository:     postRepository,
		userRepository:     userRepository,
		categoryRepository: categoryRepository,
	}
}

func (pdu *postGetUsecase) GetPost(postID int) (*models.Post, error) {
	post, err := pdu.postRepository.GetPost(postID)
	if err != nil {
		return nil, err
	}

	user, err := pdu.userRepository.GetUser(post.User.ID)
	if err != nil {
		return nil, err
	}

	post.User = user

	err = pdu.categoryRepository.GetCategoriesForPost(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (pdu *postGetUsecase) GetPostLikes(postID int) (int, error) {
	likes, err := pdu.reactionRepository.GetPostLikes(postID)
	if err != nil {
		return 0, err
	}

	return likes, nil
}

func (pdu *postGetUsecase) GetPostDislikes(postID int) (int, error) {
	dislikes, err := pdu.reactionRepository.GetPostDislikes(postID)
	if err != nil {
		return 0, err
	}

	return dislikes, nil
}

func (pdu *postGetUsecase) GetCommentsByPostId(post_id int) ([]*models.Comment, error) {
	comments, err := pdu.commentRepository.GetPostComments(post_id)
	if err != nil {
		return nil, err
	}

	for _, comment := range comments {
		replies, err := pdu.commentRepository.GetCommentRepliesCount(comment.ID)
		if err != nil {
			return nil, err
		}

		comment.ReplyCount = replies
	}

	return comments, nil
}
