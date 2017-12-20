param(
    [String]$ClientId,
    [String]$ClientSecret
)

$classifierName = "test-classifier"
$applicationFolderPath = Join-Path -Path "~" -ChildPath ".ch360"
$applicationFolderPathBackup = "$applicationFolderPath" + "_backup"
$configFilePath = Join-Path -Path $applicationFolderPath -ChildPath "config.yaml"

function Invoke-App {
    $ErrorActionPreference = "Continue"
    try {
        & ch360 $args
    } catch [System.Management.Automation.RemoteException] {
        # Catch exceptions for messages redirected from stderr and
        # write out the messages to stdout
        Write-Output $Error[0].Message
    }
}

function New-Classifier([string]$classifierName, [Io.FileInfo]$samples) {
    Invoke-App create classifier $classifierName --samples-zip=$samples 2>&1
}

function Get-Classifiers {
    Invoke-App list classifiers 2>&1
}

function Remove-Classifier([Parameter(ValueFromPipeline=$true)]$classifierName) {
    Invoke-App delete classifier $classifierName 2>&1
}

function Format-MultilineOutput([Parameter(ValueFromPipeline=$true)]$input){
    $input -join [Environment]::NewLine
}

function String-Starting([string]$input) {
    ([Regex]::Escape($input) + ".*")
}

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
  
Describe "classifiers" {
    BeforeAll {
        Backup-ApplicationFolder
        
        ch360 login --client-id="$ClientId" --client-secret="$ClientSecret" | Should -Be "Logging in... [OK]"
        $LASTEXITCODE | Should -Be 0
        
        Get-Content -Path $configFilePath | Format-MultilineOutput | Should -BeLike "*clientId: $ClientId*"
        Get-Content -Path $configFilePath | Format-MultilineOutput | Should -BeLike "*clientSecret: $ClientSecret*"    
        Write-Host "Ran ch360 login"        
    }
    
    BeforeEach {
        Get-Classifiers | Should -Be "No classifiers found."
    }

    It "should be created from a zip file of samples" {
        $samples = (Join-Path $PSScriptRoot "samples.zip")
        New-Classifier $classifierName $samples | Format-MultilineOutput | Should -Be @"
Creating classifier '$classifierName'... [OK]
Adding samples from file '$samples'... [OK]
"@
        $LASTEXITCODE | Should -Be 0

        # Verify
        Get-Classifiers | Format-MultilineOutput | Should -Contain $classifierName

        # Teardown
        Remove-Classifier $classifierName
    }

    It "should not be created from an invalid zip file of samples" {
        $samples = (Join-Path $PSScriptRoot "invalid.zip")
        New-Classifier $classifierName $samples | Format-MultilineOutput | Should -Match (String-Starting @"
Creating classifier '$classifierName'... [OK]
Adding samples from file '$samples'... [FAILED]
"@)

        $LASTEXITCODE | Should -Be 1

        #Teardown
        Remove-Classifier $classifierName
    }

    It "should not be created from a non-existent zip file of samples" {
        $samples = (Join-Path $PSScriptRoot "non-existent.zip")
        New-Classifier $classifierName $samples | Format-MultilineOutput | Should -Be @"
Creating classifier '$classifierName'... [OK]
Adding samples from file '$samples'... [FAILED]
The file '$samples' could not be found.
"@

        $LASTEXITCODE | Should -Be 1
    }

    AfterAll {
        # Tidy up any leftover classifiers in the account
        Get-Classifiers | Remove-Classifier
        
        # Restore user's original application folder
        Restore-ApplicationFolder
    }
}