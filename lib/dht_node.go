package dht
import (
	"fmt"
	"strconv"
	"math/big"
	"encoding/json"
	"crypto/sha1"
	"os"
	"time"
)

const BITS int = 160 //set the number of bits for the ring size


type DHTNode struct {
	nodeId      big.Int					// the id of the node
	address 	string					// the address of this node
	target 		string					// the address of the boost node
	Pre			string					// the id of the predecessor node
	PreAdr		string					// the address of the predecessor node
	offset		int					// the offset used to avoide a nore from responding from utdated command
	FingerTable [BITS][BITS]string		// the figer table
	storage  map[string]string			// the inmeory storage for storing data
}


// GLOBAL INITIALIZATION: 
var me DHTNode						//create this clients node
//set functions
func SetAddress(adrs string){ 		// assign address
	me.address=adrs
}
func SetTarget(tgt string){			//set the address of the boost node for finger table exchager
	me.target=tgt
}
func SetOffset(ofst int){			// update the offters to the current offset of the ring
	me.offset=ofst
}
func IncrementOffset(){				//increment ofsset before initiating a grolal action
	me.offset++
}
	
//CODES FOR IMPLEMENTATION
/*
	The Join function
	The first fuction to be executed
	Alaows it gives the hash and id and request the finger table for the boost node
*/
func Join(test bool){ 				//
	
	//generate id
	me.nodeId=Index(generateNodeId(), BITS)
	nodeIdInt:= me.nodeId
	me.storage=map[string]string{}
		
	//prepare finger table
	if(test){ //firstnode in the ring
		me.offset=0
		for i,_:= range me.FingerTable {
			//compute the ith node
			two := big.NewInt(2)
			dist := big.Int{}
			dist.Exp(two, big.NewInt(int64(i)), nil)
			var ithNode big.Int
			ithNode.Add(&dist,&nodeIdInt)
			x:=ringSize(BITS)
			ithNode.Mod(&ithNode,&x) 
			//fill the fingr table row
			me.FingerTable[i][0]= ithNode.String()
			me.FingerTable[i][1]=me.nodeId.String()
			me.FingerTable[i][2]=me.address //this node address	
		}
	}else{	//not first node in the ring
		//initialize offset
		GetOffsetClient(me.target)		
		updateFingerTable(nodeIdInt,me.target) //target server required
		go infornMySuccessor()
		go AutoUpdate()						// initialize auto update
	}
		
}
/*
	Used to update the fingertable of a node by assing for node responsible for each finger in the table
*/
func updateFingerTable(nodeIdInt big.Int, questioneeAddress string){// node's assigned id, questionee node address
	
	for i,_:= range me.FingerTable {
		//compute the ith node
		two := big.NewInt(2)
		dist := big.Int{}
		dist.Exp(two, big.NewInt(int64(i)), nil)
		var ithNode big.Int
		ithNode.Add(&dist,&nodeIdInt)
		x:=ringSize(BITS)
		ithNode.Mod(&ithNode,&x)
		//fill the fingr table row
		me.FingerTable[i][0]=ithNode.String()
		suc,address:=askForAddress(ithNode.String(),questioneeAddress)
		me.FingerTable[i][1]=suc
		me.FingerTable[i][2]=address
		go infornMySuccessor()
	}		
}	

/*
	A recursive function that asks other nodes in ithe ring for the address of a specific node,
	if the asked node do not have information about it replays with the other node in the ring to ask
*/
func askForAddress(id string, address string) (string, string){ // successor, address
	//send the request and receive responce
	msg:=Msg{Request:"getAddress",Dst:address,Id:id}
	msgresponce:=poster(msg)		
	//test the response
	if(msgresponce.Inf=="found"){		//found
		return msgresponce.Id, msgresponce.Address
	}	//not found, ask another node that is not me
	//fmt.Println("asking ",msgresponce.Id)
	return askForAddress(msgresponce.Id, msgresponce.Address)	
}

/*
	Amx (stands for "I am X"), this function used by the node to introduce it self to other nodes in ine ring just after joining the ring.
*/
func Amx(id string, myaddress string, destAddress string, offset string){//id,sourceaddress,destAddress
	//send the request and receive responce
	msg:=Msg{Request:"Joined",Dst:destAddress,Src:myaddress,Id:id,Inf:offset}
	_=poster(msg)				
}


/*
 	Used to callculat distance betweeen two node, 
	used by node when making decision on finding the closest node.
*/
func nodeDistance(s int,d int) int{
	souece:= big.NewInt(int64(s))
	dest:= big.NewInt(int64(d))
	r, _ := strconv.ParseInt(distance(souece.Bytes(), dest.Bytes(), BITS).String(), 10, 0)
	return int(r)
}

