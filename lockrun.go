/**
 * Open Source Initiative OSI - The MIT License (MIT):Licensing
 *
 * The MIT License (MIT)
 * Copyright (c) 2014 José Gil (josgilmo@gmail.com)
 * Copyright (c) 2014 Diego Campoy (manrash@gmail.com)
 *
 * Strongly based on Stephen J. Friedl lockrun (http://www.unixwiz.net/tools/)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
 * of the Software, and to permit persons to whom the Software is furnished to do
 * so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type params struct {
	showHelp     bool
	showVersion  bool
	lockFilename string
	verbose      bool
	quiet        bool
	wait         bool
	sleep        int
	retries      int
	maxtime      int
	commandInfo  []string
}

func parse() params {

	var p params

	flag.BoolVar(&p.showHelp, "help", false, "Show this brief help listing")
	flag.BoolVar(&p.showVersion, "version", false, "Report the version info")
	flag.BoolVar(&p.verbose, "verbose", false, "Show a bit more runtime debugging")
	flag.BoolVar(&p.quiet, "quiet", false, "Exit quietly (and with success) if locked")
	flag.StringVar(&p.lockFilename, "lockfile", "", "File for lock")
	flag.BoolVar(&p.wait, "wait", false, "Wait for lockfile to appear (else exit on lock)")

	flag.IntVar(&p.sleep, "sleep", 2, "Sleep for <T> seconds on each wait loop")
	flag.IntVar(&p.retries, "retries", 5, "Attempts <N> retries in each wait loop")
	flag.IntVar(&p.maxtime, "maxtime", 10, "Wait for at most <T> seconds for a lock, then exit")

	flag.Parse()
	p.commandInfo = flag.Args()

    if p.showHelp || len(p.commandInfo)==0 {
        showHelp() 
        os.Exit(1)
    }

    if p.showVersion {
        showVersion()
        os.Exit(1)
    }

	if p.lockFilename == "" {
		fmt.Println("--lockfile option can't be empty")
		os.Exit(1)
	}

	return p
}

func main() {

	p := parse()

	if p.showHelp {
		showHelp()
	} else if p.showVersion {
		showVersion()
	} else {
		tryLockAndRun(p)
	}
}

func tryLockAndRun(p params) {
	// Opening log file
	file, err := os.OpenFile(p.lockFilename, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Unable to write lockfile: %s", p.lockFilename)
	}

	// Trying to lock file
	attempts := 0
	for {
		attempts++
		err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
		if err == nil {
            if p.verbose == true { 
    			log.Printf("Locking...")
            }
			break
		}

		if p.wait == false {
			fmt.Printf("ERROR: cannot launch %s - run is locked", "TODO")
			os.Exit(1)
		} else {
            if p.verbose  {
    			log.Printf("Attempt %d failed - sleeping %d seconds", attempts, p.sleep)
            }
			time.Sleep(time.Duration(p.sleep) * time.Second)

			if attempts >= p.retries {
				fmt.Printf("ERROR: cannot launch %s - run is locked (after %d attempts", "TODO", attempts)
				os.Exit(1)
			}
		}
	}

	if err != nil {
		log.Fatalf("Locking error: %v", err)
	}

	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}

	command, err4 := exec.LookPath(p.commandInfo[0])
	if err4 != nil {
		command = p.commandInfo[0]
	}

	process, err2 := os.StartProcess(command, p.commandInfo, &procAttr)
	if err2 != nil {
		fmt.Printf("ERROR: %s\n", err2)
	} else {
		_, err3 := process.Wait()
		if err3 != nil {
			fmt.Printf("ERROR: %s\n", err3)
		}
	}

	log.Printf("Finish!")
}

func showHelp() {
	fmt.Println(`Usage: lockrun [options] -- command args...

    --help        Show this brief help listing
    --version     Report the version info
    --            Mark the end of lockrun options; command follows
    --verbose     Show a bit more runtime debugging
    --quiet       Exit quietly (and with success) if locked
    --lockfile=F  Specify lockfile as file <F>
    --wait        Wait for lockfile to appear (else exit on lock)

 Options with --wait:
    --sleep=T     Sleep for <T> seconds on each wait loop
    --retries=N   Attempt <N> retries in each wait loop
    --maxtime=T   Wait for at most <T> seconds for a lock, then exit
`)
}

func showVersion() {
	fmt.Println("v0.0.1")
}
