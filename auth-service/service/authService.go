package service

import (
	"auth-service/dto"
	"auth-service/model"
	"auth-service/repository"
)

type AuthService interface {
	UserRegister(registerRequestDTO dto.UserRegisterRequestDTO) (model.User, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (s *authService) UserRegister(registerRequestDTO dto.UserRegisterRequestDTO) (model.User, error) {
	user := &model.User{
		Email: registerRequestDTO.Email,
		Name: registerRequestDTO.Name,
		Cpf: registerRequestDTO.Cpf,
		Password: registerRequestDTO.Password,
	}
	err := s.userRepository.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}
	return *user, nil
}

