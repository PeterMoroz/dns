package resolver

import "testing"
import "fmt"

func TestSearchByName(t *testing.T) {
	fmt.Println(" - load records from file")
	LoadFromFile("..\\dns-master.txt")

	fmt.Printf("load %d records\n", len(resourceRecords))
	if len(resourceRecords) == 0 {
		t.Errorf("Could not load records from file")
	} else {
		t.Logf("load %d records", len(resourceRecords))
	}
	
	fmt.Println(" - search record by name")
	record := searchByName("ns1.student.test")
	
	if record == nil {
		t.Error("Could not find record by name")
	} else {
		t.Logf("The location: %s", record.GetLocation())
		// fmt.Printf("The location: %s\n", record.GetLocation())
	}
}