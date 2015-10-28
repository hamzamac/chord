package dht
import (
	"fmt"
	"strconv"
	"math/big"
	"encoding/json"
	"crypto/sha1"
	//"math/rand"
)
const SIZE int = 8 
const BITS int = 3


type DHTNode struct {
	nodeId      string
	address 	string
	target 		string
	successor   *DHTNode
	predecessor *DHTNode
	contact     Contact
	offset		int
	FingerTable [BITS][BITS]string		//figer table array
	storage  map[string]string			//inmeory storage
}


// GLOBAL INITIALIZATION: in memory storage
	var me DHTNode						//this clients node
	//me.storage=map[string]string{}
	//set functions
	func SetAddress(adrs string){
		me.address=adrs
	}
	func SetTarget(tgt string){
		me.target=tgt
	}
	func SetOffset(ofst int){
		me.offset=ofst
	}
	func IncrementOffset(){
		me.offset++
	}	
// NEW CODES FOR LIVE IMPLEMENTATION
	//request Join

	func Join(test bool){
		//fmt.Println("Join")
		//generate id
		me.nodeId=Index(generateNodeId(), BITS)
		nodeIdInt,_ := strconv.ParseInt(me.nodeId, 10, 0)
		me.storage=map[string]string{}
		fmt.Println("assigned id: ",nodeIdInt)		
		//prepare finger table
		if(test){ //firstnode in the ring
				me.offset=0
			for i,_:= range me.FingerTable {
				//compute the ith node
				two := big.NewInt(2)
				dist := big.Int{}
				dist.Exp(two, big.NewInt(int64(i)), nil)
				distInt,_ := strconv.ParseInt(dist.String(), 10, 0)
				ithNode:=(distInt+nodeIdInt)%8 //correction
				//fill the fingr table row
				me.FingerTable[i][0]= strconv.FormatInt(ithNode, 10)
				//suc,address:=askForAddress(strconv.FormatInt(ithNode, 10),"localhost:1200")
				me.FingerTable[i][1]=me.nodeId
				me.FingerTable[i][2]=me.address //this node address
			}
		}else{	//not first node
			//initialize offset
			GetOffsetClient(me.target)		
			updateFingerTable(nodeIdInt,me.target) //target server required
		}
		fmt.Println("Finger Table: ", me.FingerTable)	
	}

	func updateFingerTable(nodeIdInt int64, questioneeAddress string){// node's assigned id, questionee node address
		//fmt.Println("updateFingerTable")
		for i,_:= range me.FingerTable {
			//compute the ith node
			two := big.NewInt(2)
			dist := big.Int{}
			dist.Exp(two, big.NewInt(int64(i)), nil)
			distInt,_ := strconv.ParseInt(dist.String(), 10, 0)

			ithNode:=(distInt+nodeIdInt)%8 //correction to 8
			//fill the fingr table row
			me.FingerTable[i][0]=strconv.FormatInt(ithNode, 10)
			suc,address:=askForAddress(strconv.FormatInt(ithNode, 10),questioneeAddress)
			me.FingerTable[i][1]=suc
			me.FingerTable[i][2]=address
		}		
	}	
	func askForAddress(id string, address string) (string, string){ // successor, address
		//send the request and receive responce
		msg:=Msg{Request:"getAddress",Dst:address,Id:id}
		t:=Transport{}
		resp:=t.Send(&msg)
		var msgresponce	Msg	
		err:= json.Unmarshal(resp, &msgresponce)
		if err==nil {}		
		//test the response
		if(msgresponce.Inf=="found"){		//found
			return msgresponce.Id, msgresponce.Address
		}	//not found, ask another node that is not me
		fmt.Println("asking ",msgresponce.Id)
		return askForAddress(msgresponce.Id, msgresponce.Address)	
	}

	func Amx(id string, myaddress string, destAddress string, offset string){//id,sourceaddress,destAddress
		//send the request and receive responce
		msg:=Msg{Request:"Joined",Dst:destAddress,Src:myaddress,Id:id,Inf:offset}
		t:=Transport{}
		resp:=t.Send(&msg)
		var msgresponce	Msg	
		err:= json.Unmarshal(resp, &msgresponce)
		if err==nil {}		
		//test the response		
	}

	func nodeDistance(s int,d int) int{
		souece:= big.NewInt(int64(s))
		dest:= big.NewInt(int64(d))
		r, _ := strconv.ParseInt(distance(souece.Bytes(), dest.Bytes(), BITS).String(), 10, 0)
		return int(r)
	}
	func GetOffsetClient(destAddress string){
		//send the request and receive responce
		msg:=Msg{Request:"offset",Dst:destAddress}
		t:=Transport{}
		resp:=t.Send(&msg)
		var msgresponce	Msg	
		err:= json.Unmarshal(resp, &msgresponce)
		if err==nil {}		
		//test the response
		o, _ := strconv.ParseInt(msgresponce.Inf, 10, 0)	
		SetOffset(int(o))		
	}


