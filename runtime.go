package go_shono

import (
	"github.com/twmb/franz-go/pkg/kgo"
)

type Runtime struct {
	kc        *kgo.Client
	resources map[string]Resource[any]
}

func (r *Runtime) Resource(id string) any {
	res, fnd := r.resources[id]
	if !fnd {
		return nil
	}

	return res.Client()
}

//func (r *Runtime) Failed(ctx context.Context, evt EventMeta, err error) {
//	r.Send(ctx, evt, events.OperationFailed{Reason: err.Error()})
//}
//
//func (r *Runtime) Send(ctx context.Context, evt EventMeta, payload any, opts ...SendOpt) {
//	val, err := evt.Encode(payload)
//	if err != nil {
//		panic(fmt.Sprintf("error marshaling value: %v", err))
//	}
//
//	// -- construct the topic to write to
//	org := OrganizationFromContext(ctx)
//	if org == "" {
//		panic("no organization in context")
//	}
//
//	key := KeyFromContext(ctx)
//	if key == "" {
//		panic("no key in context")
//	}
//
//	record := &kgo.Record{
//		Key:   []byte(key),
//		Value: val,
//		Topic: fmt.Sprintf("%s.%s", org, evt.Domain),
//		Headers: []kgo.RecordHeader{
//			{Key: events.KindHeader, Value: []byte(evt.String())},
//		},
//	}
//
//	for _, opt := range opts {
//		opt(record)
//	}
//
//	if pr := r.kc.ProduceSync(ctx, record); pr.FirstErr() != nil {
//		panic(fmt.Sprintf("error producing record: %v", pr.FirstErr()))
//	}
//}
