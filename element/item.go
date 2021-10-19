package element

import (
	"errors"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/fatih/color"
)

//An Item holds all the information for a single task, and a collection events tied to that task
type Item struct {
	ID     int
	Note   string
	Events []Event
	Group  string
	Due    time.Time
	IsDue  bool
}

//CheckIndex ensures that the given index is within the range of a collection (Array) of Items
func CheckIndex(items []Item, index int) bool {
	if index < 0 || index > len(items)-1 {
		color.Red("Item %d not found\n", index)
		return false
	} else {
		return true
	}
}

//Prints a collection of Items, and if given a non-empty group, will filter on that group
func PrintItems(items []Item, group string) {
	for i, item := range items {
		if item.Note != "" && (group == "" || item.Group == group) {
			fmt.Printf("%d:\t", i)
			item.printItem()
		}
	}
}

//Output a collection of Items in comma seperated value format
func PrintAllToCSV(items []Item) {
	for _, item := range items {
		item.PrintCSV()
	}
}

//Find an item in a collection of Items by ID
func FindByID(items []Item, ID int) (int, error) {
	for i, item := range items {
		if item.ID == ID {
			return i, nil
		}
	}
	return -1, errors.New("ID not found")

}

//Sets an Item given by index to being due and adds the due date to it based a string.
func SetDue(items []Item, index int, when string) {

	item := &items[index]

	if when == "tomorrow" {
		year, month, day := time.Now().Date()
		item.Due = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Local().Location()).Add(time.Hour * 24)
		item.IsDue = true
		fmt.Printf("%s due date is %02d/%02d/%d\n", items[index].Note, items[index].Due.Day(), items[index].Due.Month(), items[index].Due.Year())
	} else if when == "today" {
		year, month, day := time.Now().Date()
		item.Due = time.Date(year, month, day, 0, 0, 0, 0, time.Now().Local().Location())
		item.IsDue = true
		fmt.Printf("%s due date is %02d/%02d/%d\n", items[index].Note, items[index].Due.Day(), items[index].Due.Month(), items[index].Due.Year())
	} else {
		myDate, err := dateparse.ParseAny(when)
		if err != nil {
			fmt.Printf("%s - %s isn't invalid due date", err, when)
			return
		}
		item.Due = myDate
		item.IsDue = true
		fmt.Printf("%s due date is %02d/%02d/%d\n", items[index].Note, items[index].Due.Day(), items[index].Due.Month(), items[index].Due.Year())
	}
}

//Set the group an Item, given by index, to a named group. Use blank to remove from a group
func SetGroup(items []Item, index int, group string) {
	item := &items[index]
	item.Group = group
}

//Add a named Event to an Item, given by index and datestamp that event
func AddEvent(items []Item, index int, eventType string) {

	event := new(Event)
	event.Stamp = time.Now()
	event.EventType = eventType

	item := &items[index]
	item.Events = append(item.Events, *event)
	color.Green("%d: %s is %s.\n", item.ID, item.Note, event.EventType)

}

//Outputs the details of a given
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

//Insert an item into a given index position
func Insert(a []Item, index int, value Item) []Item {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

//Remove an item at a given index in a collection
func Remove(items []Item, index int) []Item {
	return append(items[:index], items[index+1:]...)
}

//SetPrioirity inserts an exsisting indexed item into a new position. Setting a priority greater the length of the collection pushes it to last place.
func SetPriority(items []Item, index int, priority int) {
	if priority < 0 {
		priority = 0
	}
	if priority > len(items) {
		priority = len(items)
	}
	itemCopy := items[index]
	items = Remove(items, index)
	if priority > len(items) {
		items = append(items, itemCopy)
	} else {
		items = Insert(items, priority, itemCopy)
	}

}

//Finds the Time of the earliest event in a collection
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

//Output any overdue items
func Nag(items []Item) {
	for i, item := range items {
		if item.IsDue && item.Due.Before(time.Now()) {
			color.Red("\n%d: %s is overdue.", i, item.Note)
		}
	}
}

//Calculate the average Lead time - started to done - of a collection of items
func AverageCycleTime(items []Item) int {

	total := 0
	count := 0
	for _, item := range items {
		var start, done time.Time
		var hasStart, hasEnd bool = false, false
		for _, event := range item.Events {
			if event.EventType == "Started" {
				start = event.Stamp
				hasStart = true
			}
			if event.EventType == "Done" {
				done = event.Stamp
				hasEnd = true
			}
		}
		if hasStart && hasEnd {
			total += int(done.Sub(start).Hours())
			count++
		}

	}

	return total / count
}

//Calculate the average Lead time - created to done - of a collection of items
func AverageLeadTime(items []Item) int {

	total := 0
	count := 0
	for _, item := range items {
		var start, done time.Time
		var hasStart, hasEnd bool = false, false
		for _, event := range item.Events {
			if event.EventType == "Created" {
				start = event.Stamp
				hasStart = true
			}
			if event.EventType == "Done" {
				done = event.Stamp
				hasEnd = true
			}
		}
		if hasStart && hasEnd {
			total += int(done.Sub(start).Hours())
			count++
		}

	}

	return total / count
}

//Returns true if the item has an event with the given event name
func (item Item) HasEvent(EventName string) bool {
	for _, event := range item.Events {
		if event.EventType == EventName {
			return true
		}
	}
	return false
}

//Output a colour coded version of the given item
func (item Item) printItem() {
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

//Output a CSV version of the item
func (item Item) PrintCSV() {

	fmt.Printf("%d, %s, %s ", item.ID, item.Note, item.Group)
	for _, event := range item.Events {

		fmt.Printf(", %s", event.EventType)

		t := event.Stamp
		fmt.Printf(", %02d/%02d/%d", t.Day(), t.Month(), t.Year())

	}
	fmt.Println()

}
