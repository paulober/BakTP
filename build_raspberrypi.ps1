# !does not work anymore
# instead use SystemPropertiesAdvanced.exe and add 
# the values below in the per user environment vars

# Set-Variable GOOS=linux
# Set-Variable GOARCH=arm
# Set-Variable GOARM=7

### END old setup ###

### BEGIN new setup ###

# Alternative would be "Set-Location Env:" and 
# then "Set-Content -Path Test -Value 'Test value'"
# or
# [System.Environment]::SetEnvironmentVariable("GOOS", "linux", [System.EnvironmentVariableTarget]::Process)
Set-Item -Path Env:GOOS -Value "linux"
Set-Item -Path Env:GOARCH -Value "arm"
Set-Item -Path Env:GOARM -Value 7

go build

Remove-Item -Path Env:GOOS
Remove-Item -Path Env:GOARCH
Remove-Item -Path Env:GOARM
