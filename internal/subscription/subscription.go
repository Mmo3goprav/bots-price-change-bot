package subscription

type Subscription struct {
	id     int64
	state  State
	charts map[string]struct{}
}

func NewSubscription(id int64) *Subscription {
	return &Subscription{
		id:     id,
		state:  StateNew,
		charts: make(map[string]struct{}),
	}
}

func (s *Subscription) GetState() State {
	return s.state
}

func (s *Subscription) SetState(state State) {
	s.state = state
}

func (s *Subscription) SatisfyChart(chart string) bool {
	_, ok := s.charts[chart]

	return ok
}

func (s *Subscription) AddChart(chart string) {
	s.charts[chart] = struct{}{}
}

func (s *Subscription) RemoveChart(chart string) {
	delete(s.charts, chart)
}

func (s *Subscription) GetID() int64 {
	return s.id
}
