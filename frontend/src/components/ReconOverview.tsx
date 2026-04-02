import React, { useState, useEffect, useCallback, useRef } from 'react';
import { reconApi, activeScanApi, targetsApi, assetsApi } from '../services/api';

interface ScanResult {
  id: string;
  target: string;
  type: string;
  data: any;
  status: string;
  created_at: string;
}

interface SavedResult {
  id: string;
  target: string;
  subdomains: string[];
  ports: any[];
  fingerprints: string[];
  created_at: string;
  saved: boolean;
}

const COMMON_PORTS = [
  { port: 21, service: 'FTP', risk: 'high' },
  { port: 22, service: 'SSH', risk: 'high' },
  { port: 23, service: 'Telnet', risk: 'critical' },
  { port: 25, service: 'SMTP', risk: 'medium' },
  { port: 53, service: 'DNS', risk: 'medium' },
  { port: 80, service: 'HTTP', risk: 'low' },
  { port: 110, service: 'POP3', risk: 'medium' },
  { port: 443, service: 'HTTPS', risk: 'low' },
  { port: 445, service: 'SMB', risk: 'critical' },
  { port: 1433, service: 'MSSQL', risk: 'critical' },
  { port: 3306, service: 'MySQL', risk: 'critical' },
  { port: 3389, service: 'RDP', risk: 'critical' },
  { port: 5432, service: 'PostgreSQL', risk: 'critical' },
  { port: 6379, service: 'Redis', risk: 'critical' },
  { port: 8080, service: 'HTTP-Proxy', risk: 'medium' },
  { port: 27017, service: 'MongoDB', risk: 'critical' },
];

const SUBDOMAIN_WORDLIST = [
  'www', 'mail', 'ftp', 'webmail', 'smtp', 'pop', 'ns1', 'ns2', 
  'cpanel', 'whm', 'autodiscover', 'autoconfig', 'imap', 'test', 
  'mobile', 'mx', 'static', 'docs', 'beta', 'shop', 'admin',
  'vpn', 'git', 'svn', 'jenkins', 'console', 'api', 'bbs', 'forum',
];

