package memphis

import (
	"fmt"
	"github.com/memphisdev/memphis.go"
	go_shono "github.com/shono-io/go-shono"
	"github.com/shono-io/go-shono/events"
	"github.com/sirupsen/logrus"
)

var (
	HostProp     = "memphis.host"
	UsernameProp = "memphis.username"
	PasswordProp = "memphis.password"
)

func NewMemphisBackbone(id string, props map[string]any) (*MemphisBackbone, error) {
	// -- get the host from the config
	host, ok := props[HostProp]
	if !ok || host == "" {
		return nil, fmt.Errorf("missing property %s", HostProp)
	}

	// -- get the username from the config
	username, ok := props[UsernameProp]
	if !ok || username == "" {
		return nil, fmt.Errorf("missing property %s", UsernameProp)
	}

	// -- get the password from the config
	password, ok := props[PasswordProp]
	if !ok || password == "" {
		return nil, fmt.Errorf("missing property %s", PasswordProp)
	}

	c, err := memphis.Connect(host.(string), username.(string), memphis.Password(password.(string)))
	if err != nil {
		return nil, err
	}

	return &MemphisBackbone{
		id: id,
		c:  c,
		w:  NewWriter(id, c),
	}, nil
}

type MemphisBackbone struct {
	id  string
	c   *memphis.Conn
	w   *Writer
	run *Runner
}

func (b *MemphisBackbone) Close() {
	if b.run != nil {
		b.run.Close()
	}

	b.c.Close()
}

func (b *MemphisBackbone) MustWrite(evt *go_shono.EventMeta, payload any) {
	b.w.MustWrite(evt, payload)
}

func (b *MemphisBackbone) Write(evt *go_shono.EventMeta, payload any) error {
	return b.w.Write(evt, payload)
}

func (b *MemphisBackbone) Listen(r *go_shono.Router) error {
	if b.run != nil {
		return fmt.Errorf("already listening")
	}

	b.run = NewRunner(b.id, r, b.c)
	return b.run.Run()
}

func (b *MemphisBackbone) Apply(eid go_shono.EventId, event any) error {
	switch eid {
	case events.ScopeCreated.EventId:
		return b.onScopeCreated(event.(*events.ScopeCreatedEvent))
	case events.ScopeDeleted.EventId:
		return b.onScopeDeleted(event.(*events.ScopeDeletedEvent))
	}

	return nil
}

func (b *MemphisBackbone) onScopeCreated(event *events.ScopeCreatedEvent) error {
	stationName := fmt.Sprintf("%s.%s", event.Organization, event.Code)
	logrus.Debugf("creating station: %s", stationName)
	_, err := b.c.CreateStation(stationName)
	return err
}

func (b *MemphisBackbone) onScopeDeleted(event *events.ScopeDeletedEvent) error {
	stationName := fmt.Sprintf("%s.%s", event.Organization, event.Code)
	logrus.Debugf("creating station: %s", stationName)
	s, err := b.c.CreateStation(stationName)
	if err != nil {
		return fmt.Errorf("failed to retrieve station: %v", err)
	}
	return s.Destroy()
}
