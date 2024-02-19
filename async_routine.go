package async

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"gitlab.cee.redhat.com/service/uhc-clusters-service/pkg/opid"
)

type RoutineStatus string

const (
	RoutineStatusCreated                 RoutineStatus = "created"
	RoutineStatusRunning                 RoutineStatus = "running"
	RoutineStatusExceededTimebox         RoutineStatus = "running_exceeded_timebox"
	RoutineStatusFinished                RoutineStatus = "finished"
	RoutineStatusFinishedExceededTimebox RoutineStatus = "finished_exceeded_timebox"
)

var _ AsyncRoutine = (*asyncRoutine)(nil)

// AsyncRoutine is the public interface for managed async routines
type AsyncRoutine interface {
	Name() string
	CreatedAt() time.Time
	StartedAt() *time.Time
	FinishedAt() *time.Time
	Status() RoutineStatus
	OpId() string
	OriginatorOpId() string
}

type asyncRoutine struct {
	name             string
	routine          func()
	routineWithError func() error
	createdAt        time.Time
	startedAt        *time.Time
	finishedAt       *time.Time
	timebox          *time.Duration
	errGroup         *errgroup.Group
	status           RoutineStatus
	ctx              context.Context
	originatorOpId   string
}

func (r *asyncRoutine) Name() string {
	return r.name
}

func (r *asyncRoutine) CreatedAt() time.Time {
	return r.createdAt
}

func (r *asyncRoutine) StartedAt() *time.Time {
	return r.startedAt
}

func (r *asyncRoutine) FinishedAt() *time.Time {
	return r.finishedAt
}

func (r *asyncRoutine) Status() RoutineStatus {
	return r.status
}

func (r *asyncRoutine) OpId() string {
	return opid.FromContext(r.ctx)
}

func (r *asyncRoutine) OriginatorOpId() string {
	return r.originatorOpId
}

func (r *asyncRoutine) isStarted() bool {
	return r.startedAt != nil
}

func (r *asyncRoutine) isFinished() bool {
	return r.finishedAt != nil
}

func (r *asyncRoutine) isRunning() bool {
	return r.isStarted() && !r.isFinished()
}

func (r *asyncRoutine) hasExceededTimebox() bool {
	return r.isRunning() && r.timebox != nil && time.Now().UTC().After(r.startedAt.Add(*r.timebox))
}

func (r *asyncRoutine) run() {
	if r.isStarted() {
		// already running
		return
	}
	now := time.Now().UTC()
	r.startedAt = &now
	r.status = RoutineStatusRunning

	updateFinishedRoutine := func(r *asyncRoutine) {
		finishedAt := time.Now().UTC()
		r.finishedAt = &finishedAt

		if r.status == RoutineStatusExceededTimebox {
			r.status = RoutineStatusFinishedExceededTimebox
		} else {
			r.status = RoutineStatusFinished
		}
	}

	manager.notify(func(observer RoutinesObserver) {
		observer.RoutineStarted(r)
	})

	if r.errGroup != nil {
		r.errGroup.Go(func() error {
			err := r.routineWithError()
			updateFinishedRoutine(r)
			return err
		},
		)
		return
	}

	go func() {
		r.routine()
		updateFinishedRoutine(r)
	}()
}
