package rabbitmq

import "github.com/shopspring/decimal"

type TradeType int

const (
	TradeType_None TradeType = iota
	TradeType_Buy
	TradeType_Sell
)

type OpenPositionModel struct {
	ID           uint64          `valid:"required" json:"id"`
	MemberID     uint64          `valid:"required" json:"memberID"`
	ExchangeCode string          `valid:"required" json:"exchangeCode"`
	ProductCode  string          `valid:"required" json:"productCode"`
	TradeType    TradeType       `valid:"required" json:"tradeType"`
	Amount       decimal.Decimal `valid:"required" json:"amount"`
}
