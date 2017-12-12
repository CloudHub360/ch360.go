param(
    [String]$ClientId,
    [String]$ClientSecret
)

$classifierName = "test-classifier"

Describe "ch360 delete classifier" {
    It "should output success" {
        # Run delete classifier first to ensure it's not already present
        ch360 delete classifier $classifierName --id="$ClientId" --secret="$ClientSecret"
        
        ch360 create classifier $classifierName --id="$ClientId" --secret="$ClientSecret" | Should -Be "Creating classifier '$classifierName'... [OK]"
        ch360 delete classifier $classifierName --id="$ClientId" --secret="$ClientSecret" | Should -Be "Deleting classifier '$classifierName'... [OK]"
        $LASTEXITCODE | Should -Be 0
    }

    It "should output failure when the classifier does not exist" {
        ch360 delete classifier $classifierName --id="$ClientId" --secret="$ClientSecret" | Should -Be "Deleting classifier '$classifierName'... [FAILED]"
        $LASTEXITCODE | Should -Be 1
    }

    It "should output failure when client id or secret are invalid" {
        ch360 delete classifier $classifierName --id="invalid-id" --secret="$ClientSecret" | Should -Be "Deleting classifier '$classifierName'... [FAILED]"
        $LASTEXITCODE | Should -Be 1
    }
}

Describe "ch360 create classifier" {
    It "should output success when creating a new classifier" {
        # Run delete classifier first to ensure it's not already present
        ch360 delete classifier $classifierName --id="$ClientId" --secret="$ClientSecret"
        ch360 create classifier $classifierName --id="$ClientId" --secret="$ClientSecret" | Should -Be "Creating classifier '$classifierName'... [OK]"
        $LASTEXITCODE | Should -Be 0
    }

    It "should output failure when attempting to create a classifier that already exists" {
        # Run delete classifier first to ensure it's not already present
        ch360 create classifier $classifierName --id="$ClientId" --secret="$ClientSecret"
        ch360 create classifier $classifierName --id="$ClientId" --secret="$ClientSecret" | Should -Be "Creating classifier '$classifierName'... [FAILED]"
        $LASTEXITCODE | Should -Be 1
    }

    It "should output failure when client id or secret are invalid" {
        ch360 create classifier $classifierName --id="invalid-id" --secret="$ClientSecret" | Should -Be "Creating classifier '$classifierName'... [FAILED]"
        $LASTEXITCODE | Should -Be 1
    }

    It "should train the classifier with the provided ZIP file" {
        # Run delete classifier first to ensure it's not already present
        ch360 delete classifier $classifierName --id="$ClientId" --secret="$ClientSecret"
        $samples = (Join-Path $pwd "samples.zip")

        ch360 create classifier $classifierName --id="$ClientId" --secret="$ClientSecret" --samples-zip=$samples `
            | Should -Be `
@"
Creating classifier '$classifierName'... [OK]
Adding samples from file '$samples'... [OK]
"@

        $LASTEXITCODE | Should -Be 0
    }

    It "should output failure when the sample ZIP file does not exist at the specified path" {
        $samples = (Join-Path $pwd "non-existent.zip")
        ch360 create classifier $classifierName --id="$ClientId" --secret="$ClientSecret" --samples-zip=$samples `
            | Should -Be `
@"
Creating classifier '$classifierName'... [FAILED]
The file '$samples' could not be found.
"@

        $LASTEXITCODE | Should -Be 1
    }

    It "should output failure when the sample ZIP file does not match the required format " {
        $samples = (Join-Path $pwd "invalid.zip")
        ch360 create classifier $classifierName --id="$ClientId" --secret="$ClientSecret" --samples-zip=$samples `
            | Should -Be `
@"
Creating classifier '$classifierName'... [OK]
Adding samples from file '$samples'... [FAILED]
The samples file is invalid.
"@

        $LASTEXITCODE | Should -Be 1
    }

    AfterAll {
        ch360 delete classifier $classifierName --id="$ClientId" --secret="$ClientSecret"
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