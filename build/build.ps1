param(
  [String]$BuildDate = (Get-Date -Format "yyyy.MM.dd"),
  [String]$GitRev = "$(git rev-parse --short HEAD)",
  [String]$BuildNumber = "$([int]$env:BUILD_NUMBER)"
)

$RootDir = Join-Path $PsScriptRoot ".."

Task PackageRestore {
  try {
    pushd $RootDir
    exec { go get -t ./... }
  } finally {
    popd
  }
}

Task Gen {
  try {
    pushd $RootDir
    exec { go get github.com/vektra/mockery/.../ }
    exec { go generate }
  } finally {
    popd
  }
}

Task Build PackageRestore, {
  try {
    pushd $RootDir

    $version="${BuildDate}-${GitRev}:${BuildNumber}"

    exec { go install -ldflags "-X github.com/CloudHub360/ch360.go.Version=$version" ./... }
  } finally {
    popd
  }
}

Task Test Build, {
  try {
    pushd $RootDir
    exec { go test -v -race ./... }

    $env:PATH += "$([Io.Path]::PathSeparator)$env:GOPATH/bin"
    assert ((Invoke-Pester -PassThru).FailedCount -eq 0)
  } finally {
    popd
  }
}

Task . PackageRestore, Build
