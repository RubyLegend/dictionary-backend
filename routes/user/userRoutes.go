package userRoutes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RubyLegend/dictionary-backend/middleware/cors"
	userHelper "github.com/RubyLegend/dictionary-backend/middleware/users"
	userRepo "github.com/RubyLegend/dictionary-backend/repository/users"

	"github.com/julienschmidt/httprouter"
)

func UserLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)

	var userData userRepo.User
	_ = json.NewDecoder(r.Body).Decode(&userData)
	resp := make(map[string]any)

	if userData.Email == "" {
		resp["error"] = []string{"email not provided. cannot authorize"}
		w.WriteHeader(http.StatusNotFound)
	} else {
		user, err := userHelper.VerifyCredentials(userData)
		if err != nil {
			var errors []string
			for _, v := range err {
				errors = append(errors, v.Error())
			}
			resp["error"] = errors
			w.WriteHeader(http.StatusForbidden)
		} else {
			token, err := userHelper.GenerateJWT(user.Username)

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				resp["access_token"] = token
				resp["userData"] = user
				w.WriteHeader(http.StatusOK)
			}
		}
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func UserSignup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	var userData userRepo.User
	_ = json.NewDecoder(r.Body).Decode(&userData)
	resp := make(map[string]any)

	err := userRepo.AddUser(userData)

	if err != nil {
		var errors []string
		for _, v := range err {
			errors = append(errors, v.Error())
		}

		resp["error"] = errors
		w.WriteHeader(http.StatusNotAcceptable)
	} else {
		resp["status"] = "User added successfully"
		token, err := userHelper.GenerateJWT(userData.Username)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			resp["access_token"] = token
			resp["userData"] = userData
			w.WriteHeader(http.StatusOK)
		}

	}

	_ = json.NewEncoder(w).Encode(resp)
}

func UserLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	resp := make(map[string]any)

	userHelper.LogoutJWT(w, r, resp)

	_ = json.NewEncoder(w).Encode(resp)
}

func UserStatus(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	cors.Setup(w, r)
	var userData userRepo.User
	resp := make(map[string]any)

	claims := userHelper.VerifyJWT(w, r, resp)
	if resp["error"] == nil {
		userData.Username = claims["username"].(string)
		userData, err := userRepo.GetUser(userData)

		if err != nil {
			resp["error"] = []string{err.Error()}
			w.WriteHeader(http.StatusNotFound)
		} else {
			resp["claims"] = claims
			userData.Password = ""
			resp["userData"] = userData
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

	claims := userHelper.VerifyJWT(w, r, resp)

	if resp["error"] == nil {
		userData.Username = claims["username"].(string)
		err := userRepo.DeleteUser(userData)

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

	claims := userHelper.VerifyJWT(w, r, resp)

	if resp["error"] == nil {
		username := claims["username"].(string)
		err := userRepo.EditUser(username, userData)

		if err == nil {
			resp["status"] = "Success"
			userHelper.LogoutJWT(w, r, resp)
			token, err := userHelper.GenerateJWT(userData.Username)

			if err != nil {
				resp["error"] = []string{err.Error()}
				w.WriteHeader(http.StatusBadGateway)
			} else {
				resp["token"] = token
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
