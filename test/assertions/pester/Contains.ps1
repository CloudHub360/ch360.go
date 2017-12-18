<#
.SYNOPSIS
Tests whether the provided value is contained within the collection
#>
function Contains {
    [CmdletBinding()]
    param(
        $ActualValue,
        $ExpectedItem,
        [switch]$Negate
    )

    [bool]$pass = $ActualValue -contains $ExpectedItem
    if ($Negate) { $pass = -not $pass }

    if (-not $pass) {
        if ($Negate) {
            $failureMessage = 'Expected: collection {{{0}}} to not contain the item {{{1}}} but it did.' -f $ActualValue, $ExpectedItem
        } else {
            $failureMessage = 'Expected: collection {{{0}}} to contain the item {{{1}}} but it was not found.' -f $ActualValue, $ExpectedItem
        }
    }

    return New-Object PSObject -Property @{
        Succeeded = $pass
        FailureMessage = $failureMessage
    }
}

Add-AssertionOperator -Name Contain -Test $Function:Contains -SupportsArrayInput