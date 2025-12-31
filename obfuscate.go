// Copyright 2020-present Kuei-chun Chen. All rights reserved.
// obfuscate.go

package gox

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"
)

// Pre-compiled regex patterns for PII detection
var (
	RePort   = regexp.MustCompile(`:\d{2,}`)
	ReEmail  = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)
	ReIP     = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	ReIPCIDR = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}(/\d{1,2})?$`)
	ReFQDN   = regexp.MustCompile(`([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z]{2,63}`)
	ReNS     = regexp.MustCompile(`[^@$.\n]*\.[^^@.\n]*([.][^^@.\n]*)?`)
	ReDigit  = regexp.MustCompile("[0-9.]")
	ReSSN    = regexp.MustCompile(`\d{3}-\d{2}-\d{4}`)
	ReMAC    = regexp.MustCompile(`([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2}`)
	ReDate   = regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)
	ReMRN    = regexp.MustCompile(`(?i)(mrn|acct|id)[:\s#]*\d{6,}`)
	RePhone  = regexp.MustCompile(`(\+\d{1,3}[-.\s]?)?(\(?\d{3}\)?[-.\s]?)?\d{3}[-.\s]?\d{4}`)
	ReCard   = regexp.MustCompile(`\d{4}[-\s]?\d{4}[-\s]?\d{4}[-\s]?\d{4}`)
)

// City and flower names for human-readable obfuscation
var (
	Cities = []string{
		"Atlanta", "Berlin", "Chicago", "Dublin", "ElPaso",
		"Foshan", "Giza", "Hongkong", "Istanbul", "Jakarta",
		"London", "Miami", "NewYork", "Orlando", "Paris",
		"Queens", "Rome", "Sydney", "Taipei", "Utica",
		"Vancouver", "Warsaw", "Xiamen", "Yonkers", "Zurich",
	}
	Flowers = []string{
		"Aster", "Begonia", "Carnation", "Daisy", "Erica",
		"Freesia", "Gardenia", "Hyacinth", "Iris", "Jasmine",
		"Kalmia", "Lavender", "Marigold", "Narcissus", "Orchid",
		"Peony", "Rose", "Sunflower", "Tulip", "Ursinia",
		"Violet", "Wisteria", "Xylobium", "Yarrow", "Zinnia",
	}
)

// IPStyle defines how IP addresses are obfuscated
type IPStyle int

const (
	// IPStyleKeepEnds keeps first and last octets: 192.168.1.100 → 192.X.X.100
	IPStyleKeepEnds IPStyle = iota
	// IPStylePrivate maps to 10.x.x.x range: 192.168.1.100 → 10.X.X.X
	IPStylePrivate
)

// NameStyle defines how names are obfuscated
type NameStyle int

const (
	// NameStyleReadable uses city/flower names for human readability
	NameStyleReadable NameStyle = iota
	// NameStyleHash uses hash-based prefixes (host-abc123)
	NameStyleHash
)

// Obfuscator handles PII obfuscation with consistent mappings
// Uses deterministic hashing so the same input always produces the same output
type Obfuscator struct {
	// Configuration
	Coefficient float64   // Multiplier for numeric obfuscation (default 0.917)
	DateOffset  int       // Days to shift dates (default -42)
	IPStyle     IPStyle   // How to obfuscate IPs
	NameStyle   NameStyle // How to obfuscate names

	// Mapping caches for consistency
	CardMap     map[string]string
	HostnameMap map[string]string
	IDMap       map[string]string
	IntMap      map[int]int
	IPMap       map[string]string
	MACMap      map[string]string
	NameMap     map[string]string
	NumberMap   map[string]float64
	PhoneMap    map[string]string
	ReplSetMap  map[string]string
	SSNMap      map[string]string
}

// NewObfuscator creates a new Obfuscator with default settings
func NewObfuscator() *Obfuscator {
	return &Obfuscator{
		Coefficient: 0.917,
		DateOffset:  -42,
		IPStyle:     IPStyleKeepEnds,
		NameStyle:   NameStyleReadable,
		CardMap:     make(map[string]string),
		HostnameMap: make(map[string]string),
		IDMap:       make(map[string]string),
		IntMap:      make(map[int]int),
		IPMap:       make(map[string]string),
		MACMap:      make(map[string]string),
		NameMap:     make(map[string]string),
		NumberMap:   make(map[string]float64),
		PhoneMap:    make(map[string]string),
		ReplSetMap:  make(map[string]string),
		SSNMap:      make(map[string]string),
	}
}

// --- Deterministic Hash Functions ---

// HashIndex returns a deterministic index (0 to max-1) based on the input string
// Uses FNV-1a hash for speed
func HashIndex(s string, max int) int {
	if max <= 0 {
		return 0
	}
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32()) % max
}

// HashOctet returns a deterministic octet (0-255) based on input and position
func HashOctet(s string, pos int) int {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%s:%d", s, pos)))
	return int(h.Sum32()) % 256
}

// HashString returns a deterministic hex string based on input
// Uses SHA-256 for cryptographic quality
func HashString(s string, length int) string {
	hash := sha256.Sum256([]byte(s))
	hex := hex.EncodeToString(hash[:])
	if length > 0 && length < len(hex) {
		return hex[:length]
	}
	return hex
}

// --- PII Detection Functions ---

// ContainsIP checks if string contains an IP address
func ContainsIP(s string) bool {
	return ReIP.MatchString(s)
}

// ContainsEmail checks if string contains an email address
func ContainsEmail(s string) bool {
	return ReEmail.MatchString(s)
}

// ContainsFQDN checks if string contains a fully qualified domain name
func ContainsFQDN(s string) bool {
	return ReFQDN.MatchString(s)
}

// ContainsSSN checks if string contains a Social Security Number
func ContainsSSN(s string) bool {
	return ReSSN.MatchString(s)
}

// ContainsMAC checks if string contains a MAC address
func ContainsMAC(s string) bool {
	return ReMAC.MatchString(s)
}

// ContainsPhoneNo checks if string contains a phone number
func ContainsPhoneNo(s string) bool {
	if !RePhone.MatchString(s) {
		return false
	}
	// Count digits - phone numbers have 10-15 digits
	digits := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			digits++
		}
	}
	return digits >= 10 && digits <= 15
}

