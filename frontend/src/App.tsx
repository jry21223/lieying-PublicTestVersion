import { useState, useEffect } from "react";
import { 
  Target, 
  Search, 
  ShieldAlert, 
  Database, 
  Bot, 
  Settings,
  Play,
  Scan,
  ChevronRight,
  Globe,
  Activity,
  FileText,
  GraduationCap,
  Loader2,
  Filter,
  Download,
  MessageSquare,
  Cpu,
  BarChart3,
  List,
  Wifi,
  Sun,
  Moon,
  Target as TargetIcon
} from "lucide-react";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import EduPanel from "./components/EduPanel";
import ReconPanel from "./components/ReconPanel";
import ScanPanel from "./components/ScanPanel";
import AssetOverview from "./components/AssetOverview";
import ReconOverview from "./components/ReconOverview";
import VulnScanPanel from "./components/VulnScanPanel";
import VulnResultsPanel from "./components/VulnResultsPanel";
import VulnDashboard from "./components/VulnDashboard";
import VulnManagement from "./components/VulnManagement";
import AIChat from "./components/AIChat";
import AIConfigPanel from "./components/AIConfig";
import EduUniversityManager from "./components/EduUniversityManager";
import EduSRCManager from "./components/EduSRCManager";
import ReportGenerator from "./components/ReportGenerator";
import NetworkSettings from "./components/NetworkSettings";
import ActiveScan from "./components/ActiveScan";
import RequestHistory from "./components/RequestHistory";
import ProxyPool from "./components/ProxyPool";
import api from "./services/api";

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

const severityConfig = {
  critical: { color: "bg-red-600", bg: "bg-red-600/10", border: "border-red-600", label: "严重" },
  high: { color: "bg-orange-600", bg: "bg-orange-600/10", border: "border-orange-600", label: "高危" },
  medium: { color: "bg-yellow-600", bg: "bg-yellow-600/10", border: "border-yellow-600", label: "中危" },
  low: { color: "bg-blue-600", bg: "bg-blue-600/10", border: "border-blue-600", label: "低危" },
  info: { color: "bg-gray-600", bg: "bg-gray-600/10", border: "border-gray-600", label: "信息" },
};

const SidebarItem = ({ 
  icon: Icon, 
  label, 
  active, 
  onClick,
  theme
}: { 
  icon: any; 
  label: string; 
  active?: boolean; 
  onClick: () => void;
  theme?: 'dark' | 'light';
}) => (
  <button
    onClick={onClick}
    className={cn(
      "flex items-center w-full gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-200",
      active 
        ? "bg-blue-600 text-white shadow-lg shadow-blue-900/20"
        : theme === 'dark'
          ? "text-gray-400 hover:text-white hover:bg-gray-800"
          : "text-gray-600 hover:text-gray-900 hover:bg-gray-200"
    )}
  >
    <Icon className="w-5 h-5" />
    <span>{label}</span>
  </button>
);

const AssetTreeNode = ({ 
  name, 
  type, 
  vulnCount, 
  isOpen, 
  onToggle 
}: { 
  name: string; 
  type: 'domain' | 'ip' | 'url';
  vulnCount?: number;
  isOpen?: boolean;
  onToggle: () => void;
}) => (
  <div 
    className="group flex items-center justify-between px-2 py-1.5 text-sm hover:bg-gray-800 rounded cursor-pointer"
    onClick={onToggle}
  >
    <div className="flex items-center gap-2">
      <ChevronRight className={cn("w-3 h-3 text-gray-500 transition-transform", isOpen && "rotate-90")} />
      {type === 'domain' ? <Globe className="w-4 h-4 text-blue-400" /> : 
       type === 'ip' ? <Activity className="w-4 h-4 text-green-400" /> :
       <FileText className="w-4 h-4 text-purple-400" />}
      <span className="text-gray-300">{name}</span>
    </div>
    {vulnCount && vulnCount > 0 && (
      <span className="bg-red-600 text-white text-[10px] px-1.5 py-0.5 rounded-full font-bold">
        {vulnCount}
      </span>
    )}
  </div>
);

