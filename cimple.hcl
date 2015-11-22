cimple {
  version = "0.0.1"
}

name = "Cimple"
description = "Cimple CI"

task "test" {
  description = "Run Cimple tests"

  env {
    GOPATH = "{{.WorkingDir}}/_vendor"
  }

  command "go" {
    command = "go"
    args = ["test", "-v", "./...", "-cover"]
  }
}
