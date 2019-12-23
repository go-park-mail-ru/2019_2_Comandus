package model

type HttpError struct {
	ClientErr error
	LogErr    error
	HttpCode  int
}
