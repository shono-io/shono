package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/benthosdev/benthos/v4/public/bloblang"
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/shono-io/shono/commons"
	"github.com/shono-io/shono/graph"
	"github.com/shono-io/shono/systems"
	"github.com/sirupsen/logrus"
)

func Register(env graph.Environment) {
	err := service.RegisterProcessor("store", storeProcConfig(), func(conf *service.ParsedConfig, mgr *service.Resources) (service.Processor, error) {
		return procFromConfig(env, conf)
	})
	if err != nil {
		panic(err)
	}
}

func storeProcConfig() *service.ConfigSpec {
	spec := service.NewConfigSpec().
		Beta().
		Categories("Integration")

	return spec.
		Field(service.NewStringField("store_key").
			Description("The key of the store to use")).
		Field(service.NewStringField("operation").
			Description("The operation to perform, one of: 'list', 'get', 'add', 'set' or 'delete'")).
		Field(service.NewInterpolatedStringField("key").
			Description("The key to use. This is only applicable for 'get', 'add', 'set' and 'delete'").
			Optional()).
		Field(service.NewBloblangField("value").
			Description("The value to use. This is only applicable for 'add' and 'set'").
			Optional()).
		Field(service.NewInterpolatedStringMapField("filters").
			Description("A map of filters to apply to the operation. This is only applicable for 'list'").
			Optional())
}

func procFromConfig(env graph.Environment, conf *service.ParsedConfig) (*storeProc, error) {
	proc := &storeProc{}

	ssk, err := conf.FieldString("store_key")
	if err != nil {
		return nil, fmt.Errorf("failed to get the store key: %w", err)
	}

	sk, err := commons.ParseString(ssk)
	if err != nil {
		return nil, fmt.Errorf("invalid store key: %w", err)
	}

	proc.st, err = env.GetStore(sk)
	if err != nil {
		return nil, fmt.Errorf("no store with key %q: %w", sk.String(), err)
	}

	storage, err := env.GetStorage(proc.st.StorageKey())
	if err != nil {
		return nil, fmt.Errorf("no storage with key %q: %w", proc.st.StorageKey(), err)
	}

	// -- once we have the storage, we can get the storage system
	ss, fnd := systems.Storage[storage.Key().String()]
	if !fnd {
		return nil, fmt.Errorf("no storage system with key %q", storage.Key().String())
	}

	// -- get a client for the storage system
	proc.cl, err = ss.GetClient(storage.Config())
	if err != nil {
		return nil, fmt.Errorf("failed to get client for storage system %q: %w", storage.Key().String(), err)
	}

	proc.operation, err = conf.FieldString("operation")
	if err != nil {
		return nil, fmt.Errorf("failed to get operation: %w", err)
	}

	if conf.Contains("key") {
		proc.key, err = conf.FieldInterpolatedString("key")
		if err != nil {
			return nil, fmt.Errorf("failed to get key: %w", err)
		}
	}

	if conf.Contains("filters") {
		proc.filters, err = conf.FieldInterpolatedStringMap("filters")
		if err != nil {
			return nil, fmt.Errorf("failed to get filters: %w", err)
		}
	}

	if conf.Contains("value") {
		proc.value, err = conf.FieldBloblang("value")
		if err != nil {
			return nil, fmt.Errorf("failed to get value: %w", err)
		}
	}

	return proc, nil
}

type storeProc struct {
	cl graph.StorageClient
	st *graph.Store

	operation string
	key       *service.InterpolatedString
	value     *bloblang.Executor
	filters   map[string]*service.InterpolatedString
}

func (s *storeProc) Process(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	switch s.operation {
	case "get":
		return s.processGet(ctx, message)
	case "add":
		return s.processAdd(ctx, message)
	case "set":
		return s.processReplace(ctx, message)
	case "delete":
		return s.processDelete(ctx, message)
	case "list":
		return s.processList(ctx, message)
	default:
		return nil, fmt.Errorf("unknown operation: %s", s.operation)
	}
}

func (s *storeProc) Close(ctx context.Context) error {
	return s.cl.Close()
}

