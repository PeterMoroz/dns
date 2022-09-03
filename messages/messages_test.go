package messages

import "testing"

func TestRdataFromIPv4(t *testing.T) {
	addresses := [6]string{ "127.0.0.1", "8.8.8.8", "8.8.4.4", "9.9.9.9", "149.112.112.112", "192.175.105.179" }
	expected := [6][4]byte{ { 127, 0, 0, 1 }, { 8, 8, 8, 8 }, { 8, 8, 4, 4 }, { 9, 9, 9, 9 }, { 149, 112, 112, 112 }, { 192, 175, 105, 179 } }
	
	for i, _ := range addresses {
		data := rdataFromIPv4(addresses[i])
		if len(data) != 4 {
			t.Errorf("IP address: %s, data length: %d\n", addresses[i], len(data))
		} else {
			isEqual := true
			for j, _ := range data {
				if data[j] != expected[i][j] {
					t.Errorf("IP address: %s. Mismatch at %d octet. \n", addresses[i], j)
					isEqual = false
				}
			}
			if isEqual {
				t.Logf("IP address: %s - OK\n", addresses[i])
			}
		}
	}
}

func TestRdataFromDomainName(t *testing.T) {
	dnames := [4]string{ "cloudns1.global-nameservers.com", "cloudns2.global-nameservers.com", "mail.global-nameservers.com", "cloudns1.alphost.com" }
	for i, _ := range dnames {
		data := rdataFromDomainName(dnames[i])
		if len(data) != len(dnames[i]) + 1 {
			t.Errorf("Domain name: %s (length %d), data length: %d\n", dnames[i], len(dnames[i]), len(data))
		} else {
			t.Logf("Domain name: %s - OK\n", dnames[i])
		}
	}
}
