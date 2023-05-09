# be-pubsub


## how to use with RABBITMQ  


### register receiver
```go
if sub, err := rabbitmqOpenPosition.SubscribeAndListen(
	ctx,
	config.GetString("RABBITMQ_USERNAME"),
	config.GetString("RABBITMQ_PASSWORD"),
	config.GetString("RABBITMQ_HOST"),
	config.GetString("RABBITMQ_VIRTUAL_HOST"),
	config.GetString("SERVICE_NAME"),
	<YOUR_FUNCTION>,
); err != nil {
	logging.Error(ctx, "SubscribeAndListen error %v", err)
	panic(err)
} else {
	subscribers = append(subscribers, sub)
}
```
 

### register sender
```go
if publisher, err := rabbitmqOpenPosition.NewPublisher(
	config.GetString("RABBITMQ_USERNAME"),
	config.GetString("RABBITMQ_PASSWORD"),
	config.GetString("RABBITMQ_HOST"),
	config.GetString("RABBITMQ_VIRTUAL_HOST")); err == nil {
	if err = registerPublisher[*rabbitmqOpenPosition.OpenPositionModel](publisher); err != nil {
		logging.Error(ctx, "registerPublisher error %v", err)
	}
} else {
	logging.Error(ctx, "NewPublisher error %v", err)
	panic(err)
}
```

### send
```go
message := &rabbitmqOpenPosition.OpenPositionModel{
	ID:           model.ID,
	MemberID:     in.MemberID,
	ExchangeCode: in.ExchangeCode,
	ProductCode:  in.ProductCode,
	TradeType:    rabbitmqOpenPosition.TradeType(in.TradeType),
	Amount:       amount,
}

_, err = pubsub.GetPublisher[*rabbitmqOpenPosition.OpenPositionModel](ctx).Produce(ctx, message)
```

