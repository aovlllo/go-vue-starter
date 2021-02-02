package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/aovlllo/vue-template/pkg/api/response"
	"github.com/aovlllo/vue-template/pkg/model"

	"github.com/dgrijalva/jwt-go"
)

// loginHandler handles user authentication
func (a *API) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	} else if u.Email == "" {
		response.Errorf(w, r, nil, http.StatusBadRequest, "Email address is missing")
		return
	} else if u.Password == "" {
		response.Errorf(w, r, nil, http.StatusBadRequest, "Password is missing")
		return
	}

	u.Email = strings.ToLower(u.Email)

	user, err := a.db.GetUserByEmail(u.Email)
	if err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	} else if user == nil {
		response.Errorf(w, r, err, http.StatusBadRequest, "Invalid email address or password")
		return
	}

	if !user.MatchPassword(u.Password) {
		response.Errorf(w, r, err, http.StatusBadRequest, "Invalid email address or password")
		return
	}

	user.Password = ""

	user.Token, err = a.createJWT(jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	if err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Write(w, r, user)
	return
}

// signupHandler handles user sign up
func (a *API) userSignupHandler(w http.ResponseWriter, r *http.Request) {
	var u model.User
	err := json.NewDecoder(r.Body).Decode(&u)
	switch {
	case err != nil:
		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	case u.Password == "":
		response.Errorf(w, r, nil, http.StatusBadRequest, "Password is missing")
		return
	}
	if err = u.VerifyFields(); err != nil {
		response.Errorf(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	u.Email = strings.ToLower(u.Email)

	err = u.HashPassword()
	if err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	err = a.db.CreateUser(&u)
	if err != nil {
		if err.Error() == "user already exists" {
			response.Errorf(w, r, err, http.StatusBadRequest, "Email address already exists")
			return
		}

		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	u.Password = ""

	response.Write(w, r, u)
	return
}

func (a *API) userUpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(ContextJWTKey).(jwt.MapClaims)
	if !ok {
		response.Errorf(w, r, nil, http.StatusUnauthorized, "Unauthorized")
		return
	}

	id, err := a.getID(claims)
	if err != nil {
		response.Errorf(w, r, nil, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var u model.User
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if err = u.VerifyFields(); err != nil {
		response.Errorf(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	user, err := a.db.GetUser(id)
	if err != nil || user == nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user.Name = u.Name
	user.SecondName = u.SecondName
	user.Email = u.Email
	user.Birth = u.Birth
	user.City = u.City
	user.Sex = u.Sex
	user.Interests = u.Interests

	err = a.db.UpdateUser(user)
	if err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user.Password = ""

	user.Token, err = a.createJWT(jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	if err != nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.Write(w, r, user)
	return
}

func (a *API) userProfileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(ContextJWTKey).(jwt.MapClaims)
	if !ok {
		response.Errorf(w, r, nil, http.StatusInternalServerError, "Internal Server Error")
	}

	id, err := a.getID(claims)
	if err != nil {
		response.Errorf(w, r, nil, http.StatusUnauthorized, "Unauthorized")
	}

	user, err := a.db.GetUser(id)
	if err != nil || user == nil {
		response.Errorf(w, r, err, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user.Password = ""

	response.Write(w, r, user)
	return
}

