package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/ahugofreire/go-b3/internal/infra/kafka"
	"github.com/ahugofreire/go-b3/internal/market/dto"
	"github.com/ahugofreire/go-b3/internal/market/entity"
	"github.com/ahugofreire/go-b3/internal/market/transformer"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

/*
	{"order_id":"1","investor_id":"john","asset_id":"B3SA3","current_shares":10,"shares":5,"price":12,"order_type":"SELL"}
	{"investor_id":"Ana","asset_id":"B3SA3","current_shares":10,"shares":0,"price":12,"order_type":"BUY"}
*/

func main() {
	ordersIn := make(chan *entity.Order)
	ordersOut := make(chan *entity.Order)
	waitGroup := &sync.WaitGroup{}
	defer waitGroup.Wait()

	kafkaMsgChan := make(chan *ckafka.Message)
	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "orders-consumer-id",
		"auto.offset.reset": "latest",
	}

	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewConsumer(configMap, []string{"input"})

	go kafka.Consume(kafkaMsgChan) //Thread 2

	// recebe do canal do kafka, joga no input, processa e joga no output e depois publica no kafka
	book := entity.NewBook(ordersIn, ordersOut, waitGroup)
	go book.Trade() //Thread 3

	go func() {
		for msg := range kafkaMsgChan {
			log.Println(string(msg.Value))

			waitGroup.Add(1)
			tradeInput := dto.TradeInput{}
			err := json.Unmarshal(msg.Value, &tradeInput)
			if err != nil {
				panic(err)
			}
			order := transformer.TransformInput(tradeInput)
			ordersIn <- order
		}
	}()

	for res := range ordersOut {
		output := transformer.TransformOutput(res)
		outputJson, err := json.MarshalIndent(output, "", " ")
		fmt.Println(string(outputJson))
		if err != nil {
			fmt.Println(err)
		}
		producer.Publish(string(outputJson), []byte("orders"), "output")
	}
}
