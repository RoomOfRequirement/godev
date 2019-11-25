package fsm

import (
	"fmt"
	"sync"
)

// FSM for finite state machine
//	https://en.wikipedia.org/wiki/Finite-state_machine
//	http://erlang.org/documentation/doc-7.3/doc/design_principles/fsm.html
//	http://cloudmark.github.io/FSM/
//	reference: https://github.com/dyrkin/fsm/blob/master/fsm.go
//	https://github.com/arunma/AkkaFSM/blob/master/src/main/scala/me/rerun/akka/fsm/CoffeeMachine.scala
//	State(S) x Event(E) -> Actions (A), State(S')
type FSM struct {
	initialState State
	currentState State

	initialData Data
	currentData Data

	stateFunctions      map[State]StateFunction
	transitionFunctions map[transitionFunctionKey]func(from, to State)

	mutex *sync.Mutex

	defaultStateFunction StateFunction
}

type transitionFunctionKey struct {
	from State
	to   State
}

// State type
type State string

// Data type
type Data interface{}

// StateFunction type
type StateFunction func(*Event) *NextState

// NextState struct
type NextState struct {
	state State
	data  Data
}

// Event struct
type Event struct {
	Message interface{}
	Data    Data

	BeforeCallback, AfterCallback func(*Event)
}

// NewEvent creates a new event
func NewEvent(message interface{}, data Data, before, after func(*Event)) *Event {
	return &Event{
		Message:        message,
		Data:           data,
		BeforeCallback: before,
		AfterCallback:  after,
	}
}

// NewFSM creates a new finite state machine
func NewFSM() *FSM {
	return &FSM{
		initialState:        "",
		currentState:        "",
		stateFunctions:      map[State]StateFunction{},
		transitionFunctions: map[transitionFunctionKey]func(from State, to State){},
		mutex:               &sync.Mutex{},
		defaultStateFunction: func(*Event) *NextState {
			panic("default state transfer function undefined")
		},
	}
}

// Init initialize FSM
func (fsm *FSM) Init(state State, data Data, defaultStateFunction StateFunction) {
	fsm.initialState = state
	fsm.currentState = state
	fsm.initialData = data
	fsm.currentData = data
	fsm.defaultStateFunction = defaultStateFunction
	fsm.stateFunctions[state] = defaultStateFunction
}

// When add state function to state
func (fsm *FSM) When(state State) func(StateFunction) *FSM {
	return func(f StateFunction) *FSM {
		fsm.stateFunctions[state] = f
		return fsm
	}
}

// DefaultStateFunction returns default state function
func (fsm *FSM) DefaultStateFunction() StateFunction {
	return fsm.defaultStateFunction
}

// SetDefaultHandler sets default state function
func (fsm *FSM) SetDefaultHandler(defaultStateFunction StateFunction) {
	fsm.defaultStateFunction = defaultStateFunction
}

func makeTransition(fsm *FSM, nextState *NextState) {
	if f, found := fsm.transitionFunctions[transitionFunctionKey{
		from: fsm.currentState,
		to:   nextState.state,
	}]; found {
		f(fsm.currentState, nextState.state)
	}
	fsm.currentState = nextState.state
	fsm.currentData = nextState.data
}

// Send sends message
func (fsm *FSM) Send(message interface{}) {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()
	currentState := fsm.currentState
	stateFunction := fsm.stateFunctions[currentState]
	nextState := stateFunction(NewEvent(message, fsm.currentData, nil, nil))
	makeTransition(fsm, nextState)
}

// Goto goes to state
func (fsm *FSM) Goto(state State) *NextState {
	if _, ok := fsm.stateFunctions[state]; ok {
		return &NextState{state: state, data: fsm.currentData}
	}
	panic(fmt.Sprintf("Unknown state: %q", state))
}

// Stay stays current state
func (fsm *FSM) Stay() *NextState {
	return &NextState{state: fsm.currentState, data: fsm.currentData}
}

// With changes state data
func (ns *NextState) With(data Data) *NextState {
	ns.data = data
	return ns
}

// FireEvent fires event
func (fsm *FSM) FireEvent(event *Event) {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()

	if event.BeforeCallback != nil {
		event.BeforeCallback(event)
	}
	nextState := fsm.stateFunctions[fsm.currentState](event)
	makeTransition(fsm, nextState)
	if event.AfterCallback != nil {
		event.AfterCallback(event)
	}
}

// AddTransitionFunction add transition function
func (fsm *FSM) AddTransitionFunction(from, to State, f func(from, to State)) {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()
	fsm.transitionFunctions[transitionFunctionKey{
		from: from,
		to:   to,
	}] = f
}

// CurrentState returns current state
func (fsm *FSM) CurrentState() State {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()
	return fsm.currentState
}

// CurrentData returns current data
func (fsm *FSM) CurrentData() Data {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()
	return fsm.currentData
}
