package project

import (
	"bytes"
	"io"
	"log"
	exec "os/exec"
	"text/template"
)

func Run(config *Config, stdout io.Writer, stderr io.Writer) {
	for _, task := range config.Tasks {
		executeTask(task, stdout, stderr)
	}
}

func executeTask(task Task, stdout io.Writer, stderr io.Writer) {
	for _, cmd := range task.Commands {
		if cmd.Skip {
			log.Printf("Skipping task %s", task.Name)
			continue
		}
		executeCommand(task.Env, cmd, stdout, stderr)
	}
}

func executeCommand(taskEnvs map[string]string, command Command, stdout io.Writer, stderr io.Writer) {
	var cmd = exec.Command(command.Command, command.Args...)

	// Clear out env for command so not to inherit current process's environment.
	cmd.Env = []string{}

	envs := make(map[string]string)

	for k, v := range taskEnvs {
		envs[k] = v
	}

	for k, v := range command.Env {
		envs[k] = v
	}

	for k, v := range envs {
		// TODO: Extract templating environment variables
		tmpl, err := template.New("t").Parse(k + "=" + v)
		if err != nil {
			panic(err)
		}
		vari, err := GetProjectVariables()
		var doc bytes.Buffer
		tmpl.Execute(&doc, vari)
		cmd.Env = append(cmd.Env, doc.String())
	}

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Run()
}
