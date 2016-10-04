# Cimple

## Getting started

1. Install [Glide](https://github.com/Masterminds/glide)
2. Run `glide install`

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

##### Accessing variables

Cimple makes a number of environment variables available to the scripts that are run. These
are prefixed with `CIMPLE_`.

- `CIMPLE_VERSION` - the version of Cimple
- `CIMPLE_PROJECT_NAME` - the value specified by the `name` field
- `CIMPLE_PROJECT_VERSION` - the value specified by the `version` field
- `CIMPLE_WORKING_DIR` - the working directory scripts are executed within
- `CIMPLE_TASK_NAME` - the name of the task being executed
- `CIMPLE_VCS` - the name of the vcs
- `CIMPLE_VCS_BRANCH` - the name of the vcs branch
- `CIMPLE_VCS_REVISION` - the revision of the vcs
- `CIMPLE_VCS_REMOTE_URL` - the url of the vcs remote
- `CIMPLE_VCS_REMOTE_NAME` - the name of the vcs remote

These values are also accessible within the `cimple.hcl` file using go templating. These
are accessed removing the `CIMPLE_` and replacing the `_` in the environment variable name
with `.`.

Example:

```hcl
script example {
  body = "echo {{index .Project.Name}}"
}
```

### Running a Server/Agent

The Cimple CLI can be run in either Server mode or Agent mode.

```shell
cimple server
cimple agent
```

#### Running in Docker

To run the server:

```
docker run -p 8080:8080 -p 1514:1514 cimple-server
```

To run the agent:

```
docker run cimple-agent --server-addr 192.168.99.100
```

During development you can use `scripts/cimple-server.sh` and `scripts/cimple-agent.sh`.

#### Triggering builds

When running in Server/Agent mode the Server will schedule tasks across the available agent pool.

By design the Server does not poll or listen to hooks from SCM systems. Instead external processes send builds into the
Cimple Server.

See `scripts/trigger_build.sh` for an example of how a build can be pushed into Cimple Server using curl. Or use the client api

```shell
cimple builds submit https://github.com/cimple-ci/cimple-ruby-example.git master
```

### Help

To get help on the Cimple commands run `cimple --help`.
