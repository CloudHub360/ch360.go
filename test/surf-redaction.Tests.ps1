param(
    [String]$ClientId,
    [String]$ClientSecret
)

. "test/common.ps1"

$extractorName = "test-extractor"
$documentsPath = (Join-Path $PSScriptRoot (Join-Path "documents" "redaction"))

Describe "redaction" {
    BeforeAll {
        Login-Surf
    }

    BeforeEach {
        Remove-UserExtractors
    }

    It "should redact a file using an extractor" {
        $docFile = (Join-Path $documentsPath "document1.pdf")
        New-ExtractorFromModules $extractorName waives.name | Format-MultilineOutput | Should -Be @"
Creating extractor '$extractorName'... [OK]
"@
        Invoke-App redact with-extractor $extractorName $docFile -o redacted.pdf

        # Verify
        $LASTEXITCODE | Should -Be 0
        Test-PDFFile $docFile | Should -Be $true
        Remove-Item -Path redacted.pdf
    }

    AfterAll {
        Remove-UserExtractors
        Restore-ApplicationFolder
    }
}
