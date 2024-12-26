package handlers

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/entity"
	customLogger "CourseProject/auth_service/pkg/log"
	"CourseProject/auth_service/pkg/managers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type UserServer struct {
	userStorage    database.UserStorage
	refreshStorage database.RefreshStorage
	tokenManager   managers.TokenManager
	logger         *customLogger.Logger
}

func NewUserServer(userStorage database.UserStorage, tokenManager managers.TokenManager, refreshStorage database.RefreshStorage, logger *customLogger.Logger) *UserServer {
	return &UserServer{userStorage: userStorage, tokenManager: tokenManager, refreshStorage: refreshStorage, logger: logger}
}

// @Summary Register user
// @Description register user by credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entity.User true "user to register"
// @Success 201 {object} entity.AuthResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/register [post]
func (us *UserServer) RegisterUser(ctx *gin.Context) {
	var newUser entity.User
	if err := ctx.BindJSON(&newUser); err != nil {
		us.logger.ErrorLogger.Error().Msgf("unable to parse user from JSON: %s\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user's data"})
		return
	}
	newUser.ID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("unable to hash the password")
		return
	}
	newUser.Password = string(hashedPassword)
	if err = us.userStorage.Post(newUser); err != nil {
		us.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unable to post new user"})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "id": newUser.ID})
}

// @Summary Login user
// @Description login user by credentials
// @Tags auth
// @Accept json
// @Produce json
// @Param user body entity.User true "user to login"
// @Success 200 {object} entity.AuthResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /auth/login [post]
func (us *UserServer) LoginUser(ctx *gin.Context) {
	var user entity.User
	if err := ctx.BindJSON(&user); err != nil {
		us.logger.ErrorLogger.Error().Msgf("unable to read user from context for login: %s\n", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	userFromDB, err := us.userStorage.GetByName(user.Username)
	if err != nil {
		us.logger.ErrorLogger.Error().Msgf("unable to find user in database: %s", err)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password))
	if err != nil {
		us.logger.ErrorLogger.Error().Msgf("incorrect password: %s", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := us.tokenManager.GenerateBothTokens(userFromDB.ID)
	if err != nil {
		us.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	hashedRefreshToken, err := us.tokenManager.GetHashedRefreshToken(refreshToken)
	if err != nil {
		us.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	expTime := time.Now().Add(time.Minute * 43200).UTC() // 30 days refresh_token is valid
	refreshTokenEntity := entity.RefreshToken{
		UserID:       userFromDB.ID,
		RefreshToken: hashedRefreshToken,
		ExpiresAt:    expTime,
	}
	if err = us.refreshStorage.Post(refreshTokenEntity); err != nil {
		us.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// @Summary Logout user
// @Description logout user with token's validation
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {nil} nil "Token is valid"
// @Failure 401 {nil} nil "User is unauthorized"
// @Failure 400 {nil} nil "Invalid token is sent"
// @Router /auth/logout [post]
func (us *UserServer) LogoutUser(ctx *gin.Context) {
	// Delete refresh token from storage
	userID, ok := us.tokenManager.GetUserID(ctx, us.logger)
	if !ok {
		return
	}
	if err := us.refreshStorage.Delete(userID); err != nil {
		us.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Delete access_token from browser's cookie
	ctx.SetCookie("Authorization", "", -1, "/", "", false, true)
}
