package proxy

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// dnsResourceRecord describes individual records in the request and response of the DNS payload body
type dnsResourceRecord struct {
	DomainName         string
	Type               uint16
	Class              uint16
	TimeToLive         uint32
	ResourceDataLength uint16
	ResourceData       []byte
}

// dnsHeader describes the request/response DNS header
type dnsHeader struct {
	TransactionID  uint16
	Flags          uint16
	NumQuestions   uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

func parseDNSRequest(raw []byte) (dnsHeader, []dnsResourceRecord, error) {
	var requestBuffer = bytes.NewBuffer(raw)
	var queryHeader dnsHeader
	var queryResourceRecords []dnsResourceRecord

	err := binary.Read(requestBuffer, binary.BigEndian, &queryHeader) // network byte order is big endian

	if err != nil {
		return queryHeader, queryResourceRecords, fmt.Errorf("Error decoding header: ", err.Error())
	}

	queryResourceRecords = make([]dnsResourceRecord, queryHeader.NumQuestions)

	for i := range queryResourceRecords {
		qrr := dnsResourceRecord{}
		qrr.DomainName, err = readDomainName(requestBuffer)

		if err != nil {
			return queryHeader, queryResourceRecords, fmt.Errorf("Error decoding label: ", err.Error())
		}

		qrr.Type = binary.BigEndian.Uint16(requestBuffer.Next(2))
		qrr.Class = binary.BigEndian.Uint16(requestBuffer.Next(2))

		queryResourceRecords[i] = qrr
	}

	return queryHeader, queryResourceRecords, nil
}

// RFC1035: "Domain names in messages are expressed in terms of a sequence
// of labels. Each label is represented as a one octet length field followed
// by that number of octets.  Since every domain name ends with the null label
// of the root, a domain name is terminated by a length byte of zero."
func readDomainName(requestBuffer *bytes.Buffer) (string, error) {
	var domainName string

	b, err := requestBuffer.ReadByte()

	for ; b != 0 && err == nil; b, err = requestBuffer.ReadByte() {
		labelLength := int(b)
		labelBytes := requestBuffer.Next(labelLength)
		labelName := string(labelBytes)

		if len(domainName) == 0 {
			domainName = labelName
		} else {
			domainName += "." + labelName
		}
	}

	return domainName, err
}
