package rabbitmq

type NewPositionModel struct {
	Account         string `valid:"required" json:"account"`
	Amount          string `valid:"required" json:"amount"`
	Gateway         string `valid:"required" json:"gateway"`
	PaymentRecordId uint64 `valid:"required" json:"paymentRecordId"`
	Status          string `valid:"required" json:"status"`
}
