package userRoutes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	userHelper "github.com/RubyLegend/dictionary-backend/middleware/users"
	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type userHelpMock struct{}

var (
	JWT              = "someValue"
	GenerateJWTError error
)

func (u userHelpMock) GenerateJWT(username string) (string, error) {
	return JWT, GenerateJWTError
}

func (u userHelpMock) VerifyAuthorizationToken(tokenString string) bool {
	return userHelper.VerifyAuthorizationToken(tokenString)
}

var (
	claims          = make(jwt.MapClaims)
	failToVerifyJWT bool
)

func (u userHelpMock) VerifyJWT(w http.ResponseWriter, r *http.Request, resp map[string]any) jwt.MapClaims {
	if failToVerifyJWT {
		resp["error"] = []string{"Failed to verify JWT"}
		w.WriteHeader(http.StatusUnauthorized)
	}
	return claims
}

func (u userHelpMock) LogoutJWT(w http.ResponseWriter, r *http.Request, resp map[string]any) {
	resp["status"] = "Success"
}

var errors = []error{}

func (u userHelpMock) VerifyCredentials(userData userRepo.User) (userRepo.User, []error) {
	return userRepo.User{}, errors
}

type userRepoMock struct{}

var failToGetUser bool

func (u userRepoMock) GetUser(userData userRepo.User) (userRepo.User, error) {
	if failToGetUser {
		return userRepo.User{}, fmt.Errorf("Failed to get user")
	} else {
		return userRepo.User{
				UserId:   1,
				Username: "test",
				Password: "test",
				Email:    "test@example.com",
			},
			nil
	}
}

var AddUserError []error

func (u userRepoMock) AddUser(userData userRepo.User) []error {
	return AddUserError
}

var DeleteUserError error

func (u userRepoMock) DeleteUser(userData userRepo.User) error {
	return DeleteUserError
}

var EditUserError []error

func (u userRepoMock) EditUser(currentUsername string, userData userRepo.User) []error {
	return EditUserError
}

func init() {
	userHelp = userHelpMock{}
	userRepoW = userRepoMock{}
}

