param(
    [String]$ClientId,
    [String]$ClientSecret
)

. "test/common.ps1"

$extractorName = "test-classifier"

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

    It "should be created from a zip file of samples" {

    }

    It "should not be created from an invalid zip file of samples" {

    }

    It "should not be created from a non-existent zip file of samples" {

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