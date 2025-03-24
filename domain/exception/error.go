package exception

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrUserAlreadyExist = errors.New("user already exist")

	ErrUserNoPassword    = errors.New("account has no password, probably using social login")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrFileAlreadyExists = errors.New("file already exists")
	ErrFileSizeLimit     = errors.New("file size limit exceeded")
	ErrTypeConversion    = errors.New("could not convert variable type")
	ErrFileIsNull        = errors.New("file could not be null")
	ErrFileCountLimit    = errors.New("file count limit exceeded")

	ErrPermissionDenied     = errors.New("no permssion")
	ErrParamsMissing        = errors.New("some params are missing")
	ErrCouldNotUpdateStatus = errors.New("could not update status yet")
	ErrJWTSecretIsEmpty     = errors.New("jwt secret is empty")

	ErrProductRequestAlreadyPaid = errors.New("product request already paid")

	ErrDuplicateTravelerPayoutAccount = errors.New("duplicate traveler payout account")
	ErrBankNotFound                   = errors.New("bank not found")

	ErrOfferNotFound     = errors.New("offer not found")
	ErrCouldNotSelfOffer = errors.New("could not commit self offering")
)
