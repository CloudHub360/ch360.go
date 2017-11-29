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

}

Task . PackageRestore, Build
