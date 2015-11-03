/* UDPDaytimeClient
 */
package main

import (
	"os"
	"fmt"
	"flag"
	"github.com/hamzamac/chord/lib"
)



func main() {
	// flag definition
	address := flag.String("address", "localhost:1200", "ip:port")
	 join := flag.Bool("join", true, "true|false")
	 bnode := flag.String("bnode", "", "address of the boosting node")
	 flag.Parse()
	 //initial settings
	 dht.SetAddress(*address)
	 dht.SetTarget(*bnode)
	//startserver part ie Listen
	ts:=dht.Transport{*address}
	go ts.Listen()
	
	//get id and join the ring	
	if(*join){ //join existing ring
		dht.Join(false)
		dht.IncrementOffset()
		dht.Amx(dht.GetId(), *address, *bnode,dht.GetOffsetServer())			
	}else { //first node in the ring
		dht.Join(true)
	}
	for{
			dht.Menu()
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}