package v1

import (
	"minjust-website/internal/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	usecase domain.EmployeeUsecase
}

func NewEmployeeHandler(u domain.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{usecase: u}
}

func (h *EmployeeHandler) CreateAccountH(c *gin.Context) {
	var input domain.EmployeeAccount
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	if err := h.usecase.CreateAccount(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "учетная запись сотрудника успешно создана",
		"employee_id": input.ID,
	})
}

func (h *EmployeeHandler) GetHandbookH(c *gin.Context) {
	managers, err := h.usecase.GetManagementHandbook()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить справочник"})
		return
	}

	if managers == nil {
		managers = []domain.Employee{}
	}

	c.JSON(http.StatusOK, managers)
}
func (h *EmployeeHandler) GetAllAccountsH(c *gin.Context) {
	accounts, err := h.usecase.GetAllAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить учетные записи сотрудников"})
		return
	}

	if accounts == nil {
		accounts = []domain.EmployeeAccount{}
	}

	c.JSON(http.StatusOK, accounts)
}
func (h *EmployeeHandler) GetProfileH(c *gin.Context) {
	userIDInterface, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не авторизован"})
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		if floatVal, isFloat := userIDInterface.(float64); isFloat {
			userID = int64(floatVal)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный тип userId в контексте"})
			return
		}
	}

	profile, err := h.usecase.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить профиль из БД: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}
