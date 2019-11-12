package model

type BodyPassword struct {
	Password string `valid:"required"`
	NewPassword string `valid:"required, length(6|100)"`
	NewPasswordConfirmation string `valid:"required"`
}