param(
    [String]$ClientId,
    [String]$ClientSecret
)

$classifierName = "test-classifier"

function New-Classifier([string]$classifierName, [Io.FileInfo]$samples) {
    $ErrorActionPreference = "Continue"
    try {
        ch360 create classifier $classifierName `
            --id="$ClientId" `
            --secret="$ClientSecret" `
            --samples-zip=$samples 2>&1
    } catch [System.Management.Automation.RemoteException] {
        # Catch exceptions for messages redirected from stderr and
        # write out the messages to stdout
        Write-Output $Error[0].Message
    }
}

function Format-MultilineOutput([Parameter(ValueFromPipeline=$true)]$input){
    $input -join [Environment]::NewLine
}

Describe "classifiers" {
    BeforeEach {
        ch360 delete classifier $classifierName `
            --id="$ClientId" `
            --secret="$ClientSecret"
    }

    It "should be created from a zip file of samples" {
        $samples = (Join-Path $PSScriptRoot "samples.zip")
        New-Classifier $classifierName $samples | Format-MultilineOutput | Should -Be @"
Creating classifier '$classifierName'... [OK]
Adding samples from file '$samples'... [OK]
"@

        $LASTEXITCODE | Should -Be 0

        # Teardown
        ch360 delete classifier $classifierName --id="$ClientId" --secret="$ClientSecret"
    }

    It "should not be created from an invalid zip file of samples" {
        $samples = (Join-Path $PSScriptRoot "invalid.zip")
        New-Classifier $classifierName $samples | Format-MultilineOutput | Should -BeLike @"
Creating classifier '$classifierName'... [OK]
Adding samples from file '$samples'... [FAILED]
*
"@

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
}

Describe "ch360 list classifiers" {
    It "should list the names of all existing classifiers" {
        # Run delete classifier first to ensure the test classifiers are not already present        
        ch360 delete classifier "${classifierName}1" --id="$ClientId" --secret="$ClientSecret"
        ch360 delete classifier "${classifierName}2" --id="$ClientId" --secret="$ClientSecret"
        
        ch360 create classifier "${classifierName}1" --id="$ClientId" --secret="$ClientSecret"
        ch360 create classifier "${classifierName}2" --id="$ClientId" --secret="$ClientSecret"

        ch360 list classifiers --id="$ClientId" --secret="$ClientSecret" | Should -Be  @(
            "test-classifier1",
            "test-classifier2"
            )            
        $LASTEXITCODE | Should -Be 0
    }

    It "should output 'No classifiers found.' when there are no classifiers" {
        # Run delete classifier first to ensure the test classifiers are not already present        
        ch360 delete classifier "${classifierName}1" --id="$ClientId" --secret="$ClientSecret"
        ch360 delete classifier "${classifierName}2" --id="$ClientId" --secret="$ClientSecret"
        
        ch360 list classifiers --id="$ClientId" --secret="$ClientSecret" | Should -Be "No classifiers found."
        $LASTEXITCODE | Should -Be 0
    }
 
    AfterAll {
        ch360 delete classifier "${classifierName}1" --id="$ClientId" --secret="$ClientSecret"
        ch360 delete classifier "${classifierName}2" --id="$ClientId" --secret="$ClientSecret"
    }
}