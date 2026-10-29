package auth

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	wcerrors "github.com/WebCommanderCMS/webcommander-sdk-go/v4/errors"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/http"
	"github.com/WebCommanderCMS/webcommander-sdk-go/v4/redact"
)

// FileTokenStore is token storage in a JSON file, safe to share across processes.
type FileTokenStore struct {
	// shared is this sidecar path's gate, descriptor and depth, shared by every store naming the path.
	shared    *sharedFileLock
	path      string
	lockPath  string
	namespace string
	lockWait  time.Duration
}

// DefaultTokenCacheFilename is the file BuildDefaultTokenStore caches tokens in, per user.
const DefaultTokenCacheFilename = "wc_go_token.json" //#nosec G101 -- a filename, not a credential

const (
	// DefaultLockWait is how long an acquisition waits for another process before it errors.
	DefaultLockWait = 5 * time.Minute

	lockPoll = 50 * time.Millisecond

	replaceRetryCount   = 5
	replaceRetryInitial = 10 * time.Millisecond
)

// sharedFileLock is what every store in this process addressing one sidecar path shares.
type sharedFileLock struct {
	// gate is the in-process half of the refresh lock.
	gate refreshGate
	// mutex guards the three fields below.
	mutex sync.Mutex
	// handle holds the operating-system lock, or is nil when nothing holds it.
	handle *os.File
	// depth is how many scopes are outstanding.
	depth int
	// generation counts how many times this path's lock has been taken.
	generation uint64
}

// sharedLocks holds one sharedFileLock per canonical sidecar path, for the life of the process.
var (
	sharedLocksMutex sync.Mutex
	sharedLocks      = map[string]*sharedFileLock{}
)

// sharedLockFor is the shared lock a sidecar path resolves to, created on first use.
func sharedLockFor(lockPath string) *sharedFileLock {
	key := canonicalLockKey(lockPath)
	sharedLocksMutex.Lock()
	defer sharedLocksMutex.Unlock()
	if existing, held := sharedLocks[key]; held {
		return existing
	}
	if aliased := aliasOf(key); aliased != nil {
		// Recorded under this spelling too, so the next lookup skips the filesystem walk.
		sharedLocks[key] = aliased
		return aliased
	}
	created := &sharedFileLock{}
	sharedLocks[key] = created
	return created
}

// aliasOf is the entry already registered for the same sidecar under another name, or nil.
func aliasOf(key string) *sharedFileLock {
	base := foldLockKey(filepath.Base(key))
	for otherKey, other := range sharedLocks {
		if foldLockKey(filepath.Base(otherKey)) != base {
			continue
		}
		if sameDirectory(filepath.Dir(key), filepath.Dir(otherKey)) {
			return other
		}
	}
	return nil
}

// sameDirectory reports whether two paths name one directory, by the platform's own file identity.
func sameDirectory(left, right string) bool {
	leftInfo, err := os.Stat(left)
	if err != nil {
		return false
	}
	rightInfo, err := os.Stat(right)
	if err != nil {
		return false
	}
	return os.SameFile(leftInfo, rightInfo)
}

// canonicalLockKey is the key two stores must agree on to share a lock: the sidecar's folded real path.
func canonicalLockKey(lockPath string) string {
	absolute, err := filepath.Abs(lockPath)
	if err != nil {
		absolute = lockPath
	}
	if real, err := filepath.EvalSymlinks(absolute); err == nil {
		absolute = real
	}
	return foldLockKey(absolute)
}

// enter records a freshly taken OS lock as this path's, at depth one, and reports its generation.
func (l *sharedFileLock) enter(handle *os.File) uint64 {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.handle = handle
	l.depth = 1
	l.generation++
	return l.generation
}

// nest adds a reentrant level over the lock already held and reports that hold's generation.
func (l *sharedFileLock) nest() uint64 {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.depth++
	return l.generation
}

// holds reports whether the hold taken at generation is still the live one.
func (l *sharedFileLock) holds(generation uint64) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.depth > 0 && l.generation == generation
}

