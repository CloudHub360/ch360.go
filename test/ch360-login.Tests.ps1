param(
  [String]$ClientId,
  [String]$ClientSecret
)

$configFolderPath = Join-Path -Path "~" -ChildPath ".ch360"
$configFolderPathBackup = "$configFolderPath" + "_backup"

function Format-MultilineOutput([Parameter(ValueFromPipeline=$true)]$input){
  $input -join [Environment]::NewLine
}

function Backup-ConfigFolder() {
  if (Test-Path $configFolderPath) {
    Remove-BackupFolder
    Copy-Folder $configFolderPath $configFolderPathBackup
    Remove-ConfigFolder
  }
}

function Restore-ConfigFolder() {
  if (!(Test-Path $configFolderPathBackup)) {
    return
  }

  Remove-ConfigFolder
  Copy-Folder $configFolderPathBackup $configFolderPath
  Remove-BackupFolder
}

function Remove-ConfigFolder() {
  if (Test-Path $configFolderPath) {
    Remove-Item $configFolderPath -Recurse -Force
  }
}

function Remove-BackupFolder() {
  if (Test-Path $configFolderPathBackup) {
    Remove-Item $configFolderPathBackup -Recurse -Force
  }
}

function Copy-Folder($source, $destination) {
  if (!(Test-Path $destination)) {
    New-Item -ItemType Directory $destination
  }
  Get-ChildItem -Path $source | Copy-Item -Destination $destination -Recurse -Container
}

Describe "ch360 --login" {
  BeforeEach {
    Backup-ConfigFolder
  }

  It "should output ok on success" {
    ch360 login --client-id="$ClientId" --client-secret="$ClientSecret" | Should -Be "Logging in... [OK]"
    $LASTEXITCODE | Should -Be 0
  }

  It "should write credentials to config file" {
    $expectedConfigFilePath = Join-Path -Path $configFolderPath -ChildPath "config.yaml"

    ch360 login --client-id="$ClientId" --client-secret="$ClientSecret"
    Get-Content -Path $expectedConfigFilePath | Format-MultilineOutput | Should -BeLike "*clientId: $ClientId*"
    Get-Content -Path $expectedConfigFilePath | Format-MultilineOutput | Should -BeLike "*clientSecret: $ClientSecret*"
  }

  AfterEach {
    Restore-ConfigFolder
  }
}