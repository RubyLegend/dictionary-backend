package users

import (
  "time"
  "strings"
  "fmt"
  "log"
  "net/http"

  userRepo "github.com/RubyLegend/dictionary-backend/repository/users"

  jwt "github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("sampleSecretKeyYouShouldNeverShare")

var expirationTime = 10 // minutes

func GenerateJWT(username string) (string, error) {
  token := jwt.New(jwt.SigningMethodHS512)
  currentTime := time.Now()
  expireTime := currentTime.Add(time.Duration(expirationTime) * time.Minute)

  claims := token.Claims.(jwt.MapClaims)
  claims["username"] = username
  claims["authorized"] = true
  claims["expiresAt"] = expireTime.Unix()

  tokenString, err := token.SignedString(secretKey)

  if err != nil {
    log.Println(err)
    return "", err
  }
  
  log.Println("User " + username + " authorized.")
  return tokenString, nil
}

func getClaims(tokenString string) (jwt.MapClaims, error) {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      //log.Println("Token failed parsing on signing method.")
      return nil, fmt.Errorf("Token failed parsing.")
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
    return nil, fmt.Errorf("Cannot parse claims")
  }

  expireTime := claims["expiresAt"].(float64)
  if time.Now().Unix() >= int64(expireTime) {
    //log.Println("Token expired")
    return nil, fmt.Errorf("Token expired")
  }
  
  return claims, nil
}

func verifyAuthorizationToken(tokenString string) bool {
  return strings.HasPrefix(tokenString, "Bearer ")
}

func VerifyJWT(w http.ResponseWriter, r *http.Request, resp map[string]any) (jwt.MapClaims) {
	if r.Header["Authorization"] == nil {
	  resp["error"] = []string{"Authorization header not found"}
    w.WriteHeader(http.StatusForbidden)
    return nil
	} else {
	  tokenString := r.Header["Authorization"][0]
	  if ok := verifyAuthorizationToken(tokenString); !ok {
  		resp["error"] = "Token doesn't start with 'Bearer '. Token incorrect."
      w.WriteHeader(http.StatusNotAcceptable)
      return nil
	  } else { 
	  	tokenClear := tokenString[7:]
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

func VerifyCredentials(userData userRepo.User) (userRepo.User, []error) {
  user, _ := userRepo.GetUser(userData)

  var errors []error
  if user.Username != userData.Username {
    errors = append(errors, fmt.Errorf("User not found."))
  }

  if user.Password != userData.Password {
    errors = append(errors, fmt.Errorf("Password incorrect."))
  }

  if len(errors) != 0 {
    return userRepo.User{}, errors
  }
  
  user.Password = ""
  return user, nil

}
