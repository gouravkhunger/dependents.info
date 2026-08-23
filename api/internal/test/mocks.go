package test

import (
	"context"
	"errors"
	"sync"
	"time"

	"dependents.info/internal/models"
)

var ErrNotFound = errors.New("key not found")

type MockStore struct {
	mu   sync.RWMutex
	Data map[string][]byte
}

func NewMockStore() *MockStore {
	return &MockStore{Data: make(map[string][]byte)}
}

func (m *MockStore) Get(key string, out *string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.Data[key]
	if !ok {
		return ErrNotFound
	}
	*out = string(v)
	return nil
}

func (m *MockStore) Save(key string, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Data[key] = data
	return nil
}

func (m *MockStore) SaveWithTTL(key string, data []byte, _ time.Duration) error {
	return m.Save(key, data)
}

func (m *MockStore) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Data, key)
	return nil
}

func (m *MockStore) IterateKeys(callback func(key string)) {
	m.mu.RLock()
	keys := make([]string, 0, len(m.Data))
	for k := range m.Data {
		keys = append(keys, k)
	}
	m.mu.RUnlock()
	for _, k := range keys {
		callback(k)
	}
}

func (m *MockStore) Close() error { return nil }
func (m *MockStore) Sync() error  { return nil }

type MockRenderer struct {
	SVGResult     []byte
	SVGErr        error
	PageResult    []byte
	PageErr       error
	SitemapResult []byte
	SitemapErr    error
}

func (m *MockRenderer) RenderSVG(_ models.IngestRequest) ([]byte, error) {
	return m.SVGResult, m.SVGErr
}

func (m *MockRenderer) RenderPage(_ models.RepoPage) ([]byte, error) {
	return m.PageResult, m.PageErr
}

func (m *MockRenderer) RenderSitemap(_ []string) ([]byte, error) {
	return m.SitemapResult, m.SitemapErr
}

type MockDependentsTasker struct {
	NewTaskFn func(ctx context.Context, repo, id, kind string, callback func(int, []byte)) error
}

func (m *MockDependentsTasker) NewTask(ctx context.Context, repo, id, kind string, callback func(int, []byte)) error {
	if m.NewTaskFn != nil {
		return m.NewTaskFn(ctx, repo, id, kind, callback)
	}
	return nil
}

type MockOIDCVerifier struct {
	Err error
}

func (m *MockOIDCVerifier) VerifyToken(_ context.Context, _ string, _ string) error {
	return m.Err
}
