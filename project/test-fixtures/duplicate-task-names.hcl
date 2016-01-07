cimple {
  version = "0.0.1"
}

name = "Cimple"
description = "Testing duplicate task names exist"

task "echo" {
  description = "Description of the echo task"

  command "echo_hello_world" {
    command = "echo"
    args = ["hello world"]
  }
}

task "echo" {
  description = "Description of the echo task"

  command "echo_hello_world" {
    command = "echo"
    args = ["hello world"]
  }
}
