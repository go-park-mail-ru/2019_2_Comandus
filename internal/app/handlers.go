package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (s *server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateUser:")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	user := new(model.User)
	err := decoder.Decode(user)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	fmt.Println(user)

	if err := user.Validate(); err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-Validate:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := user.BeforeCreate(); err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-BeforeCreate:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.store.Mu.Lock()
	_, err = s.store.User().Create(user)
	s.store.Mu.Unlock()

	log.Println("USER ID: ", user.ID)

	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-CreateUser:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	f := model.Freelancer{
		AccountId:        user.ID,
		RegistrationDate: time.Now(),
	}

	s.store.Mu.Lock()
	err = s.store.Freelancer().Create(&f)
	s.store.Mu.Unlock()

	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-CreateFreelancer:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	m := model.HireManager{
		AccountID:        user.ID,
		RegistrationDate: time.Now(),
	}
	//Здесь нужны Мьютексы ?
	s.store.Mu.Lock()
	err = s.store.Manager().Create(&m)
	s.store.Mu.Unlock()

	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-CreateManager:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-sessionGet:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	session.Values["user_id"] = user.ID

	if err := s.sessionStore.Save(r, w, session); err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-sessionSave:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	// TODO: why we need manager cookie?
	cookie := http.Cookie{Name: hireManagerIdCookieName, Value: strconv.Itoa(1)} // m.Id
	http.SetCookie(w, &cookie)

	s.respond(w, r, http.StatusCreated, user)
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		s.store.Mu.Lock()
		u, err := s.store.User().Find(id.(int))
		s.store.Mu.Unlock()

		if err != nil {
			s.error(w, r, http.StatusNotFound, err)
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, &u)))
	})
}

func (s *server) HandleSessionCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleSessionCreate<-rBodyClose:")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()
	decoder := json.NewDecoder(r.Body)
	user := new(model.User)
	err := decoder.Decode(user)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-DecodeUser:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	log.Println("input user:", user)

	s.store.Mu.Lock()
	u, err := s.store.User().FindByEmail(user.Email)
	s.store.Mu.Unlock()

	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-FindByEmail:")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-sessionGet:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	session.Values["user_id"] = u.ID
	if err := s.sessionStore.Save(r, w, session); err != nil {
		err = errors.Wrapf(err, "HandleSessionCreate<-sessionSave:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleLogout<-sessionGet:")
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		err = errors.Wrapf(err, "HandleLogout<-sessionSave:")
		s.error(w, r, http.StatusExpectationFailed, err)
	}
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleSetUserType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type Input struct {
		UserType string `json:"type"`
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleSetUserType<-rBodyClose:")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()
	decoder := json.NewDecoder(r.Body)
	newInput := new(Input)
	err := decoder.Decode(newInput)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-Decode:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-GetUserFromRequest:")
		s.error(w, r, codeStatus, err)
		return
	}

	err = user.SetUserType(newInput.UserType)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-SetUserType:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	s.store.Mu.Lock()
	err = s.store.User().Edit(user)
	s.store.Mu.Unlock()

	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-sessionEdit:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-sessionGet:")
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	err = session.Save(r, w)
	if err != nil {
		err = errors.Wrapf(err, "HandleSetUserType<-sessionSave:")
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	s.respond(w, r, http.StatusOK, user.UserType)
}

func (s *server) HandleShowProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleShowProfile<-GetUserFromRequest:")
		s.error(w, r, codeStatus, err)
		return
	}
	s.respond(w, r, http.StatusOK, user)
}