/* 
	Used by a node after joining the ring to ask for the current offet valur in the ring. 
*/
func GetOffsetClient(destAddress string){
	//send the request and receive responce
	msg:=Msg{Request:"offset",Dst:destAddress}
	msgresponce:=poster(msg)		
	//test the response
	o, _ := strconv.ParseInt(msgresponce.Inf, 10, 0)	
	SetOffset(int(o))		
}

/*
	Used by a node to leav the ring
	first it tells all the node on itsf finger table that its leaving.
	then it exist the the ring	
*/
func leav(){
	
	//increment offset
	IncrementOffset()
	//tell all nodes other than me in my finger table to update
	var test = false
	idx:=0
	var tempft[BITS][BITS] string
	for i,_:=range me.FingerTable {
		if me.nodeId.String()!=me.FingerTable[i][1] {
			test= true
			//add none me finger to temp array
			tempft[idx]=me.FingerTable[i]
			idx++
		}
	}	
	
	if test {
		for x,_:=range tempft{
			//send a leave request
			msg:=Msg{Request:"ileft",Dst:tempft[x][2],Id:me.nodeId.String(),Suc:me.FingerTable[0][1],SucAdr:me.FingerTable[0][2],Pre:me.Pre,PreAdr:me.PreAdr,Inf:GetOffsetServer()}//id=id key=pre inf=suc	
			_=poster(msg)				
		}
	}else{
		//do nothing
	}
	//shutdown
	os.Exit(0)		
}

/*
	Used by a node to tell its successor that "i am your predecessor"
	called everytine the finger table is updated
*/
func infornMySuccessor(){
		msg:=Msg{Request:"predecessor",Dst:me.FingerTable[0][2],Id:me.nodeId.String(),Inf:me.address}
		_=poster(msg)		
}


// SERVER SIDE CODES
/*
	Used by the server part of a node to respond to askForAddress()
	It searched its finger table for the requested idif not forn it suggest other node in the ring to be asked
*/
func Lookup(id string,fingertable [BITS][BITS]string)(string,string,string){//successor,address,inf
	// Are you looking for me?
	if id == me.nodeId.String() {
		//fmt.Println("looking for me? ", id)
		return me.nodeId.String(), me.address, "found"
	}
	l:=len(fingertable)-1
	//find if id is beween the ringer references
	answer4s:=false
	answer4t:=between(str2byte(fingertable[0][0]),str2byte(fingertable[l][0]),str2byte(id))
	if fingertable[0][1]!=fingertable[l][1] {
		answer4s=between(str2byte(fingertable[0][1]),str2byte(fingertable[l][1]),str2byte(id))
	}
			
	//fmt.Print("looking in ",me.nodeId)
	if answer4s {
		//fmt.Print("found in s ",id)
		//find the closest successor and get her
		s, _ := strconv.Atoi(fingertable[0][1])
		n, _ := strconv.Atoi(id)
		var closest[BITS][BITS] string
		min:=nodeDistance(n,s)
		closest[0]=fingertable[0]
		for i,_ :=range fingertable{
			//conver successor to interger
			s, _ := strconv.Atoi(fingertable[i][1])
			n, _ := strconv.Atoi(id)

			if min > nodeDistance(n,s) {
				closest[0]=fingertable[i]
				min=nodeDistance(n,s)
			}	
		}
		
		//fmt.Println(" found", closest[0][1])
		return closest[0][1],closest[0][2], "found"	
		// end of find closest				
		
	}else if answer4t {
		//fmt.Print("found in t ",id)
		//find closest node and get successor
		t, _ := strconv.Atoi(fingertable[0][0])
		n, _ := strconv.Atoi(id)
		var closest[BITS][BITS] string
		min:=nodeDistance(n,t)
		closest[0]=fingertable[0]
		for i,_ :=range fingertable{
			//convert successor to interger
			t, _ := strconv.Atoi(fingertable[i][0])
			n, _ := strconv.Atoi(id)

			if min > nodeDistance(n,t) {
				closest[0]=fingertable[i]
				min=nodeDistance(n,t)
			}	
		}
		//fmt.Println(" found", closest[0][1])
		return closest[0][1],closest[0][2], "found"	
		// end of find closest		
	}else{
		//chek if the node is between the last (l) finger (t) and its successor (s)
		lenx:=len(fingertable)-1
		answer4l:=between(str2byte(fingertable[lenx][0]),str2byte(fingertable[lenx][1]),str2byte(id))
		if answer4l {
			return fingertable[lenx][1],fingertable[lenx][2], "found"
		}
		// check is the successor poits to other nodes than this
		//fmt.Println("all false")
		var test = false
		idx:=0
		var tempft[BITS][BITS] string
		for i,_:=range fingertable {
			if me.nodeId.String()!=fingertable[i][1] {
				test= true
				//add none me finger to temp array
				tempft[idx]=fingertable[i]
				idx++
			}
		}
		
		le:=len(tempft)-1
		//ask another node that is not me CORRECTION NEEDED
		if test {// ask the successor of the another node that is not me
			return id,tempft[le][2], "notfound"
		}else{// its only me here
			return me.nodeId.String(),me.address, "found"
		}			
	}
}

