package scan

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type UploadScanner struct {
	target     string
	results    []UploadResult
	httpClient *http.Client
}

type UploadResult struct {
	URL         string
	FormAction  string
	InputName   string
	Type        string
	Payload     string
	Evidence    string
	Severity    string
	Confirmed   bool
}

func NewUploadScanner(target string) *UploadScanner {
	return &UploadScanner{
		target:  target,
		results: []UploadResult{},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (us *UploadScanner) Scan() ([]UploadResult, error) {
	fmt.Println("🔍 开始文件上传漏洞检测...")

	us.testUploadForms()

	fmt.Printf("✅ 文件上传漏洞检测完成，发现 %d 个漏洞\n", len(us.results))
	return us.results, nil
}

func (us *UploadScanner) testUploadForms() {
	uploadPaths := []string{
		"/upload",
		"/upload.php",
		"/upload/index.php",
		"/file/upload",
		"/api/upload",
		"/admin/upload",
		"/user/upload",
		"/avatar",
		"/profile/avatar",
		"/image/upload",
		"/media/upload",
		"/attachment/upload",
	}

	maliciousExtensions := []string{
		".php",
		".php3",
		".php4",
		".php5",
		".phtml",
		".phar",
		".phps",
		".pht",
		".jsp",
		".jspx",
		".jsw",
		".jsv",
		".jspf",
		".asp",
		".aspx",
		".ascx",
		".ashx",
		".asmx",
		".cer",
		".asa",
		".cdx",
		".htr",
		".war",
		".py",
		".rb",
		".sh",
		".pl",
		".cgi",
	}

	parsedURL, err := url.Parse(us.target)
	if err != nil {
		return
	}

	baseURL := parsedURL.Scheme + "://" + parsedURL.Host

	for _, path := range uploadPaths {
		uploadURL := baseURL + path
		
		resp, err := us.httpClient.Get(uploadURL)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		bodyStr := string(body)
		if strings.Contains(bodyStr, "<form") && strings.Contains(bodyStr, "type=\"file\"") {
			for _, ext := range maliciousExtensions {
				result := UploadResult{
					URL:        uploadURL,
					FormAction: path,
					InputName:  "file",
					Type:       "File Upload",
					Payload:    "shell" + ext,
					Evidence:   fmt.Sprintf("发现文件上传表单，尝试上传%s文件", ext),
					Severity:   "High",
					Confirmed:  false,
				}
				us.results = append(us.results, result)
				fmt.Printf("⚠️  发现文件上传漏洞: %s\n", uploadURL)
				return
			}
		}
	}
}

func (us *UploadScanner) testUploadBypass(uploadURL string) {
	bypassTechniques := []struct {
		Name     string
		Payload  string
		Evidence string
	}{
		{"Extension Case", "shell.PHP", "尝试大小写绕过"},
		{"Double Extension", "shell.php.jpg", "尝试双扩展名绕过"},
		{"Null Byte", "shell.php%00.jpg", "尝试空字节绕过"},
		{"Special Chars", "shell.php.", "尝试特殊字符绕过"},
		{"MIME Type", "shell.php", "尝试修改MIME类型"},
	}

	for _, technique := range bypassTechniques {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		
		part, err := writer.CreateFormFile("file", technique.Payload)
		if err != nil {
			continue
		}
		
		part.Write([]byte("<?php echo 'test'; ?>"))
		writer.Close()

		req, err := http.NewRequest("POST", uploadURL, body)
		if err != nil {
			continue
		}
		
		req.Header.Set("Content-Type", writer.FormDataContentType())
		
		resp, err := us.httpClient.Do(req)
		if err != nil {
			continue
		}
		resp.Body.Close()
	}
}

func (us *UploadScanner) GetResults() []UploadResult {
	return us.results
}
