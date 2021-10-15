# DoThing
A simple single user CLI based task tracking tool.  

Run uncompiled as go run dothing.go {command}

Commands are
NEW - creates a brand new dothing task list
add "your task details" {optional group} - adds a task 
list {optional group} - lists all currently active tasks
Listall - lists all active and done tasks
done {index} - Moves a task into the done collection
undone {index} - Returns a done task to the active list
priority {index} {priority} - Sets a task as priority x. 0 is top priority, any priority greater than the size of the list goes to the bottom.
due {when} - Gives a task a due date (currently only supports 'tomorrow')
event {index} {event name} - Adds an event to a task
detail {index} - Displays full details of an active task
csv - Displays all tasks, active and done, in a comma seperated format
raw - Displays a raw JSON-like output
stats - Displays stats on how many items are done or active, etc,
