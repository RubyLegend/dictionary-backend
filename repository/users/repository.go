package users

import (
  "time"
  "log"
  "errors"
  "fmt"
  "bytes"
)

type User struct {
  UserId int `json:"userId"`
  Email string `json:"email"`
  Username string `json:"username"`
  Password string `json:"password"`
  CreatedAt time.Time `json:"createdAt"`
}

var Users []User

func checkUserExistance(userData User) (error) {
  for _, v := range Users {
    if v.Email == userData.Email {
      return errors.New("Email already registered.")
    } else if v.Username == userData.Username {
      return errors.New("Username already registered.")
    }
  }

  return nil
}

func validation(userData User) {
  var buf bytes.Buffer
  if(len(userData.Username) == 0){
    buf.WriteString("Username is required field\n")

  } 
  if(len(userData.Email) == 0){
    buf.WriteString("Email is required field\n")
  }
  if(len(userData.Password) == 0){
    buf.WriteString("Password is required field\n")
  }
  if buf.Len() > 0 {
    fmt.Println(buf.String())
  }
  return 
}

func findUser(userData User) (int, error) {
  for i, v := range Users {
    if v.Email == userData.Email ||
       v.Username == userData.Username {
         return i, nil
       }
  }

  return -1, errors.New("User not found")
}

func GetUser(userData User) User {
  realUserId, err := findUser(userData)

  if err != nil {
    log.Println(err)
    return User{}
  }

  return Users[realUserId]
}

func AddUser(userData User) (error) {

  // request validation
  validation(userData)
  
  err := checkUserExistance(userData)
  if err == nil {
    lastElementIndex := len(Users) - 1
    if lastElementIndex < 0 {
      userData.UserId = 0
    } else {
      userData.UserId = Users[lastElementIndex].UserId + 1
    }

    userData.CreatedAt = time.Now()
    Users = append(Users, userData)
    return nil
  } else {
    return err
  }
}

func DeleteUser(userData User) (bool, error) {
  i, err := findUser(userData)
  
  if err != nil {
    Users = append(Users[:i], Users[i+1])
    return true, nil
  }

  return false, err
}
