package main

import (
	"fmt"
	"log"

	temporal "go.temporal.io/sdk/client"

	"github.com/yegor86/tumbler-doll/cmd"
	"github.com/yegor86/tumbler-doll/internal/cryptography"
	"github.com/yegor86/tumbler-doll/internal/env"
	"github.com/yegor86/tumbler-doll/internal/jenkins/jobs"
)

func main() {
	env.LoadEnvVars()
	crypto := cryptography.GetInstance()
	crypto.LoadOrSeedCrypto()

	wfClient, err := temporal.Dial(temporal.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalf("Unable to create Workflow client: %v", err)
	}
	defer wfClient.Close()

	jobDb := jobs.GetInstance()
	jobs, err := jobDb.LoadJobs()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Jobs: %v", jobs)

	cmd.Execute(wfClient)
}