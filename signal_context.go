package daemon

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SignalHandlerFunc is the interface for signal handler functions.
type SignalHandlerWithContextFunc func(context context.Context, sig os.Signal) (err error)

// SetSigHandler sets handler for the given signals.
// SIGTERM has the default handler, he returns ErrStop.
func SetSigHandlerWithContext(handler SignalHandlerWithContextFunc, signals ...os.Signal) {
	for _, sig := range signals {
		handlersCtx[sig] = handler
	}
}

// ServeSignals calls handlers for system signals.
func ServeSignalsWithContext(ctx context.Context) (err error) {
	signals := make([]os.Signal, 0, len(handlersCtx))
	for sig := range handlersCtx {
		signals = append(signals, sig)
	}

	ch := make(chan os.Signal, 8)
	signal.Notify(ch, signals...)

	for sig := range ch {
		err = handlersCtx[sig](ctx, sig)
		if err != nil {
			break
		}
	}

	signal.Stop(ch)

	if err == ErrStop {
		err = nil
	}

	return
}

var handlersCtx = make(map[os.Signal]SignalHandlerWithContextFunc)

func init() {
	handlersCtx[syscall.SIGTERM] = sigtermDefaultHandlerWithContext
}

func sigtermDefaultHandlerWithContext(ctx context.Context, sig os.Signal) error {
	return ErrStop
}
