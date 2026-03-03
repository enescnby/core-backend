package services

import (
	"core-backend/internal/dto"
	"core-backend/internal/repositories"
)

type UserService interface {
	GetUserForLookup(coreGuardID string) (*dto.LookupResponse, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserForLookup(coreGuardID string) (*dto.LookupResponse, error) {
	user, err := s.repo.GetUserForLookup(coreGuardID)
	if err != nil {
		return nil, err
	}

	return &dto.LookupResponse{
		UserID:    user.UserID.String(),
		PublicKey: user.Key.PublicKey,
	}, err
}
