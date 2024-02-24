package async

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"gitlab.cee.redhat.com/service/uhc-clusters-service/pkg/opid"
)

var _ AsyncRoutineBuilder = (*asyncRoutineBuilder)(nil)

type asyncRoutineBuilder struct {
	asyncRoutine asyncRoutine
}

// AsyncRoutineBuilder builds a managed routine
type AsyncRoutineBuilder interface {
	// Timebox assigns a timebox to the routine. If the timebox is reached the routine is moved to warning
	// state. This value is optional.
	Timebox(duration time.Duration) AsyncRoutineBuilder

	Build() AsyncRoutine

	// Run runs the routine
	Run(manager AsyncRoutineManager)
}

// NewAsyncRoutine instantiates a new AsyncRoutineBuilder
// name is the name to assign to the routine
// routine is the function to be executed asynchronously
func NewAsyncRoutine(
	name string,
	ctx context.Context,
	routine func()) AsyncRoutineBuilder {
	return &asyncRoutineBuilder{
		asyncRoutine: asyncRoutine{
			name:           name,
			routine:        routine,
			createdAt:      time.Now().UTC(),
			status:         RoutineStatusCreated,
			ctx:            opid.NewContext(),
			originatorOpId: opid.FromContext(ctx),
		},
	}
}

func NewAsyncRoutineWithErrGroup(
	name string,
	ctx context.Context,
	errGroup *errgroup.Group,
	routine func() error) AsyncRoutineBuilder {
	return &asyncRoutineBuilder{
		asyncRoutine: asyncRoutine{
			name:             name,
			routineWithError: routine,
			errGroup:         errGroup,
			createdAt:        time.Now().UTC(),
			status:           RoutineStatusCreated,
			ctx:              opid.NewContext(),
			originatorOpId:   opid.FromContext(ctx),
		},
	}
}

func (b *asyncRoutineBuilder) Timebox(duration time.Duration) AsyncRoutineBuilder {
	b.asyncRoutine.timebox = &duration
	return b
}

func (b *asyncRoutineBuilder) Run(manager AsyncRoutineManager) {
	manager.Run(&b.asyncRoutine)
}

func (b *asyncRoutineBuilder) Build() AsyncRoutine {
	return &b.asyncRoutine
}
