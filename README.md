
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

### JENKINS file structure
```sh
JENKINS_HOME
 +- config.xml     (jenkins root configuration)
 +- *.xml          (other site-wide configuration files)
 +- userContent    (files in this directory will be served under your http://server/userContent/)
 +- fingerprints   (stores fingerprint records)
 +- nodes          (slave configurations)
 +- plugins        (stores plugins)
 +- secrets        (secretes needed when migrating credentials to other servers)
 +- workspace (working directory for the version control system)
     +- [JOBNAME] (sub directory for each job)
 +- jobs
     +- [JOBNAME]      (sub directory for each job)
         +- config.xml     (job configuration file)
         +- latest         (symbolic link to the last successful build)
         +- builds
             +- [BUILD_ID]     (for each build)
                 +- build.xml      (build result summary)
                 +- log            (log file)
                 +- changelog.xml  (change log)
```