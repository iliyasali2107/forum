package service

import "forum/internal/repository"

type VoteService interface{}

type voteService struct {
	Repository repository.VoteRepository
}

func NewVoteService(repository repository.VoteRepository) VoteService {
	return &voteService{
		Repository: repository,
	}
}
