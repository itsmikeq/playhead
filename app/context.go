package app

import (
	"encoding/json"
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

func (a *App) ExtractToken(tokenString string) *model.User {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(string(a.Config.JwtKey)), nil
	})
	jsonbody, err := json.Marshal(claims)
	if err != nil {
		// do error check
		fmt.Println(err)
		return nil
	}
	user := model.User{}
	// for key, val := range claims {
	// 	fmt.Printf("Key: %v, value: %v\n", key, val)
	// }
	if err := json.Unmarshal(jsonbody, &user); err != nil {
		logrus.Println(err)
	}

	// fmt.Printf("GOT USER: %+v\n", user)
	// fmt.Printf("GOT CLAIMS USER: %+v\n", claims["user_id"])
	user.UserID = fmt.Sprintf("%v", claims["user_id"])

	return &user
}
func (ctx *Context) WithUser(user *model.User) *Context {
	// fmt.Println("Got User ", user)
	ret := *ctx
	ret.User = user
	return &ret
}

func (ctx *Context) AuthorizationError() *UserPlayheadError {
	return &UserPlayheadError{Message: "unauthorized", StatusCode: http.StatusForbidden}
}
