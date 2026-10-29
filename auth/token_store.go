package auth

import (
	"context"
	"maps"
	"sync"
	"sync/atomic"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
)

// TokenRefreshScope is the token-store operations that run under an already-held refresh lock.
type TokenRefreshScope interface {
	// ReadTokenRecord reads the stored token under the held lock.
	ReadTokenRecord() (map[string]any, error)
	// WriteTokenRecord persists a token under the held lock.
	WriteTokenRecord(record map[string]any) error
	// ClearTokenRecord discards the stored token under the held lock.
	ClearTokenRecord() error
	// Release releases the lock.
	Release()
}

// NestedTokenRefreshScope is the optional half of TokenRefreshScope: a second scope over the held lock.
type NestedTokenRefreshScope interface {
	TokenRefreshScope
	// EnterNested adds a level to the held lock and returns a scope for it.
	EnterNested() (TokenRefreshScope, error)
}

// TokenStore is where an access token lives between calls: persistence plus refresh coordination.
type TokenStore interface {
	// ReadTokenRecord returns the stored token, or nil when there is none.
	ReadTokenRecord() (map[string]any, error)
	// WriteTokenRecord persists a token, replacing whatever was there.
	WriteTokenRecord(record map[string]any) error
	// ClearTokenRecord discards the stored token; a no-op when there is none.
	ClearTokenRecord() error
	// AcquireTokenRefreshLock holds the refresh lock for the lifetime of the scope.
	AcquireTokenRefreshLock() (TokenRefreshScope, error)
}

// ContextTokenStore is the optional half of TokenStore: a lock wait abandoned when the context is done.
type ContextTokenStore interface {
	TokenStore
	// AcquireTokenRefreshLockContext is AcquireTokenRefreshLock, abandoning the wait itself when ctx is done.
	AcquireTokenRefreshLockContext(ctx context.Context) (TokenRefreshScope, error)
}

// refreshGate is the refresh lock: a one-slot channel, since a sync.Mutex wait cannot be abandoned.
type refreshGate struct {
	once sync.Once
	slot chan struct{}
}

// channel is the slot, created on first use.
func (g *refreshGate) channel() chan struct{} {
	g.once.Do(func() { g.slot = make(chan struct{}, 1) })
	return g.slot
}

// acquire waits for the lock with no way to abandon it, as the context-free TokenStore interface needs.
func (g *refreshGate) acquire() { g.channel() <- struct{}{} }

// acquireContext waits for the lock and gives up when ctx is done, returning ctx.Err().
func (g *refreshGate) acquireContext(ctx context.Context) error {
	// Tested before the wait too, so a caller already withdrawn never takes a lock it cannot use.
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case g.channel() <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (g *refreshGate) release() { <-g.channel() }

// MemoryTokenStore is in-process token storage, isolated per instance, with nothing on disk.
type MemoryTokenStore struct {
	gate   refreshGate
	record map[string]any
	// mutex guards generation.
	mutex sync.Mutex
	// generation counts how often this store's lock was taken, so a stale scope is told from its successor.
	generation uint64
}

// NewMemoryTokenStore builds an in-memory token store.
func NewMemoryTokenStore() *MemoryTokenStore { return &MemoryTokenStore{} }

// ReadTokenRecord returns a shallow copy of the stored token, whose values are scalars, or nil.
func (s *MemoryTokenStore) ReadTokenRecord() (map[string]any, error) {
	s.gate.acquire()
	defer s.gate.release()
	return s.readCore(), nil
}

// WriteTokenRecord persists a token; takes and releases the lock.
func (s *MemoryTokenStore) WriteTokenRecord(record map[string]any) error {
	scope, err := s.AcquireTokenRefreshLock()
	if err != nil {
		return err
	}
	defer scope.Release()
	return scope.WriteTokenRecord(record)
}

// ClearTokenRecord discards the stored token; takes and releases the lock.
func (s *MemoryTokenStore) ClearTokenRecord() error {
	scope, err := s.AcquireTokenRefreshLock()
	if err != nil {
		return err
	}
	defer scope.Release()
	return scope.ClearTokenRecord()
}

// AcquireTokenRefreshLock holds the in-process lock for the lifetime of the scope.
func (s *MemoryTokenStore) AcquireTokenRefreshLock() (TokenRefreshScope, error) {
	s.gate.acquire()
	return &memoryScope{store: s, generation: s.enter()}, nil
}

// AcquireTokenRefreshLockContext is AcquireTokenRefreshLock with the wait itself bounded by ctx.
func (s *MemoryTokenStore) AcquireTokenRefreshLockContext(ctx context.Context) (TokenRefreshScope, error) {
	if err := s.gate.acquireContext(ctx); err != nil {
		return nil, err
	}
	return &memoryScope{store: s, generation: s.enter()}, nil
}

// enter records one more hold and reports the generation it is.
func (s *MemoryTokenStore) enter() uint64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.generation++
	return s.generation
}

// holds reports whether the hold taken at generation is still the live one.
func (s *MemoryTokenStore) holds(generation uint64) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.generation == generation
}

// readCore is the read both the store and its scope answer with: nil, or a shallow copy.
func (s *MemoryTokenStore) readCore() map[string]any {
	if s.record == nil {
		return nil
	}
	return copyMap(s.record)
}

// memoryScope is one hold of the memory store's lock.
type memoryScope struct {
	store *MemoryTokenStore
	// generation is the hold this scope belongs to, so a scope outliving its hold is refused.
	generation uint64
	// released is atomic because Release need not run on the goroutine that took the scope.
	released atomic.Bool
}

func (m *memoryScope) live(operation string) error {
	return scopeLiveness(operation, m.released.Load(), m.store.holds(m.generation))
}

// scopeLiveness refuses an operation on a released scope or on one whose lock was given up.
func scopeLiveness(operation string, released, held bool) error {
	if released {
		return wcerrors.NewUsageError(operation + ": the token refresh scope was already released; acquire a new one")
	}
	if !held {
		return wcerrors.NewUsageError(operation + ": the token refresh scope's lock was given up; acquire a new one")
	}
	return nil
}

// ReadTokenRecord reads under the held lock, with MemoryTokenStore.ReadTokenRecord's copy rule.
func (m *memoryScope) ReadTokenRecord() (map[string]any, error) {
	if err := m.live("ReadTokenRecord"); err != nil {
		return nil, err
	}
	return m.store.readCore(), nil
}

// WriteTokenRecord stores a shallow copy, so a later write to the caller's map cannot replace the token.
func (m *memoryScope) WriteTokenRecord(record map[string]any) error {
	if err := m.live("WriteTokenRecord"); err != nil {
		return err
	}
	m.store.record = copyMap(record)
	return nil
}

func (m *memoryScope) ClearTokenRecord() error {
	if err := m.live("ClearTokenRecord"); err != nil {
		return err
	}
	m.store.record = nil
	return nil
}

func (m *memoryScope) Release() {
	if m.released.CompareAndSwap(false, true) {
		m.store.gate.release()
	}
}

// copyMap is the one shallow copy this store makes, in both directions.
func copyMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	maps.Copy(out, m)
	return out
}
