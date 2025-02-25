package register

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/get_data_token_invite"
	"io"
	"net/http"
)

type RegisterHandler struct {
	registerUserUseCase      *RegisterUserUseCase
	getInviteLinkDataUseCase *get_data_token_invite.GetInviteLinkDataUseCase
}

func NewRegisterHandler(useCase *RegisterUserUseCase, getInviteLinkDataUseCase *get_data_token_invite.GetInviteLinkDataUseCase) *RegisterHandler {
	return &RegisterHandler{registerUserUseCase: useCase, getInviteLinkDataUseCase: getInviteLinkDataUseCase}
}

type RegisterRequest struct {
	Token                string `json:"token" binding:"required"`
	Password             string `json:"password" binding:"required,min=3"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=Password"`
}

func (l RegisterHandler) Handle(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		if errors.Is(err, io.EOF) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Request body is empty.",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if request.Password != request.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password and password confirmation are not the same.",
		})
		return
	}

	inviteLinkData, err := l.getInviteLinkDataUseCase.GetInviteLinkData(c, request.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userRequest, err := l.RegisterUserRequest(request, inviteLinkData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Cannot parse request",
		})
		return
	}

	err = l.registerUserUseCase.Execute(c, *userRequest)
	if err != nil {
		if errors.Is(err, UserAlreadyExistsErr) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":          err.Error(),
				"messageDisplay": "Un utilisateur existe déjà avec cet email.",
				"level":          "warning",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (l RegisterHandler) RegisterUserRequest(request RegisterRequest, inviteLinkData *entity.InviteLink) (*RegisterUserRequest, error) {

	return &RegisterUserRequest{
		FirstName: inviteLinkData.FirstName,
		LastName:  inviteLinkData.LastName,
		Email:     inviteLinkData.Email,
		Password:  request.Password,
	}, nil
}