// SERVER SIDE CODES
	func Lookup(id string,fingertable [BITS][BITS]string)(string,string,string){//successor,address,inf
	//if are you looking for me?
	if id == me.nodeId {
		fmt.Println("looking for me? ", id)
		return me.nodeId, me.address, "found"
	}
	l:=len(fingertable)-1
		//find if id is beween the ringer references
		answer4s:=false
		answer4t:=between(str2byte(fingertable[0][0]),str2byte(fingertable[l][0]),str2byte(id))
		if fingertable[0][1]!=fingertable[l][1] {
			answer4s=between(str2byte(fingertable[0][1]),str2byte(fingertable[l][1]),str2byte(id))
		}
				
		fmt.Print("looking in ",me.nodeId)
		if answer4s {
			fmt.Print("found in s ",id)
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
			fmt.Println(" found", closest[0][1])
			return closest[0][1],closest[0][2], "found"	
			// end of find closest				
			
		}else if answer4t {
			fmt.Print("found in t ",id)
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
			fmt.Println(" found", closest[0][1])
			return closest[0][1],closest[0][2], "found"	
			// end of find closest		
		}else{
			// check is the successor poits to other nodes than this
			fmt.Println("all false")
			
			var test = false
			idx:=0
			var tempft[BITS][BITS] string
			for i,_:=range fingertable {
				if me.nodeId!=fingertable[i][1] {
					test= true
					//add none me finger to temp array
					tempft[idx]=fingertable[i]
					idx++
				}
			}
			
			le:=len(tempft)-1
			//ask another node that is not me CORRECTION NEEDED
			if test {// ask the successor of the another node that is not me
				//find the last Index, l
					
				//return the last successor
				//return fingertable[l][1],fingertable[l][2], "notfound"
				fmt.Println("not found ask ", tempft[le][1])
				return id,tempft[le][2], "notfound"
			}else{// its only me here
				return me.nodeId,me.address, "found"
			}			
		}
	}

