package register

import (
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type RegisterHandler struct {
	registerUserUseCase *RegisterUserUseCase
}

func NewRegisterHandler(useCase *RegisterUserUseCase) *RegisterHandler {
	return &RegisterHandler{registerUserUseCase: useCase}
}

type RegisterRequest struct {
	Email                string `json:"email" binding:"required"`
	FirstName            string `json:"firstName" binding:"required"`
	LastName             string `json:"lastName" binding:"required"`
	Pseudo               string `json:"pseudo" binding:"required"`
	Password             string `json:"password" binding:"required,min=3"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,eqfield=Password"`
	BirthDate            string `json:"birthDate" binding:"required,ISO8601date"`
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

	userRequest, err := l.RegisterUserRequest(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid birth date format. Must be in RFC3339 format.",
		})
		return
	}

	err = l.registerUserUseCase.Execute(c, *userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (l RegisterHandler) RegisterUserRequest(request RegisterRequest) (*RegisterUserRequest, error) {
	parsedBirthDate, err := time.Parse("2006-01-02", request.BirthDate)
	if err != nil {
		return nil, err
	}

	return &RegisterUserRequest{
		Email:     request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Pseudo:    request.Pseudo,
		Password:  request.Password,
		BirthDate: parsedBirthDate,
	}, nil
}
