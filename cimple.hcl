cimple {
  version = "0.0.1"
}

name = "Cimple"
description = "Cimple CI build tasks"
version = "0.0.2-alpha"

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
  description = "Packages Cimple for release"
  depends = ["test", "fix"]

  script goxc {
    command = "goxc"
    body = <<BODY
goxc -build-ldflags \
  "-X main.VERSION={{index .Project.Version}} -X main.BuildDate={{index .FormattedBuildDate}} -X main.Revision={{index .Vcs.Revision}}" \
  -pv {{index .Project.Version}} \
  -br {{index .Vcs.Branch}} \
BODY
  }

  script build-cimple-docker {
    body = <<SCRIPT
docker build --build-arg CIMPLE_VERSION={{index .Project.Version}}-{{index .Vcs.Branch}} -t cimple -f Dockerfile .
SCRIPT
  }

  script tag-cimple-docker {
    body = <<SCRIPT
docker tag cimple cimpleci/cimple:latest
docker tag cimpleci/cimple:latest cimpleci/cimple:{{index .Project.Version}}
SCRIPT
  }
}

task publish {
  depends = ["package"]
  limit_to = "server"

  script publish-cimple-docker {
    body = <<SCRIPT
# docker push cimpleci/cimple:latest
# docker push cimpleci/cimple:{{index .Project.Version}}
echo pushing
SCRIPT
  }

  artifact "publish-binaries" {
    destination bintray {
      subject = "lukesmith"
      repository = "cimple"
      package = "cimple"
      username = "lukesmith"
    }
    file = "output/downloads/{{index .Project.Version}}/cimple_{{index .Project.Version}}_linux_386.tar.gz"
  }
}
