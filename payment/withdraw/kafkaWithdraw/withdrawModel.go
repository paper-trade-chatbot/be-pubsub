package kafkaWithdraw

type WithdrawModel struct {
	Account     string `valid:"required" json:"account"`
	Amount      string `valid:"required" json:"amount"`
	BankName    string `valid:"required" json:"bank_name"`
	BankAccount string `valid:"required" json:"bank_account"`
}
