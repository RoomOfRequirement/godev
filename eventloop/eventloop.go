package eventloop

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"os"
	"sync"
	"time"
)

// Loop struct
type Loop struct {
	events    chan *Event
	listeners map[string][]*EventListener

	running bool
	logger  *zap.Logger
	wg      sync.WaitGroup
}

// EventListener struct
type EventListener struct {
	callback    Callback
	timeout     time.Duration // 0 = never
	expiredTime time.Time
	expired     bool
}

// NewLoop creates a new event loop
func NewLoop(bufferedEvents int, logLevel string) EventLoop {
	if bufferedEvents < 0 {
		panic("invalid buffered events size")
	}
	return &Loop{
		events:    make(chan *Event, bufferedEvents),
		listeners: make(map[string][]*EventListener),
		running:   false,
		logger:    newLogger(logLevel),
		wg:        sync.WaitGroup{},
	}
}

// Push pushes event into event loop
func (l *Loop) Push(event *Event) {
	l.events <- event
}

// Emit broadcasts event to its listener
func (l *Loop) Emit(eventName string, data interface{}) {
	l.events <- &Event{
		Name: eventName,
		Data: data,
	}
}

// Start starts the event loop
func (l *Loop) Start() {
	l.logger.Info("Loop starting")
	l.running = true
	l.wg.Add(1)
	grp, ctx := errgroup.WithContext(context.Background())

	go func() {
		l.logger.Info("Loop listening")
		var err error

		for l.running {
			select {
			case event := <-l.events:
				l.logger.Debug("Event received", zap.String("name", event.Name), zap.Any("data", event.Data))
				var removedListener []*EventListener

				for _, listener := range l.listeners[event.Name] {
					if listener.timeout != 0 && listener.expiredTime.Before(time.Now()) {
						l.logger.Info("listener expired", zap.String("name", event.Name))
						listener.expired = true
						removedListener = append(removedListener, listener)
						continue
					}
					grp.Go(func() error {
						return listener.callback(ctx, event.Data.(interface{}))
					})
				}

				if len(removedListener) > 0 {
					var listeners []*EventListener
					for _, listener := range l.listeners[event.Name] {
						if !listener.expired {
							listeners = append(listeners, listener)
						}
					}
					l.listeners[event.Name] = listeners
				}

				err = grp.Wait()

			default:
				// do nothing
			}

			if err != nil {
				l.logger.Error("Event error", zap.String("error", err.Error()))
				l.running = false
			}
		}

		l.wg.Done()
	}()
}

// Stop stops event loop
func (l *Loop) Stop() {
	l.logger.Info("Loop stopping")
	l.running = false
	l.logger.Debug("Waiting for loop stopped")
	close(l.events)
	l.wg.Wait()
	l.logger.Info("Loop stopped")
}

func (l *Loop) addEventListener(eventName string, callback Callback, timeout time.Duration) error {
	if callback == nil {
		return errors.New("invalid nil callback")
	}

	if _, found := l.listeners[eventName]; !found {
		l.listeners[eventName] = []*EventListener{}
	}
	l.listeners[eventName] = append(l.listeners[eventName], &EventListener{
		callback:    callback,
		timeout:     timeout,
		expiredTime: time.Now().Add(timeout),
		expired:     false,
	})

	return nil
}

// On adds callback to corresponding event
func (l *Loop) On(eventName string, callback Callback) error {
	return l.addEventListener(eventName, callback, 0)
}

// OnWithTimeout adds callback to corresponding event with timeout
func (l *Loop) OnWithTimeout(eventName string, callback Callback, timeout time.Duration) error {
	return l.addEventListener(eventName, callback, timeout)
}

func newLogger(level string) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	atom := zap.NewAtomicLevel()
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	defer logger.Sync()
	var l zapcore.Level
	switch level {
	case "debug":
		l = zap.DebugLevel
	case "info":
		l = zap.InfoLevel
	case "warn":
		l = zap.WarnLevel
	case "error":
		l = zap.ErrorLevel
	case "fatal":
		l = zap.FatalLevel
	case "panic":
		l = zap.PanicLevel
	default:
		l = zap.DebugLevel
		logger.Warn("Log level '" + level + "' not recognized. Default set to Debug")
	}
	atom.SetLevel(l)
	return logger
}
