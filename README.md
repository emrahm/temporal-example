# Temporal Workflow Testing Project

This project is for testing Temporal workflows, learning about scheduling, and running a single Temporal workflow.

## Prerequisites

- [Temporal CLI](https://docs.temporal.io/cli) (or Docker)
- [Go](https://golang.org/dl/) installed
- Temporal server running locally

## Setup Instructions

### 1. Start Temporal Server

You need to run a Temporal server before executing workflows. Choose one of these methods:

**Using Docker:**

```bash
docker run --rm -it --network=host temporalio/auto-setup:1.22.2
```

Using Temporal CLI:

```bash
temporal server start-dev
```

This will start Temporal with a default namespace default.

### 2. Run the Worker

Workflows won't execute without a worker running. Start the worker with:

```bash
go run main.go run-worker
```

### 2. Trigger Workflows

You can trigger workflows using these commands:

Run a single workflow:
**Run a single workflow:**

```bash
go run main.go run delivery "test message"
```

**Schedule a recurring workflow:**

```bash
go run main.go schedule delivery 1001 "test message" --interval 5s
```

## VS Code Debug Configuration

To run/debug directly from VS Code, use these configurations in your .vscode/launch.json:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Run Worker",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceRoot}/main.go",
      "args": ["run-worker"]
    },
    {
      "name": "Trigger",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceRoot}/main.go",
      "args": ["run", "delivery", "test message"]
    },
    {
      "name": "Schedule",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceRoot}/main.go",
      "args": [
        "schedule",
        "delivery",
        "1001",
        "test message",
        "--interval",
        "5s"
      ]
    }
  ]
}
```

## Notes

**The worker must be running before triggering any workflows**
**Schedules are identified by ID (1001 in the example)**
**Interval format supports seconds (s), minutes (m), hours (h), etc.**
