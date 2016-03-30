package runner

import (
	"io"
	"log"
	"os"

	"github.com/lukesmith/cimple/build"
	"github.com/lukesmith/cimple/journal"
	"github.com/lukesmith/cimple/project"
	"github.com/lukesmith/cimple/vcs"
	"github.com/lukesmith/syslog"
	"path"
	"path/filepath"
	"time"
)

type RunOptions struct {
	LogServer string
}

func Run(options *RunOptions, explicitTasks []string) {
	buildId := buildId()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	projectName := cfg.Project.Name

	fileWriter, err := createOutputPathWriter(projectName, buildId)
	if err != nil {
		log.Fatal(err)
	}
	defer fileWriter.Close()

	r := loadRepositoryInfo()

	writers := []io.Writer{os.Stderr, fileWriter}

	journalFileWriter := journal.NewFileJournalWriter(journalPath(projectName, buildId))
	journalWriters := []journal.JournalWriter{journalFileWriter}

	if options.LogServer != "" {
		log.Print("Connecting to syslog")
		s, err := syslog.Dial("tcp", options.LogServer, syslog.LOG_INFO, "Runner", nil)
		if err != nil {
			log.Print("Failed to connect to syslog server")
		}
		defer s.Close()

		writers = append(writers, s)
		journalWriters = append(journalWriters, journal.NewSyslogWriter(s))
	}

	logWriter := io.MultiWriter(writers...)

	journal := journal.NewJournal(journalWriters)
	buildConfig := build.NewBuildConfig(buildId, logWriter, journal, cfg, *r)
	buildConfig.ExplicitTasks = explicitTasks

	err = executeBuild(buildConfig)
	if err != nil {
		log.Fatal(err)
	}
}

var loadRepositoryInfo = func() *vcs.VcsInformation {
	r, err := vcs.LoadVcsInformation()
	if err != nil {
		log.Fatal(err)
	}

	return r
}

var loadConfig = func() (*project.Config, error) {
	return project.LoadConfig("cimple.hcl")
}

var executeBuild = func(buildConfig *build.BuildConfig) error {
	build, err := build.NewBuild(buildConfig)
	if err != nil {
		return err
	}

	err = build.Run()
	if err != nil {
		return err
	}

	return nil
}

func buildId() string {
	return time.Now().Format(time.RFC3339)
}

func createOutputPathWriter(projectName string, buildId string) (*os.File, error) {
	path := outputPath(projectName, buildId)
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	fileWriter, err := os.Create(outputPath(projectName, buildId))
	if err != nil {
		return nil, err
	}

	return fileWriter, nil
}

func journalPath(projectName string, runId string) string {
	return path.Join(cimplePath(projectName, runId), "journal")
}

func outputPath(projectName string, runId string) string {
	return path.Join(cimplePath(projectName, runId), "output")
}

func cimplePath(projectName string, runId string) string {
	return path.Join(".", ".cimple", projectName, runId)
}
