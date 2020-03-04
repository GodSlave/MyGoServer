package test

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"fmt"
	"time"
	"github.com/GodSlave/MyGoServer/log"
	"sync"
	"os"
	"os/signal"
	"testing"
)

func TestSendAndReceiver(test *testing.T)  {
	topic := []string{"test"}
	Address :=[]string{"0.0.0.0:2181"}

	AsyncProducer("test")

	var wg = &sync.WaitGroup{}
	wg.Add(2)
	//广播式消费：消费者1
	go clusterConsumer(wg, Address, topic, "group-1")
	//广播式消费：消费者2
	go clusterConsumer(wg, Address, topic, "group-2")




}

func AsyncProducer(topic string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner(topic)
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V2_1_0_0
	fmt.Println("start make producer")

	producer, e := sarama.NewAsyncProducer([]string{"0.0.0.0:32770"}, config)
	if e != nil {
		log.Error(e.Error())
		return
	}
	defer producer.AsyncClose()

	log.Info("Start goroutine")
	go func(p sarama.AsyncProducer) {

		for {
			select {
			case <-p.Successes():
			case fail := <-p.Errors():
				log.Error(fail.Error())
			}
		}

	}(producer)

	var value string
	for i := 0; ; i++ {
		time.Sleep(500 * time.Microsecond)
		time11 := time.Now()
		value = "this is a message " + time11.Format("")
		msg := &sarama.ProducerMessage{
			Topic: "test",
		}
		msg.Value = sarama.StringEncoder(value)
		producer.Input() <- msg
	}

}

func clusterConsumer(wg *sync.WaitGroup,brokers []string,topics []string,groupId string)  {
	defer wg.Done()
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer,err := cluster.NewConsumer(brokers,groupId,topics,config)
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer consumer.Close()

	signals := make(chan os.Signal,1)
	signal.Notify(signals,os.Interrupt)

	go func() {

		for ntf:= range consumer.Notifications() {
			log.Info("%v",ntf)
		}
	}()

	
	var successes int 
	Loop:
		for   {
			select {
				case msg,ok := <-consumer.Messages():
					if ok {
						log.Info("%s:%s/%d/%d\t%s\t%s\n", groupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
						consumer.MarkOffset(msg,"")
						successes++
					}
			case <-signals:
				break Loop
			}
		}
	fmt.Fprintf(os.Stdout, "%s consume %d messages \n", groupId, successes)

}
