package dht
import (
	"fmt"
	"strconv"
	"math/big"
)

type Contact struct {
	ip   string
	port string
}

type DHTNode struct {
	nodeId      string
	successor   *DHTNode
	predecessor *DHTNode
	contact     Contact
}

const SIZE int = 8 
const BITS int = 3 

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

func nodeDistance(s int,d int) int{
	souece:= big.NewInt(int64(s))
	dest:= big.NewInt(int64(d))
	r, _ := strconv.ParseInt(distance(souece.Bytes(), dest.Bytes(), BITS).String(), 10, 0)
	return int(r)
}

func (dhtNode *DHTNode) fingerLookup(target int,bits int) *DHTNode {
	fmt.Println("->lookup start")
	bits--
	//jump half of the ring in invers exponential
	id1, _ := strconv.ParseInt(dhtNode.nodeId, 10, 0)
	id1bint:= big.NewInt(int64(id1))
	
	two := big.NewInt(2)
	id2 := big.Int{}
	id2.Exp(two, big.NewInt(int64(bits)), nil)
	
	x1, _ := strconv.ParseInt(id2.String(), 10, 0)
	x2, _ := strconv.ParseInt(id1bint.String(), 10, 0)
	x:=x1 + x2
	id2bint:= big.NewInt(int64(x)%int64(SIZE))
	
	fmt.Println(id2bint)
	
	key:= big.NewInt(int64(target))
	
	result:=between(id1bint.Bytes(), id2bint.Bytes(), key.Bytes())
	fmt.Println(result)
	id2int,_ := strconv.ParseInt(id2.String(), 10, 0)
	//if result is false continue lookup
	if result==false {
		fmt.Println("false, look for")
		fmt.Println(Ring[id2int].nodeId)
		return Ring[id2int].fingerLookup(target,bits)
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
		
/*		//calculate the difference between target and source
		d0:=d1-d2
		//add positive difference to array
		if d0>=0 {
			a[d0]=int(destId)
		}
*/
	}
	return Ring[referencId]
}
/*
func minimumpositive(x [10] int) *DHTNode{
	fmt.Println("->closest start")
	min := 0
    for i,_:= range x {
        if i < min {
            min = i
        }
    }
	
	return Ring[x[min]]
}
*/



func (dhtNode *DHTNode) testCalcFingers(m int, bits int) {
/*	 idBytes, _ := hex.DecodeString(dhtNode.nodeId)
	fingerHex, _ := calcFinger(idBytes, m, bits)
	fingerSuccessor := dhtNode.lookup(fingerHex)
	fingerSuccessorBytes, _ := hex.DecodeString(fingerSuccessor.nodeId)
	fmt.Println("successor    " + fingerSuccessor.nodeId)

	dist := distance(idBytes, fingerSuccessorBytes, bits)
	fmt.Println("distance     " + dist.String()) */
}