// leave gives back one level, dropping the operating-system lock and the gate on the last one.
func (l *sharedFileLock) leave() {
	l.mutex.Lock()
	if l.depth > 1 {
		l.depth--
		l.mutex.Unlock()
		return
	}
	l.depth = 0
	handle := l.handle
	l.handle = nil
	l.mutex.Unlock()
	if handle != nil {
		unlockExclusive(handle)
		_ = handle.Close()
	}
	l.gate.release()
}

// NewFileTokenStore opens a store over a file path, keyed by a namespace, waiting DefaultLockWait.
func NewFileTokenStore(filePath, namespace string) (*FileTokenStore, error) {
	return NewFileTokenStoreWithLockWait(filePath, namespace, DefaultLockWait)
}

// NewFileTokenStoreWithLockWait is NewFileTokenStore with an explicit lock wait; non-positive is the default.
func NewFileTokenStoreWithLockWait(filePath, namespace string, lockWait time.Duration) (*FileTokenStore, error) {
	if namespace == "" {
		namespace = "default"
	}
	if lockWait <= 0 {
		lockWait = DefaultLockWait
	}
	dataPath, err := resolveStorePath(filePath)
	if err != nil {
		return nil, err
	}
	lockPath := dataPath + ".lock"
	return &FileTokenStore{
		shared:    sharedLockFor(lockPath),
		path:      dataPath,
		lockPath:  lockPath,
		namespace: namespace,
		lockWait:  lockWait,
	}, nil
}

// resolveStorePath settles, once, the path the data file, sidecar and in-process gate derive from.
func resolveStorePath(filePath string) (string, error) {
	absolute, err := filepath.Abs(filePath)
	if err != nil {
		return "", withDetail(wcerrors.NewConfigError(
			"tokenCachePath: the path could not be resolved against the working directory; "+
				"name an absolute path, or keep the token in memory"), filePath, err)
	}
	dir := filepath.Dir(absolute)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", withDetail(wcerrors.NewConfigError(fmt.Sprintf(
			"tokenCachePath: the directory for %s could not be created; name a path whose parent "+
				"exists, or keep the token in memory", redact.DescribePath(absolute))), dir, err)
	}
	// The whole path first: on a case-insensitive volume it returns the name the volume holds.
	if real, err := filepath.EvalSymlinks(absolute); err == nil {
		return real, nil
	}
	realDir, evalErr := filepath.EvalSymlinks(dir)
	if evalErr != nil {
		// A directory that will not resolve keys itself, still absolute.
		realDir = dir
	}
	return filepath.Join(realDir, filepath.Base(absolute)), nil
}

// BuildDefaultTokenCachePath is the per-user cache path for this platform, under a webcommander directory.
func BuildDefaultTokenCachePath(filename string) (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		if home, homeErr := os.UserHomeDir(); homeErr == nil {
			base = home
		}
	}
	if !filepath.IsAbs(base) {
		// os.UserCacheDir returns a relative LocalAppData, and os.UserHomeDir a relative home, with a nil error.
		base = os.TempDir()
	}
	root := filepath.Join(base, "webcommander")
	if err := os.MkdirAll(root, 0o700); err != nil {
		return "", withDetail(wcerrors.NewConfigError(
			"the default token cache directory could not be created; pass a path, or keep the "+
				"token in memory"), root, err)
	}
	// Applied even to an existing directory: MkdirAll's mode covers only directories it creates.
	restrictToOwner(root, true)
	return filepath.Join(root, filename), nil
}

// restrictToOwner is a best-effort 0700 on a directory or 0600 on a file.
func restrictToOwner(path string, directory bool) {
	mode := os.FileMode(0o600)
	if directory {
		mode = 0o700
	}
	_ = os.Chmod(path, mode)
}

// ReadTokenRecord reads under the in-process gate only: atomic writes make a reader see a whole file.
func (s *FileTokenStore) ReadTokenRecord() (map[string]any, error) {
	s.shared.gate.acquire()
	defer s.shared.gate.release()
	return s.readNamespaceRecord()
}

// WriteTokenRecord persists a token; takes and releases the lock.
func (s *FileTokenStore) WriteTokenRecord(record map[string]any) error {
	scope, err := s.AcquireTokenRefreshLock()
	if err != nil {
		return err
	}
	defer scope.Release()
	return scope.WriteTokenRecord(record)
}