const VulnerabilityCard = ({ 
  title, 
  severity, 
  url, 
  confirmed, 
  onVerify 
}: { 
  title: string; 
  severity: 'critical' | 'high' | 'medium' | 'low' | 'info';
  url: string;
  confirmed: boolean;
  onVerify: () => void;
}) => {
  const severityColors = {
    critical: "border-l-red-500 bg-red-500/10",
    high: "border-l-orange-500 bg-orange-500/10",
    medium: "border-l-yellow-500 bg-yellow-500/10",
    low: "border-l-blue-500 bg-blue-500/10",
    info: "border-l-gray-500 bg-gray-500/10",
  };

  const severityLabels = {
    critical: "严重",
    high: "高危",
    medium: "中危",
    low: "低危",
    info: "信息",
  };

  return (
    <div className={cn("p-4 rounded-lg border-l-4 bg-gray-800/50 mb-3", severityColors[severity])}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <h4 className="font-semibold text-gray-100">{title}</h4>
            <span className={cn(
              "text-[10px] px-2 py-0.5 rounded-full font-bold",
              severity === 'critical' ? "bg-red-600 text-white" :
              severity === 'high' ? "bg-orange-600 text-white" :
              severity === 'medium' ? "bg-yellow-600 text-black" :
              "bg-gray-600 text-white"
            )}>
              {severityLabels[severity]}
            </span>
            {confirmed && (
              <span className="text-green-500 text-xs flex items-center gap-1">
                <ShieldAlert className="w-3 h-3" /> 已确认
              </span>
            )}
          </div>
          <p className="text-xs text-gray-400 font-mono truncate">{url}</p>
        </div>
        <button 
          onClick={onVerify}
          className="ml-4 px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white text-xs rounded-md transition-colors flex items-center gap-1"
        >
          <Bot className="w-3.5 h-3.5" />
          AI验证
        </button>
      </div>
    </div>
  );
};

// AI助手组件
const AIAssistant = () => {
  const [messages, setMessages] = useState<{role: 'user' | 'ai', content: string}[]>([
    { role: 'ai', content: '你好！我是「小影」，你的AI安全助手。我可以帮你：\n\n• 分析漏洞并生成POC\n• 生成SRC报告\n• 推荐攻击路径\n• WAF绕过辅助\n\n有什么我可以帮你的吗？' }
  ]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSend = async () => {
    if (!input.trim()) return;
    
    const userMessage = input.trim();
    setMessages(prev => [...prev, { role: 'user', content: userMessage }]);
    setInput('');
    setLoading(true);

    // 模拟AI回复
    setTimeout(() => {
      setMessages(prev => [...prev, { 
        role: 'ai', 
        content: `我收到了你的问题："${userMessage}"\n\n作为AI安全助手，我可以帮助你分析这个问题。不过目前AI功能需要连接到Ollama服务才能提供完整的智能分析。\n\n你可以尝试：\n1. 启动本地Ollama服务\n2. 使用POC生成功能\n3. 查看漏洞分析建议` 
      }]);
      setLoading(false);
    }, 1500);
  };

  return (
    <div className="max-w-3xl mx-auto h-full flex flex-col">
      <div className="flex-1 bg-gray-900 border border-gray-800 rounded-xl p-6 mb-4 overflow-y-auto">
        {messages.map((msg, index) => (
          <div key={index} className={`flex items-start gap-4 mb-6 ${msg.role === 'user' ? 'flex-row-reverse' : ''}`}>
            <div className={`w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 ${
              msg.role === 'ai' 
                ? 'bg-gradient-to-br from-blue-500 to-purple-600' 
                : 'bg-gray-700'
            }`}>
              {msg.role === 'ai' ? <Bot className="w-5 h-5 text-white" /> : <span className="text-white text-sm">我</span>}
            </div>
            <div className={`rounded-2xl p-4 text-sm max-w-[80%] whitespace-pre-wrap ${
              msg.role === 'ai' 
                ? 'bg-gray-800 rounded-tl-none text-gray-200' 
                : 'bg-blue-600 rounded-tr-none text-white'
            }`}>
              {msg.content}
            </div>
          </div>
        ))}
        {loading && (
          <div className="flex items-start gap-4 mb-6">
            <div className="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center flex-shrink-0">
              <Bot className="w-5 h-5 text-white" />
            </div>
            <div className="bg-gray-800 rounded-2xl rounded-tl-none p-4 text-sm text-gray-200">
              <Loader2 className="w-4 h-4 animate-spin" />
            </div>
          </div>
        )}
      </div>
      <div className="bg-gray-900 border border-gray-800 rounded-xl p-3 flex gap-3">
        <input 
          type="text" 
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && handleSend()}
          placeholder="输入你的问题..."
          className="flex-1 bg-transparent border-none text-gray-100 placeholder-gray-500 focus:ring-0 text-sm"
        />
        <button 
          onClick={handleSend}
          disabled={loading}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 text-white rounded-lg text-sm font-medium transition-colors"
        >
          发送
        </button>
      </div>
    </div>
  );
};

