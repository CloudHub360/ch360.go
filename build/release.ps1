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

function Get-PlatformName([string]$path){
    $platform = Get-ParentDirectoryName $path

    if ($platform -eq "bin") {
        $platform = "$(go env GOHOSTOS)_$(go env GOHOSTARCH)"
    }

    return $platform
}

function Get-AssetName([Io.FileSystemInfo]$path) {
    $binaryName = $path.Name
    $extensionIndex = $binaryName.IndexOf(".exe")
    $platform = Get-PlatformName $path

    if ($extensionIndex -ge 0){
        return $binaryName.Insert($extensionIndex, "-$platform")
    } else {
        return "$binaryName-$platform"
    }
}

function Get-UploadUrl([string]$urlTemplate, [string]$assetName) {
    $queryParamsPosition = $urlTemplate.LastIndexOf("{")
    $url = $urlTemplate.Substring(0, $queryParamsPosition)

    return "$($url)?name=$assetName"
}

Task Release {
    $releaseMetadata = [PSCustomObject] @{
        "tag_name" = "v$VersionString"
        "target_commitish" = $Commit
        "name" = $VersionString
        "body" = $ReleaseNotes
    }

    $release = [PSCustomObject] (Invoke-RestMethod `
        -Method POST `
        -Headers @{ "Authorization"="token $GitHubToken" } `
        -Uri $GITHUB_RELEASES_URI `
        -ContentType "application/json" `
        -Body (ConvertTo-Json $releaseMetadata))

    $binaries = Get-ChildItem $GO_BIN -Recurse -Include surf*

    $binaries |% {
        $uploadUrl = Get-UploadUrl $release.upload_url (Get-AssetName $_)
        Write-Information "Creating asset at $uploadUrl"
        Invoke-RestMethod `
            -Method POST `
            -Headers @{ "Authorization"="token $GitHubToken" } `
            -Uri $uploadUrl `
            -ContentType "application/octet-stream" `
            -InFile $_.FullName
    }
}

Task . Release