# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

### Building and Testing
- `make build` - Build the application executable
- `make test` - Run all tests using Ginkgo test framework
- `make generate` - Generate mocks and other generated code
- `make clean` - Clean build artifacts

### Development Tools
- `make tools` - Install required development tools (ginkgo)
- `make mockgen-install` - Install mockgen for generating mocks

### Running Single Tests
- `ginkgo -r -v --focus="specific test pattern"` - Run specific tests matching pattern
- `go test -v ./...` - Alternative way to run all tests

## Architecture Overview

This is a Go library for managing asynchronous goroutines with enhanced visibility, monitoring, and lifecycle management. The core architecture consists of:

### Core Components

**AsyncRoutineManager** (`async_routine_manager.go`): Central singleton that manages all async routines. Provides registration/deregistration, observer pattern for notifications, and optional monitoring/snapshotting capabilities.

**AsyncRoutine** (`async_routine.go`): Represents a managed goroutine with metadata (name, creation time, status, timebox limits, operation IDs). Supports both regular functions and error group integration.

**AsyncRoutineBuilder** (`async_routine_builder.go`): Fluent builder pattern for creating routines with optional configuration (timebox, custom data, error groups).

### Key Patterns

**Observer Pattern**: `RoutinesObserver` interface allows pluggable monitoring (logging, metrics). Observers get notified on routine start/finish/timeout events.

**Operation ID Tracking** (`opid/`): Each routine gets a unique operation ID for distributed tracing, with support for originator context propagation.

**Prometheus Integration** (`metrics/`): Built-in metrics observer that tracks running routine counts and instances by name using Prometheus gauges.

### Usage Pattern
```go
async.NewAsyncRoutine("routine-name", ctx, func() {
    // your async work
}).Timebox(5*time.Minute).WithData("key", "value").Run()
```

The manager can be disabled globally, falling back to standard goroutines. Monitoring includes periodic snapshots and timebox violation detection.