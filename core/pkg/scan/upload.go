package scan

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
)

type UploadScanner struct {
	target     string
	results    []UploadResult
	httpClient *http.Client
}

type UploadResult struct {
	URL        string
	FormAction string
	InputName  string
	Type       string
	Payload    string
	Evidence   string
	Severity   string
	Confirmed  bool
}

type uploadForm struct {
	actionURL string
	inputName string
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
		".pht",
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

		body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
		resp.Body.Close()
		if err != nil {
			continue
		}

		form, ok := parseUploadForm(uploadURL, string(body))
		if !ok {
			continue
		}

		for _, ext := range maliciousExtensions {
			filename := fmt.Sprintf("shell-%d%s", time.Now().UnixNano(), ext)
			result, ok := us.tryDangerousUpload(form.actionURL, form.inputName, filename)
			if ok {
				us.results = append(us.results, result)
				fmt.Printf("⚠️  发现文件上传漏洞: %s\n", form.actionURL)
				return
			}
		}
	}
}

func parseUploadForm(pageURL, body string) (uploadForm, bool) {
	formBlockRegex := regexp.MustCompile(`(?is)<form\b[^>]*>.*?</form>`)
	actionRegex := regexp.MustCompile(`(?is)\baction\s*=\s*["']?([^"'\s>]+)`)
	inputTagRegex := regexp.MustCompile(`(?is)<input\b[^>]*>`)
	fileTypeRegex := regexp.MustCompile(`(?is)\btype\s*=\s*["']?file["']?`)
	nameRegex := regexp.MustCompile(`(?is)\bname\s*=\s*["']?([^"'\s>]+)`)

	forms := formBlockRegex.FindAllString(body, -1)
	for _, formBlock := range forms {
		lowerForm := strings.ToLower(formBlock)
		if !strings.Contains(lowerForm, "type=\"file\"") && !strings.Contains(lowerForm, "type='file'") && !strings.Contains(lowerForm, "type=file") {
			continue
		}

		actionURL := pageURL
		if matches := actionRegex.FindStringSubmatch(formBlock); len(matches) > 1 {
			resolved, ok := resolveFormAction(pageURL, matches[1])
			if !ok {
				continue
			}
			if !isSameOrigin(pageURL, resolved) {
				continue
			}
			actionURL = resolved
		}

		inputName := "file"
		inputTags := inputTagRegex.FindAllString(formBlock, -1)
		for _, inputTag := range inputTags {
			if !fileTypeRegex.MatchString(inputTag) {
				continue
			}
			if matches := nameRegex.FindStringSubmatch(inputTag); len(matches) > 1 {
				inputName = matches[1]
			}
			break
		}

		return uploadForm{actionURL: actionURL, inputName: inputName}, true
	}

	return uploadForm{}, false
}

func resolveFormAction(pageURL, action string) (string, bool) {
	page, err := url.Parse(pageURL)
	if err != nil {
		return "", false
	}
	actionURL, err := url.Parse(action)
	if err != nil {
		return "", false
	}
	return page.ResolveReference(actionURL).String(), true
}

func (us *UploadScanner) tryDangerousUpload(uploadURL, inputName, filename string) (UploadResult, bool) {
	executionMarker := executionMarkerForFilename(filename)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(inputName, filename)
	if err != nil {
		return UploadResult{}, false
	}

	_, err = part.Write([]byte("<?php echo md5('" + markerForFilename(filename) + "'); ?>"))
	if err != nil {
		return UploadResult{}, false
	}
	if err := writer.Close(); err != nil {
		return UploadResult{}, false
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL, body)
	if err != nil {
		return UploadResult{}, false
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := us.httpClient.Do(req)
	if err != nil {
		return UploadResult{}, false
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return UploadResult{}, false
	}

	bodyLower := strings.ToLower(string(respBody))
	if resp.StatusCode != http.StatusOK {
		return UploadResult{}, false
	}
	if !strings.Contains(bodyLower, "success") && !strings.Contains(bodyLower, "uploaded") && !strings.Contains(bodyLower, "path") && !strings.Contains(bodyLower, filename) {
		return UploadResult{}, false
	}

	uploadedFileURL, ok := extractUploadedFileURL(uploadURL, string(respBody))
	if !ok {
		return UploadResult{}, false
	}
	if !us.isUploadedFileReachable(uploadURL, uploadedFileURL, executionMarker) {
		return UploadResult{}, false
	}

	return UploadResult{
		URL:        uploadedFileURL,
		FormAction: uploadURL,
		InputName:  inputName,
		Type:       "File Upload",
		Payload:    filename,
		Evidence:   fmt.Sprintf("上传危险扩展文件后返回成功迹象，且文件可访问并包含上传标记: %s", uploadedFileURL),
		Severity:   "High",
		Confirmed:  true,
	}, true
}

func markerForFilename(filename string) string {
	base := strings.TrimSuffix(path.Base(filename), path.Ext(filename))
	return strings.Replace(base, "shell-", "lieying-upload-marker-", 1)
}

func executionMarkerForFilename(filename string) string {
	sum := md5.Sum([]byte(markerForFilename(filename)))
	return fmt.Sprintf("%x", sum)
}

func extractUploadedFileURL(uploadURL, responseBody string) (string, bool) {
	pathRegex := regexp.MustCompile(`(?i)(?:"|')(?:path|url)(?:"|')\s*:\s*(?:"|')([^"']+)(?:"|')`)
	matches := pathRegex.FindStringSubmatch(responseBody)
	if len(matches) <= 1 {
		return "", false
	}

	resolved, ok := resolveFormAction(uploadURL, matches[1])
	if !ok {
		return "", false
	}
	if !isSameOrigin(uploadURL, resolved) {
		return "", false
	}
	return resolved, true
}

func isSameOrigin(baseURL, targetURL string) bool {
	base, err := url.Parse(baseURL)
	if err != nil {
		return false
	}
	target, err := url.Parse(targetURL)
	if err != nil {
		return false
	}
	return base.Scheme == target.Scheme && effectivePort(base) == effectivePort(target) && base.Hostname() == target.Hostname()
}

func effectivePort(u *url.URL) string {
	if port := u.Port(); port != "" {
		return port
	}
	if u.Scheme == "https" {
		return "443"
	}
	if u.Scheme == "http" {
		return "80"
	}
	return ""
}

func (us *UploadScanner) isUploadedFileReachable(uploadURL, uploadedFileURL, marker string) bool {
	if !isSameOrigin(uploadURL, uploadedFileURL) {
		return false
	}

	client := *us.httpClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Get(uploadedFileURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}
	if !isSameOrigin(uploadURL, resp.Request.URL.String()) {
		return false
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return false
	}

	return strings.Contains(string(body), marker)
}

func (us *UploadScanner) GetResults() []UploadResult {
	return us.results
}
