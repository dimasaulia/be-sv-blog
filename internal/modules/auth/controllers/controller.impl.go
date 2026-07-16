package controllers

import (
	"net/http"

	"github.com/sv-blog/internal/modules/auth/dto"
	"github.com/sv-blog/internal/platform/logger"
	"github.com/sv-blog/internal/shared/requestctx"
	"github.com/sv-blog/internal/shared/response"
)

type AuthControllerImpl struct {
	response *response.Sender
	log      *logger.LayerLogger
}

func NewAuthController(sender *response.Sender, appLogger *logger.Logger) AuthController {
	return &AuthControllerImpl{
		response: sender,
		log:      appLogger.Layer("controller.auth"),
	}
}

func (c *AuthControllerImpl) CurrentUser(w http.ResponseWriter, r *http.Request) {
	end := c.log.Start(r.Context(), "CurrentUser")

	user, ok := requestctx.CurrentUser(r.Context())
	if !ok {
		end(nil)
		c.response.Error(w, r, http.StatusUnauthorized, "auth.unauthorized", nil)
		return
	}

	end(nil)
	c.response.Success(w, r, http.StatusOK, "auth.me.success", dto.CurrentUserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		DisplayName: user.DisplayName,
	})
}
