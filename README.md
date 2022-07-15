# command-trigger
Just trigger one command and get its result by HMACed UDP

## example
- trigger
```bash
go run main.go main_windows.go -D -t 127.0.0.1:8000 -k "抿淀檆健"  -m "hello world"
[DEBUG] send trigger to 127.0.0.1:8000 : hello world
[INFO] 127.0.0.1:8000 reply: 123
```
- executer
```bash
go run main.go main_windows.go -t 127.0.0.1:8000 -e "cmd /c echo 123" -D -k "抿淀檆健"  
[INFO] 127.0.0.1:56626 triggered wit
h message: hello world
[DEBUG] exec cmd: C:\WINDOWS\system3
2\cmd.exe /c echo 123
[DEBUG] get result: 123
```
