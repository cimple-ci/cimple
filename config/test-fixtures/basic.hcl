cimple {
  version = "0.0.1"
}

name = "Cimple"
description = "Project description"

task "echo" {
  description = "Description of the echo task"

  command "echo_hello_world" {
    condition {
      // Command will only run on local build locations
      attribute = "$attr.build.location"
      value = "local"
    }

    command = "echo"
    args = ["hello world"]
  }

  command "echo" {
    command = "echo"
    args = ["moo >> cow.txt"]
  }

  command "cat" {
    command = "cat"
    args = ["cow.txt"]
  }

  archive = ["cow.txt"]
}
