package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2019_2_Comandus/internal/model"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

func (s *server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	s.usersdb.Mu.Lock()
	var id int
	var idf int
	var idc int

	if len(s.usersdb.Users) > 0 {
		id = s.usersdb.Users[len(s.usersdb.Users)-1].ID + 1
		idf = s.usersdb.Freelancers[len(s.usersdb.Freelancers)-1].ID + 1
		idc = s.usersdb.HireManagers[len(s.usersdb.HireManagers)-1].ID + 1
	}

	user := model.User{
		ID:              id,
		FirstName:            newUserInput.Name,
		Email:           newUserInput.Email,
		Password:        newUserInput.Password,
	}

	err = user.BeforeCreate()
	if err != nil {
		s.respond(w, r, http.StatusInternalServerError, newUserInput)
	}
	s.usersdb.Users = append(s.usersdb.Users, user)


	s.usersdb.Freelancers = append(s.usersdb.Freelancers, model.Freelancer {
		ID:       idf,
		AccountId: id,
	})

	s.usersdb.HireManagers = append(s.usersdb.HireManagers, model.HireManager {
		ID:       idc,
		AccountID : id,
	})

	fmt.Println(s.usersdb.Users[id])
	s.usersdb.Mu.Unlock()
	s.respond(w, r, http.StatusCreated, newUserInput)
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		var u *model.User
		var found bool

		s.usersdb.Mu.Lock()
		for i := 0; i < len(s.usersdb.Users); i++ {
			if id == s.usersdb.Users[i].ID {
				u = &s.usersdb.Users[i]
				found = true
			}
		}
		s.usersdb.Mu.Unlock()

		if !found {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) HandleSessionCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	fmt.Println(newUserInput)

	s.usersdb.Mu.Lock()
	for i:=0; i < len(s.usersdb.Users); i++ {
		if s.usersdb.Users[i].Email == newUserInput.Email &&
			s.usersdb.Users[i].ComparePassword(newUserInput.Password) {

			u := s.usersdb.Users[i]
			session, err := s.sessionStore.Get(r, sessionName)
			if err != nil {
				s.usersdb.Mu.Unlock()
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			session.Values["user_id"] = u.ID
			session.Values["user_type"] = s.userType
			//session.Values["user_type"] = userFreelancer
			if err := s.sessionStore.Save(r, w, session); err != nil {
				s.usersdb.Mu.Unlock()
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.usersdb.Mu.Unlock()
			s.respond(w, r, http.StatusOK, nil)
			return
		}
	}
	s.usersdb.Mu.Unlock()
	s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
}

func (s *server) HandleLogout(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, sessionName)
	fmt.Println(session)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if session == nil {
		s.error(w, r, http.StatusNotFound, errors.New("failed to delete session"))
		return
	}

	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		s.error(w, r, http.StatusExpectationFailed, errors.New("failed to delete session"))
	}
	fmt.Println("logout")
	http.Redirect(w, r, "/", http.StatusUnauthorized)
}

func (s * server) HandleSetUserType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//s.respond(w,r, http.StatusOK, nil)
	//return
	// TODO check if input user type invalid
	type Input struct {
		UserType     string `json:"type"`
	}
	defer func() {
		// TODO: handle err
		r.Body.Close()
	}()

	decoder := json.NewDecoder(r.Body)
	newInput := new(Input)
	err := decoder.Decode(newInput)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Values["user_type"] = newInput.UserType
	session.Save(r,w)
}

func (s *server) HandleShowProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user , SendErr , CodeStatus := s.GetUserFromRequest(r)
	if SendErr != nil {
		s.error(w, r, CodeStatus, SendErr)
		return
	}
	s.respond(w, r, http.StatusOK, user)
}

func (s *server) HandleEditProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user , SendErr , CodeStatus := s.GetUserFromRequest(r)
	if SendErr != nil {
		s.error(w, r, CodeStatus, SendErr)
		return
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(user)
	fmt.Println(user)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		SendErr := fmt.Errorf("invalid format of data")
		s.error(w, r, http.StatusBadRequest, SendErr)
		return
	}
	s.respond(w, r, http.StatusOK, nil)

}

