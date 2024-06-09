#Requires -RunAsAdministrator

$CUR_PATH = Get-Location
$APP_DIR = ".sherry"
$CONFIG_PATH = "$HOME\$APP_DIR"
$BIN_PATH = "$CONFIG_PATH\bin"

Set-Location $HOME
mkdir $APP_DIR
Set-Location $CONFIG_PATH

Remove-Item -Recurse -Force $BIN_PATH
New-Item -ItemType Directory -Path $BIN_PATH -Force
Invoke-WebRequest https://github.com/sherry-sync/cli/releases/latest/download/shr.exe -OutFile "$BIN_PATH\shr.exe"
Invoke-WebRequest https://github.com/sherry-sync/cli/releases/latest/download/sherry-windows-install.ps1 -OutFile "$BIN_PATH\shr-update.ps1"
Invoke-WebRequest https://github.com/sherry-sync/demon/releases/latest/download/sherry-demon.exe -OutFile "$BIN_PATH\sherry-demon.exe"
if ([Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User) -split ';' -notcontains $BIN_PATH)
{
    [Environment]::SetEnvironmentVariable(
            "Path",
            [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User) + $BIN_PATH,
            [EnvironmentVariableTarget]::User
    )
}

$TARGET = [System.EnvironmentVariableTarget]::User

[System.Environment]::SetEnvironmentVariable("SHERRY_CONFIG_PATH", $CONFIG_PATH, $TARGET)
[System.Environment]::SetEnvironmentVariable("SHERRY_API_URL", "http://localhost:3000", $TARGET)
[System.Environment]::SetEnvironmentVariable("SHERRY_SOCKET_URL", "http://localhost:3001", $TARGET)

$STARTUP_PATH = "HKLM:\SOFTWARE\Microsoft\Windows\CurrentVersion\Run"
$STARTUP_NAME = "sherry-demon"

$STARTUP_SCRIPT = '
$JOB_NAME = "SherryStartup"
Start-Job -Name $JOB_NAME -ScriptBlock {
    $TEMP = [System.Environment]::GetEnvironmentVariable("TEMP", [System.EnvironmentVariableTarget]::User)
    ' + $BIN_PATH + '\shr.exe -c "' + $CONFIG_PATH + '" service start -y >> "$TEMP/SherryStartup.txt"
}
Wait-Job -Name $JOB_NAME
'

Write-Output $STARTUP_SCRIPT > "$BIN_PATH\sherry_startup.ps1"
if (-Not (Get-ItemProperty -Path $STARTUP_PATH -Name $STARTUP_NAME))
{
    $confirmation = Read-Host "Do you want to start shr service automatically? (Y/n)"
    if ($confirmation -ne "n")
    {
        New-ItemProperty -Path $STARTUP_PATH -Name $STARTUP_NAME -Value "powershell.exe ""$BIN_PATH\sherry_startup.ps1""" -PropertyType "String"
    }
}

RefreshEnv

powershell.exe -Command "$BIN_PATH\shr.exe" -c $CONFIG_PATH service stop
(Start-Process -FilePath "$BIN_PATH\sherry-demon.exe" -NoNewWindow -PassThru -WorkingDirectory $BIN_PATH -ArgumentList @("-c", $CONFIG_PATH, "-s")).Id > "$CONFIG_PATH/pid"

Set-Location $CUR_PATH
