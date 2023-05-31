package dictionaryRoutes

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	dictionaryRepo "github.com/RubyLegend/dictionary-backend/repository/dictionary"
	dictionaryToWordsRepo "github.com/RubyLegend/dictionary-backend/repository/dictionaryToWords"
	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"
	wordRepo "github.com/RubyLegend/dictionary-backend/repository/words"
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
		resp["error"] = []string{"failed to verify JWT"}
		w.WriteHeader(http.StatusUnauthorized)
	}
	return claims
}

type userRepoMock struct{}

var failToGetUser bool

func (u userRepoMock) GetUser(userData userRepo.User) (userRepo.User, error) {
	if failToGetUser {
		return userRepo.User{}, fmt.Errorf("failed to get user")
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

type dictionaryRepoMock struct{}

var failGetDictionary bool

func (d dictionaryRepoMock) GetDictionary(UserId int) ([]dictionaryRepo.Dictionary, error) {
	if failGetDictionary {
		return nil, fmt.Errorf("failed to get dictionary")
	}
	return []dictionaryRepo.Dictionary{
		{
			DictionaryId: 1,
			UserId:       1,
			Name:         "test",
			Total:        1,
		},
	}, nil
}

var failAddDictionary bool

func (d dictionaryRepoMock) AddDictionary(UserId int, dictionaryData dictionaryRepo.DictionaryPost) []error {
	if failAddDictionary {
		return []error{fmt.Errorf("failed to add dictionary")}
	}
	return nil
}

var failDeleteDictionary bool

func (d dictionaryRepoMock) DeleteDictionary(UserId int, DictionaryId int) error {
	if failDeleteDictionary {
		return fmt.Errorf("failed to delete dictionary")
	}
	return nil
}

var (
	failUpdateDictionary bool
	UpdatedDictionary    dictionaryRepo.Dictionary
)

func (d dictionaryRepoMock) UpdateDictionary(UserId int, DictionaryId int, dictionaryData dictionaryRepo.Dictionary) (dictionaryRepo.Dictionary, error) {
	if failUpdateDictionary {
		return dictionaryRepo.Dictionary{}, fmt.Errorf("failed to update dictionary")
	}
	return UpdatedDictionary, nil
}

type dtwrMock struct{}

var failToGetWords bool

func (d dtwrMock) GetWords(DictionaryId int, page int, limit int) ([]dictionaryToWordsRepo.DictionaryToWords, int, error) {
	if failToGetWords {
		return []dictionaryToWordsRepo.DictionaryToWords{}, -1, fmt.Errorf("failed to get words")
	}
	return []dictionaryToWordsRepo.DictionaryToWords{
			{
				DictionaryId: 1,
				WordId:       1,
			},
			{
				DictionaryId: 1,
				WordId:       2,
			},
		},
		1,
		nil
}

type wrMock struct{}

var failWordIDtoWords bool

func (w wrMock) WordIDtoWords(dictToWords []dictionaryToWordsRepo.DictionaryToWords) ([]wordRepo.Word, error) {
	if failWordIDtoWords {
		return nil, fmt.Errorf("failed to convert words")
	}
	return []wordRepo.Word{
		{
			WordId:      1,
			Name:        "test",
			IsLearned:   false,
			Translation: "test_transl",
		},
	}, nil
}

func init() {
	userHelp = userHelpMock{}
	userRepoW = userRepoMock{}
	drW = dictionaryRepoMock{}
	dtwr = dtwrMock{}
	wr = wrMock{}
}

func TestDictionaryGet(t *testing.T) {
	// Create a new GET request
	req, err := http.NewRequest("GET", "/api/v1/dictionary", nil)
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

	// Mock the GetDictionary function to return a dictionary without error
	failGetDictionary = false

	// Call the DictionaryGet handler function
	DictionaryGet(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody []interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(respBody[0])
	// Check the expected response
	assert.Equal(t, "test", respBody[0].(map[string]interface{})["name"])
	assert.Equal(t, float64(1), respBody[0].(map[string]interface{})["total"])
}

func TestDictionaryGet_FailJWT(t *testing.T) {
	// Create a new GET request
	req, err := http.NewRequest("GET", "/api/v1/dictionary", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = true

	// Call the DictionaryGet handler function
	DictionaryGet(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to verify JWT"}, respBody["error"])
}

func TestDictionaryGet_UserNotFound(t *testing.T) {
	// Create a new GET request
	req, err := http.NewRequest("GET", "/api/v1/dictionary", nil)
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

	// Call the DictionaryGet handler function
	DictionaryGet(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"failed to get user"}, respBody["error"])
}

func TestDictionaryGet_GetDictionaryError(t *testing.T) {
	// Create a new GET request
	req, err := http.NewRequest("GET", "/api/v1/dictionary", nil)
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

	// Mock the GetDictionary function to return an error
	failGetDictionary = true

	// Call the DictionaryGet handler function
	DictionaryGet(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"failed to get dictionary"}, respBody["error"])
}

func TestDictionaryGetWords(t *testing.T) {
	// Create a new GET request with query parameters
	req, err := http.NewRequest("GET", "/api/v1/dictionary/1?page=2&limit=10", nil)
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

	// Mock the GetWords function to return word IDs and count without error
	failToGetWords = false

	// Mock the WordIDtoWords function to return words without error
	failWordIDtoWords = false

	// Preparing httprouter params
	// Call the DictionaryGetWords handler function
	DictionaryGetWords(rr, req, httprouter.Params{{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Len(t, respBody["words"].([]interface{})[0], 5)
	assert.Equal(t, 1, int(respBody["count"].(float64)))
	assert.Equal(t, 10, int(respBody["limit"].(float64)))
	assert.Equal(t, 2, int(respBody["page"].(float64)))
	assert.Equal(t, math.Ceil(float64(1)/float64(10)), respBody["pages"].(float64))
}

func TestDictionaryGetWords_WithoutURLParams(t *testing.T) {
	// Create a new GET request with query parameters
	req, err := http.NewRequest("GET", "/api/v1/dictionary/1", nil)
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

	// Mock the GetWords function to return word IDs and count without error
	failToGetWords = false

	// Mock the WordIDtoWords function to return words without error
	failWordIDtoWords = false

	// Preparing httprouter params
	// Call the DictionaryGetWords handler function
	DictionaryGetWords(rr, req, httprouter.Params{{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Len(t, respBody["words"].([]interface{})[0], 5)
	assert.Equal(t, 1, int(respBody["count"].(float64)))
	assert.Equal(t, int(^uint(0)>>1)*(-1)-1, int(respBody["limit"].(float64))) // Due to float64 conversions I'm getting lower bound of int, instead of upper
	assert.Equal(t, 1, int(respBody["page"].(float64)))
	assert.Equal(t, math.Ceil(float64(1)/float64(10)), respBody["pages"].(float64))
}

func TestDictionaryGetWords_GetUserError(t *testing.T) {
	// Create a new GET request with query parameters
	req, err := http.NewRequest("GET", "/api/v1/dictionary/1?page=2&limit=10", nil)
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

	// Mock the GetWords function to return word IDs and count without error
	failToGetWords = false

	// Mock the WordIDtoWords function to return words without error
	failWordIDtoWords = false

	// Preparing httprouter params
	// Call the DictionaryGetWords handler function
	DictionaryGetWords(rr, req, httprouter.Params{{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to get user"}, respBody["error"])
}

func TestDictionaryGetWords_IDConversionError(t *testing.T) {
	// Create a new GET request with query parameters
	req, err := http.NewRequest("GET", "/api/v1/dictionary/sadfasf?page=2&limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return an error
	failToGetUser = false

	// Mock the GetWords function to return word IDs and count without error
	failToGetWords = false

	// Mock the WordIDtoWords function to return words without error
	failWordIDtoWords = false

	// Preparing httprouter params
	// Call the DictionaryGetWords handler function
	DictionaryGetWords(rr, req, httprouter.Params{{Key: "id", Value: "sadfasf"}})

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"strconv.Atoi: parsing \"sadfasf\": invalid syntax"}, respBody["error"])
}

func TestDictionaryGetWords_GetWordsError(t *testing.T) {
	// Create a new GET request with query parameters
	req, err := http.NewRequest("GET", "/api/v1/dictionary/1?page=2&limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return an error
	failToGetUser = false

	// Mock the GetWords function to return word IDs and count without error
	failToGetWords = true

	// Mock the WordIDtoWords function to return words without error
	failWordIDtoWords = false

	// Preparing httprouter params
	// Call the DictionaryGetWords handler function
	DictionaryGetWords(rr, req, httprouter.Params{{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to get words"}, respBody["error"])
}

func TestDictionaryGetWords_WordIDtoWordsError(t *testing.T) {
	// Create a new GET request with query parameters
	req, err := http.NewRequest("GET", "/api/v1/dictionary/1?page=2&limit=10", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return an error
	failToGetUser = false

	// Mock the GetWords function to return word IDs and count without error
	failToGetWords = false

	// Mock the WordIDtoWords function to return words without error
	failWordIDtoWords = true

	// Preparing httprouter params
	// Call the DictionaryGetWords handler function
	DictionaryGetWords(rr, req, httprouter.Params{{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to convert words"}, respBody["error"])
}

func TestDictionaryPost(t *testing.T) {
	// Create a sample request body
	requestBody := `{"name": "My Dictionary"}`

	// Create a new POST request with the request body
	req, err := http.NewRequest("POST", "/api/v1/dictionary", strings.NewReader(requestBody))
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

	// Mock the AddDictionary function to return no error
	failAddDictionary = false

	// Call the DictionaryPost handler function
	DictionaryPost(rr, req, nil)

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

func TestDictionaryPost_FailedJWT(t *testing.T) {
	// Create a sample request body
	requestBody := `{"name": "My Dictionary"}`

	// Create a new POST request with the request body
	req, err := http.NewRequest("POST", "/api/v1/dictionary", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = true

	// Mock the GetUser function to return a user without error
	failToGetUser = false

	// Mock the AddDictionary function to return no error
	failAddDictionary = false

	// Call the DictionaryPost handler function
	DictionaryPost(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to verify JWT"}, respBody["error"])
}

func TestDictionaryPost_GetUserError(t *testing.T) {
	// Create a sample request body
	requestBody := `{"name": "My Dictionary"}`

	// Create a new POST request with the request body
	req, err := http.NewRequest("POST", "/api/v1/dictionary", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return a user without error
	failToGetUser = true

	// Mock the AddDictionary function to return no error
	failAddDictionary = false

	// Call the DictionaryPost handler function
	DictionaryPost(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to get user"}, respBody["error"])
}

func TestDictionaryPost_AddDictionaryError(t *testing.T) {
	// Create a sample request body
	requestBody := `{"name": "My Dictionary"}`

	// Create a new POST request with the request body
	req, err := http.NewRequest("POST", "/api/v1/dictionary", strings.NewReader(requestBody))
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

	// Mock the AddDictionary function to return no error
	failAddDictionary = true

	// Call the DictionaryPost handler function
	DictionaryPost(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to add dictionary"}, respBody["error"])
}

func TestDictionaryPatch(t *testing.T) {
	// Create a sample request body
	requestBody := `{"name": "Updated Dictionary"}`

	// Create a new PATCH request with the request body
	req, err := http.NewRequest("PATCH", "/api/v1/dictionary/1", strings.NewReader(requestBody))
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

	// Mock the UpdateDictionary function to return the updated dictionary without error
	UpdatedDictionary = dictionaryRepo.Dictionary{
		DictionaryId: 1,
		Name:         "Updated Dictionary",
		Total:        1,
	}
	failUpdateDictionary = false

	// Call the DictionaryPatch handler function
	DictionaryPatch(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]dictionaryRepo.Dictionary
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, UpdatedDictionary, respBody["dictionary"])
}

func TestDictionaryPatch_FailedJWT(t *testing.T) {
	// Create a sample request body
	requestBody := `{"name": "Updated Dictionary"}`

	// Create a new PATCH request with the request body
	req, err := http.NewRequest("PATCH", "/api/v1/dictionary/1", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = true

	// Mock the GetUser function to return a user without error
	failToGetUser = false

	// Mock the UpdateDictionary function to return the updated dictionary without error
	UpdatedDictionary = dictionaryRepo.Dictionary{
		DictionaryId: 1,
		Name:         "Updated Dictionary",
		Total:        1,
	}
	failUpdateDictionary = false

	// Call the DictionaryPatch handler function
	DictionaryPatch(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to verify JWT"}, respBody["error"])
}

func TestDictionaryPatch_GetUserError(t *testing.T) {
	// Create a sample request body
	requestBody := `{"name": "Updated Dictionary"}`

	// Create a new PATCH request with the request body
	req, err := http.NewRequest("PATCH", "/api/v1/dictionary/1", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return a user without error
	failToGetUser = true

	// Mock the UpdateDictionary function to return the updated dictionary without error
	UpdatedDictionary = dictionaryRepo.Dictionary{
		DictionaryId: 1,
		Name:         "Updated Dictionary",
		Total:        1,
	}
	failUpdateDictionary = false

	// Call the DictionaryPatch handler function
	DictionaryPatch(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to get user"}, respBody["error"])
}

func TestDictionaryPatch_URLParamError(t *testing.T) {
	// Create a sample request body
	requestBody := `{"name": "Updated Dictionary"}`

	// Create a new PATCH request with the request body
	req, err := http.NewRequest("PATCH", "/api/v1/dictionary/abc", strings.NewReader(requestBody))
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

	// Mock the UpdateDictionary function to return the updated dictionary without error
	UpdatedDictionary = dictionaryRepo.Dictionary{
		DictionaryId: 1,
		Name:         "Updated Dictionary",
		Total:        1,
	}
	failUpdateDictionary = false

	// Call the DictionaryPatch handler function
	DictionaryPatch(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: "abc"}})

	// Check the response status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"strconv.Atoi: parsing \"abc\": invalid syntax"}, respBody["error"])
}

func TestDictionaryPatch_UpdateDictionaryError(t *testing.T) {
	// Create a sample request body
	requestBody := `{"name": "Updated Dictionary"}`

	// Create a new PATCH request with the request body
	req, err := http.NewRequest("PATCH", "/api/v1/dictionary/1", strings.NewReader(requestBody))
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

	// Mock the UpdateDictionary function to return the updated dictionary without error
	UpdatedDictionary = dictionaryRepo.Dictionary{
		DictionaryId: 1,
		Name:         "Updated Dictionary",
		Total:        1,
	}
	failUpdateDictionary = true

	// Call the DictionaryPatch handler function
	DictionaryPatch(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to update dictionary"}, respBody["error"])
}

func TestDictionaryDelete(t *testing.T) {
	// Create a new DELETE request
	req, err := http.NewRequest("DELETE", "/dictionary/1", nil)
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

	// Mock the DeleteDictionary function to return no error
	failDeleteDictionary = false

	// Call the DictionaryDelete handler function
	DictionaryDelete(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, "success", respBody["status"])
}

func TestDictionaryDelete_FailedJWT(t *testing.T) {
	// Create a new DELETE request
	req, err := http.NewRequest("DELETE", "/dictionary/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = true

	// Mock the GetUser function to return a user without error
	failToGetUser = false

	// Mock the DeleteDictionary function to return no error
	failDeleteDictionary = false

	// Call the DictionaryDelete handler function
	DictionaryDelete(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to verify JWT"}, respBody["error"])
}

func TestDictionaryDelete_GetUUserError(t *testing.T) {
	// Create a new DELETE request
	req, err := http.NewRequest("DELETE", "/dictionary/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the GetUser function to return a user without error
	failToGetUser = true

	// Mock the DeleteDictionary function to return no error
	failDeleteDictionary = false

	// Call the DictionaryDelete handler function
	DictionaryDelete(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusNotFound, rr.Code)

	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to get user"}, respBody["error"])
}

func TestDictionaryDelete_URLParamError(t *testing.T) {
	// Create a new DELETE request
	req, err := http.NewRequest("DELETE", "/dictionary/abc", nil)
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

	// Mock the DeleteDictionary function to return no error
	failDeleteDictionary = false

	// Call the DictionaryDelete handler function
	DictionaryDelete(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: "abc"}})

	// Check the response status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"strconv.Atoi: parsing \"abc\": invalid syntax"}, respBody["error"])
}

func TestDictionaryDelete_DeleteError(t *testing.T) {
	// Create a new DELETE request
	req, err := http.NewRequest("DELETE", "/dictionary/1", nil)
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

	// Mock the DeleteDictionary function to return no error
	failDeleteDictionary = true

	// Call the DictionaryDelete handler function
	DictionaryDelete(rr, req, httprouter.Params{httprouter.Param{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response
	assert.Equal(t, []interface{}{"failed to delete dictionary"}, respBody["error"])
}
