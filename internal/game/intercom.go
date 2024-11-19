package game

import "sync"

type entry struct {
	character string
}

type intercomService struct {
	data map[string][]entry
	mu   sync.RWMutex // protects data
}

func newIntercomService() *intercomService {
	return &intercomService{
		data: make(map[string][]entry),
	}
}

func (s *intercomService) Set(character string, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range s.data[name] {
		if e.character == character {
			return
		}
	}

	s.data[name] = append(s.data[name], entry{character: character})
}

func (s *intercomService) UnSet(character string, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if character == "" {
		delete(s.data, name)
		return
	}

	data, ok := s.data[name]
	if !ok {
		return
	}

	new := []entry{}
	for _, e := range data {
		if e.character != character {
			new = append(new, e)
		}
	}
	s.data[name] = new
}

func (s *intercomService) Get(character string, name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, ok := s.data[name]
	if !ok {
		return false
	}

	if character == "" {
		return len(data) > 0
	}

	for _, e := range data {
		if e.character == character {
			return true
		}
	}

	return false
}
