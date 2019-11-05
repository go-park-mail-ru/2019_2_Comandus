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
			err = errors.Wrapf(err, "HandleCreateUser:<-Body.Close")
			s.error(w, r, http.StatusInternalServerError, err)
		}
	}()

	decoder := json.NewDecoder(r.Body)
	user := new(model.User)
	err := decoder.Decode(user)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-Decode:")
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

	err = s.store.User().Create(user)

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

	err = s.store.Freelancer().Create(&f)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateUser<-CreateFreelancer:")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	m := model.HireManager{
		AccountID:        user.ID,
		RegistrationDate: time.Now(),
	}

	err = s.store.Manager().Create(&m)

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

		u, err := s.store.User().Find(id.(int64))

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

	u, err := s.store.User().FindByEmail(user.Email)

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

	err = s.store.User().Edit(user)

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

	err = s.store.User().Edit(user)
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

	err = s.store.User().Edit(user)

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
	uid := uidInterface.(int64)

	user, err := s.store.User().Find(uid)

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

	err = s.store.User().Edit(user)

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

	_, err = s.store.Manager().FindByUser(user.ID)
	if err != nil {
		err = errors.Wrapf(err, "HandleRoles<-ManagerFind:")
		s.error(w, r, http.StatusNotFound, err)
	}

	// TODO: rewrite after Roles and Companies db interfaces realization

	//company := s.usersdb.Companies[hireManager.ID]

	company := model.Company{
		ID:          0,
		CompanyName: "default company",
		Site:        "company.ru",
		Description: "default company",
		Country:     "Russia",
		City:        "Moscow",
		Address:     "Red square street",
		Phone:       "88888888888",
	}

	var roles []*model.Role
	clientRole := &model.Role{
		Role:   "client",
		Label:  company.CompanyName,
		Avatar: "../store//default.png",
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

	manager, err := s.store.Manager().FindByUser(user.ID)

	if err != nil {
		log.Println("fail find manager", err)
		err = errors.Wrapf(err, "HandleCreateJob<-FindByUser: ")
		s.error(w, r, http.StatusNotFound, err)
	}

	err = s.store.Job().Create(job, manager)
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

	job, err := s.store.Job().Find(int64(id))
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

	freelancer, err := s.store.Freelancer().Find(int64(id))
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

	user, err := s.store.User().Find(int64(id))
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

	jobs, err := s.store.Job().List()
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

	job, err := s.store.Job().Find(int64(id))
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

func (s *server) HandleResponseJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
	}
	jobId := int64(id)

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-GetUserFromRequest: ")
		s.error(w, r, codeStatus, err)
		return
	}

	if user.IsManager() {
		err = errors.Wrapf(errors.New("to response user need to be freelancer"),
			"HandleResponseJob<-IsManager: ")
		s.error(w, r, codeStatus, err)
		return
	}

	freelancer, err := s.store.Freelancer().FindByUser(user.ID)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-Freelancer().FindByUser: ")
		s.error(w, r, codeStatus, err)
		return
	}

	// TODO: get files from request
	response := model.Response{
		ID:               0,
		FreelancerId:     freelancer.ID,
		JobId:            jobId,
		Files:            "",
		Date:             time.Now(),
		StatusManager:    model.ResponseStatusReview,
		StatusFreelancer: model.ResponseStatusBlock,
	}

	if err := response.Validate(0); err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-Validate: ")
		s.error(w, r, http.StatusBadRequest, err)
	}

	if err := s.store.Response().Create(&response); err != nil {
		err = errors.Wrapf(err, "HandleResponseJob<-Response().Create")
		s.error(w, r, http.StatusInternalServerError, err)
	}

	s.respond(w, r, http.StatusOK, struct {}{})
}


// TODO : to another file
func (s * server) getManagerResponses(userId int64) (*[]model.Response, error){
	manager, err := s.store.Manager().FindByUser(userId)
	if err != nil {
		err = errors.Wrapf(err, " GetManagerResponses<-Manager().FindByUser: ")
		return nil, err
	}

	responses, err := s.store.Response().ListForManager(manager.ID)
	if err != nil {
		err = errors.Wrapf(err, "GetManagerResponses<-Responses().ListForManager: ")
		return nil, err
	}
	return &responses, nil
}

func (s * server) getFreelancerResponses(userId int64) (*[]model.Response, error){
	freelancer, err := s.store.Freelancer().FindByUser(userId)
	if err != nil {
		err = errors.Wrapf(err, " GetManagerResponses<-Manager().FindByUser: ")
		return nil, err
	}

	responses, err := s.store.Response().ListForManager(freelancer.ID)
	if err != nil {
		err = errors.Wrapf(err, "GetManagerResponses<-Responses().ListForManager: ")
		return nil, err
	}
	return &responses, nil
}

