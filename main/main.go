package main

import (
	"fmt"
	"dns/messages"
	"dns/resolver"
	"dns/server"
)

func main() {
	fmt.Println("DNS server")
	resolver.LoadFromFile("..\\dns-master.txt")
	
	var handler server.Handler = func(request []byte) []byte {
	/*
		header, question := messages.ParseQuery(request)
		
		var qname string = question.GetQname()
		
		fmt.Printf("DNS header {id: %x, flags: %x, qdcount: %d, ancount: %d, nscount: %d, arcount: %d}\nqueried name: %s\n",
			header.GetId(), header.GetFlags(), header.GetQDcount(), header.GetANcount(), header.GetNScount(), header.GetARcount(), qname)
			
		record := records.SearchByName(resourceRecords, qname)
		matchedResourceRecords := make([]*records.ResourceRecord, 0)
		
		if record != nil {
			matchedResourceRecords = append(matchedResourceRecords, record)
			for record.GetType() == "A" {
				record := records.SearchByName(resourceRecords, record.GetLocation())
				if record == nil {
					break
				}
				matchedResourceRecords = append(matchedResourceRecords, record)
			}
		}
		
		authorityResourceRecords := records.GetAllByType(resourceRecords, "NS")
		
		additionalResourceRecords := make([]*records.ResourceRecord, 0)
		if len(authorityResourceRecords) != 0 {
			rec := additionalResourceRecords[0]
			for rec != nil {
				additionalResourceRecords = append(additionalResourceRecords, rec)
				rec = records.SearchByName(resourceRecords, rec.GetLocation())
			}
		}
	*/
	
		var query messages.Query
		query.Deserialize(request)
		
		var header messages.Header
		var question messages.Question = query.Question		
		var response messages.Response
		
		var qname string = question.GetName()
		arecords := resolver.GetAddressRecords(qname)
		
		if len(arecords) == 0 {
			header.SetId(query.Header.GetId())
			header.SetQdcount(1)
			header.SetIsResponse(true)
			header.SetRcode(messages.Rcode_NameError)			
			
			response = messages.NewResponse(header, question, nil, nil, nil)
		} else {
			nsrecords := resolver.GetNameServerRecords()
			
			answers := make([]messages.Answer, len(arecords))
			for i, rec := range arecords {
				answers[i] = messages.NewAnswer(rec.GetName(), rec.GetType(), rec.GetTTL(), rec.GetLocation())
			}
			
			authorities := make([]messages.Answer, len(nsrecords))
			for i, rec := range nsrecords {
				authorities[i] = messages.NewAnswer(rec.GetName(), rec.GetType(), rec.GetTTL(), rec.GetLocation())
			}
			
			header.SetId(query.Header.GetId())
			header.SetQdcount(1)
			header.SetIsResponse(true)
			header.SetRcode(messages.Rcode_NoError)
			header.SetAA(true)
			header.SetAncount(uint16(len(arecords)))
			header.SetNscount(uint16(len(nsrecords)))
			
			response = messages.NewResponse(header, question, answers, authorities, nil)
		}

		
		data := response.Serialize()
		return data
	}
	
	server := server.NewServer(9000, 1024, handler)
	server.Run()
}