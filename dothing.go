package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
)

type Event struct {
	EventType string
	Stamp     time.Time
}
type Item struct {
	ID     int
	Note   string
	Events []Event
	Group  string
	Due    time.Time
	IsDue  bool
}

type DoThing struct {
	Items  []Item
	Done   []Item
	LastID int
}

func main() {

	//createSaveTestData()

	command := "help"
	body := ""
	mod := ""

	save := false

	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	if len(os.Args) > 2 {
		body = os.Args[2]
	}
	if len(os.Args) > 3 {
		mod = os.Args[3]
	}

	dothing := new(DoThing)
	if command != "NEW" {
		dothing.Load()
	} else {
		dothing = createNewDothing()
		save = true
	}
	fmt.Println()
	switch command {
	case "add":
		AddItem(dothing, body, mod)
		save = true
	case "list":
		PrintItems(dothing.Items, body)
	case "listall":
		PrintAllItems(dothing.Items)
		color.Cyan("Done items\n")
		PrintAllItems(dothing.Done)
	case "done":
		index, _ := strconv.Atoi(body)
		if CheckIndex(dothing.Items, index) {
			Done(dothing, index)
			PrintItems(dothing.Items, "")
			save = true
		}

	case "help":
		fmt.Println("Commands are add <note> <group>, list <group>, listall, done <ID>, raw, csv, due <ID> <when>, group <ID> <group>, event <ID> <event>")
	case "raw":
		fmt.Println(dothing)
	case "csv":
		PrintAllToCSV(dothing.Items)
	case "due":
		index, _ := strconv.Atoi(body)
		if CheckIndex(dothing.Items, index) {
			SetDue(dothing.Items, index, mod)
			Detail(dothing.Items, index)
			save = true
		}
	case "group":
		index, _ := strconv.Atoi(body)
		if CheckIndex(dothing.Items, index) {
			SetGroup(dothing.Items, index, mod)
			Detail(dothing.Items, index)
			save = true
		}
	case "event":
		index, _ := strconv.Atoi(body)
		if CheckIndex(dothing.Items, index) {

			AddEvent(dothing.Items, index, mod)
			Detail(dothing.Items, index)
			save = true
		}
	case "detail":
		index, _ := strconv.Atoi(body)
		if CheckIndex(dothing.Items, index) {

			Detail(dothing.Items, index)
		}
	case "priority":
		index, _ := strconv.Atoi(body)
		if CheckIndex(dothing.Items, index) {

			priority, _ := strconv.Atoi(mod)
			SetPriority(dothing.Items, index, priority)
			PrintItems(dothing.Items, "")

			save = true
		}
	case "stats":
		PrintStats(dothing)
	case "undone":
		index, _ := strconv.Atoi(body)
		if CheckIndex(dothing.Done, index) {
			Undone(dothing, index)
			PrintItems(dothing.Items, "")

			save = true
		}

	case "PURGE_DONE":
		PurgeDone(dothing)
		save = true

	}

	if save {
		dothing.Save()
		fmt.Println("dothing updated.")
	}
	fmt.Println()
}

func (dothing *DoThing) Save() {

	file, err := os.Create("dothing.gob")

	filetime := time.Now()
	timename := fmt.Sprintf("hist/hist_dothing_%d%02d%02d%02d%02d%02d.gob", filetime.Year(), filetime.Month(), filetime.Day(), filetime.Hour(), filetime.Minute(), filetime.Second())

	hist_file, hist_err := os.Create(timename)

	defer file.Close()
	defer hist_file.Close()

	encoder := gob.NewEncoder(file)
	hist_encoder := gob.NewEncoder(hist_file)
	encoder.Encode(dothing)
	encoder.Encode(hist_encoder)
	if err != nil {
		fmt.Printf("%s", err)
	}
	if hist_err != nil {
		fmt.Printf("%s", hist_err)
	}

}

