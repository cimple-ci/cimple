cimple {
  version = "0.0.1"
}

name = "Cimple"
description = "Project description"

task "echo" {
  command "A" {
    command = "echo"
  }

  command "." {
    command = "echo"
  }

  command "" {
    command = "echo"
  }
}
