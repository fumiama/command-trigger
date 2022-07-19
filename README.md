# command-trigger
Just trigger one command and get its result by HMACed UDP

## usage
```powershell
  -D    debug-level log output
  -d int
        max valid time diff (default 5)
  -e string
        execute this command on triggered
  -g    generate a random key and exit
  -h    display this help
  -i int
        min execute interval (default 10)
  -k string
        64 bytes hmac key in base16384 format
  -m string
        send additional message
  -t string
        send/recv trigger to this addr:port
  -w uint
        max wait seconds for reply of executer (default 16)
```

## example
> Windows
- trigger
```powershell
go run main.go main_windows.go -t 127.0.0.1:8000 -k "抿淀檆健" -m "hello world"
[INFO] send trigger to 127.0.0.1:8000 : hello world
[INFO] 127.0.0.1:8000 reply: 123
```
- executer
```powershell
go run main.go main_windows.go -t 127.0.0.1:8000 -k "抿淀檆健" -e "cmd /c echo 123" 
[INFO] 127.0.0.1:56626 triggered with message: hello world
[INFO] get result: 123
```
> unix
- trigger
```powershell
go run main.go -t 127.0.0.1:8000 -k "抿淀檆健" -m "hello world"
[INFO] send trigger to 127.0.0.1:8000 : hello world
[INFO] 127.0.0.1:8000 reply: 123
```
- executer
```powershell
go run main.go -t 127.0.0.1:8000 -k "抿淀檆健" -e "cmd /c echo 123" 
[INFO] 127.0.0.1:56626 triggered with message: hello world
[INFO] get result: 123
```
