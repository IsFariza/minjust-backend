package v1

import (
	"minjust-website/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RequestHandler struct {
	usecase domain.PasswordRequestUsecase
}

func NewRequestHandler(u domain.PasswordRequestUsecase) *RequestHandler {
	return &RequestHandler{usecase: u}
}

func (h *RequestHandler) RequestResetPasswordH(c *gin.Context) {
	var req domain.PasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	if err := h.usecase.RequestResetPassword(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "password reset request created",
		"request_id": req.ID,
		"status":     req.Status,
	})
}

func (h *RequestHandler) CheckStatusH(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	req, err := h.usecase.GetRequestStatus(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

type ProcessRequestInput struct {
	Status            string `json:"status" binding:"required"`
	GeneratedPassword string `json:"generated_password"`
	AdminComment      string `json:"admin_comment"`
}

func (h *RequestHandler) ProcessRequestH(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	var input ProcessRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid process request format"})
		return
	}

	err = h.usecase.ProcessRequest(id, input.Status, input.GeneratedPassword, input.AdminComment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "request processed by admin"})
}
