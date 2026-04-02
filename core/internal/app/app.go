package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kunlun-sec/lunying/internal/api"
	"github.com/kunlun-sec/lunying/internal/models"
	"github.com/kunlun-sec/lunying/internal/repository"
	"github.com/kunlun-sec/lunying/internal/service"
	"github.com/kunlun-sec/lunying/pkg/logger"
	"github.com/kunlun-sec/lunying/pkg/network"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var (
	db            *gorm.DB
	networkEngine *network.NetworkEngine
)

func InitializeDB(dbPath string) (*gorm.DB, error) {
	sqliteDB := sqlite.Open(dbPath)

	db, err := gorm.Open(sqliteDB, &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	if err := db.AutoMigrate(
		&models.Target{},
		&models.Vulnerability{},
		&models.VulnerabilityComment{},
		&models.Asset{},
		&models.University{},
		&models.EduSystem{},
		&models.EduScanResult{},
		&models.VulnTechnique{},
		&models.VulnPracticeResult{},
		&models.AIConfig{},
		&models.ProxyConfig{},
		&models.NetworkConfig{},
		&models.ActiveScanTask{},
		&models.ScanResult{},
		&models.InfoCollectionTask{},
		&models.NetworkRequest{},
	); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	initializeSampleData(db)

	return db, nil
}

func InitializeNetworkEngine(config *models.NetworkConfig) error {
	engineConfig := &network.EngineConfig{
		Proxy:          "",
		Concurrency:   config.Concurrency,
		QPS:            float64(config.RateLimit),
		Timeout:        config.Timeout,
		UserAgent:      config.UserAgent,
		FollowRedirects: config.FollowRedirects,
	}

	var err error
	networkEngine, err = network.NewNetworkEngine(engineConfig)
	return err
}

func GetDB() *gorm.DB {
	return db
}

func GetNetworkEngine() *network.NetworkEngine {
	return networkEngine
}

func RunServer() {
	fmt.Println("=====================================")
	fmt.Println("  猎影渗透测试平台 - Core")
	fmt.Println("  Lieying Penetration Testing Platform")
	fmt.Println("=====================================")
	fmt.Println("  昆仑安全实验室(前逍遥安全实验室-逍遥)")
	fmt.Println("  KunLun Security Lab (Former XiaoYao Security Lab - XiaoYao)")
	fmt.Println("=====================================")

	dbPath := "data/lieying.db"
	
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatalf("创建数据目录失败: %v", err)
	}

	var err error
	db, err = InitializeDB(dbPath)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	fmt.Printf("✅ 数据库连接成功\n")

	var netConfig models.NetworkConfig
	if err := db.First(&netConfig).Error; err != nil {
		netConfig = models.NetworkConfig{
			Concurrency:     10,
			RateLimit:       50,
			Timeout:         30,
			Retries:         3,
			UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			FollowRedirects: true,
		}
		db.Create(&netConfig)
	}

	if err := InitializeNetworkEngine(&netConfig); err != nil {
		log.Printf("⚠️  网络引擎初始化失败: %v", err)
	} else {
		fmt.Println("✅ 网络引擎初始化完成")
	}

	if err := logger.Init("info", "logs", ""); err != nil {
		log.Printf("⚠️  日志初始化失败: %v", err)
	} else {
		fmt.Println("✅ 日志系统初始化完成")
	}

	fmt.Println("✅ 猎影平台初始化完成")
	fmt.Println()
	fmt.Println("🚀 猎影API服务器启动中...")
	fmt.Println("   监听地址: :8081")
	fmt.Println(fmt.Sprintf("   数据库: %s", dbPath))
	fmt.Println()

	StartAPIServer()
}

