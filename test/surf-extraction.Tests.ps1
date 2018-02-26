param(
    [String]$ClientId,
    [String]$ClientSecret
)

. "test/common.ps1"

$extractorName = "test-extractor"

function New-Extractor([string]$extractorName, [Io.FileInfo]$extractorDefinition) {
    Invoke-App create extractor $extractorName $extractorDefinition 2>&1
}

function Get-Extractors {
    Invoke-App list extractors 2>&1
}

function Remove-Extractor([Parameter(ValueFromPipeline=$true)]$extractorName) {
    Invoke-App delete extractor $extractorName 2>&1
}

Describe "extractors" {
    BeforeAll {
        Login-Surf
    }

    BeforeEach {
        # Tidy up any leftover classifiers in the account
        Get-Extractors | Remove-Extractor
    }

    It "should be created from an fpxlc definition file" {
        $extractorDefinition = (Join-Path $PSScriptRoot "extract-amount.fpxlc")
        New-Extractor $extractorName $extractorDefinition | Format-MultilineOutput | Should -Be @"
Creating extractor '$extractorName'... [OK]
"@
        $LASTEXITCODE | Should -Be 0

        # Verify
        Get-Extractors | Format-MultilineOutput | Should -Contain $extractorName
    }

    It "should not be created from an invalid fpxlc definition file" {
        $extractorDefinition = (Join-Path $PSScriptRoot "invalid.fpxlc")
        New-Extractor $extractorName $extractorDefinition | Format-MultilineOutput | Should -Be @"
Creating extractor 'test-extractor'... [FAILED]
The file supplied is not a valid extractor configuration file.
"@
        $LASTEXITCODE | Should -Be 1

        # Verify
        Get-Extractors | Format-MultilineOutput | Should -Not -Contain $extractorName
    }

    It "should not be created from a non-existent fpxlc definition file" {
        $extractorDefinition = (Join-Path $PSScriptRoot "non-existent.fpxlc")
        New-Extractor $extractorName $extractorDefinition | Format-MultilineOutput | Should -Be @"
open F:\code\go\src\github.com\CloudHub360\ch360.go\test\non-existent.fpxlc: The system cannot find the file specified.
"@
        $LASTEXITCODE | Should -Be 1

        # Verify
        Get-Extractors | Format-MultilineOutput | Should -Not -Contain $extractorName
    }

    AfterAll {
        # Tidy up any leftover classifiers in the account
        Get-Extractors | Remove-Extractor

        # Restore user's original application folder
        Restore-ApplicationFolder
    }
}

Describe "extraction" {
    BeforeAll {
        Login-Surf
    }

    BeforeEach {
        # Tidy up any leftover classifiers in the account
        Get-Extractors | Remove-Extractor
    }

    It "should attempt to classify a file" {

    }

    It "should classify files and write to a multiple csv results files" {

    }

    It "should classify files and write to a single csv results file" {

    }

    It "should classify a file and write a json result file" {

    }

    AfterAll {
        # Tidy up any leftover classifiers in the account
        Get-Extractors | Remove-Extractor

        # Restore user's original application folder
        Restore-ApplicationFolder
    }
}