cimple {
  version = "0.0.1"
}

name = "Cimple"
description = ""

task "echo" {
  description = "Description of the echo task"

  command "echo" {
    condition {
      // Command will only run on local build locations
      attribute = "$attr.build.location"
      value = "local"
    }

    args = ["hello world"]
  }

  command "echo" {
    args = ["moo >> cow.txt"]
  }

  command "cat" {
    args = ["cow.txt"]
  }

  // Will archive the paths specified, making available to subsequent tasks
  archive = ["cow.txt"]
}

task "build" {
  // Specifies the image to use for all commands
  base_image = "technosophos/goglide:1.5-0-onbuild"

  command "go" {
    description = "Build Cimple"
    args = ["build", "cmd/cli/main.go"]
  }
}

task "docs" {
  command "godoc" {
  }
}

task "noop" {
  // Will skip all commands in the task
  skip = true

  command "noop" {
    // Will skip the command
    skip = true
  }
}

task "release" {
  // This task will only run if "build" completed successfully
  requires_tasks = ["build"]

  condition {
    // The task will only run if the source branch is master
    attribute = "$attr.source.branch"
    value = "master"
  }

  condition {
    // The task will only run if the build is performed via the server
    attribute = "$attr.build.location"
    value = "server"
  }

  command "upload" {
  }
}
