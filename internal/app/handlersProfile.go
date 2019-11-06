package apiserver

import (
	"bytes"
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
)

func (s *server) HandleShowProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err, codeStatus := s.GetUserFromRequest(r)
	if err != nil {
		err = errors.Wrapf(err, "HandleShowProfile<-GetUserFromRequest:")
		s.error(w, r, codeStatus, err)
		return
	}
	user.Sanitize(s.sanitizer)
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
	userInput := new(model.User)
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(userInput)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		err = errors.Wrapf(err, "HandleEditProfile<-Decode:")
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	userInput.ID = user.ID
	userInput.Email = user.Email

	err = s.store.User().Edit(userInput)
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
	uid := uidInterface.(int)

	user, err := s.store.User().Find(int64(uid))

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