package app

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"

	"github.com/sirupsen/logrus"

	"playhead/db"
	"playhead/model"
)

type Context struct {
	Logger        logrus.FieldLogger
	RemoteAddress string
	Database      *db.Database
	UserPlayhead  *model.UserPlayhead
	User          *model.User
}

func (ctx *Context) WithLogger(logger logrus.FieldLogger) *Context {
	ret := *ctx
	ret.Logger = logger
	return &ret
}

func (ctx *Context) WithRemoteAddress(address string) *Context {
	ret := *ctx
	ret.RemoteAddress = address
	return &ret
}

// func (ctx *Context) WithBearerToken(token string) (*Context, error) {
// 	user, failureReason, err := ctx.App.ValidateBearerToken(token)
// 	if err != nil {
// 		return nil, err
// 	} else if failureReason != "" {
// 		ctx.Logger.WithField("reason", failureReason).Info("bearer token validation failure")
// 		return nil, nil
// 	}
//
// 	return ctx.WithUser(user), nil
// }
type CustomClaims struct {
	jwt.Claims
	// additional claims apart from standard claims
	extra map[string]interface{}
}

func (a *App) ExtractToken(tokenString string) (*model.User, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(string(a.Config.JwtKey)), nil
	})
	fmt.Println("TOKEN: ", token)
	fmt.Println("ERR: ", err)
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val)
	}
	user := model.User{}
	return &user, nil
}
func (ctx *Context) WithUser(authString string) *Context {
	ret := *ctx
	ret.User = &model.User{}
	return &ret
}

func (ctx *Context) AuthorizationError() *UserPlayheadError {
	return &UserPlayheadError{Message: "unauthorized", StatusCode: http.StatusForbidden}
}