func (s *server) HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditProfile<-GetUserFromRequest:")
		s.error(w, r, codeStatus, err)
		return
	}

	// TODO: validate edited user

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditProfile<-rBodyClose:")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(user)

	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		err = errors.Wrapf(err, "HandleEditProfile<-Decode:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	s.store.Mu.Lock()
	err = s.store.User().Edit(user)
	s.store.Mu.Unlock()

	if err != nil {
		err = errors.Wrapf(err, "HandleEditProfile<-userEdit")
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	log.Println("edited profile: ", user)
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleEditPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditPassword<-GetUserFromRequest:")
		s.error(w, r, codeStatus, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditPassword<-rBodyClose:")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()

	bodyPassword := new(model.BodyPassword)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(bodyPassword)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditPassword<-Decode: ")
		s.error(w, r, http.StatusInternalServerError, err)
	}

	if bodyPassword.NewPassword != bodyPassword.NewPasswordConfirmation {
		err = fmt.Errorf("new Passwords are different")
		err = errors.Wrapf(err, "HandleEditPassword<-PasswordAreDifferent:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if !user.ComparePassword(bodyPassword.Password) {
		err = fmt.Errorf("new Passwords are different")
		err = errors.Wrapf(err, "HandleEditPassword<-PasswordAreDifferent:")
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("wrong password"))
		return
	}

	newEncryptPassword, err := model.EncryptString(bodyPassword.NewPasswordConfirmation)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditPassword<-EncryptString:")
		s.error(w, r, http.StatusInternalServerError, fmt.Errorf("error in updating password"))
		return
	}
	user.EncryptPassword = newEncryptPassword

	s.store.Mu.Lock()
	err = s.store.User().Edit(user)
	s.store.Mu.Unlock()

	if err != nil {
		err = errors.Wrapf(err, "HandleEditPassword<-userEdit:")
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) GetUserFromRequest(r *http.Request) (*model.User, error, int) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		sendErr := fmt.Errorf("user isn't authorized")
		return nil, sendErr, http.StatusUnauthorized
	}

	uidInterface := session.Values["user_id"]
	uid := uidInterface.(int)

	s.store.Mu.Lock()
	user, err := s.store.User().Find(uid)
	s.store.Mu.Unlock()

	if err != nil {
		sendErr := fmt.Errorf("can't find user with id:" + strconv.Itoa(int(uid)))
		return nil, sendErr, http.StatusBadRequest
	}
	return user, nil, http.StatusOK
}

// TODO: fix after creating Notifications table
func (s *server) HandleEditNotifications(w http.ResponseWriter, r *http.Request) {
	/*w.Header().Set("Content-Type", "application/json")

	user, sendErr, codeStatus := s.GetUserFromRequest(r)
	if sendErr != nil {
		sendErr = errors.Wrapf(sendErr, "HandleEditNotifications<-GetUserFromRequest:")
		s.error(w, r, codeStatus, sendErr)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditNotifications<-rBodyClose:")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()

	userNotification := s.usersdb.Notifications[user.ID]
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(userNotification)
	log.Println(user)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		err = errors.Wrapf(err, "HandleEditNotifications<-Decode:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}
	s.respond(w, r, http.StatusOK, struct{}{})*/
}

func (s *server) HandleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-ParseMultipartForm:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-FormFile:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-GetUserFromRequest:")
		s.error(w, r, codeStatus, err)
		return
	}

	image := bytes.NewBuffer(nil)
	_, err = io.Copy(image, file)
	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-ioCopy:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	user.Avatar = image.Bytes()

	//s.store.Mu.Lock()
	err = s.store.User().Edit(user)
	//s.usersdb.Mu.Unlock()

	if err != nil {
		err = errors.Wrapf(err, "HandleUploadAvatar<-userEdit:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleDownloadAvatar(w http.ResponseWriter, r *http.Request) {

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleDownloadAvatar<-GetUserFromRequest:")
		s.error(w, r, codeStatus, err)
		return
	}

	if user.Avatar != nil {
		image := user.Avatar
		if _, err := w.Write(image); err != nil {
			err = errors.Wrapf(err, "HandleDownloadAvatar<-Write(image):")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(image)))
		w.Header().Set("Content-Type", "multipart/form-data")

	} else {
		filename := "internal/store/avatars/default.png"

		var openFile *os.File
		openFile, err = os.Open(filename)
		if err != nil {
			err = errors.Wrapf(err, "HandleDownloadAvatar<-Open:")
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		defer func() {
			if err := openFile.Close(); err != nil {
				err = errors.Wrapf(err, "HandleDownloadAvatar<-Close:")
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}()

		fileHeader := make([]byte, 100000) // max image size!!!
		if _, err := openFile.Read(fileHeader); err != nil {
			err = errors.Wrapf(err, "HandleDownloadAvatar<-Read:")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		fileContentType := http.DetectContentType(fileHeader)
		fileStat, _ := openFile.Stat()
		fileSize := strconv.FormatInt(fileStat.Size(), 10)

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", fileContentType)
		w.Header().Set("Content-Length", fileSize)

		if _, err := openFile.Seek(0, 0); err != nil {
			err = errors.Wrapf(err, "HandleDownloadAvatar<-Seek:")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if _, err := io.Copy(w, openFile); err != nil {
			err = errors.Wrapf(err, "HandleDownloadAvatar<-Copy:")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleRoles(w http.ResponseWriter, r *http.Request) {
	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleRoles<-GetUserFromRequest:")
		s.error(w, r, codeStatus, err)
		return
	}

	hireManager, err := s.store.Manager().FindByUser(user.ID)
	if err != nil {
		err = errors.Wrapf(err, "HandleRoles<-ManagerFind:")
		s.error(w, r, http.StatusNotFound, err)
	}

	// TODO: rewrite after Roles and Companies db interfaces realization
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

func (s *server) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateJob<-Close: ")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	job := new(model.Job)
	err := decoder.Decode(job)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateJob<-Decode: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateJob<-GetUserFromRequest: ")
		s.error(w, r, codeStatus, err)
		return
	}

	if !user.IsManager() {
		err = errors.New("HandleCreateJob:current user is not a manager : ")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.store.Mu.Lock()
	manager, err := s.store.Manager().FindByUser(user.ID)
	s.store.Mu.Unlock()

	if err != nil {
		log.Println("fail find manager", err)
		err = errors.Wrapf(err, "HandleCreateJob<-FindByUser: ")
		s.error(w, r, http.StatusNotFound, err)
	}

	s.store.Mu.Lock()
	err = s.store.Job().Create(job, manager)
	s.store.Mu.Unlock()

	if err != nil {
		log.Println("fail create job", err)
		err = errors.Wrapf(err, "HandleCreateJob<-Create: ")
		s.error(w, r, http.StatusInternalServerError, err)
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Atoi(wrong id): ")
		s.error(w, r, http.StatusBadRequest, err)
	}

	job, err := s.store.Job().Find(id)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Find: ")
		s.error(w, r, http.StatusNotFound, err)
	}

	s.respond(w, r, http.StatusOK, &job)
}

func (s *server) HandleEditFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-GetUserFromRequest: ")
		s.error(w, r, codeStatus, err)
		return
	}

	freelancer, err := s.store.Freelancer().FindByUser(user.ID)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-FindByUser: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleEditFreelancer<-rBodyClose: ")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(freelancer)

	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-Decode: ")
		s.error(w, r, http.StatusBadRequest, errors.New("invalid format of data"))
		return
	}
	// TODO: validate freelancer

	err = s.store.Freelancer().Edit(freelancer)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditFreelancer<-Edit: ")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s *server) HandleGetFreelancer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-Atoi(wrong id): ")
		s.error(w, r, http.StatusBadRequest, err)
	}

	freelancer, err := s.store.Freelancer().Find(id)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetFreelancer<-Find: ")
		s.error(w, r, http.StatusNotFound, err)
	}

	s.respond(w, r, http.StatusOK, &freelancer)
}

