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
    exec { go build -race -o "$env:GOPATH/bin/ch360" }
  } finally {
    popd
  }
}

Task Test {
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
