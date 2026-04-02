package recon

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type PortScanner struct {
	host      string
	ports     []int
	openPorts []int
	mu        sync.Mutex
	timeout   time.Duration
	workers   int
}

func NewPortScanner(host string) *PortScanner {
	return &PortScanner{
		host:      host,
		ports:     []int{},
		openPorts: []int{},
		timeout:   500 * time.Millisecond,
		workers:   100,
	}
}

func (ps *PortScanner) SetTimeout(timeout time.Duration) {
	ps.timeout = timeout
}

func (ps *PortScanner) SetWorkers(workers int) {
	ps.workers = workers
}

func (ps *PortScanner) AddPort(port int) {
	ps.ports = append(ps.ports, port)
}

func (ps *PortScanner) AddPortRange(start, end int) {
	for port := start; port <= end; port++ {
		ps.ports = append(ps.ports, port)
	}
}

func (ps *PortScanner) AddCommonPorts() {
	commonPorts := []int{
		21, 22, 23, 25, 53, 80, 110, 135, 139, 143,
		443, 445, 993, 995, 1433, 1521, 3306, 3389, 5432, 5900,
		6379, 8000, 8080, 8443, 8888, 9000, 9090, 27017, 27018, 27019,
	}
	for _, port := range commonPorts {
		ps.ports = append(ps.ports, port)
	}
}

func (ps *PortScanner) Scan() ([]int, error) {
	fmt.Printf("🔍 开始端口扫描: %s\n", ps.host)
	fmt.Printf("📊 扫描端口数量: %d\n", len(ps.ports))

	if len(ps.ports) == 0 {
		return []int{}, nil
	}

	portChan := make(chan int, len(ps.ports))
	resultChan := make(chan int, len(ps.ports))
	var wg sync.WaitGroup

	for i := 0; i < ps.workers; i++ {
		wg.Add(1)
		go ps.worker(portChan, resultChan, &wg)
	}

	for _, port := range ps.ports {
		portChan <- port
	}
	close(portChan)

	wg.Wait()
	close(resultChan)

	for port := range resultChan {
		ps.openPorts = append(ps.openPorts, port)
	}

	fmt.Printf("✅ 发现开放端口: %v\n", ps.openPorts)
	return ps.openPorts, nil
}

func (ps *PortScanner) worker(portChan <-chan int, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range portChan {
		if ps.scanPort(port) {
			ps.mu.Lock()
			resultChan <- port
			ps.mu.Unlock()
		}
	}
}

func (ps *PortScanner) scanPort(port int) bool {
	address := fmt.Sprintf("%s:%d", ps.host, port)
	conn, err := net.DialTimeout("tcp", address, ps.timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func (ps *PortScanner) GetServiceName(port int) string {
	serviceMap := map[int]string{
		21:    "FTP",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		53:    "DNS",
		80:    "HTTP",
		110:   "POP3",
		135:   "MSRPC",
		139:   "NetBIOS",
		143:   "IMAP",
		443:   "HTTPS",
		445:   "SMB",
		993:   "IMAPS",
		995:   "POP3S",
		1433:  "MSSQL",
		1521:  "Oracle",
		3306:  "MySQL",
		3389:  "RDP",
		5432:  "PostgreSQL",
		5900:  "VNC",
		6379:  "Redis",
		8000:  "HTTP-Proxy",
		8080:  "HTTP-Proxy",
		8443:  "HTTPS-Alt",
		8888:  "HTTP-Alt",
		9000:  "HTTP-Alt",
		9090:  "HTTP-Alt",
		27017: "MongoDB",
		27018: "MongoDB",
		27019: "MongoDB",
	}

	if service, ok := serviceMap[port]; ok {
		return service
	}
	return "Unknown"
}

func (ps *PortScanner) GetOpenPorts() []int {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	result := make([]int, len(ps.openPorts))
	copy(result, ps.openPorts)
	return result
}
