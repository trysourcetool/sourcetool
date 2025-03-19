package sourcetool

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gofrs/uuid/v5"

	"github.com/trysourcetool/sourcetool-go/internal/logger"
)

type Sourcetool struct {
	Router
	apiKey   string
	endpoint string
	runtime  *runtime
	pages    map[uuid.UUID]*page
	mu       sync.RWMutex
}

func New(config *Config) *Sourcetool {
	hostParts := strings.Split(config.Endpoint, "://")
	if len(hostParts) != 2 {
		panic("invalid host")
	}
	namespaceDNS := strings.Split(hostParts[1], ":")[0]
	s := &Sourcetool{
		apiKey:   config.APIKey,
		endpoint: fmt.Sprintf("%s/ws", config.Endpoint),
		pages:    make(map[uuid.UUID]*page),
	}
	s.Router = newRouter(s, namespaceDNS)
	return s
}

func (s *Sourcetool) Listen() error {
	if err := s.validatePages(); err != nil {
		return err
	}

	if err := logger.Init(); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	s.mu.RLock()
	r, err := startRuntime(s.apiKey, s.endpoint, s.pages)
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	defer r.wsClient.Close()

	s.runtime = r

	return r.wsClient.Wait()
}

func (s *Sourcetool) Close() error {
	return s.runtime.wsClient.Close()
}

func (s *Sourcetool) validatePages() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pagesByRoute := make(map[string]uuid.UUID)
	for id, p := range s.pages {
		pagesByRoute[p.route] = id
	}

	newPages := make(map[uuid.UUID]*page)
	for _, id := range pagesByRoute {
		newPages[id] = s.pages[id]
	}
	s.pages = newPages

	return nil
}

func (s *Sourcetool) addPage(id uuid.UUID, p *page) {
	s.mu.Lock()
	s.pages[id] = p
	s.mu.Unlock()
}
