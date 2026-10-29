package client

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	wccfg "github.com/WebCommanderCMS/webcommander-sdk-go/v4/config"
	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
)

// ClientRegistry holds one cached authenticated client per registered credentials block, built on first use.
type ClientRegistry struct {
	// mutex guards every field below, reads included: AddClient grows names, labels and blocks.
	mutex sync.Mutex
	// blocks are validated lazily by GetClient, so one bad block cannot disable the others.
	blocks   []any
	names    map[string]int
	labels   []string
	clients  map[int]*WebCommanderClient
	supplied map[int]bool
}

// NewClientRegistry registers one or more credentials blocks.
func NewClientRegistry(credentials ...any) (*ClientRegistry, error) {
	return newClientRegistry("ClientRegistry", credentials...)
}

// newClientRegistry is NewClientRegistry with the entry point named for the needs-at-least-one message.
func newClientRegistry(entryPoint string, credentials ...any) (*ClientRegistry, error) {
	if len(credentials) == 0 {
		return nil, wcerrors.NewUsageError(entryPoint + " needs at least one credentials block")
	}
	r := &ClientRegistry{
		names:    map[string]int{},
		clients:  map[int]*WebCommanderClient{},
		supplied: map[int]bool{},
	}
	r.blocks = append(r.blocks, credentials...)
	// Names come off the raw blocks, before validation.
	for i, block := range r.blocks {
		name := wccfg.ReadClientName(block)
		label := name
		if name == "" {
			label = fmt.Sprintf(wccfg.UnnamedIndexPattern, i)
		} else if earlier, ok := r.names[name]; ok {
			return nil, wcerrors.NewUsageError(fmt.Sprintf(
				"two instances share clientName '%s' (positions %d and %d); names must be unique so one "+
					"can be told from the other", name, earlier, i))
		} else {
			r.names[name] = i
		}
		r.labels = append(r.labels, label)
	}
	return r, nil
}

// asConfig turns one registered entry into a validated config; index names which entry was wrong.
func asConfig(block any, index int) (wccfg.SdkConfig, error) {
	var raw wccfg.SdkConfig
	switch typedBlock := block.(type) {
	case wccfg.SdkConfig:
		raw = typedBlock
	case wccfg.AuthConfig:
		raw = wccfg.SdkConfig{Auth: typedBlock}
	case map[string]any:
		// config owns the mapping rules and validates before normalising, so no second pass here.
		return wccfg.SdkConfigFromCredentials([]any{typedBlock}, nil, nil)
	default:
		return wccfg.SdkConfig{}, wcerrors.NewUsageError(fmt.Sprintf(
			"credentials[%d]: entries must be an SdkConfig, an AuthConfig or a mapping of "+
				"config fields, got %T", index, block))
	}
	// Validated raw, then normalised once: normaliseAPIVersion is not idempotent.
	if err := raw.Validate(); err != nil {
		return wccfg.SdkConfig{}, err
	}
	return raw.Normalised(), nil
}

// ResolveIndex turns an instance selector into a position: nil, an int (negative counts from the end) or a name.
func (r *ClientRegistry) ResolveIndex(name any) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.resolveIndex(name)
}

// resolveIndex is the body of ResolveIndex; the caller holds r.mutex.
func (r *ClientRegistry) resolveIndex(name any) (int, error) {
	if name == nil {
		return 0, nil
	}
	switch s := name.(type) {
	case int:
		size := len(r.blocks)
		if s >= -size && s < size {
			return ((s % size) + size) % size, nil
		}
		return 0, wcerrors.NewUsageError(fmt.Sprintf("no instance at index %d; %d registered (0-%d)", s, size, size-1))
	case string:
		if index, ok := r.names[s]; ok {
			return index, nil
		}
		return 0, wcerrors.NewUsageError(fmt.Sprintf("unknown instance '%s'; registered: %s", s, strings.Join(r.labels, ", ")))
	default:
		return 0, wcerrors.NewUsageError("name must be nil, an index or a clientName")
	}
}

// GetClient returns a client for one instance, built and cached on first use.
func (r *ClientRegistry) GetClient(name any) (*WebCommanderClient, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	index, err := r.resolveIndex(name)
	if err != nil {
		return nil, err
	}
	if c, ok := r.clients[index]; ok {
		return c, nil
	}
	config, err := asConfig(r.blocks[index], index)
	if err != nil {
		return nil, err
	}
	c, err := NewClient(config)
	if err != nil {
		return nil, err
	}
	r.clients[index] = c
	return c, nil
}

