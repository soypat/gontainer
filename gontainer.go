package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

var ( // flags
	chroot, chdir string
	loud          bool
)

// cleanup tasks
var cntcmd, runcmd *exec.Cmd
var wg = &sync.WaitGroup{}
var shutDownChan chan os.Signal
var args []string

// Flag and argument parsing
func init() {
	pflag.StringVarP(&chroot, "chrt", "", "", "Where to chroot to. Should contain a linux filesystem. Alpine is recommended. GONTAINER_FS environment is default if not set")
	pflag.StringVarP(&chdir, "chdr", "", "/usr", "Initial chdir executed when running container")
	pflag.BoolVar(&loud, "loud", false, "Suppresses not container output. Debugging purposes")
	pflag.Parse()
	args = pflag.Args()
	if chroot == "" {
		chroot = os.Getenv("GONTAINER_FS")
		if chroot == "" {
			fatalf("chroot (--chrt flag) is required. got args: %v", args)
		}
	}
	if len(args) < 2 {
		fatalf("too few arguments. got: %v", args)
	}
}

func main() {
	switch args[0] {
	case "run":
		run()
	case "child":
		wg.Add(1)
		go child()
		c := make(chan os.Signal, 1)
		shutDownChan = make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		select {
		case <-c:
		case <-shutDownChan:
		}
		wg.Wait()
	default:
		panic("bad command")
	}

}

func run() {
	infof("run as [%d] : running %v", os.Getpid(), args[1:])
	lst := append([]string{"--chrt", chroot, "--chdr", chdir, "child"}, args[1:]...)
	runcmd = exec.Command("/proc/self/exe", lst...)
	runcmd.Stdin = os.Stdin
	runcmd.Stdout = os.Stdout
	runcmd.Stderr = os.Stderr
	runcmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	runcmd.Run()
}

// This child function runs a command in a containerized
// linux filesystem so it can't hurt you.
func child() {
	defer cleanup()
	infof("child as [%d]: chrt: %s,  chdir:%s", os.Getpid(), chroot, chdir)
	infof("running %v", args[1:])
	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot(chroot), "error in 'chroot ", chroot+"'")
	syscall.Mkdir(chdir, 0600)

	// initial chdir is necessary so dir pointer is in chroot dir when proc mount is called
	must(syscall.Chdir("/"), "error in 'chdir /'")
	must(syscall.Mount("proc", "proc", "proc", 0, ""), "error in proc mount")
	must(syscall.Chdir(chdir), "error in 'chdir ", chdir+"'")
	cntcmd = exec.Command(args[1], args[2:]...)
	cntcmd.Stdin = os.Stdin
	cntcmd.Stdout = os.Stdout
	cntcmd.Stderr = os.Stderr
	must(cntcmd.Run(), fmt.Sprintf("run %v return error", args[1:]))
	syscall.Unmount("/proc", 0)
}

func cleanup() {
	if cntcmd != nil {
		cntcmd.Process.Signal(os.Interrupt)
		time.Sleep(time.Millisecond * 1)
		cntcmd.Process.Signal(os.Kill)
		syscall.Unmount("/proc", 0)
	}
	shutDownChan <- os.Interrupt
	wg.Done()
}

func must(err error, s ...string) {
	if err != nil {
		loud = true
		errorf("%s : %v", err, s)
		os.Exit(1)
	}
}

func infof(format string, args ...interface{}) { logf("inf", format, args) }

//func printf(format string, args ...interface{}) { logf("prn", format, args) }
func errorf(format string, args ...interface{}) { logf("err", format, args) }
func fatalf(format string, args ...interface{}) { loud = true; logf("fat", format, args); os.Exit(1) }
func logf(tag, format string, args []interface{}) {
	if loud {
		msg := fmt.Sprintf(format, args...)
		if args == nil {
			msg = fmt.Sprintf(format)
		}
		fmt.Println(fmt.Sprintf("[%s] %s", tag, msg))
	}
}
