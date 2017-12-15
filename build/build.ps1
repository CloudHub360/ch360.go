param(
  [String]$BuildDate = (Get-Date -Format "yyyy.MM.dd"),
  [String]$GitRev = "$(git rev-parse --short HEAD)",
  [String]$BuildNumber = "$([int]$env:BUILD_NUMBER)",
  [String]$ClientId,
  [String]$ClientSecret
)

$RootDir = Join-Path $PsScriptRoot ".."

Task PackageRestore {
  try {
    pushd $RootDir
    exec { go get -t -d ./... }
  } finally {
    popd
  }
}

Task GenerateMocks {
  try {
    pushd $RootDir
    exec { go get github.com/vektra/mockery/.../ }
    exec { go generate ./... }
  } finally {
    popd
  }
}

Task Build PackageRestore, {
  try {
    pushd $RootDir

    $version="${BuildDate}-${GitRev}-${BuildNumber}"

    exec { go install -ldflags "-X github.com/CloudHub360/ch360.go/ch360.Version=$version" ./... }
  } finally {
    popd
  }
}

function Register-CustomAssertions([string]$path) {
  Write-Host -ForegroundColor Green "Registering custom assertions from '$path'"
  Get-ChildItem (Resolve-Path "$PsScriptRoot/../$path") -Exclude "*.Tests.ps1" |% {
    $name = ($_.Name)
    try {
      . $_
      Write-Host -ForegroundColor DarkGreen "  Registered $name"
    } catch {
      Write-Host -ForegroundColor Red "  Failed to register $name"
    }
  }
}

Task Test Build, {
  try {
    pushd $RootDir
    exec { go test -v -race ./... }

    Register-CustomAssertions "test/assertions"

    $testResults = Invoke-Pester -PassThru -Script @{
      Path="test";
      Parameters = @{ClientId = $ClientId; ClientSecret = $ClientSecret}
    }
    assert ($testResults.FailedCount -eq 0)
  } finally {
    popd
  }
}

Task . PackageRestore, Build
