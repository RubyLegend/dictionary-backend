package historyRoutes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	historyRepo "github.com/RubyLegend/dictionary-backend/repository/history"
	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

type userHelpMock struct{}

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

type historyRepoMock struct{}

var failToGetHistory bool

func (h historyRepoMock) GetHistory(UserId int) ([]historyRepo.History, error) {
	if failToGetHistory {
		return nil, fmt.Errorf("failed to get history")
	}
	return []historyRepo.History{
		{
			HistoryId: 1,
			UserId:    1,
			WordId:    1,
			IsCorrect: true,
		},
	}, nil
}

var failToDeleteHistory bool

func (h historyRepoMock) DeleteHistory(UserId int, HistoryId int) error {
	if failToDeleteHistory {
		return fmt.Errorf("Failed to delete history")
	}
	return nil
}

func init() {
	userHelp = userHelpMock{}
	userRepoW = userRepoMock{}
	historyRepoW = historyRepoMock{}
}

func TestHistoryGet(t *testing.T) {
	// Create a new request with a GET method
	req, err := http.NewRequest("GET", "/history", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return a user without error
	failToGetUser = false

	// Mock the GetHistory function to return the history without error
	failToGetHistory = false

	// Call the HistoryGet handler function
	HistoryGet(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var history []historyRepo.History
	err = json.NewDecoder(rr.Body).Decode(&history)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected history data
	assert.Equal(t, []historyRepo.History{
		{
			HistoryId: 1,
			UserId:    1,
			WordId:    1,
			IsCorrect: true,
		},
	}, history)
}

func TestHistoryGet_FailJWT(t *testing.T) {
	// Create a new request with a GET method
	req, err := http.NewRequest("GET", "/history", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = true

	// Call the HistoryGet handler function
	HistoryGet(rr, req, nil)

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

func TestHistoryGet_UserNotFound(t *testing.T) {
	// Create a new request with a GET method
	req, err := http.NewRequest("GET", "/history", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return an error
	failToGetUser = true

	// Call the HistoryGet handler function
	HistoryGet(rr, req, nil)

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

func TestHistoryGet_GetHistoryError(t *testing.T) {
	// Create a new request with a GET method
	req, err := http.NewRequest("GET", "/history", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return a user without error
	failToGetUser = false

	// Mock the GetHistory function to return an error
	failToGetHistory = true

	// Call the HistoryGet handler function
	HistoryGet(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"failed to get history"}, respBody["error"])
}

func TestHistoryDelete(t *testing.T) {
	// Create a new request with a DELETE method and a valid history ID
	historyID := 1
	req, err := http.NewRequest("DELETE", "/history/"+strconv.Itoa(historyID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return a user without error
	failToGetUser = false

	// Mock the DeleteHistory function to return no error
	failToDeleteHistory = false

	// Call the HistoryDelete handler function
	HistoryDelete(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: strconv.Itoa(historyID)}})

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, "success", respBody["status"])
}

func TestHistoryDelete_UserNotFound(t *testing.T) {
	// Create a new request with a DELETE method and a valid history ID
	historyID := 123
	req, err := http.NewRequest("DELETE", "/history/"+strconv.Itoa(historyID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return an error
	failToGetUser = true

	// Call the HistoryDelete handler function
	HistoryDelete(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: strconv.Itoa(historyID)}})

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

func TestHistoryDelete_InvalidHistoryID(t *testing.T) {
	// Create a new request with a DELETE method and an invalid history ID
	invalidHistoryID := "invalid"
	req, err := http.NewRequest("DELETE", "/history/"+invalidHistoryID, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return a user without error
	failToGetUser = false

	// Call the HistoryDelete handler function
	HistoryDelete(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: invalidHistoryID}})

	// Check the response status code
	assert.Equal(t, http.StatusNotAcceptable, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"strconv.Atoi: parsing \"invalid\": invalid syntax"}, respBody["error"])
}

func TestHistoryDelete_DeleteHistoryError(t *testing.T) {
	// Create a new request with a DELETE method and a valid history ID
	historyID := 123
	req, err := http.NewRequest("DELETE", "/history/"+strconv.Itoa(historyID), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return a user without error
	failToGetUser = false

	// Mock the DeleteHistory function to return an error
	failToDeleteHistory = true

	// Call the HistoryDelete handler function
	HistoryDelete(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: strconv.Itoa(historyID)}})

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"Failed to delete history"}, respBody["error"])
}
