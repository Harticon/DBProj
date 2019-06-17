package DBproj

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
)

func UserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	return func(ctx echo.Context) error {
		t := ctx.Request().Header.Get("Authorization")
		token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			fmt.Println("You are not logged in")
			return ctx.JSON(http.StatusUnauthorized, err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if claims.Valid() != nil || !ok {
			fmt.Println("not valid claims")
			return ctx.JSON(http.StatusUnauthorized, "user does not exists")
		}

		ctx.Set("id", int(claims["id"].(float64)))

		return next(ctx)
	}

}
