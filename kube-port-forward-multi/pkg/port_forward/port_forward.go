package port_forward

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Target struct {
	LocalPort  string
	Resource   string
	RemotePort string
}

type Options struct {
	Targets    []string
	Namespace  string
	Context    string
	Kubeconfig string
}

func parseTarget(raw string) (Target, error) {
	parts := strings.Split(raw, ":")
	if len(parts) != 3 {
		return Target{}, fmt.Errorf(
			"invalid target %q, expected <local-port>:<svc/name|pod/name>:<remote-port>",
			raw,
		)
	}
	return Target{
		LocalPort:  parts[0],
		Resource:   parts[1],
		RemotePort: parts[2],
	}, nil
}

func Run(opts Options) error {
	if _, err := exec.LookPath("kubectl"); err != nil {
		return fmt.Errorf("kubectl not found in PATH: %w", err)
	}

	targets := make([]Target, 0, len(opts.Targets))
	for _, raw := range opts.Targets {
		t, err := parseTarget(raw)
		if err != nil {
			return err
		}
		targets = append(targets, t)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nStopping port forwards...")
		cancel()
	}()

	var wg sync.WaitGroup
	for _, t := range targets {
		wg.Add(1)
		go func(t Target) {
			defer wg.Done()
			runForward(ctx, opts, t)
		}(t)
	}
	wg.Wait()

	return nil
}

func runForward(ctx context.Context, opts Options, t Target) {
	label := fmt.Sprintf("[%s -> %s:%s]", t.LocalPort, t.Resource, t.RemotePort)

	args := []string{"port-forward", t.Resource, fmt.Sprintf("%s:%s", t.LocalPort, t.RemotePort)}
	if opts.Namespace != "" {
		args = append(args, "-n", opts.Namespace)
	}
	if opts.Context != "" {
		args = append(args, "--context", opts.Context)
	}
	if opts.Kubeconfig != "" {
		args = append(args, "--kubeconfig", opts.Kubeconfig)
	}

	for ctx.Err() == nil {
		cmd := exec.CommandContext(ctx, "kubectl", args...)
		cmd.Stdout = &prefixWriter{prefix: label, w: os.Stdout}
		cmd.Stderr = &prefixWriter{prefix: label, w: os.Stderr}

		fmt.Printf("%s starting port-forward\n", label)
		err := cmd.Run()
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			fmt.Printf("%s port-forward exited: %v, retrying in 2s...\n", label, err)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}

// prefixWriter prefixes every line written to it with a label, so output from
// concurrent port-forwards stays distinguishable.
type prefixWriter struct {
	prefix string
	w      io.Writer
}

func (p *prefixWriter) Write(b []byte) (int, error) {
	for _, line := range strings.Split(strings.TrimRight(string(b), "\n"), "\n") {
		if line == "" {
			continue
		}
		fmt.Fprintf(p.w, "%s %s\n", p.prefix, line)
	}
	return len(b), nil
}
