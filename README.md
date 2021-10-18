# DoThing
A simple single user CLI based task tracking tool.  


go get github.com/SashaBilton/dothing

Run uncompiled as go run dothing.go {command}

Commands are
NEW - creates a brand new dothing task list - do this first<br>
add "your task details" {optional group} - adds a task <br>
list {optional group} - lists all currently active tasks<br>
Listall - lists all active and done tasks<br>
done {index} - Moves a task into the done collection<br>
undone {index} - Returns a done task to the active list<br>
priority {index} {priority} - Sets a task as priority x. 0 is top priority, any priority greater than the size of the list goes to the bottom.<br>
due {when} - Gives a task a due date currently supports today, tomorrom, & mm-dd-yyyy (US style)<br>
event {index} {event name} - Adds an event to a task<br>
detail {index} - Displays full details of an active task<br>
csv - Displays all tasks, active and done, in a comma seperated format<br>
raw - Displays a raw JSON-like output<br>
stats - Displays stats on how many items are done or active, etc,.<br>
nag - turns on the display of overdue items on every command.<br>
unnage - turns off the display of overdue items on every command.<br> 

