package httpSender

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"server/internal/domain"
)

type Handler struct {
	senderUsecase domain.SenderUsecase
}

func NewSenderHandler(su domain.SenderUsecase) *Handler {
	return &Handler{
		senderUsecase: su,
	}
}

func (h *Handler) GetFile(c echo.Context) error {
	fileName := c.QueryParam("fileName")

	f, err := h.senderUsecase.ReadFile(fileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "file doesn't exist")
	}

	return c.Stream(http.StatusOK, "multipart/form-data", f)
}
