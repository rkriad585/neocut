param(
    [string]$Name = "neocut",
    [string]$OutputDir = (Join-Path $PSScriptRoot "bin")
)

$Version = (Get-Content (Join-Path $PSScriptRoot ".version") -Raw).Trim()
$Commit = & git -C $PSScriptRoot rev-parse --short HEAD 2>$null
if (-not $Commit) { $Commit = "unknown" }

$PublisherName = "rkriad585"
$PublisherEmail = "rkriad585@gmail.com"

$LdFlags = @(
    "-X 'neocut/internal/config.Commit=$Commit'",
    "-X 'neocut/internal/config.Version=$Version'",
    "-X 'neocut/internal/config.PublisherName=$PublisherName'",
    "-X 'neocut/internal/config.PublisherEmail=$PublisherEmail'"
) -join " "

$Platforms = @(
    @{ OS = "windows"; Arch = "amd64"; Ext = ".exe" },
    @{ OS = "windows"; Arch = "arm64"; Ext = ".exe" },
    @{ OS = "linux";   Arch = "amd64"; Ext = "" },
    @{ OS = "linux";   Arch = "arm64"; Ext = "" },
    @{ OS = "darwin";  Arch = "amd64"; Ext = "" },
    @{ OS = "darwin";  Arch = "arm64"; Ext = "" }
)

Write-Host "  Generating embedded assets..."
& go generate ./internal/config/ 2>&1

New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null

Write-Host "╭──────────────── neocut build ────────────────╮"
Write-Host "│  Version : $($Version.PadRight(20))│"
Write-Host "│  Commit  : $($Commit.PadRight(20))│"
Write-Host "│  Publisher: $($PublisherName.PadRight(17))│"
Write-Host "│  Email   : $($PublisherEmail.PadRight(19))│"
Write-Host "╰─────────────────────────────────────────────╯"
Write-Host ""

$count = 0
foreach ($p in $Platforms) {
    $binary = "$Name-$($p.OS)-$($p.Arch)$($p.Ext)"
    $path = Join-Path $OutputDir $binary

    Write-Host "  [$($count+1)/$($Platforms.Length)] $binary" -NoNewline

    $env:GOOS = $p.OS
    $env:GOARCH = $p.Arch

    $result = & go build -ldflags $LdFlags -o $path ./cmd/neocut/ 2>&1
    if ($LASTEXITCODE -eq 0) {
        $size = (Get-Item $path).Length
        $sizeStr = if ($size -gt 1MB) { "{0:N1} MB" -f ($size / 1MB) } else { "{0:N0} KB" -f ($size / 1KB) }
        Write-Host "  ✓ $sizeStr" -ForegroundColor Green
        $count++
    } else {
        Write-Host "  ✗ FAILED" -ForegroundColor Red
        Write-Host "    $result" -ForegroundColor DarkRed
    }
}

Remove-Item Env:GOOS, Env:GOARCH -ErrorAction SilentlyContinue

Write-Host ""
if ($count -eq $Platforms.Length) {
    Write-Host "  All $count binaries built successfully in $OutputDir"
} else {
    Write-Host "  $count/$($Platforms.Length) binaries built (see errors above)"
}
Write-Host ""
