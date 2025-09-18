I've analyzed the codebase and created a comprehensive CLAUDE.md file that covers:

**Development Commands:**
- Build: `make build`
- Test: `make test` (uses Ginkgo)
- Generate mocks: `make generate`
- Tools setup: `make tools`

**Architecture Overview:**
The AsyncRoutineManager is a Go framework for managing goroutines with monitoring capabilities. Key components include:

1. **AsyncRoutineManager**: Central routine tracking and observer management
2. **AsyncRoutine**: Managed goroutine interface with status tracking
3. **AsyncRoutineBuilder**: Fluent API for routine creation
4. **RoutinesObserver**: Observer pattern for lifecycle events
5. **AsyncRoutineMonitor**: Background monitoring system

**Key Technical Details:**
- Uses Ginkgo/Gomega for testing
- Prometheus integration for metrics
- Operation ID tracking via `opid` package
- Thread-safe concurrent maps
- Auto-generated mocks for testing

The CLAUDE.md provides future Claude Code instances with essential information about the build system, testing framework, and the high-level architecture needed to work effectively in this codebase.