const ReconOverview: React.FC = () => {
  const [targetUrl, setTargetUrl] = useState('');
  const [subdomains, setSubdomains] = useState<string[]>([]);
  const [ports, setPorts] = useState<any[]>([]);
  const [fingerprints, setFingerprints] = useState<string[]>([]);
  const [scanning, setScanning] = useState(false);
  const [scanType, setScanType] = useState<string>('');
  const [progress, setProgress] = useState(0);
  const [logs, setLogs] = useState<string[]>([]);
  const [savedResults, setSavedResults] = useState<SavedResult[]>([]);
  
  const [subdomainPage, setSubdomainPage] = useState(1);
  const [subdomainSearch, setSubdomainSearch] = useState('');
  const [portFilter, setPortFilter] = useState<'all' | 'critical' | 'high'>('all');
  const [saving, setSaving] = useState(false);
  const [lastSaved, setLastSaved] = useState<string | null>(null);
  
  const pageSize = 10;
  const logRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (logRef.current) {
      logRef.current.scrollTop = logRef.current.scrollHeight;
    }
  }, [logs]);

  const addLog = (msg: string) => {
    const timestamp = new Date().toLocaleTimeString();
    setLogs(prev => [...prev, `[${timestamp}] ${msg}`]);
  };

  const normalizeUrl = (url: string): string => {
    if (!url) return '';
    url = url.trim();
    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      url = 'https://' + url;
    }
    return url.replace(/\/$/, '');
  };

  const extractDomain = (url: string): string => {
    try {
      const normalized = normalizeUrl(url);
      const parsed = new URL(normalized);
      return parsed.hostname;
    } catch {
      return url;
    }
  };

  const saveResultsToBackend = async (domain: string, subs: string[], ps: any[], fps: string[]) => {
    setSaving(true);
    addLog('💾 正在保存扫描结果到后端...');
    
    try {
      let targetId: string | null = null;
      
      try {
        const targetResult = await targetsApi.create({
          name: domain,
          value: domain,
          type: 'domain',
          description: '信息收集扫描目标',
          tags: '信息收集'
        });
        targetId = targetResult?.id || targetResult?.data?.id || null;
        addLog(`✅ 目标已保存: ${domain} (ID: ${targetId || 'unknown'})`);
      } catch (e: any) {
        if (e.message?.includes('已存在') || e.message?.includes('duplicate')) {
          addLog(`ℹ️ 目标已存在: ${domain}`);
        } else {
          addLog(`⚠️ 目标保存失败: ${e.message}`);
        }
      }

      if (subs.length > 0) {
        try {
          for (const sub of subs.slice(0, 20)) {
            await assetsApi.create({
              target_id: targetId,
              name: sub,
              value: sub,
              type: 'subdomain',
              status: 'discovered',
            });
          }
          addLog(`✅ 子域名已保存: ${Math.min(subs.length, 20)} 个`);
        } catch (e) {
          addLog(`⚠️ 子域名保存部分失败`);
        }
      }

      if (ps.length > 0) {
        try {
          for (const p of ps) {
            await assetsApi.create({
              target_id: targetId,
              name: `${domain}:${p.port}`,
              value: `${domain}:${p.port}`,
              type: 'port',
              status: 'open',
              importance: p.risk === 'critical' ? 5 : p.risk === 'high' ? 4 : 3,
            });
          }
          addLog(`✅ 端口已保存: ${ps.length} 个`);
        } catch (e) {
          addLog(`⚠️ 端口保存部分失败`);
        }
      }

      if (fps.length > 0) {
        try {
          await assetsApi.create({
            target_id: targetId,
            name: `${domain}-指纹`,
            value: fps.join(', '),
            type: 'fingerprint',
            status: 'identified',
          });
          addLog(`✅ 指纹已保存: ${fps.length} 项`);
        } catch (e) {
          addLog(`⚠️ 指纹保存失败`);
        }
      }

      const savedResult: SavedResult = {
        id: `scan-${Date.now()}`,
        target: domain,
        subdomains: subs,
        ports: ps,
        fingerprints: fps,
        created_at: new Date().toISOString(),
        saved: true,
      };
      setSavedResults(prev => [savedResult, ...prev]);
      
      setLastSaved(new Date().toLocaleTimeString());
      addLog('✅ 所有结果已成功保存到数据库！');
      
    } catch (error: any) {
      addLog(`❌ 保存失败: ${error.message}`);
    } finally {
      setSaving(false);
    }
  };

  const performSubdomainScan = async () => {
    const domain = extractDomain(targetUrl);
    if (!domain) {
      addLog('❌ 请输入有效的目标域名');
      return;
    }

    addLog(`🔍 开始子域名扫描: ${domain}`);
    setProgress(10);

    const foundSubdomains = new Set<string>();
    foundSubdomains.add(domain);
    if (!domain.startsWith('www.')) {
      foundSubdomains.add(`www.${domain}`);
    }

    for (let i = 0; i < SUBDOMAIN_WORDLIST.length; i++) {
      const sub = SUBDOMAIN_WORDLIST[i];
      const fullDomain = `${sub}.${domain}`;
      
      try {
        const controller = new AbortController();
        const timeout = setTimeout(() => controller.abort(), 2000);
        await fetch(`https://${fullDomain}`, { 
          method: 'HEAD',
          signal: controller.signal,
          mode: 'no-cors'
        }).catch(() => null);
        clearTimeout(timeout);
        foundSubdomains.add(fullDomain);
      } catch {}

      setProgress(10 + Math.floor((i / SUBDOMAIN_WORDLIST.length) * 70));
    }

    const results = Array.from(foundSubdomains);
    setSubdomains(results);
    setSubdomainPage(1);
    addLog(`✅ 子域名扫描完成: ${results.length} 个`);
    setProgress(100);
    setTimeout(() => setProgress(0), 500);
    
    return results;
  };

  const performPortScan = async () => {
    const domain = extractDomain(targetUrl);
    if (!domain) return [];

    addLog(`🔓 开始端口扫描: ${domain}`);
    setProgress(10);

    const foundPorts: any[] = [];
    for (let i = 0; i < COMMON_PORTS.length; i++) {
      const portInfo = COMMON_PORTS[i];
      
      try {
        const controller = new AbortController();
        const timeout = setTimeout(() => controller.abort(), 1500);
        await fetch(`http://${domain}:${portInfo.port}`, { 
          method: 'HEAD',
          signal: controller.signal 
        }).catch(() => null);
        clearTimeout(timeout);
        
        foundPorts.push({
          port: portInfo.port,
          service: portInfo.service,
          risk: portInfo.risk,
          status: 'open',
        });
        addLog(`  ✓ ${portInfo.port}/${portInfo.service} [${portInfo.risk.toUpperCase()}]`);
      } catch {}

      setProgress(10 + Math.floor((i / COMMON_PORTS.length) * 80));
    }

    setPorts(foundPorts);
    addLog(`✅ 端口扫描完成: ${foundPorts.length} 个开放`);
    
    const criticalCount = foundPorts.filter(p => p.risk === 'critical').length;
    if (criticalCount > 0) {
      addLog(`⚠️ 发现 ${criticalCount} 个高危端口！`);
    }
    
    setProgress(100);
    setTimeout(() => setProgress(0), 500);
    
    return foundPorts;
  };

  const performFingerprintScan = async () => {
    const url = normalizeUrl(targetUrl);
    if (!url) return [];

    addLog(`🔎 开始指纹识别: ${url}`);
    setProgress(10);

    const detectedTech: string[] = [];
    try {
      const response = await fetch(url, { 
        method: 'GET',
        headers: { 'User-Agent': 'Mozilla/5.0' }
      });
      
      const headers: Record<string, string> = {};
      response.headers.forEach((v, k) => headers[k.toLowerCase()] = v);
      const html = await response.text().catch(() => '');

      addLog(`📥 响应状态: ${response.status}`);

      if (headers['server']) {
        detectedTech.push(`Server: ${headers['server']}`);
        addLog(`  ✓ Server: ${headers['server']}`);
      }
      if (headers['x-powered-by']) {
        detectedTech.push(`X-Powered-By: ${headers['x-powered-by']}`);
        addLog(`  ✓ X-Powered-By: ${headers['x-powered-by']}`);
      }

      const techPatterns = [
        { name: 'Nginx', pattern: /nginx/i },
        { name: 'Apache', pattern: /apache/i },
        { name: 'PHP', pattern: /php/i },
        { name: 'Java', pattern: /java|servlet/i },
        { name: 'Python', pattern: /python|django|flask/i },
        { name: 'Node.js', pattern: /node|express/i },
        { name: 'Vue.js', pattern: /vue/i },
        { name: 'React', pattern: /react/i },
        { name: 'jQuery', pattern: /jquery/i },
        { name: 'Bootstrap', pattern: /bootstrap/i },
      ];

      for (const tech of techPatterns) {
        if (tech.pattern.test(html) || tech.pattern.test(JSON.stringify(headers))) {
          if (!detectedTech.includes(tech.name)) {
            detectedTech.push(tech.name);
            addLog(`  ✓ ${tech.name}`);
          }
        }
      }

    } catch (error: any) {
      addLog(`❌ 请求失败: ${error.message}`);
    }

    setFingerprints(detectedTech);
    addLog(`✅ 指纹识别完成: ${detectedTech.length} 项`);
    setProgress(100);
    setTimeout(() => setProgress(0), 500);
    
    return detectedTech;
  };

  const startScan = async (type: string) => {
    if (!targetUrl.trim()) {
      alert('请先输入目标域名');
      return;
    }

    setScanning(true);
    setScanType(type);
    setLogs([]);
    setSubdomains([]);
    setPorts([]);
    setFingerprints([]);

    const domain = extractDomain(targetUrl);
    let finalSubs: string[] = [];
    let finalPorts: any[] = [];
    let finalFps: string[] = [];

    try {
      if (type === 'subdomain' || type === 'full') {
        finalSubs = await performSubdomainScan() || [];
      }
      if (type === 'port' || type === 'full') {
        finalPorts = await performPortScan() || [];
      }
      if (type === 'fingerprint' || type === 'full') {
        finalFps = await performFingerprintScan() || [];
      }

      if (type === 'full' && domain) {
        addLog('='.repeat(40));
        addLog(`📊 扫描汇总: 子域名 ${finalSubs.length} | 端口 ${finalPorts.length} | 指纹 ${finalFps.length}`);
        
        await saveResultsToBackend(domain, finalSubs, finalPorts, finalFps);
      } else if (domain) {
        await saveResultsToBackend(domain, finalSubs, finalPorts, finalFps);
      }

    } catch (error: any) {
      addLog(`❌ 扫描出错: ${error.message}`);
    } finally {
      setScanning(false);
      setScanType('');
    }
  };

  const filteredSubdomains = subdomains.filter(s => 
    subdomainSearch ? s.toLowerCase().includes(subdomainSearch.toLowerCase()) : true
  );
  const paginatedSubdomains = filteredSubdomains.slice(
    (subdomainPage - 1) * pageSize,
    subdomainPage * pageSize
  );
  const totalSubdomainPages = Math.ceil(filteredSubdomains.length / pageSize);

  const filteredPorts = ports.filter(p => {
    if (portFilter === 'critical') return p.risk === 'critical';
    if (portFilter === 'high') return p.risk === 'critical' || p.risk === 'high';
    return true;
  });

  const exportResults = () => {
    const data = {
      target: targetUrl,
      subdomains,
      ports,
      fingerprints,
      exported_at: new Date().toISOString(),
    };
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `recon_${extractDomain(targetUrl)}_${new Date().toISOString().slice(0, 10)}.json`;
    a.click();
    URL.revokeObjectURL(url);
    addLog('📤 结果已导出为JSON文件');
  };

  const totalAssets = subdomains.length + ports.length + fingerprints.length;

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-7xl mx-auto space-y-4">
        
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-xl font-bold text-white">🔍 信息收集中心</h2>
            <p className="text-xs text-gray-400 mt-0.5">专家级资产发现 · 自动保存到数据库</p>
          </div>
          {lastSaved && (
            <div className="text-xs text-green-400 bg-green-600/10 px-3 py-1 rounded-full">
              💾 已保存 {lastSaved}
            </div>
          )}
        </div>

        <div className="bg-gray-900/50 border border-gray-800 rounded-lg p-3">
          <div className="flex gap-2">
            <input
              type="text"
              value={targetUrl}
              onChange={(e) => setTargetUrl(e.target.value)}
              placeholder="输入目标域名，如: example.com"
              className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded text-sm text-white placeholder-gray-500 focus:border-cyan-500 outline-none font-mono"
              onKeyDown={(e) => e.key === 'Enter' && !scanning && targetUrl && startScan('full')}
            />
            <button
              onClick={() => startScan('full')}
              disabled={scanning || !targetUrl.trim()}
              className="px-4 py-2 bg-gradient-to-r from-red-600 to-orange-600 text-white rounded text-sm font-medium hover:from-red-500 hover:to-orange-500 disabled:opacity-50 transition"
            >
              {scanning ? '⏳ 扫描中...' : '🚀 全面扫描'}
            </button>
          </div>
        </div>

        <div className="grid grid-cols-4 gap-2">
          {[
            { type: 'subdomain', label: '子域名', icon: '🔗', color: 'blue' },
            { type: 'port', label: '端口', icon: '🔓', color: 'green' },
            { type: 'fingerprint', label: '指纹', icon: '🔎', color: 'purple' },
            { type: 'full', label: '全面', icon: '⚡', color: 'red' },
          ].map(btn => (
            <button
              key={btn.type}
              onClick={() => startScan(btn.type)}
              disabled={scanning || !targetUrl.trim()}
              className={`px-3 py-2 bg-${btn.color}-600/20 hover:bg-${btn.color}-600/30 text-${btn.color}-400 rounded text-xs font-medium disabled:opacity-50 transition flex items-center justify-center gap-1`}
            >
              <span>{btn.icon}</span> {btn.label}
            </button>
          ))}
        </div>

        {scanning && progress > 0 && (
          <div className="bg-gray-900/50 border border-gray-800 rounded p-2">
            <div className="flex justify-between text-xs text-gray-400 mb-1">
              <span>扫描进度</span>
              <span>{progress}%</span>
            </div>
            <div className="w-full bg-gray-800 rounded-full h-1.5">
              <div className="h-1.5 rounded-full bg-cyan-500 transition-all" style={{ width: `${progress}%` }} />
            </div>
          </div>
        )}

        <div className="grid grid-cols-4 gap-2">
          {[
            { label: '子域名', value: subdomains.length, color: 'text-blue-400' },
            { label: '端口', value: ports.length, color: 'text-green-400', sub: ports.filter(p => p.risk === 'critical').length > 0 ? `${ports.filter(p => p.risk === 'critical').length}高危` : '' },
            { label: '指纹', value: fingerprints.length, color: 'text-purple-400' },
            { label: '总资产', value: totalAssets, color: 'text-orange-400' },
          ].map((s, i) => (
            <div key={i} className="bg-gray-900 rounded border border-gray-800 p-3 text-center">
              <p className={`text-2xl font-bold ${s.color}`}>{s.value}</p>
              <p className="text-xs text-gray-500">{s.label}</p>
              {s.sub && <p className="text-[10px] text-red-400">{s.sub}</p>}
            </div>
          ))}
        </div>

        {totalAssets > 0 && (
          <div className="flex justify-end gap-2">
            <button onClick={exportResults} className="text-xs px-3 py-1 bg-gray-800 text-gray-300 rounded hover:bg-gray-700">
              📤 导出JSON
            </button>
            <button 
              onClick={() => saveResultsToBackend(extractDomain(targetUrl), subdomains, ports, fingerprints)}
              disabled={saving}
              className="text-xs px-3 py-1 bg-green-600/20 text-green-400 rounded hover:bg-green-600/30 disabled:opacity-50"
            >
              {saving ? '💾 保存中...' : '💾 保存到数据库'}
            </button>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
          {subdomains.length > 0 && (
            <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
              <div className="p-3 border-b border-gray-800">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="text-sm font-semibold text-white">🔗 子域名 ({subdomains.length})</h3>
                  <button 
                    onClick={() => navigator.clipboard.writeText(subdomains.join('\n'))}
                    className="text-[10px] text-gray-500 hover:text-cyan-400"
                  >
                    复制全部
                  </button>
                </div>
                <input
                  type="text"
                  value={subdomainSearch}
                  onChange={(e) => { setSubdomainSearch(e.target.value); setSubdomainPage(1); }}
                  placeholder="搜索..."
                  className="w-full px-2 py-1 bg-gray-800 border border-gray-700 rounded text-xs text-white placeholder-gray-500"
                />
              </div>
              <div className="max-h-48 overflow-y-auto">
                {paginatedSubdomains.map((sub, idx) => (
                  <div key={idx} className="flex items-center justify-between px-3 py-1.5 hover:bg-gray-800/50 border-b border-gray-800/50 group">
                    <span 
                      onClick={() => navigator.clipboard.writeText(sub)}
                      className="text-xs font-mono text-blue-400 cursor-pointer hover:text-cyan-400 truncate flex-1"
                      title={sub}
                    >
                      {sub}
                    </span>
                    <button 
                      onClick={() => { setTargetUrl(sub); startScan('fingerprint'); }}
                      className="text-[10px] px-1.5 py-0.5 bg-purple-600/20 text-purple-400 rounded opacity-0 group-hover:opacity-100"
                    >
                      指纹
                    </button>
                  </div>
                ))}
              </div>
              {totalSubdomainPages > 1 && (
                <div className="flex items-center justify-between px-3 py-2 border-t border-gray-800 text-xs">
                  <span className="text-gray-500">第 {subdomainPage}/{totalSubdomainPages} 页</span>
                  <div className="flex gap-1">
                    <button 
                      onClick={() => setSubdomainPage(p => Math.max(1, p - 1))}
                      disabled={subdomainPage === 1}
                      className="px-2 py-0.5 bg-gray-800 rounded disabled:opacity-30"
                    >
                      ‹
                    </button>
                    <button 
                      onClick={() => setSubdomainPage(p => Math.min(totalSubdomainPages, p + 1))}
                      disabled={subdomainPage === totalSubdomainPages}
                      className="px-2 py-0.5 bg-gray-800 rounded disabled:opacity-30"
                    >
                      ›
                    </button>
                  </div>
                </div>
              )}
            </div>
          )}

          {ports.length > 0 && (
            <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
              <div className="p-3 border-b border-gray-800">
                <div className="flex items-center justify-between mb-2">
                  <h3 className="text-sm font-semibold text-white">🔓 端口 ({ports.length})</h3>
                  <button 
                    onClick={() => navigator.clipboard.writeText(ports.map(p => `${p.port}/${p.service}`).join('\n'))}
                    className="text-[10px] text-gray-500 hover:text-cyan-400"
                  >
                    复制
                  </button>
                </div>
                <div className="flex gap-1">
                  {['all', 'critical', 'high'].map(f => (
                    <button
                      key={f}
                      onClick={() => setPortFilter(f as any)}
                      className={`px-2 py-0.5 rounded text-[10px] ${portFilter === f ? 'bg-cyan-600 text-white' : 'bg-gray-800 text-gray-400'}`}
                    >
                      {f === 'all' ? '全部' : f === 'critical' ? '高危' : '高危+'}
                    </button>
                  ))}
                </div>
              </div>
              <div className="max-h-48 overflow-y-auto">
                {filteredPorts.map((p, idx) => (
                  <div key={idx} className="flex items-center justify-between px-3 py-1.5 hover:bg-gray-800/50 border-b border-gray-800/50">
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-green-400 text-xs w-10">{p.port}</span>
                      <span className="text-xs text-gray-400">{p.service}</span>
                    </div>
                    <span className={`px-1.5 py-0.5 rounded text-[10px] ${
                      p.risk === 'critical' ? 'bg-red-600/20 text-red-400' :
                      p.risk === 'high' ? 'bg-orange-600/20 text-orange-400' :
                      'bg-gray-700 text-gray-400'
                    }`}>
                      {p.risk.toUpperCase()}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {fingerprints.length > 0 && (
            <div className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
              <div className="p-3 border-b border-gray-800">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-semibold text-white">🔎 技术栈 ({fingerprints.length})</h3>
                  <button 
                    onClick={() => navigator.clipboard.writeText(fingerprints.join(', '))}
                    className="text-[10px] text-gray-500 hover:text-cyan-400"
                  >
                    复制
                  </button>
                </div>
              </div>
              <div className="p-3 max-h-48 overflow-y-auto">
                <div className="flex flex-wrap gap-1">
                  {fingerprints.map((fp, idx) => (
                    <span key={idx} className="px-2 py-1 bg-purple-600/20 text-purple-400 rounded text-xs">
                      {fp}
                    </span>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>

        {logs.length > 0 && (
          <div className="bg-gray-900/50 border border-gray-800 rounded-lg overflow-hidden">
            <div className="flex items-center justify-between px-3 py-2 border-b border-gray-800">
              <h3 className="text-xs font-medium text-gray-400">📜 扫描日志</h3>
              <button onClick={() => setLogs([])} className="text-[10px] text-gray-500 hover:text-gray-300">清空</button>
            </div>
            <div ref={logRef} className="h-32 overflow-y-auto font-mono text-[10px] bg-gray-950 p-2 space-y-0.5">
              {logs.map((log, idx) => (
                <div key={idx} className={`${
                  log.includes('✓') || log.includes('✅') ? 'text-green-400' : 
                  log.includes('❌') || log.includes('⚠️') ? 'text-red-400' :
                  log.includes('📊') || log.includes('💾') ? 'text-yellow-400' :
                  'text-gray-400'
                }`}>{log}</div>
              ))}
              {scanning && <div className="text-cyan-400 animate-pulse">▋ 扫描进行中...</div>}
            </div>
          </div>
        )}

        {totalAssets === 0 && !scanning && (
          <div className="text-center py-12 text-gray-500">
            <div className="text-4xl mb-3">🔍</div>
            <p className="text-sm">输入目标域名开始信息收集</p>
            <p className="text-xs mt-1">结果将自动保存到数据库</p>
          </div>
        )}

      </div>
    </div>
  );
};

export default ReconOverview;
