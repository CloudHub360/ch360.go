param(
    [string]$VersionString,
    [string]$GitHubToken,
    [string]$Commit = "master",
    [string]$ReleaseNotes
)

$GITHUB_RELEASES_URI = "https://api.github.com/repos/CloudHub360/ch360.go/releases"
$GO_BIN = (Join-Path $env:GOPATH bin)

function Get-ParentDirectoryName([string]$path) {
    Split-Path (Split-Path $path -Parent) -Leaf
}

function Get-AssetName([Io.FileSystemInfo]$path) {
    $binaryName = $path.Name
    $extensionIndex = $binaryName.IndexOf(".exe")
    $platform = Get-ParentDirectoryName $path

    if ($extensionIndex -ge 0){
        return $binaryName.Insert($extensionIndex, "-$platform")
    } else {
        return "$binaryName-$platform"
    }
}

Task Release {
    $releaseMetadata = [PSCustomObject] @{
        "tag_name" = "v$VersionString"
        "target_commitish" = $Commit
        "name" = $VersionString
        "body" = $ReleaseNotes
    }

    $release = Invoke-RestMethod `
        -Method POST `
        -Headers @{ "Authorization"="token $GitHubToken" } `
        -Uri $GITHUB_RELEASES_URI `
        -ContentType "application/json" `
        -Body (ConvertTo-Json $releaseMetadata)

    $binaries = Get-ChildItem $GO_BIN -Recurse -Include ch360*

    $binaries |% {
        Invoke-RestMethod `
            -Method POST `
            -Headers @{ "Authorization"="token $GitHubToken" } `
            -Uri "$GITHUB_RELEASES_URI/$($release.id)/assets?name=$(Get-AssetName $_)" `
            -ContentType "application/octet-stream" `
            -Body (Get-Content $_.FullName)
    }
}

Task . Release