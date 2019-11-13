package model

type BodyPassword struct {
	Password string `valid:"required"`
	NewPassword string `valid:"required, length(6|100)"`
	NewPasswordConfirmation string `valid:"required"`
}

// TODO: delete
type InnerInfo struct {
	UserID int64 `json:"user_id"`
	WhoSeeProfile string `json:"who_see_profile"`
	ControlQuestion string `json:"control_question"`
	ControlAnswer string `json:"-"`
}

type Notification struct {
	UserID int64 `json:"-"`
	NewMessages bool `json:"new_messages"`
	NewProjects bool `json:"new_projects"`
	NewsFromService bool `json:"news_service"`
}