func (s *storeProc) processGet(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	// -- get the key from the message
	key, err := s.key.TryString(message)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	res, err := s.cl.Get(ctx, s.st.Collection(), key)
	if err != nil {
		return nil, fmt.Errorf("unable to read document with key %q: %w", key, err)
	}

	result := service.NewMessage(nil)
	result.SetStructured(res)

	CopyMeta(message, result)

	return service.MessageBatch{result}, nil

}

func (s *storeProc) processAdd(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	// -- get the key from the message
	key, err := s.key.TryString(message)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	data, err := s.getMessagePayload(message)
	if err != nil {
		return nil, err
	}

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		b, _ := json.Marshal(data)
		logrus.Tracef("adding document %q as %s", key, b)
	}

	if err := s.cl.Add(ctx, s.st.Collection(), key, data); err != nil {
		return nil, fmt.Errorf("unable to add document with key %s: %w", key, err)
	}

	result := service.NewMessage(nil)
	result.SetStructured(data)

	CopyMeta(message, result)

	return service.MessageBatch{result}, nil
}

func (s *storeProc) processReplace(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	// -- get the key from the message
	key, err := s.key.TryString(message)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	data, err := s.getMessagePayload(message)
	if err != nil {
		return nil, err
	}

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		b, _ := json.Marshal(data)
		logrus.Tracef("setting document %q to %s", key, b)
	}

	if err := s.cl.Set(ctx, s.st.Collection(), key, data); err != nil {
		return nil, fmt.Errorf("unable to set document with key %s: %w", key, err)
	}

	result := service.NewMessage(nil)
	result.SetStructured(data)

	CopyMeta(message, result)

	return service.MessageBatch{result}, nil
}

func (s *storeProc) processDelete(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	// -- get the key from the message
	key, err := s.key.TryString(message)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		logrus.Tracef("removing document %q", key)
	}

	// -- get the document so we can return it
	data, err := s.cl.Get(ctx, s.st.Collection(), key)
	if err != nil {
		return nil, fmt.Errorf("unable to read document with key %q: %w", key, err)
	}

	if err := s.cl.Delete(ctx, s.st.Collection(), key); err != nil {
		return nil, fmt.Errorf("unable to delete document with key %s: %w", key, err)
	}

	result := service.NewMessage(nil)
	result.SetStructured(data)

	CopyMeta(message, result)

	return service.MessageBatch{result}, nil
}

func (s *storeProc) processList(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	// -- construct the filters
	filters := make(map[string]interface{}, len(s.filters))
	for k, v := range s.filters {
		value, err := v.TryString(message)
		if err != nil {
			return nil, fmt.Errorf("failed to interpolate filter: %w", err)
		}

		filters[k] = value
	}

	cur, err := s.cl.List(ctx, s.st.Collection(), filters, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}

	var res []*service.Message
	for cur.HasNext() {
		doc, err := cur.Read()
		if err != nil {
			return nil, fmt.Errorf("failed to read document: %w", err)
		}

		result := service.NewMessage(nil)
		result.SetStructured(doc)

		CopyMeta(message, result)

		res = append(res, result)
	}

	return res, nil
}

func (s *storeProc) getMessagePayload(message *service.Message) (map[string]any, error) {
	sd, err := s.value.Query(message)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the value: %w", err)
	}

	switch data := sd.(type) {
	case map[string]any:
		return data, nil
	case *service.Message:
		m, err := data.AsStructuredMut()
		if err != nil {
			return nil, fmt.Errorf("failed to get the value from the message: %w", err)
		}

		switch dt := m.(type) {
		case map[string]any:
			return dt, nil
		default:
			return nil, fmt.Errorf("unsupported mapped message payload type: %T", sd)
		}

	default:
		return nil, fmt.Errorf("unsupported message payload type: %T", sd)
	}
}

func CopyMeta(src, dst *service.Message) {
	_ = src.MetaWalk(func(k string, v string) error {
		dst.MetaSet(k, v)
		return nil
	})
}
