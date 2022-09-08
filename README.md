# CHORD
## About the Project
The main purpose of this assignment is to build a distributed system by implementing a Distributed Hash Table (DHT) algorithm called Chord. Chord is a distributed lookup protocol designed to locate a node in a network of nodes where a particular data object is stored. 



## Project Report
https://1drv.ms/w/s!AmDAJNCe6kNuo2-C8BEuq5-zQLh6

## Usage

To start the application run the `node.go` file

Each instance of a node.go represents an individual node on the ring

#### NODE FLAGS

	-address		a string specifying the address (`IP:PORT`) of this node
	-join			a boolean specifying if the node is joining (`true`) or starting a new ring (false)
	-bnode			a string specifying the address (`IP:PORT`) on a boost node used to join the ring (only used in joining a an exesting ring)
	
#### NODE OPERATION

	INSERT			allows insert of a text to the ring for storage
	FIND				requires a hash to find its corresponding valur rome the ring
	FINGER TABLE		prints currents node finger table
	LEAVE				allaows the node to safely leave the ring
	UPDATE			updates the fingertable of the current node
	INVENTORY			display the in-memory storage contents on the current node
	DELETE			delet a given value fron the ring using the hash

## Docker
You can also find a docker image version of the code here https://hub.docker.com/r/jeanpok8/skychord/

### Instructions to run the project with Docker containers
1. Run `docker pull jeanpok8/skychord`(for 160 bits) 
2. Start the first container and run a node inside with following commands: 
	- open the terminal 
 	- `docker run -i -t --name Container1 jeanpok8/skychord /bin/bash`
 	- `cd bin`
 	- `./chord -address=Container1:1201 -join=false`

3. start another node on container2

	- open the terminal
 	- `docker run -ti-p 127.0.0.1::1202/udp --name Container2 --link Container1 jeanpok8/skychord /bin/bash`
 	- `cd bin`
 	- `./chord -address=Container2:1202 -join=true -bnode=Container1:1201`

N.B:you may run additional nodes to the ring by repeating step 3 but by changing:

	- Port number
	- Name of container
