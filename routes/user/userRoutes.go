package userRoutes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"
	"github.com/RubyLegend/dictionary-backend/middleware/httphelper"
	userHelper "github.com/RubyLegend/dictionary-backend/middleware/users"
	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"
	"github.com/golang-jwt/jwt/v5"

	"github.com/julienschmidt/httprouter"
)

type userHelperWrapper interface {
	GenerateJWT(string) (string, error)
	VerifyCredentials(userRepo.User) (userRepo.User, []error)
	VerifyAuthorizationToken(string) bool
	VerifyJWT(http.ResponseWriter, *http.Request, map[string]any) jwt.MapClaims
	LogoutJWT(http.ResponseWriter, *http.Request, map[string]any)
}

type userHelpWrap struct{}

func (u userHelpWrap) GenerateJWT(username string) (string, error) {
	return userHelper.GenerateJWT(username)
}

func (u userHelpWrap) VerifyAuthorizationToken(tokenString string) bool {
	return userHelper.VerifyAuthorizationToken(tokenString)
}

func (u userHelpWrap) VerifyJWT(w http.ResponseWriter, r *http.Request, resp map[string]any) jwt.MapClaims {
	return userHelper.VerifyJWT(w, r, resp)
}

func (u userHelpWrap) LogoutJWT(w http.ResponseWriter, r *http.Request, resp map[string]any) {
	userHelper.LogoutJWT(w, r, resp)
}

func (u userHelpWrap) VerifyCredentials(userData userRepo.User) (userRepo.User, []error) {
	return userHelper.VerifyCredentials(userData)
}

type userRepoWrapper interface {
	GetUser(userRepo.User) (userRepo.User, error)
	AddUser(userRepo.User) []error
	DeleteUser(userRepo.User) error
	EditUser(string, userRepo.User) []error
}

type userRepoWrap struct{}

func (u userRepoWrap) GetUser(userData userRepo.User) (userRepo.User, error) {
	return userRepo.GetUser(userData)
}

func (u userRepoWrap) AddUser(userData userRepo.User) []error {
	return userRepo.AddUser(userData)
}

func (u userRepoWrap) DeleteUser(userData userRepo.User) error {
	return userRepo.DeleteUser(userData)
}

func (u userRepoWrap) EditUser(currentUsername string, userData userRepo.User) []error {
	return userRepo.EditUser(currentUsername, userData)
}

var (
	userHelp  userHelperWrapper
	userRepoW userRepoWrapper
)

func init() {
	userHelp = userHelpWrap{}
	userRepoW = userRepoWrap{}
}

func UserLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)

	var userData userRepo.UserLogin
	_ = json.NewDecoder(r.Body).Decode(&userData)
	resp := make(map[string]any)

	if userData.Email == "" {
		resp["error"] = []string{"email not provided. cannot authorize"}
		w.WriteHeader(http.StatusNotFound)
	} else {
		userData.Email = strings.ToLower(userData.Email)
		user, err := userHelp.VerifyCredentials(userData.ConvertToUser())
		if err != nil {
			var errors []string
			for _, v := range err {
				errors = append(errors, v.Error())
			}
			resp["error"] = errors
			w.WriteHeader(http.StatusForbidden)
		} else {
			token, err := userHelp.GenerateJWT(user.Username)

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				resp["access_token"] = token
				httphelper.UnpackToResp(user, resp)
				w.WriteHeader(http.StatusOK)
			}
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func UserSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	var userData userRepo.UserCreate
	_ = json.NewDecoder(r.Body).Decode(&userData)
	resp := make(map[string]any)

	if userData.Password != userData.ConfirmPassword {
		resp["error"] = []string{"passwords doesn't match"}
		w.WriteHeader(http.StatusNotAcceptable)
	} else {
		err := userRepoW.AddUser(userData.ConvertToUser())

		if err != nil {
			var errors []string
			for _, v := range err {
				errors = append(errors, v.Error())
			}

			resp["error"] = errors
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			resp["status"] = "User added successfully"
			token, err := userHelp.GenerateJWT(userData.Username)

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				resp["access_token"] = token
				httphelper.UnpackToResp(userData, resp)
				w.WriteHeader(http.StatusOK)
			}

		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func UserLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	resp := make(map[string]any)

	userHelp.LogoutJWT(w, r, resp)

	_ = json.NewEncoder(w).Encode(resp)
}

func UserStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	var userData userRepo.User
	resp := make(map[string]any)

	claims := userHelp.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		userData.Username = claims["username"].(string)
		userData, err := userRepoW.GetUser(userData)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotFound)
		} else {
			userData.Password = ""
			httphelper.UnpackToResp(userData, resp)
		}
	}
	_ = json.NewEncoder(w).Encode(resp)

}

func UserRestoreUsername(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cors.Setup(w, r)
	fmt.Fprintf(w, "Not Implemented\n")
}

func UserRestorePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cors.Setup(w, r)
	fmt.Fprintf(w, "Not Implemented\n")
}

func UserDelete(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	var userData userRepo.User
	resp := make(map[string]any)

	claims := userHelp.VerifyJWT(w, r, resp)

	if resp["error"] == nil {
		userData.Username = claims["username"].(string)
		err := userRepoW.DeleteUser(userData)

		if err == nil {
			resp["status"] = "Success"
		} else {
			resp["status"] = "Failed"
			resp["error"] = []string{err.Error()}
		}
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func UserPatch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	var userData userRepo.User
	_ = json.NewDecoder(r.Body).Decode(&userData)
	resp := make(map[string]any)

	claims := userHelp.VerifyJWT(w, r, resp)

	if resp["error"] == nil {
		username := claims["username"].(string)
		err := userRepoW.EditUser(username, userData)

		if err == nil {
			resp["status"] = "Success"
			userHelp.LogoutJWT(w, r, resp)
			token, err := userHelp.GenerateJWT(userData.Username)

			if err != nil {
				resp["status"] = "Failed"
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusBadGateway)
			} else {
				resp["access_token"] = token
				w.WriteHeader(http.StatusOK)
			}

		} else {
			resp["status"] = "Failed"
			var errors []string
			for _, v := range err {
				errors = append(errors, v.Error())
			}
			resp["error"] = errors
			w.WriteHeader(http.StatusNotAcceptable)
		}

	}

	_ = json.NewEncoder(w).Encode(resp)
}
