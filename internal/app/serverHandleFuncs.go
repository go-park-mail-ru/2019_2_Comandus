package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (s *server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", s.clientUrl)
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	defer func() {
		// TODO: handle err
		r.Body.Close()
	}()
	decoder := json.NewDecoder(r.Body)
	newUserInput := new(model.UserInput)
	err := decoder.Decode(newUserInput)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	err = newUserInput.CheckEmail()
	if err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	s.userType = newUserInput.UserType
	if s.userType != userFreelancer && s.userType != userCustomer {
		s.userType = userFreelancer
	}
	fmt.Println(newUserInput)

	s.usersdb.Mu.Lock()

	user := model.User{
		FirstName: newUserInput.Name,
		SecondName: newUserInput.Surname,
		Email:     newUserInput.Email,
	}

	err = user.BeforeCreate()
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, err)
	}

	err = s.store.User().Create(&user)
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, err)
	}

	fmt.Println("user id: ", user.ID)

	f := model.Freelancer{
		AccountId:        user.ID,
		RegistrationDate: time.Now(),
	}
	err = s.store.Freelancer().Create(&f)
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, err)
	}

	m := model.HireManager{
		AccountID:        user.ID,
		RegistrationDate: time.Now(),
	}
	err = s.store.Manager().Create(&m)
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, err)
	}

	s.usersdb.Mu.Unlock()

	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	session.Values["user_id"] = user.ID
	session.Values["user_type"] = s.userType

	if err := s.sessionStore.Save(r, w, session); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	// TODO: rename cookie2
	cookie := http.Cookie{Name: userTypeCookieName, Value: s.userType}
	cookie2 := http.Cookie{Name: hireManagerIdCookieName, Value: strconv.Itoa(1)} // m.Id
	http.SetCookie(w, &cookie)
	http.SetCookie(w, &cookie2)

	s.respond(w, r, http.StatusCreated, newUserInput)
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.clientUrl)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		s.usersdb.Mu.Lock()

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w,r,http.StatusNotFound, err)
		}

		s.usersdb.Mu.Unlock()

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, &u)))
	})
}

func (s *server) HandleSessionCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", s.clientUrl)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	decoder := json.NewDecoder(r.Body)
	newUserInput := new(model.UserInput)
	err := decoder.Decode(newUserInput)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	log.Println(newUserInput)

	s.usersdb.Mu.Lock()
	for _, v := range s.usersdb.Users {
	//for i := 0; i < len(s.usersdb.Users); i++ {
		if v.Email == newUserInput.Email &&
			v.ComparePassword(newUserInput.Password) {

			u := v
			session, err := s.sessionStore.Get(r, sessionName)
			if err != nil {
				s.usersdb.Mu.Unlock()
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			session.Values["user_id"] = u.ID
			session.Values["user_type"] = s.userType
			if err := s.sessionStore.Save(r, w, session); err != nil {
				s.usersdb.Mu.Unlock()
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.usersdb.Mu.Unlock()
			s.respond(w, r, http.StatusOK, struct {
			}{})
			return
		}
	}
	s.usersdb.Mu.Unlock()
	s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
}

func (s *server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", s.clientUrl)
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	if session == nil {
		s.error(w, r, http.StatusNotFound, errors.New("failed to delete session"))
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		s.error(w, r, http.StatusExpectationFailed, err)
	}
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleSetUserType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//s.respond(w,r, http.StatusOK, nil)
	//return
	// TODO check if input user type invalid
	type Input struct {
		UserType string `json:"type"`
	}
	decoder := json.NewDecoder(r.Body)
	newInput := new(Input)
	err := decoder.Decode(newInput)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	if newInput.UserType != userCustomer && newInput.UserType != userFreelancer {
		s.error(w,r, http.StatusBadRequest, errors.New("user type may be only customer or freelancer"))
	}
	session.Values["user_type"] = newInput.UserType
	session.Save(r, w)
}

func (s *server) HandleShowProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user, sendErr, codeStatus := s.GetUserFromRequest(r)
	if sendErr != nil {
		s.error(w, r, codeStatus, sendErr)
		return
	}
	s.respond(w, r, http.StatusOK, user)
}

func (s *server) HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", s.clientUrl)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	user, sendErr, codeStatus := s.GetUserFromRequest(r)
	if sendErr != nil {
		s.error(w, r, codeStatus, sendErr)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(user)
	log.Println(user)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		SendErr := fmt.Errorf("invalid format of data")
		s.error(w, r, http.StatusBadRequest, SendErr)
		return
	}
	s.respond(w, r, http.StatusOK, struct{}{})

}

