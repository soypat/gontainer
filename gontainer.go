//go:build !legacy

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type configuration struct {
	timeout     time.Duration
	chroot      string
	chdir       string
	verbose     int
	childInputs []string
}

func main() {
	var cfg configuration
	flag.StringVar(&cfg.chroot, "chrt", os.Getenv("GONTAINER_FS"), "Where to chroot to. Should contain a linux filesystem. Alpine is recommended. GONTAINER_FS environment is default if not set")
	flag.StringVar(&cfg.chdir, "chdir", "/usr", "Initial chdir executed when running container")
	flag.DurationVar(&cfg.timeout, "timeout", 0, "Timeout before ending program. If 0 then never ends")
	flag.IntVar(&cfg.verbose, "v", 0, "If v is set container command output not suppressed. Debugging purposes")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		flag.Usage()
		cfg.fatalf("need 2 arguments. got: %v", args)
	}
	if cfg.chroot == "" {
		flag.Usage()
		cfg.fatalf("chroot flag is required (--chrt flag). got flags/args: %v", os.Args)
	}
	if args[0] != "child" {
		flag.VisitAll(func(f *flag.Flag) {
			cfg.childInputs = append(cfg.childInputs, "-"+f.Name, f.Value.String())
		})
	}
	var ctx context.Context
	var cancel func()
	if cfg.timeout == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), cfg.timeout)
	}
	cleanup := func() {
		// Not sure how to fix console breaking on 1s timeout. Following is not working.
		// os.Stdin = os.NewFile(uintptr(syscall.Stdin), "/dev/stdin")
		// os.Stdout = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
		// os.Stderr = os.NewFile(uintptr(syscall.Stderr), "/dev/stderr")
		cancel()
	}
	defer cleanup()
	// Ctrl+C interrupt channel
	intChan := make(chan os.Signal, 1)
	signal.Notify(intChan, os.Interrupt)
	go func() {
		<-intChan
		cancel()
		// Exit program forcefully on second Interrupt receive.
		<-intChan
		cleanup()
		cfg.fatalf("forced program exit")
	}()
	var err error
	switch args[0] {
	case "run":
		err = run(ctx, args, cfg)
	case "child":
		err = child(ctx, args, cfg)
	default:
		cfg.fatalf("unknown argument " + args[0])
	}
	cleanup()
	if err != nil {
		cfg.fatalf(err.Error())
	}
	cfg.infof("%s finished", args[0])
}

func run(ctx context.Context, args []string, cfg configuration) error {
	cfg.infof("run as [%d] : running %v", os.Getpid(), args[1:])
	lst := append(append(cfg.childInputs, "child"), args[1:]...)

	cmd := exec.CommandContext(ctx, "/proc/self/exe", lst...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	return cmd.Run()
}

func child(ctx context.Context, args []string, cfg configuration) error {
	cfg.infof("child as [%d]: chrt: %s,  chdir:%s, args: ", os.Getpid(), cfg.chroot, cfg.chdir, args[1:])

	cfg.must(syscall.Sethostname([]byte("container")), "setting hostname")
	cfg.must(syscall.Chroot(cfg.chroot), "during chroot")
	syscall.Mkdir(cfg.chdir, 0600)

	// initial chdir is necessary so dir pointer is in chroot dir when proc mount is called
	cfg.must(syscall.Chdir("/"), "during chdir before proc mount")
	cfg.must(syscall.Mount("proc", "proc", "proc", 0, ""), "error in proc mount")
	cfg.must(syscall.Chdir(cfg.chdir), "during chdir")

	cmd := exec.CommandContext(ctx, args[1], args[2:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	defer syscall.Unmount("/proc", 0)
	return cmd.Run()
}

func (cfg *configuration) infof(format string, args ...interface{})  { cfg.logf("inf", format, args) }
func (cfg *configuration) errorf(format string, args ...interface{}) { cfg.logf("err", format, args) }
func (cfg *configuration) fatalf(format string, args ...interface{}) {
	cfg.verbose = 1
	cfg.logf("fat", format, args)
	os.Exit(1)
}
func (cfg *configuration) logf(tag, format string, args []interface{}) {
	if cfg.verbose != 0 {
		msg := fmt.Sprintf(format, args...)
		if args == nil {
			msg = fmt.Sprintf(format)
		}
		fmt.Println(fmt.Sprintf("[%s] %s", tag, msg))
	}
}
func (cfg *configuration) must(err error, msg string) {
	if err != nil {
		cfg.fatalf(msg + ": " + err.Error())
	}
}
