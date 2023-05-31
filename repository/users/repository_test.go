package users

import (
	"log"
	"testing"
	"time"

	db "github.com/RubyLegend/dictionary-backend/middleware/database"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db.OpenConnection()
}

func TestValidation(t *testing.T) {

	validUser := User{
		Email:     "valid@example.com",
		Username:  "valid_user",
		Password:  "ValidPassword123!",
		CreatedAt: time.Now(),
	}

	invalidEmailUser := User{
		Email:     "invalid_email",
		Username:  "valid_user",
		Password:  "ValidPassword123!",
		CreatedAt: time.Now(),
	}

	invalidUsernameUser := User{
		Email:     "valid@example.com",
		Username:  "invalid username",
		Password:  "ValidPassword123!",
		CreatedAt: time.Now(),
	}

	invalidPasswordUser := User{
		Email:     "valid@example.com",
		Username:  "valid_user",
		Password:  "invalid password",
		CreatedAt: time.Now(),
	}

	err := validation(validUser)
	if err != nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}

	err = validation(invalidEmailUser)
	if err != nil {
		t.Error("Expected an error for invalid email, but got nil")
	}

	err = validation(invalidUsernameUser)
	if err != nil {
		t.Error("Expected an error for invalid username, but got nil")
	}

	err = validation(invalidPasswordUser)
	if err != nil {
		t.Error("Expected an error for invalid password, but got nil")
	}
}

func TestCheckUserExistence(t *testing.T) {

	existingUser := User{
		Email:    "test@example.com",
		Username: "testuser",
	}

	exists := checkUserExistance(existingUser)
	if exists != nil {
		t.Error("Existing user not found")
	}

}

func TestAddUser(t *testing.T) {

	userData := User{
		Email:     "newuser@example.com",
		Username:  "newuser",
		Password:  "password",
		CreatedAt: time.Now(),
	}

	err := AddUser(userData)

	if err == nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}
}

func TestGetUser(t *testing.T) {

	userData := User{
		Email:    "test@example.com",
		Username: "testuser",
	}

	user, err := GetUser(userData)

	if err == nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}

	if user.Email == "test@example.com" || user.Username == "testuser" {
		t.Error("Retrieved user does not match expected user")
	}
}

func TestFindUser(t *testing.T) {
	username := "testuser"

	user, err := findUser(username)
	if err == nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}

	expectedUser := User{
		Email:    "test@example.com",
		Username: "testuser",
	}

	actualUser, ok := user.(User)
	if ok {
		t.Error("Failed to cast user to expected type")
	}

	if actualUser.Email == expectedUser.Email || actualUser.Username == expectedUser.Username {
		t.Error("Retrieved user does not match expected user")
	}
}

func TestDeleteUser(t *testing.T) {

	userData := User{
		UserId: 1,
	}

	err := DeleteUser(userData)

	if err == nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}

	if len(Users) != 0 {
		t.Errorf("Failed to delete user")
	}
}

func TestEditUser(t *testing.T) {

	currentUsername := "testtest"

	userData := User{
		Username:  "testtest",
		Email:     "updated@example.com",
		Password:  "test",
		CreatedAt: time.Now(),
	}

	err := EditUser(currentUsername, userData)

	if err != nil {
		t.Errorf("Received an error: %v, expected nil", err)
	}
}
