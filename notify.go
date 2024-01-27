package notify

import (
	"golang.org/x/sync/errgroup"
)

type Notifier interface {
	Send(string, string) error
}

type Notify struct {
	Disabled  bool
	notifiers []Notifier
}

func New() *Notify {
	return &Notify{
		Disabled:  false,
		notifiers: []Notifier{},
	}
}

func (n *Notify) useService(s Notifier) {
	if s != nil {
		n.notifiers = append(n.notifiers, s)
	}
}

func (n *Notify) UseService(service ...Notifier) {
	for _, s := range service {
		n.useService(s)
	}
}

func (n *Notify) Send(subject, message string) error {
	if n.Disabled {
		return nil
	}

	var eg errgroup.Group

	for _, service := range n.notifiers {
		if service != nil {
			s := service
			eg.Go(func() error {
				return s.Send(subject, message)
			})
		}
	}

	return eg.Wait()
}
