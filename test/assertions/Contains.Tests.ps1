Describe "Contains assertions with int arrays" {
    It "5 should be in the array 1..10" {
        1..10 | Should -Contain 5
    }

    It "11 should not be in the array 1..10" {
        1..10 | Should -Not -Contain 11
    }
}
