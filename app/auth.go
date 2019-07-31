package app

import (
	"context"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"

	"rest_api_gorm/utils"

	"rest_api_gorm/models"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		//List of endpoints that dont require auth:
		notAuth := []string{"/api/user/new", "/api/user/login", "/"}
		requestPath := req.URL.Path //current request path

		//chenck if request does not need authentication, serve the request if it doesnt need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(res, req)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := req.Header.Get("Authorization") //grab the token from the header

		//if token is missing, return 403
		if tokenHeader == "" {
			response = utils.Message(false, "Missing auth token")
			res.WriteHeader(http.StatusForbidden)
			res.Header().Add("Content-Type", "application/json")
			utils.Respond(res, response)
			return
		}

		//The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			response = utils.Message(false, "Invalid/Malformed auth token")
			res.WriteHeader(http.StatusForbidden)
			res.Header().Add("Content-Type", "application/json")
			utils.Respond(res, response)
			return
		}
		//get only the token, which is the second guy:
		tokenPart := splitted[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		//Malformed token, returns with http code 403 as usual
		if err != nil {
			response = utils.Message(false, "Malformed authentication token")
			res.WriteHeader(http.StatusForbidden)
			res.Header().Add("Content-Type", "application/json")
			utils.Respond(res, response)
			return
		}
		//Token is invalid, maybe not signed in the server
		if !token.Valid {
			response = utils.Message(false, "Token is not valid")
			res.WriteHeader(http.StatusForbidden)
			res.Header().Add("Content-Type", "application/json")
			utils.Respond(res, response)
			return
		}
		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		// fmt.Sprintf("User %", tk.username)
		ctx := context.WithValue(req.Context(), "user", tk.UserId)
		req = req.WithContext(ctx)
		next.ServeHTTP(res, req) ////proceed in the middleware chain!
	})
}
