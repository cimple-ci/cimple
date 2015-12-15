package cli

import (
	"io"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/lukesmith/cimple/build"
	"github.com/lukesmith/cimple/journal"
	"github.com/lukesmith/cimple/logging"
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
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:  "task",
				Usage: "a specific task to run. Note that if the task is set to `skip` it will be run.",
			},
		},
		Action: func(c *cli.Context) {
			runId := runId()
			fileWriter, err := createOutputPathWriter(runId)
			if err != nil {
				log.Fatal(err)
			}
			defer fileWriter.Close()

			logWriter := io.MultiWriter(os.Stdout, fileWriter)
			logger := logging.CreateLogger("cli", logWriter)

			cfg, err := loadConfig()
			if err != nil {
				log.Fatal(err)
			}

			skipNonSpecificTasks(logger, c.StringSlice("task"), cfg.Tasks)

			err = executeBuild(runId, logWriter, c, cfg)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
}

var loadConfig = func() (*project.Config, error) {
	return project.LoadConfig("cimple.hcl")
}

var executeBuild = func(runId string, out io.Writer, c *cli.Context, cfg *project.Config) error {
	journal, _ := createJournal(runId)
	buildConfig := build.NewBuildConfig(out, journal, &cfg.Project, cfg.Tasks)

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

func runId() string {
	return time.Now().Format(time.RFC3339)
}

func createJournal(runId string) (journal.Journal, error) {
	journalWriter := journal.NewFileJournalWriter(journalPath(runId))
	journal := journal.NewJournal(journalWriter)
	return journal, nil
}

func skipNonSpecificTasks(log *log.Logger, specificTasks []string, tasks map[string]*project.Task) {
	if len(specificTasks) != 0 {
		for _, t := range tasks {
			if contains(specificTasks, t.Name) {
				if t.Skip {
					log.Printf("Unskipping task %s as explicitly specified", t.Name)
					t.Skip = false
				}
			} else {
				log.Printf("Skipping task %s as not explicitly specified", t.Name)
				t.Skip = true
			}
		}
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
