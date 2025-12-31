// Copyright 2020-present Kuei-chun Chen. All rights reserved.
// obfuscate_test.go

package gox

import (
	"testing"
)

func TestHashIndex(t *testing.T) {
	// Test determinism - same input always produces same output
	idx1 := HashIndex("test", 100)
	idx2 := HashIndex("test", 100)
	if idx1 != idx2 {
		t.Errorf("HashIndex not deterministic: got %d and %d", idx1, idx2)
	}

	// Test range
	for i := 0; i < 100; i++ {
		idx := HashIndex("test"+string(rune(i)), 10)
		if idx < 0 || idx >= 10 {
			t.Errorf("HashIndex out of range: got %d, expected 0-9", idx)
		}
	}

	// Test edge case
	if HashIndex("test", 0) != 0 {
		t.Error("HashIndex(s, 0) should return 0")
	}
}

func TestHashOctet(t *testing.T) {
	// Test determinism
	oct1 := HashOctet("192.168.1.1", 1)
	oct2 := HashOctet("192.168.1.1", 1)
	if oct1 != oct2 {
		t.Errorf("HashOctet not deterministic: got %d and %d", oct1, oct2)
	}

	// Test range (0-255)
	for i := 0; i < 100; i++ {
		oct := HashOctet("test", i)
		if oct < 0 || oct > 255 {
			t.Errorf("HashOctet out of range: got %d, expected 0-255", oct)
		}
	}
}

func TestHashString(t *testing.T) {
	// Test determinism
	hash1 := HashString("test", 8)
	hash2 := HashString("test", 8)
	if hash1 != hash2 {
		t.Errorf("HashString not deterministic: got %s and %s", hash1, hash2)
	}

	// Test length
	hash := HashString("test", 16)
	if len(hash) != 16 {
		t.Errorf("HashString wrong length: got %d, expected 16", len(hash))
	}
}

