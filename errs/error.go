package errs

import "errors"

var ErrNoUesrWitThisId = errors.New("no user with this id")
var ErrNoUesrWitThisEmail = errors.New("no user with this id")
var ErrNoUesrWitThisUsername = errors.New("no user with this username")
var ErrNoUserWithId = errors.New("no user with this id")
var ErrNoUserWithUsername = errors.New("no username with this username")
var ErrNoUserWithEmail = errors.New("un registered email")
var ErrFailToDecode = errors.New("fail to decode")
var ErrCursorDuringItr = errors.New("error occur during iteration")
var ErrUnidentifiedToken = errors.New("unidentified token")
var ErrIsnvalidID = errors.New("invalied id")
var ErrIncorrectPassword = errors.New("incorrect password")
var ErrUserExistWithThisEmail = errors.New("user exist witht this email")
var ErrUserExistWithThisUsername = errors.New("user exist witht this username")
var ErrCantCreateUser = errors.New("can't create user")
var ErrSuccess = errors.New("succesfully registered")
var ErrFailToCreateUser = errors.New("fail to create user")
var ErrFailToDelete = errors.New("fail to delert the doucument")
var ErrAccountNotActivated = errors.New("account not activated")
var ErrFailToSendEmail = errors.New("failed to send an email")
