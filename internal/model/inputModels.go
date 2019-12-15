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

type SearchParams struct {
	MinGrade         int64   `json:"minGrade"`
	MaxGrade         int64   `json:"maxGrade"`
	MinPaymentAmount float64 `json:"minPaymentAmount"`
	MaxPaymentAmount float64 `json:"maxPaymentAmount"`
	Country          int64   `json:"country"`
	City             int64   `json:"city"`
	Proposals        int64   `json:"proposals"`
	ExperienceLevel  [3]bool `json:"experienceLevel"`
	Desc             bool    `json:"desc"`
	Limit			 int64	 `json:"limit"`
}