func (s *server) HandleEditPassword(w http.ResponseWriter, r *http.Request) {
	type BodyPassword struct {
		Password string
		NewPassword string
		NewPasswordConfirmation string
	}
	w.Header().Set("Content-Type", "application/json")
	var err error
	user, sendErr, codeStatus := s.GetUserFromRequest(r)
	if sendErr != nil {
		s.error(w, r, codeStatus, sendErr)
		return
	}
	bodyPassword := new(BodyPassword)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(bodyPassword)
	if bodyPassword.NewPassword != bodyPassword.NewPasswordConfirmation {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("new Passwords are different"))
		return
	}

	if !user.ComparePassword(bodyPassword.Password) {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("wrong password"))
		return
	}

	newEncryptPassword, err := model.EncryptString(bodyPassword.NewPasswordConfirmation)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, fmt.Errorf("error in updating password"))
		return
	}
	user.EncryptPassword = newEncryptPassword
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) GetUserFromRequest(r *http.Request) (*model.User, error, int) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		sendErr := fmt.Errorf("user isn't authorized")
		return nil, sendErr, http.StatusUnauthorized
	}
	uidInteface := session.Values["user_id"]
	uid := uidInteface.(int)

	s.usersdb.Mu.Lock()
	user, ok := s.usersdb.Users[uid]
	s.usersdb.Mu.Unlock()

	if !ok {
		sendErr := fmt.Errorf("can't find user with id:" + strconv.Itoa(int(uid)))
		return nil, sendErr, http.StatusBadRequest
	}
	return &user, nil, http.StatusOK
}

// TODO:
func (s *server) HandleEditNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user, sendErr, codeStatus := s.GetUserFromRequest(r)
	if sendErr != nil {
		s.error(w, r, codeStatus, sendErr)
		return
	}
	userNotification := s.usersdb.Notifications[user.ID]
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(userNotification)
	log.Println(user)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		sendErr := fmt.Errorf("invalid format of data")
		s.error(w, r, http.StatusBadRequest, sendErr)
		return
	}
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, errors.New("error retrieving the file"))
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errors.New("error retrieving the file"))
		return
	}
	defer file.Close()

	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	uidInterface := session.Values["user_id"]
	uid, ok := uidInterface.(int)
	if !ok {
		s.error(w,r, http.StatusInternalServerError, errors.New("cookie value not set"))
	}

	s.usersdb.Mu.Lock()

	u, ok := s.usersdb.Users[uid]
	if !ok {
		s.error(w,r,http.StatusNotFound, errors.New("no such "))
	}

	//u.Avatar = true
	delete(s.usersdb.Users, uid)
	s.usersdb.Users[uid] = u

	image := bytes.NewBuffer(nil)
	io.Copy(image, file)
	s.usersdb.ImageStore[uid] = image.Bytes()
	u.Avatar = image.Bytes()

	s.usersdb.Mu.Unlock()

	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleDownloadAvatar(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		s.error(w,r,http.StatusUnauthorized, err)
		return
	}
	uidInterface := session.Values["user_id"]
	uid, ok := uidInterface.(int)
	if !ok {
		s.error(w,r, http.StatusInternalServerError, errors.New("cookie value not set"))
	}

	s.usersdb.Mu.Lock()
	user := s.usersdb.Users[uid]
	s.usersdb.Mu.Unlock()

	var openfile *os.File
	if user.Avatar != nil {
		s.usersdb.Mu.Lock()
		image := s.usersdb.ImageStore[uid]
		w.Header().Set("Content-Type", "multipart/form-data")
		w.Header().Set("Content-Length", strconv.Itoa(len(image)))
		if _, err := w.Write(image); err != nil {
			s.error(w,r,http.StatusInternalServerError, err)
		}
		s.usersdb.Mu.Unlock()
	} else {
		filename := "internal/store/avatars/default.png"
		openfile, err = os.Open(filename)
		defer openfile.Close()
		if err != nil {
			s.error(w, r, http.StatusNotFound, errors.New("cant open file"))
			return
		}
		fileHeader := make([]byte, 100000) // max image size!!!
		openfile.Read(fileHeader)
		fileContentType := http.DetectContentType(fileHeader)

		fileStat, _ := openfile.Stat()
		fileSize := strconv.FormatInt(fileStat.Size(), 10)

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", fileContentType)
		w.Header().Set("Content-Length", fileSize)

		openfile.Seek(0, 0)
		io.Copy(w, openfile)
	}
	s.respond(w,r,http.StatusOK, struct{}{})
}

func (s *server) HandleRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", s.clientUrl)
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	user, sendErr, codeStatus := s.GetUserFromRequest(r)
	if sendErr != nil {
		s.error(w, r, codeStatus, sendErr)
		return
	}
	hireManager := s.usersdb.HireManagers[user.ID]
	company := s.usersdb.Companies[hireManager.ID]
	var roles []*model.Role
	clientRole := &model.Role{
		Role:   "client",
		Label:  company.CompanyName,
		Avatar: "/default.png",
	}
	freelanceRole := &model.Role{
		Role:   "freelancer",
		Label:  user.FirstName + " " + user.SecondName,
		Avatar: "/default.png",
	}
	roles = append(roles, clientRole)
	roles = append(roles, freelanceRole)
	s.respond(w, r, http.StatusOK, roles)
}

func (s *server) HandleGetAuthHistory(w http.ResponseWriter, r *http.Request) {
	// TODO: get auth history
}

