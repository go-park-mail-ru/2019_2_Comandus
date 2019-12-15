package model

type ReviewInput struct {
	Grade   int    `json:"grade"`
	Comment string `json:"comment"`
}

type BodyPassword struct {
	Password                string `json:"password" valid:"required"`
	NewPassword             string `json:"newPassword" valid:"required, length(6|100)"`
	NewPasswordConfirmation string `json:"newPasswordConfirmation" valid:"required"`
}

// TODO: delete
type InnerInfo struct {
	UserID          int64  `json:"user_id"`
	WhoSeeProfile   string `json:"who_see_profile"`
	ControlQuestion string `json:"control_question"`
	ControlAnswer   string `json:"-"`
}

type Notification struct {
	UserID          int64 `json:"-"`
	NewMessages     bool  `json:"new_messages"`
	NewProjects     bool  `json:"new_projects"`
	NewsFromService bool  `json:"news_service"`
}