func StartAPIServer() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))

	targetRepo := repository.NewTargetRepository(db)
	vulnRepo := repository.NewVulnerabilityRepository(db)
	assetRepo := repository.NewAssetRepository(db)

	targetService := service.NewTargetService(targetRepo)
	vulnService := service.NewVulnerabilityService(vulnRepo)
	assetService := service.NewAssetService(assetRepo)

	apiHandler := api.NewAPI(
		targetService,
		vulnService,
		assetService,
		targetRepo,
		vulnRepo,
		assetRepo,
		db,
		networkEngine,
	)

	apiHandler.SetupRoutes(r)

	fmt.Println("API端点:")
	fmt.Println("  GET  /api/targets          - 列出所有目标")
	fmt.Println("  POST /api/targets          - 创建目标")
	fmt.Println("  DELETE /api/targets/:id    - 删除目标")
	fmt.Println("  GET  /api/vulnerabilities  - 列出所有漏洞")
	fmt.Println("  POST /api/vulnerabilities - 创建漏洞")
	fmt.Println("  POST /api/vulnerabilities/:id/confirm - 确认漏洞")
	fmt.Println("  GET  /api/targets/:id/assets/tree - 获取资产树")
	fmt.Println("  GET  /api/assets/:id      - 获取资产详情")
	fmt.Println("  PATCH /api/assets/:id      - 更新资产标签/备注")
	fmt.Println("  DELETE /api/assets/:id    - 删除资产")
	fmt.Println()
	fmt.Println("  新功能API端点:")
	fmt.Println("  POST /api/v1/info-collection    - 创建信息收集任务")
	fmt.Println("  POST /api/v1/vuln-scans         - 创建漏洞扫描任务")
	fmt.Println("  POST /api/v1/active-scan        - 主动扫描")
	fmt.Println("  GET  /api/v1/vuln-techniques    - 获取漏洞技巧列表")
	fmt.Println("  GET  /api/v1/vuln-techniques/:id - 获取技巧详情")
	fmt.Println("  POST /api/v1/vuln-techniques/practice - 执行漏洞实践")
	fmt.Println()

	log.Fatal(http.ListenAndServe(":8081", r))
}

