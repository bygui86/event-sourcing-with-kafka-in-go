package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	"github.com/backlin/resources/event-sourcing-in-go/examples/orchestrator"
	"github.com/bsm/sarama-cluster"
	"github.com/go-kit/kit/log"
)

type Worker struct {
	Logger       log.Logger
	Consumer     *cluster.Consumer
	Producer     sarama.SyncProducer
	TrackerTopic string

	offset int64
}

func (w Worker) Run() error {
	for {

		inMsg, open := <-w.Consumer.Messages()

		if !open {
			w.Logger.Log("consumer closed, shutting down")
			return nil
		}

		w.offset = inMsg.Offset

		id, work, err := orchestrator.UnmarshalWork(inMsg)
		if err != nil {
			w.Logger.Log("offset", inMsg.Offset, "msg", "invalid message", "err", err)
			continue
		}

		w.Logger.Log("offset", inMsg.Offset, "batch_id", id, "job_id", work.JobId, "status", orchestrator.Event_RUNNING)
		status := processWork(work)

		w.sendResponse(work, status)
		w.Consumer.MarkOffset(inMsg, "")
		w.Logger.Log("offset", inMsg.Offset, "batch_id", id, "job_id", work.JobId, "status", status)

	}
}

func processWork(work orchestrator.Work) orchestrator.Event_Status {
	// Do the actual work
	time.Sleep(time.Duration(work.Duration) * time.Millisecond)

	if rand.Float32() < work.FailureRate {
		return orchestrator.Event_FAILURE
	}

	return orchestrator.Event_SUCCESS
}

func (w *Worker) sendResponse(work orchestrator.Work, status orchestrator.Event_Status) error {
	// Package the output
	value := orchestrator.Event{
		BatchId:     work.BatchId,
		JobId:       work.JobId,
		StatusLevel: orchestrator.Event_JOB,
		Status:      status,
	}

	b, err := value.Marshal()
	if err != nil {
		return fmt.Errorf("could not marshal %s job event: %s", value.Status, err)
		// The application is broken, always fail hard
	}

	outMsg := &sarama.ProducerMessage{
		Topic: w.TrackerTopic,
		Value: sarama.ByteEncoder(b),
	}

	// Send with simple retry logic
	for {
		_, _, err := w.Producer.SendMessage(outMsg)
		if err == nil {
			break
		}

		w.Logger.Log("offset", w.offset, "err", err, "msg", "failed sending event message")
		time.Sleep(5 * time.Second)
	}

	return nil
}
