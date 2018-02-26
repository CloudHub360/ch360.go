param(
    [String]$ClientId,
    [String]$ClientSecret
)

. "test/common.ps1"

$classifierName = "test-classifier"

function New-Classifier([string]$classifierName, [Io.FileInfo]$samples) {
    Invoke-App create classifier $classifierName $samples 2>&1
}

function Get-Classifiers {
    Invoke-App list classifiers 2>&1
}

function Remove-Classifier([Parameter(ValueFromPipeline=$true)]$classifierName) {
    Invoke-App delete classifier $classifierName 2>&1
}

function Invoke-Classifier(
    [Parameter(ParameterSetName="SingleFile", Position=0)]
    [Io.FileInfo]$file,

    [Parameter(Mandatory=$true, Position=1)]
    [string]$classifierName,

    [Parameter(ParameterSetName="MultipleFiles", Position=0)]
    [string]$filePattern,

    [Parameter(ParameterSetName="MultipleFiles", Position=2)]
    [Io.FileInfo]$outputFile,

    [Parameter(ParameterSetName="MultipleFiles", Position=3)]
    [string]$Format = "csv")
{
    if ($file -ne $null) {
        Invoke-App classify $($file.FullName) $classifierName
    } else {
        if ($outputFile -ne $null) {
            Invoke-App classify "`"$filePattern`"" $classifierName -o $outputFile -f $Format
        } else {
            Invoke-App classify "`"$filePattern`"" $classifierName -m -f $Format
        }
    }
}

Describe "classifiers" {
    BeforeAll {
        Login-Surf
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
        Invoke-Classifier -File $document $classifierName | Format-MultilineOutput | Should -Be @"
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

        Invoke-Classifier $filePattern $classifierName -Format "csv"

        Get-Content $document2OutputFile | ConvertFrom-WaivesCsv `
          | where {($_.file -like "*document2.pdf" -and $_.documenttype -eq "Notice of Lien" -and $_.score -eq 1.177 -and $_.confident)} `
          | Should -Not -Be $null

        Get-Content $document3OutputFile | ConvertFrom-WaivesCsv `
          | where {($_.file -like "*document3.pdf" -and $_.documenttype -eq "Notice of Default" -and $_.score -eq 3.351 -and $_.confident)} `
          | Should -Not -Be $null

        Remove-Item -Path $document2OutputFile
        Remove-Item -Path $document3OutputFile
    }

    It "should classify files and write to a single csv results file" {
        $samples = (Join-Path $PSScriptRoot "samples.zip")
        New-Classifier $classifierName $samples

        $filePattern = (Join-Path $PSScriptRoot "documents/subfolder1/*.pdf")
        $outputFile = New-TemporaryFile
        Invoke-Classifier $filePattern $classifierName -OutputFile $outputFile -Format "csv"

        $results = Get-Content $outputFile | ConvertFrom-WaivesCsv

        $results.length | Should -Be 3
        $doc2result = ($results | where {($_.file -like "*document2.pdf" -and $_.documenttype -eq "Notice of Lien" -and $_.score -eq 1.177 -and $_.confident)})
        $doc2result | Should -Not -Be $null

        $doc3result = ($results | where {($_.file -like "*document3.pdf" -and $_.documenttype -eq "Notice of Default" -and $_.score -eq 3.351 -and $_.confident)})
        $doc3result | Should -Not -Be $null

        Remove-Item -Path $outputFile
    }

    It "should classify a file and write a json result file" {
        $samples = (Join-Path $PSScriptRoot "samples.zip")
        New-Classifier $classifierName $samples

        $filePattern = (Join-Path $PSScriptRoot "documents/subfolder1/document2.pdf")
        $outputFile = (Join-Path $PSScriptRoot "documents/subfolder1/document2.json")

        Invoke-Classifier $filePattern $classifierName -OutputFile $outputFile -Format "json"

        $results = Get-Content -Path $outputFile | ConvertFrom-Json
        Write-Host $results
        Write-Host $results.file

        $results.filename | Should -BeLike "*document2.pdf"
        $results.classification_results.document_type | Should -Be "Notice of Lien"
        $results.classification_results.is_confident | Should -Be True
        $results.classification_results.relative_confidence | Should -Be 1.17683673

        Remove-Item -Path $outputFile
    }

    AfterAll {
        # Tidy up any leftover classifiers in the account
        Get-Classifiers | Remove-Classifier

        # Restore user's original application folder
        Restore-ApplicationFolder
    }
}