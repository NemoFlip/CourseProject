package handlers

import (
	"CourseProject/auth_service/internal/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type PassRecoveryServer struct {
	verifyCodeStorage *database.VerifyCodeStorage
	userStorage       database.UserStorage
}

func NewPassRecoveryServer(verifyCodeStorage *database.VerifyCodeStorage, userStorage database.UserStorage) *PassRecoveryServer {
	return &PassRecoveryServer{verifyCodeStorage: verifyCodeStorage, userStorage: userStorage}
}

type codeInput struct {
	Email string
	Code  string
}

// @Summary Verify Code
// @Description compare passed code with the saved one
// @Tags recovery
// @Accept json
// @Produce json
// @Param code body codeInput true "code and email of the user for recovery"
// @Success 200 {nil} nil "code is valid"
// @Failure 400 {nil} nil "invalid code"
// @Router /password/verify [post]
func (rs *PassRecoveryServer) VerifyCode(ctx *gin.Context) {
	var input codeInput
	if err := ctx.BindJSON(&input); err != nil {
		log.Printf("unable to get input for verification: %s", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	codeFromDB, err := rs.verifyCodeStorage.Get(input.Email)
	if err != nil {
		log.Println(err)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if input.Code != codeFromDB {
		log.Println("invalid code from user")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = rs.verifyCodeStorage.Delete(input.Email); err != nil {
		log.Println(err)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"email": input.Email,
	})

}

type passwordInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Update Password
// @Description update password for registered user
// @Tags recovery
// @Accept json
// @Param password body passwordInput true "password and email of the user for recovery"
// @Success 200 {nil} nil "password is valid"
// @Failure 400 {nil} nil "invalid code"
// @Router /password/update [post]
func (rs *PassRecoveryServer) UpdatePassword(ctx *gin.Context) {
	var input passwordInput
	if err := ctx.BindJSON(&input); err != nil {
		log.Printf("unable to get input for updating password: %s", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("unable to hash the password")
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = rs.userStorage.Update(input.Email, string(hashedPassword)); err != nil {
		log.Println(err)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}
