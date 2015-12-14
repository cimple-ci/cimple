package cli

import (
	"io"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/build"
	"github.com/lukesmith/cimple/journal"
	"github.com/lukesmith/cimple/project"
	"path"
	"path/filepath"
	"time"
)

func Run() cli.Command {
	return cli.Command{
		Name:    "run",
		Aliases: []string{"r"},
		Usage:   "Run Cimple against the current directory",
		Action: func(c *cli.Context) {
			cfg, err := project.LoadConfig("cimple.hcl")
			if err != nil {
				log.Fatal(err)
				panic(err)
			}

			runId := time.Now().Format(time.RFC3339)

			journalWriter := journal.NewFileJournalWriter(journalPath(runId))
			journal := journal.NewJournal(journalWriter)

			fileWriter, err := createOutputPathWriter(runId)
			if err != nil {
				log.Fatal(err)
			}
			defer fileWriter.Close()

			logWriter := io.MultiWriter(os.Stdout, fileWriter)
			buildConfig := build.NewBuildConfig(logWriter, journal, &cfg.Project, cfg.Tasks)

			build, err := build.NewBuild(buildConfig)
			if err != nil {
				log.Fatal(err)
			}
			build.Run()
		},
	}
}

func createOutputPathWriter(runId string) (*os.File, error) {
	path := outputPath(runId)
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	fileWriter, err := os.Create(outputPath(runId))
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
