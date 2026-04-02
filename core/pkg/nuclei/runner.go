package nuclei

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type NucleiResult struct {
	TemplateID       string   `json:"template-id"`
	TemplatePath     string   `json:"template-path,omitempty"`
	TemplateURL      string   `json:"template-url,omitempty"`
	Info             Info     `json:"info"`
	Type             string   `json:"type"`
	Host             string   `json:"host"`
	MatchedAt        string   `json:"matched-at"`
	ExtractedResults []string `json:"extracted-results,omitempty"`
	Request          string   `json:"request,omitempty"`
	Response         string   `json:"response,omitempty"`
	Timestamp        string   `json:"timestamp,omitempty"`
	CurlCommand      string   `json:"curl-command,omitempty"`
	MatcherName      string   `json:"matcher-name,omitempty"`
	MatcherStatus    bool     `json:"matcher-status,omitempty"`
}

type Info struct {
	Name           string   `json:"name"`
	Author         string   `json:"author,omitempty"`
	Severity       string   `json:"severity"`
	Description    string   `json:"description,omitempty"`
	Reference      []string `json:"reference,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	Classification *Class   `json:"classification,omitempty"`
	Metadata       Metadata `json:"metadata,omitempty"`
}

type Class struct {
	CveID     string `json:"cve-id,omitempty"`
	CweID     string `json:"cwe-id,omitempty"`
	CvssScore string `json:"cvss-score,omitempty"`
}

type Metadata struct {
	Product  string `json:"product,omitempty"`
	Vendor   string `json:"vendor,omitempty"`
	Verified bool   `json:"verified,omitempty"`
}

type RunnerConfig struct {
	TemplatePaths  []string
	Tags           []string
	Severity       []string
	ExcludeTags    []string
	ExcludeSeverity []string
	Concurrency    int
	RateLimit      int
	Timeout        int
	Retries        int
	Verbose        bool
	Silent         bool
	JSON           bool
	OutputFile     string
}

type Runner struct {
	config RunnerConfig
	mu     sync.Mutex
}

func NewRunner(config RunnerConfig) *Runner {
	if config.Concurrency <= 0 {
		config.Concurrency = 25
	}
	if config.RateLimit <= 0 {
		config.RateLimit = 150
	}
	if config.Timeout <= 0 {
		config.Timeout = 10
	}
	if config.Retries <= 0 {
		config.Retries = 3
	}
	return &Runner{config: config}
}

func (r *Runner) Run(target string) ([]NucleiResult, error) {
	args := r.buildArgs(target)
	
	cmd := exec.Command("nuclei", args...)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("nuclei execution failed: %w, output: %s", err, string(output))
	}

	return r.parseOutput(string(output))
}

func (r *Runner) RunAsync(target string, resultChan chan<- []NucleiResult, errChan chan<- error) {
	go func() {
		results, err := r.Run(target)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- results
	}()
}

func (r *Runner) buildArgs(target string) []string {
	args := []string{"-target", target, "-json"}

	if r.config.Silent {
		args = append(args, "-silent")
	} else if r.config.Verbose {
		args = append(args, "-v")
	}

	if len(r.config.TemplatePaths) > 0 {
		for _, path := range r.config.TemplatePaths {
			args = append(args, "-t", path)
		}
	}

	if len(r.config.Tags) > 0 {
		args = append(args, "-tags", strings.Join(r.config.Tags, ","))
	}

	if len(r.config.ExcludeTags) > 0 {
		args = append(args, "-etags", strings.Join(r.config.ExcludeTags, ","))
	}

	if len(r.config.Severity) > 0 {
		args = append(args, "-severity", strings.Join(r.config.Severity, ","))
	}

	if len(r.config.ExcludeSeverity) > 0 {
		args = append(args, "-eseverity", strings.Join(r.config.ExcludeSeverity, ","))
	}

	if r.config.Concurrency > 0 {
		args = append(args, "-c", fmt.Sprintf("%d", r.config.Concurrency))
	}

	if r.config.RateLimit > 0 {
		args = append(args, "-rl", fmt.Sprintf("%d", r.config.RateLimit))
	}

	if r.config.Timeout > 0 {
		args = append(args, "-timeout", fmt.Sprintf("%d", r.config.Timeout))
	}

	if r.config.Retries > 0 {
		args = append(args, "-retries", fmt.Sprintf("%d", r.config.Retries))
	}

	if r.config.OutputFile != "" {
		args = append(args, "-o", r.config.OutputFile)
	}

	args = append(args, "-no-color")

	return args
}

func (r *Runner) parseOutput(output string) ([]NucleiResult, error) {
	var results []NucleiResult
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var result NucleiResult
		if err := json.Unmarshal([]byte(line), &result); err == nil {
			if result.TemplateID != "" {
				results = append(results, result)
			}
		}
	}

	return results, nil
}

func (r *Runner) CheckNucleiAvailable() bool {
	cmd := exec.Command("nuclei", "-version")
	_, err := cmd.CombinedOutput()
	return err == nil
}

func (r *Runner) GetNucleiVersion() (string, error) {
	cmd := exec.Command("nuclei", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get nuclei version: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func (r *Runner) UpdateTemplates() error {
	cmd := exec.Command("nuclei", "-update-templates")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to update templates: %w", err)
	}
	return nil
}

type ScanStats struct {
	TotalResults   int
	SeverityCounts map[string]int
	TemplateCounts map[string]int
	ScanDuration   time.Duration
}

func CalculateStats(results []NucleiResult, duration time.Duration) ScanStats {
	stats := ScanStats{
		TotalResults:   len(results),
		SeverityCounts: make(map[string]int),
		TemplateCounts: make(map[string]int),
		ScanDuration:   duration,
	}

	for _, result := range results {
		stats.SeverityCounts[result.Info.Severity]++
		stats.TemplateCounts[result.TemplateID]++
	}

	return stats
}

func FilterResultsBySeverity(results []NucleiResult, severities []string) []NucleiResult {
	if len(severities) == 0 {
		return results
	}

	severityMap := make(map[string]bool)
	for _, s := range severities {
		severityMap[strings.ToLower(s)] = true
	}

	var filtered []NucleiResult
	for _, result := range results {
		if severityMap[strings.ToLower(result.Info.Severity)] {
			filtered = append(filtered, result)
		}
	}
	return filtered
}

func FilterResultsByTags(results []NucleiResult, tags []string) []NucleiResult {
	if len(tags) == 0 {
		return results
	}

	var filtered []NucleiResult
	for _, result := range results {
		for _, tag := range tags {
			for _, resultTag := range result.Info.Tags {
				if strings.EqualFold(tag, resultTag) {
					filtered = append(filtered, result)
					goto nextResult
				}
			}
		}
	nextResult:
	}
	return filtered
}

func GetUniqueHosts(results []NucleiResult) []string {
	hostMap := make(map[string]bool)
	for _, result := range results {
		hostMap[result.Host] = true
	}

	var hosts []string
	for host := range hostMap {
		hosts = append(hosts, host)
	}
	return hosts
}
