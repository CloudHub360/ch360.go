#Requires -Version 5.0

$global:ErrorActionPreference = "Stop"
$global:ProgressPreference = 'SilentlyContinue' # Hide progress bars

function global:Test-CommandVersion([string]$command, [string]$commandVersion) {
  $oldErrorPreference = $ErrorActionPreference
  $ErrorActionPreference = 'stop'
  try {
    return (& $command --version) -eq $commandVersion
  } catch {
    return $false
  } finally {
    $ErrorActionPreference = $oldErrorPreference
  }
}

function global:Install-PowerShellDependency([string]$ModuleName) {
  Write-Host "Testing for $ModuleName..."
  if ((Get-Module -ListAvailable $ModuleName) -eq $null) {
    Write-Host "Installing $ModuleName..."
    Install-Module -Force $ModuleName -Scope CurrentUser
    Get-Command -Module $ModuleName # This seems to set the necessary aliases?!
  }
}

function global:RestoreBuildLevelPackages() {
  try {
    Push-Location (Join-Path $PsScriptRoot "..")

    Install-PowerShellDependency "InvokeBuild"
    Install-PowerShellDependency "Pester"
  } finally {
    Pop-Location
  }
}

<#
.SYNOPSIS
Build.

.DESCRIPTION
This is really a wrapper around build.ps1 (build.ps1 is our actual build script.
2 main steps:

    1 - Restore build-time dependencies via paket.
    2 - Execute the build. (psake build.ps1)

In theory, Teamcity will also use this build command. Probably like this:
`build -Task Build`

.EXAMPLE
build
Run the build script with default values for all parameters

.EXAMPLE
build -Task Clean
Run the build script and execute only the 'Clean' task.
#>
function global:build() {
  [CmdletBinding()]
  param(
      # The Tasks to execute. An empty list runs the default task, as defined in build.ps1
      [Parameter(Position=0)]
      [string[]] $Tasks = @(),

      [Parameter(Position=1)]
      [string] $BuildNumber = '0'
  )

  RestoreBuildLevelPackages

  Invoke-Build `
      -File "build\build.ps1" `
      -Task $Tasks
}

Write-Host "This is the CloudHub360 Platform repo. And here are the available commands:" -Fore Magenta
Write-Host "`t build" -Fore Green
Write-Host "For more information about the commands, use Get-Help <command-name>" -Fore Magenta
Write-Host "To learn view the tasks exposed by each command, use <command-name> help" -Fore Magenta
