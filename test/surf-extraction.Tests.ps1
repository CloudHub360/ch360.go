param(
    [String]$ClientId,
    [String]$ClientSecret
)

. "test/common.ps1"

$extractorName = "test-extractor"
$documentsPath = (Join-Path $PSScriptRoot (Join-Path "documents" "extraction"))

function New-Extractor([string]$extractorName, [Io.FileInfo]$extractorDefinition) {
    Invoke-App upload extractor $extractorName $extractorDefinition 2>&1
}

function New-ExtractorFromModules([string]$extractorName,
    [parameter(Position=0, ValueFromRemainingArguments=$true)] $moduleIds) {
    Invoke-App create extractor $extractorName @moduleIds 2>&1
}

function New-ExtractorFromTemplate([string]$extractorName, [Io.FileInfo]$extractorTemplate) {
    Invoke-App create extractor $extractorName --from-template $extractorTemplate 2>&1
}

function Get-Extractors {
    Invoke-App list extractors 2>&1
}

function Remove-Extractor([Parameter(ValueFromPipeline=$true)]$extractorName) {
    Write-Debug "Deleting extractor: $extractorName"
    Invoke-App delete extractor $extractorName 2>&1
}

function Remove-UserExtractors() {
    Get-Extractors | Where-Object { !$_.StartsWith("waives.") } | Remove-Extractor
}

Describe "extractors" {
    BeforeAll {
        Login-Surf
    }

    BeforeEach {
        Remove-UserExtractors
    }

    It "should be created from an fpxlc definition file" {
        $extractorDefinition = (Join-Path $PSScriptRoot "extract-amount.fpxlc")
        New-Extractor $extractorName $extractorDefinition | Format-MultilineOutput | Should -Be @"
Uploading extractor '$extractorName'... [OK]
"@
        $LASTEXITCODE | Should -Be 0

        # Verify
        Get-Extractors | Format-MultilineOutput | Should -Match $extractorName
    }

   It "should not be created from an invalid fpxlc definition file" {
       $extractorDefinition = (Join-Path $PSScriptRoot "invalid.fpxlc")
       New-Extractor $extractorName $extractorDefinition | Format-MultilineOutput | Should -Be @"
Uploading extractor 'test-extractor'... [FAILED]
The file supplied is not a valid extractor configuration file.
"@
       $LASTEXITCODE | Should -Be 1

       # Verify
       Get-Extractors | Format-MultilineOutput | Should -Not -Match $extractorName
   }

    It "should not be created from a non-existent fpxlc definition file" {
        $extractorDefinition = (Join-Path $PSScriptRoot "non-existent.fpxlc")
        New-Extractor $extractorName $extractorDefinition | Format-MultilineOutput | Should -Be "The file '$extractorDefinition' could not be found."
        $LASTEXITCODE | Should -Be 1

        # Verify
        Get-Extractors | Format-MultilineOutput | Should -Not -Match $extractorName
    }

    It "should be created from a list of module IDs " {
        New-ExtractorFromModules $extractorName "waives.invoice_number" "waives.name" | Format-MultilineOutput | Should -Be @"
Creating extractor '$extractorName'... [OK]
"@
        $LASTEXITCODE | Should -Be 0

        # Verify
        Get-Extractors | Format-MultilineOutput | Should -Match $extractorName
    }

    It "should be created from a template" {
        $extractorTemplate = (Join-Path $PSScriptRoot "extractor-template.json")
        New-ExtractorFromTemplate $extractorName $extractorTemplate | Format-MultilineOutput | Should -Be @"
Creating extractor '$extractorName'... [OK]
"@
        $LASTEXITCODE | Should -Be 0

        # Verify
        Get-Extractors | Format-MultilineOutput | Should -Match $extractorName
    }

    It "should list available modules" {
        $output = Invoke-App list modules 2>&1 | Format-MultilineOutput

        $output | Should -Match "Name\s+ID\s+Summary"
        $output | Should -Match "Name\s+waives.name\s+.*"
        ($output | Measure-Object -Line).Lines | Should -BeGreaterThan 20
    }

    It "should not be created from a missing extractor template" {
        $extractorTemplate = "missing.json"
        New-ExtractorFromTemplate $extractorName $extractorTemplate | Format-MultilineOutput | Should -Be @"
Error when opening template file 'missing.json': no such file or directory
"@
        $LASTEXITCODE | Should -Be 1
    }

    It "should not be created from invalid json" {
        $extractorTemplate = (Join-Path $PSScriptRoot "invalid.json")
        New-ExtractorFromTemplate $extractorName $extractorTemplate | Format-MultilineOutput | Should -Be @"
Error when reading json template '$extractorTemplate': invalid character 'b' looking for beginning of value
"@
        $LASTEXITCODE | Should -Be 1
    }

    AfterAll {
        Remove-UserExtractors
        Restore-ApplicationFolder
    }
}

Describe "extraction" {
    BeforeAll {
        Login-Surf
    }

    BeforeEach {
        Remove-UserExtractors

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
File                                Amount
document1.pdf                       `$5.50
"@
    }

    It "should extract data and write to a multiple csv results files" {
        $filePattern = (Join-Path $documentsPath "/subfolder1/*.pdf")
        $document2OutputFile = (Join-Path $documentsPath "subfolder1/document2.csv")
        $document3OutputFile = (Join-Path $documentsPath "subfolder1/document3.csv")

        Invoke-Extractor $filePattern $extractorName -Format "csv"

        Get-Content $document2OutputFile | ConvertFrom-ExtractionCsv `
          | where { $_.file -like "*document2.pdf" -and $_.amount -eq "`$5.50|`$4" } `
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
        $results | where { $_.file -like "*document2.pdf" -and $_.amount -eq "`$5.50|`$4" } | Should -Not -Be $null
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
        Remove-UserExtractors
        Restore-ApplicationFolder
    }
}