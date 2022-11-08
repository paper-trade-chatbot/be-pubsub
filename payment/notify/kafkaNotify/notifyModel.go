package kafkaNotify

type NotifyModel struct {
	Account     string `valid:"required" json:"account"`
	Amount      string `valid:"required" json:"amount"`
	Action      string `valid:"required" json:"action"`
	BankName    string `valid:"-" json:"bankName"`
	BankAccount string `valid:"-" json:"bankAccount"`
	Gateway     string `valid:"-" json:"gateway"`
	Email       string `valid:"-" json:"email"`
	PhoneNumber string `valid:"-" json:"phoneNumber"`
	FcmToken    string `valid:"-" json:"fcmToken"`
}
