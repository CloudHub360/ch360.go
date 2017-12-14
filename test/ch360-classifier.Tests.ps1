param(
    [String]$ClientId,
    [String]$ClientSecret
)

$classifierName = "test-classifier"

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
    Invoke-App create classifier $classifierName `
        --id="$ClientId" `
        --secret="$ClientSecret" `
        --samples-zip=$samples 2>&1
}

function Get-Classifiers {
    Invoke-App list classifiers `
        --id="$ClientId" `
        --secret="$ClientSecret" 2>&1
}

function Remove-Classifier([Parameter(ValueFromPipeline=$true)]$classifierName) {
    Invoke-App delete classifier $classifierName `
        --id="$ClientId" `
        --secret="$ClientSecret" 2>&1
}

function Format-MultilineOutput([Parameter(ValueFromPipeline=$true)]$input){
    $input -join [Environment]::NewLine
}

function String-Starting([string]$input) {
    ([Regex]::Escape($input) + ".*")
}

Describe "classifiers" {
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
    }
}