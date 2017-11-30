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

    if ($env:GOOS -eq $null) {
      throw "Environment variable GOOS is not set. Possible values are listed here: https://golang.org/doc/install/source#environment"
    }

    if ($env:GOARCH -eq $null) {
      throw "Environment variable GOARCH is not set. Possible values are listed here: https://golang.org/doc/install/source#environment"
    }

    $version="${BuildDate}-${GitRev}:${BuildNumber}"
    $outputDir = "../../../../bin/$env:GOOS-$env:GOARCH"
    
    if ($env:GOOS -eq "windows") {
      $outputFile = Join-Path $outputDir "ch360.exe"
    }
    else {
      $outputFile = Join-Path $outputDir "ch360"
    }

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