// AddClient registers a pre-built client under a name; its lifetime stays the caller's.
func (r *ClientRegistry) AddClient(name string, client *WebCommanderClient) error {
	if strings.TrimSpace(name) == "" {
		return wcerrors.NewUsageError("addClient needs a non-blank name")
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, ok := r.names[name]; ok {
		return wcerrors.NewUsageError(fmt.Sprintf("instance '%s' is already registered", name))
	}
	r.blocks = append(r.blocks, client.Config()) // already a deep copy
	index := len(r.blocks) - 1
	r.names[name] = index
	r.labels = append(r.labels, name)
	r.clients[index] = client
	r.supplied[index] = true
	return nil
}

// Close closes every client this registry built, trying all even when one fails; a supplied client is skipped.
func (r *ClientRegistry) Close() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	// Sorted indexes, not map order, so the summary and the wrapped first cause are stable.
	built := make([]int, 0, len(r.clients))
	for index := range r.clients {
		if !r.supplied[index] {
			built = append(built, index)
		}
	}
	sort.Ints(built)
	var failures []string
	var first error
	for _, index := range built {
		if err := r.clients[index].Close(); err != nil {
			// The first failure is wrapped below, so errors.Is/As reach the underlying cause.
			failures = append(failures, fmt.Sprintf("%s: %T: %v", r.labels[index], err, err))
			if first == nil {
				first = err
			}
		}
		delete(r.clients, index)
	}
	if len(failures) > 0 {
		return &wcerrors.SdkError{
			Message: fmt.Sprintf(
				"%d of %d clients failed to close (%s); all of them were attempted and dropped, "+
					"so the registry is usable again",
				len(failures), len(built), strings.Join(failures, "; ")),
			Cause: first,
		}
	}
	return nil
}

var (
	registeredMu sync.Mutex
	registered   []wccfg.SdkConfig
)

// RegisterInstances registers every instance this process can talk to; the first is the default.
func RegisterInstances(credentials ...any) error {
	declared, err := newClientRegistry("RegisterInstances", credentials...)
	if err != nil {
		return err
	}
	blocks := make([]wccfg.SdkConfig, len(declared.blocks))
	for index, block := range declared.blocks {
		config, err := asConfig(block, index)
		if err != nil {
			return err
		}
		blocks[index] = config
	}
	registeredMu.Lock()
	defer registeredMu.Unlock()
	registered = blocks
	return nil
}

// BuildClient is the one-call entry point: a client for a registered instance.
func BuildClient(instance any, overrides ...func(*wccfg.SdkConfig)) (*WebCommanderClient, error) {
	registeredMu.Lock()
	snapshot := append([]wccfg.SdkConfig(nil), registered...)
	registeredMu.Unlock()
	if len(snapshot) == 0 {
		return nil, errNoInstancesRegistered()
	}
	resolver, err := NewClientRegistry(toAnySlice(snapshot)...)
	if err != nil {
		return nil, err
	}
	index, err := resolver.ResolveIndex(instance)
	if err != nil {
		return nil, err
	}
	// Deep-copied: a struct copy would alias DefaultHeaders and the code slices to the registration.
	chosen := snapshot[index].DeepCopy()
	for _, override := range overrides {
		if override != nil {
			override(&chosen)
		}
	}
	return NewClient(chosen)
}

// ListInstanceNames returns one label per registered instance, default first.
func ListInstanceNames() ([]string, error) {
	blocks, err := GetRegisteredInstances()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(blocks))
	for i, block := range blocks {
		names[i] = block.Auth.ClientName
	}
	return names, nil
}

// GetRegisteredInstances returns a deep copy of the registered blocks, or an error when none are.
func GetRegisteredInstances() ([]wccfg.SdkConfig, error) {
	registeredMu.Lock()
	defer registeredMu.Unlock()
	if len(registered) == 0 {
		return nil, errNoInstancesRegistered()
	}
	out := make([]wccfg.SdkConfig, len(registered))
	for i, block := range registered {
		out[i] = block.DeepCopy()
	}
	return out, nil
}

// errNoInstancesRegistered is the one message all three package-level entry points give.
func errNoInstancesRegistered() error {
	return wcerrors.NewUsageError("no instances registered — call RegisterInstances(...) first, or build a " +
		"registry of your own with NewClientRegistry(credentials...)")
}

func toAnySlice(configs []wccfg.SdkConfig) []any {
	out := make([]any, len(configs))
	for i, c := range configs {
		out[i] = c
	}
	return out
}
