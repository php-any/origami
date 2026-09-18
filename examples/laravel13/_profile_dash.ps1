$ErrorActionPreference = 'Continue'
Stop-Process -Name laravel13 -Force -ErrorAction SilentlyContinue
Set-Location d:\gitcode.com\origami\examples\laravel13
$p = Start-Process -FilePath .\laravel13.exe -ArgumentList @('run','tests\origami\dash_only_probe.php') -RedirectStandardOutput dash-slow-out.txt -RedirectStandardError dash-slow-err.txt -PassThru -NoNewWindow
$deadline = (Get-Date).AddSeconds(90)
while (-not $p.HasExited -and (Get-Date) -lt $deadline) {
    Start-Sleep -Seconds 5
    $p.Refresh()
    $sec = [int]((Get-Date) - $p.StartTime).TotalSeconds
    $mb = [int]($p.WorkingSet64 / 1MB)
    Write-Host "t=${sec}s cpu=$($p.CPU) ws=${mb}MB"
}
if (-not $p.HasExited) {
    Write-Host 'KILL'
    $p.Kill()
    Start-Sleep -Seconds 1
}
Write-Host 'OUT:'
Get-Content dash-slow-out.txt -ErrorAction SilentlyContinue
Write-Host 'ERR:'
Get-Content dash-slow-err.txt -ErrorAction SilentlyContinue | Select-Object -First 20
Write-Host "EXIT=$($p.ExitCode)"
