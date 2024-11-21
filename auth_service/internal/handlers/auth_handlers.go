package handlers

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/entity"
	"CourseProject/auth_service/pkg/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

type UserServer struct {
	userStorage  database.UserStorage
	tokenManager auth.TokenManager
}

func NewUserServer(userStorage database.UserStorage, tokenManager auth.TokenManager) *UserServer {
	return &UserServer{userStorage: userStorage, tokenManager: tokenManager}
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
// @Router /registration [post]
func (us *UserServer) RegisterUser(ctx *gin.Context) {
	var newUser entity.User
	if err := ctx.BindJSON(&newUser); err != nil {
		log.Printf("unable to parse user from JSON: %s\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user's data"})
		return
	}
	newUser.ID = uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 4)
	if err != nil {
		log.Printf("unable to hash the password")
		return
	}
	newUser.Password = string(hashedPassword)
	if err := us.userStorage.Post(newUser); err != nil {
		log.Println(err)
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
// @Router /login [post]
func (us *UserServer) LoginUser(ctx *gin.Context) {
	var user entity.User
	if err := ctx.BindJSON(&user); err != nil {
		log.Printf("unable to read user from context for login: %s\n", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	userFromDB, err := us.userStorage.Get(user.Username)
	if err != nil {
		log.Printf("unable to find user in database: %s", err)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFromDB.Password), []byte(user.Password))
	if err != nil {
		log.Printf("incorrect password")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	jwtTokenString, err := us.tokenManager.SignToken(token)
	if err != nil {
		log.Println(err)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", jwtTokenString))
	//ctx.JSON(http.StatusOK, gin.H{
	//	"message": "login successfully",
	//	"token":   jwtTokenString,
	//})
}

func (us *UserServer) LogoutUser(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")

	if len(authHeader) == 0 {
		log.Printf("authroization token is absent")
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	bearerToken := strings.Split(authHeader, " ")
	if bearerToken[0] != "Bearer" || len(bearerToken) != 2 {
		log.Printf("invalid format of auth token")
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString := bearerToken[1]
	err := us.tokenManager.ValidateToken(tokenString)
	if err != nil {
		log.Println(err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
}