// ClearTokenRecord discards the stored token; takes and releases the lock.
func (s *FileTokenStore) ClearTokenRecord() error {
	scope, err := s.AcquireTokenRefreshLock()
	if err != nil {
		return err
	}
	defer scope.Release()
	return scope.ClearTokenRecord()
}

// AcquireTokenRefreshLock holds the in-process and cross-process locks for the scope's lifetime.
func (s *FileTokenStore) AcquireTokenRefreshLock() (TokenRefreshScope, error) {
	return s.AcquireTokenRefreshLockContext(context.Background())
}

// AcquireTokenRefreshLockContext is AcquireTokenRefreshLock with both lock waits also bounded by ctx.
func (s *FileTokenStore) AcquireTokenRefreshLockContext(ctx context.Context) (TokenRefreshScope, error) {
	if err := s.shared.gate.acquireContext(ctx); err != nil {
		return nil, err
	}
	handle, err := s.takeTheFileLock(ctx)
	if err != nil {
		// Every exit that hands no scope back must release the gate, or later acquisitions wait forever.
		s.shared.gate.release()
		return nil, err
	}
	return &fileScope{store: s, generation: s.shared.enter(handle)}, nil
}

// takeTheFileLock opens the sidecar and takes the OS lock on it, waiting out a peer that holds it.
func (s *FileTokenStore) takeTheFileLock(ctx context.Context) (*os.File, error) {
	handle, err := s.openTheSidecar()
	if err != nil {
		return nil, err
	}
	deadline := time.Now().Add(s.lockWait)
	for {
		locked, lockErr := tryLockExclusive(handle)
		if locked {
			return handle, nil
		}
		if lockErr != nil {
			_ = handle.Close()
			return nil, withDetail(&wcerrors.SdkError{
				Message: fmt.Sprintf(
					"token cache lock beside %s could not be taken; use a lockable path or TokenCacheMemoryOnly",
					redact.DescribePath(s.path)),
			}, s.lockPath, lockErr)
		}
		if time.Now().After(deadline) {
			_ = handle.Close()
			return nil, withDetail(&wcerrors.SdkError{
				Message: fmt.Sprintf("token cache lock beside %s held longer than %s",
					redact.DescribePath(s.path), s.lockWait),
			}, s.lockPath, nil)
		}
		select {
		case <-ctx.Done():
			_ = handle.Close()
			return nil, ctx.Err()
		case <-time.After(lockPoll):
		}
	}
}

func (s *FileTokenStore) openTheSidecar() (*os.File, error) {
	remadeTheDirectory := false
	for {
		handle, err := os.OpenFile(s.lockPath, os.O_CREATE|os.O_RDWR, 0o600) //#nosec G304 -- the sidecar path is derived from the caller's own cache path
		if err == nil {
			return handle, nil
		}
		// A tmp-reaper may have removed the directory since: the one failure here that can still come good.
		if errors.Is(err, fs.ErrNotExist) && !remadeTheDirectory {
			remadeTheDirectory = true
			mkdirErr := os.MkdirAll(filepath.Dir(s.lockPath), 0o700)
			if mkdirErr == nil {
				continue
			}
			return nil, withDetail(&wcerrors.SdkError{
				Message: fmt.Sprintf(
					"the directory holding the token cache lock file beside %s is gone and could "+
						"not be recreated, so the cache cannot be locked; every process sharing it "+
						"would refresh at once", redact.DescribePath(s.path)),
			}, s.lockPath, mkdirErr)
		}
		return nil, withDetail(wcerrors.NewConfigError(fmt.Sprintf(
			"tokenCachePath: the token cache lock file beside %s could not be opened, so the cache "+
				"cannot be locked and this call was not attempted. Name a path this process can "+
				"write, or keep the token in memory", redact.DescribePath(s.path))), s.lockPath, err)
	}
}

// cacheState is which of five states a reading of the cache file found it in.
type cacheState int

const (
	cacheAbsent cacheState = iota
	cacheEmpty
	cacheRecords
	cacheCorrupt
	cacheUnreadable
)

