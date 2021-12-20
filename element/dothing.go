package element

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

//A DoThing houses a whole collection of items, both done and todo as well as storing the last ID used
type DoThing struct {
	Items  []Item
	Done   []Item
	LastID int
	Nag    bool
}

//CLI entrance

//Saves a dothing collection as a serial gob file, and also creates a historical entry in the hist directory
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

//Loads the local dothing.gob serial collection
func (dothing *DoThing) Load() {
	file, err := os.Open("dothing.gob")
	if err != nil {
		color.Red("Failed to load donothing.gob database with - %s\nIf you don't have an existing dothing list, create one with the command NEW. ", err)

	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	_ = decoder.Decode(dothing)
}

//ItemDone move the Item at the given index from the active Items list to the Done list
func (dothing *DoThing) ItemDone(index int) {

	done := new(Event)
	done.Stamp = time.Now()
	done.EventType = "Done"

	item := &dothing.Items[index]
	item.Events = append(item.Events, *done)

	itemCopy := dothing.Items[index]
	dothing.Done = append(dothing.Done, itemCopy)

	dothing.Items = Remove(dothing.Items, index)

	color.Green("%s is done.\n", itemCopy.Note)

}

//Undone moves a Item back from Done to the active Items collection
func (dothing *DoThing) Undone(index int) {
	done := new(Event)
	done.Stamp = time.Now()
	done.EventType = "Undone"
	item := &dothing.Done[index]
	item.Events = append(item.Events, *done)

	dothing.Items = append(dothing.Items, *item)
	dothing.Done = Remove(dothing.Done, index)

	color.Green("%s is undone.\n", item.Note)

}

//Adds
func (dothing *DoThing) AddItem(note string, group string) {

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
	fmt.Printf("%d: %s added to %s\n", len(dothing.Items)-1, addItem.Note, addItem.Group)

}

//Outputs statistics on active and Done items.
func (dothing *DoThing) PrintStats() {

	items := len(dothing.Items)
	done := len(dothing.Done)
	total := items + done

	incomplete := 100.0 / float32(total) * float32(done)

	fmt.Printf("To do: %d Done: %d Total: %d Complete: %.00f%% \n", items, done, total, incomplete)

	earliestItem := Earliest(dothing.Items)
	earliestDone := Earliest(dothing.Done)
	now := time.Now()

	//fmt.Printf("Earliest Item %02d/%02d/%d\n", earliestItem.Day(), earliestItem.Month(), earliestItem.Year())

	itemDays := now.Sub(earliestItem).Hours() / 24
	doneDays := now.Sub(earliestDone).Hours() / 24

	fmt.Printf("Oldest item: %.00f days.  Oldest Done item: %.00f days\n", itemDays, doneDays)
	donePerDay := doneDays / float64(done)
	daysToDoItems := float64(items) / donePerDay
	color.HiGreen("Done per day: %.2f Days to complete all active items: %.0f", donePerDay, daysToDoItems)
	cyclehrs := AverageCycleTime(dothing.Done)
	color.HiGreen("Average cycle time of started and finished items: %d hrs or %.1f days", cyclehrs, float32(cyclehrs)/24)
	leadhrs := AverageLeadTime(dothing.Done)
	color.HiGreen("Average lead time of created and finished items: %d hrs or %.1f days", leadhrs, float32(leadhrs)/24)

}
