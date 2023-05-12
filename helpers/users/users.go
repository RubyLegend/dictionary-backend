package users

import (
  "time"
  "strings"
  "fmt"
  "log"

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

func VerifyJWT(tokenString string) (string, error) {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      log.Println("Token failed parsing on signing method.")
      return nil, fmt.Errorf("Token failed parsing.")
    }

    return secretKey, nil
  })

  if token == nil {
    log.Println("Token error.")
    return "", err
  }

  claims, ok := token.Claims.(jwt.MapClaims)

  if !ok {
    log.Println("Cannot parse claims")
    return "", fmt.Errorf("Cannot parse claims")
  }

  expireTime := claims["expiresAt"].(float64)
  if time.Now().Unix() >= int64(expireTime) {
    log.Println("Token expired")
    return "", fmt.Errorf("Token expired")
  }
  
  username := claims["username"].(string)
  return username, nil
}

func VerifyAuthorizationToken(tokenString string) bool {
  return strings.HasPrefix(tokenString, "Bearer ")
}
