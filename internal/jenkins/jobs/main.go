package jobs

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Job struct {
	Name        string
	Description string
	Status      string
	Script      string
	IsDir       bool
	Children    []*Job
}

type JobDatabase struct {
	Root *Job
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
		instance = &JobDatabase{
			Root: &Job{},
		}
	})
	return instance
}

func (jdb *JobDatabase) ListJobs(jobPath string) []*Job {
	
	jobs := jdb._listJobs(jobPath, jdb.Root)

	jobs = deepCopy(jobs)
	for _, job := range jobs {
		job.Name = strings.ReplaceAll(job.Name, "/jobs", "")
		job.Script = ""
		job.Children = nil
	}
	return jobs
}

func (jdb *JobDatabase) _listJobs(prefix string, root *Job) []*Job {
	node := jdb._findSubtree(prefix, root)
	if node != nil && node.IsDir {
		return node.Children
	} else if node != nil && !node.IsDir {
		return []*Job{node}
	}
	return []*Job{}
}

func (jdb *JobDatabase) _findSubtree(prefix string, root *Job) *Job {
	if root.Name == prefix || strings.Contains(root.Name, prefix) {
		return root
	}
	
	for _, child := range root.Children {
		node := jdb._findSubtree(prefix, child)
		if node != nil {
			return node
		}
	}

	return nil
}

func (jdb *JobDatabase) LoadJobs() (*Job, error) {
	jenkinsHome := os.Getenv("JENKINS_HOME")
	jobDir := filepath.Join(jenkinsHome, "jobs")
	jobs, err := jdb._loadJobs(jobDir)
	if err != nil {
		return nil, err
	}
	jdb.Root = &Job{
		Name: "/jobs/",
		IsDir: true,
		Children: jobs,
	}
	
	return jdb.Root, nil
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
		job, err := parseRootBased(processedData, jobsDir, folder.Name())
		if err != nil {
			fmt.Printf("Error parsing config.xml for job %s: %v\n", folder.Name(), err)
		}

		job.Children, err = jdb._loadJobs(job.Name)
		job.Name = strings.TrimPrefix(job.Name, os.Getenv("JENKINS_HOME"))
		if err != nil {
			continue
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// Function to parse XML data dynamically
func parseRootBased(data []byte, jobDir string, folder string) (*Job, error) {
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
			Name:        filepath.Join(jobDir, folder),
			Description: workflow.Description,
			Status:      "pending",
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

func deepCopy(src []*Job) []*Job {
	dest := make([]*Job, len(src))
	for i, item := range src {
		dest[i] = &Job{
			Name:        item.Name,
			Description: item.Description,
			Status:      item.Status,
			Script:      item.Script,
			IsDir:       item.IsDir,
		}
	}
	return dest
}
