package eventloop

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Loop struct
type Loop struct {
	events    chan *Event
	listeners *eventListenerMap

	running uint32
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

type eventListenerMap struct {
	lock sync.RWMutex
	m    map[string][]*EventListener
}

// NewLoop creates a new event loop
func NewLoop(bufferedEvents int, logLevel string) EventLoop {
	if bufferedEvents < 0 {
		panic("invalid buffered events size")
	}
	return &Loop{
		events: make(chan *Event, bufferedEvents),
		listeners: &eventListenerMap{
			lock: sync.RWMutex{},
			m:    make(map[string][]*EventListener),
		},
		running: 0,
		logger:  newLogger(logLevel),
		wg:      sync.WaitGroup{},
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
	if atomic.LoadUint32(&l.running) != 0 {
		l.logger.Info("Loop already running")
		return
	}
	l.logger.Info("Loop starting")
	atomic.StoreUint32(&l.running, 1)
	l.wg.Add(1)

	go func() {
		l.logger.Info("Loop listening")
		var err error
		grp, ctx := errgroup.WithContext(context.Background())

		for atomic.LoadUint32(&l.running) != 0 {
			select {
			case <-ctx.Done():
				err = ctx.Err()
			case event := <-l.events:
				l.logger.Debug("Event received", zap.String("name", event.Name), zap.Any("data", event.Data))
				var removedListener []*EventListener

				l.listeners.lock.RLock()
				for _, listener := range l.listeners.m[event.Name] {
					if listener.timeout != 0 && listener.expiredTime.Before(time.Now()) {
						l.logger.Info("listener expired", zap.String("name", event.Name))
						listener.expired = true
						removedListener = append(removedListener, listener)
						continue
					} else {
						// closure
						listener, event := listener, event
						grp.Go(func() error {
							return listener.callback(ctx, event.Data)
						})
					}
				}
				l.listeners.lock.RUnlock()

				if len(removedListener) > 0 {
					var listeners []*EventListener
					l.listeners.lock.RLock()
					for _, listener := range l.listeners.m[event.Name] {
						if !listener.expired {
							listeners = append(listeners, listener)
						}
					}
					l.listeners.lock.RUnlock()
					l.listeners.lock.Lock()
					l.listeners.m[event.Name] = listeners
					l.listeners.lock.Unlock()
				}

				err = grp.Wait()

			default:
				// do nothing
			}

			if err != nil {
				l.logger.Error("Event error", zap.String("error", err.Error()))
				atomic.StoreUint32(&l.running, 0)
			}
		}

		l.wg.Done()
	}()
}

// Stop stops event loop
func (l *Loop) Stop() {
	l.logger.Info("Loop stopping")
	atomic.StoreUint32(&l.running, 0)
	l.logger.Debug("Waiting for loop stopped")
	close(l.events)
	l.wg.Wait()
	l.logger.Info("Loop stopped")
}

func (l *Loop) addEventListener(eventName string, callback Callback, timeout time.Duration) error {
	if callback == nil {
		return errors.New("invalid nil callback")
	}

	l.listeners.lock.Lock()
	if _, found := l.listeners.m[eventName]; !found {
		l.listeners.m[eventName] = []*EventListener{}
	}
	l.listeners.m[eventName] = append(l.listeners.m[eventName], &EventListener{
		callback:    callback,
		timeout:     timeout,
		expiredTime: time.Now().Add(timeout),
		expired:     false,
	})
	l.listeners.lock.Unlock()

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
