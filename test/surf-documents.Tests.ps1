param(
    [String]$ClientId,
    [String]$ClientSecret
)

. "test/common.ps1"

$documentsPath = (Join-Path $PSScriptRoot (Join-Path "documents" "classification"))

function New-Document([Io.FileInfo]$document) {
    Invoke-App create document $document 2>&1
}

function Get-Documents {
    Invoke-App list documents 2>&1
}

function Remove-Document([Parameter(ValueFromPipeline=$true)]$documentId) {
    Invoke-App delete document $documentId 2>&1
}

Describe "documents" {
    BeforeAll {
        Login-Surf
    }

    BeforeEach {
        # Tidy up any leftover documents in the account
        Remove-Document --all
    }

    It "should be created from a file" {
        $documentPath = (Join-Path $documentsPath "document1.pdf")
        New-Document $documentPath | Format-MultilineOutput | Should -Match @"
Creating document...
  File                                      ID                      Size   Type               SHA256\s*
--------------------------------------------------------------------------------------------------------------------------------------------------------------------
  ...ocuments/classification/document1.pdf  ......................  50248  PDF:ImagePlusText  650c75913be04fa0f790abdcaddae6c9093b1d575cffbed2e098eb0de0e1d4b1..
"@
        $LASTEXITCODE | Should -Be 0
    }

    It "should be listed once created" {
        $documentPath = (Join-Path $documentsPath "document1.pdf")
        New-Document $documentPath
        $LASTEXITCODE | Should -Be 0

        Get-Documents | Format-MultilineOutput | Should -Match @"
  ID                      Size   Type               SHA256\s*
-------------------------------------------------------------------------------------------------------------------------
  ......................  50248  PDF:ImagePlusText  650c75913be04fa0f790abdcaddae6c9093b1d575cffbed2e098eb0de0e1d4b1\s*
"@
    }

    It "should print a useful message if no documents present" {
        Get-Documents | Format-MultilineOutput | Should -Match @"
No documents found.
"@
    }

    AfterAll {
        # Tidy up any leftover documents in the account
        Remove-Document --all

        # Restore user's original application folder
        Restore-ApplicationFolder
    }
}