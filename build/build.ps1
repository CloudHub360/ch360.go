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
    $today=Get-Date -format "yyyy.MM.dd"
    $rev="$(git rev-parse --short HEAD)"
    $buildNum=([int]$env:BUILD_NUMBER)
    $version="${today}-${rev}:${buildNum}"
    exec { go install -ldflags "-X github.com/CloudHub360/ch360.go.Version=$version" ./... }
  } finally {
    popd
  }
}

Task Test Build, {
  try {
    pushd $RootDir
    exec { go test -v -race ./... }

    $env:PATH += ";$($env:GOPATH)bin"
    Invoke-Pester
  } finally {
    popd
  }
}

Task . PackageRestore, Build
