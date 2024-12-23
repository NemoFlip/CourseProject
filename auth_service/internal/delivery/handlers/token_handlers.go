package handlers

//
//type TokenServer struct {
//}

// @Summary Refresh tokens
// @Description get access and refresh tokens via user_id
// @Tags tokens
// @Param token body refreshInput true "Данные для регистрации пользователя"
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 500
// @Router /refresh [post]
//func (as *A) RefreshTokens(ctx *gin.Context) {
//	userID, inputIP, ok := as.tokenManager.GetClaims(ctx)
//	if !ok {
//		log.Println("unable to get user_id and input_ip from context")
//		ctx.Writer.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	// Get token from request body for refreshing
//	var inputToken refreshInput
//	if err := ctx.BindJSON(&inputToken); err != nil {
//		fmt.Printf("invalid input refresh token")
//		ctx.Writer.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	// Compare it with already saved token
//	as.compareTokens(ctx, userID, inputToken)
//
//	ip := ctx.ClientIP()
//	if ip != inputIP {
//		if as.mailSender == nil {
//			log.Printf("unable to send message to user: mailSender is not created")
//			ctx.Writer.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//
//		err := as.mailSender.SendMessage(ip)
//		if err != nil {
//			log.Println(err)
//			ctx.Writer.WriteHeader(http.StatusInternalServerError)
//		}
//		ctx.Writer.WriteHeader(http.StatusBadRequest)
//		return
//	}
//
//	as.tokenManager.GenerateBothTokens(ctx, as.tokenStorage, userID, ip)
//}
