import React, { useState, useEffect } from 'react';
import { eduApi, vulnTechniquesApi, activeScanApi } from '../services/api';

interface ScanResult {
  id: string;
  domain: string;
  subdomain: string;
  ip: string;
  ports: string;
  vuln_count: number;
  last_scan: string;
}

interface University {
  id: string;
  name: string;
  level: string;
  province: string;
  type: string;
  domain: string;
  created_at: string;
}

interface VulnTechnique {
  id: string;
  title: string;
  category: string;
  severity: string;
  summary: string;
  description: string;
  payload: string;
  fix_suggestion: string;
  examples: string;
  references: string;
  tags: string;
  impact: string;
  difficulty: string;
  video_url?: string;
  created_at: string;
}

interface PracticeResult {
  id: string;
  technique_id: string;
  target: string;
  url: string;
  payload: string;
  result: string;
  status: string;
  notes: string;
  created_at: string;
}

interface POC {
  id: string;
  name: string;
  cve_id: string;
  severity: string;
  target: string;
  description: string;
  payload: string;
  impact: string;
  remediation: string;
  references: string[];
  created_at: string;
}

interface SRCPlatform {
  id: string;
  name: string;
  url: string;
  type: string;
  bounty_range: string;
  rank: string;
  domain: string;
}

interface ActiveScanTask {
  id: string;
  name: string;
  target: string;
  type: string;
  status: string;
  progress: number;
  created_at: string;
}

interface ScanResultItem {
  id: string;
  task_id: string;
  type: string;
  target: string;
  data: string;
  severity: string;
  confirmed: boolean;
  created_at: string;
}

