package users

import (
  "time"
  "log"
  "errors"
)

type User struct {
  UserId int `json:"userId"`
  Email string `json:"email"`
  Username string `json:"username"`
  Password string `json:"password"`
  CreatedAt time.Time `json:"createdAt"`
}

var Users []User

func checkUserExistance(userData User) ([]error) {
  var Errors []error
  for _, v := range Users {
    if v.Email == userData.Email {
      Errors = append(Errors, errors.New("Email already registered."))
    }
    if v.Username == userData.Username {
      Errors = append(Errors, errors.New("Username already registered."))
    }
  }

  return Errors
}

func validation(userData User) ([]error) {
  var err []error

  if(len(userData.Username) == 0){
    err = append(err, errors.New("Username is required field"))
  } 
  if(len(userData.Email) == 0){
    err = append(err, errors.New("Email is required field"))
  }
  if(len(userData.Password) == 0){
    err = append(err, errors.New("Password is required field"))
  }

  return err
}

func findUser(params ...interface{}) (int, error) {
  userData, ok := params[0].(User)
  if ok == true {
    for i, v := range Users {
     if v.Email == userData.Email ||
        v.Username == userData.Username {
          return i, nil
        }
     }
  } else {
    username, ok := params[0].(string)
    if ok == true {
      for i, v := range Users {
       if v.Username == username {
            return i, nil
          }
       }
    } else {
      return -1, errors.New("Unknown parameter passed")
    }
  }

  return -1, errors.New("User not found")
}

func GetUser(userData User) (User, error) {
  realUserId, err := findUser(userData)

  if err != nil {
    log.Println(err)
    return User{}, err
  }

  return Users[realUserId], nil
}

func AddUser(userData User) ([]error) {
  var err []error

  // request validation
  err = append(err, validation(userData)...)
  
  err2 := checkUserExistance(userData)

  if err2 != nil {
    for _, v := range err2 {
      err = append(err, v)
    }
  }

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

func DeleteUser(userData User) error {
  i, err := findUser(userData)
  
  log.Println(i, err)

  if err == nil {
    Users = append(Users[:i], Users[i+1:]...)
    log.Println(Users)
    return nil
  }

  return err
}

func EditUser(currentUsername string, userData User) []error {
  var err []error
  // request validation
  err = append(err, validation(userData)...)
  
  i, err2 := findUser(currentUsername)
  
  if err2 != nil {
    err = append(err, err2)
  } else {
    user := Users[i]

    err2 := checkUserExistance(userData)

    if err2 != nil {
      for _, v := range err2 {
        if v.Error() == "Email already registered." && user.Email != userData.Email {
          err = append(err, v)
        } else if v.Error() == "Username already registere." && user.Username != userData.Username {
          err = append(err, v)
        }
      }
    }
  
    if err == nil {
      userData.CreatedAt = user.CreatedAt
      temp := append([]User{}, Users[:i]...)
      temp = append(temp, userData)
      Users = append(temp, Users[i+1:]...)
      return nil
    }
  }

  return err
}
