package service

import (
	"auth-service/dto"
	"auth-service/model"
	"auth-service/repository"
	"auth-service/utils"
	"fmt"
)

type AuthService interface {
	UserRegister(registerRequestDTO dto.UserRegisterRequestDTO) (model.User, error)
	UserLogin(loginRequestDTO dto.UserLoginRequestDTO) (model.User, error)
}

type authService struct {
	userRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

func (s *authService) UserRegister(registerRequestDTO dto.UserRegisterRequestDTO) (model.User, error) {

	hashedPassword, err := utils.HashPassword(registerRequestDTO.Password)
	if err != nil {
		return model.User{}, fmt.Errorf("Erro ao criptografar senha: %w.", err)
	}

	user := &model.User{
		Email: registerRequestDTO.Email,
		Name: registerRequestDTO.Name,
		Cpf: registerRequestDTO.Cpf,
		Password: hashedPassword,
	}
	err = s.userRepository.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}
	return *user, nil
}

func (s *authService) UserLogin(loginRequestDTO dto.UserLoginRequestDTO) (model.User, error) {
	user, err := s.userRepository.GetUserByEmail(loginRequestDTO.Email)
	if err != nil {
		return model.User{}, err
	}

	if !utils.ComparePassword(user.Password, loginRequestDTO.Password) {
		return model.User{}, fmt.Errorf("Senha incorreta. Tente novamente.")
	}

	return user, nil
}