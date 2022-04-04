package pkg

var (
	ErrWalletNotFound = "wallet not found"
	ErrWalletFunds    = "wallet funds not enough"
	ErrWalletID       = "invalid wallet id"
	ErrUserID         = "invalid user id"
	ErrInvalidAction  = "unavailable action"
	ErrWalletDeposit  = "cannot deposit negative funds"
	ErrWalletWithdraw = "cannot withdraw negative funds"
	ErrDuplicated     = "user already has a wallet"
)

// HttpError represents http server error
type HTTPError interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code
type StatusError struct {
	Code   int
	ErrMsg string
}

func (s StatusError) Error() string {
	return s.ErrMsg
}

func (s StatusError) Status() int {
	return s.Code
}
