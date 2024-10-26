
## Prerequisites
### Install Temporal CLI

    brew install temporal

### Start Temporal dev server
This command automatically starts the Web UI, creates the default Namespace, and uses an in-memory database.
The Temporal Server should be available on localhost:7233 and the Temporal Web UI should be available at http://localhost:8233.

    temporal server start-dev

### Run API server

    go run main.go api

### Run Workflow executor with default TaskQueue `dsl`

    go run main.go wf

### Upload workflow [DSL sample](configs/workflow1.yaml)
