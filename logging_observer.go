package async

import (
	"context"

	"github.com/openshift-online/ocm-sdk-go/logging"
)

var _ RoutinesObserver = &loggingObserver{}

type loggingObserver struct {
	logging.Logger
	ctx context.Context
}

func (l *loggingObserver) RoutineStarted(routine AsyncRoutine) {
	l.Logger.Info(
		l.ctx,
		"New go-routine created, original opid='%s', new routine opid='%s', name='%s'",
		routine.OriginatorOpId(),
		routine.OpId(),
		routine.Name())
}

func (l *loggingObserver) RoutineFinished(routine AsyncRoutine) {
	l.Logger.Info(
		l.ctx,
		"Routine finished, routine opid='%s', name='%s', started=%v, ended=%s, status=%s]",
		routine.OpId(), routine.Name(), routine.StartedAt(), routine.FinishedAt(), routine.Status())
}

func (l *loggingObserver) RoutineExceededTimebox(routine AsyncRoutine) {
	l.Logger.Warn(l.ctx, "Routine exceeded timebox, routine opid='%s', name='%s', started=%v",
		routine.OpId(), routine.Name(), routine.StartedAt())
}

func (l *loggingObserver) RunningRoutineCount(count int) {
	l.Logger.Info(l.ctx, "Current managed routines count: %d", count)
}

func (l *loggingObserver) RunningRoutineByNameCount(name string, count int) {
	l.Logger.Info(l.ctx, "Current managed routines count (by name): name='%s', count=%d", name, count)
}

func NewLoggingObserver(ctx context.Context, logger logging.Logger) RoutinesObserver {
	return &loggingObserver{
		logger,
		ctx,
	}
}
