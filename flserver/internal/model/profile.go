package model

import (
	"net/url"
	"time"
	"mime/multipart"
)

type Profile struct {
	ID int `json:"profile_id"`
	CompanyName string `json:"company_name, omitempty"`
	FirstName string `json:"first_name"`
	SecondName string `json:"second_name"`
	Email string `json:"email"`
	Specializations []string `json:"specializations"`
	ContactInformation ContactInfo `json:"contact_info"`
	AdditionalInfo AdditionalInformation `json:"additional_info"`
	InnerInformation InnerInfo `json:"inner_info"`
	Avatar multipart.File  `json:"avatar"`
//	What about Rating
}

type ContactInfo struct {
	Owner string  `json:"owner, omitempty"`
	Country string `json:"country"`
	City string `json:"city"`
	Address string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type InnerInfo struct {
	WhoSeeProfile string `json:"who_see_profile"`
	ExperienceLevel string `json:"experience_level"`
}

type AdditionalInformation struct {
	BirthdayDate time.Time `json:"birthday_date , omitempty"`
	Language string `json:"language"`
	PersonalSite url.URL `json:"personal_site"`
	Status string `json:"status"`
	Description string `json:"description"`
	Specializations []string `json:"specializations"`
	ControlQuestion string `json:"control_question"`
	ControlAnswer string `json:"-"`
}

type NotificationInfo struct {
	NewMessages bool `json:"new_messages"`
	NewProjects bool `json:"new_projects"`
	NewsFromService bool `json:"news_service"`
}
