package users

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"

	jwt "github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("")

var expirationTime = 60 // minutes

var loggedOut = make([]string, 0)

func GenerateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	currentTime := time.Now()
	expireTime := currentTime.Add(time.Duration(expirationTime) * time.Minute)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["authorized"] = true
	claims["expiresAt"] = expireTime.Unix()

	if len(secretKey) == 0 {
		secretKey = []byte(os.Getenv("JWT_SECRET"))
	}

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

func getClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			//log.Println("Token failed parsing on signing method.")
			return nil, fmt.Errorf("token failed parsing")
		}

		if len(secretKey) == 0 {
			secretKey = []byte(os.Getenv("JWT_SECRET"))
		}
		return secretKey, nil
	})

	if token == nil {
		//log.Println("Token error.")
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		//log.Println("Cannot parse claims")
		return nil, fmt.Errorf("cannot parse claims")
	}

	expireTime := claims["expiresAt"].(float64)
	if time.Now().Unix() >= int64(expireTime) {
		//log.Println("Token expired")
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}

func VerifyAuthorizationToken(tokenString string) bool {
	return strings.HasPrefix(tokenString, "Bearer ")
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func VerifyJWT(w http.ResponseWriter, r *http.Request, resp map[string]any) jwt.MapClaims {
	if r.Header["Authorization"] == nil {
		resp["error"] = []string{"Authorization header not found"}
		w.WriteHeader(http.StatusForbidden)
		return nil
	} else {
		tokenString := r.Header["Authorization"][0]
		if ok := VerifyAuthorizationToken(tokenString); !ok {
			resp["error"] = "Token doesn't start with 'Bearer '. Token incorrect."
			w.WriteHeader(http.StatusNotAcceptable)
			return nil
		} else {
			tokenClear := tokenString[7:]
			if contains(loggedOut, tokenClear) {
				resp["error"] = []string{"Token logged out"}
				w.WriteHeader(http.StatusForbidden)
				return nil
			} else {
				claims, err := getClaims(tokenClear)
				if err != nil {
					var errors = []string{"Errors while parsing token."}
					errors = append(errors, err.Error())
					resp["error"] = errors
					w.WriteHeader(http.StatusForbidden)
					return nil
				} else {
					return claims
				}
			}
		}
	}
}

func LogoutJWT(w http.ResponseWriter, r *http.Request, resp map[string]any) {
	if r.Header["Authorization"] == nil {
		resp["error"] = []string{"Authorization header not found"}
		w.WriteHeader(http.StatusForbidden)
	} else {
		tokenString := r.Header["Authorization"][0]
		if ok := VerifyAuthorizationToken(tokenString); !ok {
			resp["error"] = "Token doesn't start with 'Bearer '. Token incorrect."
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			tokenClear := tokenString[7:]
			_, err := getClaims(tokenClear)
			if err != nil {
				var errors = []string{"Errors while parsing token."}
				errors = append(errors, err.Error())
				resp["error"] = errors
				w.WriteHeader(http.StatusForbidden)
			} else {
				if !contains(loggedOut, tokenClear) {
					loggedOut = append(loggedOut, tokenClear)
					resp["status"] = "Success"
				} else {
					resp["error"] = []string{"Already logged out"}
				}
			}
		}
	}

}

func VerifyCredentials(userData userRepo.User) (userRepo.User, []error) {
	user, _ := userRepo.GetUser(userData)

	var errors []error
	if user.Email != userData.Email {
		errors = append(errors, fmt.Errorf("user not found"))
		return userRepo.User{}, errors
	}

	if user.Password != userData.Password {
		errors = append(errors, fmt.Errorf("password incorrect"))
		return userRepo.User{}, errors
	}

	if len(errors) != 0 {
		return userRepo.User{}, errors
	}

	user.Password = ""
	return user, nil

}

func LogoutMonitor() {
	log.Println("Logout monitor started")
	for {
		var temp []string
		for _, v := range loggedOut {
			_, err := getClaims(v)

			if err == nil {
				temp = append(temp, v)
			}
		}
		loggedOut = temp
		time.Sleep(30 * time.Second)
	}
}
