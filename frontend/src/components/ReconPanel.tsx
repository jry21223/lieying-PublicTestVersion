import React, { useState, useEffect } from 'react';
import { targetsApi, activeScanApi } from '../services/api';

interface Target {
  id: string;
  name: string;
  value: string;
  type: string;
  description?: string;
  tags?: string;
  created_at: string;
}

interface Subdomain {
  id: string;
  domain: string;
  subdomain: string;
  ip: string;
  ports: string;
  status: string;
}

interface ScanResult {
  id: string;
  type: string;
  target: string;
  data: string;
  severity: string;
  created_at: string;
}

const ReconPanel: React.FC = () => {
  const [targets, setTargets] = useState<Target[]>([]);
  const [selectedTarget, setSelectedTarget] = useState<Target | null>(null);
  const [subdomains, setSubdomains] = useState<Subdomain[]>([]);
  const [scanResults, setScanResults] = useState<ScanResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [scanType, setScanType] = useState<'subdomain' | 'port' | 'vuln' | 'web'>('subdomain');
  const [activeTaskId, setActiveTaskId] = useState<string | null>(null);
  const [stats, setStats] = useState({ subdomains: 0, ips: 0, ports: 0, vulns: 0 });

  useEffect(() => {
    loadTargets();
  }, []);

  const loadTargets = async () => {
    try {
      const data = await targetsApi.list() as Target[];
      setTargets(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('加载目标失败:', error);
    }
  };

  const handleScan = async () => {
    if (!selectedTarget) {
      alert('请先选择一个目标');
      return;
    }

    setLoading(true);
    try {
      const task = await activeScanApi.createTask({
        name: `信息收集 - ${selectedTarget.name}`,
        target: selectedTarget.value,
        type: scanType,
      }) as { id: string };

      setActiveTaskId(task.id);

      setTimeout(async () => {
        try {
          const results = await activeScanApi.getResults(task.id) as ScanResult[];
          setScanResults(Array.isArray(results) ? results : []);

          if (scanType === 'subdomain') {
            const newSubdomains: Subdomain[] = results
              .filter(r => r.type === 'subdomain_found')
              .map((r, i) => ({
                id: `sub-${i}`,
                domain: selectedTarget.value,
                subdomain: r.target,
                ip: extractIP(r.data),
                ports: extractPorts(r.data),
                status: 'pending',
              }));
            setSubdomains(prev => [...prev, ...newSubdomains]);
            setStats(prev => ({ ...prev, subdomains: newSubdomains.length }));
          }

          if (scanType === 'vuln') {
            setStats(prev => ({ ...prev, vulns: results.filter(r => r.type === 'vuln_detected').length }));
          }
        } catch (e) {
          console.error('获取结果失败:', e);
        }
      }, 3000);
    } catch (error: any) {
      alert(`扫描启动失败: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const extractIP = (data: string): string => {
    const match = data.match(/\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}/);
    return match ? match[0] : '未知';
  };

  const extractPorts = (data: string): string => {
    const match = data.match(/端口:\s*(\d+)/);
    return match ? match[1] : '';
  };

  const getSeverityColor = (severity: string) => {
    switch (severity.toLowerCase()) {
      case 'critical': return 'bg-red-600 text-white';
      case 'high': return 'bg-red-500/20 text-red-400 border border-red-500/30';
      case 'medium': return 'bg-yellow-500/20 text-yellow-400 border border-yellow-500/30';
      case 'low': return 'bg-blue-500/20 text-blue-400 border border-blue-500/30';
      default: return 'bg-gray-500/20 text-gray-400 border border-gray-500/30';
    }
  };

  return (
    <div className="h-full flex flex-col bg-gray-950">
      <div className="p-4 bg-gray-900 border-b border-gray-800">
        <h2 className="text-lg font-semibold text-white">信息收集</h2>
        <p className="text-sm text-gray-400 mt-1">专家级信息收集与子域名探测</p>
      </div>

      <div className="flex-1 overflow-auto p-6">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <p className="text-sm text-gray-400">子域名</p>
            <p className="text-3xl font-bold text-blue-400 mt-1">{stats.subdomains}</p>
          </div>
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <p className="text-sm text-gray-400">IP资产</p>
            <p className="text-3xl font-bold text-green-400 mt-1">{stats.ips}</p>
          </div>
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <p className="text-sm text-gray-400">端口服务</p>
            <p className="text-3xl font-bold text-purple-400 mt-1">{stats.ports}</p>
          </div>
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <p className="text-sm text-gray-400">发现漏洞</p>
            <p className="text-3xl font-bold text-red-400 mt-1">{stats.vulns}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="space-y-4">
            <div className="bg-gray-900 rounded-xl border border-gray-800 p-4">
              <h3 className="text-sm font-medium text-gray-400 mb-3">选择目标</h3>
              {targets.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  <div className="text-4xl mb-4">🎯</div>
                  <p>暂无目标</p>
                  <p className="text-sm mt-2">请先在资产概览中添加目标</p>
                </div>
              ) : (
                <div className="space-y-2 max-h-60 overflow-y-auto">
                  {targets.map((t) => (
                    <div
                      key={t.id}
                      onClick={() => setSelectedTarget(t)}
                      className={`p-3 rounded-lg cursor-pointer transition ${
                        selectedTarget?.id === t.id
                          ? 'bg-blue-600/20 border border-blue-500'
                          : 'bg-gray-800 hover:bg-gray-750'
                      }`}
                    >
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-white font-medium">{t.name}</p>
                          <p className="text-sm text-gray-400">{t.value}</p>
                        </div>
                        <span className={`px-2 py-0.5 rounded text-xs ${
                          t.type === 'domain' ? 'bg-blue-600/20 text-blue-400' :
                          t.type === 'ip' ? 'bg-green-600/20 text-green-400' :
                          'bg-purple-600/20 text-purple-400'
                        }`}>
                          {t.type}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            <div className="bg-gray-900 rounded-xl border border-gray-800 p-4">
              <h3 className="text-sm font-medium text-gray-400 mb-3">扫描类型</h3>
              <div className="grid grid-cols-2 gap-2">
                {[
                  { key: 'subdomain', label: '🔍 子域名探测', desc: '发现子域名' },
                  { key: 'port', label: '🚪 端口扫描', desc: '扫描常见端口' },
                  { key: 'vuln', label: '🛡️ 漏洞检测', desc: 'OWASP Top 10' },
                  { key: 'web', label: '🌐 Web分析', desc: 'Web信息收集' },
                ].map((type) => (
                  <button
                    key={type.key}
                    onClick={() => setScanType(type.key as any)}
                    className={`p-3 rounded-lg text-left transition ${
                      scanType === type.key
                        ? 'bg-blue-600/20 border border-blue-500'
                        : 'bg-gray-800 hover:bg-gray-750'
                    }`}
                  >
                    <p className="text-white font-medium">{type.label}</p>
                    <p className="text-xs text-gray-400">{type.desc}</p>
                  </button>
                ))}
              </div>
            </div>

            <button
              onClick={handleScan}
              disabled={!selectedTarget || loading}
              className={`w-full px-4 py-3 rounded-lg text-white font-medium transition ${
                !selectedTarget || loading
                  ? 'bg-gray-600 cursor-not-allowed'
                  : 'bg-green-600 hover:bg-green-500'
              }`}
            >
              {loading ? '扫描中...' : `🚀 开始${scanType === 'subdomain' ? '子域名探测' : scanType === 'port' ? '端口扫描' : scanType === 'vuln' ? '漏洞检测' : 'Web分析'}`}
            </button>
          </div>

          <div className="space-y-4">
            <div className="bg-gray-900 rounded-xl border border-gray-800 p-4">
              <h3 className="text-sm font-medium text-gray-400 mb-3">扫描结果</h3>
              {scanResults.length === 0 ? (
                <div className="text-center py-8 text-gray-500">
                  <div className="text-4xl mb-4">📋</div>
                  <p>暂无扫描结果</p>
                </div>
              ) : (
                <div className="space-y-2 max-h-96 overflow-y-auto">
                  {scanResults.map((result) => (
                    <div key={result.id} className="p-3 bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between mb-1">
                        <span className={`px-2 py-0.5 rounded text-xs ${getSeverityColor(result.severity)}`}>
                          {result.severity.toUpperCase()}
                        </span>
                        <span className="text-xs text-gray-500">
                          {new Date(result.created_at).toLocaleString('zh-CN')}
                        </span>
                      </div>
                      <p className="text-white text-sm">{result.target}</p>
                      <p className="text-xs text-gray-400 mt-1 line-clamp-2">{result.data}</p>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {subdomains.length > 0 && (
              <div className="bg-gray-900 rounded-xl border border-gray-800 p-4">
                <h3 className="text-sm font-medium text-gray-400 mb-3">发现子域名 ({subdomains.length})</h3>
                <div className="space-y-2 max-h-60 overflow-y-auto">
                  {subdomains.map((sub) => (
                    <div key={sub.id} className="p-3 bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-white font-medium">{sub.subdomain}</p>
                          <p className="text-xs text-gray-400">IP: {sub.ip} | 端口: {sub.ports || '未知'}</p>
                        </div>
                        <button
                          onClick={() => {
                            setSelectedTarget({ id: sub.id, name: sub.subdomain, value: sub.subdomain, type: 'domain', created_at: '' });
                            setScanType('vuln');
                          }}
                          className="px-2 py-1 bg-red-600/20 text-red-400 rounded text-xs hover:bg-red-600/30"
                        >
                          漏洞扫描
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ReconPanel;
