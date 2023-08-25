package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type userSignUpInput struct {
	Name     string `json:"name" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=8,max=64"`
}

type receiveInput struct {
	guid string `json:"guid"`
}

type tokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type refreshInput struct {
	Token string `json:"token" binding:"required"`
}

// func (h *Handler) signUp(c *gin.Context) {
// 	var inp userSignUpInput
// 	if err := c.BindJSON(&inp); err != nil {
// 		newResponse(c, http.StatusBadRequest, "invalid input body")

// 		return
// 	}

// 	fmt.Println(inp)
// 	if err := h.services.Users.SignUp(c.Request.Context(), service.UserSignUpInput{
// 		Name:     inp.Name,
// 		Email:    inp.Email,
// 		Password: inp.Password,
// 	}); err != nil {
// 		if errors.Is(err, domain.ErrUserAlreadyExists) {
// 			newResponse(c, http.StatusBadRequest, err.Error())

// 			return
// 		}

// 		newResponse(c, http.StatusInternalServerError, err.Error())

// 		return
// 	}

// 	c.Status(http.StatusCreated)
// }

// func (h *Handler) signIn(c *gin.Context) {
// 	var inp signInInput

// 	if err := c.BindJSON(&inp); err != nil {
// 		newResponse(c, http.StatusBadRequest, "invalid input body")

// 		return
// 	}

// 	id, err := getUserId(c)
// 	if err != nil {
// 		newResponse(c, http.StatusInternalServerError, err.Error())

// 		return
// 	}

// 	res, err := h.services.Users.SignIn(c.Request.Context(), service.UserSignInInput{
// 		Email:    inp.Email,
// 		Password: inp.Password,
// 	}, id)
// 	if err != nil {
// 		if errors.Is(err, domain.ErrUserNotFound) {
// 			newResponse(c, http.StatusBadRequest, err.Error())

// 			return
// 		}

// 		newResponse(c, http.StatusInternalServerError, err.Error())

// 		return
// 	}

// 	c.JSON(http.StatusOK, tokenResponse{
// 		AccessToken:  res.AccessToken,
// 		RefreshToken: res.RefreshToken,
// 	})
// }

func (h *Handler) receive(c *gin.Context) {
	var input receiveInput
	input.guid = c.Query("guid")
	if input.guid == "" {
		newResponse(c, http.StatusBadRequest, "url param is missing")

		return
	}

	res, err := h.services.Users.GetTokens(c.Request.Context(), input.guid)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.SetCookie("access", res.AccessToken, int(h.accessTokenTTL), "/", "localhost", true, true)
	c.SetCookie("refresh", res.RefreshToken, int(h.refreshTokenTTL), "/", "localhost", true, true)

	c.JSON(http.StatusOK, "authentication was successful!")
}

func (h *Handler) refresh(c *gin.Context) {
	var tokens tokenResponse
	var err error
	tokens.AccessToken, err = c.Cookie("access")
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())

		return
	}

	tokens.RefreshToken, err = c.Cookie("refresh")
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())

		return
	}

	newTokens, err := h.services.Users.RefreshTokens(c, tokens.AccessToken, tokens.RefreshToken)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())

		return
	}

	c.SetCookie("access", newTokens.AccessToken, int(h.accessTokenTTL), "/", "localhost", true, true)
	c.SetCookie("refresh", newTokens.RefreshToken, int(h.refreshTokenTTL), "/", "localhost", true, true)
}