func (s *server) HandleGetAvatar(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetAvatar<-Atoi(wrong id): ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := s.store.User().Find(id)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetAvatar<-userFind: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	if user.Avatar != nil {
		image := user.Avatar
		if _, err := w.Write(image); err != nil {
			err = errors.Wrapf(err, "HandleGetAvatar<-Write: ")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("Content-Type", "multipart/form-data")
		w.Header().Set("Content-Length", strconv.Itoa(len(image)))
	} else {
		filename := "internal/store/avatars/default.png"

		var openFile *os.File
		openFile, err = os.Open(filename)
		if err != nil {
			err = errors.Wrapf(err, "HandleGetAvatar<-Open : ")
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		defer func() {
			if err := openFile.Close(); err != nil {
				err = errors.Wrapf(err, "HandleGetAvatar<-Close: ")
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}()

		fileHeader := make([]byte, 100000)
		_, err := openFile.Read(fileHeader)
		if err != nil {
			err = errors.Wrapf(err, "HandleGetAvatar<-Read: ")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		fileContentType := http.DetectContentType(fileHeader)
		fileStat, _ := openFile.Stat()
		fileSize := strconv.FormatInt(fileStat.Size(), 10)
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", fileContentType)
		w.Header().Set("Content-Length", fileSize)

		_, err = openFile.Seek(0, 0)
		if err != nil {
			err = errors.Wrapf(err, "HandleGetAvatar<-Seek: ")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		_, err = io.Copy(w, openFile)
		if err != nil {
			err = errors.Wrapf(err, "HandleGetAvatar<-ioCopy: ")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}


func (s *server) HandleGetAllJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	jobs, err := s.store.Job().GetAllJobs()
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Find: ")
		s.error(w, r, http.StatusNotFound, err)
	}

	s.respond(w, r, http.StatusOK, &jobs)
}

func (s *server) HandleUpdateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	inputJob := new(model.Job)
	err := decoder.Decode(inputJob)

	// Validate Job
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Atoi(wrong type id): ")
		s.error(w, r, http.StatusBadRequest, err)
	}

	job, err := s.store.Job().Find(id)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetJob<-Find: ")
		s.error(w, r, http.StatusNotFound, err)
	}
	inputJob.ID = job.ID
	inputJob.HireManagerId = job.HireManagerId
	err = s.store.Job().Edit(inputJob)
	if err != nil {
		err = errors.Wrapf(err, "HandleEditProfile<-JobEdit")
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	s.respond(w, r, http.StatusOK, struct {}{})
}
