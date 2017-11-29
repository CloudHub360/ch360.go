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
  } finally {
    popd
  }
}

Task Test {
  try {
    pushd $RootDir
    exec { go test -v -race ./... }
  } finally {
    popd
  }
}

Task . PackageRestore, Build
