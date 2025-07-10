package providers

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// FileProviderOptions holds configuration for the file provider.
type FileProviderOptions struct {
	// Path is the filesystem path to the JSON file storing flags.
	Path string
}

// FileProvider implements IFeatureFlagProvider using a JSON file.
type FileProvider struct {
	path  string
	flags map[string]bool
	mu    sync.RWMutex
}

// NewFileProvider creates a FileProvider, loading existing flags (or creating an empty file).
func NewFileProvider(opts FileProviderOptions) IFeatureFlagProvider {
	p := &FileProvider{
		path:  opts.Path,
		flags: make(map[string]bool),
	}

	// Ensure directory exists
	if dir := filepath.Dir(p.path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Fatalf("cannot create directories for feature‐flags file: %v", err)
		}
	}

	if err := p.load(); err != nil {
		log.Fatalf("failed to load feature‐flags from file: %v", err)
	}
	return p
}

func (p *FileProvider) GetListOfFeatureFlags(ctx context.Context) (map[string]bool, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	// Create a copy of the flags map to avoid race conditions
	flagsCopy := make(map[string]bool)
	for key, value := range p.flags {
		flagsCopy[key] = value
	}
	
	return flagsCopy, nil
}

func (p *FileProvider) load() error {
	data, err := os.ReadFile(p.path)
	if err != nil {
		if os.IsNotExist(err) {
			// write empty JSON object
			return p.save()
		}
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return json.Unmarshal(data, &p.flags)
}

func (p *FileProvider) save() error {
	p.mu.RLock()
	data, err := json.MarshalIndent(p.flags, "", "  ")
	p.mu.RUnlock()
	if err != nil {
		return err
	}
	// Write to temp and then atomically rename
	tmp := p.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, p.path)
}

// SetFeatureFlagStatus sets a flag in memory and persists the file.
func (p *FileProvider) SetFeatureFlagStatus(ctx context.Context, key string, value bool) error {
	p.mu.Lock()
	p.flags[key] = value
	p.mu.Unlock()
	if err := p.save(); err != nil {
		log.Printf("error saving feature flag file: %v", err)
		return err
	}
	return nil
}

// GetFeatureFlagStatus reads a flag from the in‐memory map.
func (p *FileProvider) GetFeatureFlagStatus(ctx context.Context, key string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.flags[key]
}
