package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Target struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	Tags        string    `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Target) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

type NetworkRequest struct {
	ID              string         `gorm:"primaryKey" json:"id"`
	Module          string         `json:"module"`
	Target          string         `json:"target"`
	Method          string         `json:"method"`
	URL             string         `json:"url"`
	RequestHeaders  string         `gorm:"type:text" json:"request_headers"`
	RequestBody     string         `gorm:"type:text" json:"request_body"`
	ResponseStatus  int            `json:"response_status"`
	ResponseHeaders string         `gorm:"type:text" json:"response_headers"`
	ResponseBody    string         `gorm:"type:text" json:"response_body"`
	Duration        int64          `json:"duration"`
	Error           string         `gorm:"type:text" json:"error,omitempty"`
	Timestamp       time.Time      `json:"timestamp"`
}

func (n *NetworkRequest) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}

type Vulnerability struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	FixStatus   string    `json:"fix_status"`
	URL         string    `json:"url"`
	CVEID       string    `json:"cve_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (v *Vulnerability) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	if v.Status == "" {
		v.Status = "pending"
	}
	if v.FixStatus == "" {
		v.FixStatus = "pending"
	}
	return nil
}

type Asset struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	TargetID    string    `json:"target_id"`
	ParentID    string    `json:"parent_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Value       string    `json:"value"`
	Status      string    `json:"status"`
	Importance  int       `json:"importance"`
	Description string    `json:"description"`
	Tags        string    `json:"tags"`
	Notes       string    `json:"notes"`
	FirstSeen   *time.Time `json:"first_seen"`
	LastSeen    *time.Time `json:"last_seen"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (a *Asset) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	if a.Status == "" {
		a.Status = "active"
	}
	if a.Importance == 0 {
		a.Importance = 50
	}
	return nil
}

type InfoCollectionTask struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Target    string    `json:"target"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"`
	Results   string    `gorm:"type:text" json:"results"`
	Error     string    `gorm:"type:text" json:"error,omitempty"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (i *InfoCollectionTask) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	if i.Status == "" {
		i.Status = "pending"
	}
	return nil
}

type VulnerabilityComment struct {
	ID              string    `gorm:"primaryKey" json:"id"`
	VulnerabilityID string    `json:"vulnerability_id"`
	Content         string    `json:"content"`
	CreatedAt       time.Time `json:"created_at"`
}

func (c *VulnerabilityComment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

type AIConfig struct {
	ID                   string    `gorm:"primaryKey" json:"id"`
	Mode                 string    `json:"mode"`
	Active               bool      `json:"active"`
	LocalBaseURL         string    `json:"local_base_url"`
	LocalModel           string    `json:"local_model"`
	LocalTemperature     float64   `json:"local_temperature"`
	LocalMaxTokens       int       `json:"local_max_tokens"`
	CloudProvider        string    `json:"cloud_provider"`
	CloudAPIKey          string    `json:"cloud_api_key"`
	CloudBaseURL         string    `json:"cloud_base_url"`
	CloudModel           string    `json:"cloud_model"`
	CloudTemperature     float64   `json:"cloud_temperature"`
	CloudMaxTokens       int       `json:"cloud_max_tokens"`
	SystemPrompt         string    `json:"system_prompt"`
	EnableFunctionCalling bool      `json:"enable_function_calling"`
	EnableActiveSuggestions bool    `json:"enable_active_suggestions"`
	EnableContextAwareness bool    `json:"enable_context_awareness"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (c *AIConfig) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	if c.Mode == "" {
		c.Mode = "local"
	}
	if c.LocalBaseURL == "" {
		c.LocalBaseURL = "http://localhost:11434"
	}
	if c.LocalModel == "" {
		c.LocalModel = "llama3"
	}
	if c.LocalTemperature == 0 {
		c.LocalTemperature = 0.7
	}
	if c.LocalMaxTokens == 0 {
		c.LocalMaxTokens = 4096
	}
	if c.SystemPrompt == "" {
		c.SystemPrompt = "你是一个专业的网络安全助手。"
	}
	return nil
}

type ReportTemplate struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *ReportTemplate) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

type Report struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Target    string    `json:"target"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

type University struct {
	ID       string    `gorm:"primaryKey" json:"id"`
	Name     string    `json:"name"`
	Level    string    `json:"level"`
	Province string    `json:"province"`
	Type     string    `json:"type"`
	Domain   string    `json:"domain"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *University) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type EduSystem struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	Name       string    `json:"name"`
	Category   string    `json:"category"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *EduSystem) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

type VulnTechnique struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	Title         string    `json:"title"`
	Category      string    `json:"category"`
	Severity      string    `json:"severity"`
	Summary       string    `json:"summary"`
	Description   string    `gorm:"type:text" json:"description"`
	Payload       string    `gorm:"type:text" json:"payload"`
	Fix建议       string    `gorm:"type:text" json:"fix_suggestion"`
	Examples      string    `gorm:"type:text" json:"examples"`
	References    string    `gorm:"type:text" json:"references"`
	Tags          string    `json:"tags"`
	Impact        string    `json:"impact"`
	Difficulty    string    `json:"difficulty"`
	VideoURL     string    `json:"video_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (v *VulnTechnique) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	if v.Severity == "" {
		v.Severity = "medium"
	}
	if v.Difficulty == "" {
		v.Difficulty = "intermediate"
	}
	return nil
}