func TestUserLogin(t *testing.T) {
	// Create a request body with the required JSON payload
	requestBody := []byte(`{"Email": "test@example.com"}`)

	// Create a new request with the request body
	req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock results
	errors = nil

	// Call the UserLogin handler function
	UserLogin(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.NotEmpty(t, respBody["access_token"])
	assert.Empty(t, respBody["error"])
}

func TestUserLogin_JWT_problem(t *testing.T) {
	// Create a request body with the required JSON payload
	requestBody := []byte(`{"Email": "test@example.com"}`)

	// Create a new request with the request body
	req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock results
	errors = nil
	JWT = ""
	GenerateJWTError = fmt.Errorf("Failed to generate JWT")

	// Call the UserLogin handler function
	UserLogin(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Empty(t, respBody["access_token"])
	assert.Equal(t, []interface{}([]interface{}{"Failed to generate JWT"}), respBody["error"])
}

func TestUserLogin_InvalidEmail(t *testing.T) {
	// Create a request body with an empty email field
	requestBody := []byte(`{"Email": ""}`)

	// Create a new request with the request body
	req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the UserLogin handler function
	UserLogin(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}([]interface{}{"email not provided. cannot authorize"}), respBody["error"])
}

func TestUserLogin_InvalidCredentials(t *testing.T) {
	// Create a request body with valid email but invalid credentials
	requestBody := []byte(`{"Email": "test@example.com"}`)

	// Create a new request with the request body
	req, err := http.NewRequest("POST", "/api/v1/user/login", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyCredentials function to return an error
	errors = []error{fmt.Errorf("invalid credentials")}

	// Call the UserLogin handler function
	UserLogin(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusForbidden, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}([]interface{}{"invalid credentials"}), respBody["error"])
}

func TestUserSignup(t *testing.T) {
	// Create a request body with the required JSON payload
	requestBody := []byte(`{"email": "test@example.com", "username": "test", "password": "test", "confirmPassword": "test"}`)

	// Create a new request with the request body
	req, err := http.NewRequest("POST", "/api/v1/user/signup", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock values
	errors = nil
	JWT = "someValue"
	GenerateJWTError = nil

	// Call the UserSignup handler function
	UserSignup(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, "User added successfully", respBody["status"])
	assert.NotEmpty(t, respBody["access_token"])
	assert.Empty(t, respBody["error"])
}

func TestUserSignup_JWT_problem(t *testing.T) {
	// Create a request body with the required JSON payload
	requestBody := []byte(`{"Email": "test@example.com"}`)

	// Create a new request with the request body
	req, err := http.NewRequest("POST", "/api/v1/user/signup", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock results
	errors = nil
	JWT = ""
	GenerateJWTError = fmt.Errorf("Failed to generate JWT")

	// Call the UserLogin handler function
	UserSignup(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Empty(t, respBody["access_token"])
	assert.Equal(t, []interface{}([]interface{}{"Failed to generate JWT"}), respBody["error"])
}

func TestUserSignup_PasswordMissmatch(t *testing.T) {
	// Create a request body with the required JSON payload
	requestBody := []byte(`{"email": "test@example.com", "username": "test", "password": "test", "confirmPassword": "test1"}`)

	// Create a new request with the request body
	req, err := http.NewRequest("POST", "/api/v1/user/signup", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the UserSignup handler function
	UserSignup(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusNotAcceptable, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Empty(t, respBody["status"])
	assert.Empty(t, respBody["access_token"])
	assert.Equal(t, []interface{}([]interface{}{"passwords doesn't match"}), respBody["error"])
}

func TestUserSignup_AddUser_Error(t *testing.T) {
	// Create a request body with the required JSON payload
	requestBody := []byte(`{"email": "test@example.com", "username": "test", "password": "test", "confirmPassword": "test"}`)

	// Create a new request with the request body
	req, err := http.NewRequest("POST", "/api/v1/user/signup", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	AddUserError = []error{fmt.Errorf("User already exist")}

	// Call the UserSignup handler function
	UserSignup(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusNotAcceptable, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Empty(t, respBody["status"])
	assert.Empty(t, respBody["access_token"])
	assert.Equal(t, []interface{}([]interface{}{"User already exist"}), respBody["error"])
}

func TestUserLogout(t *testing.T) {
	// Create a new request with the request body
	req, err := http.NewRequest("POST", "/api/v1/user/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Set the request Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock values
	errors = nil
	JWT = "someValue"
	GenerateJWTError = nil

	// Call the UserSignup handler function
	UserLogout(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, "Success", respBody["status"])
}

func TestUserStatus(t *testing.T) {
	// Create a new request with a GET method
	req, err := http.NewRequest("GET", "/api/v1/user/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock variables
	claims["username"] = "test"

	// Call the UserStatus handler function
	UserStatus(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.NotEmpty(t, respBody["username"])
	assert.Empty(t, respBody["error"])
}

func TestUserStatus_Unauthorized(t *testing.T) {
	// Create a new request with a GET method
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock variables
	failToVerifyJWT = true

	// Call the UserStatus handler function
	UserStatus(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"Failed to verify JWT"}, respBody["error"])
}

func TestUserStatus_GetUserError(t *testing.T) {
	// Create a new request with a GET method
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock variables
	claims["username"] = "test"
	failToVerifyJWT = false
	failToGetUser = true

	// Call the UserStatus handler function
	UserStatus(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"Failed to get user"}, respBody["error"])
}

func TestUserDelete(t *testing.T) {
	// Create a new request with a DELETE method
	req, err := http.NewRequest("DELETE", "/delete", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock variables
	failToVerifyJWT = false
	claims["username"] = "test"
	DeleteUserError = nil

	// Call the UserDelete handler function
	UserDelete(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, "Success", respBody["status"])
	assert.Empty(t, respBody["error"])
}

func TestUserDelete_Unauthorized(t *testing.T) {
	// Create a new request with a DELETE method
	req, err := http.NewRequest("DELETE", "/delete", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock variables
	failToVerifyJWT = true
	claims["username"] = nil
	DeleteUserError = nil

	// Call the UserDelete handler function
	UserDelete(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"Failed to verify JWT"}, respBody["error"])
}

func TestUserDelete_DeleteUserError(t *testing.T) {
	// Create a new request with a DELETE method
	req, err := http.NewRequest("DELETE", "/delete", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock variables
	failToVerifyJWT = false
	claims["username"] = "test"
	DeleteUserError = fmt.Errorf("failed to delete user")

	// Call the UserDelete handler function
	UserDelete(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, "Failed", respBody["status"])
	assert.Equal(t, []interface{}{"failed to delete user"}, respBody["error"])
}

func TestUserPatch(t *testing.T) {
	// Create a JSON payload for the request body
	requestBody := []byte(`{"firstName": "John", "lastName": "Doe"}`)

	// Create a new request with a PATCH method and the payload
	req, err := http.NewRequest("PATCH", "/patch", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock variables
	failToVerifyJWT = false
	claims["username"] = "test"

	// Call the UserPatch handler function
	UserPatch(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, "Success", respBody["status"])
	assert.NotEmpty(t, respBody["access_token"])
	assert.Empty(t, respBody["error"])
}

func TestUserPatch_Unauthorized(t *testing.T) {
	// Create a JSON payload for the request body
	requestBody := []byte(`{"firstName": "John", "lastName": "Doe"}`)

	// Create a new request with a PATCH method and the payload
	req, err := http.NewRequest("PATCH", "/patch", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock variables
	failToVerifyJWT = true

	// Call the UserPatch handler function
	UserPatch(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"Failed to verify JWT"}, respBody["error"])
}

func TestUserPatch_EditUserError(t *testing.T) {
	// Create a JSON payload for the request body
	payload := `{"firstName": "John", "lastName": "Doe"}`

	// Create a new request with a PATCH method and the payload
	req, err := http.NewRequest("PATCH", "/patch", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock variables
	failToVerifyJWT = false
	claims["username"] = "test"
	EditUserError = []error{fmt.Errorf("failed to edit user")}

	// Call the UserPatch handler function
	UserPatch(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusNotAcceptable, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, "Failed", respBody["status"])
	assert.Equal(t, []interface{}{"failed to edit user"}, respBody["error"])
}

func TestUserPatch_GenerateJWTError(t *testing.T) {
	// Create a JSON payload for the request body
	payload := `{"firstName": "John", "lastName": "Doe"}`

	// Create a new request with a PATCH method and the payload
	req, err := http.NewRequest("PATCH", "/patch", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Setting up mock variables
	failToVerifyJWT = false
	claims["username"] = "test"
	EditUserError = nil
	GenerateJWTError = fmt.Errorf("Failed to generate new JWT token")

	// Call the UserPatch handler function
	UserPatch(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusBadGateway, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, "Failed", respBody["status"])
	assert.Equal(t, []interface{}{"Failed to generate new JWT token"}, respBody["error"])
}