/*
	Used by the server part of the node to update its finger tabler on new node joining the ring or on other events
*/
func Update(id string, address string, fingertable *[BITS][BITS]string, offset string) string{ //id, address, fingertable
	//fmt.Println("Update")
	// confirm the ofset is grater
	ofst, _ := strconv.ParseInt(offset, 10, 0)
	if me.offset< int(ofst) {
		me.offset=int(ofst)
		for i,_:= range fingertable{
			// neither me of the sender -> notigy on the new node
			if fingertable[i][1]!=id && fingertable[i][1]!=me.nodeId.String() {
				//spoof the Amx request
				go Amx(id,address, fingertable[i][2],offset)
			}
			//convert contenders to intergers
			t, _ := strconv.Atoi(fingertable[i][0])
			s, _ := strconv.Atoi(fingertable[i][1])
			n, _ := strconv.Atoi(id)
			//
			if (nodeDistance(t,s)>=nodeDistance(t,n)) { //fing closest to the reference
				fingertable[i][1]=id
				fingertable[i][2]=address
			}
		}		
	}	
	return "ack"
}

/*
	return the id of the current instance of a node
*/
func GetId() string{
	return me.nodeId.String()
}

/*
	Userd by the server part of the noder to get the the previous offset 
	For comparizon with the current offset to avoiderepeted requests
*/
func GetOffsetServer() string{
	o:= strconv.FormatInt(int64(me.offset), 10)
	return o		
}

/*
	Used by the server part of the node to store the received key value pair to the node in memory storage 
	It also replicate the key to its successor where necessacry for consistancy
*/
func store(key string, value string,replicate bool){
	//copie the value to inmemory storage
	me.storage[key]=value
	if replicate {
		//update the successor
		msg:=Msg{Request:"Replicate",Dst:me.FingerTable[0][2],Key:key,Inf:value}
		_=poster(msg)			
	}
}

/*
	used by the server part of a node to delete the requested key value pair from the inmemory storage
	it also delete the backup key on the successor where necessary
*/
func deleteKey(hash string, dereplicate bool){
	//delete value from  inmemory storage
	delete(me.storage, hash)
	if dereplicate {
		//update the successor
		msg:=Msg{Request:"dereplicate",Dst:me.FingerTable[0][2],Key:hash}
		_=poster(msg)			
	}			
}

/*
	Used by the server part of the node in response to a hash query
*/
func getValue(hash string) string{
	//retreave value from  inmemory storage
	return me.storage[hash]
}

/*
	Used by the server pard of the node in response to infornMySuccessor()
	set the value received from predecessor
*/
func setPredecessor(id,address string){
	me.Pre=id
	me.PreAdr=address
}

/*
	Used by the server part of the of the node i responsr to leav()
	checks if the nide is in the finger table and remove it
*/
func hasLeft(id, suc,sucAdr,pre,preAdr,offset string){
	// check the offset
	ofst, _ := strconv.ParseInt(offset, 10, 0)
	if me.offset< int(ofst) {
		me.offset=int(ofst)	
		
		for i,_:=range me.FingerTable{
			//neither me nor the sender
			if me.FingerTable[i][1]!=id && me.FingerTable[i][1]!=me.nodeId.String() {
				//spoof the leav() request
				msg:=Msg{Request:"ileft",Dst:me.FingerTable[i][2],Id:id,Suc:suc,SucAdr:sucAdr,Pre:pre,PreAdr:preAdr,Inf:offset}//id=id key=pre inf=suc	
				go poster(msg)				
			}			
			//for each successor, look if it maches the id
			if id == me.FingerTable[i][1] {
				//compare t with id
				if me.FingerTable[i][0]<= id { //give successor
					me.FingerTable[i][1]=suc
					me.FingerTable[i][2]=sucAdr
				}else{// give predecessor
					me.FingerTable[i][1]=pre
					me.FingerTable[i][2]=preAdr					
				}
			}
		}		
	}
		
}
	
