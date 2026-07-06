package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/ilexPar/simple-kube/pkg/base"
	"github.com/ilexPar/simple-kube/pkg/core"
	skerr "github.com/ilexPar/simple-kube/pkg/errors"
)

// The core package holds the CRUD flow shared by every resource in both scopes,
// so it is tested once here against fakes instead of once per resource.

type fakeRaw struct{ Name string }
type fakeSimple struct{ Name string }

type fakeCodec struct {
	decodeErr error
}

func (c fakeCodec) Encode(in fakeSimple) (*fakeRaw, error) {
	return &fakeRaw{Name: in.Name}, nil
}

func (c fakeCodec) Decode(raw *fakeRaw, out *fakeSimple) error {
	if c.decodeErr != nil {
		return c.decodeErr
	}
	out.Name = raw.Name
	return nil
}

type fakeBackend struct {
	getObj    *fakeRaw
	getErr    error
	listObjs  []fakeRaw
	listOpts  base.QueryOpts
	created   *fakeRaw
	updated   *fakeRaw
	deletedID string
	events    *[]string
}

func (b *fakeBackend) Get(id string) (*fakeRaw, error) {
	b.record("get")
	return b.getObj, b.getErr
}

func (b *fakeBackend) List(opts base.QueryOpts) ([]fakeRaw, error) {
	b.listOpts = opts
	return b.listObjs, nil
}

func (b *fakeBackend) Create(obj *fakeRaw) error {
	b.record("send")
	b.created = obj
	return nil
}

func (b *fakeBackend) Update(obj *fakeRaw) error {
	b.record("send")
	b.updated = obj
	return nil
}

func (b *fakeBackend) Delete(id string) error {
	b.deletedID = id
	return nil
}

func (b *fakeBackend) record(ev string) {
	if b.events != nil {
		*b.events = append(*b.events, ev)
	}
}

func newAction(codec core.Codec[fakeSimple, fakeRaw], backend core.Backend[fakeRaw]) *core.Action[fakeSimple, fakeRaw] {
	return core.NewAction[fakeSimple, fakeRaw](codec, backend)
}

func TestGet(t *testing.T) {
	t.Run("decodes the fetched object", func(t *testing.T) {
		act := newAction(fakeCodec{}, &fakeBackend{getObj: &fakeRaw{Name: "obj"}})
		out, err := act.Get("obj").Run()
		assert.NoError(t, err)
		assert.Equal(t, "obj", out.Name)
	})
	t.Run("runs DataHandler before decode", func(t *testing.T) {
		var order []string
		act := newAction(fakeCodec{}, &fakeBackend{getObj: &fakeRaw{Name: "obj"}, events: &order})
		out, err := act.Get("obj").
			DataHandler(func(raw *fakeRaw) error {
				order = append(order, "callback")
				raw.Name = "mutated"
				return nil
			}).
			Run()
		assert.NoError(t, err)
		assert.Equal(t, "mutated", out.Name)
		assert.Equal(t, []string{"get", "callback"}, order)
	})
	t.Run("aborts on DataHandler error", func(t *testing.T) {
		act := newAction(fakeCodec{}, &fakeBackend{getObj: &fakeRaw{Name: "obj"}})
		_, err := act.Get("obj").
			DataHandler(func(*fakeRaw) error { return errors.New("boom") }).
			Run()
		assert.EqualError(t, err, "boom")
	})
	t.Run("normalizes not-found errors", func(t *testing.T) {
		notFound := kerrors.NewNotFound(schema.GroupResource{Resource: "things"}, "obj")
		act := newAction(fakeCodec{}, &fakeBackend{getErr: notFound})
		_, err := act.Get("obj").Run()
		assert.EqualError(t, err, skerr.ERROR_NOT_FOUND)
	})
}

func TestPut(t *testing.T) {
	t.Run("create encodes then sends", func(t *testing.T) {
		var order []string
		backend := &fakeBackend{events: &order}
		err := newAction(fakeCodec{}, backend).Create(fakeSimple{Name: "obj"}).Run()
		assert.NoError(t, err)
		assert.Equal(t, &fakeRaw{Name: "obj"}, backend.created)
	})
	t.Run("update binds to backend Update", func(t *testing.T) {
		backend := &fakeBackend{}
		err := newAction(fakeCodec{}, backend).Update(fakeSimple{Name: "obj"}).Run()
		assert.NoError(t, err)
		assert.Equal(t, &fakeRaw{Name: "obj"}, backend.updated)
		assert.Nil(t, backend.created)
	})
	t.Run("runs DataHandler between encode and send", func(t *testing.T) {
		var order []string
		backend := &fakeBackend{events: &order}
		err := newAction(fakeCodec{}, backend).
			Create(fakeSimple{Name: "obj"}).
			DataHandler(func(*fakeRaw) error {
				order = append(order, "callback")
				return nil
			}).
			Run()
		assert.NoError(t, err)
		assert.Equal(t, []string{"callback", "send"}, order)
	})
	t.Run("aborts before send on DataHandler error", func(t *testing.T) {
		backend := &fakeBackend{}
		err := newAction(fakeCodec{}, backend).
			Create(fakeSimple{Name: "obj"}).
			DataHandler(func(*fakeRaw) error { return errors.New("boom") }).
			Run()
		assert.EqualError(t, err, "boom")
		assert.Nil(t, backend.created)
	})
}

func TestList(t *testing.T) {
	t.Run("decodes every item", func(t *testing.T) {
		backend := &fakeBackend{listObjs: []fakeRaw{{Name: "a"}, {Name: "b"}}}
		out, err := newAction(fakeCodec{}, backend).List().Run()
		assert.NoError(t, err)
		assert.Equal(t, []fakeSimple{{Name: "a"}, {Name: "b"}}, out)
	})
	t.Run("propagates decode errors", func(t *testing.T) {
		backend := &fakeBackend{listObjs: []fakeRaw{{Name: "a"}}}
		_, err := newAction(fakeCodec{decodeErr: errors.New("bad")}, backend).List().Run()
		assert.EqualError(t, err, "bad")
	})
	t.Run("FilterByLabels forwards a selector", func(t *testing.T) {
		backend := &fakeBackend{}
		_, err := newAction(fakeCodec{}, backend).
			List().
			FilterByLabels(map[string]string{"app": "nginx"}).
			Run()
		assert.NoError(t, err)
		assert.Equal(t, "app=nginx", backend.listOpts.List.LabelSelector)
	})
}

func TestDelete(t *testing.T) {
	backend := &fakeBackend{}
	err := newAction(fakeCodec{}, backend).Delete("obj").Run()
	assert.NoError(t, err)
	assert.Equal(t, "obj", backend.deletedID)
}
