package handlers

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/entity"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
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
	ctx.Status(http.StatusCreated)
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

	tokenCredentials := fmt.Sprintf("%s:%s", user.Username, user.Password)
	token := base64.StdEncoding.EncodeToString([]byte(tokenCredentials))
	ctx.Writer.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// TODO: implement sessions
}
