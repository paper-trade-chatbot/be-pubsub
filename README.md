# be-pubsub



## unable to import private repo
>  https://blog.wu-boy.com/2020/03/read-private-module-in-golang/
> 
>  git config --global --add url."git@github.com:".insteadOf "https://github.com/"
>  go env -w GOPRIVATE=github.com/paper-trade-chatbot

## how to know topic name?

```go
temp, err := xxx.NewSubscriber()

fmt.Println(temp.GetTopicName())
```

---

## how to use this lib

### publisher

1. initialize publishers in main
```go
import (
	"context"
	"reflect"
	"sync"

	"github.com/Funzhou-tech/rf-payment/config"
	"github.com/Funzhou-tech/rf-payment/logging"
	rfPubsub "github.com/paper-trade-chatbot/be-pubsub"
)

var publishers = map[string]rfPubsub.Publisher{}
var publisherLock sync.RWMutex

//Initialize
// please register all instance creator of publisher here
func Initialize(ctx context.Context) {

	// ==============================
	// |   initialize publishers    |
	// ==============================
	func() {
		publisherLock.Lock()
		defer publisherLock.Unlock()
		if withdraw, err := kafkaWithdraw.GetPublisher([]string{config.GetString("KAFKA")}); err != nil {
			registerPublisher(payment.WithdrawModel{}, withdraw)
		}
	}()
}

func Finalize(ctx context.Context) {
	publisherLock.Lock()
	defer publisherLock.Unlock()
	for k, v := range publishers {
		err := v.Close()
		if err != nil {
			logging.Error(ctx, "pubsub Finalize error %v", err)
		}
		delete(publishers, k)
	}
}

func GetPublisher(ctx context.Context, model interface{}) rfPubsub.Publisher {

	if !publisherLock.TryRLock() {
		logging.Error(ctx, "GetPublisher: not initialized yet.")
		return nil
	}
	defer publisherLock.RUnlock()

	modelType := reflect.TypeOf(model).String()

	if _, ok := publishers[modelType]; !ok {
		logging.Error(ctx, "GetPublisher: no such publisher.", modelType)
		return nil
	}
	return publishers[modelType]
}

func registerPublisher(model interface{}, publisher rfPubsub.Publisher) error {
	publishers[reflect.TypeOf(model).String()] = publisher
	return nil
} 


```

2. call pulisher in your function

```go
pub := GetPublisher(ctx, payment.WithdrawModel{})
if pub == nil {
    // failed...
}

_, err := pub.Produce(context.Background(), payment.WithdrawModel{
    // data...
})
if err != nil {
    // failed...
}

```

### subscriber

1. setup your callback, you need to cast interface to your model

```go
func myCallback(model interface{}) error {
    myModel := model.(<model type>)
    ...
}

```

2. subscribe in main

```go

import (
	"context"
	"reflect"
	"sync"

	"github.com/Funzhou-tech/rf-payment/config"
	"github.com/Funzhou-tech/rf-payment/logging"
	rfPubsub "github.com/paper-trade-chatbot/be-pubsub"
	"github.com/paper-trade-chatbot/be-pubsub/model/payment"
	"github.com/paper-trade-chatbot/be-pubsub/payment/withdraw/kafkaWithdraw"
)

var subscribers = []rfPubsub.Subcriber{}

//Initialize
// please register all instance creator of publisher here
func Initialize(ctx context.Context) {

	// ==============================
	// |    register subscribers    |
	// ==============================

	if sub, err := kafkaWithdraw.SubscribeAndListen(
		config.GetString("PROJECT_NAME"),
		[]string{config.GetString("KAFKA")},
		&<your model>,
		myCallback,
	); err != nil {
		logging.Error(ctx, "SubscribeAndListen error %v", err)
	} else {
		subscribers = append(subscribers, sub)
	}

}

func Finalize(ctx context.Context) {

	for _, s := range subscribers {
		s.Close()
	}
}
```

---

## listen topics to debug

- make sure you are under the right environment

```bash
kubectl exec --tty -i kafka-0 -- /opt/bitnami/kafka/bin/kafka-console-consumer.sh \
            --bootstrap-server kafka:9092 \
            --topic <your topic name> \
            --from-beginning
```

- list topic

```bash
kubectl exec --tty -i kafka-0 -- /opt/bitnami/kafka/bin/kafka-topics.sh --list --bootstrap-server kafka:9092
```