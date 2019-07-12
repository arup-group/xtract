### Xtract 
Runs backups on gpj file (gINT database) at 15min intervals
After 32 copies are made we move the latest to a timestamped subfolder


### Compile
env GOOS=windows GOARCH=amd64 go install -v github.com/ArupAus/xtract