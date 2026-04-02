package portscanner

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

type PortResult struct {
	Port     int    `json:"port"`
	Status   string `json:"status"`
	Service  string `json:"service,omitempty"`
	Banner   string `json:"banner,omitempty"`
}

type ScannerConfig struct {
	Timeout     time.Duration
	Concurrency int
}

type Scanner struct {
	config ScannerConfig
}

func NewScanner(config ScannerConfig) *Scanner {
	if config.Timeout == 0 {
		config.Timeout = 2 * time.Second
	}
	if config.Concurrency <= 0 {
		config.Concurrency = 50
	}
	return &Scanner{config: config}
}

func (s *Scanner) ScanPort(ip string, port int) *PortResult {
	target := net.JoinHostPort(ip, strconv.Itoa(port))
	
	conn, err := net.DialTimeout("tcp", target, s.config.Timeout)
	if err != nil {
		return &PortResult{
			Port:   port,
			Status: "closed",
		}
	}
	defer conn.Close()

	result := &PortResult{
		Port:   port,
		Status: "open",
	}

	result.Service = getCommonService(port)
	result.Banner = s.getBanner(conn)

	return result
}

func (s *Scanner) ScanPorts(ip string, ports []int) []PortResult {
	var wg sync.WaitGroup
	resultsChan := make(chan *PortResult, len(ports))
	sem := make(chan struct{}, s.config.Concurrency)

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			sem <- struct{}{}
			resultsChan <- s.ScanPort(ip, p)
			<-sem
		}(port)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var results []PortResult
	for result := range resultsChan {
		results = append(results, *result)
	}

	return results
}

func (s *Scanner) ScanPortRange(ip string, startPort, endPort int) []PortResult {
	var ports []int
	for port := startPort; port <= endPort; port++ {
		ports = append(ports, port)
	}
	return s.ScanPorts(ip, ports)
}

func (s *Scanner) getBanner(conn net.Conn) string {
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return ""
	}
	
	return string(buffer[:n])
}

func getCommonService(port int) string {
	services := map[int]string{
		21:    "ftp",
		22:    "ssh",
		23:    "telnet",
		25:    "smtp",
		53:    "dns",
		80:    "http",
		110:   "pop3",
		143:   "imap",
		443:   "https",
		465:   "smtps",
		587:   "submission",
		993:   "imaps",
		995:   "pop3s",
		1433:  "mssql",
		1521:  "oracle",
		3306:  "mysql",
		3389:  "rdp",
		5432:  "postgresql",
		5900:  "vnc",
		6379:  "redis",
		8080:  "http-proxy",
		8443:  "https-alt",
		8888:  "sun-answerbook",
		27017: "mongodb",
		27018: "mongodb",
		27019: "mongodb",
	}

	if service, ok := services[port]; ok {
		return service
	}
	return "unknown"
}

func CommonPorts() []int {
	return []int{
		21, 22, 23, 25, 53, 80, 110, 143, 443, 465, 587, 993, 995,
		1433, 1521, 3306, 3389, 5432, 5900, 6379,
		8080, 8443, 8888, 27017, 27018, 27019,
	}
}

func FullPorts() []int {
	var ports []int
	for port := 1; port <= 65535; port++ {
		ports = append(ports, port)
	}
	return ports
}

func ParsePorts(portStr string) ([]int, error) {
	var ports []int
	
	parts := splitPortString(portStr)
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		if isRange(part) {
			start, end, err := parseRange(part)
			if err != nil {
				return nil, err
			}
			for port := start; port <= end; port++ {
				ports = append(ports, port)
			}
		} else {
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("port out of range: %d", port)
			}
			ports = append(ports, port)
		}
	}
	
	return uniqueInts(ports), nil
}

func splitPortString(s string) []string {
	var result []string
	var current []rune
	
	for _, r := range s {
		if r == ',' || r == ' ' || r == '\t' || r == '\n' {
			if len(current) > 0 {
				result = append(result, string(current))
				current = []rune{}
			}
		} else {
			current = append(current, r)
		}
	}
	
	if len(current) > 0 {
		result = append(result, string(current))
	}
	
	return result
}

func isRange(s string) bool {
	for _, r := range s {
		if r == '-' {
			return true
		}
	}
	return false
}

func parseRange(s string) (int, int, error) {
	parts := splitRange(s)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid range format: %s", s)
	}
	
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start port: %s", parts[0])
	}
	
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end port: %s", parts[1])
	}
	
	if start < 1 || start > 65535 {
		return 0, 0, fmt.Errorf("start port out of range: %d", start)
	}
	
	if end < 1 || end > 65535 {
		return 0, 0, fmt.Errorf("end port out of range: %d", end)
	}
	
	if start > end {
		return 0, 0, fmt.Errorf("start port greater than end port: %d > %d", start, end)
	}
	
	return start, end, nil
}

func splitRange(s string) []string {
	var result []string
	var current []rune
	
	for _, r := range s {
		if r == '-' {
			if len(current) > 0 {
				result = append(result, string(current))
				current = []rune{}
			}
		} else {
			current = append(current, r)
		}
	}
	
	if len(current) > 0 {
		result = append(result, string(current))
	}
	
	return result
}

func uniqueInts(ints []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	
	for _, i := range ints {
		if !seen[i] {
			seen[i] = true
			result = append(result, i)
		}
	}
	
	return result
}
