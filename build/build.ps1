param(
  [String]$BuildDate = (Get-Date -Format "yyyy.MM.dd"),
  [String]$GitRev = "$(git rev-parse --short HEAD)",
  [String]$BuildNumber = "$([int]$env:BUILD_NUMBER)"
)

$RootDir = Join-Path $PsScriptRoot ".."

Task PackageRestore {
  try {
    pushd $RootDir
    exec { go get -t }
  } finally {
    popd
  }
}

Task Build PackageRestore, {
  try {
    pushd $RootDir
    
    $OS = $(go env GOOS)
    $arch = $(go env GOARCH)
    $suffix = $(go env GOEXE)

    $version="${BuildDate}-${GitRev}:${BuildNumber}"
    $outputFile = Join-Path $env:GOPATH -ChildPath "bin" | Join-Path -ChildPath "${OS}-${arch}" | Join-Path -ChildPath "ch360${suffix}"

    exec { go build -ldflags "-X github.com/CloudHub360/ch360.go.Version=$version"  -o $outputFile ./cmd/ch360 }
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
