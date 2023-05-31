package wordRoutes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"

	dictionaryToWordsRepo "github.com/RubyLegend/dictionary-backend/repository/dictionaryToWords"

	wordRepo "github.com/RubyLegend/dictionary-backend/repository/words"
)

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
		},
		1,
		nil
}

var AddConnectionError error

func (d dtwrMock) AddConnection(DictionaryId int, WordId int) error {
	return AddConnectionError
}

type wrMock struct{}

var failAddWord bool

func (w wrMock) AddWord(wordData wordRepo.Word) (int, wordRepo.Word, error) {
	if failAddWord {
		return -1, wordRepo.Word{}, fmt.Errorf("failed to add word")
	}
	return 1, wordRepo.Word{WordId: 1, Name: "test", IsLearned: false, Translation: "test_trans"}, nil
}

var UpdateWordError error

func (w wrMock) UpdateWord(wordData wordRepo.Word) error {
	return UpdateWordError
}

var DeleteWordError error

func (w wrMock) DeleteWord(wordData wordRepo.Word) error {
	return DeleteWordError
}

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

func init() {
	dtwr = dtwrMock{}
	wr = wrMock{}
	userHelperW = userHelpMock{}
}

func TestWordPost(t *testing.T) {
	// Create a JSON payload for the request body
	payload := `{"name": "test", "dictionaryId": 1, "translation": "test_trans"}`

	// Create a new request with a POST method and the payload
	req, err := http.NewRequest("POST", "/api/v1/word", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the WordPost handler function
	WordPost(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, "success", respBody["status"])
	assert.Empty(t, respBody["error"])
	// Additional assertions based on the response body structure and data
}

func TestWordPost_Unauthorized(t *testing.T) {
	// Create a JSON payload for the request body
	payload := `{"word": "apple", "dictionaryId": 1}`

	// Create a new request with a POST method and the payload
	req, err := http.NewRequest("POST", "/words", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return an error
	failToVerifyJWT = true

	// Call the WordPost handler function
	WordPost(rr, req, nil)

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

func TestWordPost_AddWordError(t *testing.T) {
	// Create a JSON payload for the request body
	payload := `{"word": "apple", "dictionaryId": 1}`

	// Create a new request with a POST method and the payload
	req, err := http.NewRequest("POST", "/words", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the AddWord function to return an error
	failToVerifyJWT = false
	claims["username"] = "test"
	failAddWord = true

	// Call the WordPost handler function
	WordPost(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusNotAcceptable, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, []interface{}{"failed to add word"}, respBody["error"])
}

func TestWordPost_AddWordConnectionError(t *testing.T) {
	// Create a JSON payload for the request body
	payload := `{"word": "apple", "dictionaryId": 1}`

	// Create a new request with a POST method and the payload
	req, err := http.NewRequest("POST", "/words", strings.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the AddWord function to return an error
	failToVerifyJWT = false
	claims["username"] = "test"
	failAddWord = false
	AddConnectionError = fmt.Errorf("failed to add dictionary to word connection")

	// Call the WordPost handler function
	WordPost(rr, req, nil)

	// Check the response status code
	assert.Equal(t, http.StatusNotAcceptable, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, []interface{}{"failed to add dictionary to word connection"}, respBody["error"])
}

func TestWordDelete(t *testing.T) {
	// Create a new request with a DELETE method
	req, err := http.NewRequest("DELETE", "/words/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the DeleteWord function to return no error
	DeleteWordError = nil

	// Call the WordDelete handler function
	WordDelete(rr, req, httprouter.Params{{Key: "id", Value: "1"}})

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, "success", respBody["status"])
	assert.Empty(t, respBody["error"])
}

func TestWordDelete_InvalidID(t *testing.T) {
	// Create a new request with a DELETE method
	req, err := http.NewRequest("DELETE", "/words/abc", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the WordDelete handler function
	WordDelete(rr, req, httprouter.Params{{Key: "id", Value: "abc"}})

	// Check the response status code
	assert.Equal(t, http.StatusNotAcceptable, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"strconv.Atoi: parsing \"abc\": invalid syntax"}, respBody["error"])
}

func TestWordDelete_DeleteWordError(t *testing.T) {
	// Create a new request with a DELETE method
	req, err := http.NewRequest("DELETE", "/words/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the DeleteWord function to return an error
	DeleteWordError = fmt.Errorf("failed to delete word")

	// Call the WordDelete handler function
	WordDelete(rr, req, httprouter.Params{{Key: "id", Value: "123"}})

	// Check the response status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"failed to delete word"}, respBody["error"])
}

func TestWordPatch(t *testing.T) {
	// Create a sample wordData JSON payload
	wordDataJSON := `{"wordId": 123, "word": "updated word", "meaning": "updated meaning"}`

	// Create a new request with a PATCH method and the sample payload
	req, err := http.NewRequest("PATCH", "/words/123", strings.NewReader(wordDataJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the UpdateWord function to return no error
	UpdateWordError = nil

	// Call the WordPatch handler function
	WordPatch(rr, req, httprouter.Params{{Key: "id", Value: "123"}})

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected response fields
	assert.Equal(t, "success", respBody["status"])
	assert.Empty(t, respBody["error"])
}

func TestWordPatch_InvalidID(t *testing.T) {
	// Create a sample wordData JSON payload
	wordDataJSON := `{"wordId": 123, "word": "updated word", "meaning": "updated meaning"}`

	// Create a new request with a PATCH method and the sample payload
	req, err := http.NewRequest("PATCH", "/words/abc", strings.NewReader(wordDataJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the WordPatch handler function
	WordPatch(rr, req, httprouter.Params{{Key: "id", Value: "abc"}})

	// Check the response status code
	assert.Equal(t, http.StatusNotAcceptable, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"strconv.Atoi: parsing \"abc\": invalid syntax"}, respBody["error"])
}

func TestWordPatch_UpdateWordError(t *testing.T) {
	// Create a sample wordData JSON payload
	wordDataJSON := `{"wordId": 123, "word": "updated word", "meaning": "updated meaning"}`

	// Create a new request with a PATCH method and the sample payload
	req, err := http.NewRequest("PATCH", "/words/123", strings.NewReader(wordDataJSON))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Mock the VerifyJWT function to return no error
	failToVerifyJWT = false
	claims["username"] = "test"

	// Mock the UpdateWord function to return an error
	UpdateWordError = fmt.Errorf("update word error")

	// Call the WordPatch handler function
	WordPatch(rr, req, httprouter.Params{{Key: "id", Value: "123"}})

	// Check the response status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Decode the response body
	var respBody map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	// Check the expected error message
	assert.Equal(t, []interface{}{"update word error"}, respBody["error"])
}
