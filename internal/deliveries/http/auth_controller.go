package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/yourusername/hotel-booking-service/internal/data"
	"github.com/yourusername/hotel-booking-service/internal/usecases"
)

type AuthController struct {
	authUsecase *usecases.AuthUsecase
}

func NewAuthController(authUsecase *usecases.AuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
	}
}

func (c *AuthController) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required,min=6"`
		}
		
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: err.Error()})
			return
		}
		
		user, err := c.authUsecase.Register(ctx, req.Email, req.Password)
		if err != nil {
			status := http.StatusInternalServerError
			if strings.Contains(err.Error(), "already exists") {
				status = http.StatusConflict
			}
			ctx.JSON(status, data.ApiError{Error: err.Error()})
			return
		}
		
		ctx.JSON(http.StatusCreated, user)
	}
}

func (c *AuthController) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req data.LoginRequest
		
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: err.Error()})
			return
		}
		
		token, user, err := c.authUsecase.Login(ctx, req.Email, req.Password)
		if err != nil {
			status := http.StatusInternalServerError
			
			if err == usecases.ErrUserNotFound {
				status = http.StatusNotFound
			} else if err == usecases.ErrInvalidCredentials {
				status = http.StatusUnauthorized
			}
			
			ctx.JSON(status, data.ApiError{Error: err.Error()})
			return
		}
		
		ctx.JSON(http.StatusOK, data.LoginResponse{
			Token: token,
			User:  *user,
		})
	}
}

func (c *AuthController) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, data.ApiError{Error: "missing authorization header"})
			return
		}
		
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, data.ApiError{Error: "invalid authorization header format"})
			return
		}
		
		claims, err := c.authUsecase.VerifyToken(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, data.ApiError{Error: "invalid token"})
			return
		}
		
		ctx.Set("userID", claims.UserID)
		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context) (int64, bool) {
	userID, exists := ctx.Get("userID")
	if !exists {
		return 0, false
	}
	
	id, ok := userID.(int64)
	return id, ok
}