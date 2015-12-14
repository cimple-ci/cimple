cimple {
  version = "0.0.1"
}

name = "Cimple"
description = "Cimple CI"

task "fix" {
  env {
    GOPATH = "{{.WorkingDir}}/_vendor"
  }

  command "gofmt" {
    command = "go"
    args = ["fmt", "./..."]
  }

  command "govet" {
    command = "go"
    args = ["vet", "./..."]
  }
}

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
