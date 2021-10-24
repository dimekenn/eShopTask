package middlewares

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type AuthorizationMiddleware struct {

}

func NewAuthorizationMiddleware() *AuthorizationMiddleware {
	return &AuthorizationMiddleware{}
}

func (a *AuthorizationMiddleware) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.Contains(c.Request().RequestURI, "/signin"){
			return next(c)
		}

		if strings.Contains(c.Request().RequestURI, "/signup"){
			return next(c)
		}

		r := c.Request()

		tokenStr := r.Header.Get("Authorization")
		if tokenStr == ""{
			return echo.NewHTTPError(http.StatusUnauthorized, "token is missing")
		}

		if len(tokenStr) < len("Bearer"){
			return echo.NewHTTPError(http.StatusUnauthorized, "token is missing")
		}

		var parts []string
		if strings.Index(tokenStr, " ") > 0 {
			parts = strings.Split(tokenStr, " ")
		} else {
			parts = strings.Fields(tokenStr)
		}

		token := strings.TrimSpace(parts[1])

		userId, username, roles, err := ValidateToken(token)
		if err != nil{
			return err
		}

		rolesInt := roles.([]interface{})
		rolesStr := make([]string, len(rolesInt))
		for i, v := range rolesInt{
			rolesStr[i] = v.(string)
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		ctx = context.WithValue(ctx, "roles", rolesStr)
		ctx = context.WithValue(ctx, "username", username)
		c.SetRequest(r.WithContext(ctx))

		return next(c)
	}
}

func ValidateToken(str string) (interface{}, interface{}, interface{}, error)  {
	token, err := jwt.Parse(
		str,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
			}
			return []byte("secret_key"), nil
		},
	)
	if err != nil{
		return 0, "",nil, echo.NewHTTPError(http.StatusUnauthorized, err)
	}

	if !token.Valid{
		return 0, "",nil, echo.NewHTTPError(http.StatusUnauthorized, "token is not valid")
	}

	claims, claimsErr := extractClaims(token)
	if claimsErr != nil{
		return 0, "", nil, claimsErr
	}

	mapClaims := claims.(jwt.MapClaims)
	userId := mapClaims["user_id"]
	roles := mapClaims["role"]
	username := mapClaims["username"]

	if userId == "" || username == ""{
		return 0, "", nil, echo.NewHTTPError(http.StatusUnauthorized, "token is not valis, userId or role is empty")
	}
	return userId, username, roles, nil
}

func extractClaims(token *jwt.Token) (jwt.Claims, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok{
		return claims, nil
	}else{
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "cant extract claims from token")
	}
}
