package workflow

import (
	"testing"

	"github.com/google/go-cmp/cmp"

)

func TestParseSingleStep(t *testing.T) {

	jenkinsfile := `
    pipeline {
        agent none
        stages {
            stage('Test') {
                steps {
                    sh 'node --version'
                }
            }
        }
    }
    `

	want := &Pipeline{
		Agent: &Agent{
			Docker: nil,
		},
		Stages: []*Stage{
			{
				Name: "Test",
				Steps: []*Step{
					{
						SingleKV: &SingleKVCommand{
							Command: "sh",
							Value:   "node --version",
						},
					},
				},
			},
		},
	}

	dslParser := DslParser{}
	pipeline, err := dslParser.Parse(jenkinsfile)
	if err != nil {
		t.Fatalf("Failed to parse Jenkinsfile: %v", err)
	}

	if diff := cmp.Diff(pipeline, want); diff != "" {
		t.Errorf("Structs are not equal (-got +want):\n%s", diff)
	}
}

func TestParseMultipleSteps(t *testing.T) {

	jenkinsfile := `
    pipeline {
		agent none
		stages {
			stage('Example Build') {
				agent { docker 'maven:3.9.3-eclipse-temurin-17' }
				steps {
					echo 'Hello, Maven'
					sh 'mvn --version'
					git branch: 'main',
						credentialsId: '12345-1234-4696-af25-123455',
						url: 'https://github.com/yegor86/tumbler-doll.git'
				}
			}
		}
	}
    `

	want := &Pipeline{
		Agent: &Agent{
			Docker: nil,
		},
		Stages: []*Stage{
			{
				Name: "Example Build",
				Agent: &Agent{
					Docker: &Docker{
						Image: "maven:3.9.3-eclipse-temurin-17",
					},
				},
				Steps: []*Step{
					{
						SingleKV: &SingleKVCommand{
							Command: "echo",
							Value:   "Hello, Maven",
						},
					},
					{
						SingleKV: &SingleKVCommand{
							Command: "sh",
							Value:   "mvn --version",
						},
					},
					{
						MultiKV: &MultiKVCommand{
							Command: "git",
							Params: []Param{
								{Key: "branch", Value: "main"},
								{Key: "credentialsId", Value: "12345-1234-4696-af25-123455"},
								{Key: "url", Value: "https://github.com/yegor86/tumbler-doll.git"},
							},
						},
					},
				},
			},
		},
	}

	dslParser := DslParser{}
	pipeline, err := dslParser.Parse(jenkinsfile)
	if err != nil {
		t.Fatalf("Failed to parse Jenkinsfile: %v", err)
	}

	if diff := cmp.Diff(pipeline, want); diff != "" {
		t.Errorf("Structs are not equal (-got +want):\n%s", diff)
	}
}

func TestParseParallelStage(t *testing.T) {

	jenkinsfile := `
    pipeline {
		agent none
		stages {
			stage('Example Build') {
				agent { docker 'maven:3.9.3-eclipse-temurin-17' }
				steps {
					sh 'mvn --version'
				}
			}
			stage('Parallel Stage') {
				failFast true
				parallel {
					stage('Branch A') {
						steps {
							echo 'On Branch A'
						}
					}
					stage('Branch B') {
						steps {
							echo 'On Branch B'
						}
					}
				}
			}
		}
	}
    `

	want := &Pipeline{
		Agent: &Agent{
			Docker: nil,
		},
		Stages: []*Stage{
			{
				Name: "Example Build",
				Agent: &Agent{
					Docker: &Docker{
						Image: "maven:3.9.3-eclipse-temurin-17",
					},
				},
				Steps: []*Step{
					{
						SingleKV: &SingleKVCommand{
							Command: "sh",
							Value:   "mvn --version",
						},
					},
				},
			},
			{
				Name: "Parallel Stage",
				FailFast: func(b bool) *bool { return &b }(true),
				Parallel: Parallel {
					{
						Name: "Branch A",
						Steps: []*Step{
							{
								SingleKV: &SingleKVCommand{
									Command: "echo",
									Value:   "On Branch A",
								},
							},
						},
					},
					{
						Name: "Branch B",
						Steps: []*Step{
							{
								SingleKV: &SingleKVCommand{
									Command: "echo",
									Value:   "On Branch B",
								},
							},
						},
					},
				},
			},
		},
	}

	dslParser := DslParser{}
	pipeline, err := dslParser.Parse(jenkinsfile)
	if err != nil {
		t.Fatalf("Failed to parse Jenkinsfile: %v", err)
	}

	if diff := cmp.Diff(pipeline, want); diff != "" {
		t.Errorf("Structs are not equal (-got +want):\n%s", diff)
	}
}