/*
	Used by the client part of the node to store a "file" on the ring
	it calculates the hash and id of the "file"
	finds a responsible mode
	sends the "file" to the responsible node
*/
func insert(file string){

	//generate hash and id
	hash:=hash(file)
	fmt.Println("Hash key: ", hash)
	id:=Index(hash, BITS)
	//fing the responsible key for the id
	s,a:=askForAddress(id.String(), me.address)
	fmt.Print("Reposible: ", s)
	fmt.Println(" Address: ", a)
	//send the value to the responsible node	
	msg:=Msg{Request:"Insert",Dst:a,Key:hash,Inf:file}
	_=poster(msg)
}

/*
	Used by the nodes to post messanges using UDP to each other
*/
func poster(msg Msg) Msg{
	//send the request and receive responce
	t:=Transport{}
	resp:=t.Send(&msg)
	var msgresponce	Msg	
	err:= json.Unmarshal(resp, &msgresponce)
	if err==nil {}		
	return msgresponce	
}

/*
	Used to calculate the hash of a "file"
*/
func hash(file string) string{
	hasher := sha1.New()
	hasher.Write([]byte(file))
	return fmt.Sprintf("%x", hasher.Sum(nil))	
}

/*
	Used to find a "file" in the ring using a hash
*/
func find(hash string){
	//get the id
	id:=Index(hash, BITS)	
	//find responible node
	s,a:=askForAddress(id.String(), me.address)
	fmt.Print("Reposible: ", s)
	fmt.Println(" Address: ", a)	
	msg:=Msg{Request:"find",Dst:a,Key:hash}
	msgresponce:=poster(msg)
	fmt.Println("Found : ", msgresponce.Inf)	
}

/*
	Used by remove a "file" from the ring
*/
func remove(file string){
	//generate hash of the file
	hash:=hash(file)
	//get the id
	id:=Index(hash, BITS)	
	//find responible node
	s,a:=askForAddress(id.String(), me.address)
	fmt.Print("Reposible node: ", s)
	fmt.Println(" address: ", a)	
	msg:=Msg{Request:"remove",Dst:a,Key:hash}
	_=poster(msg)	
}
	
// USER INTERACTION CODES
func Menu(){
	fmt.Println("\nSELECT ACTION ON NODE ", me.nodeId.String())
	fmt.Println("1 = INSERT")
	fmt.Println("2 = FIND")
	fmt.Println("3 = FINGER TABLE")
	fmt.Println("4 = LEAV")
	fmt.Println("5 = UPDATE")
	fmt.Println("6 = INVENTORY")
	fmt.Println("7 = DELETE")
	var choice string
	fmt.Scanln(&choice)
	selection(choice)
}

func selection(choice string){
	//determine and execute the selected operation
		switch choice {
			case "1":
				{
					fmt.Println("Enter the name of the file to Insert:")
					var file string
    				fmt.Scanln(&file)
					insert(file)
				}
			case "2":
					fmt.Println("Enter the name of the file to Search:")
					var file string
    				fmt.Scanln(&file)
					find(file)
			case "3":
				fmt.Println(me.FingerTable)
			case "4":
					leav()			
			case "5":
			{
				updateFingerTable(me.nodeId,me.FingerTable[0][2]) //nodeid, target server (successor) address required				
			}
			case "6":
				fmt.Println(me.storage)	
			case "7":
				{
					fmt.Println("Enter the name of the file to delete:")
					var file string
    				fmt.Scanln(&file)
					remove(file)
				}				
						
		}	
	//clear console screen
}
/*
	Used by a node to self update after every interval
	improve the consistancy of the nodes on the sring
*/
func AutoUpdate(){
	time.Sleep(20000 * time.Millisecond)
	updateFingerTable(me.nodeId,me.FingerTable[0][2]) 
	AutoUpdate()
}
/*

	Used to detect crashing nodes nad re direct the request to the successor node
	failen implementation untill the puplishing time of these codes
	difficulties in detecting successor and predecessor on the crashed node
	

 func Stabilize(id string,msg *Msg) {
	 fmt.Print("Stabilize ")
	 fmt.Print(" id ",id)
	 // convert the id to big int
	 one := big.NewInt(2)
	 i := new(big.Int)
	 i.SetString(id, 10)	 
	 //find the successor ie. imcrement by err1
	 i.Add(i,one) 		//incorrect way to find successor
	 //get the address of this successor
	 s,a:=askForAddress(i.String(), me.address)
	//modify message DestIdfmt.
	fmt.Print("s ", s)
	fmt.Print(" a ", a)
	msg.DstId=s
	msg.Dst=a
	//redend message
	fmt.Println(" msg ", *msg)
	_=poster(*msg)
 }
 */
//