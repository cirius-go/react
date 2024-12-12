package main

import (
	"fmt"

	"github.com/cirius-go/react"
)

func main() {
	type Appt struct {
		ID     int
		Status int
	}

	var appt = &Appt{}
	at := react.NewAtom[*Appt]()

	react.RegisterState(at, "checkedIn", func(ent *Appt) bool {
		return ent.Status == 1
	})

	react.RegisterState(at, "processing", func(ent *Appt) bool {
		return ent.Status == 2
	})

	react.OnState(at, "checkedIn", func() error {
		fmt.Println("onState: checkedIn")
		return nil
	})

	react.OnState(at, "processing", func() error {
		fmt.Println("onState: processing")
		return nil
	})

	react.OnStates(at, []string{"checkedIn", "processing"}, func() error {
		fmt.Println("onState: checkedIn|processing")
		return nil
	})

	react.OnTransitionState(at, "checkedIn", "processing", func() error {
		fmt.Println("transition: checkedIn->processing")
		return nil
	})

	appt.Status = 1
	if f, err := react.React(at, appt); err != nil {
		panic(fmt.Sprintf("failed on %s, error %v", f, err))
	}

	if f, err := react.React(at, appt); err != nil {
		panic(fmt.Sprintf("failed on %s, error %v", f, err))
	}

	appt.Status = 2
	if f, err := react.React(at, appt); err != nil {
		panic(fmt.Sprintf("failed on %s, error %v", f, err))
	}

	appt.Status = 3
	if f, err := react.React(at, appt); err != nil {
		panic(fmt.Sprintf("failed on %s, error %v", f, err))
	}
}
