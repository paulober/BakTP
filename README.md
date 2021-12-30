# BakTP

Simple SFTP backup tool for files.

#### config.example.json
---
Contains an example how to backup a database.

### This application can be added to `crontab -e` for a timed backup
`6 */12 * * * /home/pi/baktp -c /home/pi/baktp-config.json`
This for example makes a backup every 12h + 6min (0:06am and 12:06am)
