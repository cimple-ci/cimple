# Cimple

## Getting started

1. Install [Glide](https://github.com/Masterminds/glide)
2. Run `source bootstrap.sh`
3. Run `glide install`

To run Cimple against itself run `go run main.go run`.

### Building Cimple

To build a release binary of the Cimple command line run `make build`

## Documentation

### Project Configuration

Projects are configured using `hcl`. A configuration is made up of several parts

- cimple metadata
- project information
- tasks
- steps

#### Cimple metadata

The Cimple metadata consists of the following

- version - The version of the Cimple schema *required*

Example:

```hcl
cimple {
  version = "0.0.1"
}
```

#### Project information

- name - The name of the project
- description - A description of the project

Example:

```hcl
name = "Cimple"
description = "Cimple CI build tasks"
```

#### Tasks

Tasks are a group of related steps that run together. There can be one or more steps
within a task.

Task names must be unique within a configuration. They can only contain the following
characters: `a-z`, `_`, `0-9`.

```hcl
task task_name {
}
```

#### Steps

Steps specify what should happen. There are two ways to specify a step:

##### Command steps

Command steps will run a specific command with optional arguments.

Step names must be unique within their parent task. They can only contain the following
characters: `a-z`, `_`, `0-9`.

Example:

```hcl
command glide {
  command = "glide"
  args = ["install"]
}
```

##### Script steps

In cases where more complex steps need to be performed the script step can be used to
specify a shell script to be run.

Example:

```hcl
script gotest {
  body = "go test -v $(go list ./... | grep -v /vendor/) -cover"
}
```

### Help

To get help on the Cimple commands run `cimple --help`.
