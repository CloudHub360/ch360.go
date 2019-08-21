function Invoke-App {
    $ErrorActionPreference = "Continue"
    try {
        Write-Information "Running: surf $args"
        & surf $args
    } catch [System.Management.Automation.RemoteException] {
        # Catch exceptions for messages redirected from stderr and
        # write out the messages to stdout
        Write-Output $Error[0].Message
    }
}

function Format-MultilineOutput([Parameter(ValueFromPipeline=$true)]$input){
    $input -join [Environment]::NewLine
}

function String-Starting([string]$input) {
    ([Regex]::Escape($input) + ".*")
}

$applicationFolderPath = Join-Path -Path "~" -ChildPath ".surf"
$applicationFolderPathBackup = "$applicationFolderPath" + "_backup"

function Backup-ApplicationFolder() {
    if (!(Test-Path $applicationFolderPath)) {
        return
    }

    Remove-Item $applicationFolderPathBackup -Recurse -Force -ErrorAction SilentlyContinue
    Copy-Folder $applicationFolderPath $applicationFolderPathBackup
    Remove-Item $applicationFolderPath -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "Backed up application folder"
}

function Restore-ApplicationFolder() {
    if (!(Test-Path $applicationFolderPathBackup)) {
        return
    }

    Remove-Item $applicationFolderPath -Recurse -Force -ErrorAction SilentlyContinue
    Copy-Folder $applicationFolderPathBackup $applicationFolderPath
    Remove-Item $applicationFolderPathBackup -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "Restored application folder"
}

function Copy-Folder($source, $destination) {
    if (!(Test-Path $destination)) {
        New-Item -ItemType Directory $destination
    }
    Get-ChildItem -Path $source | Copy-Item -Destination $destination -Recurse -Container
}

function Login-Surf() {
    Backup-ApplicationFolder

    $configFilePath = Join-Path -Path $applicationFolderPath -ChildPath "config.yaml"

    Invoke-App login --client-id="$ClientId" --client-secret="$ClientSecret" 2>&1 | Should -Be "Logging in... [OK]"
    $LASTEXITCODE | Should -Be 0

    Get-Content -Path $configFilePath | Format-MultilineOutput | Should -BeLike "*clientId: $ClientId*"
    Get-Content -Path $configFilePath | Format-MultilineOutput | Should -BeLike "*clientSecret: $ClientSecret*"
    Write-Host "Ran surf login"
}