func (c cacheState) description() string {
	switch c {
	case cacheAbsent:
		return "not there"
	case cacheEmpty:
		return "empty"
	case cacheRecords:
		return "readable"
	case cacheCorrupt:
		return "present but not a JSON object"
	default:
		return "present but could not be read"
	}
}

// known reports whether the file's contents are known: what a read may answer from and a write replace.
func (c cacheState) known() bool {
	return c == cacheAbsent || c == cacheEmpty || c == cacheRecords
}

// cacheReading is one reading of the cache file: its state, entries, and what made it unreadable.
type cacheReading struct {
	state   cacheState
	records map[string]any
	cause   error
}

// mutableRecords is a copy of the entries a write can build on, empty unless the state is cacheRecords.
func (r cacheReading) mutableRecords() map[string]any {
	return copyMap(r.records)
}

// unknownCache is the error a read or a write reports for a file whose contents are unknown.
func (r cacheReading) unknownCache(path, refused string) error {
	return withDetail(&wcerrors.SdkError{
		Message: fmt.Sprintf("token cache at %s is %s, so %s; clear the token to rebuild it",
			redact.DescribePath(path), r.state.description(), refused),
	}, path, r.cause)
}

// MaxTokenCacheBytes caps how much of the token cache file is read into memory: 4 MiB.
const MaxTokenCacheBytes = 4 << 20

func (s *FileTokenStore) readTheCache() cacheReading {
	data, err := http.ReadCappedFile(s.path, MaxTokenCacheBytes)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return cacheReading{state: cacheAbsent}
		}
		return cacheReading{state: cacheUnreadable, cause: err}
	}
	if len(strings.TrimSpace(string(data))) == 0 {
		return cacheReading{state: cacheEmpty}
	}
	decoded, ok := http.DecodeJSON(data)
	if !ok {
		return cacheReading{state: cacheCorrupt}
	}
	records, ok := decoded.(map[string]any)
	if !ok {
		return cacheReading{state: cacheCorrupt}
	}
	return cacheReading{state: cacheRecords, records: records}
}

// readNamespaceRecord is this namespace's entry, or nil when the cache holds none for it.
func (s *FileTokenStore) readNamespaceRecord() (map[string]any, error) {
	reading := s.readTheCache()
	if !reading.state.known() {
		return nil, reading.unknownCache(s.path, "nothing was read from it")
	}
	stored, present := reading.records[s.namespace]
	if !present || stored == nil {
		return nil, nil
	}
	entry, isObject := stored.(map[string]any)
	if !isObject {
		return nil, corruptRecordError(fmt.Sprintf("it is %s rather than an object", http.JSONKind(stored)))
	}
	return entry, nil
}

// writeNamespaceRecord stores this namespace's entry, keeping every sibling namespace intact.
func (s *FileTokenStore) writeNamespaceRecord(record map[string]any) error {
	reading := s.readTheCache()
	if !reading.state.known() {
		return reading.unknownCache(s.path, "this entry was not written")
	}
	stored := reading.mutableRecords()
	stored[s.namespace] = record
	return s.writeAllNamespaces(stored)
}

// clearNamespaceRecord discards this namespace's entry, and the file itself once it holds no other.
func (s *FileTokenStore) clearNamespaceRecord() error {
	reading := s.readTheCache()
	if !reading.state.known() {
		return deletionOutcome(s.path, os.Remove(s.path))
	}
	stored := reading.mutableRecords()
	if _, held := stored[s.namespace]; !held {
		return nil
	}
	delete(stored, s.namespace)
	if len(stored) > 0 {
		return s.writeAllNamespaces(stored)
	}
	return deletionOutcome(s.path, os.Remove(s.path))
}

func deletionOutcome(path string, err error) error {
	if err == nil || errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return withDetail(&wcerrors.SdkError{
		Message: fmt.Sprintf("token cache file %s could not be deleted", redact.DescribePath(path)),
	}, path, err)
}