func TestContainsIP(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"host:192.168.1.1:27017", true},
		{"not an ip", false},
		{"192.168.1", false},
		{"192.168.1.1.1", true}, // Contains valid IP pattern 192.168.1.1
	}

	for _, tc := range tests {
		result := ContainsIP(tc.input)
		if result != tc.expected {
			t.Errorf("ContainsIP(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestContainsEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"user@example.com", true},
		{"user.name@domain.org", true},
		{"user+tag@example.co.uk", true},
		{"not an email", false},
		{"@example.com", false},
		{"user@", false},
	}

	for _, tc := range tests {
		result := ContainsEmail(tc.input)
		if result != tc.expected {
			t.Errorf("ContainsEmail(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestContainsSSN(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"123-45-6789", true},
		{"SSN: 123-45-6789", true},
		{"12345-6789", false},
		{"123456789", false},
		{"not a ssn", false},
	}

	for _, tc := range tests {
		result := ContainsSSN(tc.input)
		if result != tc.expected {
			t.Errorf("ContainsSSN(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestContainsPhoneNo(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"555-123-4567", true},
		{"(555) 123-4567", true},
		{"+1-555-123-4567", true},
		{"5551234567", true},
		{"555-1234", false},
		{"123", false},
	}

	for _, tc := range tests {
		result := ContainsPhoneNo(tc.input)
		if result != tc.expected {
			t.Errorf("ContainsPhoneNo(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestIsNamespace(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"mydb.mycollection", true},
		{"mydb.mycollection.index", true},
		{"admin.system.version", true},
		{"mydb", false},
		{"/path/to/file", false},
		{"", false},
		{".collection", false},
		{"db.", false},
	}

	for _, tc := range tests {
		result := IsNamespace(tc.input)
		if result != tc.expected {
			t.Errorf("IsNamespace(%q) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestObfuscateIP(t *testing.T) {
	o := NewObfuscator()

	// Test determinism
	ip1 := o.ObfuscateIP("192.168.1.100")
	ip2 := o.ObfuscateIP("192.168.1.100")
	if ip1 != ip2 {
		t.Errorf("ObfuscateIP not deterministic: got %s and %s", ip1, ip2)
	}

	// Test that localhost/zero are preserved
	if o.ObfuscateIP("127.0.0.1") != "127.0.0.1" {
		t.Error("127.0.0.1 should not be obfuscated")
	}
	if o.ObfuscateIP("0.0.0.0") != "0.0.0.0" {
		t.Error("0.0.0.0 should not be obfuscated")
	}

	// Test IPStyleKeepEnds (default)
	obfuscated := o.ObfuscateIP("192.168.1.100")
	if obfuscated == "192.168.1.100" {
		t.Error("IP should be obfuscated")
	}
	// First and last octets should be preserved
	if obfuscated[:4] != "192." {
		t.Errorf("First octet should be preserved, got %s", obfuscated)
	}
	if obfuscated[len(obfuscated)-3:] != "100" {
		t.Errorf("Last octet should be preserved, got %s", obfuscated)
	}

	// Test IPStylePrivate
	o2 := NewObfuscator()
	o2.IPStyle = IPStylePrivate
	obfuscated2 := o2.ObfuscateIP("192.168.1.100")
	if obfuscated2[:3] != "10." {
		t.Errorf("IPStylePrivate should produce 10.x.x.x, got %s", obfuscated2)
	}
}

func TestObfuscateHostname(t *testing.T) {
	o := NewObfuscator()

	// Test determinism
	h1 := o.ObfuscateHostname("server1.example.com")
	h2 := o.ObfuscateHostname("server1.example.com")
	if h1 != h2 {
		t.Errorf("ObfuscateHostname not deterministic: got %s and %s", h1, h2)
	}

	// Test NameStyleReadable (default)
	obfuscated := o.ObfuscateHostname("server1.example.com")
	if obfuscated == "server1.example.com" {
		t.Error("Hostname should be obfuscated")
	}
	// Should end with .local and contain flower/city names
	if obfuscated[len(obfuscated)-6:] != ".local" {
		t.Errorf("NameStyleReadable should produce .local suffix, got %s", obfuscated)
	}

	// Test NameStyleHash
	o2 := NewObfuscator()
	o2.NameStyle = NameStyleHash
	obfuscated2 := o2.ObfuscateHostname("server1.example.com")
	if obfuscated2[:5] != "host-" {
		t.Errorf("NameStyleHash should produce host- prefix, got %s", obfuscated2)
	}
}

func TestObfuscateHostPort(t *testing.T) {
	o := NewObfuscator()

	// Test hostname:port
	result := o.ObfuscateHostPort("server1.example.com:27017")
	if result == "server1.example.com:27017" {
		t.Error("HostPort should be obfuscated")
	}
	// Port should be preserved
	if result[len(result)-6:] != ":27017" {
		t.Errorf("Port should be preserved, got %s", result)
	}

	// Test IP:port
	result2 := o.ObfuscateHostPort("192.168.1.1:27017")
	if result2 == "192.168.1.1:27017" {
		t.Error("IP:Port should be obfuscated")
	}
	if result2[len(result2)-6:] != ":27017" {
		t.Errorf("Port should be preserved, got %s", result2)
	}
}

func TestObfuscateEmail(t *testing.T) {
	o := NewObfuscator()

	// Test determinism
	e1 := o.ObfuscateEmail("user@example.com")
	e2 := o.ObfuscateEmail("user@example.com")
	if e1 != e2 {
		t.Errorf("ObfuscateEmail not deterministic: got %s and %s", e1, e2)
	}

	// Test format
	obfuscated := o.ObfuscateEmail("user@example.com")
	if obfuscated == "user@example.com" {
		t.Error("Email should be obfuscated")
	}
	// Should contain @ and .com
	if !ContainsEmail(obfuscated) {
		t.Errorf("Obfuscated email should still look like email, got %s", obfuscated)
	}
}

func TestObfuscateSSN(t *testing.T) {
	o := NewObfuscator()

	// Test determinism
	s1 := o.ObfuscateSSN("123-45-6789")
	s2 := o.ObfuscateSSN("123-45-6789")
	if s1 != s2 {
		t.Errorf("ObfuscateSSN not deterministic: got %s and %s", s1, s2)
	}

	// Test format preservation
	obfuscated := o.ObfuscateSSN("123-45-6789")
	if obfuscated == "123-45-6789" {
		t.Error("SSN should be obfuscated")
	}
	// Should maintain XXX-XX-XXXX format
	if len(obfuscated) != 11 || obfuscated[3] != '-' || obfuscated[6] != '-' {
		t.Errorf("SSN format not preserved, got %s", obfuscated)
	}
}

func TestObfuscatePhoneNo(t *testing.T) {
	o := NewObfuscator()

	// Test determinism
	p1 := o.ObfuscatePhoneNo("555-123-4567")
	p2 := o.ObfuscatePhoneNo("555-123-4567")
	if p1 != p2 {
		t.Errorf("ObfuscatePhoneNo not deterministic: got %s and %s", p1, p2)
	}

	// Test that first 5 digits are preserved
	obfuscated := o.ObfuscatePhoneNo("555-123-4567")
	// First 5 digits should be preserved (555-1)
	if obfuscated[:5] != "555-1" {
		t.Errorf("First 5 digits should be preserved, got %s", obfuscated)
	}
}

func TestObfuscateInt(t *testing.T) {
	o := NewObfuscator()

	// Test small values preserved
	if o.ObfuscateInt(0) != 0 {
		t.Error("0 should not be obfuscated")
	}
	if o.ObfuscateInt(1) != 1 {
		t.Error("1 should not be obfuscated")
	}

	// Test larger values
	result := o.ObfuscateInt(100)
	expected := int(100 * o.Coefficient)
	if result != expected {
		t.Errorf("ObfuscateInt(100) = %d, expected %d", result, expected)
	}

	// Test determinism
	r1 := o.ObfuscateInt(1000)
	r2 := o.ObfuscateInt(1000)
	if r1 != r2 {
		t.Errorf("ObfuscateInt not deterministic: got %d and %d", r1, r2)
	}
}

func TestObfuscateDate(t *testing.T) {
	o := NewObfuscator()

	// Test date shifting
	result := o.ObfuscateDate("2024-06-15")
	if result == "2024-06-15" {
		t.Error("Date should be shifted")
	}

	// Test format preservation
	if len(result) != 10 || result[4] != '-' || result[7] != '-' {
		t.Errorf("Date format not preserved, got %s", result)
	}
}

func TestObfuscateMap(t *testing.T) {
	o := NewObfuscator()

	doc := map[string]interface{}{
		"ip":    "192.168.1.1",
		"email": "user@example.com",
		"count": 100,
		"nested": map[string]interface{}{
			"host": "server1.example.com",
		},
	}

	result := o.ObfuscateMap(doc)

	// Check that values are obfuscated
	if result["ip"] == "192.168.1.1" {
		t.Error("IP should be obfuscated in Map")
	}
	if result["email"] == "user@example.com" {
		t.Error("Email should be obfuscated in Map")
	}
	if result["count"] == 100 {
		t.Error("Count should be obfuscated in Map")
	}

	nested, ok := result["nested"].(map[string]interface{})
	if !ok {
		t.Error("Nested should be map[string]interface{}")
	}
	if nested["host"] == "server1.example.com" {
		t.Error("Nested host should be obfuscated")
	}
}

func TestObfuscateSlice(t *testing.T) {
	o := NewObfuscator()

	arr := []interface{}{
		"192.168.1.1",
		"user@example.com",
		100,
	}

	result := o.ObfuscateSlice(arr)

	if result[0] == "192.168.1.1" {
		t.Error("IP should be obfuscated in Slice")
	}
	if result[1] == "user@example.com" {
		t.Error("Email should be obfuscated in Slice")
	}
	if result[2] == 100 {
		t.Error("Number should be obfuscated in Slice")
	}
}

func TestGetMappings(t *testing.T) {
	o := NewObfuscator()

	// Generate some mappings
	o.ObfuscateIP("192.168.1.1")
	o.ObfuscateHostname("server1.example.com")
	o.ObfuscateEmail("user@example.com")

	mappings := o.GetMappings()

	// Check that maps are present
	if _, ok := mappings["ip_map"]; !ok {
		t.Error("Mappings should contain ip_map")
	}
	if _, ok := mappings["hostname_map"]; !ok {
		t.Error("Mappings should contain hostname_map")
	}
	if _, ok := mappings["name_map"]; !ok {
		t.Error("Mappings should contain name_map")
	}
}

func TestReset(t *testing.T) {
	o := NewObfuscator()

	// Generate some mappings
	o.ObfuscateIP("192.168.1.1")
	o.ObfuscateHostname("server1.example.com")

	// Verify mappings exist
	if len(o.IPMap) == 0 {
		t.Error("IPMap should have entries before reset")
	}

	// Reset
	o.Reset()

	// Verify mappings cleared
	if len(o.IPMap) != 0 {
		t.Error("IPMap should be empty after reset")
	}
	if len(o.HostnameMap) != 0 {
		t.Error("HostnameMap should be empty after reset")
	}
}

