package domain

type Admin struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type AuthRepository interface {
	GetByUsername(username string) (*Admin, error)
}

type AuthUsecase interface {
	Login(username, password, role string) (string, error)
}
