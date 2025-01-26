package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shreyanshsharma88/golang-bank/auth"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker auth.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authParts := strings.Fields(authHeader)
		if len(authParts) != 2 {
			err := errors.New("authorization header is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		if strings.ToLower(authParts[0]) != authorizationTypeBearer {
			err := errors.New("authorization type is not valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authToken := authParts[1]
		payload, err := tokenMaker.VerifyToken(authToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		c.Set(authorizationPayloadKey, payload)
		c.Next()

	}
}
