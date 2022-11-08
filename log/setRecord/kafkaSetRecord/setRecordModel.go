package kafkaSetRecord

type SetRecordModel struct {
	Time     int64  `valid:"required" json:"time"`
	Account  string `valid:"required" json:"account"`
	Msgtype  int32  `valid:"required" json:"msgtype"`
	Detail   string `valid:"required" json:"detail"`
	Modifier string `valid:"required" json:"modifier"`
	Ip       string `valid:"-" json:"ip"`
}

type DetailModel struct {
	Type            int64   `valid:"required" json:"type"`
	Amount          float64 `valid:"required" json:"amount"`
	Gateway         string  `valid:"required" json:"gateway"`
	AccountRole     string  `valid:"required" json:"accountRole"`
	MainAccount     string  `valid:"required" json:"mainAccount"`
	SubAccount      string  `valid:"required" json:"subAccount"`
	MainAccountRole string  `valid:"required" json:"mainAccountRole"`
	SubAccountRole  string  `valid:"required" json:"subAccountRole"`
}
