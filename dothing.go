package main

import (
	"dothing/element"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
)

type DoThing struct {
	Items  []element.Item
	Done   []element.Item
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
		element.PrintItems(dothing.Items, body)
	case "listall":
		element.PrintAllItems(dothing.Items)
		color.Cyan("Done items\n")
		element.PrintAllItems(dothing.Done)
	case "done":
		index, _ := strconv.Atoi(body)
		if element.CheckIndex(dothing.Items, index) {
			Done(dothing, index)
			element.PrintItems(dothing.Items, "")
			save = true
		}

	case "help":
		fmt.Println("Commands are add <note> <group>, list <group>, listall, done <ID>, raw, csv, due <ID> <when>, group <ID> <group>, event <ID> <event>")
	case "raw":
		fmt.Println(dothing)
	case "csv":
		element.PrintAllToCSV(dothing.Items)
	case "due":
		index, _ := strconv.Atoi(body)
		if element.CheckIndex(dothing.Items, index) {
			element.SetDue(dothing.Items, index, mod)
			element.Detail(dothing.Items, index)
			save = true
		}
	case "group":
		index, _ := strconv.Atoi(body)
		if element.CheckIndex(dothing.Items, index) {
			element.SetGroup(dothing.Items, index, mod)
			element.Detail(dothing.Items, index)
			save = true
		}
	case "event":
		index, _ := strconv.Atoi(body)
		if element.CheckIndex(dothing.Items, index) {

			element.AddEvent(dothing.Items, index, mod)
			element.Detail(dothing.Items, index)
			save = true
		}
	case "detail":
		index, _ := strconv.Atoi(body)
		if element.CheckIndex(dothing.Items, index) {

			element.Detail(dothing.Items, index)
		}
	case "priority":
		index, _ := strconv.Atoi(body)
		if element.CheckIndex(dothing.Items, index) {

			priority, _ := strconv.Atoi(mod)
			element.SetPriority(dothing.Items, index, priority)
			element.PrintItems(dothing.Items, "")

			save = true
		}
	case "stats":
		PrintStats(dothing)
	case "undone":
		index, _ := strconv.Atoi(body)
		if element.CheckIndex(dothing.Done, index) {
			Undone(dothing, index)
			element.PrintItems(dothing.Items, "")

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

func Done(dothing *DoThing, index int) {

	done := new(element.Event)
	done.Stamp = time.Now()
	done.EventType = "Done"

	item := &dothing.Items[index]
	item.Events = append(item.Events, *done)

	itemCopy := dothing.Items[index]
	dothing.Done = append(dothing.Done, itemCopy)

	dothing.Items = element.Remove(dothing.Items, index)

	color.Green("%s is done.\n", itemCopy.Note)

}

func Undone(dothing *DoThing, index int) {
	done := new(element.Event)
	done.Stamp = time.Now()
	done.EventType = "Undone"
	item := &dothing.Done[index]
	item.Events = append(item.Events, *done)

	dothing.Items = append(dothing.Items, *item)
	dothing.Done = element.Remove(dothing.Done, index)

	color.Green("%s is undone.\n", item.Note)

}

func AddItem(dothing *DoThing, note string, group string) {

	item := new(element.Item)
	item.Note = note
	item.Group = group
	dothing.LastID++
	item.ID = dothing.LastID

	created := new(element.Event)
	created.Stamp = time.Now()
	created.EventType = "Created"
	item.Events = []element.Event{*created}
	dothing.Items = append(dothing.Items, *item)

	addItem := dothing.Items[len(dothing.Items)-1]
	fmt.Printf("%d: %s added to %s\n", addItem.ID, addItem.Note, addItem.Group)

}

func PrintStats(dothing *DoThing) {

	items := len(dothing.Items)
	done := len(dothing.Done)
	total := items + done

	incomplete := 100.0 / float32(total) * float32(done)

	fmt.Printf("To do: %d Done: %d Total : %d Complete: %.00f%% \n", items, done, total, incomplete)

	earliestItem := element.Earliest(dothing.Items)
	earliestDone := element.Earliest(dothing.Done)
	now := time.Now()

	//fmt.Printf("Earliest Item %02d/%02d/%d\n", earliestItem.Day(), earliestItem.Month(), earliestItem.Year())

	itemDays := now.Sub(earliestItem).Hours() / 24
	doneDays := now.Sub(earliestDone).Hours() / 24

	fmt.Printf("Days since first Item: %.00f Days since first Done %.00f\n", itemDays, doneDays)
	donePerDay := doneDays / float64(done)
	daysToDoItems := float64(items) / donePerDay
	color.HiGreen("Done per day: %.2f Days to complete items: %.0f", donePerDay, daysToDoItems)
}

func PurgeDone(dothing *DoThing) {
	for i, item := range dothing.Items {
		if item.HasEvent("Done") {
			Done(dothing, i)
		}
	}

}

func createNewDothing() *DoThing {
	dothing := new(DoThing)
	dothing.LastID = 0
	dothing.Items = []element.Item{}
	return dothing

}