const EduPanel: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'overview' | 'poc' | 'src' | 'batch'>('overview');
  const [stats, setStats] = useState({
    totalUniversities: 0,
    totalDomains: 0,
    totalVulns: 0,
    srcPlatforms: 0,
  });
  const [recentScans, setRecentScans] = useState<ScanResult[]>([]);
  const [universities, setUniversities] = useState<University[]>([]);
  const [loading, setLoading] = useState(false);

  const [techniques, setTechniques] = useState<VulnTechnique[]>([]);
  const [selectedTechnique, setSelectedTechnique] = useState<VulnTechnique | null>(null);
  const [practiceMode, setPracticeMode] = useState<'sync' | 'async'>('sync');
  const [practiceTarget, setPracticeTarget] = useState('');
  const [practiceURL, setPracticeURL] = useState('');
  const [practicePayload, setPracticePayload] = useState('');
  const [practiceResults, setPracticeResults] = useState<PracticeResult[]>([]);
  const [isPracticing, setIsPracticing] = useState(false);
  const [practiceLoading, setPracticeLoading] = useState(false);
  const [filterCategory, setFilterCategory] = useState('');
  const [filterSeverity, setFilterSeverity] = useState('');
  const [filterKeyword, setFilterKeyword] = useState('');

  const [pocs, setPocs] = useState<POC[]>([]);
  const [selectedPOC, setSelectedPOC] = useState<POC | null>(null);
  const [pocLoading, setPocLoading] = useState(false);
  const [pocFilter, setPocFilter] = useState({ severity: '', search: '' });

  const [srcPlatforms] = useState<SRCPlatform[]>([
    { id: '1', name: 'EDU SRC', url: 'https://src.edu.cn', type: '教育SRC', bounty_range: '积分制', rank: 'Top', domain: 'edu.cn' },
    { id: '2', name: '漏洞盒子', url: 'https://www.boxing.com', type: '民间SRC', bounty_range: '500-10万', rank: 'Top', domain: 'ebox' },
    { id: '3', name: '补天', url: 'https://butian.360.cn', type: '民间SRC', bounty_range: '500-50万', rank: 'Top', domain: 'butian' },
    { id: '4', name: 'CNVD', url: 'https://www.cnvd.org.cn', type: '国家队', bounty_range: '证书', rank: '官方', domain: 'cnvd' },
    { id: '5', name: '教育漏洞响应平台', url: 'https://src.edu.cn', type: '教育SRC', bounty_range: '积分制', rank: '官方', domain: 'src' },
  ]);

  const [activeTasks, setActiveTasks] = useState<ActiveScanTask[]>([]);
  const [scanResults, setScanResults] = useState<Record<string, ScanResultItem[]>>({});

  useEffect(() => {
    loadData();
    loadTechniques();
    loadPracticeResults();
    loadPOCs();
    loadActiveTasks();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [statsData, universitiesData] = await Promise.all([
        eduApi.getStats(),
        eduApi.getUniversities(),
      ]);

      const unis = universitiesData.universities || universitiesData || [];
      setUniversities(unis);

      const eduStats = statsData as any;
      setStats({
        totalUniversities: eduStats.total || unis.length,
        totalDomains: eduStats.total || unis.length,
        totalVulns: 0,
        srcPlatforms: (eduStats.by_level?.['985'] || 0) + (eduStats.by_level?.['211'] || 0),
      });
    } catch (error) {
      console.error('加载数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadTechniques = async () => {
    try {
      const data = await vulnTechniquesApi.list({
        category: filterCategory || undefined,
        severity: filterSeverity || undefined,
        keyword: filterKeyword || undefined,
      });
      setTechniques(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('加载漏洞技巧失败:', error);
      setTechniques([]);
    }
  };

  const loadPracticeResults = async () => {
    try {
      const data = await vulnTechniquesApi.getPracticeResults();
      setPracticeResults(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('加载实践结果失败:', error);
    }
  };

  const loadPOCs = () => {
    const expertPOCs: POC[] = [
      {
        id: 'poc-001',
        name: 'Apache Struts2远程代码执行',
        cve_id: 'CVE-2021-31805',
        severity: 'critical',
        target: 'Apache Struts2',
        description: 'S2-062远程代码执行漏洞，Apache Struts2框架存在远程代码执行漏洞，攻击者可利用该漏洞在目标服务器上执行任意代码。',
        payload: '%{(#req=#context.get(\'coommons.io.JdbcRowSetImpl\')).(#dm=@ognl.OgnlContext@defaultMemberAccess).(#ct=#req.getContainer()).(#ot=#ct.getValue()).(#bt=#ot.findParameter(\'accessObject\')).(#a=#bt.getClass()).(#b=#a.getClass()).(#m=#b.getMethod(\'toString\')).(#poc=#m.invoke(#a,\'test\'))}',
        impact: '服务器完全控制，敏感数据泄露，横向移动',
        remediation: '升级到Struts2 2.5.30或更高版本，使用WAF防护',
        references: ['https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-31805'],
        created_at: '2024-01-15',
      },
      {
        id: 'poc-002',
        name: 'Spring Core RCE',
        cve_id: 'CVE-2022-22965',
        severity: 'critical',
        target: 'Spring Framework',
        description: 'Spring Framework远程代码执行漏洞（Log4j之前最严重的Java漏洞），允许通过参数绑定机制实现远程代码执行。',
        payload: 'class.module.classLoader.resources.context.parent.pipeline.first.pattern=%{cže}',
        impact: '服务器完全控制，敏感数据泄露，挖矿木马植入',
        remediation: '升级到Spring Framework 5.3.18+或5.2.20+，禁用参数绑定',
        references: ['https://spring.io/blog/2022/03/31/spring-framework-rce-early-announcement/'],
        created_at: '2024-01-20',
      },
      {
        id: 'poc-003',
        name: 'Fastjson反序列化',
        cve_id: 'CVE-2022-25845',
        severity: 'high',
        target: 'Fastjson',
        description: 'Fastjson存在反序列化漏洞，攻击者可通过构造恶意JSON实现远程代码执行。',
        payload: '{"@type":"com.alibaba.fastjson.parser.ParserConfig","name":"test"}',
        impact: '远程代码执行，服务器沦陷',
        remediation: '升级到Fastjson 1.2.83+，使用安全框架替代',
        references: ['https://help.aliyun.com/noti/'],
        created_at: '2024-01-25',
      },
      {
        id: 'poc-004',
        name: 'Shiro反序列化',
        cve_id: 'CVE-2022-32532',
        severity: 'high',
        target: 'Apache Shiro',
        description: 'Apache Shiro RememberMe反序列化漏洞，攻击者可利用Shiro的RememberMe功能实现远程代码执行。',
        payload: 'rememberMe=admin|123456',
        impact: '账户接管，远程代码执行',
        remediation: '升级到Shiro 1.9.1+，更换加密密钥',
        references: ['https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-32532'],
        created_at: '2024-02-01',
      },
      {
        id: 'poc-005',
        name: 'Weblogic反序列化',
        cve_id: 'CVE-2023-21839',
        severity: 'critical',
        target: 'Oracle WebLogic Server',
        description: 'WebLogic Server存在远程代码执行漏洞，通过IIOP协议可实现未授权访问和RCE。',
        payload: 't3 protocol payload',
        impact: '服务器完全控制，数据库拖库',
        remediation: '安装Oracle最新安全补丁，禁用IIOP协议',
        references: ['https://www.oracle.com/security-alerts/cpujan2023.html'],
        created_at: '2024-02-10',
      },
      {
        id: 'poc-006',
        name: 'Confluence OGNL注入',
        cve_id: 'CVE-2022-26134',
        severity: 'critical',
        target: 'Atlassian Confluence',
        description: 'Confluence Server和Data Center存在OGNL注入漏洞，攻击者可实现无需认证的远程代码执行。',
        payload: '/%24%7B%40java.lang.Runtime%40getRuntime%28%29.exec%28%22whoami%22%29%7D',
        impact: '服务器完全控制，敏感文档泄露',
        remediation: '升级到Confluence 7.19.17+或8.0+，部署WAF防护',
        references: ['https://confluence.atlassian.com/doc/confluence-security-advisory-2022-06-02-1101413739.html'],
        created_at: '2024-02-15',
      },
      {
        id: 'poc-007',
        name: 'F5 BIG-IP RCE',
        cve_id: 'CVE-2022-1388',
        severity: 'critical',
        target: 'F5 BIG-IP',
        description: 'F5 BIG-IP iControl REST接口存在未授权远程代码执行漏洞。',
        payload: 'POST /mgmt/shared/authn/login',
        impact: '网络设备控制，流量劫持',
        remediation: '升级到BIG-IP 16.1.2.2+或17.0.0+，限制管理接口访问',
        references: ['https://support.f5.com/csp/article/K23605346'],
        created_at: '2024-02-20',
      },
      {
        id: 'poc-008',
        name: 'VMware vCenter RCE',
        cve_id: 'CVE-2021-21972',
        severity: 'critical',
        target: 'VMware vCenter Server',
        description: 'vCenter Server插件存在远程代码执行漏洞，攻击者可通过未授权访问实现RCE。',
        payload: 'POST /ui/vropspluginui/rest/services/uploadova',
        impact: '虚拟化基础设施控制，虚拟机逃逸',
        remediation: '安装VMware安全补丁，禁用插件',
        references: ['https://www.vmware.com/security/advisories/VMSA-2021-0010.html'],
        created_at: '2024-02-25',
      },
      {
        id: 'poc-009',
        name: 'Jenkins RCE',
        cve_id: 'CVE-2024-23897',
        severity: 'high',
        target: 'Jenkins',
        description: 'Jenkins CLI存在任意文件读取漏洞，可导致远程代码执行。',
        payload: 'java -jar jenkins-cli.jar -s http://target:8080 who-am-i',
        impact: '凭据泄露，RCE',
        remediation: '升级到Jenkins 2.442+或LTS 2.426.3+',
        references: ['https://www.jenkins.io/security/advisory/2024-01-24/'],
        created_at: '2024-03-01',
      },
      {
        id: 'poc-010',
        name: 'GitLab RCE',
        cve_id: 'CVE-2023-2825',
        severity: 'critical',
        target: 'GitLab',
        description: 'GitLab社区版和企业版存在任意文件读取漏洞，可导致RCE。',
        payload: 'uploads/user 或 attachments/../../../etc/passwd',
        impact: '敏感文件读取，服务器沦陷',
        remediation: '升级到GitLab 16.0.1+或15.11.8+',
        references: ['https://about.gitlab.com/releases/2023/05/22/critical-security-release-gitlab-16-0-1-released/'],
        created_at: '2024-03-05',
      },
    ];
    setPocs(expertPOCs);
  };

  const loadActiveTasks = async () => {
    try {
      const tasks = await activeScanApi.listTasks() as ActiveScanTask[];
      setActiveTasks(Array.isArray(tasks) ? tasks.slice(0, 10) : []);
    } catch (error) {
      console.error('加载扫描任务失败:', error);
    }
  };

  const executePOC = async (poc: POC) => {
    setPocLoading(true);
    try {
      const task = await activeScanApi.createTask({
        name: `POC扫描 - ${poc.name}`,
        target: poc.target,
        type: 'vuln',
      }) as ActiveScanTask;

      setTimeout(async () => {
        try {
          const results = await activeScanApi.getResults(task.id) as ScanResultItem[];
          setScanResults(prev => ({ ...prev, [task.id]: results }));
        } catch (e) {
          console.error('获取结果失败:', e);
        }
      }, 3000);

      alert(`POC扫描任务已启动!\n目标: ${poc.target}\nCVE: ${poc.cve_id}\n请稍后查看扫描结果`);
    } catch (error: any) {
      alert(`执行失败: ${error.message}`);
    } finally {
      setPocLoading(false);
    }
  };

  const handlePractice = async () => {
    if (!selectedTechnique || !practiceURL) {
      alert('请选择漏洞技巧并输入目标URL');
      return;
    }

    setIsPracticing(true);
    setPracticeLoading(true);

    try {
      const payload = practicePayload || selectedTechnique.payload;
      const result = await vulnTechniquesApi.practice({
        technique_id: selectedTechnique.id,
        target: practiceTarget,
        url: practiceURL,
        payload: payload,
        mode: practiceMode,
      });

      if (practiceMode === 'sync') {
        alert(`同步执行完成!\n任务ID: ${(result as any).task_id}`);
      } else {
        alert(`异步任务已启动!\n任务ID: ${(result as any).task_id}\n请稍后查看结果`);
      }

      setTimeout(() => loadPracticeResults(), 1000);
    } catch (error: any) {
      alert(`执行失败: ${error.message}`);
    } finally {
      setIsPracticing(false);
      setPracticeLoading(false);
    }
  };

  const handleUniversityScan = async (uni: University) => {
    try {
      const task = await activeScanApi.createTask({
        name: `高校扫描 - ${uni.name}`,
        target: uni.domain,
        type: 'vuln',
      }) as ActiveScanTask;

      alert(`扫描任务已启动!\n高校: ${uni.name}\n域名: ${uni.domain}\n请稍后查看结果`);
      loadActiveTasks();
    } catch (error: any) {
      alert(`扫描启动失败: ${error.message}`);
    }
  };

  const getSeverityColor = (severity: string) => {
    switch (severity) {
      case 'critical': return 'bg-red-600 text-white';
      case 'high': return 'bg-red-500/20 text-red-400 border border-red-500/30';
      case 'medium': return 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30';
      case 'low': return 'bg-blue-500/20 text-blue-400 border border-blue-500/30';
      default: return 'bg-gray-500/20 text-gray-400 border border-gray-500/30';
    }
  };

  const getDifficultyColor = (difficulty: string) => {
    switch (difficulty) {
      case 'beginner': return 'bg-green-500/20 text-green-400';
      case 'intermediate': return 'bg-yellow-500/20 text-yellow-400';
      case 'advanced': return 'bg-red-500/20 text-red-400';
      default: return 'bg-gray-500/20 text-gray-400';
    }
  };

  const filteredPOCs = pocs.filter(poc => {
    if (pocFilter.severity && poc.severity !== pocFilter.severity) return false;
    if (pocFilter.search) {
      const search = pocFilter.search.toLowerCase();
      return poc.name.toLowerCase().includes(search) ||
             poc.cve_id.toLowerCase().includes(search) ||
             poc.target.toLowerCase().includes(search);
    }
    return true;
  });

  return (
    <div className="h-full flex flex-col bg-gray-950">
      <div className="p-4 bg-gray-900 border-b border-gray-800">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-white">教育SRC</h2>
            <p className="text-sm text-gray-400">高校信息收集与漏洞挖掘 · 专家级POC库</p>
          </div>
        </div>
      </div>

      <div className="border-b border-gray-800 bg-gray-900">
        <nav className="flex px-4">
          {[
            { key: 'overview' as const, label: '🔓 专项POC库' },
            { key: 'poc' as const, label: '🏛️ SRC平台' },
            { key: 'src' as const, label: '📚 挖洞技巧' },
            { key: 'batch' as const, label: '🚀 批量扫描' },
          ].map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`py-4 px-4 border-b-2 font-medium text-sm transition-colors ${
                activeTab === tab.key
                  ? 'border-blue-500 text-blue-400'
                  : 'border-transparent text-gray-500 hover:text-gray-300'
              }`}
            >
              {tab.label}
            </button>
          ))}
        </nav>
      </div>

      <div className="flex-1 overflow-auto p-6">
        {activeTab === 'overview' && (
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-400">POC总数</p>
                    <p className="text-3xl font-bold text-white mt-1">{pocs.length}</p>
                  </div>
                  <div className="w-12 h-12 bg-red-500/20 rounded-lg flex items-center justify-center">
                    <span className="text-2xl">🛡️</span>
                  </div>
                </div>
              </div>
              <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-400">严重漏洞</p>
                    <p className="text-3xl font-bold text-red-500 mt-1">
                      {pocs.filter(p => p.severity === 'critical' || p.severity === 'high').length}
                    </p>
                  </div>
                  <div className="w-12 h-12 bg-red-500/20 rounded-lg flex items-center justify-center">
                    <span className="text-2xl">⚠️</span>
                  </div>
                </div>
              </div>
              <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-400">扫描任务</p>
                    <p className="text-3xl font-bold text-blue-400 mt-1">{activeTasks.length}</p>
                  </div>
                  <div className="w-12 h-12 bg-blue-500/20 rounded-lg flex items-center justify-center">
                    <span className="text-2xl">🔍</span>
                  </div>
                </div>
              </div>
              <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm text-gray-400">SRC平台</p>
                    <p className="text-3xl font-bold text-purple-500 mt-1">{srcPlatforms.length}</p>
                  </div>
                  <div className="w-12 h-12 bg-purple-500/20 rounded-lg flex items-center justify-center">
                    <span className="text-2xl">🏛️</span>
                  </div>
                </div>
              </div>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
                <h3 className="text-lg font-semibold text-white mb-4">🔥 热点POC</h3>
                <div className="space-y-3">
                  {pocs.slice(0, 5).map((poc) => (
                    <div key={poc.id} className="p-3 bg-gray-800 rounded-lg hover:bg-gray-750 transition cursor-pointer"
                         onClick={() => setSelectedPOC(poc)}>
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-2">
                          <span className={`px-2 py-0.5 rounded text-xs ${getSeverityColor(poc.severity)}`}>
                            {poc.severity === 'critical' ? '严重' : poc.severity === 'high' ? '高危' : '中危'}
                          </span>
                          <span className="text-white font-medium text-sm">{poc.name}</span>
                        </div>
                        <span className="text-xs text-gray-500">{poc.cve_id}</span>
                      </div>
                      <p className="text-xs text-gray-400 mt-1">影响: {poc.target}</p>
                    </div>
                  ))}
                </div>
              </div>

              <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
                <h3 className="text-lg font-semibold text-white mb-4">📊 最近扫描任务</h3>
                {activeTasks.length === 0 ? (
                  <div className="text-center py-8 text-gray-500">
                    <div className="text-4xl mb-4">🔍</div>
                    <p>暂无扫描任务</p>
                  </div>
                ) : (
                  <div className="space-y-3">
                    {activeTasks.slice(0, 5).map((task) => (
                      <div key={task.id} className="p-3 bg-gray-800 rounded-lg">
                        <div className="flex items-center justify-between mb-1">
                          <span className="text-white font-medium text-sm">{task.name}</span>
                          <span className={`px-2 py-0.5 rounded text-xs ${
                            task.status === 'completed' ? 'bg-green-600/20 text-green-400' :
                            task.status === 'running' ? 'bg-yellow-600/20 text-yellow-400' :
                            'bg-gray-600/20 text-gray-400'
                          }`}>
                            {task.status === 'completed' ? '完成' :
                             task.status === 'running' ? '运行中' : '等待'}
                          </span>
                        </div>
                        <div className="w-full bg-gray-700 rounded-full h-1.5 mt-2">
                          <div
                            className="bg-blue-500 h-1.5 rounded-full transition-all"
                            style={{ width: `${task.progress}%` }}
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>

            {selectedPOC && (
              <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
                <div className="flex items-center justify-between mb-4">
                  <h3 className="text-lg font-semibold text-white">{selectedPOC.name}</h3>
                  <button onClick={() => setSelectedPOC(null)} className="text-gray-400 hover:text-white">✕</button>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-3">
                    <div className="flex items-center gap-2">
                      <span className={`px-2 py-0.5 rounded text-xs ${getSeverityColor(selectedPOC.severity)}`}>
                        {selectedPOC.severity === 'critical' ? '严重' : '高危'}
                      </span>
                      <span className="text-sm text-gray-400">{selectedPOC.cve_id}</span>
                    </div>
                    <div>
                      <p className="text-xs text-gray-400 mb-1">影响目标</p>
                      <p className="text-sm text-white">{selectedPOC.target}</p>
                    </div>
                    <div>
                      <p className="text-xs text-gray-400 mb-1">漏洞描述</p>
                      <p className="text-sm text-gray-300">{selectedPOC.description}</p>
                    </div>
                    <div>
                      <p className="text-xs text-gray-400 mb-1">安全影响</p>
                      <p className="text-sm text-red-400">{selectedPOC.impact}</p>
                    </div>
                  </div>
                  <div className="space-y-3">
                    <div>
                      <p className="text-xs text-gray-400 mb-1">修复方案</p>
                      <p className="text-sm text-green-400">{selectedPOC.remediation}</p>
                    </div>
                    <div>
                      <p className="text-xs text-gray-400 mb-1">POC Payload</p>
                      <pre className="text-xs text-red-400 bg-gray-800 p-2 rounded overflow-x-auto">
                        {selectedPOC.payload}
                      </pre>
                    </div>
                    <button
                      onClick={() => executePOC(selectedPOC)}
                      disabled={pocLoading}
                      className="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-500 transition disabled:opacity-50"
                    >
                      {pocLoading ? '执行中...' : '🚀 执行POC扫描'}
                    </button>
                  </div>
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === 'poc' && (
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-semibold text-white">🏛️ SRC漏洞响应平台</h3>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {srcPlatforms.map((platform) => (
                <div key={platform.id} className="bg-gray-900 rounded-xl border border-gray-800 p-5 hover:border-blue-500/50 transition">
                  <div className="flex items-start justify-between">
                    <div>
                      <h4 className="font-semibold text-white">{platform.name}</h4>
                      <span className="px-2 py-0.5 bg-blue-600/20 text-blue-400 rounded text-xs mt-1 inline-block">
                        {platform.type}
                      </span>
                    </div>
                    <span className="px-2 py-0.5 bg-purple-600/20 text-purple-400 rounded text-xs">
                      {platform.rank}
                    </span>
                  </div>
                  <p className="text-sm text-gray-400 mt-3">{platform.url}</p>
                  <div className="flex items-center justify-between mt-3">
                    <span className="text-xs text-green-400">{platform.bounty_range}</span>
                    <a
                      href={platform.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="px-3 py-1 bg-blue-600/20 text-blue-400 rounded text-xs hover:bg-blue-600/30 transition"
                    >
                      访问平台 →
                    </a>
                  </div>
                </div>
              ))}
            </div>

            <div className="bg-gray-900 rounded-xl border border-gray-800 p-5 mt-6">
              <h3 className="text-lg font-semibold text-white mb-4">💡 SRC挖掘建议</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="p-4 bg-gray-800 rounded-lg">
                  <h4 className="font-medium text-white mb-2">🎯 高校SRC挖洞技巧</h4>
                  <ul className="text-sm text-gray-400 space-y-1">
                    <li>• 关注教务系统、选课系统、成绩管理系统</li>
                    <li>• 测试VPN、统一身份认证(CAS)系统</li>
                    <li>• 检查邮件系统和FTP服务的弱口令</li>
                    <li>• 关注校企合作平台的安全问题</li>
                  </ul>
                </div>
                <div className="p-4 bg-gray-800 rounded-lg">
                  <h4 className="font-medium text-white mb-2">🛡️ 漏洞评级参考</h4>
                  <ul className="text-sm text-gray-400 space-y-1">
                    <li>• 严重：Getshell、RCE、SQL注入大量数据</li>
                    <li>• 高危：敏感信息泄露、任意文件读取</li>
                    <li>• 中危：存储型XSS、CSRF、弱口令</li>
                    <li>• 低危：反射型XSS、路径遍历、信息收集</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'src' && (
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <h3 className="text-lg font-semibold text-white">漏洞挖掘技巧库</h3>
              <button
                onClick={() => { loadTechniques(); loadPracticeResults(); }}
                className="px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition"
              >
                刷新数据
              </button>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              <div className="lg:col-span-2 space-y-4">
                <div className="bg-gray-900 rounded-xl border border-gray-800 p-4">
                  <div className="flex flex-wrap gap-3 mb-4">
                    <select
                      value={filterCategory}
                      onChange={(e) => setFilterCategory(e.target.value)}
                      className="px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm"
                    >
                      <option value="">全部分类</option>
                      <option value="注入类">注入类</option>
                      <option value="跨站类">跨站类</option>
                      <option value="认证类">认证类</option>
                      <option value="文件处理">文件处理</option>
                      <option value="信息收集">信息收集</option>
                      <option value="权限类">权限类</option>
                      <option value="命令注入">命令注入</option>
                    </select>
                    <select
                      value={filterSeverity}
                      onChange={(e) => setFilterSeverity(e.target.value)}
                      className="px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm"
                    >
                      <option value="">全部严重性</option>
                      <option value="critical">严重</option>
                      <option value="high">高危</option>
                      <option value="medium">中危</option>
                      <option value="low">低危</option>
                    </select>
                    <input
                      type="text"
                      placeholder="搜索关键词..."
                      value={filterKeyword}
                      onChange={(e) => setFilterKeyword(e.target.value)}
                      className="px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm flex-1"
                    />
                    <button
                      onClick={loadTechniques}
                      className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500 transition"
                    >
                      搜索
                    </button>
                  </div>

                  {techniques.length === 0 ? (
                    <div className="text-center py-8 text-gray-500">
                      <div className="text-4xl mb-4">📚</div>
                      <p>暂无漏洞技巧数据</p>
                    </div>
                  ) : (
                    <div className="space-y-3 max-h-96 overflow-y-auto">
                      {techniques.map((tech) => (
                        <div
                          key={tech.id}
                          onClick={() => setSelectedTechnique(tech)}
                          className={`p-4 rounded-lg border cursor-pointer transition ${
                            selectedTechnique?.id === tech.id
                              ? 'bg-blue-600/20 border-blue-500'
                              : 'bg-gray-800 border-gray-700 hover:border-gray-600'
                          }`}
                        >
                          <div className="flex items-start justify-between">
                            <div className="flex-1">
                              <h4 className="font-medium text-white">{tech.title}</h4>
                              <p className="text-sm text-gray-400 mt-1 line-clamp-2">{tech.summary}</p>
                              <div className="flex flex-wrap gap-2 mt-2">
                                <span className={`px-2 py-0.5 rounded text-xs ${getSeverityColor(tech.severity)}`}>
                                  {tech.severity === 'critical' ? '严重' : tech.severity === 'high' ? '高危' : tech.severity === 'medium' ? '中危' : '低危'}
                                </span>
                                <span className={`px-2 py-0.5 rounded text-xs ${getDifficultyColor(tech.difficulty)}`}>
                                  {tech.difficulty === 'beginner' ? '入门' : tech.difficulty === 'intermediate' ? '进阶' : '高级'}
                                </span>
                              </div>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>

                {selectedTechnique && (
                  <div className="bg-gray-900 rounded-xl border border-gray-800 p-4">
                    <h4 className="font-medium text-white mb-3">实践结果</h4>
                    {practiceResults.filter(r => r.technique_id === selectedTechnique.id).length === 0 ? (
                      <p className="text-gray-500 text-sm">暂无实践记录</p>
                    ) : (
                      <div className="space-y-2 max-h-48 overflow-y-auto">
                        {practiceResults
                          .filter(r => r.technique_id === selectedTechnique.id)
                          .slice(0, 5)
                          .map((result) => (
                            <div key={result.id} className="p-3 bg-gray-800 rounded-lg">
                              <div className="flex items-center justify-between mb-1">
                                <span className={`px-2 py-0.5 rounded text-xs ${
                                  result.status === 'success' ? 'bg-green-600/20 text-green-400' :
                                  result.status === 'error' ? 'bg-red-600/20 text-red-400' :
                                  'bg-gray-600/20 text-gray-400'
                                }`}>
                                  {result.status === 'success' ? '成功' : result.status === 'error' ? '错误' : '执行中'}
                                </span>
                                <span className="text-xs text-gray-500">
                                  {new Date(result.created_at).toLocaleString('zh-CN')}
                                </span>
                              </div>
                              <p className="text-sm text-gray-300 truncate">{result.url}</p>
                            </div>
                          ))}
                      </div>
                    )}
                  </div>
                )}
              </div>

              <div className="space-y-4">
                {selectedTechnique ? (
                  <>
                    <div className="bg-gray-900 rounded-xl border border-gray-800 p-4">
                      <h4 className="font-medium text-white mb-3">{selectedTechnique.title}</h4>
                      <div className="space-y-3">
                        <div>
                          <p className="text-xs text-gray-400 mb-1">描述</p>
                          <p className="text-sm text-gray-300">{selectedTechnique.description}</p>
                        </div>
                        <div>
                          <p className="text-xs text-gray-400 mb-1">影响范围</p>
                          <p className="text-sm text-gray-300">{selectedTechnique.impact}</p>
                        </div>
                        <div>
                          <p className="text-xs text-gray-400 mb-1">Payload 示例</p>
                          <pre className="text-xs text-red-400 bg-gray-800 p-2 rounded overflow-x-auto whitespace-pre-wrap">
                            {selectedTechnique.payload}
                          </pre>
                        </div>
                        <div>
                          <p className="text-xs text-gray-400 mb-1">修复建议</p>
                          <p className="text-sm text-gray-300 whitespace-pre-wrap">{selectedTechnique.fix_suggestion}</p>
                        </div>
                      </div>
                    </div>

                    <div className="bg-gray-900 rounded-xl border border-gray-800 p-4">
                      <h4 className="font-medium text-white mb-3">实践测试</h4>
                      <div className="space-y-3">
                        <div>
                          <p className="text-xs text-gray-400 mb-1">执行模式</p>
                          <div className="flex gap-2">
                            <button
                              onClick={() => setPracticeMode('sync')}
                              className={`flex-1 px-3 py-2 rounded-lg text-sm transition ${
                                practiceMode === 'sync' ? 'bg-blue-600 text-white' : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                              }`}
                            >
                              同步
                            </button>
                            <button
                              onClick={() => setPracticeMode('async')}
                              className={`flex-1 px-3 py-2 rounded-lg text-sm transition ${
                                practiceMode === 'async' ? 'bg-blue-600 text-white' : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                              }`}
                            >
                              异步
                            </button>
                          </div>
                        </div>
                        <div>
                          <p className="text-xs text-gray-400 mb-1">目标URL *</p>
                          <input
                            type="text"
                            value={practiceURL}
                            onChange={(e) => setPracticeURL(e.target.value)}
                            placeholder="https://target.com"
                            className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm"
                          />
                        </div>
                        <button
                          onClick={handlePractice}
                          disabled={isPracticing || !practiceURL}
                          className={`w-full px-4 py-2 rounded-lg text-white transition ${
                            isPracticing || !practiceURL ? 'bg-gray-600 cursor-not-allowed' : 'bg-green-600 hover:bg-green-500'
                          }`}
                        >
                          {isPracticing ? '执行中...' : '开始实践'}
                        </button>
                      </div>
                    </div>
                  </>
                ) : (
                  <div className="bg-gray-900 rounded-xl border border-gray-800 p-8 text-center">
                    <div className="text-4xl mb-4">📚</div>
                    <p className="text-gray-400">请从左侧选择一个漏洞技巧</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}

        {activeTab === 'batch' && (
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <p className="text-gray-400">共 {universities.length} 所高校</p>
              <button
                onClick={() => loadData()}
                className="px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition"
              >
                刷新数据
              </button>
            </div>
            {loading ? (
              <div className="text-center py-12 text-gray-500">
                <div className="animate-spin text-4xl mb-4">⏳</div>
                <p>加载中...</p>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {universities.map((uni) => (
                  <div key={uni.id} className="bg-gray-900 rounded-lg border border-gray-800 p-4">
                    <div className="flex items-start justify-between">
                      <div>
                        <h3 className="font-semibold text-white">{uni.name}</h3>
                        <p className="text-sm text-gray-400">{uni.province}</p>
                      </div>
                    </div>
                    <div className="mt-3 flex items-center gap-2">
                      <span className="px-2 py-1 bg-blue-600/20 text-blue-400 rounded text-xs">
                        {uni.level}
                      </span>
                      <span className="px-2 py-1 bg-gray-600/20 text-gray-400 rounded text-xs">
                        {uni.type}
                      </span>
                    </div>
                    <div className="mt-3 flex gap-2">
                      <button className="flex-1 px-3 py-1.5 text-sm bg-blue-600/20 text-blue-400 rounded hover:bg-blue-600/30 transition">
                        查看详情
                      </button>
                      <button
                        onClick={() => handleUniversityScan(uni)}
                        className="flex-1 px-3 py-1.5 text-sm bg-green-600/20 text-green-400 rounded hover:bg-green-600/30 transition"
                      >
                        开始扫描
                      </button>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default EduPanel;
