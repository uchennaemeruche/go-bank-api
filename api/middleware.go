package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	api "github.com/uchennaemeruche/go-bank-api/api/util"
	"github.com/uchennaemeruche/go-bank-api/token"
)

const (
	authHeaderKey  = "authorization"
	bearerAuthType = "bearer"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authenticationHeader := ctx.GetHeader(authHeaderKey)
		if len(authenticationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(err))
			return
		}

		fields := strings.Fields(authenticationHeader)
		if len(fields) < 2 || len(fields[1]) == 0 {
			err := errors.New("invalid authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(err))
			return
		}

		if bearerAuthType != strings.ToLower(fields[0]) {
			err := fmt.Errorf("unsupported authentication type %s", fields[0])
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			err := errors.New("invalid token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, api.ErrorResponse(err))
			return
		}
		ctx.Set(api.AuthPayloadKey, payload)
		ctx.Next()

	}
}

/**
* Gin Middlewares are designed as higher order functions that will return an anonymous function that receives the *Context
* The anonymous function is in fact the middleware that does the job.


** Auth Middleware Implementation steps.
* Define Middlware Function that recieves a tokenMaker of type token.Maker and returns gin.HandlerFunc
	- The function returns an anonymous function that receives the *gin.Context
* Extract the authentication header from the context
* check if the authentication header is empty or null.
	- If empty, abort the process using ctx.AbortWithStatusJSON

* split the authorizationHeader using strings.Field - The strings.Field splits a string using the space delimeter and returns an slice of sptrin []string

* check if the returned split result is less than 2
	 - If less than 2, abort the process and return an error

* check the authorization type - either bearer, openssl etc.
	 - If the authorization type is not supported, abort the process and throw an error.

* call the VerifyToken function of the tokenMaker to verify the access token.
	- If it returns an error, abort the process and return error.
* If the verification in the last step is correct, extract the payload.
* set the payload in the context using an auth payload key.
* call the Next() middleware
**/
