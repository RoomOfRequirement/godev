package fsm

import (
	"fmt"
	"testing"
)

//states
const InitialState = "Initial"
const AwaitFromState = "AwaitFrom"
const AwaitToState = "AwaitTo"
const DoneState = "Done"

//messages
type Transfer struct {
	source chan int
	target chan int
	amount int
}

const Done = "Done"
const Failed = "Failed"

//data
type WireTransferData struct {
	source chan int
	target chan int
	amount int
	client *FSM
}

func newWireTransfer(transferred chan bool) *FSM {
	wt := NewFSM()

	wt.Init(InitialState, nil, nil)

	wt.When(InitialState)(
		func(event *Event) *NextState {
			transfer, transferOk := event.Message.(*Transfer)
			if transferOk && event.Data == nil {
				transfer.source <- transfer.amount
				return wt.Goto(AwaitFromState).With(
					&WireTransferData{transfer.source, transfer.target, transfer.amount, wt},
				)
			}
			return wt.DefaultStateFunction()(event)
		})

	wt.When(AwaitFromState)(
		func(event *Event) *NextState {
			data, dataOk := event.Data.(*WireTransferData)
			if dataOk {
				switch event.Message {
				case Done:
					data.target <- data.amount
					return wt.Goto(AwaitToState)
				case Failed:
					//go data.client.Send(Failed)
					go data.client.FireEvent(NewEvent(Failed, data, func(event *Event) {
						fmt.Printf("Event launched: %v\n", event)
					}, func(event *Event) {
						fmt.Printf("Event completed: %v\n", event)
					}))
					return wt.Stay()
				}
			}
			return wt.DefaultStateFunction()(event)
		})

	wt.When(AwaitToState)(
		func(event *Event) *NextState {
			data, dataOk := event.Data.(*WireTransferData)
			if dataOk {
				switch event.Message {
				case Done:
					transferred <- true
					return wt.Stay()
				case Failed:
					go data.client.Stay()
				}
			}
			return wt.DefaultStateFunction()(event)
		})
	return wt
}

func TestNewFSM(t *testing.T) {
	transferred := make(chan bool)

	wireTransfer := newWireTransfer(transferred)

	if wireTransfer.DefaultStateFunction() == nil {
		wireTransfer.SetDefaultHandler(nil)
	}

	wireTransfer.AddTransitionFunction(InitialState, AwaitFromState, func(from, to State) {
		fmt.Println(from, to)
	})

	transfer := &Transfer{
		source: make(chan int),
		target: make(chan int),
		amount: 30,
	}

	source := func() {
		withdrawAmount := <-transfer.source
		fmt.Printf("Withdrawn from source account: %d\n", withdrawAmount)
		// wireTransfer.Send(Done)
		wireTransfer.FireEvent(NewEvent(Done, wireTransfer.CurrentData(), func(event *Event) {
			fmt.Printf("Event launched: %v\n", event)
		}, func(event *Event) {
			fmt.Printf("Event completed: %v\n", event)
		}))
	}

	target := func() {
		topupAmount := <-transfer.target
		fmt.Printf("ToppedUp target account: %d\n", topupAmount)
		wireTransfer.Send(Done)
	}

	go source()
	go target()

	go wireTransfer.Send(transfer)

	if done := <-transferred; !done {
		panic("Something went wrong")
	}

	fmt.Println(wireTransfer.CurrentState())

	fmt.Println("DONE")
}
