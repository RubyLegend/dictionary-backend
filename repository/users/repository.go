package users

import (
	"errors"
	"log"
	"time"

	db "github.com/RubyLegend/dictionary-backend/middleware/database"
)

type User struct {
	UserId    int       `json:"userId"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

var Users []User

func checkUserExistance(userData User) []error {
	var Errors []error
	dbCon := db.GetConnection()
	var res, res2 int

	err := dbCon.QueryRow("select count(*) from Users where email = ?", userData.Email).Scan(&res)
	if err != nil {
		Errors = append(Errors, err)
		return Errors
	}

	err2 := dbCon.QueryRow("select count(*) from Users where username = ?", userData.Username).Scan(&res2)
	if err2 != nil {
		Errors = append(Errors, err2)
		return Errors
	}

	if res != 0 {
		Errors = append(Errors, errors.New("email already registered"))
	}
	if res2 != 0 {
		Errors = append(Errors, errors.New("username already registered"))
	}

	return Errors
}

func validation(userData User) []error {
	var err []error

	if len(userData.Username) == 0 {
		err = append(err, errors.New("username is required field"))
	}
	if len(userData.Email) == 0 {
		err = append(err, errors.New("email is required field"))
	}
	if len(userData.Password) == 0 {
		err = append(err, errors.New("password is required field"))
	}

	return err
}

func findUser(params ...interface{}) (interface{}, error) {
	userData, ok := params[0].(User)
	if ok {
		dbCon := db.GetConnection()
		var user User

		rows, err := dbCon.Query("select * from Users where email = ? or username = ?", userData.Email, userData.Username)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		// Under normal circumstances, there will be only one record
		rows.Next()
		rows.Scan(&user.UserId, &user.Email, &user.Username, &user.Password, &user.CreatedAt)
		return user, nil

	} else {
		username, ok := params[0].(string)
		if ok {
			dbCon := db.GetConnection()
			var user User

			rows, err := dbCon.Query("select * from Users where username = ?", username)
			if err != nil {
				return -1, err
			}

			defer rows.Close()

			// Under normal circumstances, there will be only one record
			rows.Next()
			rows.Scan(&user.UserId, &user.Email, &user.Username, &user.Password, &user.CreatedAt)
			return user, nil
		} else {
			return nil, errors.New("unknown parameter passed")
		}
	}
}

func GetUser(userData User) (User, error) {
	user, err := findUser(userData)

	if err != nil {
		log.Println(err)
		return User{}, err
	}

	return user.(User), nil
}

func AddUser(userData User) []error {
	var Errors []error

	// request validation
	Errors = append(Errors, validation(userData)...)

	Errors = append(Errors, checkUserExistance(userData)...)

	if Errors == nil {
		dbCon := db.GetConnection()

		_, err := dbCon.Exec("insert into Users values (default, ?, ?, ?, CURRENT_TIME())", userData.Email, userData.Username, userData.Password)

		if err != nil {
			Errors = append(Errors, err)
			return Errors
		}

		return nil
	} else {
		return Errors
	}
}

func DeleteUser(userData User) error {
	user, err := findUser(userData)

	if err != nil {
		return err
	}

	dbCon := db.GetConnection()

	_, err = dbCon.Exec("delete from Users where userID = ?", user.(User).UserId)

	if err != nil {
		return err
	}

	return nil
}

func EditUser(currentUsername string, userData User) []error {
	var Errors []error
	// request validation
	Errors = append(Errors, validation(userData)...)

	user, err := findUser(currentUsername)

	if err != nil {
		Errors = append(Errors, err)
		return Errors
	}

	if user.(User).Password != userData.Password {
		Errors = append(Errors, errors.New("password doesn't match"))
		return Errors
	}

	err_array := checkUserExistance(userData)
	if err != nil {
		Errors = append(Errors, err_array...)
		return Errors
	}

	dbCon := db.GetConnection()

	_, err = dbCon.Exec("update Users set email = ?, username = ?, password = ? where userID = ?", userData.Email, userData.Username, userData.Password, user.(User).UserId)

	if err != nil {
		Errors = append(Errors, err)
		return Errors
	}

	return nil

}
