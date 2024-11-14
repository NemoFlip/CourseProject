package handlers

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/entity"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

type UserServer struct {
	userStorage database.UserStorage
}

func NewUserServer(userStorage database.UserStorage) *UserServer {
	return &UserServer{userStorage: userStorage}
}

func (us *UserServer) RegisterHandler(ctx *gin.Context) {
	var newUser entity.User
	if err := ctx.BindJSON(&newUser); err != nil {
		log.Printf("unable to parse user from JSON: %s\n", err)
		return
	}
	newUser.ID = uuid.New().String()
	if err := us.userStorage.Post(newUser); err != nil {
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "id": newUser.ID})
}

func (us *UserServer) LoginUser(ctx *gin.Context) {
	var user entity.User
	if err := ctx.BindJSON(&user); err != nil {
		log.Printf("unable to read user from context for login: %s\n", err)
		return
	}
	userFromDB, err := us.userStorage.Get(user.Username)
	if err != nil {
		log.Println(err)
		return
	}
	if user.Password != userFromDB.Password {
		log.Printf("incorrect password")
		return
	}
	payload := jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	jwtSecret, exists := os.LookupEnv("JWT_SECRET_KEY")
	if !exists {
		log.Printf("JWT_SECRET_KEY is not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	jwtTokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("unable to sign jwt token: %s", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", jwtTokenString))
	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successfully",
		"token":   jwtTokenString,
	})
}
