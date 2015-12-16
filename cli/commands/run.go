package cli

import (
	"io"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/build"
	"github.com/lukesmith/cimple/journal"
	"github.com/lukesmith/cimple/project"
	"github.com/lukesmith/cimple/vcs"
	"path"
	"path/filepath"
	"time"
)

func Run() cli.Command {
	return cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run Cimple against the current directory",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "task",
				Usage: "a specific task to run. Note that if the task is set to `skip` it will be run.",
			},
		},
		Action: func(c *cli.Context) {
			buildId := buildId()
			fileWriter, err := createOutputPathWriter(buildId)
			if err != nil {
				log.Fatal(err)
			}
			defer fileWriter.Close()

			logWriter := io.MultiWriter(os.Stdout, fileWriter)

			cfg, err := loadConfig()
			if err != nil {
				log.Fatal(err)
			}

			r := loadRepositoryInfo()

			journal, _ := createJournal(buildId)
			buildConfig := build.NewBuildConfig(buildId, logWriter, journal, cfg, *r)
			buildConfig.ExplicitTasks = c.StringSlice("task")

			err = executeBuild(buildConfig)
			if err != nil {
				log.Fatal(err)
			}
		},
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

func createJournal(buildId string) (journal.Journal, error) {
	journalWriter := journal.NewFileJournalWriter(journalPath(buildId))
	journal := journal.NewJournal(journalWriter)
	return journal, nil
}

func createOutputPathWriter(buildId string) (*os.File, error) {
	path := outputPath(buildId)
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	fileWriter, err := os.Create(outputPath(buildId))
	if err != nil {
		return nil, err
	}

	return fileWriter, nil
}

func journalPath(runId string) string {
	return path.Join(cimplePath(runId), "journal")
}

func outputPath(runId string) string {
	return path.Join(cimplePath(runId), "output")
}

func cimplePath(runId string) string {
	return path.Join(".", ".cimple", runId)
}
