package service

import (
	"auth-service/dto"
	"auth-service/model"
	"auth-service/repository"
	"auth-service/utils"
	"fmt"
	"time"
	"github.com/google/uuid"
)

type AuthService interface {
	UserRegister(registerRequestDTO dto.UserRegisterRequestDTO) (model.User, error)
	UserLogin(loginRequestDTO dto.UserLoginRequestDTO) (model.User, string, error)
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
		ID: uuid.New().String(),
		Email: registerRequestDTO.Email,
		Name: registerRequestDTO.Name,
		Cpf: registerRequestDTO.Cpf,
		Role: "USER",
		Password: hashedPassword,
		AvatarUrl: registerRequestDTO.AvatarUrl,
	}

	err = s.userRepository.CreateUser(user)
	if err != nil {
		return model.User{}, err
	}

	return *user, nil
}

func (s *authService) UserLogin(loginRequestDTO dto.UserLoginRequestDTO) (model.User, string, error) {
	user, err := s.userRepository.GetUserByEmail(loginRequestDTO.Email)
	if err != nil {
		return model.User{}, "", err
	}

	if !utils.ComparePassword(user.Password, loginRequestDTO.Password) {
		return model.User{}, "", fmt.Errorf("Senha incorreta. Tente novamente.")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, user.Name, time.Hour*168)
	if err != nil {
		return model.User{}, "", fmt.Errorf("erro ao gerar token: %w", err)
	}

	return user, token, nil
}