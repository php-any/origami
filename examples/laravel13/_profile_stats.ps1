$ErrorActionPreference = 'Continue'
Stop-Process -Name laravel13 -Force -ErrorAction SilentlyContinue
Set-Location d:\gitcode.com\origami\examples\laravel13
$p = Start-Process -FilePath .\laravel13.exe -ArgumentList @('run','tests\origami\stats_mount_only_probe.php') -RedirectStandardOutput stats-m-out.txt -RedirectStandardError stats-m-err.txt -PassThru -NoNewWindow
$deadline = (Get-Date).AddSeconds(70)
while (-not $p.HasExited -and (Get-Date) -lt $deadline) {
    Start-Sleep -Seconds 5
    try { $p.Refresh() } catch {}
    Write-Host "alive=$(-not $p.HasExited) cpu=$($p.CPU)"
}
if (-not $p.HasExited) {
    Write-Host 'KILL'
    Stop-Process -Id $p.Id -Force
}
Write-Host 'OUT:'
Get-Content stats-m-out.txt -ErrorAction SilentlyContinue
Write-Host 'ERR:'
Get-Content stats-m-err.txt -ErrorAction SilentlyContinue | Select-Object -First 25
