package resolver

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"strconv"
)

var recordAuthDomain string
var recordGlobalTTL uint32 = math.MaxUint32
var resourceRecords []ResourceRecord = make([]ResourceRecord, 2)

type ResourceRecord struct {
	name string
	typ string
	ttl uint32
	location string
}

func (record *ResourceRecord) GetName() string {
	return record.name
}

func (record *ResourceRecord) GetType() string {
	return record.typ
}

func (record *ResourceRecord) GetTTL() uint32 {
	return record.ttl
}

func (record *ResourceRecord) GetLocation() string {
	return record.location
}


func printRecord(record *ResourceRecord) {
	fmt.Printf("Name: %s\nType: %s\nTTL: %d\nLocation: %s\n\n",
		record.name, record.typ, record.ttl, record.location);
}

func PrintRecords(records []ResourceRecord) {
	for _, r := range records {
		printRecord(&r)
	}
}

func LoadFromFile(filename string) {
	
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Could not open file %s - %s", filename, err)
	}
	
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	
	// resourceRecords = make([]ResourceRecord, 2)
	
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, "\n")
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		
		// assume that the first record is Auth Domain, 
		// the second one is global TTL followed by other records
		if len(recordAuthDomain) == 0 {
			recordAuthDomain = line
		} else if (recordGlobalTTL == math.MaxUint32) {
			ttl, err := strconv.Atoi(line)
			if err != nil {
				log.Fatalf("Could not parse global TTL record (%s) - %s", line, err)
			}
			recordGlobalTTL = uint32(ttl)
		} else {
			fields := strings.Fields(line)
			if len(fields) != 3 {
				log.Fatalf("Invalid record '%s'", line)
			}
			record := ResourceRecord{name: fields[0], typ: fields[1], location: fields[2], ttl: recordGlobalTTL}
			resourceRecords = append(resourceRecords, record)
		}
	}
}

func searchByName(name string) *ResourceRecord {
	for i, _ := range resourceRecords {
		if resourceRecords[i].name == name {
			return &resourceRecords[i]
		}
	}
	return nil
}

func GetAddressRecords(qname string) []ResourceRecord {
	records := make([]ResourceRecord, 0)
	record := searchByName(qname)
	if record != nil && record.GetType() == "A" {
		records = append(records, *record)
		record = searchByName(record.GetLocation())
	}
	return records
}

func GetNameServerRecords() []ResourceRecord {
	records := make([]ResourceRecord, 0)
	for _, record := range resourceRecords {
		if record.GetType() == "NS" {
			records = append(records, record)
		}
	}
	return records
}
