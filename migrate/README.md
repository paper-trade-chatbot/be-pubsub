# pubsub scripts

## create your script

1. please copy the `template.sh` and rename to `<yyyymmdd>_<your topic name>.sh`
2. sudo chmod +x `<yyyymmdd>_<your topic name>.sh`
3. change `TOPIC_NAME` to your topic name

## run the script

1. make sure you are under the right environment
2. your may check existing topics using command below:
```bash
kubectl exec --tty -i kafka-0 -- /opt/bitnami/kafka/bin/kafka-topics.sh --list --bootstrap-server kafka:9092
```
3. run the script

## delete topic

- message won't be deleted even you delete topic

```bash
kubectl exec --tty -i kafka-0 -- /opt/bitnami/kafka/bin/kafka-topics.sh --bootstrap-server kafka:9092 --delete --topic <your topic name>
```

## how to know topic name?

- put this in your code and run it to print the topic name

```go
temp, err := xxx.NewSubscriber()

fmt.Println(temp.GetTopicName())
```