func (s *server) HandleGetSecQuestion(w http.ResponseWriter, r *http.Request) {
	// TODO: get sec question
}

func (s *server) HandleEditSecQuestion(w http.ResponseWriter, r *http.Request) {
	// TODO: edit sec question
}

func (s *server) HandleCheckSecQuestion(w http.ResponseWriter, r *http.Request) {
	// TODO: check seq question
}

func (s *server) HandleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", s.clientUrl)
	w.Header().Set("Access-Control-Allow-Methods", "POST,PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-Lol")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	s.respond(w, r, http.StatusOK, nil)
}

func (s *server) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	if session == nil {
		s.error(w, r, http.StatusNotFound, errors.New("session is nil"))
		return
	}

	uti := session.Values["user_type"]
	ut, ok := uti.(string)
	if !ok {
		s.error(w,r, http.StatusInternalServerError, errors.New("cookie value not set"))
	}

	log.Println(ut)

	// TODO: add test for this case
	//if ut != userCustomer {
	//	s.error(w, r, http.StatusBadRequest, errors.New("current user is not a manager"))
	//	return
	//}

	defer func() {
		// TODO: handle err
		r.Body.Close()
	}()
	decoder := json.NewDecoder(r.Body)
	newJob := new(model.Job)
	err = decoder.Decode(newJob)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	newJob.ID = len(s.usersdb.Jobs)

	uidi := session.Values["user_id"]
	uid, ok := uidi.(int)
	if !ok {
		s.error(w,r, http.StatusInternalServerError, errors.New("cookie value not set"))
	}

	// TODO write getbyID func for Hire Manager
	for i := 0; i < len(s.usersdb.HireManagers); i++ {
		if s.usersdb.HireManagers[i].AccountID == uid {
			newJob.HireManagerId = s.usersdb.HireManagers[i].ID
			//s.usersdb.Jobs = append(s.usersdb.Jobs, *newJob)
			id := len(s.usersdb.Jobs) + 1
			s.usersdb.Jobs[id] = *newJob
			s.respond(w, r, http.StatusOK, newJob)
			return
		}
	}
	s.error(w, r, http.StatusInternalServerError, errors.New("client not found"))
}

func (s *server) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, errors.New("wrong id"))
	}

	job, ok := s.usersdb.Jobs[id]
	if !ok {
		s.error(w, r, http.StatusNotFound, errors.New("job not found"))
	}

	s.respond(w, r, http.StatusOK, &job)
}


func (s *server) HandleEditFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", s.clientUrl)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	user, sendErr, codeStatus := s.GetUserFromRequest(r)
	if sendErr != nil {
		s.error(w, r, codeStatus, sendErr)
		return
	}
	profile, ok := s.usersdb.Freelancers[user.ID]
	if !ok {
		s.error(w,r,http.StatusNotFound, errors.New("no such freelancer"))
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(profile)
	fmt.Println(user)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		sendErr := fmt.Errorf("invalid format of data")
		s.error(w, r, http.StatusBadRequest, sendErr)
		return
	}
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleGetFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", s.clientUrl)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	vars := mux.Vars(r)
	freelancerIDStr := vars["freelancerId"]
	freelancerID, err := strconv.Atoi(freelancerIDStr)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("id isn't number"))
	}
	profile, ok := s.usersdb.Freelancers[freelancerID]
	if !ok {
		s.error(w, r, http.StatusBadRequest, err)
	}
	s.respond(w, r, http.StatusOK, &profile)
}
func (s *server) HandleGetAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, errors.New("wrong id"))
	}

	s.usersdb.Mu.Lock()
	user, ok := s.usersdb.Users[id]
	if !ok {
		s.error(w,r,http.StatusNotFound, errors.New("no such user in database"))
	}

	var openfile *os.File
	if user.Avatar != nil {
		s.usersdb.Mu.Lock()
		image := s.usersdb.ImageStore[id]
		w.Header().Set("Content-Type", "multipart/form-data")
		w.Header().Set("Content-Length", strconv.Itoa(len(image)))
		if _, err := w.Write(image); err != nil {
			s.error(w,r,http.StatusInternalServerError, err)
		}
		s.usersdb.Mu.Unlock()
	} else {
		filename := "internal/store/avatars/default.png"
		openfile, err = os.Open(filename)
		if err != nil {
			s.error(w, r, http.StatusNotFound, errors.New("cant open file"))
			return
		}
		defer openfile.Close()

		fileHeader := make([]byte, 100000)
		openfile.Read(fileHeader)

		fileContentType := http.DetectContentType(fileHeader)
		fileStat, _ := openfile.Stat()
		fileSize := strconv.FormatInt(fileStat.Size(), 10)

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", fileContentType)
		w.Header().Set("Content-Length", fileSize)

		openfile.Seek(0, 0)
		io.Copy(w, openfile)
	}
	s.usersdb.Mu.Unlock()
	s.respond(w,r,http.StatusOK, struct{}{})
}
