package dht
import (
	"net"
	"encoding/json"
	"fmt"
	"time"
)


type Transport struct{
	BindAddress string
}

type Msg struct{
	Request string
	Key string
	Id string
	Address string
	Inf string
	Src string
	Dst	string
	DstId	string
	Suc string
	Pre string
	SucAdr string
	PreAdr string
}

/*
	Controll the server part of a node, 
	Always listsening for new connection
*/
func(transport *Transport) Listen(){
	//fmt.Println("Listening")
	udpAddr,err:=net.ResolveUDPAddr("udp",transport.BindAddress)
	checkError(err)	
	conn,err:=net.ListenUDP("udp",udpAddr)
	checkError(err)
	defer conn.Close()

	for{
		 handleClient(conn)	
	}
	checkError(err)
}

/*
	Used by the client part of a node
	Send messages to other node
*/
func(transport *Transport) Send(msg	*Msg) []byte{
	//send request
	//fmt.Println("send")
	udpAddr,err:=net.ResolveUDPAddr("udp",msg.Dst)
	conn,err:=net.DialUDP("udp",nil,udpAddr)
	if err!=nil {
	//	Stabilize(msg.DstId,msg)		
	}
	defer conn.Close()
	jmsg, _ := json.Marshal(msg)
	_,err=conn.Write([]byte(jmsg))
	checkError(err)
	
	//response
	var buf [512]byte
	conn.SetReadDeadline(time.Now().Add(120000 * time.Millisecond))
	n, err := conn.Read(buf[0:])
	checkError(err)
	
	return buf[0:n]
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Not Alive ", err)
		//os.Exit(1)
	}
}

/*
	Used by the server part of the node to handle individual connections
	initiate a response for a given task
*/
func handleClient(conn *net.UDPConn) {
	//fmt.Println("handleClient")
	//receive message
	var buf [512]byte

	n, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}
	
	//decode message into message object
	var msg *Msg
	err1 := json.Unmarshal(buf[0:n], &msg)
	checkError(err1)
	
	//process the message object
		switch msg.Request {
			case "getAddress":
				msg.Id,msg.Address,msg.Inf=Lookup(msg.Id, me.FingerTable)
			case "offset":
				msg.Inf=GetOffsetServer()
			case "Joined":
				msg.Inf=Update(msg.Id, msg.Src, &me.FingerTable,msg.Inf)
			case "Insert":
				store(msg.Key,msg.Inf, true)
			case "Replicate":
				store(msg.Key,msg.Inf, false)
			case "remove":
				deleteKey(msg.Key, true)
			case "dereplicate":
				deleteKey(msg.Key, false)
			case "find":
				msg.Inf=getValue(msg.Key)
			case "ileft":
				hasLeft(msg.Id,msg.Suc,msg.SucAdr,msg.Pre,msg.PreAdr,msg.Inf)	//id=id key=pre inf=suc	
			case "predecessor":
				setPredecessor(msg.Id,msg.Inf)																				
		}
	jmsg, _ := json.Marshal(msg)
	//return response
	conn.WriteToUDP([]byte(jmsg), addr)
}

