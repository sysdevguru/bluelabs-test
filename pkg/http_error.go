package pkg

var (
	ErrWalletNotFound = "wallet not found"
	ErrWalletBalance  = "wallet balance not enough"
	ErrWalletID       = "invalid wallet id"
	ErrUserID         = "invalid user id"
	ErrInvalidAction  = "unavailable action"
	ErrWalletFund     = "cannot update balance with nagetive fund"
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
