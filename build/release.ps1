param(
    [string]$VersionString,
    [string]$GitHubToken,
    [string]$Commit = "master",
    [string]$ReleaseNotes
)

$GITHUB_RELEASES_URI = "https://api.github.com/repos/CloudHub360/ch360.go/releases"

Task Release {
    $releaseMetadata = [PSCustomObject] @{
        "tag_name" = "v$VersionString"
        "target_commitish" = $Commit
        "name" = $VersionString
        "body" = $ReleaseNotes
    }

    Invoke-RestMethod `
        -Method POST `
        -Headers @{ "Authorization"="token $GitHubToken" } `
        -Uri $GITHUB_RELEASES_URI `
        -ContentType "application/json" `
        -Body (ConvertTo-Json $releaseMetadata)
}

Task . Release