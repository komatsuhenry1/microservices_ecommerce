package handler

import (
	"auth-service/dto"
	"auth-service/utils"
	"net/http"
	"auth-service/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) UserRegister(c *gin.Context) {
	// 1. Criar o DTO e preenchê-lo com os dados do formulário
	var userRequestDTO dto.UserRegisterRequestDTO
	// userRequestDTO.Name = c.PostForm("email")
	// userRequestDTO.Email = c.PostForm("email")
	// userRequestDTO.Cpf = c.PostForm("cpf")
	// userRequestDTO.Password = c.PostForm("password")

	if err := c.ShouldBindJSON(&userRequestDTO); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.authService.UserRegister(userRequestDTO)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "usuário criado com sucesso", gin.H{"user": createdUser})
}

func (h *AuthHandler) UserLogin(c *gin.Context) {
	// 1. Criar o DTO e preenchê-lo com os dados do formulário
	var userRequestDTO dto.UserLoginRequestDTO

	if err := c.ShouldBindJSON(&userRequestDTO); err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.authService.UserLogin(userRequestDTO)
	if err != nil {
		utils.SendErrorResponse(c, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendSuccessResponse(c, "usuário logado com sucesso", gin.H{"user": user, "token": token})
}
