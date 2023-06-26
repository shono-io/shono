package kafka

type Opt func(config map[string]any)

func WithSeedBroker(seedBroker string) Opt {
	return func(config map[string]any) {
		sb, fnd := config["seed_brokers"]
		if !fnd {
			sb = []string{}
		}
		sb = append(sb.([]string), seedBroker)
		config["seed_brokers"] = sb
	}
}

func WithConsumerGroup(consumerGroup string) Opt {
	return func(config map[string]any) {
		config["consumer_group"] = consumerGroup
	}
}

func WithInputTopics(topics ...string) Opt {
	return func(config map[string]any) {
		config["topics"] = topics
	}
}

func WithOutputTopic(topic string) Opt {
	return func(config map[string]any) {
		config["topic"] = topic
	}
}