// ContainsCreditCardNo checks if string contains a credit card number
func ContainsCreditCardNo(s string) bool {
	if !ReCard.MatchString(s) {
		return false
	}
	// Basic Luhn check could be added here
	digits := 0
	for _, c := range s {
		if c >= '0' && c <= '9' {
			digits++
		}
	}
	return digits >= 13 && digits <= 19
}

// IsNamespace checks if string looks like a MongoDB namespace (db.collection)
func IsNamespace(s string) bool {
	if strings.Contains(s, "/") || strings.Contains(s, "\\") {
		return false
	}
	parts := strings.Split(s, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return false
	}
	for _, part := range parts {
		if len(part) == 0 {
			return false
		}
	}
	return true
}

// LooksLikeHostname checks if string looks like a hostname
func LooksLikeHostname(s string) bool {
	if strings.Contains(s, " ") {
		return false
	}
	if strings.Contains(s, ".") {
		return true
	}
	if strings.Contains(s, "-") {
		return true
	}
	return false
}

// LooksLikeHostPort checks if string matches hostname:port pattern
func LooksLikeHostPort(s string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._-]+:\d+$`, s)
	return matched
}

// --- Core Obfuscation Methods ---

// ObfuscateIP obfuscates an IP address consistently
func (o *Obfuscator) ObfuscateIP(ip string) string {
	if !ContainsIP(ip) {
		return ip
	}

	matches := ReIP.FindStringSubmatch(ip)
	if len(matches) == 0 {
		return ip
	}

	matched := matches[0]
	if matched == "0.0.0.0" || matched == "127.0.0.1" {
		return ip
	}

	// Handle CIDR notation
	cidrSuffix := ""
	baseIP := matched
	if idx := strings.Index(ip, "/"); idx != -1 {
		cidrSuffix = ip[idx:]
	}

	if cached, exists := o.IPMap[baseIP]; exists {
		return strings.Replace(ip, matched, cached, -1) + cidrSuffix
	}

	var newIP string
	octets := strings.Split(baseIP, ".")
	if len(octets) != 4 {
		return ip
	}

	switch o.IPStyle {
	case IPStylePrivate:
		hash := sha256.Sum256([]byte(baseIP))
		newIP = fmt.Sprintf("10.%d.%d.%d", hash[0], hash[1], hash[2])
	case IPStyleKeepEnds:
		fallthrough
	default:
		newIP = octets[0] + "." + strconv.Itoa(HashOctet(baseIP, 1)) + "." +
			strconv.Itoa(HashOctet(baseIP, 2)) + "." + octets[3]
	}

	o.IPMap[baseIP] = newIP
	return strings.Replace(ip, matched, newIP, -1) + cidrSuffix
}

// ObfuscateHostname obfuscates a hostname consistently
func (o *Obfuscator) ObfuscateHostname(hostname string) string {
	if hostname == "" {
		return hostname
	}

	if cached, exists := o.HostnameMap[hostname]; exists {
		return cached
	}

	var obfuscated string
	switch o.NameStyle {
	case NameStyleHash:
		hash := HashString(hostname, 8)
		obfuscated = fmt.Sprintf("host-%s.local", hash)
	case NameStyleReadable:
		fallthrough
	default:
		city := Cities[HashIndex(hostname, len(Cities))]
		flower := Flowers[HashIndex(hostname+"flower", len(Flowers))]
		obfuscated = strings.ToLower(flower + "." + city + ".local")
	}

	o.HostnameMap[hostname] = obfuscated
	return obfuscated
}

// ObfuscateHostPort obfuscates hostname:port strings
func (o *Obfuscator) ObfuscateHostPort(value string) string {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return o.ObfuscateHostname(value)
	}

	host := parts[0]
	port := parts[1]

	if ContainsIP(host) {
		return o.ObfuscateIP(host) + ":" + port
	}

	return o.ObfuscateHostname(host) + ":" + port
}

// ObfuscateReplSet obfuscates a replica set name consistently
func (o *Obfuscator) ObfuscateReplSet(name string) string {
	if name == "" {
		return name
	}

	if cached, exists := o.ReplSetMap[name]; exists {
		return cached
	}

	var obfuscated string
	switch o.NameStyle {
	case NameStyleHash:
		hash := HashString(name, 8)
		obfuscated = fmt.Sprintf("rs-%s", hash)
	case NameStyleReadable:
		fallthrough
	default:
		city := Cities[HashIndex(name, len(Cities))]
		obfuscated = strings.ToLower("rs-" + city)
	}

	o.ReplSetMap[name] = obfuscated
	return obfuscated
}

// ObfuscateEmail obfuscates an email address consistently
func (o *Obfuscator) ObfuscateEmail(email string) string {
	if !ContainsEmail(email) {
		return email
	}

	matches := ReEmail.FindStringSubmatch(email)
	if len(matches) == 0 {
		return email
	}

	matched := matches[0]
	if cached, exists := o.NameMap[matched]; exists {
		return strings.Replace(email, matched, cached, -1)
	}

	city := Cities[HashIndex(matched, len(Cities))]
	flower := Flowers[HashIndex(matched+"flower", len(Flowers))]
	newValue := strings.ToLower(flower + "@" + city + ".com")

	o.NameMap[matched] = newValue
	o.NameMap[newValue] = newValue // Prevent re-obfuscation
	return strings.Replace(email, matched, newValue, -1)
}

// ObfuscateFQDN obfuscates a fully qualified domain name consistently
func (o *Obfuscator) ObfuscateFQDN(fqdn string) string {
	if strings.Contains(fqdn, "/") || strings.Contains(fqdn, "\\") {
		return fqdn
	}
	if !ContainsFQDN(fqdn) {
		return fqdn
	}

	matches := ReFQDN.FindStringSubmatch(fqdn)
	if len(matches) == 0 {
		return fqdn
	}

	matched := matches[0]
	if cached, exists := o.NameMap[matched]; exists {
		return strings.Replace(fqdn, matched, cached, -1)
	}

	newValue := o.generateObfuscatedName(matched)
	o.NameMap[matched] = newValue
	o.NameMap[newValue] = newValue
	return strings.Replace(fqdn, matched, newValue, -1)
}

// ObfuscateNamespace obfuscates a MongoDB namespace (db.collection)
func (o *Obfuscator) ObfuscateNamespace(ns string) string {
	if strings.Contains(ns, "/") || strings.Contains(ns, "\\") {
		return ns
	}
	if !IsNamespace(ns) {
		return ns
	}

	chars := ReDigit.ReplaceAllString(ns, "")
	if len(chars) == 0 {
		return ns
	}

	matches := ReNS.FindStringSubmatch(ns)
	if len(matches) == 0 {
		return ns
	}

	matched := matches[0]
	if cached, exists := o.NameMap[matched]; exists {
		return strings.Replace(ns, matched, cached, -1)
	}

	newValue := o.generateObfuscatedName(matched)
	o.NameMap[matched] = newValue
	o.NameMap[newValue] = newValue
	return strings.Replace(ns, matched, newValue, -1)
}

// ObfuscateSSN obfuscates a Social Security Number consistently
func (o *Obfuscator) ObfuscateSSN(ssn string) string {
	if !ContainsSSN(ssn) {
		return ssn
	}

	matches := ReSSN.FindStringSubmatch(ssn)
	if len(matches) == 0 {
		return ssn
	}

	matched := matches[0]
	if cached, exists := o.SSNMap[matched]; exists {
		return strings.Replace(ssn, matched, cached, -1)
	}

	digits := []byte{}
	for _, c := range matched {
		if c >= '0' && c <= '9' {
			digits = append(digits, byte(c))
		}
	}

	// Deterministic shuffle using hash
	for i := len(digits) - 1; i > 0; i-- {
		j := HashIndex(matched+strconv.Itoa(i), i+1)
		digits[i], digits[j] = digits[j], digits[i]
	}

	newValue := string(digits[:3]) + "-" + string(digits[3:5]) + "-" + string(digits[5:])
	o.SSNMap[matched] = newValue
	return strings.Replace(ssn, matched, newValue, -1)
}

// ObfuscateMAC obfuscates a MAC address consistently (keeps vendor prefix)
func (o *Obfuscator) ObfuscateMAC(value string) string {
	matches := ReMAC.FindStringSubmatch(value)
	if len(matches) == 0 {
		return value
	}

	matched := matches[0]
	if cached, exists := o.MACMap[matched]; exists {
		return strings.Replace(value, matched, cached, -1)
	}

	sep := ":"
	if strings.Contains(matched, "-") {
		sep = "-"
	}

	parts := strings.FieldsFunc(matched, func(r rune) bool { return r == ':' || r == '-' })
	if len(parts) != 6 {
		return value
	}

	// Keep vendor prefix (first 3 octets), obfuscate device ID (last 3)
	newParts := make([]string, 6)
	copy(newParts[:3], parts[:3])
	for i := 3; i < 6; i++ {
		newParts[i] = fmt.Sprintf("%02X", HashOctet(matched, i))
	}

	newValue := strings.Join(newParts, sep)
	o.MACMap[matched] = newValue
	return strings.Replace(value, matched, newValue, -1)
}

// ObfuscatePhoneNo obfuscates a phone number consistently
func (o *Obfuscator) ObfuscatePhoneNo(phoneNo string) string {
	if !ContainsPhoneNo(phoneNo) {
		return phoneNo
	}

	if cached, exists := o.PhoneMap[phoneNo]; exists {
		return cached
	}

	obfuscated := make([]byte, len(phoneNo))
	n := 0
	for i := range obfuscated {
		if phoneNo[i] >= '0' && phoneNo[i] <= '9' {
			n++
			if n > 5 {
				obfuscated[i] = byte(HashIndex(phoneNo+strconv.Itoa(i), 10) + '0')
			} else {
				obfuscated[i] = phoneNo[i]
			}
		} else {
			obfuscated[i] = phoneNo[i]
		}
	}

	o.PhoneMap[phoneNo] = string(obfuscated)
	return string(obfuscated)
}

// ObfuscateCreditCardNo obfuscates a credit card (masks all but last 4 digits)
func (o *Obfuscator) ObfuscateCreditCardNo(cardNo string) string {
	if !ContainsCreditCardNo(cardNo) {
		return cardNo
	}

	if cached, exists := o.CardMap[cardNo]; exists {
		return cached
	}

	lastFourDigits := cardNo[len(cardNo)-4:]
	obfuscated := make([]rune, len(cardNo)-4)
	for i, c := range cardNo[:len(cardNo)-4] {
		if c >= '0' && c <= '9' {
			obfuscated[i] = '*'
		} else {
			obfuscated[i] = c
		}
	}

	result := string(obfuscated) + lastFourDigits
	o.CardMap[cardNo] = result
	return result
}

// ObfuscateDate shifts dates by DateOffset days
func (o *Obfuscator) ObfuscateDate(value string) string {
	matches := ReDate.FindAllString(value, -1)
	if len(matches) == 0 {
		return value
	}

	for _, matched := range matches {
		year, _ := strconv.Atoi(matched[0:4])
		month, _ := strconv.Atoi(matched[5:7])
		day, _ := strconv.Atoi(matched[8:10])

		day += o.DateOffset
		for day < 1 {
			month--
			if month < 1 {
				month = 12
				year--
			}
			day += 30
		}
		for day > 28 {
			day -= 28
			month++
			if month > 12 {
				month = 1
				year++
			}
		}

		newValue := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
		value = strings.Replace(value, matched, newValue, 1)
	}
	return value
}

// ObfuscateInt obfuscates an integer using the coefficient
func (o *Obfuscator) ObfuscateInt(value int) int {
	if value <= 1 {
		return value
	}
	if cached, exists := o.IntMap[value]; exists {
		return cached
	}
	newValue := int(float64(value) * o.Coefficient)
	o.IntMap[value] = newValue
	return newValue
}

// ObfuscateNumber obfuscates a float using the coefficient
func (o *Obfuscator) ObfuscateNumber(value float64) float64 {
	key := fmt.Sprintf("%f", value)
	if cached, exists := o.NumberMap[key]; exists {
		return cached
	}
	newValue := value * o.Coefficient
	o.NumberMap[key] = newValue
	return newValue
}

// --- Generic Traversal Methods ---

// ObfuscateMap recursively obfuscates a map[string]interface{}
func (o *Obfuscator) ObfuscateMap(doc map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(doc))
	for k, v := range doc {
		result[k] = o.ObfuscateValue(v)
	}
	return result
}

// ObfuscateSlice recursively obfuscates a []interface{}
func (o *Obfuscator) ObfuscateSlice(arr []interface{}) []interface{} {
	result := make([]interface{}, len(arr))
	for i, elem := range arr {
		result[i] = o.ObfuscateValue(elem)
	}
	return result
}

// ObfuscateValue obfuscates a value based on its type
func (o *Obfuscator) ObfuscateValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		return o.ObfuscateMap(v)
	case []interface{}:
		return o.ObfuscateSlice(v)
	case string:
		return o.ObfuscateString(v)
	case int:
		return o.ObfuscateInt(v)
	case int32:
		return int32(o.ObfuscateInt(int(v)))
	case int64:
		return int64(o.ObfuscateInt(int(v)))
	case float32:
		return float32(o.ObfuscateNumber(float64(v)))
	case float64:
		return o.ObfuscateNumber(v)
	default:
		return value
	}
}

// ObfuscateString applies all string obfuscation rules
func (o *Obfuscator) ObfuscateString(value string) string {
	// Port numbers
	if matches := RePort.FindStringSubmatch(value); len(matches) > 0 {
		matched := matches[0]
		port := ToInt(matched[1:])
		newValue := fmt.Sprintf(":%v", int(float64(port)*o.Coefficient))
		value = strings.Replace(value, matched, newValue, -1)
	}

	// Credit cards
	if ContainsCreditCardNo(value) {
		value = o.ObfuscateCreditCardNo(value)
	}

	// Order matters for these
	value = o.ObfuscateEmail(value)
	value = o.ObfuscateNamespace(value)
	value = o.ObfuscateFQDN(value)
	value = o.ObfuscateIP(value)
	value = o.ObfuscateMAC(value)
	value = o.ObfuscateSSN(value)
	value = o.ObfuscatePhoneNo(value)
	value = o.ObfuscateDate(value)

	return value
}

// --- Utility Methods ---

// generateObfuscatedName generates an obfuscated name from city and flower
func (o *Obfuscator) generateObfuscatedName(matched string) string {
	city := Cities[HashIndex(matched, len(Cities))]
	flower := Flowers[HashIndex(matched+"flower", len(Flowers))]
	parts := strings.Split(matched, ".")
	if len(parts) > 2 {
		tail := parts[len(parts)-1]
		return strings.ToLower(flower + "." + city + "." + tail)
	}
	return strings.ToLower(city + "." + flower)
}

// GetMappings returns all obfuscation mappings (for debugging/reference)
func (o *Obfuscator) GetMappings() map[string]interface{} {
	// Filter out self-mappings from NameMap
	filteredNameMap := make(map[string]string)
	for k, v := range o.NameMap {
		if k != v {
			filteredNameMap[k] = v
		}
	}

	return map[string]interface{}{
		"coefficient":  o.Coefficient,
		"date_offset":  o.DateOffset,
		"card_map":     o.CardMap,
		"hostname_map": o.HostnameMap,
		"id_map":       o.IDMap,
		"ip_map":       o.IPMap,
		"mac_map":      o.MACMap,
		"name_map":     filteredNameMap,
		"phone_map":    o.PhoneMap,
		"replset_map":  o.ReplSetMap,
		"ssn_map":      o.SSNMap,
	}
}

// Reset clears all obfuscation mappings
func (o *Obfuscator) Reset() {
	o.CardMap = make(map[string]string)
	o.HostnameMap = make(map[string]string)
	o.IDMap = make(map[string]string)
	o.IntMap = make(map[int]int)
	o.IPMap = make(map[string]string)
	o.MACMap = make(map[string]string)
	o.NameMap = make(map[string]string)
	o.NumberMap = make(map[string]float64)
	o.PhoneMap = make(map[string]string)
	o.ReplSetMap = make(map[string]string)
	o.SSNMap = make(map[string]string)
}

