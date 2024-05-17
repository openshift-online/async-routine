package async

type Snapshot struct {
	totalRoutineCount int
	routineByName     map[string][]AsyncRoutine
	routineNames      []string
	timedOutRoutines  []AsyncRoutine
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
	if r.hasExceededTimebox() {
		s.timedOutRoutines = append(s.timedOutRoutines, r)
	}
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

// GetTimedOutRoutines returns the list of routines that are still running and that have exceeded the configured
// time box
func (s *Snapshot) GetTimedOutRoutines() []AsyncRoutine {
	return s.timedOutRoutines
}
