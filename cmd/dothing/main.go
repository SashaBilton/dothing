package main

import (
	"dothing/element"
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
)

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

	dothing := new(element.DoThing)
	if command != "NEW" {
		dothing.Load()
	} else {
		dothing = createNewDothing()
		save = true
	}
	fmt.Println()
	switch command {
	case "add":
		dothing.AddItem(body, mod)
		save = true
	case "list":
		element.PrintItems(dothing.Items, body)
	case "listall":
		element.PrintItems(dothing.Items, "")
		color.Cyan("Done items\n")
		element.PrintItems(dothing.Done, "")
	case "done":
		index, _ := strconv.Atoi(body)
		if element.CheckIndex(dothing.Items, index) {
			dothing.ItemDone(index)
			element.PrintItems(dothing.Items, "")
			save = true
		}

	case "help":
		fmt.Println("Commands are add <note> <group>, list <group>, listall, done <ID>, raw, csv, due <ID> <when>, group <ID> <group>, event <ID> <event>")
	case "raw":
		fmt.Println(dothing)
	case "csv":
		element.PrintAllToCSV(dothing.Items)
		element.PrintAllToCSV(dothing.Done)

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
		dothing.PrintStats()
	case "undone":
		index, _ := strconv.Atoi(body)
		if element.CheckIndex(dothing.Done, index) {
			dothing.Undone(index)
			element.PrintItems(dothing.Items, "")
			save = true
		}
	case "nag":
		dothing.Nag = true
		save = true
	case "unnag":
		dothing.Nag = false
		save = true
	}

	if dothing.Nag {
		element.Nag(dothing.Items)

	}

	if save {
		dothing.Save()
		fmt.Println("\ndothing updated.")
	}
	fmt.Println()
}

//Creates a new empty DoThing object
func createNewDothing() *element.DoThing {
	dothing := new(element.DoThing)
	dothing.LastID = 0
	dothing.Items = []element.Item{}
	return dothing

}