type VulnPracticeResult struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	TechniqueID string    `json:"technique_id"`
	Target      string    `json:"target"`
	URL         string    `json:"url"`
	Payload     string    `json:"payload"`
	Result      string    `gorm:"type:text" json:"result"`
	Status      string    `json:"status"`
	Notes       string    `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p *VulnPracticeResult) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

type EduScanResult struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	University string    `json:"university"`
	System     string    `json:"system"`
	URL        string    `json:"url"`
	Status     string    `json:"status"`
	Severity   string    `json:"severity"`
	TaskType   string    `json:"task_type"`
	Progress   int       `json:"progress"`
	Results    string    `gorm:"type:text" json:"results"`
	Error      string    `gorm:"type:text" json:"error,omitempty"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

func (r *EduScanResult) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	if r.Status == "" {
		r.Status = "pending"
	}
	if r.Progress == 0 {
		r.Progress = 0
	}
	return nil
}

type ActiveScanTask struct {
	ID             string         `gorm:"primaryKey" json:"id"`
	Name           string         `json:"name"`
	Target         string         `json:"target"`
	Type           string         `json:"type"`
	Status         string         `json:"status"`
	Options        string         `gorm:"type:text" json:"options"`
	Progress       int            `json:"progress"`
	TotalItems     int            `json:"total_items"`
	ProcessedItems int            `json:"processed_items"`
	Error          string         `gorm:"type:text" json:"error"`
	StartedAt      *time.Time     `json:"started_at"`
	FinishedAt     *time.Time     `json:"finished_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *ActiveScanTask) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	if t.Status == "" {
		t.Status = "pending"
	}
	if t.Progress == 0 {
		t.Progress = 0
	}
	return nil
}

type ScanResult struct {
	ID              string         `gorm:"primaryKey" json:"id"`
	TaskID          string         `gorm:"index" json:"task_id"`
	Type            string         `json:"type"`
	Target          string         `json:"target"`
	URL             string         `json:"url"`
	Title           string         `json:"title"`
	Description     string         `gorm:"type:text" json:"description"`
	Data            string         `gorm:"type:text" json:"data"`
	Severity        string         `json:"severity"`
	Confirmed       bool           `json:"confirmed"`
	CVEID           string         `json:"cve_id"`
	CNVDID          string         `json:"cnvd_id"`
	Payload         string         `gorm:"type:text" json:"payload"`
	Request         string         `gorm:"type:text" json:"request"`
	Response        string         `gorm:"type:text" json:"response"`
	ReproduceSteps  string         `gorm:"type:text" json:"reproduce_steps"`
	FixSuggestion   string         `gorm:"type:text" json:"fix_suggestion"`
	References      string         `gorm:"type:text" json:"references"`
	Tags            string         `json:"tags"`
	Proof           string         `gorm:"type:text" json:"proof"`
	RiskLevel       int            `json:"risk_level"`
	AffectedVersion string         `json:"affected_version"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r *ScanResult) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

type ProxyConfig struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Enabled   bool           `json:"enabled"`
	Type      string         `json:"type"`
	Host      string         `json:"host"`
	Port      int            `json:"port"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (p *ProxyConfig) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.Type == "" {
		p.Type = "http"
	}
	if p.Port == 0 {
		p.Port = 8080
	}
	return nil
}

type NetworkConfig struct {
	ID              string         `gorm:"primaryKey" json:"id"`
	Concurrency     int            `json:"concurrency"`
	RateLimit       int            `json:"rate_limit"`
	Timeout         int            `json:"timeout"`
	Retries         int            `json:"retries"`
	UserAgent       string         `json:"user_agent"`
	FollowRedirects bool           `json:"follow_redirects"`
	SkipTLSVerify   bool           `json:"skip_tls_verify"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (nc *NetworkConfig) BeforeCreate(tx *gorm.DB) error {
	if nc.ID == "" {
		nc.ID = uuid.New().String()
	}
	if nc.Concurrency == 0 {
		nc.Concurrency = 10
	}
	if nc.RateLimit == 0 {
		nc.RateLimit = 50
	}
	if nc.Timeout == 0 {
		nc.Timeout = 30
	}
	if nc.Retries == 0 {
		nc.Retries = 3
	}
	if nc.UserAgent == "" {
		nc.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
	}
	return nil
}
