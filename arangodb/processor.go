package arangodb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/benthosdev/benthos/v4/public/bloblang"
	"github.com/benthosdev/benthos/v4/public/service"
	"github.com/sirupsen/logrus"
	"strings"
)

func arangodbProcConfig() *service.ConfigSpec {
	spec := service.NewConfigSpec().
		Beta().
		Categories("Integration")

	for _, f := range clientFields() {
		spec = spec.Field(f)
	}

	return spec.
		Field(service.NewStringField("database").
			Default("_system").
			Description("The database to use when connecting to the ArangoDB server")).
		Field(service.NewStringField("operation").
			Description("The operation to perform, one of: 'list', 'get', 'add', 'set' or 'delete'")).
		Field(service.NewStringField("collection").
			Description("The collection to use.").
			Optional()).
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

func procFromConfig(conf *service.ParsedConfig) (*arangodbProc, error) {
	adb, err := getClientFromConfig(conf)
	if err != nil {
		return nil, err
	}

	proc := &arangodbProc{}

	proc.operation, err = conf.FieldString("operation")
	if err != nil {
		return nil, fmt.Errorf("failed to get operation: %w", err)
	}

	dn, err := conf.FieldString("database")
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	} else {
		proc.db, err = adb.Database(context.Background(), dn)
		if err != nil {
			return nil, fmt.Errorf("failed to get arangodb database: %w", err)
		}
	}

	if conf.Contains("collection") {
		cn, err := conf.FieldString("collection")
		if err != nil {
			return nil, fmt.Errorf("failed to get collection: %w", err)
		}

		proc.col, err = proc.db.Collection(context.Background(), cn)
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

type arangodbProc struct {
	db  driver.Database
	col driver.Collection

	operation string
	key       *service.InterpolatedString
	value     *bloblang.Executor
	filters   map[string]*service.InterpolatedString
}

func (a arangodbProc) Process(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	switch a.operation {
	case "get":
		return a.processGet(ctx, message)
	case "add":
		return a.processAdd(ctx, message)
	case "set":
		return a.processReplace(ctx, message)
	case "delete":
		return a.processDelete(ctx, message)
	case "list":
		return a.processList(ctx, message)
	default:
		return nil, fmt.Errorf("unknown operation: %s", a.operation)
	}
}

func (a arangodbProc) processGet(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	// -- get the key from the message
	key, err := a.key.TryString(message)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	var doc map[string]any
	dm, err := a.col.ReadDocument(ctx, key, &doc)
	if err != nil {
		return nil, fmt.Errorf("failed to read document: %w", err)
	}

	result := service.NewMessage(nil)
	result.MetaSet("adb_key", dm.Key)
	result.MetaSet("adb_revision", dm.Rev)
	result.SetStructured(doc)

	return service.MessageBatch{result}, nil
}

func (a arangodbProc) processAdd(ctx context.Context, message *service.Message) (service.MessageBatch, error) {

	data, err := a.getMessagePayload(message)
	if err != nil {
		return nil, err
	}

	if logrus.IsLevelEnabled(logrus.TraceLevel) {
		b, _ := json.Marshal(data)
		logrus.Tracef("storing document %s", b)
	}

	result := service.NewMessage(nil)

	dm, err := a.col.CreateDocument(ctx, data)
	if err != nil {
		if driver.IsConflict(err) {
			err2 := fmt.Errorf("a document with key %q already exists", data["_key"])
			result.SetError(err2)
			return nil, err2
		}

		return nil, fmt.Errorf("failed to add document: %w", err)
	}

	_ = message.MetaWalk(func(k string, v string) error {
		result.MetaSet(k, v)
		return nil
	})

	result.MetaSet("adb_key", dm.Key)
	result.MetaSet("adb_revision", dm.Rev)
	result.SetStructured(data)

	return service.MessageBatch{result}, nil
}

func (a arangodbProc) processReplace(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	// -- get the key from the message
	key, err := a.key.TryString(message)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	result := message.Copy()

	data, err := a.getMessagePayload(message)
	if err != nil {
		return nil, err
	}

	dm, err := a.col.ReplaceDocument(ctx, key, data)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			result.SetError(err)
			return nil, nil
		}

		return nil, fmt.Errorf("failed to replace document: %w", err)
	}

	result.MetaSet("adb_key", dm.Key)
	result.MetaSet("adb_revision", dm.Rev)
	result.MetaSet("adb_old_revision", dm.OldRev)

	return service.MessageBatch{result}, nil
}

func (a arangodbProc) processDelete(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	// -- get the key from the message
	key, err := a.key.TryString(message)
	if err != nil {
		return nil, fmt.Errorf("failed to get key: %w", err)
	}

	result := message.Copy()

	var oldData map[string]any
	ctx = driver.WithReturnOld(ctx, &oldData)

	dm, err := a.col.RemoveDocument(ctx, key)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			result.SetError(err)
			return nil, nil
		}

		return nil, fmt.Errorf("failed to delete document: %w", err)
	}

	result.MetaSet("adb_key", dm.Key)
	result.MetaSet("adb_revision", dm.Rev)
	result.SetStructured(oldData)

	return service.MessageBatch{result}, nil
}

func (a arangodbProc) processList(ctx context.Context, message *service.Message) (service.MessageBatch, error) {
	// -- construct the filters
	bindVars := make(map[string]interface{}, len(a.filters))
	filters := make([]string, 0, len(a.filters))
	for k, v := range a.filters {
		value, err := v.TryString(message)
		if err != nil {
			return nil, fmt.Errorf("failed to interpolate filter: %w", err)
		}

		param := strings.ReplaceAll(k, ".", "_")
		filters = append(filters, fmt.Sprintf("%s == @%s", k, param))
		bindVars[param] = value
	}

	filterString := ""
	if len(filters) > 0 {
		filterString = "FILTER " + strings.Join(filters, " AND ")
	}

	cur, err := a.db.Query(ctx, fmt.Sprintf("FOR d IN %s %s RETURN d", a.col.Name(), filterString), bindVars)
	if err != nil {
		return nil, fmt.Errorf("failed to query collection: %w", err)
	}

	res := make(service.MessageBatch, cur.Count())
	current := 0
	for cur.HasMore() {
		var doc map[string]any
		dm, err := cur.ReadDocument(ctx, &doc)
		if err != nil {
			return nil, fmt.Errorf("failed to read document: %w", err)
		}

		result := service.NewMessage(nil)
		result.MetaSet("adb_key", dm.Key)
		result.MetaSet("adb_revision", dm.Rev)
		result.SetStructured(doc)
		res[current] = result

		current++
	}

	return res, nil
}

func (a arangodbProc) getMessagePayload(message *service.Message) (map[string]any, error) {
	sd, err := a.value.Query(message)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve the value: %w", err)
	}

	switch data := sd.(type) {
	case map[string]any:
		// -- add the key to the payload
		key, err := a.key.TryString(message)
		if err != nil {
			return nil, fmt.Errorf("failed to get the key from the message: %w", err)
		}
		data["_key"] = key

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

func (a arangodbProc) Close(ctx context.Context) error {
	return nil
}
