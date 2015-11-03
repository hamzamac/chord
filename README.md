# CHORD
To start the application run the node.go file
Each instance of a node.go represents an individual node on the ring

#NODE FLAGS
	-addrass		a string specifying the address (IP:PORT) of this node
	-join			a boolean specifying if the node is joining (true) or starting a new ring (false)
	-bnode			a string specifying the address (IP:PORT) on a boost node used to join the ring (only used in joining a an exesting ring)
	
#NODE OPERATION
	1 INSERT			allows insert of a text to the ring for storage
	2 FIND				requires a hash to find its corresponding valur rome the ring
	3 FINGER TABLE		prints currents node finger table
	4 LEAVE				allaows the node to safely leave the ring
	5 UPDATE			updates the fingertable of the current node
	6 INVENTORY			display the in-memory storage contents on the current node
	7 DELETE			delet a given value fron the ring using the hash

