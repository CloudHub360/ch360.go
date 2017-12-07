param(
    [String]$ClientId,
    [String]$ClientSecret
)

Describe "ch360 delete classifier" {
    It "should output success" {
        $classifierName = "test-classifier"
        ch360 delete classifier $classifierName --id="$ClientId" --secret="$ClientSecret" | Should -Be "Deleting classifier '$classifierName'... [OK]"
        $LASTEXITCODE | Should -Be 0
    }

    It "should output failure when client id or secret are invalid" {
        $classifierName = "test-classifier"
        ch360 delete classifier $classifierName --id="invalid-id" --secret="$ClientSecret" | Should -Be "Deleting classifier '$classifierName'... [FAILED]"
        $LASTEXITCODE | Should -Be 1
    }
}

Describe "ch360 create classifier" {
    It "should output success when creating a new classifier" {
        $classifierName = "test-classifier"
        # Run delete classifier first to ensure it's not already present
        ch360 delete classifier $classifierName --id="$ClientId" --secret="$ClientSecret"
        ch360 create classifier $classifierName --id="$ClientId" --secret="$ClientSecret" | Should -Be "Creating classifier '$classifierName'... [OK]"
        $LASTEXITCODE | Should -Be 0

        # Clean up
        ch360 delete classifier $classifierName --id="$ClientId" --secret="$ClientSecret"
    }

    It "should output failure when attempting to create a classifier that already exists" {
        $classifierName = "test-classifier"
        # Run delete classifier first to ensure it's not already present
        ch360 create classifier $classifierName --id="$ClientId" --secret="$ClientSecret"
        ch360 create classifier $classifierName --id="$ClientId" --secret="$ClientSecret" | Should -Be "Creating classifier '$classifierName'... [FAILED]"
        $LASTEXITCODE | Should -Be 1
    }

    It "should output failure when client id or secret are invalid" {
        $classifierName = "test-classifier"
        ch360 create classifier $classifierName --id="invalid-id" --secret="$ClientSecret" | Should -Be "Creating classifier '$classifierName'... [FAILED]"
        $LASTEXITCODE | Should -Be 1
    }
}
