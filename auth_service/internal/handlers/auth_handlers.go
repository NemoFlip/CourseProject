package handlers

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/entity"
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
		log.Printf("unable to parse user from JSON: %s", err)
		return
	}
	newUser.ID = uuid.New().String()
	if err := us.userStorage.Post(newUser); err != nil {
		log.Println(err)
		return
	}
	ctx.Status(http.StatusCreated)
}

func LoginUser(ctx *gin.Context) {

}
