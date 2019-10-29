package model

type BodyPassword struct {
	Password string `valid:"required"`
	NewPassword string `valid:"required"`
	NewPasswordConfirmation string `valid:"required"`
}