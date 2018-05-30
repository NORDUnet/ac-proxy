
## Building

```
GOOS=windows GOARCH=amd64 go build -o ac-proxy.exe
```

## Installing as a windows service

1. Download [WinSW](https://github.com/kohsuke/winsw/releases)
2. Rename `WinSW.exe` to `ac-proxy-svc.exe`
3. Place `ac-proxy.exe` in the same directory
4. Copy over `ac-proxy-svc.xml` and edit the `-host` argument to fit the server the service is running on
5. Run `ac-proxy-svc.exe install`
6. Start the service in the Service manager