func initializeSampleData(db *gorm.DB) {
	targets := []models.Target{
		{Name: "VIPC6资源网", Value: "vipc6.com", Type: "domain", Description: "VIPC6资源网站"},
		{Name: "百度", Value: "baidu.com", Type: "domain", Description: "百度搜索"},
		{Name: "测试示例", Value: "example.com", Type: "domain", Description: "测试用域名"},
	}

	var targetCount int64
	db.Model(&models.Target{}).Count(&targetCount)
	if targetCount == 0 {
		for i := range targets {
			db.Create(&targets[i])
		}
		fmt.Println("✅ 目标数据初始化完成")
	} else {
		fmt.Println("ℹ️  检测到已有目标数据，跳过目标初始化")
	}

	vulns := []models.Vulnerability{
		{Title: "SQL注入", Severity: "high", Status: "open", Description: "测试漏洞"},
		{Title: "XSS", Severity: "medium", Status: "confirmed", Description: "跨站脚本"},
	}

	var vulnCount int64
	db.Model(&models.Vulnerability{}).Count(&vulnCount)
	if vulnCount == 0 {
		for i := range vulns {
			db.Create(&vulns[i])
		}
		fmt.Println("✅ 漏洞数据初始化完成")
	} else {
		fmt.Println("ℹ️  检测到已有漏洞数据，跳过漏洞初始化")
	}

	universities := []models.University{
		{Name: "清华大学", Level: "985", Province: "北京", Type: "本科", Domain: "tsinghua.edu.cn"},
		{Name: "北京大学", Level: "985", Province: "北京", Type: "本科", Domain: "pku.edu.cn"},
		{Name: "浙江大学", Level: "985", Province: "浙江", Type: "本科", Domain: "zju.edu.cn"},
		{Name: "复旦大学", Level: "985", Province: "上海", Type: "本科", Domain: "fudan.edu.cn"},
		{Name: "上海交通大学", Level: "985", Province: "上海", Type: "本科", Domain: "sjtu.edu.cn"},
		{Name: "南京大学", Level: "985", Province: "江苏", Type: "本科", Domain: "nju.edu.cn"},
		{Name: "中国科学技术大学", Level: "985", Province: "安徽", Type: "本科", Domain: "ustc.edu.cn"},
		{Name: "哈尔滨工业大学", Level: "985", Province: "黑龙江", Type: "本科", Domain: "hit.edu.cn"},
		{Name: "西安交通大学", Level: "985", Province: "陕西", Type: "本科", Domain: "xjtu.edu.cn"},
		{Name: "同济大学", Level: "985", Province: "上海", Type: "本科", Domain: "tongji.edu.cn"},
		{Name: "北京航空航天大学", Level: "985", Province: "北京", Type: "本科", Domain: "buaa.edu.cn"},
		{Name: "中山大学", Level: "985", Province: "广东", Type: "本科", Domain: "sysu.edu.cn"},
		{Name: "华南理工大学", Level: "985", Province: "广东", Type: "本科", Domain: "scut.edu.cn"},
		{Name: "四川大学", Level: "985", Province: "四川", Type: "本科", Domain: "scu.edu.cn"},
		{Name: "武汉大学", Level: "985", Province: "湖北", Type: "本科", Domain: "whu.edu.cn"},
	}

	for i := range universities {
		db.Create(&universities[i])
	}

	eduSystems := []models.EduSystem{
		{Name: "统一身份认证(CAS)", Category: "认证系统"},
		{Name: "OA办公系统", Category: "办公系统"},
		{Name: "教务管理系统", Category: "教务系统"},
		{Name: "邮件系统", Category: "网络服务"},
		{Name: "VPN系统", Category: "网络服务"},
		{Name: "一卡通系统", Category: "校园服务"},
		{Name: "图书馆系统", Category: "校园服务"},
		{Name: "选课系统", Category: "教务系统"},
	}

	var eduSystemCount int64
	db.Model(&models.EduSystem{}).Count(&eduSystemCount)
	if eduSystemCount == 0 {
		for i := range eduSystems {
			db.Create(&eduSystems[i])
		}
		fmt.Println("✅ 教育系统数据初始化完成")
	} else {
		fmt.Println("ℹ️  检测到已有教育系统数据，跳过初始化")
	}

	techniques := []models.VulnTechnique{
		{
			Title:       "SQL注入漏洞挖掘",
			Category:    "注入类",
			Severity:    "high",
			Summary:     "通过构造恶意SQL语句获取数据库敏感信息",
			Description: "SQL注入是一种代码注入技术，通过在用户输入中注入恶意SQL语句，从而执行未授权的数据库操作。常见于登录框、搜索框、URL参数等用户可控输入点。",
			Payload:     "1' OR '1'='1\n1' UNION SELECT NULL--\n1' AND SLEEP(5)--",
			Fix建议:     "1. 使用参数化查询\n2. 输入过滤和验证\n3. 最小权限原则\n4. Web应用防火墙(WAF)",
			Examples:      "登录框绕过: admin' OR '1'='1\n联合查询: 1' UNION SELECT username,password FROM users--",
			References:   "https://owasp.org/www-community/attacks/SQL_Injection",
			Tags:         "SQL注入,数据库,OWASP",
			Impact:       "数据库泄露、用户数据窃取、服务器控制",
			Difficulty:   "intermediate",
		},
		{
			Title:       "XSS跨站脚本攻击",
			Category:    "跨站类",
			Severity:    "medium",
			Summary:     "在网页中注入恶意JavaScript代码",
			Description: "XSS攻击允许攻击者在受害者的浏览器中执行恶意脚本代码。分为反射型、存储型和DOM型三种。",
			Payload:     "<script>alert('XSS')</script>\n<img src=x onerror=alert('XSS')>",
			Fix建议:     "1. 输入过滤HTML特殊字符\n2. 输出编码\n3. HTTPOnly和Secure Cookie",
			Examples:      "弹窗测试: <script>alert(1)</script>",
			References:   "https://owasp.org/www-community/attacks/xss/",
			Tags:         "XSS,JavaScript,前端安全",
			Impact:       "会话劫持、钓鱼攻击、蠕虫传播",
			Difficulty:   "beginner",
		},
		{
			Title:       "CSRF跨站请求伪造",
			Category:    "认证类",
			Severity:    "medium",
			Summary:     "利用用户已登录身份发起恶意请求",
			Description: "CSRF攻击者诱导已登录用户在不知情的情况下向目标网站发送恶意请求。",
			Payload:     "<img src='http://target.com/transfer?to=attacker&amount=10000'>",
			Fix建议:     "1. CSRF Token\n2. 验证Referer/Origin\n3. SameSite Cookie",
			Examples:      "修改密码: <img src='http://bank.com/transfer?to=hacker&amount=10000'>",
			References:   "https://owasp.org/www-community/attacks/csrf",
			Tags:         "CSRF,会话,认证",
			Impact:       "账户篡改，资金转移、权限滥用",
			Difficulty:   "intermediate",
		},
		{
			Title:       "文件上传漏洞",
			Category:    "文件处理",
			Severity:    "high",
			Summary:     "上传恶意文件获取服务器权限",
			Description: "文件上传功能未严格验证上传文件类型和内容，攻击者可上传webshell等恶意文件。",
			Payload:     "<?php system($_GET['cmd']); ?>",
			Fix建议:     "1. 白名单验证文件扩展名\n2. MIME类型检测\n3. 文件内容检查\n4. 上传目录禁止执行",
			Examples:      "PHP webshell: <?php @eval($_POST['cmd']); ?>",
			References:   "https://owasp.org/www-community/vulnerabilities/Unrestricted_File_Upload",
			Tags:         "文件上传,webshell,服务器",
			Impact:       "服务器沦陷、webshell管理、数据窃取",
			Difficulty:   "intermediate",
		},
		{
			Title:       "敏感信息泄露",
			Category:    "信息收集",
			Severity:    "low",
			Summary:     "通过错误信息、调试接口暴露敏感信息",
			Description: "应用配置错误或开发遗留导致敏感信息泄露。",
			Payload:     "/.git/config\n/.env\n/admin/debug",
			Fix建议:     "1. 生产环境关闭调试模式\n2. 禁止目录遍历\n3. 统一错误页面",
			Examples:      "Git泄露: https://target.com/.git/config",
			References:   "https://owasp.org/www-project-web-security-testing-guide/",
			Tags:         "信息泄露,.git,敏感文件",
			Impact:       "源码泄露、配置暴露、进一步攻击",
			Difficulty:   "beginner",
		},
		{
			Title:       "弱口令与暴力破解",
			Category:    "认证类",
			Severity:    "high",
			Summary:     "利用弱密码或暴力破解获取账户权限",
			Description: "用户使用弱密码或系统缺乏防护机制，攻击者可通过字典或暴力破解获取账户权限。",
			Payload:     "admin/admin\nadmin/123456\nroot/root",
			Fix建议:     "1. 强密码策略\n2. 账户锁定机制\n3. 验证码防护\n4. 多因素认证",
			Examples:      "常见弱口令: admin/123456, root/root",
			References:   "https://weakpass.com/wordlist",
			Tags:         "暴力破解,弱口令,账户安全",
			Impact:       "账户沦陷、数据泄露、横向移动",
			Difficulty:   "beginner",
		},
		{
			Title:       "SSRF服务器端请求伪造",
			Category:    "注入类",
			Severity:    "high",
			Summary:     "利用服务器发起对内网资源的攻击",
			Description: "SSRF漏洞允许攻击者通过服务器向内网或外部系统发起请求。",
			Payload:     "http://localhost:80\nhttp://169.254.169.254/latest/meta-data/",
			Fix建议:     "1. URL白名单验证\n2. 禁止内网IP访问\n3. 协议限制",
			Examples:      "云元数据: http://169.254.169.254/",
			References:   "https://owasp.org/www-community/attacks/Server_Side_Request_Forgery",
			Tags:         "SSRF,内网,云安全",
			Impact:       "内网探测、云密钥泄露、服务沦陷",
			Difficulty:   "intermediate",
		},
		{
			Title:       "未授权访问",
			Category:    "权限类",
			Severity:    "high",
			Summary:     "绕过认证直接访问敏感功能",
			Description: "应用程序访问控制缺陷，攻击者可直接访问未授权的敏感功能。",
			Payload:     "/admin/users/delete?id=1\n/api/v1/user/1001/info",
			Fix建议:     "1. 权限层级清晰划分\n2. 所有接口鉴权\n3. 最小权限原则",
			Examples:      "IDOR: /api/user/1001 -> /api/user/1002",
			References:   "https://portswigger.net/web-security/access-control",
			Tags:         "未授权,越权,IDOR",
			Impact:       "数据窃取、权限滥用、管理功能滥用",
			Difficulty:   "intermediate",
		},
		{
			Title:       "命令执行漏洞",
			Category:    "命令注入",
			Severity:    "critical",
			Summary:     "通过注入系统命令获取服务器控制权",
			Description: "应用程序将用户输入传递给系统命令执行函数，攻击者可通过注入恶意命令字符执行任意系统命令。",
			Payload:     "; ls\n| cat /etc/passwd\n&& whoami",
			Fix建议:     "1. 避免使用系统命令执行函数\n2. 输入严格过滤\n3. 参数化命令执行",
			Examples:      "Ping检测: ; cat /etc/passwd",
			References:   "https://owasp.org/www-community/attacks/Command_Injection",
			Tags:         "命令注入,RCE,系统安全",
			Impact:       "服务器完全控制、数据泄露、横向移动",
			Difficulty:   "intermediate",
		},
		{
			Title:       "OAuth2.0安全漏洞",
			Category:    "认证类",
			Severity:    "high",
			Summary:     "OAuth认证实现缺陷导致账户劫持",
			Description: "OAuth2.0实现不安全，攻击者可利用回调URL验证缺陷等劫持用户会话。",
			Payload:     "redirect_uri=https://attacker.com/callback\nscope=read,write,admin",
			Fix建议:     "1. 严格验证回调URL\n2. scope最小授权\n3. state参数验证",
			Examples:      "回调绕过: redirect_uri=http://target.com..//attacker.com",
			References:   "https://oauth.net/2.0/security/",
			Tags:         "OAuth,认证,SSO",
			Impact:       "账户劫持、权限滥用、会话窃取",
			Difficulty:   "advanced",
		},
	}

	var techniqueCount int64
	db.Model(&models.VulnTechnique{}).Count(&techniqueCount)
	if techniqueCount == 0 {
		for i := range techniques {
			db.Create(&techniques[i])
		}
		fmt.Println("✅ 漏洞技巧数据初始化完成")
	} else {
		fmt.Println("ℹ️  检测到已有漏洞技巧数据，跳过初始化")
	}

	fmt.Println("✅ 示例数据初始化完成")
}