func (dothing *DoThing) Load() {
	file, err := os.Open("dothing.gob")
	if err != nil {
		color.Red("Failed to load donothing.gob database with - %s\nIf you don't have an existing dothing list, create one with the command NEW. ", err)

	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	_ = decoder.Decode(dothing)
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

func Done(dothing *DoThing, index int) {

	done := new(Event)
	done.Stamp = time.Now()
	done.EventType = "Done"

	item := &dothing.Items[index]
	item.Events = append(item.Events, *done)

	itemCopy := dothing.Items[index]
	dothing.Done = append(dothing.Done, itemCopy)

	dothing.Items = remove(dothing.Items, index)

	color.Green("%s is done.\n", itemCopy.Note)

}

func Undone(dothing *DoThing, index int) {
	done := new(Event)
	done.Stamp = time.Now()
	done.EventType = "Undone"
	item := &dothing.Done[index]
	item.Events = append(item.Events, *done)

	dothing.Items = append(dothing.Items, *item)
	dothing.Done = remove(dothing.Done, index)

	color.Green("%s is undone.\n", item.Note)

}

func AddItem(dothing *DoThing, note string, group string) {

	item := new(Item)
	item.Note = note
	item.Group = group
	dothing.LastID++
	item.ID = dothing.LastID

	created := new(Event)
	created.Stamp = time.Now()
	created.EventType = "Created"
	item.Events = []Event{*created}
	dothing.Items = append(dothing.Items, *item)

	addItem := dothing.Items[len(dothing.Items)-1]
	fmt.Printf("%d: %s added to %s\n", addItem.ID, addItem.Note, addItem.Group)

}

func is(events []Event, ofType string) bool {
	for _, event := range events {
		if event.EventType == ofType {
			return true
		}
	}
	return false
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

	if is(item.Events, "Done") {
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

func insert(a []Item, index int, value Item) []Item {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func remove(slice []Item, s int) []Item {
	return append(slice[:s], slice[s+1:]...)
}

func SetPriority(items []Item, index int, priority int) {

	if priority > len(items) {
		priority = len(items)
	}
	itemCopy := items[index]
	item := &items[index]
	fmt.Printf("Remove %d %s\n", item.ID, item.Note)
	items = remove(items, index)
	fmt.Printf("Insert %d %s\n", item.ID, item.Note)
	if priority > len(items) {
		items = append(items, itemCopy)
	} else {
		items = insert(items, priority, itemCopy)
	}

}

func HasEvent(item Item, EventName string) bool {
	for _, event := range item.Events {
		if event.EventType == EventName {
			return true
		}
	}
	return false
}

func PrintStats(dothing *DoThing) {

	items := len(dothing.Items)
	done := len(dothing.Done)
	total := items + done

	incomplete := 100.0 / float32(total) * float32(done)

	fmt.Printf("To do: %d Done: %d Total : %d Complete: %.00f%% \n", items, done, total, incomplete)

	earliestItem := Earliest(dothing.Items)
	earliestDone := Earliest(dothing.Done)
	now := time.Now()

	//fmt.Printf("Earliest Item %02d/%02d/%d\n", earliestItem.Day(), earliestItem.Month(), earliestItem.Year())

	itemDays := now.Sub(earliestItem).Hours() / 24
	doneDays := now.Sub(earliestDone).Hours() / 24

	fmt.Printf("Days since first Item: %.00f Days since first Done %.00f\n", itemDays, doneDays)
	donePerDay := doneDays / float64(done)
	daysToDoItems := float64(items) / donePerDay
	color.HiGreen("Done per day: %.2f Days to complete items: %.0f", donePerDay, daysToDoItems)
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

func PurgeDone(dothing *DoThing) {
	for i, item := range dothing.Items {
		if HasEvent(item, "Done") {
			Done(dothing, i)
		}
	}

}

func createNewDothing() *DoThing {
	dothing := new(DoThing)
	dothing.LastID = 0
	dothing.Items = []Item{}
	return dothing

}
