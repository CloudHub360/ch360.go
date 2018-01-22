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
    Invoke-App create classifier $classifierName $samples 2>&1
}

function Get-Classifiers {
    Invoke-App list classifiers 2>&1
}

function Remove-Classifier([Parameter(ValueFromPipeline=$true)]$classifierName) {
    Invoke-App delete classifier $classifierName 2>&1
}

function Invoke-Classifier([Io.FileInfo]$file, [string]$classifierName) {
    Invoke-App classify $($file.FullName) $classifierName
}

function Classify-Files-And-Write-CSV-OutputFile([string]$filePattern, [string]$classifierName, [string]$outputFile) {
    Invoke-App classify "`"$filePattern`"" $classifierName -o $outputFile -f csv
}

function Classify-Files-And-Write-Multiple-OutputFiles([string]$filePattern, [string]$classifierName, [string]$format) {
    Invoke-App classify "`"$filePattern`"" $classifierName -m -f $format
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
        # Tidy up any leftover classifiers in the account
        Get-Classifiers | Remove-Classifier
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
    }

    It "should not be created from an invalid zip file of samples" {
        $samples = (Join-Path $PSScriptRoot "invalid.zip")
        New-Classifier $classifierName $samples | Format-MultilineOutput | Should -Match (String-Starting @"
Creating classifier '$classifierName'... [OK]
Adding samples from file '$samples'... [FAILED]
"@)

        $LASTEXITCODE | Should -Be 1
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

    It "should attempt to classify a file" {
        $samples = (Join-Path $PSScriptRoot "samples.zip")
        New-Classifier $classifierName $samples
        
        $document = (Join-Path $PSScriptRoot "documents/document1.pdf")
        Invoke-Classifier $document $classifierName | Format-MultilineOutput | Should -Be @"
FILE                                 DOCUMENT TYPE                    CONFIDENT
document1.pdf                        Notice of Lien                   true
"@
    }

    It "should classify files and write to a multiple csv results files" {
        $samples = (Join-Path $PSScriptRoot "samples.zip")
        New-Classifier $classifierName $samples
        
        $filePattern = (Join-Path $PSScriptRoot "documents/subfolder1/*.pdf")
        $document2OutputFile = (Join-Path $PSScriptRoot "documents/subfolder1/document2.csv")
        $document3OutputFile = (Join-Path $PSScriptRoot "documents/subfolder1/document3.csv")
        
        Classify-Files-And-Write-Multiple-OutputFiles $filePattern $classifierName "csv"
        
        (Get-Content -Path $document2OutputFile) | Format-MultilineOutput | Should -BeLike @"
*document2.pdf,Notice of Lien,true,1.177
"@

        (Get-Content -Path $document3OutputFile) | Format-MultilineOutput | Should -BeLike @"
*document3.pdf,Notice of Default,true,3.351
"@

        Remove-Item -Path $document2OutputFile
        Remove-Item -Path $document3OutputFile
    }

    It "should classify files and write to a single csv results file" {
        $samples = (Join-Path $PSScriptRoot "samples.zip")
        New-Classifier $classifierName $samples
        
        $filePattern = (Join-Path $PSScriptRoot "documents/subfolder1/*.pdf")
        $outputFile = New-TemporaryFile
        Classify-Files-And-Write-CSV-OutputFile $filePattern $classifierName $outputFile
        
        (Get-Content -Path $outputFile) | Format-MultilineOutput | Should -BeLike @"
*document2.pdf,Notice of Lien,true,1.177
*document3.pdf,Notice of Default,true,3.351
"@
        Remove-Item -Path $outputFile
    }
    
It "should classify a file and write a json result file" {
    $samples = (Join-Path $PSScriptRoot "samples.zip")
    New-Classifier $classifierName $samples
    
    $filePattern = (Join-Path $PSScriptRoot "documents/subfolder1/document2.pdf")
    $outputFile = (Join-Path $PSScriptRoot "documents/subfolder1/document2.json")
    
    Classify-Files-And-Write-Multiple-OutputFiles $filePattern $classifierName "json"
    
    (Get-Content -Path $outputFile) | Format-MultilineOutput | Should -BeLike "{*}"

    Remove-Item -Path $outputFile
}
    AfterAll {
        # Tidy up any leftover classifiers in the account
        Get-Classifiers | Remove-Classifier

        # Restore user's original application folder
        Restore-ApplicationFolder
    }
}