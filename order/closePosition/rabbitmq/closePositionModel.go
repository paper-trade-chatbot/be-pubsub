package rabbitmq

import "github.com/shopspring/decimal"

type TradeType int

const (
	TradeType_None TradeType = iota
	TradeType_Buy
	TradeType_Sell
)

type ClosePositionModel struct {
	ID           uint64          `valid:"required" json:"id"`
	MemberID     uint64          `valid:"required" json:"memberID"`
	PositionID   uint64          `valid:"required" json:"positionID"`
	ExchangeCode string          `valid:"required" json:"exchangeCode"`
	ProductCode  string          `valid:"required" json:"productCode"`
	TradeType    TradeType       `valid:"required" json:"tradeType"`
	OpenPrice    decimal.Decimal `valid:"required" json:"openPrice"` // 開倉價
	CloseAmount  decimal.Decimal `valid:"required" json:"amount"`    // 為正數，要平多少數量的倉
}
