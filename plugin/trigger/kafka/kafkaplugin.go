package kafka

import (
	"EventFlow/common/interface/pluginbase"
	"EventFlow/common/tool/stringtool"
	"context"

	kafka "github.com/segmentio/kafka-go"
)

type KafkaPlugin struct {
	reader *kafka.Reader
	pluginbase.ActionHandler
	Setting SettingConfig
}

func (trigger *KafkaPlugin) Start() {

	trigger.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  stringtool.ToStringArray(trigger.Setting.Brokers),
		GroupID:  trigger.Setting.GroupID,
		Topic:    trigger.Setting.Topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {

		defer cancel()

		for {
			message, err := trigger.reader.FetchMessage(ctx)

			if err == nil {
				var triggerPlugin pluginbase.ITriggerPlugin = trigger
				value := string(message.Value)

				trigger.FireAction(&triggerPlugin, &value)

				trigger.reader.CommitMessages(ctx, message)
			} else {
				if err == context.Canceled {
					break
				}
			}
		}
	}()
}

func (trigger *KafkaPlugin) Stop() {
	trigger.reader.Close()
}
