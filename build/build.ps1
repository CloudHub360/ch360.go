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
    exec { go install }
    mv -Force "$env:GOPATH\bin\ch360.go.exe" "$env:GOPATH\bin\ch360.exe"
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
