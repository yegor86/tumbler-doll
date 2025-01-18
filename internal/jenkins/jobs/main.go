package jobs

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type Job struct {
	Name        string
	Description string
	Script      string
	IsDir       bool
}

type JobDatabase struct {
	jobs []*Job
}

type JobDefinition struct {
	XMLName   xml.Name `xml:"definition"`
	JobScript string   `xml:"script"`
}

type FolderDefinition struct {
	XMLName     xml.Name `xml:"com.cloudbees.hudson.plugins.folder.Folder"`
	Description string   `xml:"description"`
	DisplayName string   `xml:"displayName"`
}

type WorkflowDefinition struct {
	XMLName       xml.Name      `xml:"flow-definition"`
	Description   string        `xml:"description"`
	JobDefinition JobDefinition `xml:"definition"`
	Disabled      bool          `xml:"disabled"`
}

var (
	instance *JobDatabase
	once     sync.Once
)

func GetInstance() *JobDatabase {
	once.Do(func() {
		instance = &JobDatabase{}
	})
	return instance
}

func (jdb *JobDatabase) LoadJobs() ([]*Job, error) {
	jenkinsHome := os.Getenv("JENKINS_HOME")
	jobDir := filepath.Join(jenkinsHome, "jobs")
	jobs, err := jdb._loadJobs(jobDir)
	if err != nil {
		return nil, err
	}
	jdb.jobs = jobs
	return jobs, nil
}

func (jdb *JobDatabase) _loadJobs(jobsDir string) ([]*Job, error) {
	// List to store job details
	var jobs []*Job

	// Read all job folders
	folders, err := os.ReadDir(jobsDir)
	if err != nil {
		fmt.Printf("Error reading jobs directory: %v\n", err)
		return nil, err
	}

	for _, folder := range folders {
		if !folder.IsDir() {
			continue
		}

		// Path to the config.xml file for the job
		configPath := filepath.Join(jobsDir, folder.Name(), "config.xml")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			fmt.Printf("No config.xml found for job: %s\n", folder.Name())
			continue
		}

		configData, err := os.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Error reading config.xml for job %s: %v\n", folder.Name(), err)
			continue
		}

		processedData := stripXmlVersion(configData)
		job, err := parseRootBased(processedData, jobsDir)
		if err != nil {
			fmt.Printf("Error parsing config.xml for job %s: %v\n", folder.Name(), err)
		}

		if !job.IsDir {
			jobs = append(jobs, job)
			continue
		}
		children, err := jdb._loadJobs(job.Name)
		if err != nil {
			continue	
		}
		jobs = append(jobs, children...)
	}
	return jobs, nil
}

// Function to parse XML data dynamically
func parseRootBased(data []byte, jobDir string) (*Job, error) {
	var root struct {
		XMLName xml.Name
	}
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, err
	}

	switch root.XMLName.Local {
	case "flow-definition":
		var workflow WorkflowDefinition
		if err := xml.Unmarshal(data, &workflow); err != nil {
			return nil, fmt.Errorf("failed to parse workfow XML: %w", err)
		}
		return &Job{
			Name:        jobDir,
			Description: workflow.Description,
			Script:      workflow.JobDefinition.JobScript,
			IsDir:       false,
		}, nil
	case "com.cloudbees.hudson.plugins.folder.Folder":
		var folder FolderDefinition
		if err := xml.Unmarshal(data, &folder); err != nil {
			return nil, fmt.Errorf("failed to parse pipeline XML: %w", err)
		}
		return &Job{
			Name:        filepath.Join(jobDir, folder.DisplayName, "jobs"),
			Description: folder.Description,
			Script:      "",
			IsDir:       true,
		}, nil
	default:
		return nil, fmt.Errorf("unknown root element: %s", root.XMLName.Local)
	}
}

func stripXmlVersion(data []byte) []byte {
	res := regexp.
		MustCompile("(?m)^.*<?xml.*$").
		ReplaceAllString(string(data), "")
	return []byte(res)
}
