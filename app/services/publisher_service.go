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

	// Get signal for finish
	doneCh := make(chan struct{})

	go func() {

		// BackTrack data 3 months back, subtract about 90 days from now:
		threeMonthsAgo := time.Now().Unix() - 7776000
		backTrackData(&q, ch, threeMonthsAgo, originIDs, activityData)

		for {

			now := time.Now()
			currentInterval := getCurrentInterval(now)

			fmt.Println("Current Interval: ", currentInterval)

			for _, v := range originIDs {
				unixtime := now.Unix()
				publishMessageToQueue(&q, ch, v, currentInterval, unixtime, activityData)
			}

			time.Sleep(5 * time.Minute)
		}
	}()

	<-doneCh
}

func backTrackData(q *amqp.Queue, ch *amqp.Channel, unixtime int64, originIDs []int, activityData [5][288]models.FiveMinuteActivity) {

	timenow := time.Now().Unix()

	for unixtime < timenow {

		time := time.Unix(int64(unixtime), 0)
		currentInterval := getCurrentInterval(time)

		for _, v := range originIDs {
			publishMessageToQueue(q, ch, v, currentInterval, unixtime, activityData)
		}

		unixtime += 300 //add 5 minutes
	}
}

func publishMessageToQueue(q *amqp.Queue, ch *amqp.Channel, origin int, interval int, unixtime int64, activityData [5][288]models.FiveMinuteActivity) {

	activityStreamID := activityStream[origin]
	intervalData := activityData[activityStreamID][interval]

	intervalData.OriginID = uint64(origin)
	intervalData.Epochtime = unixtime

	fmt.Printf("originID: %d activity: %s\n", origin, intervalData.Description)
	newActivityMsg, err := json.Marshal(intervalData)
	failOnError(err, "Unable to marshall activity message")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        newActivityMsg,
		})

	failOnError(err, "Failed to publish a message")
}

func getCurrentInterval(now time.Time) int {

	currentHour := now.Hour()
	currentMinute := now.Minute()
	currentInterval := (currentHour * 12) + (currentMinute / 5)

	return currentInterval
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
