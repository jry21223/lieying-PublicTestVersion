// Mock API 服务 - 用于前端演示
// 无需后端即可展示界面效果

const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

export const mockApi = {
  // 目标管理
  targets: {
    list: async () => {
      await delay(500);
      return [
        { id: '1', name: 'Example Corp', type: 'domain', value: 'example.com', status: 'active', created_at: '2024-01-15' },
        { id: '2', name: 'Test Site', type: 'ip', value: '192.168.1.1', status: 'scanning', created_at: '2024-01-14' },
        { id: '3', name: 'University', type: 'domain', value: 'university.edu.cn', status: 'completed', created_at: '2024-01-13' },
      ];
    },
    create: async (data: any) => {
      await delay(300);
      return { id: Date.now().toString(), ...data, status: 'active', created_at: new Date().toISOString() };
    },
  },

  // 漏洞管理
  vulnerabilities: {
    list: async () => {
      await delay(600);
      return [
        { id: '1', title: 'SQL注入 - 登录接口', severity: 'critical', type: 'sqli', url: 'http://example.com/login', confirmed: true, created_at: '2024-01-15' },
        { id: '2', title: 'XSS - 评论功能', severity: 'high', type: 'xss', url: 'http://example.com/comments', confirmed: true, created_at: '2024-01-14' },
        { id: '3', title: '未授权访问 - 管理后台', severity: 'high', type: 'unauth', url: 'http://example.com/admin', confirmed: false, created_at: '2024-01-13' },
        { id: '4', title: '文件上传漏洞', severity: 'medium', type: 'upload', url: 'http://example.com/upload', confirmed: false, created_at: '2024-01-12' },
      ];
    },
    confirm: async (id: string) => {
      await delay(200);
      return { success: true };
    },
  },

  // 信息收集
  recon: {
    start: async (target: string) => {
      await delay(2000);
      return {
        target,
        subdomains: [
          { subdomain: 'www.example.com', ip: '93.184.216.34', status_code: 200, title: 'Example Domain' },
          { subdomain: 'api.example.com', ip: '93.184.216.35', status_code: 200, title: 'API Documentation' },
          { subdomain: 'admin.example.com', ip: '93.184.216.36', status_code: 200, title: 'Admin Panel' },
          { subdomain: 'mail.example.com', ip: '93.184.216.37', status_code: 403, title: '403 Forbidden' },
        ],
        ports: [
          { port: 80, service: 'HTTP', version: 'nginx/1.18.0' },
          { port: 443, service: 'HTTPS', version: 'nginx/1.18.0' },
          { port: 22, service: 'SSH', version: 'OpenSSH 8.2' },
          { port: 3306, service: 'MySQL', version: '5.7.32' },
        ],
        fingerprints: [
          { name: 'Nginx', version: '1.18.0', confidence: 95 },
          { name: 'PHP', version: '7.4.3', confidence: 88 },
          { name: 'WordPress', version: '5.8', confidence: 75 },
        ],
      };
    },
  },

  // 漏洞扫描
  scan: {
    start: async (target: string) => {
      await delay(3000);
      return {
        target,
        vulnerabilities: [
          { id: '1', title: 'SQL注入 - 登录接口', severity: 'critical', type: 'sqli', url: `${target}/login`, description: '登录接口存在SQL注入漏洞' },
          { id: '2', title: 'XSS - 评论功能', severity: 'high', type: 'xss', url: `${target}/comments`, description: '评论内容未过滤，存在XSS' },
          { id: '3', title: '未授权访问', severity: 'high', type: 'unauth', url: `${target}/admin`, description: '管理后台未授权可访问' },
        ],
      };
    },
  },

  // AI助手
  ai: {
    status: async () => {
      await delay(300);
      return { available: true, model: 'qwen2.5:7b', provider: 'ollama' };
    },
    generatePOC: async (data: any) => {
      await delay(2000);
      return {
        poc: `id: example-sqli
info:
  name: Example SQL Injection
  severity: critical
requests:
  - method: GET
    path:
      - "{{BaseURL}}/api/user?id=1' AND 1=1--"
    matchers:
      - type: word
        words:
          - "user_id"`,
      };
    },
  },

  // 教育SRC
  edu: {
    getStats: async () => {
      await delay(400);
      return {
        total: 109,
        '985': 39,
        '211': 70,
        provinces: { '北京': 26, '上海': 14, '江苏': 11, '陕西': 8, '湖北': 7 },
      };
    },
    getUniversities: async (params?: any) => {
      await delay(500);
      const universities = [
        { name: '北京大学', shortName: '北大', province: '北京', level: '985', domains: ['pku.edu.cn'] },
        { name: '清华大学', shortName: '清华', province: '北京', level: '985', domains: ['tsinghua.edu.cn'] },
        { name: '复旦大学', shortName: '复旦', province: '上海', level: '985', domains: ['fudan.edu.cn'] },
        { name: '上海交通大学', shortName: '上交', province: '上海', level: '985', domains: ['sjtu.edu.cn'] },
        { name: '浙江大学', shortName: '浙大', province: '浙江', level: '985', domains: ['zju.edu.cn'] },
        { name: '南京大学', shortName: '南大', province: '江苏', level: '985', domains: ['nju.edu.cn'] },
        { name: '西安交通大学', shortName: '西交', province: '陕西', level: '985', domains: ['xjtu.edu.cn'] },
        { name: '武汉大学', shortName: '武大', province: '湖北', level: '985', domains: ['whu.edu.cn'] },
      ];
      
      if (params?.level) {
        return universities.filter(u => u.level === params.level);
      }
      if (params?.province) {
        return universities.filter(u => u.province === params.province);
      }
      return universities;
    },
    scan: async (target: string) => {
      await delay(2500);
      return {
        target,
        results: [
          { title: 'SQL注入 - 正方教务登录', severity: 'critical', url: `${target}/login.aspx`, type: 'zfsoft' },
          { title: '信息泄露 - 学生信息', severity: 'high', url: `${target}/student/info`, type: 'info_leak' },
        ],
      };
    },
    getSystems: async () => {
      await delay(200);
      return [
        { name: '正方教务', type: 'zfsoft', description: '正方教务管理系统' },
        { name: '强智教务', type: 'qzsoft', description: '强智教务管理系统' },
        { name: '金智教务', type: 'kingo', description: '金智教育教务系统' },
        { name: 'URP教务', type: 'urp', description: 'URP综合教务系统' },
        { name: '青果教务', type: 'kingosoft', description: '青果教务管理系统' },
      ];
    },
  },
};

export default mockApi;
