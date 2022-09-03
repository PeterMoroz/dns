package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"strconv"
)

var recTypes = make(map[string]uint16)

const (
	Header_QR_offset = 15
	Header_AA_offset = 10
	Header_TC_offset = 9
	Header_RR_offset = 8
	Header_RA_offset = 7
)

const (
	Rcode_NoError = iota
	Rcode_FormatError
	Rcode_ServerFailure
	Rcode_NameError
	Rcode_NotImplemented
	Rcode_Refused
)


func init() {
	recTypes = map[string] uint16 {
		"A": 1,
		"NS": 2,
		"CNAME": 5,
	}
}

type Query struct {
	Header
	Question
}

type Response struct {
	Header
	Question
	
	answers []Answer
	authorities []Answer
	additional []Answer
}

type Header struct {
	id uint16
	flags uint16
	
	qdcount uint16
	ancount uint16
	nscount uint16
	arcount uint16
}

func (h *Header) ToString() string {
	return fmt.Sprintf("id: %05d, flags: %04X, qdcount: %d, ancount: %d, nscount: %d, arcount: %d",
		h.id, h.flags, h.qdcount, h.ancount, h.nscount, h.arcount);
}

func (h *Header) GetId() uint16 {
	return h.id
}

func (h *Header) SetId(id uint16) {
	h.id = id
}

func (h *Header) SetIsResponse(val bool) {
	if val {
		h.flags |= (1 << Header_QR_offset)
	} else  {
		h.flags = h.flags &^ (1 << Header_QR_offset)
	}
}

func (h *Header) SetAA(val bool) {
	if val {
		h.flags |= (1 << Header_AA_offset)
	} else  {
		h.flags = h.flags &^ (1 << Header_AA_offset)
	}
}

// TO DO: set other flags

func (h *Header) SetRcode(val uint8) {
	h.flags &= 0xFFF0
	h.flags |= uint16(val & 0xF)
}


func (h *Header) SetQdcount(qdcount uint16) {
	h.qdcount = qdcount
}

func (h *Header) SetAncount(ancount uint16) {
	h.ancount = ancount
}

func (h *Header) SetNscount(nscount uint16) {
	h.nscount = nscount
}

func (h *Header) SetArcount(arcount uint16) {
	h.arcount = arcount
}

type Question struct {
	qname string
	qtype uint16
	qclass uint16
}

func (q *Question) GetName() string {
	return q.qname
}

func (q *Question) GetType() uint16 {
	return q.qtype
}

func (q *Question) GetClass() uint16 {
	return q.qclass
}

type Answer struct {
	name string
	atype uint16
	aclass uint16
	ttl uint32
	rdata []byte
}

func rdataFromIPv4(addr string) []byte {
	rdata := make([]byte, 4)
	octets := strings.Split(addr, ".")
	
	if len(octets) != 4 {
		log.Fatalf("Invalid address (IPv4) string %s", addr)
	}

	for i := 0; i < 4; i++ {
		n, err := strconv.Atoi(octets[i])
		if err != nil {
			log.Fatalf("Invalid address (IPv4) string %s", addr)
		}
		rdata[i] = byte(n)
	}
	
	return rdata
}

func rdataFromDomainName(name string) []byte {
	rdata := make([]byte, len(name) + 1)
	labels := strings.Split(name, ".")
	
	idx := 0
	for i, _ := range labels {
		rdata[idx] = byte(len(labels[i]))
		idx += 1
		copy(rdata[idx:], labels[i]) // rdata = append(rdata[idx:], labels[i]...)
		idx += len(labels[i])
	}
	
	return rdata
}

func NewAnswer(name string, typ string, ttl uint32, location string) Answer {
	var data [] byte
	atype, exist := recTypes[typ]
	if !exist {
		log.Fatalf("No numeric code for DNS type %s", typ)
	}
	
	switch atype {
		case 1:	// A record
			data = rdataFromIPv4(location)
		case 2, 5:	// NS, CNAME records
			data = rdataFromDomainName(location)
		default:
			log.Fatalf("DNS type %d is not supported.", atype)
	}
	
	return Answer{name: name, atype: atype, aclass: 1, ttl: ttl, rdata: data}
}

func NewResponse(header Header, question Question, answers, authorities, additional []Answer) Response {
	return Response{Header: header, Question: question, answers: answers, authorities: authorities, additional: additional}
}

func (query *Query) Deserialize(data []byte) {
	query.Header.id = binary.BigEndian.Uint16(data[0:2])
	query.Header.flags = binary.BigEndian.Uint16(data[2:4])
	
	query.Header.qdcount = binary.BigEndian.Uint16(data[4:6])
	query.Header.ancount = binary.BigEndian.Uint16(data[6:8])
	query.Header.nscount = binary.BigEndian.Uint16(data[8:10])
	query.Header.arcount = binary.BigEndian.Uint16(data[10:12])

	var index uint8 = 12
	var labelLength uint8 = uint8(data[index])
	for labelLength != 0 {
		index += 1
		query.Question.qname += string(data[index:index + labelLength])
		index += labelLength
		labelLength = uint8(data[index])
		if labelLength != 0 {
			query.Question.qname += "."
		}
	}
	
	index += 1	
	query.Question.qtype = binary.BigEndian.Uint16(data[index:index+2])
	index += 2
	query.Question.qclass = binary.BigEndian.Uint16(data[index:index+2])
}

func (response Response) Serialize() []byte {
	// TO DO: calculate the length of the response to preallocate the buffer
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, response.Header.id)
	if err != nil {
		fmt.Printf("Error when serialize ID of the response header: %s\n", err)
		return buff.Bytes()
	}
	
	err = binary.Write(buff, binary.BigEndian, response.Header.flags)
	if err != nil {
		fmt.Printf("Error when serialize flags of the response header: %s\n", err)
		return buff.Bytes()
	}
	
	err = binary.Write(buff, binary.BigEndian, response.Header.qdcount)
	if err != nil {
		fmt.Printf("Error when serialize qdcount of the response header: %s\n", err)
		return buff.Bytes()
	}
	
	binary.Write(buff, binary.BigEndian, response.Header.ancount)
	if err != nil {
		fmt.Printf("Error when serialize ancount of the response header: %s\n", err)
		return buff.Bytes()
	}
	
	binary.Write(buff, binary.BigEndian, response.Header.nscount)
	if err != nil {
		fmt.Printf("Error when serialize nscount of the response header: %s\n", err)
		return buff.Bytes()
	}
	
	binary.Write(buff, binary.BigEndian, response.Header.arcount)
	if err != nil {
		fmt.Printf("Error when serialize arcount of the response header: %s\n", err)
		return buff.Bytes()
	}	
	

	labels := strings.Split(response.Question.qname, ".")
	for _, lbl := range labels {
		buff.WriteByte(byte(len(lbl)))
		buff.WriteString(lbl)
	}
	
	buff.WriteByte(byte(0))
	
	binary.Write(buff, binary.BigEndian, response.Question.qtype)
	binary.Write(buff, binary.BigEndian, response.Question.qclass)
	
	fmt.Printf("\nResponse::Serialize - buff.Bytes(): %x\n", buff.Bytes())
	for _, b := range buff.Bytes() {
		fmt.Printf("%x ", b)
	}
		
	return buff.Bytes()	
}
