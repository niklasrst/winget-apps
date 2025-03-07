package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

func main() {
	if len(os.Args) < 2 || (len(os.Args) > 1 && os.Args[1] == "/help") {
		fmt.Println("This tool converts EntraID group ObjectIDs to SIDs and vice versa.")
		fmt.Println("Usage: groupconverter <Object-to-Sid|Sid-to-Object> <ObjectID|SID>")
		fmt.Println("Example: groupconverter Object-to-Sid 3f2504e0-4f89-11d3-9a0c-0305e82c3301")
		return
	}

	direction := os.Args[1]
	id := os.Args[2]

	switch direction {
	case "Object-to-Sid":
		sid, err := objectToSid(id)
		if err != nil {
			log.Fatalf("Error converting ObjectID to SID: %v", err)
		}
		fmt.Println(sid)
	case "Sid-to-Object":
		objectID, err := sidToObject(id)
		if err != nil {
			log.Fatalf("Error converting SID to ObjectID: %v", err)
		}
		fmt.Println(objectID)
	default:
		log.Fatalf("Invalid direction. Please choose between Object-to-Sid or Sid-to-Object")
	}
}

func objectToSid(objectID string) (string, error) {
	u, err := uuid.Parse(objectID)
	if err != nil {
		return "", err
	}

	bytes := u[:]
	array := make([]uint32, 4)
	for i := 0; i < 4; i++ {
		array[i] = binary.LittleEndian.Uint32(bytes[i*4 : (i+1)*4])
	}

	sid := fmt.Sprintf("S-1-12-1-%d-%d-%d-%d", array[0], array[1], array[2], array[3])
	return sid, nil
}

func sidToObject(sid string) (string, error) {
	if !strings.HasPrefix(sid, "S-1-12-1-") {
		return "", errors.New("invalid SID format")
	}

	parts := strings.Split(sid[9:], "-")
	if len(parts) != 4 {
		return "", errors.New("invalid SID format")
	}

	array := make([]uint32, 4)
	for i, part := range parts {
		var value uint32
		_, err := fmt.Sscanf(part, "%d", &value)
		if err != nil {
			return "", err
		}
		array[i] = value
	}

	bytes := make([]byte, 16)
	for i := 0; i < 4; i++ {
		binary.LittleEndian.PutUint32(bytes[i*4:], array[i])
	}

	u, err := uuid.FromBytes(bytes)
	if err != nil {
		return "", err
	}

	return u.String(), nil
}
