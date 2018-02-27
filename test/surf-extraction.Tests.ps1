param(
    [String]$ClientId,
    [String]$ClientSecret
)

. "test/common.ps1"

$extractorName = "test-extractor"
$documentsPath = (Join-Path $PSScriptRoot (Join-Path "documents" "extraction"))

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
        Get-Extractors | Remove-Extractor
        Restore-ApplicationFolder
    }
}

Describe "extraction" {
    BeforeAll {
        Login-Surf
    }

    BeforeEach {
        Get-Extractors | Remove-Extractor

        $extractorDefinition = (Join-Path $PSScriptRoot "extract-amount.fpxlc")
        New-Extractor $extractorName $extractorDefinition
    }

    function Invoke-Extractor(
        [Parameter(ParameterSetName="SingleFile", Position=0)]
        [Io.FileInfo]$file,

        [Parameter(Mandatory=$true, Position=1)]
        [string]$extractorName,

        [Parameter(ParameterSetName="MultipleFiles", Position=0)]
        [string]$filePattern,

        [Parameter(ParameterSetName="MultipleFiles", Position=2)]
        [Io.FileInfo]$outputFile,

        [Parameter(ParameterSetName="MultipleFiles", Position=3)]
        [string]$Format = "csv"
    )
    {
        if ($file -ne $null) {
            Invoke-App extract $($file.FullName) $extractorName
        } else {
            if ($outputFile -ne $null) {
                Invoke-App extract "`"$filePattern`"" $extractorName -o $outputFile -f $Format
            } else {
                Invoke-App extract "`"$filePattern`"" $extractorName -m -f $Format
            }
        }
    }

    function ConvertFrom-ExtractionCsv([Parameter(ValueFromPipeline=$true)][PsObject[]]$InputObject) {
        PROCESS {
            $InputObject | ConvertFrom-Csv -Header "file", "amount"
        }
    }

    It "should attempt to extract data from a document" {
        $document = (Join-Path $documentsPath "document1.pdf")
        Invoke-Extractor -File $document $extractorName | Format-MultilineOutput | Should -Be @"
FILE                                Amount
document1.pdf                       `$5.50
"@
    }

    It "should extract data and write to a multiple csv results files" {
        $filePattern = (Join-Path $documentsPath "/subfolder1/*.pdf")
        $document2OutputFile = (Join-Path $documentsPath "subfolder1/document2.csv")
        $document3OutputFile = (Join-Path $documentsPath "subfolder1/document3.csv")

        Invoke-Extractor $filePattern $extractorName -Format "csv"

        Get-Content $document2OutputFile | ConvertFrom-ExtractionCsv `
          | where { $_.file -like "*document2.pdf" -and $_.amount -eq "`$5.50" } `
          | Should -Not -Be $null

        Get-Content $document3OutputFile | ConvertFrom-ExtractionCsv `
          | where { $_.file -like "*document3.pdf" -and $_.amount -eq "`$5.50" } `
          | Should -Not -Be $null

        Remove-Item -Path $document2OutputFile
        Remove-Item -Path $document3OutputFile
    }

    It "should extract data and write to a single csv results file" {
        $filePattern = (Join-Path $documentsPath "subfolder1/*.pdf")
        $outputFile = New-TemporaryFile
        Invoke-Extractor $filePattern $extractorName -OutputFile $outputFile -Format "csv"

        $results = Get-Content $outputFile | ConvertFrom-ExtractionCsv

        $results.length | Should -Be 3
        $results | where { $_.file -like "*document2.pdf" -and $_.amount -eq "`$5.50" } | Should -Not -Be $null
        $results | where { $_.file -like "*document3.pdf" -and $_.amount -eq "`$5.50" } | Should -Not -Be $null

        Remove-Item -Path $outputFile
    }

    It "should extract data and write a json result file" {
        $filePattern = (Join-Path $documentsPath "subfolder1/document2.pdf")
        $outputFile = (Join-Path $documentsPath "subfolder1/document2.json")

        Invoke-Extractor $filePattern $extractorName -OutputFile $outputFile -Format "json"

        $results = Get-Content -Path $outputFile | ConvertFrom-Json

        $results.filename | Should -BeLike "*document2.pdf"
        $results.field_results.field_name | Should -Be 'Amount'
        $results.field_results.result.text | Should -Be '$5.50'
        $results.field_results.result.rejected | Should -Be $false
        $results.field_results.result.areas | Should -Be -Not $null

        Remove-Item -Path $outputFile
    }

    AfterAll {
        Get-Extractors | Remove-Extractor
        Restore-ApplicationFolder
    }
}