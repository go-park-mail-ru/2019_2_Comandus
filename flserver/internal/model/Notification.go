package model

type Notification struct {
	UserID int `json:"-"`
	NewMessages bool `json:"new_messages"`
	NewProjects bool `json:"new_projects"`
	NewsFromService bool `json:"news_service"`
}
