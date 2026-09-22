package kafka

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/mgreben/glink/internal/config"
	"github.com/mgreben/glink/internal/modules/links"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/fx"
)

type ClickProducer struct {
	client *kgo.Client
	topic  string
}

func NewClickProducer(lc fx.Lifecycle, cfg *config.Config) (*ClickProducer, error) {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Kafka.Brokers...),
		kgo.DefaultProduceTopic(cfg.Kafka.ClickTopic),
		kgo.MaxBufferedRecords(cfg.Kafka.MaxBufferedRecords),
	)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{OnStop: func(context.Context) error { client.Close(); return nil }})
	return &ClickProducer{client: client, topic: cfg.Kafka.ClickTopic}, nil
}

func (p *ClickProducer) Publish(_ context.Context, event links.ClickEvent) {
	value, err := json.Marshal(event)
	if err != nil {
		return
	}
	p.client.TryProduce(context.Background(), &kgo.Record{
		Topic: p.topic,
		Key:   []byte(strconv.FormatInt(event.LinkID, 10)),
		Value: value,
	}, func(*kgo.Record, error) {})
}

type ClickStatsSink interface {
	InsertClickEvents(context.Context, []links.ClickEvent) error
}

type ClickConsumer struct {
	cfg *config.Config
}

func NewClickConsumer(cfg *config.Config) *ClickConsumer {
	return &ClickConsumer{cfg: cfg}
}

func (c *ClickConsumer) Run(ctx context.Context, sink ClickStatsSink) error {
	client, err := kgo.NewClient(
		kgo.SeedBrokers(c.cfg.Kafka.Brokers...),
		kgo.ConsumerGroup(c.cfg.Kafka.ClickConsumerGroup),
		kgo.ConsumeTopics(c.cfg.Kafka.ClickTopic),
		kgo.BlockRebalanceOnPoll(),
	)
	if err != nil {
		return err
	}
	defer client.Close()

	for {
		fetches := client.PollFetches(ctx)
		if ctx.Err() != nil || fetches.IsClientClosed() {
			return nil
		}
		if fetches.Err() != nil {
			continue
		}

		records := fetches.Records()
		batch := clickEvents(records)
		if err := sink.InsertClickEvents(ctx, batch); err != nil {
			client.AllowRebalance()
			continue
		}
		if len(records) > 0 {
			_ = client.CommitRecords(ctx, records...)
		}
		client.AllowRebalance()
	}
}

func clickEvents(records []*kgo.Record) []links.ClickEvent {
	events := make([]links.ClickEvent, 0, len(records))
	for _, record := range records {
		var event links.ClickEvent
		if json.Unmarshal(record.Value, &event) != nil || event.LinkID == 0 {
			continue
		}
		events = append(events, event)
	}

	return events
}
