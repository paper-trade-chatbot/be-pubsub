package kafka

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/segmentio/kafka-go"
)

type PubsubImpl struct {
	Config Config
	Pubsub
	LogMode bool
}

func (p *PubsubImpl) GetTopicName() string {
	name := []string{}

	name = append(name, string(p.Config.GetExported()))

	name = append(name, p.Config.GetDomain()...)
	if len(p.Config.GetDomain()) == 0 {
		name = append(name, "default")
	}

	name = append(name, p.Config.GetDataset())

	name = append(name, string(p.Config.GetDataType()))

	name = append(name, strconv.FormatUint(uint64(p.Config.GetVersion()), 10))

	topicName := strings.Join(name, ".")

	return topicName
}

func (p *PubsubImpl) GetMutex() *sync.Mutex {
	return p.Pubsub.GetMutex()
}

func (p *PubsubImpl) ListTopics() []string {

	topicSet := mapset.NewSet[string]()
	for _, b := range p.Config.GetBrokers() {
		func() {
			conn, err := kafka.Dial("tcp", b)
			if err != nil {
				panic(err.Error())
			}
			defer conn.Close()

			partitions, err := conn.ReadPartitions()
			if err != nil {
				panic(err.Error())
			}

			for _, p := range partitions {
				topicSet.Add(p.Topic)
			}
		}()
	}

	return topicSet.ToSlice()
}

func (p *PubsubImpl) SetLogMode(on bool) {
	p.LogMode = on
}

func (p *PubsubImpl) Log(format string, a ...any) {
	if p.LogMode {
		fmt.Printf("kafka: "+format, a)
		fmt.Println()
	}
}
