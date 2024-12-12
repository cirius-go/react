package react

import (
	"fmt"
	"sync"
)

// Atom represents a single entity.
type Atom[E any] struct {
	mu    sync.Mutex
	muted bool

	currentState       string
	detectorChanges    map[string]func(ent E) bool
	onStateEvents      map[string][]func() error
	onTransitionEvents map[string][]func() error
}

// NewAtom creates a new Atom.
func NewAtom[E any]() *Atom[E] {
	return &Atom[E]{
		currentState:       "",
		detectorChanges:    map[string]func(ent E) bool{},
		onStateEvents:      map[string][]func() error{},
		onTransitionEvents: map[string][]func() error{},
	}
}

// RegisterState register state of atom.
func RegisterState[E any](a *Atom[E], name string, detector func(ent E) bool) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.muted {
		return
	}

	if a.currentState == "" {
		a.currentState = name
	}
	a.detectorChanges[name] = detector
}

func OnState[E any](a *Atom[E], name string, f func() error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.muted {
		return
	}

	a.onStateEvents[name] = append(a.onStateEvents[name], f)
}

func OnStates[E any](a *Atom[E], names []string, f func() error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.muted {
		return
	}

	for _, name := range names {
		a.onStateEvents[name] = append(a.onStateEvents[name], f)
	}
}

func tranKey(from string, to string) string {
	return fmt.Sprintf("%s->%s", from, to)
}

func OnTransitionState[E any](a *Atom[E], from string, to string, f func() error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.muted {
		return
	}

	key := tranKey(from, to)
	a.onTransitionEvents[key] = append(a.onTransitionEvents[key], f)
}

// React on Atom 'a' state.
func React[E any](a *Atom[E], ent E) (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.muted = true

	for state, detector := range a.detectorChanges {
		if detector(ent) {
			events := []func() error{}
			oldState := a.currentState
			a.currentState = state

			tk := tranKey(oldState, a.currentState)
			if evs, ok := a.onTransitionEvents[tk]; ok {
				events = append(events, evs...)
			}

			if evs, ok := a.onStateEvents[a.currentState]; ok {
				events = append(events, evs...)
			}

			for i, ev := range events {
				if err := ev(); err != nil {
					fmsg := fmt.Sprintf("onState:%s", state)
					lenTs := len(a.onStateEvents[a.currentState])
					if i > lenTs-1 {
						fmsg = fmt.Sprintf("onTransition:%s", tk)
					}

					return fmsg, err
				}
			}

			continue
		}
	}
	return "", nil
}