func Update(id string, address string, fingertable *[BITS][BITS]string, offset string) string{ //id, address, fingertable
	//fmt.Println("Update")
	// confirm the ofset is grater
	ofst, _ := strconv.ParseInt(offset, 10, 0)
	if me.offset< int(ofst) {
		me.offset=int(ofst)
		for i,_:= range fingertable{
			// neither me of the sender -> notigy on the new node
			if fingertable[i][1]!=id && fingertable[i][1]!=me.nodeId {
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
	func GetId() string{
		return me.nodeId
	}

	func GetOffsetServer() string{
		o:= strconv.FormatInt(int64(me.offset), 10)
		return o		
	}
	
	func store(key string, value string,replicate bool){
		//copie the value to inmemory storage
		me.storage[key]=value
		if replicate {
			//get this nodes successor
			//update the successor
			msg:=Msg{Request:"Replicate",Dst:me.FingerTable[0][2],Key:key,Inf:value}
			_=poster(msg)			
		}
	}
	func deleteKey(hash string, dereplicate bool){
		//derete value from  inmemory storage
		delete(me.storage, hash)
		if dereplicate {
			//get this nodes successor
			//update the successor
			msg:=Msg{Request:"dereplicate",Dst:me.FingerTable[0][2],Key:hash}
			_=poster(msg)			
		}		
	}
// USER INTERACTION CODES
func Menu(){
	fmt.Println("\nSELECT ACTION")
	fmt.Println("1 = INSERT")
	fmt.Println("2 = FIND")
	fmt.Println("3 = FINGER TABLE")
	fmt.Println("4 = LEAV")
	fmt.Println("5 = UPDATE")
	fmt.Println("6 = INVENTORY")
	fmt.Println("6 = DELETE")
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
				fmt.Println("2")
			case "3":
				fmt.Println(me.FingerTable)
			case "4":
				fmt.Println("3")				
			case "5":
			{
				nodeIdInt,_ := strconv.ParseInt(me.nodeId, 10, 0)
				updateFingerTable(nodeIdInt,me.FingerTable[0][2]) //nodeid, target server (successor) address required				
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

func insert(file string){
	//generate hash and id
	/*hasher := sha1.New()
	hasher.Write([]byte(file))
	hash:=fmt.Sprintf("%x", hasher.Sum(nil))*/
	hash:=hash(file)
	fmt.Println("Hash key: ", hash)
	id:=Index(hash, BITS)
	fmt.Println("identifier: ", id)	
	//fing the responsible key for the id
	s,a:=askForAddress(id, me.address)
	fmt.Println("reposible: ", s)
	fmt.Println("address: ", a)
	//send the value to the responsible node	
	msg:=Msg{Request:"Insert",Dst:a,Key:hash,Inf:file}
	_=poster(msg)
}
func poster(msg Msg) Msg{
	//send the request and receive responce
	t:=Transport{}
	resp:=t.Send(&msg)
	var msgresponce	Msg	
	err:= json.Unmarshal(resp, &msgresponce)
	if err==nil {}		
	return msgresponce	
}
func hash(file string) string{
	hasher := sha1.New()
	hasher.Write([]byte(file))
	return fmt.Sprintf("%x", hasher.Sum(nil))	
}
func find(file string){
	//generate hash of the file
	hash:=hash(file)
	//get the id
	id:=Index(hash, BITS)	
	//find responible node
	s,a:=askForAddress(id, me.address)
	fmt.Print("reposible: ", s)
	fmt.Println(" address: ", a)	
	//msg:=Msg{Request:"find",Dst:a,Key:hash,Inf:file}
	//_=poster(msg)	
}
func remove(file string){
	//generate hash of the file
	hash:=hash(file)
	//get the id
	id:=Index(hash, BITS)	
	//find responible node
	s,a:=askForAddress(id, me.address)
	fmt.Print("reposible: ", s)
	fmt.Println(" address: ", a)	
	msg:=Msg{Request:"remove",Dst:a,Key:hash}
	_=poster(msg)	
}

// OLD CODES FROM SIMULATION
type Contact struct {
	ip   string
	port string
}


var Ring [SIZE] *DHTNode



func makeDHTNode(nodeId *string, ip string, port string) *DHTNode {
	dhtNode := new(DHTNode)
	dhtNode.contact.ip = ip
	dhtNode.contact.port = port

	if nodeId == nil {
		genNodeId := generateNodeId()
		dhtNode.nodeId = genNodeId
	} else {
		dhtNode.nodeId = *nodeId
	}

	dhtNode.successor = nil
	dhtNode.predecessor = nil

	return dhtNode
}

func (dhtNode *DHTNode) addToRing(newDHTNode *DHTNode) {
	i, _ := strconv.ParseInt(newDHTNode.nodeId, 10, 0)
	Ring[i]=newDHTNode
}

func (dhtNode *DHTNode) lookup(key string) *DHTNode {
	// TODO
	
	
	return dhtNode // XXX This is not correct obviously
}

func (dhtNode *DHTNode) acceleratedLookupUsingFingers(key string) *DHTNode {
	// TODO
	return dhtNode // XXX This is not correct obviously
}

func (dhtNode *DHTNode) responsible(key string) bool {
	// TODO
	return false
}

func (dhtNode *DHTNode) printRing() {	
	id, _ := strconv.ParseInt(dhtNode.nodeId, 10,0)
	for  i := int(id); i < int(id)+SIZE; i++  {
		if Ring[i%SIZE]!= nil{
			fmt.Println(i%SIZE)
		}
	}
}

func updateRing() {

	for  i := 0; i < SIZE; i++  { // i represent a position
		
		if Ring[i%SIZE]!= nil{ //if node exists, find it peers		
			findPeers(i) // find peers
		}
	}
}
func findPeers(nodeId int){
		var sucfound = false
		var predfound = false		
		
			for x:= 1; x<=(SIZE/2); x++ {
				
				if Ring[(nodeId+x)%SIZE]!=nil && !sucfound {  //find successor
					Ring[nodeId].successor=Ring[(nodeId+x)%SIZE]
					sucfound=true
				}
				if Ring[(nodeId-x+SIZE)%SIZE]!=nil && !predfound {	//find predecessor
					Ring[nodeId].predecessor=Ring[(nodeId-x+SIZE)%SIZE]
					predfound=true
				}
				if sucfound==true && predfound==true { // stope when all peers are found
					x=(SIZE/2)
				}
			}	
}

func (dhtNode *DHTNode) printFinger(fingerIndex int, bits int) *DHTNode {
	fmt.Println("-> printFinger")
	// convert id to interger
	nid, _ := strconv.ParseInt(dhtNode.nodeId, 10, 0)
	//convert integer to byte
	id:= big.NewInt(nid)
	//calculate finger table value
	ni,_:=calcFinger(id.Bytes(), fingerIndex, BITS)
	i, _ := strconv.ParseInt(ni, 10, 0)
	// verify finger
	return verifyFinger(int(i)) 
}

func verifyFinger(id int) *DHTNode{
	fmt.Println("-> verifyFinger")
	if Ring[id]==nil {
		id++
		return verifyFinger(id)
	}
	//fmt.Println(Ring[id].nodeId)
	return Ring[id]
}



func (dhtNode *DHTNode) fingerLookup(target int,bits int) *DHTNode {
	bits--
	fmt.Println("->lookup start")
	//fmt.Println(bits)
	
	//jump half of the ring in invers exponential
	id1, _ := strconv.ParseInt(dhtNode.nodeId, 10, 0)
	id1bint:= big.NewInt(int64(id1))
	
	two := big.NewInt(2)
	id2 := big.Int{}
	id2.Exp(two, big.NewInt(int64(bits)), nil)
	
	x1, _ := strconv.ParseInt(id2.String(), 10, 0)
	x2, _ := strconv.ParseInt(id1bint.String(), 10, 0)
	x:=(x1 + x2)%int64(SIZE)
	id2bint:= big.NewInt(x)
	
	fmt.Println(id2bint)
	
	key:= big.NewInt(int64(target))
	
	result:=between(id1bint.Bytes(), id2bint.Bytes(), key.Bytes())
	fmt.Println(result)
	//id2int,_ := strconv.ParseInt(id2.String(), 10, 0)
	//if result is false continue lookup
	if result==false {
		fmt.Print("look from ")
		fmt.Println(Ring[int(x)].nodeId)
		return Ring[int(x)].fingerLookup(target,bits)
	}
	
	//find closent from finger table
	
	closest:= Ring[int(x2)].findClosest(target)
	//tect if is the target
	return closest	
}	

func (dhtNode *DHTNode) findClosest(target int) *DHTNode {
	fmt.Println("->findClosest start")
	//var a [10]int	
	//get the caller id
	soueceId, _ := strconv.ParseInt(dhtNode.nodeId, 10, 0)
	//fmt.Println(dhtNode.nodeId)
	var referencId int64
	d1:=nodeDistance(int(soueceId),target)
	//get the finger table and compute distances from
	for i:=1; i<=BITS; i++ {
		//get finger reference id
		referencId, _ = strconv.ParseInt(dhtNode.printFinger(i, BITS).nodeId, 10, 0)
		//calculate distance from source
		d2:=nodeDistance(int(soueceId),int(referencId))
		
		if d2>=d1 {
			return Ring[referencId]
		}
		
	}
	return Ring[referencId]
}
/*
func Range(id,s string) bool{
	//convert to interger
	
	//compute the opposite

	return "Index(generateNodeId(), BITS)"
}
*/

