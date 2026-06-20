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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	if err := h.usecase.CreateAccount(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "employee account created",
		"employee_id": input.ID,
	})
}

func (h *EmployeeHandler) GetHandbookH(c *gin.Context) {
	managers, err := h.usecase.GetManagementHandbook()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get handbook"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch employee accounts"})
		return
	}

	if accounts == nil {
		accounts = []domain.EmployeeAccount{}
	}

	c.JSON(http.StatusOK, accounts)
}
