package dht

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"math/big"
	"strconv"
)

func distance(a, b []byte, bits int) *big.Int {
	var ring big.Int
	ring.Exp(big.NewInt(2), big.NewInt(int64(bits)), nil)

	var a_int, b_int big.Int
	(&a_int).SetBytes(a)
	(&b_int).SetBytes(b)

	var dist big.Int
	(&dist).Sub(&b_int, &a_int)

	(&dist).Mod(&dist, &ring)
	return &dist
}

func between(id1, id2, key []byte) bool {
	// 0 if a==b, -1 if a < b, and +1 if a > b

	if bytes.Compare(key, id1) == 0 { // key == id1
		return true
	}
	
	if bytes.Compare(key, id2) == 0 { // key == id2
		return true
	}	

	if bytes.Compare(id2, id1) == 1 { // id2 > id1
		if bytes.Compare(key, id2) == -1 && bytes.Compare(key, id1) == 1 { // key < id2 && key > id1
			return true
		} else {
			return false
		}
	} else { // id1 > id2
		if bytes.Compare(key, id1) == 1 || bytes.Compare(key, id2) == -1 { // key > id1 || key < id2
			return true
		} else {
			return false
		}
	}
}
/*
	generate a random hash a for the joining node
*/
func generateNodeId() string {
	u, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	// calculate sha-1 hash
	hasher := sha1.New()
	hasher.Write([]byte(u.String()))

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

/*
	Truncate the hash to  a number between 0 and 2^bits -1
*/
func Index(key string, bits int) big.Int{
	k:=[]byte(key)
	sum:=0
	for  _, num :=range k {
		s := fmt.Sprintf("%d", num)		
		sint, _ := strconv.Atoi(s)		
		sum+=int(sint)
	}
	
	two := big.NewInt(2)
	id2 := big.Int{}
	id2.Exp(two, big.NewInt(int64(bits)), nil)
	s:=big.NewInt(int64(sum))

	bint := big.Int{}
	bint.Mod(s,&id2)
	return bint
}

/*
	converts string to byte array
*/
func str2byte(str string) []byte{
	x, _ := strconv.ParseInt(str, 10, 0)
	y:= big.NewInt(x)
	return y.Bytes()
}

/*
	calculates 2^bits
*/
func ringSize(bits int)big.Int{
			two2 := big.NewInt(2)
			id3 := big.Int{}
			id3.Exp(two2, big.NewInt(int64(BITS)), nil)		
			return 	id3
}