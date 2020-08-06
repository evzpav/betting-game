package storage

import "gitlab.com/evzpav/betting-game/pkg/errors"

const (
	ErrUserNotFound   errors.Code = "USER_NOT_FOUND"
	ErrUserDuplicated errors.Code = "USER_DUPLICATED"
)