func (s *server) HandleGetResponses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetResponses<-GetUserFromRequest: ")
		s.error(w, r, codeStatus, err)
		return
	}

	var responses *[]model.Response
	if user.IsManager() {
		responses, err = s.getManagerResponses(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleGetResponses<-GetManagerResponses: ")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	} else {
		responses, err = s.getFreelancerResponses(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleGetResponses<-GetFreelancerResponses: ")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
	s.respond(w, r, http.StatusOK, responses)
}

func (s * server) HandleResponseAccept(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-GetUserFromRequest: ")
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	response, err := s.store.Response().Find(responseId)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-Response().Find(): ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	if user.IsManager() {
		manager, err := s.store.Manager().FindByUser(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Manager().FindByUser: ")
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		job, err := s.store.Job().Find(response.JobId)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Job.Find: ")
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if job.HireManagerId != manager.ID {
			err = errors.New("current manager cant accept this response")
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		response.StatusManager = model.ResponseStatusAccepted
		response.StatusFreelancer = model.ResponseStatusReview
	} else {
		freelancer, err := s.store.Freelancer().FindByUser(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Freelancer().FindByUser: ")
			s.error(w, r, codeStatus, err)
			return
		}

		if freelancer.ID != response.FreelancerId {
			err = errors.New("current freelancer can't accept this response")
			s.error(w, r, codeStatus, err)
			return
		}

		if response.StatusFreelancer == model.ResponseStatusBlock {
			err = errors.New("freelancer can't accept response before manager")
			s.error(w, r, codeStatus, err)
			return
		}

		response.StatusManager = model.ResponseStatusAccepted
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s * server) HandleResponseDeny(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-GetUserFromRequest: ")
		s.error(w, r, http.StatusUnauthorized, err)
		return
	}

	response, err := s.store.Response().Find(responseId)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept<-Response().Find(): ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	if user.IsManager() {
		manager, err := s.store.Manager().FindByUser(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Manager().FindByUser: ")
			s.error(w, r, http.StatusNotFound, err)
			return
		}
		job, err := s.store.Job().Find(response.JobId)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Job.Find: ")
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		if job.HireManagerId != manager.ID {
			err = errors.New("current manager cant accept this response")
			s.error(w, r, http.StatusNotFound, err)
			return
		}

		response.StatusManager = model.ResponseStatusDenied
		response.StatusFreelancer = model.ResponseStatusBlock
	} else {
		freelancer, err := s.store.Freelancer().FindByUser(user.ID)
		if err != nil {
			err = errors.Wrapf(err, "HandleResponseAccept<-Freelancer().FindByUser: ")
			s.error(w, r, codeStatus, err)
			return
		}

		if freelancer.ID != response.FreelancerId {
			err = errors.New("current freelancer can't accept this response")
			s.error(w, r, codeStatus, err)
			return
		}

		if response.StatusFreelancer == model.ResponseStatusBlock {
			err = errors.New("freelancer can't accept response before manager")
			s.error(w, r, codeStatus, err)
			return
		}

		response.StatusManager = model.ResponseStatusDenied
	}

	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s * server) HandleCreateContract(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer func() {
		if err := r.Body.Close(); err != nil {
			err = errors.Wrapf(err, "HandleCreateContract:")
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}()

	decoder := json.NewDecoder(r.Body)
	contract := new(model.Contract)
	err := decoder.Decode(contract)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleResponseAccept-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}
	responseId := int64(id)

	response, err := s.store.Response().Find(responseId)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-Response().Find: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	job, err := s.store.Job().Find(response.JobId)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-Job().Find: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	manager, err := s.store.Manager().Find(job.HireManagerId)
	if err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-Manager().Find: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	contract.ResponseID = response.ID
	contract.CompanyID = manager.CompanyID
	contract.FreelancerID = response.FreelancerId
	contract.Status = model.ContractStatusUnderDevelopment

	//TODO: uncommnet when fix response
	//contract.PaymentAmount = response.PaymentAmount
	contract.Grade = 0

	if err := s.store.Contract().Create(contract); err != nil {
		err = errors.Wrapf(err, "HandleCreateContract<-Contract().Create: ")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	// TODO: other responses should be denied
	s.respond(w, r, http.StatusOK, struct{}{})
}

func (s * server) SetStatusContract(user * model.User, contract *model.Contract, status string) error {
	// TODO: fix if add new modes
	if !user.IsManager() && status != model.ContractStatusDone {
		err := errors.New("freelancer can change status only to done status")
		return errors.Wrapf(err, "SetStatusContract<-GetUserFromRequest:")
	}

	contract.Status = status
	if err := s.store.Contract().Edit(contract); err != nil {
		return errors.Wrapf(err, "SetStatusContract<-Contract().Edit:")
	}
	return nil
}

func (s * server) HandleTickContractAsDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user, err, status := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-GetUserFromRequest: ")
		s.error(w, r, status, err)
		return
	}

	if user.IsManager() {
		err = errors.New("user must be freelancer")
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	freelancer, err := s.store.Freelancer().FindByUser(user.ID)
	if err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-Freelancer().FindByUser: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}

	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	if err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-strconv.Atoi: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	contractId := int64(id)
	contract, err := s.store.Contract().Find(contractId)
	if err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-Contract().Find: ")
		s.error(w, r, http.StatusNotFound, err)
		return
	}


	if contract.FreelancerID != freelancer.ID {
		err = errors.New("current freelancer can't manage this contract")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if err := s.SetStatusContract(user, contract, model.ContractStatusDone); err != nil {
		err = errors.Wrapf(err, "HandleTickContractAsDone<-SetStatusContract: ")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	s.respond(w,r, http.StatusOK, struct{}{})
}

func (s * server) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	sess , err := s.sessionStore.Get(r, sessionName)
	if err != nil {
		err = errors.Wrapf(err, "HandleGetToken<-sessionStore.Get :")
		s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
		return
	}

	token, err := s.token.Create(sess, time.Now().Add(24*time.Hour).Unix())
	s.respond(w, r, http.StatusOK, map[string]string{"csrf-token" : token})
}