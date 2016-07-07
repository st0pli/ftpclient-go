package main

import (
    "fmt"
    "flag"
    "log"
    "net"
    "strconv"
    "io"
    "os"
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
    addr := net.JoinHostPort(host, strconv.Itoa(port))
    client, err := ftpclient.Connect(addr)
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
    
    client.SetPasv(true)
    infos, err := client.Dir(remote)
    if err != nil {
        panic(err)
    }

    length := len(infos)
    if length == 0 {
        msg := fmt.Sprintf("file not found: %s\n", remote)
        panic(msg)
    }

    filesize := infos[0].Size()
    reader, err := client.RetrRequest(remote)
    if err != nil {
        panic(err)
    }
    defer reader.Close()
 
    file, err := os.Create(local)
    if err != nil {
        panic(err)
    }
    defer func() {
        if err := file.Close(); err != nil {
            panic(err)
        }
    }()
    
    start := time.Now()
    buf := make([]byte, 32*1024)
    for {
        nr, err := reader.Read(buf)
        if nr > 0 {
            nw, err := file.Write(buf[:nr])
            if err != nil {
                panic(err)
            }
            if nr != nw {
                panic(io.ErrShortWrite)
            }
        }
        if err == io.EOF {
            break
        }
        if err != nil {
            panic(err)
        }

    }

    now := time.Now()
    sec := (now.Sub(start)).Seconds()    
    transferbps := (float64(filesize) / sec) * 8
    fmt.Printf("Stopwatch : %f seconds, %f Mbit/s", sec, transferbps / 1048576)
}
