package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

var (
	log = logrus.New()

	cfgFile = flag.String("cfg", "copycat.json", "Configuration file to load from")
)

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	data, err := os.ReadFile(*cfgFile)
	if err != nil {
		log.Fatalf("failed to read config from '%s': %v", *cfgFile, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	} else if err := cfg.validate(); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.TargetAddr,
	})
	if err != nil {
		log.Fatalf("failed to setup producer: %v", err)
	}

	log.Infof(
		"Will be copying from topic '%s' (in %s) to topic '%s' (in %s)",
		cfg.SourceTopic, cfg.SourceAddr, cfg.TargetTopic, cfg.TargetAddr)
	st := &Streamer{
		Topic:         cfg.SourceTopic,
		Logger:        log.WithField("component", "streamer"),
		Servers:       cfg.SourceAddr,
		Workers:       cfg.Workers,
		StartOffset:   cfg.StartOffset,
		ConsumerGroup: cfg.GroupID,
		Apply: func(ctx context.Context, key, val []byte) error {

			err := kafkaPublishBlocking(ctx, p, kafka.Message{
				Key:   key,
				Value: val,
				TopicPartition: kafka.TopicPartition{
					Topic:     &cfg.TargetTopic,
					Partition: kafka.PartitionAny,
				},
			})
			if err != nil {
				log.Warnf("failed to publish message copy: %v", err)
			}
			log.Infof("published message")
			return nil
		},
	}

	if err := st.Run(ctx); err != nil {
		log.Fatalf("streamer exited with error: %v", err)
	}
	log.Infof("streamer exited successfully")
}

func kafkaPublishBlocking(ctx context.Context, p *kafka.Producer, msg kafka.Message) error {
	deliveryChan := make(chan kafka.Event)
	if err := p.Produce(&msg, deliveryChan); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return nil

	case e := <-deliveryChan:
		if m, ok := e.(*kafka.Message); ok {
			return m.TopicPartition.Error
		}
	}
	return nil
}

type Config struct {
	Workers     int    `json:"workers"`
	GroupID     string `json:"group_id"`
	StartOffset string `json:"start_offset"`
	SourceAddr  string `json:"source_addr"`
	TargetAddr  string `json:"target_addr"`
	SourceTopic string `json:"source_topic"`
	TargetTopic string `json:"target_topic"`
}

func (c *Config) validate() error {
	if c.Workers == 0 {
		c.Workers = 1
	}
	if c.GroupID == "" {
		c.GroupID = "copycat"
	}
	if c.SourceAddr == c.TargetAddr && c.SourceTopic == c.TargetTopic {
		return errors.New("source and destination are same, nothing to do")
	}
	return nil
}
