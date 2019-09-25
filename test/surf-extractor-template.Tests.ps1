param(
    [String]$ClientId,
    [String]$ClientSecret
)

. "test/common.ps1"

Describe "extractor-templates" {
    $extractorTemplatePath = "test.json"

    BeforeAll {
        Login-Surf
    }

    BeforeEach {
    }

    It "should create an extractor template from module ids" {
        Invoke-App create extractor-template waives.name 2>&1 | Format-MultilineOutput | Should -Be @"
{
  "modules": [
    {
      "id": "waives.name"
    }
  ]
}
"@
        $LASTEXITCODE | Should -Be 0
    }

    It "write the template to a file if specified" {
        Invoke-App create extractor-template waives.name -o $extractorTemplatePath

        Get-Content -Raw $extractorTemplatePath | Should -Be @"
{
  "modules": [
    {
      "id": "waives.name"
    }
  ]
}
"@
    }

    AfterEach {
        if (Test-Path $extractorTemplatePath) {
            Remove-Item $extractorTemplatePath
        }
    }

    AfterAll {
        Restore-ApplicationFolder
    }
}