// withDetail attaches the full path and the cause to a refusal, for Explain to render and Error to withhold.
func withDetail(refusal *wcerrors.SdkError, path string, cause error) error {
	refusal.Cause = cause
	if cause == nil {
		refusal.CauseText = path
		return refusal
	}
	refusal.CauseText = fmt.Sprintf("%s: %v", path, pathErrorReason(cause))
	return refusal
}

// pathErrorReason unwraps an *fs.PathError, whose Error() repeats the operation and path.
func pathErrorReason(err error) error {
	var pathError *fs.PathError
	if errors.As(err, &pathError) && pathError.Err != nil {
		return pathError.Err
	}
	return err
}

func (s *FileTokenStore) writeAllNamespaces(stored map[string]any) error {
	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	// BOM-less UTF-8 without HTML escaping, so the file stays byte-comparable and plain JSON.
	data, err := http.EncodeJSONBody(stored)
	if err != nil {
		return err
	}
	temporaryPath := fmt.Sprintf("%s.%d.tmp", s.path, time.Now().UnixNano())
	if err := writeNewFile(temporaryPath, data); err != nil {
		return err
	}
	if err := replaceFile(temporaryPath, s.path); err != nil {
		_ = os.Remove(temporaryPath)
		return err
	}
	// The replace carries 0600 across on POSIX, but not over a Windows destination or a widened umask.
	restrictToOwner(s.path, false)
	return nil
}

// writeNewFile writes data to a path that must not already exist.
func writeNewFile(path string, data []byte) error {
	handle, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600) //#nosec G304 -- the cache path is the caller's own configured path
	if err != nil {
		return err
	}
	// Removed on every failure past creation, or each failed attempt leaves a .tmp holding token bytes.
	if _, err := handle.Write(data); err != nil {
		_ = handle.Close()
		_ = os.Remove(path)
		return err
	}
	if err := handle.Close(); err != nil {
		_ = os.Remove(path)
		return err
	}
	return nil
}

func replaceFile(temporaryPath, destination string) error {
	wait := replaceRetryInitial
	for attempt := 0; ; attempt++ {
		err := os.Rename(temporaryPath, destination)
		if err == nil {
			return nil
		}
		if runtime.GOOS != "windows" || !isSharingDenial(err) || attempt >= replaceRetryCount-1 {
			return err
		}
		time.Sleep(wait)
		wait *= 2
	}
}

// isSharingDenial reports the Windows refusals a concurrent reader causes: access denied, sharing violation.
func isSharingDenial(err error) bool {
	if errors.Is(err, fs.ErrPermission) {
		return true
	}
	var errno syscall.Errno
	return errors.As(err, &errno) && errno == errorSharingViolation
}

// errorSharingViolation is Windows' ERROR_SHARING_VIOLATION.
const errorSharingViolation = syscall.Errno(32)

// fileScope is one outstanding level of the store's refresh lock.
type fileScope struct {
	store *FileTokenStore
	// generation is the hold this scope belongs to.
	generation uint64
	// released is atomic because Release need not run on the goroutine that took the scope.
	released atomic.Bool
}

func (f *fileScope) live(operation string) error {
	return scopeLiveness(operation, f.released.Load(), f.store.shared.holds(f.generation))
}

func (f *fileScope) ReadTokenRecord() (map[string]any, error) {
	if err := f.live("ReadTokenRecord"); err != nil {
		return nil, err
	}
	return f.store.readNamespaceRecord()
}

func (f *fileScope) WriteTokenRecord(record map[string]any) error {
	if err := f.live("WriteTokenRecord"); err != nil {
		return err
	}
	return f.store.writeNamespaceRecord(record)
}

func (f *fileScope) ClearTokenRecord() error {
	if err := f.live("ClearTokenRecord"); err != nil {
		return err
	}
	return f.store.clearNamespaceRecord()
}

// EnterNested adds a level to the lock this scope already holds.
func (f *fileScope) EnterNested() (TokenRefreshScope, error) {
	if err := f.live("EnterNested"); err != nil {
		return nil, err
	}
	return &fileScope{store: f.store, generation: f.store.shared.nest()}, nil
}

func (f *fileScope) Release() {
	if !f.released.CompareAndSwap(false, true) {
		return
	}
	f.store.shared.leave()
}
