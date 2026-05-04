# AGENTS.md

This file provides guidance to AI coding assistants when working with this repository.

## Project Overview

AsyncRoutineManager is a lightweight framework for managing asynchronous operations in Go. It provides goroutine lifecycle management, metrics collection, and operational ID propagation for concurrent workloads.

## Build & Test Commands

```bash
make build          # Build the project
make test           # Run all tests (uses Ginkgo)
make generate       # Regenerate mocks
make clean          # Clean build artifacts
```

## Architecture

- **Root package**: Core `AsyncRoutineManager` type for goroutine lifecycle management
- **metrics/**: Prometheus metrics integration for routine tracking
- **opid/**: Operational ID propagation through goroutine contexts
- **linter/**: Custom linter rules for async routine usage

## Key Conventions

- Uses Ginkgo/Gomega for testing
- Mock generation via `mockgen`
- Module path: `github.com/openshift-online/async-routine`
