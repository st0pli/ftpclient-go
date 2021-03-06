# ftpclient-go
Golang Ftp Client Library

## Installing

```
go get -u github.com/tsujimic/ftpclient-go
```

## Using the ftpclient-go
```
package main

import (
    "fmt"
    "flag"
    "log"
    "net"
    "strconv"
    "time"
    "github.com/tsujimic/ftpclient-go"
)

func main() {
    var host, user, pass, remote, local string
    var port int
    
    flag.StringVar(&host, "host", "", "target host name")
    flag.IntVar(&port, "port", 21, "tcp/ip port number")
    flag.StringVar(&user, "user", "", "login username")
    flag.StringVar(&pass, "pass", "", "login password")
    flag.StringVar(&remote, "remote", "", "remote file path")
    flag.StringVar(&local, "local", "", "local file path")
    flag.Parse()

    log.Println("Start")

    logger := NewDefaultLogger()
    cfg := NewConfig().WithLogger(logger)
    client := New(cfg)
    addr := net.JoinHostPort(host, strconv.Itoa(port))
    err := client.DialTimeout(addr, 30*time.Second)
    if err != nil {
        panic(err)
    }
    defer client.Quit()

    err = client.Login(user, pass)
    if err != nil {
        panic(err)
    }

    err = client.Type("I")
    if err != nil {
        panic(err)
    }

    client.SetPasv(false)
    start := time.Now()
    err = client.StorFile(local, remote)
    if err != nil {
        panic(err)
    }

    sec := (time.Now().Sub(start)).Seconds()    
    msg := fmt.Sprintf("Stopwatch : %f seconds\n", sec)
    log.Println(msg)
}
```
