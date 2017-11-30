$RootDir = Join-Path $PsScriptRoot ".."

function Build-WithVersion([DateTime]$Date, [String]$GitRev, [int]$BuildNo) {
  try {
    pushd $RootDir
    $dateStr = $Date.ToString("yyyy.MM.dd")
    $version="${dateStr}-${GitRev}:${BuildNo}"

    exec { go install -ldflags "-X github.com/CloudHub360/ch360.go.Version=$version" ./... }
  } finally {
    popd
  }
}

Task PackageRestore {
  try {
    pushd $RootDir
    exec { go get -t }
  } finally {
    popd
  }
}

Task Build PackageRestore, {
    pushd $RootDir
    $rev="$(git rev-parse --short HEAD)"
    $buildNum=([int]$env:BUILD_NUMBER)
    Build-WithVersion (Get-Date) $rev $buildNum
}

Task Test Build, {
  try {
    pushd $RootDir
    exec { go test -v -race ./... }

    $env:PATH += "$([Io.Path]::PathSeparator)$env:GOPATH/bin"
    $testResult = Invoke-Pester -PassThru
    assert ($testResult.FailedCount -eq 0)
  } finally {
    popd
  }
}

Task . PackageRestore, Build
