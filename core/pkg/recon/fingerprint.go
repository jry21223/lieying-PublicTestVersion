package recon

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Fingerprinter struct {
	url        string
	headers    http.Header
	body       string
	result     *FingerprintResult
	httpClient *http.Client
}

type FingerprintResult struct {
	URL        string
	CMS        string
	CMSVersion string
	WebServer  string
	Framework  string
	WAF        string
	TechStack  []string
}

func NewFingerprinter(url string) *Fingerprinter {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	return &Fingerprinter{
		url: url,
		result: &FingerprintResult{
			URL:       url,
			TechStack: []string{},
		},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (fp *Fingerprinter) Fingerprint() (*FingerprintResult, error) {
	fmt.Printf("🔍 开始Web指纹识别: %s\n", fp.url)

	resp, err := fp.httpClient.Get(fp.url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fp.headers = resp.Header

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fp.body = string(bodyBytes)

	fp.detectWebServer(resp.Header)
	fp.detectCMS()
	fp.detectFramework()
	fp.detectWAF(resp.Header)
	fp.detectTechStack()

	fmt.Printf("✅ 指纹识别完成:\n")
	fmt.Printf("   Web服务器: %s\n", fp.result.WebServer)
	fmt.Printf("   CMS: %s\n", fp.result.CMS)
	fmt.Printf("   框架: %s\n", fp.result.Framework)
	fmt.Printf("   WAF: %s\n", fp.result.WAF)
	fmt.Printf("   技术栈: %v\n", fp.result.TechStack)

	return fp.result, nil
}

func (fp *Fingerprinter) detectWebServer(headers http.Header) {
	server := headers.Get("Server")
	xPoweredBy := headers.Get("X-Powered-By")

	if server != "" {
		fp.result.WebServer = server
		if strings.Contains(strings.ToLower(server), "nginx") {
			fp.addTechStack("Nginx")
		} else if strings.Contains(strings.ToLower(server), "apache") {
			fp.addTechStack("Apache")
		} else if strings.Contains(strings.ToLower(server), "iis") {
			fp.addTechStack("IIS")
		} else if strings.Contains(strings.ToLower(server), "tomcat") {
			fp.addTechStack("Tomcat")
		}
	}

	if xPoweredBy != "" {
		if strings.Contains(xPoweredBy, "PHP") {
			fp.addTechStack("PHP")
		} else if strings.Contains(xPoweredBy, "ASP.NET") {
			fp.addTechStack("ASP.NET")
		} else if strings.Contains(xPoweredBy, "Node.js") {
			fp.addTechStack("Node.js")
		}
	}
}

func (fp *Fingerprinter) detectCMS() {
	bodyLower := strings.ToLower(fp.body)

	if strings.Contains(bodyLower, "wp-content") || strings.Contains(bodyLower, "wp-includes") {
		fp.result.CMS = "WordPress"
		fp.addTechStack("WordPress")
		fp.detectWordPressVersion()
	} else if strings.Contains(bodyLower, "drupal") {
		fp.result.CMS = "Drupal"
		fp.addTechStack("Drupal")
	} else if strings.Contains(bodyLower, "joomla") {
		fp.result.CMS = "Joomla"
		fp.addTechStack("Joomla")
	} else if strings.Contains(bodyLower, "thinkphp") || strings.Contains(bodyLower, "think_template") {
		fp.result.CMS = "ThinkPHP"
		fp.addTechStack("ThinkPHP")
	} else if strings.Contains(bodyLower, "typecho") {
		fp.result.CMS = "Typecho"
		fp.addTechStack("Typecho")
	} else if strings.Contains(bodyLower, "zblog") {
		fp.result.CMS = "Z-Blog"
		fp.addTechStack("Z-Blog")
	} else if strings.Contains(bodyLower, "dede") || strings.Contains(bodyLower, "dedecms") {
		fp.result.CMS = "织梦CMS"
		fp.addTechStack("织梦CMS")
	} else if strings.Contains(bodyLower, "empirecms") || strings.Contains(bodyLower, "帝国") {
		fp.result.CMS = "帝国CMS"
		fp.addTechStack("帝国CMS")
	} else if strings.Contains(bodyLower, "phpcms") {
		fp.result.CMS = "PHPCMS"
		fp.addTechStack("PHPCMS")
	}
}

func (fp *Fingerprinter) detectWordPressVersion() {
	re := regexp.MustCompile(`WordPress (\d+\.\d+(\.\d+)?)`)
	matches := re.FindStringSubmatch(fp.body)
	if len(matches) > 1 {
		fp.result.CMSVersion = matches[1]
	}
}

func (fp *Fingerprinter) detectFramework() {
	bodyLower := strings.ToLower(fp.body)

	if strings.Contains(bodyLower, "react") || strings.Contains(bodyLower, "__react") {
		fp.result.Framework = "React"
		fp.addTechStack("React")
	} else if strings.Contains(bodyLower, "vue") || strings.Contains(bodyLower, "__vue") {
		fp.result.Framework = "Vue.js"
		fp.addTechStack("Vue.js")
	} else if strings.Contains(bodyLower, "angular") {
		fp.result.Framework = "Angular"
		fp.addTechStack("Angular")
	} else if strings.Contains(bodyLower, "django") {
		fp.result.Framework = "Django"
		fp.addTechStack("Django")
	} else if strings.Contains(bodyLower, "flask") {
		fp.result.Framework = "Flask"
		fp.addTechStack("Flask")
	} else if strings.Contains(bodyLower, "spring") {
		fp.result.Framework = "Spring"
		fp.addTechStack("Spring")
	}
}

func (fp *Fingerprinter) detectWAF(headers http.Header) {
	for key, values := range headers {
		keyLower := strings.ToLower(key)
		valueLower := strings.ToLower(strings.Join(values, " "))

		if strings.Contains(keyLower, "cloudflare") || strings.Contains(valueLower, "cloudflare") {
			fp.result.WAF = "Cloudflare"
			fp.addTechStack("Cloudflare")
			return
		}

		if strings.Contains(keyLower, "akamai") || strings.Contains(valueLower, "akamai") {
			fp.result.WAF = "Akamai"
			fp.addTechStack("Akamai")
			return
		}

		if strings.Contains(keyLower, "sucuri") || strings.Contains(valueLower, "sucuri") {
			fp.result.WAF = "Sucuri"
			fp.addTechStack("Sucuri")
			return
		}

		if strings.Contains(keyLower, "x-powered-by") && strings.Contains(valueLower, "yunsuo") {
			fp.result.WAF = "云锁"
			fp.addTechStack("云锁")
			return
		}

		if strings.Contains(keyLower, "yunsuo") || strings.Contains(valueLower, "yunsuo") {
			fp.result.WAF = "云锁"
			fp.addTechStack("云锁")
			return
		}

		if strings.Contains(keyLower, "safe3") || strings.Contains(valueLower, "safe3") {
			fp.result.WAF = "安全狗"
			fp.addTechStack("安全狗")
			return
		}

		if strings.Contains(keyLower, "x-waf") || strings.Contains(valueLower, "x-waf") {
			fp.result.WAF = "Unknown WAF"
			return
		}
	}
}

func (fp *Fingerprinter) detectTechStack() {
	bodyLower := strings.ToLower(fp.body)

	if strings.Contains(bodyLower, "jquery") {
		fp.addTechStack("jQuery")
	}

	if strings.Contains(bodyLower, "bootstrap") {
		fp.addTechStack("Bootstrap")
	}

	if strings.Contains(bodyLower, "tailwind") {
		fp.addTechStack("Tailwind CSS")
	}

	if strings.Contains(bodyLower, "next.js") || strings.Contains(bodyLower, "nextjs") {
		fp.addTechStack("Next.js")
	}

	if strings.Contains(bodyLower, "nuxt") {
		fp.addTechStack("Nuxt.js")
	}

	if strings.Contains(bodyLower, "webpack") {
		fp.addTechStack("Webpack")
	}
}

func (fp *Fingerprinter) addTechStack(tech string) {
	for _, t := range fp.result.TechStack {
		if t == tech {
			return
		}
	}
	fp.result.TechStack = append(fp.result.TechStack, tech)
}

func (fp *Fingerprinter) GetResult() *FingerprintResult {
	return fp.result
}
