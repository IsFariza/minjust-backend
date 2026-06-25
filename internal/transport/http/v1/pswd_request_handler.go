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

func (h *RequestHandler) CreateRequestH(c *gin.Context) {
	empIDVal, ok1 := c.Get("userId")
	userIINVal, ok2 := c.Get("userIIN")

	if !ok1 || !ok2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неавторизованный запрос"})
		return
	}

	empID, okId := empIDVal.(int64)
	userIIN, okIin := userIINVal.(string)

	if !okId || !okIin {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка обработки данных авторизации"})
		return
	}

	var input domain.CreatePasswordRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	err := h.usecase.CreateRequest(empID, userIIN, input.EmployeeIIN, input.SystemName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "заявка успешно создана"})
}

func (h *RequestHandler) GetMyRequestsH(c *gin.Context) {
	empID, _ := c.Get("userId")

	requests, err := h.usecase.GetEmployeeRequests(empID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

func (h *RequestHandler) GetAllRequestsH(c *gin.Context) {
	role, exists := c.Get("userRole")
	if !exists || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "доступ запрещен"})
		return
	}

	requests, err := h.usecase.GetAllRequests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, requests)
}

func (h *RequestHandler) ReviewRequestH(c *gin.Context) {
	role, _ := c.Get("userRole")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "доступ запрещен"})
		return
	}

	reqID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var input domain.UpdatePasswordStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	err = h.usecase.ReviewRequest(reqID, input.Status, input.RejectionReason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "статус заявки обновлен"})
}
