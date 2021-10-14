package element

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Item struct {
	ID     int
	Note   string
	Events []Event
	Group  string
	Due    time.Time
	IsDue  bool
}

type Event struct {
	EventType string
	Stamp     time.Time
}

func CheckIndex(items []Item, index int) bool {
	if index < 0 || index > len(items)-1 {
		color.Red("Item %d not found\n", index)
		return false
	} else {
		return true
	}
}

func PrintItems(items []Item, group string) {
	for i, item := range items {
		if item.Note != "" && (group == "" || item.Group == group) {
			fmt.Printf("%d:\t", i)
			printItem(item)
		}
	}
}

func PrintAllItems(items []Item) {
	for i, item := range items {
		fmt.Printf("%d:\t", i)
		printItem(item)
	}
}

func PrintAllToCSV(items []Item) {
	for _, item := range items {
		PrintCSV(item)
	}
}

func printItem(item Item) {
	green := color.New(color.FgGreen)
	white := color.New(color.FgWhite)
	yellow := color.New(color.FgYellow)
	red := color.New(color.FgRed)

	p := white
	if item.IsDue {
		//
		t := time.Now()
		warning := item.Due.Add(-time.Hour * 48)
		//fmt.Printf("t %v d %v w %v c %v", t, item.Due, warning, critical)
		if t.After(warning) && t.Before(item.Due) {
			p = yellow
		} else if t.After(item.Due) {
			p = red
		}

	}

	if Is(item.Events, "Done") {
		p = green
	}
	if item.Group != "" {
		p.Printf("%s [%s] - ", item.Note, item.Group)
	} else {
		p.Printf("%s - ", item.Note)
	}
	if item.IsDue {
		p.Printf(" due %02d/%02d/%04d ", item.Due.Day(), item.Due.Month(), item.Due.Year())
	}

	for i, event := range item.Events {
		if i > 0 {
			p.Printf(", %s", event.EventType)
		} else {
			p.Printf("%s", event.EventType)
		}

		p.Printf(" %02d/%02d/%d", event.Stamp.Day(), event.Stamp.Month(), event.Stamp.Year())

	}
	fmt.Println()

}

func PrintCSV(item Item) {

	fmt.Printf("%d, %s, %s ", item.ID, item.Note, item.Group)
	for _, event := range item.Events {

		fmt.Printf(", %s", event.EventType)

		t := event.Stamp
		fmt.Printf(", %02d/%02d/%d", t.Day(), t.Month(), t.Year())

	}
	fmt.Println()

}

func FindByID(items []Item, ID int) (int, error) {
	for i, item := range items {
		if item.ID == ID {
			return i, nil
		}
	}
	return -1, errors.New("ID not found")

}

func SetDue(items []Item, index int, when string) {

	item := &items[index]

	if when == "tomorrow" {
		year, month, day := time.Now().Date()
		item.Due = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Local().Location()).Add(time.Hour * 24)
		item.IsDue = true
		fmt.Printf("%s due date is %02d/%02d/%d\n", items[index].Note, items[index].Due.Day(), items[index].Due.Month(), items[index].Due.Year())

	}

}

func SetGroup(items []Item, index int, group string) {
	item := &items[index]
	item.Group = group
}

func AddEvent(items []Item, index int, eventType string) {

	event := new(Event)
	event.Stamp = time.Now()
	event.EventType = eventType

	item := &items[index]
	item.Events = append(item.Events, *event)
	color.Green("%d: %s is %s.\n", item.ID, item.Note, event.EventType)

}

func Detail(items []Item, index int) {

	item := items[index]
	fmt.Printf("Item\t%s\nID\t%d\nGroup\t%s\n", item.Note, item.ID, item.Group)
	if item.IsDue {
		fmt.Printf("Due\t%02d/%02d/%d\n", item.Due.Day(), item.Due.Month(), item.Due.Year())
	}
	fmt.Printf("Pos\t%d\n", index)
	fmt.Printf("Events")
	for _, event := range item.Events {

		fmt.Printf("\t%s", event.EventType)
		fmt.Printf(" %02d/%02d/%d\n", event.Stamp.Day(), event.Stamp.Month(), event.Stamp.Year())

	}

}

func Insert(a []Item, index int, value Item) []Item {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func Remove(slice []Item, s int) []Item {
	return append(slice[:s], slice[s+1:]...)
}

func SetPriority(items []Item, index int, priority int) {

	if priority > len(items) {
		priority = len(items)
	}
	itemCopy := items[index]
	item := &items[index]
	fmt.Printf("Remove %d %s\n", item.ID, item.Note)
	items = Remove(items, index)
	fmt.Printf("Insert %d %s\n", item.ID, item.Note)
	if priority > len(items) {
		items = append(items, itemCopy)
	} else {
		items = Insert(items, priority, itemCopy)
	}

}
func Earliest(items []Item) time.Time {
	earliest := time.Now()
	for _, item := range items {
		for _, event := range item.Events {
			if event.Stamp.Before(earliest) {
				earliest = event.Stamp
			}
		}
	}
	return earliest
}

func HasEvent(item Item, EventName string) bool {
	for _, event := range item.Events {
		if event.EventType == EventName {
			return true
		}
	}
	return false
}

func Is(events []Event, ofType string) bool {
	for _, event := range events {
		if event.EventType == ofType {
			return true
		}
	}
	return false
}
