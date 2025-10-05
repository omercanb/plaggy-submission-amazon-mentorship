/*
NOTES
I think the comments below explain briefly what is going on, however go really isn't what I'm used to
so there are probably a lot of shit that is suboptiomal or just outright fucking wrong so please double check
whatever the fuck I wrote.

READ THIS FOR STRUCTURE INFO
/models is where we define types/classes/structs or whatever the fuck you want to call them
I think we need to (I'm not sure about this part tbh) define structs to convert our data to JSON format
to then send to wherever it needs to go.

/routeHandles is where our handler functions go. How to use them is below.
Check /routeHandles/detections.go for a simple example of how we send JSON.
>
*/

package main

import "github.com/plagai/plagai-backend/server"

// standard lib imports

func main() {
	server.Start()
}
