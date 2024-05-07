package async

type Snapshot struct {
	totalRoutineCount int
	routineByName     map[string][]AsyncRoutine
	routineNames      []string
}

func newSnapshot() Snapshot {
	return Snapshot{
		routineByName: make(map[string][]AsyncRoutine),
	}
}

func (s *Snapshot) registerRoutine(r AsyncRoutine) {
	s.totalRoutineCount++
	if _, ok := s.routineByName[r.Name()]; !ok {
		s.routineNames = append(s.routineNames, r.Name())
	}
	s.routineByName[r.Name()] = append(s.routineByName[r.Name()], r)
}

func (s *Snapshot) GetTotalRoutineCount() int {
	return s.totalRoutineCount
}

func (s *Snapshot) GetRunningRoutinesNames() []string {
	return s.routineNames
}

func (s *Snapshot) GetRunningRoutinesCount(routineName string) int {
	return len(s.routineByName[routineName])
}
