package main
 
import (
//    "bytes"
//    "log"
//	"bufio"
	"flag"
//	"strconv"
    "syscall"
	"fmt"
    "os"
    "C"
//	"log"
    "os/exec"
)


func main() {

    lockfile := flag.String("lockfile", "", "File for lock")

    maxtime := flag.Int("maxtime", 60, "Time for execute the process in secods") 
   
    wait := flag.Bool("wait", false, "Time for execute the process in secods") 

    sleep := flag.Int("sleep", 60, "Time for wait in seconds if the var wait is true") 

    retries := flag.Int("retries", 60, "Time for wait in seconds if the var wait is true") 

    verbose := flag.Bool("verbose", false, "Verbose option") 

    quiet := flag.Bool("quiet", false, "Quiet option") 

    flag.Parse() 

    attemtps := 0
    wait_for_lock := false
    

    if *lockfile == "" {
       fmt.Println("--lockfile option can't be empty")
       os.Exit(1) 
    }

    lfd, err := os.OpenFile("/tmp/ttt.txt", os.O_CREATE|os.O_RDWR, 0666 )  
    if err!=nil {

    }

    lfd.WriteString("locking")
    for {
        if syscall.Flock(int(lfd.Fd()), syscall.F_TLOCK) == nil {
            break
        }
        attemtps = attemtps+1
        if !wait_for_lock {
            if *quiet==true {
                panic("Bye")
            } else {
                fmt.Println("ERROR: cannot launch  - run is locked");
            }
        }
    }
// https://gist.github.com/wofeiwo/3634357
ret, ret2, err := syscall.RawSyscall(syscall.SYS_FORK, 0, 0, 0)
    if ret > 0 {
        cmd := exec.Command("sleep", "1")
        err = cmd.Run()
    }
    fmt.Println("pid", ret)
    fmt.Println(ret2)
/*
    arr := []string{"/tmp/ppp.ttt"}
    var procAttr syscall.ProcAttr
fmt.Println(procAttr)
    procAttr.Files = []*os.File{nil, nil, nil}
    pid, err := syscall.ForkExec("/bin/touch", arr, &procAttr )
    w, err := syscall.Wait4(pid,nil, 0, nil)
fmt.Println(w)
    if pid == 0 {
fmt.Println("Parent process")
    } else {

    }

fmt.Println(pid)
*/
fmt.Println(err)
    /*

    fmt.Printf("hello, world %s \n", *lockfile)
    cmd := exec.Command("sleep", "1")
    var out bytes.Buffer
    cmd.Stdout = &out
    err = cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
*/

//     fmt.Printf("in all caps: %s\n", out.String())
    fmt.Println("maxtime : ", maxtime)
    fmt.Println("wait : ", wait)
    fmt.Println("sleep : ", sleep)
    fmt.Println("retries : ", retries)
    fmt.Println("verbose : ", verbose)
    fmt.Println("quiet : ", quiet)
}
