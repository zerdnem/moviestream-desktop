# MovieStream Desktop Release Script (PowerShell)
# This script helps you create a new release automatically

$ErrorActionPreference = "Stop"

Write-Host "MovieStream Desktop - Release Creator" -ForegroundColor Green
Write-Host "==========================================" -ForegroundColor Green
Write-Host ""

# Check if git is clean
$gitStatus = git status -s
if ($gitStatus) {
    Write-Host "Error: Working directory is not clean!" -ForegroundColor Red
    Write-Host "Please commit or stash your changes first."
    exit 1
}

# Get current branch
$branch = git branch --show-current
Write-Host "Current branch: $branch" -ForegroundColor Yellow

# Get the last tag
try {
    $lastTag = git describe --tags --abbrev=0 2>$null
} catch {
    $lastTag = "v0.0.0"
}
Write-Host "Last release: $lastTag" -ForegroundColor Yellow
Write-Host ""

# Ask for version
$version = Read-Host "Enter new version (format: v1.0.0)"

# Validate version format
if ($version -notmatch '^v\d+\.\d+\.\d+$') {
    Write-Host "Error: Invalid version format!" -ForegroundColor Red
    Write-Host "Version must be in format: v1.0.0"
    exit 1
}

# Check if tag already exists
$tagExists = git tag -l $version
if ($tagExists) {
    Write-Host "Error: Tag $version already exists!" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Creating release " -NoNewline
Write-Host $version -ForegroundColor Green
Write-Host "This will:"
Write-Host "  1. Create and push a git tag"
Write-Host "  2. Trigger GitHub Actions to:"
Write-Host "     - Create a GitHub release"
Write-Host "     - Build Windows .exe"
Write-Host "     - Build Linux binary"
Write-Host "     - Build macOS binaries (Intel + ARM)"
Write-Host "     - Attach all binaries to the release"
Write-Host ""
$confirm = Read-Host "Continue? (y/n)"

if ($confirm -ne "y" -and $confirm -ne "Y") {
    Write-Host "Cancelled."
    exit 0
}

# Create and push tag
Write-Host ""
Write-Host "Creating tag..."
git tag -a $version -m "Release $version"

Write-Host "Pushing tag to remote..."
git push origin $version

Write-Host ""
Write-Host "[SUCCESS]" -ForegroundColor Green
Write-Host ""
Write-Host "Release tag created and pushed!"
Write-Host ""
Write-Host "GitHub Actions is now:"
Write-Host "  -> Creating the release"
Write-Host "  -> Building binaries for all platforms"
Write-Host "  -> Uploading binaries to the release"
Write-Host ""
Write-Host "View progress at:"
Write-Host "https://github.com/zerdnem/moviestream-desktop/actions" -ForegroundColor Yellow
Write-Host ""
Write-Host "Release will be available at:"
Write-Host "https://github.com/zerdnem/moviestream-desktop/releases/tag/$version" -ForegroundColor Yellow
Write-Host ""
Write-Host "This usually takes 5-10 minutes to complete."

