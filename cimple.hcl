cimple {
  version = "0.0.1"
}

name = "Cimple"
description = "Cimple CI build tasks"
version = "0.0.1"

env {
  GOPATH = "{{index .HostEnv \"GOPATH\"}}"
  GOROOT = "{{index .HostEnv \"GOROOT\"}}"
  # PATH required for glide command (needs access to Git).
  # Should PATH always be mapped by default?
  PATH = "{{index .HostEnv \"PATH\"}}"
}

task fix {
  script gofmt {
    body = "go fmt $(go list ./... | grep -v /vendor/)"
  }

  script govet {
    body = "go vet $(go list ./... | grep -v /vendor/)"
  }
}

task test {
  description = <<DESC
Runs the Cimple tests.

Prior to running the tests dependencies are installed.
DESC

  command glide {
    command = "glide"
    args = ["install"]
  }

  script gotest {
    body = "go test -v $(go list ./... | grep -v /vendor/) -cover"
  }
}

task package {
  description = <<DESC
Packages Cimple for release
DESC

  skip = true

  command goxc {
    command = "goxc"
    args = ["-pv", "{{index .Project.Version}}", "-br", "{{index .Vcs.Branch}}"]
  }

  script build-cimple-docker {
    body = <<SCRIPT
docker build -t cimple-base -f Dockerfile .
docker build -t cimple-server -f server.Dockerfile .
docker build -t cimple-agent -f agent.Dockerfile .
SCRIPT
  }
}
