package users

import (
  "time"
  "strings"
  "fmt"
  "log"

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

func VerifyJWT(tokenString string) (string, error) {
  token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
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

func VerifyCredentials(userData userRepo.User) (userRepo.User, error) {
  user := userRepo.GetUser(userData)

  if user.Username != userData.Username {
    return userRepo.User{}, fmt.Errorf("User not found.")
  }

  if user.Password == userData.Password {
    user.Password = ""
    return user, nil
  }

  return userRepo.User{}, fmt.Errorf("Password incorrect.")
}