// 漏洞管理组件
const VulnerabilityManager = () => {
  const [vulnerabilities, setVulnerabilities] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('all');
  const [selectedVuln, setSelectedVuln] = useState<any>(null);

  useEffect(() => {
    loadVulnerabilities();
  }, []);

  const loadVulnerabilities = async () => {
    try {
      const data = await api.vulnerabilities.list();
      setVulnerabilities(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('加载漏洞失败:', error);
      setVulnerabilities([]);
    } finally {
      setLoading(false);
    }
  };

  const handleVerify = async (id: string) => {
    try {
      await api.vulnerabilities.confirm(id);
      loadVulnerabilities();
    } catch (error) {
      console.error('确认漏洞失败:', error);
    }
  };

  const handleExportVulns = () => {
    const dataStr = JSON.stringify(vulnerabilities, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `vulnerabilities-${new Date().toISOString().split('T')[0]}.json`;
    link.click();
  };

  const filteredVulns = filter === 'all' ? vulnerabilities : vulnerabilities.filter(v => v.severity === filter);

  const severityCount = {
    critical: vulnerabilities.filter(v => v.severity === 'critical').length,
    high: vulnerabilities.filter(v => v.severity === 'high').length,
    medium: vulnerabilities.filter(v => v.severity === 'medium').length,
    low: vulnerabilities.filter(v => v.severity === 'low').length,
  };

  if (loading) {
    return <div className="flex items-center justify-center h-full"><Loader2 className="w-8 h-8 animate-spin" /></div>;
  }

  return (
    <div className="h-full flex flex-col">
      <div className="flex items-center justify-between mb-6">
        <h3 className="text-lg font-semibold text-white">漏洞列表</h3>
        <div className="flex gap-2">
          <button 
            onClick={handleExportVulns}
            className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white text-sm rounded-lg transition-colors"
          >
            导出
          </button>
          <button 
            onClick={loadVulnerabilities}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-sm rounded-lg transition-colors"
          >
            刷新
          </button>
        </div>
      </div>

      {vulnerabilities.length > 0 && (
        <div className="grid grid-cols-4 gap-4 mb-6">
          <div className="bg-red-900/20 border border-red-800/30 rounded-xl p-4">
            <p className="text-sm text-gray-400">严重</p>
            <p className="text-2xl font-bold text-red-500">{severityCount.critical}</p>
          </div>
          <div className="bg-orange-900/20 border border-orange-800/30 rounded-xl p-4">
            <p className="text-sm text-gray-400">高危</p>
            <p className="text-2xl font-bold text-orange-500">{severityCount.high}</p>
          </div>
          <div className="bg-yellow-900/20 border border-yellow-800/30 rounded-xl p-4">
            <p className="text-sm text-gray-400">中危</p>
            <p className="text-2xl font-bold text-yellow-500">{severityCount.medium}</p>
          </div>
          <div className="bg-blue-900/20 border border-blue-800/30 rounded-xl p-4">
            <p className="text-sm text-gray-400">低危</p>
            <p className="text-2xl font-bold text-blue-500">{severityCount.low}</p>
          </div>
        </div>
      )}

      {vulnerabilities.length > 0 && (
        <div className="flex items-center gap-2 mb-4 overflow-x-auto pb-2">
          <Filter className="w-4 h-4 text-gray-500 flex-shrink-0" />
          {[
            { key: 'all', label: '全部' },
            { key: 'critical', label: '严重' },
            { key: 'high', label: '高危' },
            { key: 'medium', label: '中危' },
            { key: 'low', label: '低危' },
          ].map(({ key, label }) => (
            <button
              key={key}
              onClick={() => setFilter(key)}
              className={cn(
                "px-3 py-1.5 rounded-lg text-sm whitespace-nowrap transition-colors",
                filter === key
                  ? "bg-blue-600 text-white"
                  : "bg-gray-800 text-gray-400 hover:bg-gray-700"
              )}
            >
              {label}
            </button>
          ))}
        </div>
      )}

      <div className="flex-1 overflow-y-auto">
        {filteredVulns.length > 0 ? (
          <div className="space-y-3">
            {filteredVulns.map((vuln) => {
              const config = severityConfig[vuln.severity as keyof typeof severityConfig];
              return (
                <VulnerabilityCard 
                  key={vuln.id}
                  title={vuln.title}
                  severity={vuln.severity}
                  url={vuln.url}
                  confirmed={vuln.confirmed}
                  onVerify={() => handleVerify(vuln.id)}
                />
              );
            })}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-full text-gray-500">
            <ShieldAlert className="w-16 h-16 mb-4 opacity-30" />
            <p className="text-lg">暂无漏洞数据</p>
            <p className="text-sm mt-2">开始扫描以发现漏洞</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default function App() {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [selectedDomain, setSelectedDomain] = useState<string | null>(null);
  const [stats, setStats] = useState({ targets: 0, vulns: 0, confirmed: 0 });
  const [targets, setTargets] = useState<any[]>([]);
  const [vulnerabilities, setVulnerabilities] = useState<any[]>([]);
  const [refreshKey, setRefreshKey] = useState(0);
  const [showAddModal, setShowAddModal] = useState(false);
  const [showImportModal, setShowImportModal] = useState(false);
  const [newTarget, setNewTarget] = useState({ name: '', value: '', type: 'domain', description: '', tags: '' });
  const [importText, setImportText] = useState('');
  const [importType, setImportType] = useState<'domain' | 'ip' | 'url'>('domain');
  const [scanningTargets, setScanningTargets] = useState<string[]>([]);
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [sidebarWidth, setSidebarWidth] = useState(288);
  const [isResizing, setIsResizing] = useState(false);
  const [startX, setStartX] = useState(0);
  const [startWidth, setStartWidth] = useState(288);
  const [theme, setTheme] = useState<'dark' | 'light'>('dark');

  useEffect(() => {
    loadStats();
    loadTargets();
    loadVulnerabilities();
  }, [refreshKey]);

  const loadTargets = async () => {
    try {
      const data = await api.targets.list();
      const targetsData = Array.isArray(data) ? data : (data && typeof data === 'object' ? [data] : []);
      setTargets(targetsData);
    } catch (error) {
      console.error('加载目标失败:', error);
      setTargets([]);
    }
  };

  const loadVulnerabilities = async () => {
    try {
      const data = await api.vulnerabilities.list();
      const vulnsData = Array.isArray(data) ? data : (data && typeof data === 'object' ? [data] : []);
      setVulnerabilities(vulnsData);
      
      const confirmedCount = vulnsData.filter((v: any) => v.confirmed).length;
      setStats({
        targets: targets.length,
        vulns: vulnsData.length,
        confirmed: confirmedCount
      });
    } catch (error) {
      console.error('加载漏洞失败:', error);
      setVulnerabilities([]);
    }
  };

  const handleScanComplete = () => {
    setRefreshKey(k => k + 1);
  };

  const handleMouseDown = (e: React.MouseEvent) => {
    e.preventDefault();
    setIsResizing(true);
    setStartX(e.clientX);
    setStartWidth(sidebarWidth);
  };

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!isResizing) return;
      const delta = e.clientX - startX;
      const newWidth = Math.max(200, startWidth + delta);
      setSidebarWidth(newWidth);
    };

    const handleMouseUp = () => {
      setIsResizing(false);
    };

    if (isResizing) {
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
    }

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [isResizing, startX, startWidth]);

  const handleDeleteTarget = async (id: string) => {
    try {
      await api.targets.delete(id);
      loadTargets();
    } catch (error) {
      console.error('删除目标失败:', error);
    }
  };

  const handleExportTargets = () => {
    const dataStr = JSON.stringify(targets, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `targets-${new Date().toISOString().split('T')[0]}.json`;
    link.click();
  };

  const handleAddTarget = async () => {
    if (!newTarget.name || !newTarget.value) {
      alert('请填写名称和值');
      return;
    }
    try {
      await api.targets.create({
        name: newTarget.name,
        value: newTarget.value,
        type: newTarget.type,
        description: newTarget.description,
        tags: newTarget.tags
      });
      alert('目标添加成功！');
      setShowAddModal(false);
      setNewTarget({ name: '', value: '', type: 'domain', description: '', tags: '' });
      loadTargets();
      loadStats();
    } catch (error) {
      console.error('添加目标失败:', error);
      alert('添加失败！');
    }
  };

  const handleImport = async () => {
    if (!importText.trim()) {
      alert('请输入要导入的内容');
      return;
    }
    const lines = importText.split('\n').filter(line => line.trim());
    if (lines.length === 0) {
      alert('没有有效的数据');
      return;
    }
    try {
      for (const line of lines) {
        const value = line.trim();
        if (!value) continue;
        await api.targets.create({
          name: `导入-${value}`,
          value: value,
          type: importType,
          description: '批量导入',
          tags: '批量导入'
        });
      }
      alert(`成功导入 ${lines.length} 个目标！`);
      setShowImportModal(false);
      setImportText('');
      loadTargets();
      loadStats();
    } catch (error) {
      console.error('导入失败:', error);
      alert('导入失败！');
    }
  };

  const handleStartScanForTarget = async (target: any) => {
    if (scanningTargets.includes(target.id)) {
      return;
    }
    setScanningTargets([...scanningTargets, target.id]);
    try {
      await api.recon.startRecon({
        target: target.value,
        type: target.type,
        options: {
          subdomain: true,
          portscan: true,
          fingerprint: true
        }
      });
      alert(`已为 ${target.name} 启动扫描任务！`);
      setTimeout(() => {
        setScanningTargets(scanningTargets.filter(id => id !== target.id));
        loadTargets();
      }, 3000);
    } catch (error) {
      console.error('启动扫描失败:', error);
      alert('启动扫描失败！');
      setScanningTargets(scanningTargets.filter(id => id !== target.id));
    }
  };

  const loadStats = async () => {
    try {
      const [targetsData, vulns] = await Promise.all([
        api.targets.list(),
        api.vulnerabilities.list()
      ]);
      const targetsCount = Array.isArray(targetsData) ? targetsData.length : 0;
      const vulnsData = Array.isArray(vulns) ? vulns : [];
      const confirmedCount = vulnsData.filter((v: any) => v.confirmed).length;
      setStats({
        targets: targetsCount,
        vulns: vulnsData.length,
        confirmed: confirmedCount
      });
    } catch (error) {
      console.error('加载统计数据失败:', error);
      setStats({ targets: 0, vulns: 0, confirmed: 0 });
    }
  };

  const handleStartScan = () => {
    setActiveTab('active-scan');
  };

  const handleVerifyVuln = async (id: string) => {
    try {
      await api.vulnerabilities.confirm(id);
      loadVulnerabilities();
    } catch (error) {
      console.error('确认漏洞失败:', error);
    }
  };

  const getSelectedTargetDetail = () => {
    if (!selectedDomain) return null;
    return targets.find(t => t.name === selectedDomain || t.value === selectedDomain);
  };

  return (
    <div 
      className={cn(
        "flex h-screen w-full overflow-hidden font-sans transition-colors duration-200",
        theme === 'dark' ? "dark" : ""
      )}
    >
      <aside className={cn(
        "border-r flex flex-col transition-colors duration-200",
        theme === 'dark' ? "bg-gray-900 border-gray-800" : "bg-white border-gray-200"
      )}>
        <div className="p-4 border-b border-gray-800 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center shadow-lg shadow-blue-900/50">
              <Target className="w-5 h-5 text-white" />
            </div>
            <div>
              <h1 className={cn("font-bold text-lg tracking-tight", theme === 'dark' ? "text-white" : "text-gray-900")}>猎影</h1>
              <p className={cn("text-xs", theme === 'dark' ? "text-gray-500" : "text-gray-500")}>Lieying v1.0.0</p>
            </div>
          </div>
          <button
            onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}
            className={cn(
              "p-1.5 rounded-lg transition-colors",
              theme === 'dark' ? "hover:bg-gray-800" : "hover:bg-gray-100"
            )}
            title={theme === 'dark' ? '切换白色主题' : '切换深色主题'}
          >
            {theme === 'dark' ? (
              <Sun className="w-4 h-4 text-yellow-400" />
            ) : (
              <Moon className="w-4 h-4 text-gray-600" />
            )}
          </button>
        </div>

        <div className="flex-1 py-4 px-3 space-y-1 overflow-y-auto">
          <SidebarItem 
            icon={Database} 
            label="资产概览" 
            active={activeTab === 'dashboard'}
            onClick={() => setActiveTab('dashboard')}
            theme={theme}
          />
          <SidebarItem 
            icon={Search} 
            label="信息收集" 
            active={activeTab === 'recon'}
            onClick={() => setActiveTab('recon')}
            theme={theme}
          />
          <SidebarItem 
            icon={Scan} 
            label="漏洞扫描任务" 
            active={activeTab === 'vuln-scan'}
            onClick={() => setActiveTab('vuln-scan')}
            theme={theme}
          />
          <SidebarItem 
            icon={ShieldAlert} 
            label="漏洞结果" 
            active={activeTab === 'vuln-results'}
            onClick={() => setActiveTab('vuln-results')}
            theme={theme}
          />
          <SidebarItem 
            icon={BarChart3} 
            label="漏洞看板" 
            active={activeTab === 'vuln-dashboard'}
            onClick={() => setActiveTab('vuln-dashboard')}
            theme={theme}
          />
          <SidebarItem 
            icon={List} 
            label="漏洞管理" 
            active={activeTab === 'vuln-management'}
            onClick={() => setActiveTab('vuln-management')}
            theme={theme}
          />
          <div className="pt-2 pb-1">
            <p className="px-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">本地网络渗透</p>
          </div>
          <SidebarItem 
            icon={TargetIcon} 
            label="主动扫描" 
            active={activeTab === 'active-scan'}
            onClick={() => setActiveTab('active-scan')}
            theme={theme}
          />
          <SidebarItem 
            icon={Wifi} 
            label="网络设置" 
            active={activeTab === 'network-settings'}
            onClick={() => setActiveTab('network-settings')}
            theme={theme}
          />
          <SidebarItem 
            icon={List} 
            label="请求历史" 
            active={activeTab === 'request-history'}
            onClick={() => setActiveTab('request-history')}
            theme={theme}
          />
          <SidebarItem 
            icon={Globe} 
            label="代理池" 
            active={activeTab === 'proxy-pool'}
            onClick={() => setActiveTab('proxy-pool')}
            theme={theme}
          />
          <SidebarItem 
            icon={MessageSquare} 
            label="AI对话" 
            active={activeTab === 'ai-chat'}
            onClick={() => setActiveTab('ai-chat')}
            theme={theme}
          />
          <SidebarItem 
            icon={Cpu} 
            label="AI配置" 
            active={activeTab === 'ai-config'}
            onClick={() => setActiveTab('ai-config')}
            theme={theme}
          />
          <SidebarItem 
            icon={Bot} 
            label="AI助手" 
            active={activeTab === 'ai'}
            onClick={() => setActiveTab('ai')}
            theme={theme}
          />
          <div className="pt-2 pb-1">
            <p className="px-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">教育SRC</p>
          </div>
          <SidebarItem 
            icon={GraduationCap} 
            label="高校域名库" 
            active={activeTab === 'edu-universities'}
            onClick={() => setActiveTab('edu-universities')}
            theme={theme}
          />
          <SidebarItem 
            icon={ShieldAlert} 
            label="教育SRC工具" 
            active={activeTab === 'edu-tools'}
            onClick={() => setActiveTab('edu-tools')}
            theme={theme}
          />
          <SidebarItem 
            icon={GraduationCap} 
            label="教育SRC" 
            active={activeTab === 'edu'}
            onClick={() => setActiveTab('edu')}
            theme={theme}
          />
          <div className="pt-2 pb-1">
            <p className="px-3 text-xs font-semibold text-gray-500 uppercase tracking-wider">报告工具</p>
          </div>
          <SidebarItem 
            icon={FileText} 
            label="报告生成器" 
            active={activeTab === 'report'}
            onClick={() => setActiveTab('report')}
            theme={theme}
          />
        </div>

        <div className="p-3 border-t border-gray-800">
          <SidebarItem 
            icon={Settings} 
            label="设置" 
            active={activeTab === 'settings'}
            onClick={() => setActiveTab('settings')}
            theme={theme}
          />
        </div>
      </aside>

      <main className="flex-1 flex flex-col min-w-0">
        <header className="h-14 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-6">
          <div className="flex items-center gap-2">
            <h2 className="font-semibold text-gray-100">
              {activeTab === 'dashboard' && '资产概览'}
              {activeTab === 'recon' && '信息收集'}
              {activeTab === 'scan' && '漏洞扫描'}
              {activeTab === 'vuln-scan' && '漏洞扫描任务'}
              {activeTab === 'vuln-results' && '漏洞结果'}
              {activeTab === 'vuln-dashboard' && '漏洞看板'}
              {activeTab === 'vuln-management' && '漏洞管理'}
              {activeTab === 'vulns' && '漏洞管理'}
              {activeTab === 'active-scan' && '主动扫描'}
              {activeTab === 'network-settings' && '网络设置'}
              {activeTab === 'request-history' && '请求历史'}
              {activeTab === 'proxy-pool' && '代理池'}
              {activeTab === 'ai-chat' && 'AI对话'}
              {activeTab === 'ai-config' && 'AI配置'}
              {activeTab === 'ai' && 'AI助手'}
              {activeTab === 'edu-universities' && '高校域名库'}
              {activeTab === 'edu-tools' && '教育SRC工具'}
              {activeTab === 'edu' && '教育SRC'}
              {activeTab === 'report' && '报告生成器'}
              {activeTab === 'settings' && '设置'}
            </h2>
          </div>
          <div className="flex items-center gap-3">
            <button 
              onClick={handleStartScan}
              className="flex items-center gap-2 px-4 py-1.5 bg-green-600 hover:bg-green-500 text-white text-sm font-medium rounded-md transition-colors shadow-lg shadow-green-900/20"
            >
              <Play className="w-4 h-4" />
              开始扫描
            </button>
          </div>
        </header>

        <div className="flex-1 flex overflow-hidden relative">
          <div
            className={`bg-gray-900/50 border-r border-gray-800 flex flex-col transition-all duration-150 ${sidebarCollapsed ? 'w-12' : ''}`}
            style={!sidebarCollapsed ? { width: sidebarWidth } : undefined}
          >
            {!sidebarCollapsed && (
              <div
                className="absolute left-0 top-0 bottom-0 w-1 cursor-ew-resize hover:bg-blue-500/50 transition-colors"
                onMouseDown={handleMouseDown}
              />
            )}
            <div className="p-4 border-b border-gray-800 flex items-center justify-between">
              {!sidebarCollapsed && <h3 className="text-xs font-bold text-gray-500 uppercase tracking-wider">目标资产</h3>}
              <button
                onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
                className="p-1 hover:bg-gray-800 rounded transition-colors"
                title={sidebarCollapsed ? '展开' : '折叠'}
              >
                {sidebarCollapsed ? '→' : '←'}
              </button>
            </div>
            {!sidebarCollapsed && (
              <>
                <div className="p-2 border-b border-gray-800 flex items-center gap-2">
                  <button
                    onClick={() => setShowAddModal?.(true)}
                    className="text-xs px-2 py-1 bg-green-600/20 text-green-400 hover:bg-green-600/30 rounded transition-colors"
                    title="添加目标"
                  >
                    + 添加
                  </button>
                  <button
                    onClick={() => setShowImportModal?.(true)}
                    className="text-xs px-2 py-1 bg-blue-600/20 text-blue-400 hover:bg-blue-600/30 rounded transition-colors"
                    title="批量导入"
                  >
                    📥 导入
                  </button>
                  {targets.length > 0 && (
                    <button
                      onClick={handleExportTargets}
                      className="text-xs text-blue-400 hover:text-blue-300"
                      title="导出目标"
                    >
                      导出
                    </button>
                  )}
                </div>
              <div className="space-y-1">
                {targets.length > 0 ? (
                  targets.map((t: any) => (
                    <div key={t.id || t.name} className="group flex items-center justify-between">
                      <AssetTreeNode
                        name={t.name || t.value || 'unknown'}
                        type={t.type || 'domain'}
                        vulnCount={t.vuln_count || 0}
                        isOpen={selectedDomain === (t.name || t.value)}
                        onToggle={() => setSelectedDomain(t.name || t.value)}
                      />
                      <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100">
                        <button
                          onClick={() => handleStartScanForTarget(t)}
                          disabled={scanningTargets.includes(t.id)}
                          className="text-blue-400 hover:text-blue-300 p-1 disabled:opacity-50"
                          title="启动扫描"
                        >
                          {scanningTargets.includes(t.id) ? '⏳' : '🚀'}
                        </button>
                        <button
                          onClick={() => handleDeleteTarget(t.id)}
                          className="text-red-500 hover:text-red-400 p-1"
                          title="删除目标"
                        >
                          🗑️
                        </button>
                      </div>
                    </div>
                  ))
                ) : (
                  <div className="text-center py-8">
                    <p className="text-xs text-gray-500 mb-3">暂无目标资产</p>
                    <button
                      onClick={() => setShowAddModal(true)}
                      className="text-xs px-3 py-1.5 bg-green-600/20 text-green-400 hover:bg-green-600/30 rounded transition-colors"
                    >
                      + 添加目标
                    </button>
                  </div>
                )}
              </div>
            </>
          )}
        </div>
          <div className="flex-1 flex flex-col bg-gray-950 overflow-hidden">
            <div className="flex-1 overflow-y-auto">
              {activeTab === 'dashboard' && selectedDomain ? (
                <div className="h-full">
                  {(() => {
                    const target = getSelectedTargetDetail();
                    return target ? (
                      <div className="h-full flex flex-col">
                        <div className="p-4 border-b border-gray-800 flex items-center justify-between">
                          <div>
                            <h3 className="text-lg font-semibold text-white">资产概览</h3>
                            <p className="text-sm text-gray-400 mt-1">{selectedDomain}</p>
                          </div>
                          <button
                            onClick={() => setSelectedDomain(null)}
                            className="px-3 py-1.5 text-sm text-gray-400 hover:text-white hover:bg-gray-800 rounded-md transition-colors"
                          >
                            返回概览
                          </button>
                        </div>
                        <div className="flex-1">
                          <AssetOverview targetId={target.id} />
                        </div>
                      </div>
                    ) : null;
                  })()}
                </div>
              ) : (
                <div className="p-6">
                  {activeTab === 'dashboard' && !selectedDomain && (
                    <div className="space-y-6">
                      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                          <div className="flex items-center justify-between">
                            <div>
                              <p className="text-sm text-gray-400">发现目标</p>
                              <p className="text-3xl font-bold text-white mt-1">{stats.targets}</p>
                            </div>
                            <div className="w-12 h-12 bg-blue-500/20 rounded-lg flex items-center justify-center">
                              <Target className="w-6 h-6 text-blue-400" />
                            </div>
                          </div>
                        </div>
                        <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                          <div className="flex items-center justify-between">
                            <div>
                              <p className="text-sm text-gray-400">高危漏洞</p>
                              <p className="text-3xl font-bold text-red-500 mt-1">{stats.vulns}</p>
                            </div>
                            <div className="w-12 h-12 bg-red-500/20 rounded-lg flex items-center justify-center">
                              <ShieldAlert className="w-6 h-6 text-red-400" />
                            </div>
                          </div>
                        </div>
                        <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
                          <div className="flex items-center justify-between">
                            <div>
                              <p className="text-sm text-gray-400">已确认</p>
                              <p className="text-3xl font-bold text-green-500 mt-1">{stats.confirmed}</p>
                            </div>
                            <div className="w-12 h-12 bg-green-500/20 rounded-lg flex items-center justify-center">
                              <Bot className="w-6 h-6 text-green-400" />
                            </div>
                          </div>
                        </div>
                      </div>

                      <div className="space-y-4">
                        <h3 className="text-lg font-semibold text-white">最新发现的漏洞</h3>
                        {vulnerabilities.length > 0 ? (
                          vulnerabilities.slice(0, 5).map((vuln: any) => (
                            <VulnerabilityCard 
                              key={vuln.id}
                              title={vuln.title}
                              severity={vuln.severity as any}
                              url={vuln.url}
                              confirmed={vuln.confirmed}
                              onVerify={() => handleVerifyVuln(vuln.id)}
                            />
                          ))
                        ) : (
                          <div className="text-center py-8 text-gray-500">
                            <ShieldAlert className="w-12 h-12 mx-auto mb-2 opacity-30" />
                            <p>暂无漏洞数据</p>
                            <p className="text-xs mt-1">开始扫描以发现漏洞</p>
                          </div>
                        )}
                      </div>
                      
                      <div className="bg-gray-900/50 border border-gray-800 rounded-xl p-6 text-center">
                        <Database className="w-12 h-12 mx-auto mb-3 text-gray-500" />
                        <p className="text-gray-400">选择左侧目标资产查看详细的资产概览</p>
                      </div>
                    </div>
                  )}

                  {activeTab === 'vulns' && <VulnerabilityManager />}
                  {activeTab === 'vuln-dashboard' && <VulnDashboard />}
                  {activeTab === 'vuln-management' && <VulnManagement />}
                  {activeTab === 'active-scan' && <ActiveScan />}
                  {activeTab === 'network-settings' && <NetworkSettings />}
                  {activeTab === 'request-history' && <RequestHistory />}
                  {activeTab === 'proxy-pool' && <ProxyPool />}
                  {activeTab === 'ai-chat' && <AIChat />}
                  {activeTab === 'ai-config' && <AIConfigPanel />}
                  {activeTab === 'edu-universities' && <EduUniversityManager />}
                  {activeTab === 'edu-tools' && <EduSRCManager />}
                  {activeTab === 'vuln-scan' && <VulnScanPanel />}
                  {activeTab === 'vuln-results' && <VulnResultsPanel />}
                  {activeTab === 'ai' && <AIAssistant />}
                  {activeTab === 'edu' && <EduPanel />}
                  {activeTab === 'report' && <ReportGenerator />}
                  {activeTab === 'recon' && (
                    <div className="h-full">
                      {selectedDomain ? (
                        (() => {
                          const target = getSelectedTargetDetail();
                          return target ? (
                            <div className="h-full flex flex-col">
                              <div className="p-4 border-b border-gray-800 flex items-center justify-between">
                                <div>
                                  <h3 className="text-lg font-semibold text-white">信息收集</h3>
                                  <p className="text-sm text-gray-400 mt-1">{selectedDomain}</p>
                                </div>
                                <button
                                  onClick={() => setSelectedDomain(null)}
                                  className="px-3 py-1.5 text-sm text-gray-400 hover:text-white hover:bg-gray-800 rounded-md transition-colors"
                                >
                                  返回概览
                                </button>
                              </div>
                              <div className="flex-1">
                                <ReconOverview targetId={target.id} />
                              </div>
                            </div>
                          ) : null;
                        })()
                      ) : (
                        <div className="p-6">
                          <div className="bg-gray-900/50 border border-gray-800 rounded-xl p-6 text-center">
                            <Search className="w-12 h-12 mx-auto mb-3 text-gray-500" />
                            <p className="text-gray-400">选择左侧目标资产开始信息收集</p>
                          </div>
                        </div>
                      )}
                    </div>
                  )}
                  {activeTab === 'scan' && <ScanPanel onComplete={handleScanComplete} />}
                </div>
              )}
            </div>
          </div>
        </div>
      </main>

      {showAddModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold text-white">添加目标资产</h3>
              <button onClick={() => setShowAddModal(false)} className="text-gray-400 hover:text-white">✕</button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">名称 *</label>
                <input type="text" value={newTarget.name} onChange={(e) => setNewTarget({...newTarget, name: e.target.value})}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white" placeholder="例如：VIPC6资源网" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">值 *</label>
                <input type="text" value={newTarget.value} onChange={(e) => setNewTarget({...newTarget, value: e.target.value})}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white" placeholder="例如：vipc6.com" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">类型</label>
                <select value={newTarget.type} onChange={(e) => setNewTarget({...newTarget, type: e.target.value})}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white">
                  <option value="domain">域名</option>
                  <option value="ip">IP地址</option>
                  <option value="url">URL地址</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">描述</label>
                <input type="text" value={newTarget.description} onChange={(e) => setNewTarget({...newTarget, description: e.target.value})}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white" placeholder="可选描述" />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">标签</label>
                <input type="text" value={newTarget.tags} onChange={(e) => setNewTarget({...newTarget, tags: e.target.value})}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white" placeholder="例如：重要,测试" />
              </div>
            </div>
            <div className="flex justify-end space-x-3 mt-6">
              <button onClick={() => setShowAddModal(false)} className="px-4 py-2 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600">取消</button>
              <button onClick={handleAddTarget} className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-500">添加</button>
            </div>
          </div>
        </div>
      )}

      {showImportModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 w-full max-w-2xl">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold text-white">批量导入目标</h3>
              <button onClick={() => setShowImportModal(false)} className="text-gray-400 hover:text-white">✕</button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">目标类型</label>
                <select value={importType} onChange={(e) => setImportType(e.target.value as any)}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white">
                  <option value="domain">域名</option>
                  <option value="ip">IP地址</option>
                  <option value="url">URL地址</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">导入内容 (每行一个)</label>
                <textarea value={importText} onChange={(e) => setImportText(e.target.value)}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white h-64 font-mono text-sm"
                  placeholder={importType === 'domain' ? 'vipc6.com\nexample.com\nbaidu.com' : importType === 'ip' ? '192.168.1.1\n10.0.0.1\n8.8.8.8' : 'https://example.com\nhttps://baidu.com'} />
              </div>
              <div className="bg-gray-800/50 rounded-lg p-3">
                <p className="text-sm text-gray-400">💡 支持格式：每行一个目标</p>
                <p className="text-sm text-gray-400 mt-1">📊 当前输入：{importText.split('\n').filter(l => l.trim()).length} 行</p>
              </div>
            </div>
            <div className="flex justify-end space-x-3 mt-6">
              <button onClick={() => setShowImportModal(false)} className="px-4 py-2 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600">取消</button>
              <button onClick={handleImport} className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500">开始导入</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
