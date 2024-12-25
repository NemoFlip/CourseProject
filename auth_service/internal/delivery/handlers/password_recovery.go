package handlers

import (
	"CourseProject/auth_service/internal/database"
	customLogger "CourseProject/auth_service/pkg/log"
	"CourseProject/auth_service/pkg/managers"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type PassRecoveryServer struct {
	verifyCodeStorage *database.VerifyCodeStorage
	emailManager      *managers.EmailManager
	userStorage       database.UserStorage
	logger            *customLogger.Logger

	verifyCodePrefix  string
	resetStatusPrefix string
}

func NewPassRecoveryServer(verifyCodeStorage *database.VerifyCodeStorage, userStorage database.UserStorage, logger *customLogger.Logger, emailManager *managers.EmailManager) *PassRecoveryServer {
	verifyCodePrefix := "verify_code"
	resetStatusPrefix := "reset_status"
	return &PassRecoveryServer{
		verifyCodeStorage: verifyCodeStorage,
		emailManager:      emailManager,
		userStorage:       userStorage,
		logger:            logger,
		verifyCodePrefix:  verifyCodePrefix,
		resetStatusPrefix: resetStatusPrefix,
	}
}

type inputRecovery struct {
	Email string `json:"email"`
}

// @Summary Recover password
// @Description recover your password by email code
// @Tags recovery
// @Accept json
// @Produce json
// @Param email body inputRecovery true "email of the user"
// @Success 200 {nil} nil "code was sent"
// @Failure 400 {nil} nil "invalid email"
// @Router /password/request-reset [post]
func (rs *PassRecoveryServer) PasswordRecovery(ctx *gin.Context) {
	var input inputRecovery
	if err := ctx.BindJSON(&input); err != nil {
		rs.logger.ErrorLogger.Error().Msgf("unable to get email: %s", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := rs.userStorage.GetByEmail(input.Email)
	if err != nil {
		rs.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if rs.emailManager != nil {
		generatedCode := rs.emailManager.GenerateVerifyCode()
		expTime := time.Now().Add(time.Minute * 15).UTC()
		if err = rs.verifyCodeStorage.PostWithPrefix("verify_code", input.Email, generatedCode, expTime); err != nil {
			rs.logger.ErrorLogger.Error().Msg(err.Error())
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = rs.emailManager.SendCode(user.Username, input.Email, generatedCode); err != nil {
			rs.logger.ErrorLogger.Error().Msg(err.Error())
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"email": input.Email,
		})
	} else {
		rs.logger.ErrorLogger.Error().Msg("email manager is nil")
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type codeInput struct {
	Email string
	Code  string
}

// @Summary Validate Code
// @Description compare passed code with the saved one
// @Tags recovery
// @Accept json
// @Produce json
// @Param code body codeInput true "code and email of the user for recovery"
// @Success 200 {nil} nil "code is valid"
// @Failure 400 {nil} nil "invalid code"
// @Router /password/validate-code [post]
func (rs *PassRecoveryServer) ValidateCode(ctx *gin.Context) {
	var input codeInput
	if err := ctx.BindJSON(&input); err != nil {
		rs.logger.ErrorLogger.Error().Msgf("unable to get input for verification: %s", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	codeFromDB, err := rs.verifyCodeStorage.GetWithPrefix(rs.verifyCodePrefix, input.Email)
	if err != nil {
		rs.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if input.Code != codeFromDB {
		rs.logger.ErrorLogger.Error().Msg("invalid code from user")
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = rs.verifyCodeStorage.DeleteWithPrefix(rs.verifyCodePrefix, input.Email); err != nil {
		rs.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	resetExpTime := time.Now().Add(10 * time.Minute).UTC()
	if err = rs.verifyCodeStorage.PostWithPrefix(rs.resetStatusPrefix, input.Email, "valid", resetExpTime); err != nil {
		rs.logger.ErrorLogger.Error().Msg(err.Error())
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

// @Summary Reset Password
// @Description update password for registered user
// @Tags recovery
// @Accept json
// @Param password body passwordInput true "password and email of the user for recovery"
// @Success 200 {nil} nil "password is valid"
// @Failure 400 {nil} nil "invalid code"
// @Router /password/reset [post]
func (rs *PassRecoveryServer) ResetPassword(ctx *gin.Context) {
	var input passwordInput
	if err := ctx.BindJSON(&input); err != nil {
		rs.logger.ErrorLogger.Error().Msgf("unable to get input for updating password: %s", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := rs.verifyCodeStorage.GetWithPrefix(rs.resetStatusPrefix, input.Email); err != nil {
		if errors.Is(err, redis.Nil) {
			rs.logger.ErrorLogger.Error().Msg("unable to find this email with reset_status prefix in redis: pair is not set")
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		rs.logger.ErrorLogger.Error().Msgf("error fetching reset status: %s", err)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := rs.verifyCodeStorage.DeleteWithPrefix(rs.resetStatusPrefix, input.Email); err != nil {
		rs.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		rs.logger.ErrorLogger.Error().Msgf("unable to hash the password: %s", err)
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = rs.userStorage.Update(input.Email, string(hashedPassword)); err != nil {
		rs.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}
