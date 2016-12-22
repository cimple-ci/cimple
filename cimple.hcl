cimple {
  version = "0.0.1"
}

name = "Cimple"
description = "Cimple CI build tasks"
version = "0.0.4"

env {
  GOPATH = "{{index .HostEnv \"GOPATH\"}}"
  GOROOT = "{{index .HostEnv \"GOROOT\"}}"
  # PATH required for glide command (needs access to Git).
  # Should PATH always be mapped by default?
  PATH = "{{index .HostEnv \"PATH\"}}"
  VERSION_LABEL = "{{if ne (index .Vcs.Branch) \"master\"}}{{index .Vcs.Branch}}-{{index .Vcs.Revision}}{{end}}"
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

  command clean {
    command = "rm"
    args = ["-rf", "output"]
  }

  script goxc {
    body = <<BODY
goxc -build-ldflags \
  "-X main.VERSION={{index .Project.Version}} -X main.BuildDate={{index .FormattedBuildDate}} -X main.Revision={{index .Vcs.Revision}}" \
  -pv {{index .Project.Version}} \
  -br "$VERSION_LABEL" \
BODY
  }

  script build-cimple-docker {
    body = <<SCRIPT
{{ if ne (index .Vcs.Branch) "master" }}
docker build --build-arg CIMPLE_VERSION={{index .Project.Version}}-$VERSION_LABEL -t cimple -f Dockerfile .
{{ else }}
docker build --build-arg CIMPLE_VERSION={{index .Project.Version}} -t cimple -f Dockerfile .
{{ end }}
SCRIPT
  }

  script tag-cimple-docker {
    body = <<SCRIPT
{{ if ne (index .Vcs.Branch) "master" }}
docker tag cimpleci/cimple:latest cimpleci/cimple:{{index .Project.Version}}-$VERSION_LABEL
{{ else }}
docker tag cimple cimpleci/cimple:latest
docker tag cimpleci/cimple:latest cimpleci/cimple:{{index .Project.Version}}
{{ end }}
SCRIPT
  }
}

task publish {
  depends = ["package"]
  limit_to = "server"

  script publish-cimple-docker {
    body = <<SCRIPT
{{ if ne (index .Vcs.Branch) "master" }}
docker push cimpleci/cimple:{{index .Project.Version}}-$VERSION_LABEL
{{ else }}
docker push cimpleci/cimple:latest
docker push cimpleci/cimple:{{index .Project.Version}}
{{ end }}
SCRIPT
  }

  publish binaries {
    destination bintray {
      subject = "cimpleci"
      repository = "pkgs"
      package = "cimple"
      username = "lukesmith"
    }
    files = [
      "output/downloads/{{index .Project.Version}}*/cimple_{{index .Project.Version}}*.tar.gz",
      "output/downloads/{{index .Project.Version}}*/cimple_{{index .Project.Version}}*.zip"
    ]
  }

  publish deb-packages {
    destination bintray {
      subject = "cimpleci"
      repository = "debian"
      package = "cimple"
      username = "lukesmith"
    }
    files = [
      "output/downloads/{{index .Project.Version}}*/cimple_{{index .Project.Version}}*.deb"
    ]
  }
}
