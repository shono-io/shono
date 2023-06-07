package backbone

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/shono-io/shono/commons"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl"
	kaws "github.com/twmb/franz-go/pkg/sasl/aws"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

func (k *kafkaBackbone) GetClient() (Client, error) {
	var opts []kgo.Opt
	opts = append(opts, kgo.SeedBrokers(k.config.BootstrapServers...))

	if k.config.SASL != nil {
		var ms []sasl.Mechanism
		for _, s := range k.config.SASL {
			sm, err := asSASLMechanism(s)
			if err != nil {
				return nil, err
			}
			ms = append(ms, sm)
		}
		opts = append(opts, kgo.SASL(ms...))
	}

	panic("implement me")
}

func asSASLMechanism(config SASLConfig) (sasl.Mechanism, error) {
	switch config.Mechanism {
	case "PLAIN":
		return plain.Auth{User: config.Username, Pass: config.Password}.AsMechanism(), nil
	case "AWS_MSK_IAM":
		awsSession, err := getAwsSession(*config.AwsConfig)
		if err != nil {
			return nil, err
		}

		creds := awsSession.Config.Credentials
		return kaws.ManagedStreamingIAM(func(ctx context.Context) (kaws.Auth, error) {
			val, err := creds.GetWithContext(ctx)
			if err != nil {
				return kaws.Auth{}, err
			}
			return kaws.Auth{
				AccessKey:    val.AccessKeyID,
				SecretKey:    val.SecretAccessKey,
				SessionToken: val.SessionToken,
			}, nil
		}), nil
	default:
		return nil, fmt.Errorf("unsupported SASL mechanism: %s", config.Mechanism)
	}
}

// GetSession attempts to create an AWS session based on the parsedConfig.
func getAwsSession(cfg AWSConfig, opts ...func(*aws.Config)) (*session.Session, error) {
	awsConf := aws.NewConfig()

	if cfg.Region != "" {
		awsConf = awsConf.WithRegion(cfg.Region)
	}
	if cfg.Endpoint != "" {
		awsConf = awsConf.WithEndpoint(cfg.Endpoint)
	}
	if cfg.Credentials.Profile != "" {
		awsConf = awsConf.WithCredentials(credentials.NewSharedCredentials(
			"", cfg.Credentials.Profile,
		))
	} else if cfg.Credentials.Id != "" {
		awsConf = awsConf.WithCredentials(credentials.NewStaticCredentials(
			cfg.Credentials.Id, cfg.Credentials.Secret, cfg.Credentials.Token,
		))
	}

	for _, opt := range opts {
		opt(awsConf)
	}

	sess, err := session.NewSession(awsConf)
	if err != nil {
		return nil, err
	}

	if role := cfg.Credentials.Role; role != "" {
		var opts []func(*stscreds.AssumeRoleProvider)
		if externalID := cfg.Credentials.RoleExternalId; externalID != "" {
			opts = []func(*stscreds.AssumeRoleProvider){
				func(p *stscreds.AssumeRoleProvider) {
					p.ExternalID = &externalID
				},
			}
		}
		sess.Config = sess.Config.WithCredentials(
			stscreds.NewCredentials(sess, role, opts...),
		)
	}

	if cfg.Credentials.FromEC2Role {
		sess.Config = sess.Config.WithCredentials(ec2rolecreds.NewCredentials(sess))
	}

	return sess, nil
}

type kafkaClient struct {
	kc *kgo.Client
}

func (k *kafkaClient) Produce(ctx context.Context, event commons.Key, key commons.Key, payload map[string]any) (err error) {
	var b []byte = nil

	if payload == nil {
		b, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}

	record := &kgo.Record{
		Key:   []byte(key.String()),
		Value: b,
		Headers: []kgo.RecordHeader{
			{Key: "io_shono_kind", Value: []byte(event.String())},
		},
	}

	pr := k.kc.ProduceSync(ctx, record)
	if pr.FirstErr() != nil {
		return pr.FirstErr()
	}

	return nil
}
