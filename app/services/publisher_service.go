package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/garcialuis/ActivityGenerator/app/models"
	"github.com/streadway/amqp"
)

var activityStream map[int]int

func PublishActivities(activityData [5][288]models.FiveMinuteActivity, originIDs []int) {
	// map originIDs to activityData
	originToActivityMapper(originIDs)
	//
	rabbitMQProducer(activityData, originIDs)
}

func rabbitMQProducer(activityData [5][288]models.FiveMinuteActivity, originIDs []int) {

	fmt.Println("Starting RabbitMQ Producer")
	time.Sleep(7 * time.Second)

	conn, err := amqp.Dial(brokerAddr())
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queue(), // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgCount := 0

	// Get signal for finish
	doneCh := make(chan struct{})

	go func() {

		for {

			now := time.Now()

			currentHour := now.Hour()
			currentMinute := now.Minute()
			currentInterval := (currentHour * 12) + (currentMinute / 5)

			fmt.Println("Current Interval: ", currentInterval)

			for _, v := range originIDs {

				activityStreamID := activityStream[v]
				intervalData := activityData[activityStreamID][currentInterval]

				intervalData.OriginID = uint64(v)
				intervalData.Epochtime = now.Unix()

				fmt.Printf("originID: %d activity: %s\n", v, intervalData.Description)

				newActivityMsg, err := json.Marshal(intervalData)
				failOnError(err, "Unable to marshall activity message")

				msgCount++
				body := fmt.Sprintf("Hello RabbitMQ message %v", msgCount)

				err = ch.Publish(
					"",     // exchange
					q.Name, // routing key
					false,  // mandatory
					false,  // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        newActivityMsg,
					})
				log.Printf(" [x] Sent %s", body)
				failOnError(err, "Failed to publish a message")

			}

			time.Sleep(5 * time.Minute)
		}
	}()

	<-doneCh
}

func originToActivityMapper(originIDs []int) {
	activityStream = make(map[int]int)

	for _, v := range originIDs {
		randDataStream := models.RandomNumInRage(0, 4)
		activityStream[v] = randDataStream
	}
}

func brokerAddr() string {
	brokerAddr := os.Getenv("BROKER_ADDR")
	if len(brokerAddr) == 0 {
		brokerAddr = "amqp://guest:guest@localhost:5672/"
	}
	return brokerAddr
}

func queue() string {
	queue := os.Getenv("QUEUE")
	if len(queue) == 0 {
		queue = "default-queue"
	}
	return queue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
