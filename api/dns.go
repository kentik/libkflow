package api

//go:generate msgp

type DNSQuestion struct {
	Name string
	Host []byte
}

type DNSResourceRecord struct {
	Name  string // unused
	CNAME string
	IP    []byte
	TTL   uint32
}

type DNSResponse struct {
	Question DNSQuestion
	Answers  []DNSResourceRecord
}
