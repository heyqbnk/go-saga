package transaction

type cachedStep struct {
	step *Step[any]
	res  any
}

type stepStore struct {
	cache  []*cachedStep
	cursor int
}

// Size returns count of cached steps.
func (s *stepStore) Size() int {
	return len(s.cache)
}

// Store stores step in cache.
func (s *stepStore) Store(step *Step[any], res any) {
	s.cache = append(s.cache, &cachedStep{step: step, res: res})
	s.cursor = len(s.cache) - 1
}

func (s *stepStore) Pop() *cachedStep {
	if s.cursor == -1 {
		return nil
	}
	v := s.cache[s.cursor]
	s.cursor--
	return v
}

func newStepStore() *stepStore {
	return &stepStore{cursor: -1}
}
