param(
    [switch]$SelfUninstall = $false
)

$ProjectName = "neocut"
$GitHubUser = "rkriad585"
$ConfigDir = "$env:USERPROFILE\.config\neostore\$ProjectName"
$BinDir = "$ConfigDir\bin"
$BinaryPath = "$BinDir\$ProjectName.exe"

function Write-Step {
    param([string]$Message, [string]$Status = "INFO")
    $symbol = switch ($Status) {
        "OK"   { "✓" }
        "FAIL" { "✗" }
        "SKIP" { "‣" }
        default { "·" }
    }
    Write-Host "  $symbol $Message"
}

function Test-Admin {
    $id = [System.Security.Principal.WindowsIdentity]::GetCurrent()
    $principal = New-Object System.Security.Principal.WindowsPrincipal($id)
    return $principal.IsInRole([System.Security.Principal.WindowsBuiltInRole]::Administrator)
}

function Add-ToUserPath {
    param([string]$Dir)
    $current = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($current -like "*$Dir*") { return $false }
    $new = if ($current) { "$current;$Dir" } else { $Dir }
    [Environment]::SetEnvironmentVariable("PATH", $new, "User")
    return $true
}

function Remove-FromUserPath {
    param([string]$Dir)
    $current = [Environment]::GetEnvironmentVariable("PATH", "User")
    if (-not $current) { return $false }
    $entries = $current -split ';' | Where-Object { $_ -ne $Dir }
    $new = $entries -join ';'
    if ($new -eq $current) { return $false }
    [Environment]::SetEnvironmentVariable("PATH", $new, "User")
    return $true
}

# ── Handle self-uninstall ──────────────────────────────────
if ($SelfUninstall -or ($args -contains "--selfuninstall")) {
    Write-Host "╭──────────────── $ProjectName uninstall ───────────────╮"
    Write-Host "│"
    $removed = $false
    if (Test-Path $BinaryPath) {
        Remove-Item $BinaryPath -Force
        Write-Step "$ProjectName binary removed" "OK"
        $removed = $true
    } else {
        Write-Step "No binary found at $BinaryPath" "SKIP"
    }
    if (Test-Path $BinDir) {
        $remaining = Get-ChildItem $BinDir -ErrorAction SilentlyContinue
        if (-not $remaining) {
            Remove-Item $BinDir -Force -ErrorAction SilentlyContinue
            Write-Step "Bin directory removed" "OK"
        }
    }
    if (Remove-FromUserPath $BinDir) {
        Write-Step "PATH entry removed" "OK"
    } else {
        Write-Step "No PATH entry found" "SKIP"
    }
    Write-Host "│"
    if ($removed) {
        Write-Host "╰──────────────── $ProjectName uninstalled ───────────────╯"
    } else {
        Write-Host "╰────────── $ProjectName not found — nothing to do ───────╯"
    }
    return
}

# ── Install ─────────────────────────────────────────────────
Write-Host ""
Write-Host "╭──────────────── $ProjectName installer ────────────────╮"
Write-Host "│"

# 1. Detect architecture
Write-Step "Detecting system architecture..." -Status "INFO"
$archCode = (Get-CimInstance Win32_Processor).Architecture
$arch = switch ($archCode) {
    0  { "386"; Write-Step "Unsupported: x86 (32-bit)" "FAIL"; return }
    5  { "arm"; Write-Step "Unsupported: ARM (32-bit)" "FAIL"; return }
    9  { "amd64" }
    12 { "arm64" }
    default { Write-Step "Unknown architecture (code $archCode)" "FAIL"; return }
}
Write-Step "Architecture: $arch" "OK"

# 2. Detect OS
$os = "windows"
Write-Step "OS: $os" "OK"

# 3. Fetch latest version
Write-Step "Fetching latest version..." "INFO"
try {
    $versionUrl = "https://raw.githubusercontent.com/$GitHubUser/$ProjectName/main/.version"
    $version = (Invoke-RestMethod -Uri $versionUrl -ErrorAction Stop).Trim()
    if (-not $version) { throw "empty version" }
    Write-Step "Latest version: $version" "OK"
} catch {
    Write-Step "Failed to fetch version: $_" "FAIL"
    return
}

# 4. Build download URL
$binaryName = "$ProjectName-windows-$arch.exe"
$downloadUrl = "https://github.com/$GitHubUser/$ProjectName/releases/download/$version/$binaryName"
Write-Step "Download URL: $downloadUrl" "INFO"

# 5. Ensure bin directory
New-Item -ItemType Directory -Force -Path $BinDir | Out-Null

# 6. Download binary
$tmpPath = [System.IO.Path]::GetTempFileName()
try {
    Write-Step "Downloading $ProjectName $version..." -Status "INFO"
    $wc = New-Object System.Net.WebClient
    $wc.DownloadFile($downloadUrl, $tmpPath)
    
    if (-not (Test-Path $tmpPath) -or ((Get-Item $tmpPath).Length -eq 0)) {
        Write-Step "Download failed — empty or missing file" "FAIL"
        return
    }
    Write-Step "Downloaded ($((Get-Item $tmpPath).Length / 1KB) KB)" "OK"
} catch {
    Write-Step "Download failed: $_" "FAIL"
    if (Test-Path $tmpPath) { Remove-Item $tmpPath -Force }
    return
}

# 7. Install binary
Move-Item $tmpPath $BinaryPath -Force
Write-Step "Installed to $BinaryPath" "OK"

# 8. Add to PATH
if (Add-ToUserPath $BinDir) {
    Write-Step "Added to user PATH" "OK"
    Write-Step "Restart your terminal or run:" "INFO"
    Write-Step "  `$env:PATH = `"`$env:PATH;$BinDir`"" "INFO"
} else {
    Write-Step "Already in PATH" "SKIP"
}

# 9. Test
try {
    $test = & $BinaryPath --version 2>&1
    Write-Step "Verified: $test" "OK"
} catch {
    Write-Step "Verification failed: $_" "FAIL"
}

Write-Host "│"
Write-Host "╰──────────────── $ProjectName installed ─────────────────╯"
Write-Host ""

Write-Host "  Run '$ProjectName --help' to get started."
Write-Host "  To uninstall later, run this in PowerShell:"
Write-Host ""
Write-Host "    iex `"& { `$(Invoke-RestMethod https://raw.githubusercontent.com/$GitHubUser/$ProjectName/main/installer.ps1) } --selfuninstall`""
Write-Host ""
