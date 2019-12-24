package proxy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	FlagResponse uint16 = 1 << 15
)

func writeNoFound(queryHeader dnsHeader, reqs []dnsResourceRecord) []byte {
	var responseBuffer = new(bytes.Buffer)
	var responseHeader dnsHeader

	responseHeader = dnsHeader{
		TransactionID:  queryHeader.TransactionID,
		Flags:          FlagResponse,
		NumQuestions:   uint16(len(reqs)),
		NumAnswers:     0,
		NumAuthorities: 0,
		NumAdditionals: 0,
	}

	err := binary.Write(responseBuffer, binary.BigEndian, &responseHeader)

	if err != nil {
		fmt.Println("Error writing to buffer: ", err.Error())
	}

	for _, queryResourceRecord := range reqs {
		err = writeDomainName(responseBuffer, queryResourceRecord.DomainName)

		if err != nil {
			fmt.Println("Error writing to buffer: ", err.Error())
		}

		binary.Write(responseBuffer, binary.BigEndian, queryResourceRecord.Type)
		binary.Write(responseBuffer, binary.BigEndian, queryResourceRecord.Class)
	}

	return responseBuffer.Bytes()
}

// RFC1035: "Domain names in messages are expressed in terms of a sequence
// of labels. Each label is represented as a one octet length field followed
// by that number of octets.  Since every domain name ends with the null label
// of the root, a domain name is terminated by a length byte of zero."
func writeDomainName(responseBuffer *bytes.Buffer, domainName string) error {
	labels := strings.Split(domainName, ".")

	for _, label := range labels {
		labelLength := len(label)
		labelBytes := []byte(label)

		responseBuffer.WriteByte(byte(labelLength))
		responseBuffer.Write(labelBytes)
	}

	err := responseBuffer.WriteByte(byte(0))

	return err
}
