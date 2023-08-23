package domain

import "errors"

var (
	ErrUserNotFound            = errors.New("user doesn't exists")
	ErrUserAlreadyExists       = errors.New("user with such email already exists")
	ErrUnknownCallbackType     = errors.New("unknown callback type")
	ErrSendPulseIsNotConnected = errors.New("sendpulse is not connected")
)
