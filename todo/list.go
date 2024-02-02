package todo

import "errors"

type Item struct {
	Title  string
	IsDone bool
	Id     int
}

type List struct {
	Items  []Item
	nextId int
}

func NewList() List {
	return List{}
}

func (l *List) Add(title string) {
	if len(title) == 0 {
		return
	}
	l.Items = append(l.Items, Item{title, false, l.nextId})
	l.nextId++
}

func (l *List) Toggle(id int) error {
	if id < 0 || id >= len(l.Items) {
		return errors.New("bad todo-item ID")
	}
	l.Items[id].IsDone = !l.Items[id].IsDone
	return nil
}
