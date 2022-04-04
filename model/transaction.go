package model

type ActionValue string

const (
	ActionDeposit  ActionValue = "deposit"
	ActionWithdraw ActionValue = "withdraw"
)

type Transaction struct {
	Action ActionValue `json:"action" validate:"required,oneof='deposit''withdraw'"`
	Fund   float64     `json:"fund" validate:"required,gte=0"`
}
