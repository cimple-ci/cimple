cimple {
  version = "0.0.1"
}

name = "Cimple"
description = "Project description"

env {
  project_env = "project"
}

task "echo" {
  description = "Description of the echo task"
  skip = true

  env {
    task_env = "global"
  }

  command "echo_hello_world" {
    command = "echo"
    args = ["hello world"]
  }

  command "echo" {
    command = "echo"
    args = ["moo >> cow.txt"]
    skip = true
  }

  command "cat" {
    command = "cat"
    args = ["cow.txt"]

    env {
      env = "test"
    }
  }

  archive = ["cow.txt"]
}