func (s *server) HandleEditPassword(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var err error
	user , sendErr, codeStatus := s.GetUserFromRequest(r)
	if sendErr != nil {
		s.error(w, r, codeStatus, sendErr)
		return
	}
	currPassword := r.FormValue("password")
	newPassword := r.FormValue("newPassword")
	newPasswordConfirmation := r.FormValue("newPasswordConfirmation")
	if newPassword != newPasswordConfirmation {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("new Passwords are different"))
	}
	if user.Password != currPassword {
		s.error(w, r, http.StatusBadRequest, fmt.Errorf("wrong password"))
	}
	newEncryptPassword , err := model.EncryptString(newPassword)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError , fmt.Errorf("error in updating password"))
	}
	user.Password = newPassword
	user.EncryptPassword = newEncryptPassword
	s.respond(w, r, http.StatusOK, nil)
}

func (s *server) GetUserFromRequest (r *http.Request) (*model.User , error , int) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		SendErr := fmt.Errorf( "user isn't authorized")
		return nil , SendErr , http.StatusUnauthorized
	}
	uidInteface := session.Values["user_id"]
	uid := uidInteface.(int)

	s.usersdb.Mu.Lock()
	user := s.usersdb.GetUserByID(uid)
	s.usersdb.Mu.Unlock()

	if user == nil {
		SendErr := fmt.Errorf( "can't find user with id:" + strconv.Itoa(int(uid)))
		return nil , SendErr , http.StatusBadRequest
	}
	return user, nil , http.StatusOK
}

// TODO:
func (s * server) HandleEditNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	user , sendErr, codeStatus := s.GetUserFromRequest(r)
	if sendErr != nil {
		s.error(w, r, codeStatus, sendErr)
		return
	}
	userNotification := s.usersdb.GetNotificationsByUserID(user.ID)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(userNotification)
	fmt.Println(user)
	if err != nil {
		log.Printf("error while marshalling JSON: %s", err)
		SendErr := fmt.Errorf("invalid format of data")
		s.error(w, r, http.StatusBadRequest, SendErr)
		return
	}
	s.respond(w, r, http.StatusOK, nil)

}

func (s *server) HandleUploadAvatar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("File Upload Endpoint Hit")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, errors.New("error retrieving the file"))
		return
	}

	file, _, err := r.FormFile("myFile")
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errors.New("error retrieving the file"))
		return
	}
	defer file.Close()

	tempFile, err := ioutil.TempFile("internal/store/avatars", "upload-*.png")
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, errors.New("error creating the file"))
		return
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	uidInterface := session.Values["user_id"]
	uid := uidInterface.(int)
	s.usersdb.Mu.Lock()
	s.usersdb.Users[uid].Avatar = tempFile.Name()
	s.usersdb.Mu.Unlock()
}


func (s *server) HandleDownloadAvatar(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, sessionName)
	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	uidInterface := session.Values["user_id"]
	uid := uidInterface.(int)

	s.usersdb.Mu.Lock()
	Filename := s.usersdb.Users[uid].Avatar
	s.usersdb.Mu.Unlock()

	if Filename == "" {
		Filename = "/internal/store/avatars/default.png"
	}
	fmt.Println("Client requests: " + Filename)

	Openfile, err := os.Open(Filename)
	defer Openfile.Close()
	if err != nil {
		s.error(w,r,http.StatusNotFound, errors.New("cant open file"))
		return
	}

	FileHeader := make([]byte, 100000) // max image size!!!
	Openfile.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)

	FileStat, _ := Openfile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	Openfile.Seek(0, 0)
	io.Copy(w, Openfile)
}


func (s * server) HandleGetAuthHistory(w http.ResponseWriter, r *http.Request) {

}

func (s * server) HandleGetSecQuestion(w http.ResponseWriter, r *http.Request) {

}

func (s * server) HandleEditSecQuestion(w http.ResponseWriter, r *http.Request) {

}

func (s * server) HandleCheckSecQuestion(w http.ResponseWriter, r *http.Request) {

}

/*func (s * server) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)

	s.usersdb.mu.Lock()
	err := encoder.Encode(s.usersdb.users)
	s.usersdb.mu.Unlock()


	if err != nil {
		s.error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
}*/
