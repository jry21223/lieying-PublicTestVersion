import React, { useState, useEffect } from 'react';
import { targetsApi, assetsApi, reconApi, activeScanApi } from '../services/api';

interface Asset {
  id: string;
  name: string;
  value: string;
  type: string;
  status: string;
  importance: number;
  tags: string;
  description: string;
  first_seen: string;
  last_seen: string;
}

interface Target {
  id: string;
  name: string;
  value: string;
  type: string;
  description: string;
  tags: string;
  created_at: string;
}

interface ScanResult {
  id: string;
  target: string;
  url: string;
  type: string;
  severity: string;
  title: string;
  data: string;
  description: string;
  created_at: string;
}

const AssetOverview: React.FC = () => {
  const [targets, setTargets] = useState<Target[]>([]);
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(false);
  const [showAddModal, setShowAddModal] = useState(false);
  const [showImportModal, setShowImportModal] = useState(false);
  const [newTarget, setNewTarget] = useState({
    name: '',
    value: '',
    type: 'domain',
    description: '',
    tags: ''
  });
  const [importText, setImportText] = useState('');
  const [importType, setImportType] = useState<'domain' | 'ip' | 'url'>('domain');
  const [scanningTargets, setScanningTargets] = useState<string[]>([]);
  const [expandedTargets, setExpandedTargets] = useState<Set<string>>(new Set());
  const [targetResults, setTargetResults] = useState<Record<string, ScanResult[]>>({});
  const [loadingResults, setLoadingResults] = useState<Set<string>>(new Set());

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [targetsData, assetsData] = await Promise.all([
        targetsApi.list(),
        assetsApi.list()
      ]);
      setTargets(targetsData || []);
      setAssets(assetsData || []);
    } catch (error) {
      console.error('加载数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const toggleTarget = async (targetId: string, targetValue: string) => {
    const newExpanded = new Set(expandedTargets);
    
    if (newExpanded.has(targetId)) {
      newExpanded.delete(targetId);
      setExpandedTargets(newExpanded);
    } else {
      newExpanded.add(targetId);
      setExpandedTargets(newExpanded);
      
      if (!targetResults[targetId]) {
        setLoadingResults(prev => new Set(prev).add(targetId));
        try {
          const tasksData: any = await activeScanApi.listTasks();
          const tasks = tasksData?.data || tasksData || [];
          const matchingTasks = tasks.filter((t: any) => 
            t.target?.includes(targetValue) || 
            t.name?.includes(targetValue) ||
            targetValue.includes(t.target)
          );
          
          let allResults: ScanResult[] = [];
          for (const task of matchingTasks.slice(0, 5)) {
            try {
              const resultsData: any = await activeScanApi.getResults(task.id);
              const results = resultsData?.data || resultsData?.results || resultsData || [];
              allResults = [...allResults, ...results];
            } catch (e) {
              console.warn('获取任务结果失败:', e);
            }
          }
          
          setTargetResults(prev => ({
            ...prev,
            [targetId]: allResults
          }));
        } catch (error) {
          console.error('加载扫描结果失败:', error);
          setTargetResults(prev => ({
            ...prev,
            [targetId]: []
          }));
        } finally {
          setLoadingResults(prev => {
            const next = new Set(prev);
            next.delete(targetId);
            return next;
          });
        }
      }
    }
  };

  const handleAddTarget = async () => {
    if (!newTarget.name || !newTarget.value) {
      alert('请填写名称和值');
      return;
    }

    setLoading(true);
    try {
      await targetsApi.create({
        name: newTarget.name,
        value: newTarget.value,
        type: newTarget.type,
        description: newTarget.description,
        tags: newTarget.tags
      });
      alert('目标添加成功！');
      setShowAddModal(false);
      setNewTarget({ name: '', value: '', type: 'domain', description: '', tags: '' });
      loadData();
    } catch (error) {
      console.error('添加目标失败:', error);
      alert('添加失败！');
    } finally {
      setLoading(false);
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

    setLoading(true);
    try {
      for (const line of lines) {
        const value = line.trim();
        if (!value) continue;

        await targetsApi.create({
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
      loadData();
    } catch (error) {
      console.error('导入失败:', error);
      alert('导入失败！');
    } finally {
      setLoading(false);
    }
  };

  const handleStartScan = async (target: Target) => {
    if (scanningTargets.includes(target.id)) {
      return;
    }

    setScanningTargets([...scanningTargets, target.id]);
    setLoading(true);

    try {
      await reconApi.startRecon({
        target: target.value,
        type: target.type,
        options: {
          subdomain: true,
          portscan: true,
          fingerprint: true
        }
      });
      alert(`已为 ${target.name} 启动扫描任务！`);
    } catch (error) {
      console.error('启动扫描失败:', error);
      alert('启动扫描失败！');
    } finally {
      setLoading(false);
      setTimeout(() => {
        setScanningTargets(scanningTargets.filter(id => id !== target.id));
      }, 2000);
    }
  };

  const handleBatchScan = async () => {
    if (targets.length === 0) {
      alert('没有可扫描的目标');
      return;
    }

    setLoading(true);
    try {
      for (const target of targets.slice(0, 10)) {
        await reconApi.startRecon({
          target: target.value,
          type: target.type,
          options: {
            subdomain: true,
            portscan: true,
            fingerprint: true
          }
        });
      }
      alert(`已为前10个目标启动扫描任务！`);
    } catch (error) {
      console.error('批量扫描失败:', error);
      alert('批量扫描失败！');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteTarget = async (id: string) => {
    if (!confirm('确定要删除这个目标吗？')) return;

    setLoading(true);
    try {
      await targetsApi.delete(id);
      loadData();
    } catch (error) {
      console.error('删除失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const stats = {
    total: targets.length,
    domains: targets.filter(t => t.type === 'domain').length,
    ips: targets.filter(t => t.type === 'ip').length,
    urls: targets.filter(t => t.type === 'url').length,
  };

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'domain': return '🌐';
      case 'ip': return '💻';
      case 'url': return '🔗';
      default: return '📦';
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'domain': return 'text-blue-400';
      case 'ip': return 'text-green-400';
      case 'url': return 'text-purple-400';
      default: return 'text-gray-400';
    }
  };

  const getSeverityConfig = (severity: string) => {
    const configs: Record<string, { icon: string; color: string; bg: string; border: string }> = {
      critical: { icon: '🔴', color: 'text-red-400', bg: 'bg-red-600/10', border: 'border-red-600/30' },
      high: { icon: '🟠', color: 'text-orange-400', bg: 'bg-orange-600/10', border: 'border-orange-600/30' },
      medium: { icon: '🟡', color: 'text-yellow-400', bg: 'bg-yellow-600/10', border: 'border-yellow-600/30' },
      low: { icon: '🔵', color: 'text-blue-400', bg: 'bg-blue-600/10', border: 'border-blue-600/30' },
      info: { icon: '⚪', color: 'text-gray-400', bg: 'bg-gray-600/10', border: 'border-gray-600/30' },
    };
    return configs[severity?.toLowerCase()] || configs.info;
  };

  const getTotalVulnCount = () => {
    return Object.values(targetResults).reduce((sum, results) => sum + results.length, 0);
  };

  const getCriticalCount = () => {
    return Object.values(targetResults).reduce((sum, results) => 
      sum + results.filter(r => r.severity === 'critical' || r.severity === 'high').length, 0
    );
  };

  return (
    <div className="h-full overflow-auto bg-gray-950">
      <header className="h-14 bg-gray-900 border-b border-gray-800 flex items-center justify-between px-6 sticky top-0 z-10">
        <div className="flex items-center gap-2">
          <h2 className="font-semibold text-gray-100">资产概览</h2>
          <span className="text-xs text-gray-500">|</span>
          <span className="text-xs text-gray-400">{stats.total} 目标 · {getTotalVulnCount()} 发现</span>
          {getCriticalCount() > 0 && (
            <span className="px-2 py-0.5 bg-red-600/20 text-red-400 rounded text-xs">
              {getCriticalCount()} 高危
            </span>
          )}
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => setShowImportModal(true)}
            className="flex items-center gap-1.5 px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white text-sm rounded-lg transition"
          >
            📥 导入
          </button>
          <button
            onClick={() => setShowAddModal(true)}
            className="flex items-center gap-1.5 px-3 py-1.5 bg-green-600 hover:bg-green-500 text-white text-sm rounded-lg transition"
          >
            ➕ 添加
          </button>
          <button
            onClick={handleBatchScan}
            disabled={loading || targets.length === 0}
            className="flex items-center gap-1.5 px-3 py-1.5 bg-purple-600 hover:bg-purple-500 disabled:opacity-50 text-white text-sm rounded-lg transition"
          >
            🚀 批量扫描
          </button>
        </div>
      </header>

      <div className="p-4">
        <div className="grid grid-cols-4 gap-3 mb-4">
          <div className="bg-gray-900 rounded-lg border border-gray-800 p-4">
            <p className="text-xs text-gray-500">目标总数</p>
            <p className="text-2xl font-bold text-white mt-1">{stats.total}</p>
          </div>
          <div className="bg-gray-900 rounded-lg border border-gray-800 p-4">
            <p className="text-xs text-gray-500">域名</p>
            <p className="text-2xl font-bold text-blue-400 mt-1">{stats.domains}</p>
          </div>
          <div className="bg-gray-900 rounded-lg border border-gray-800 p-4">
            <p className="text-xs text-gray-500">IP地址</p>
            <p className="text-2xl font-bold text-green-400 mt-1">{stats.ips}</p>
          </div>
          <div className="bg-gray-900 rounded-lg border border-gray-800 p-4">
            <p className="text-xs text-gray-500">URL地址</p>
            <p className="text-2xl font-bold text-purple-400 mt-1">{stats.urls}</p>
          </div>
        </div>

        {loading ? (
          <div className="text-center py-12 text-gray-500">
            <div className="animate-spin text-4xl mb-4">⏳</div>
            <p>加载中...</p>
          </div>
        ) : targets.length === 0 ? (
          <div className="text-center py-12 text-gray-500 bg-gray-900 rounded-xl border border-gray-800">
            <div className="text-4xl mb-4">📦</div>
            <p>暂无目标资产</p>
            <p className="text-sm mt-2">点击"添加"或"导入"开始</p>
          </div>
        ) : (
          <div className="space-y-2">
            {targets.map((target) => {
              const isExpanded = expandedTargets.has(target.id);
              const results = targetResults[target.id] || [];
              const isLoading = loadingResults.has(target.id);
              const criticalCount = results.filter(r => r.severity === 'critical' || r.severity === 'high').length;
              
              return (
                <div key={target.id} className="bg-gray-900 rounded-lg border border-gray-800 overflow-hidden">
                  <div 
                    className={`p-3 flex items-center justify-between cursor-pointer hover:bg-gray-800/50 transition ${isExpanded ? 'border-b border-gray-800' : ''}`}
                    onClick={() => toggleTarget(target.id, target.value)}
                  >
                    <div className="flex items-center gap-3 flex-1 min-w-0">
                      <span className={`transition-transform ${isExpanded ? 'rotate-90' : ''}`}>
                        ▶
                      </span>
                      <span className="text-xl">{getTypeIcon(target.type)}</span>
                      <div className="flex-1 min-w-0">
                        <div className="flex items-center gap-2">
                          <span className={`font-mono ${getTypeColor(target.type)} truncate`}>
                            {target.value}
                          </span>
                          <span className="text-xs text-gray-600 bg-gray-800 px-1.5 py-0.5 rounded">
                            {target.type}
                          </span>
                          {target.tags && (
                            <span className="text-xs text-gray-500 bg-gray-800 px-1.5 py-0.5 rounded">
                              {target.tags}
                            </span>
                          )}
                        </div>
                        {target.description && (
                          <p className="text-xs text-gray-500 truncate mt-0.5">{target.description}</p>
                        )}
                      </div>
                    </div>
                    
                    <div className="flex items-center gap-2">
                      {results.length > 0 && (
                        <div className="flex items-center gap-1 mr-2">
                          <span className="text-xs text-gray-400">{results.length} 发现</span>
                          {criticalCount > 0 && (
                            <span className="px-1.5 py-0.5 bg-red-600/20 text-red-400 rounded text-xs">
                              {criticalCount} 高危
                            </span>
                          )}
                        </div>
                      )}
                      
                      <button
                        onClick={(e) => { e.stopPropagation(); handleStartScan(target); }}
                        disabled={scanningTargets.includes(target.id)}
                        className="px-2 py-1 bg-blue-600/20 text-blue-400 hover:bg-blue-600/30 rounded text-xs transition disabled:opacity-50"
                      >
                        {scanningTargets.includes(target.id) ? '⏳' : '🚀'}
                      </button>
                      <button
                        onClick={(e) => { e.stopPropagation(); handleDeleteTarget(target.id); }}
                        className="px-2 py-1 bg-red-600/20 text-red-400 hover:bg-red-600/30 rounded text-xs transition"
                      >
                        🗑️
                      </button>
                    </div>
                  </div>
                  
                  {isExpanded && (
                    <div className="bg-gray-950/50">
                      {isLoading ? (
                        <div className="p-4 text-center text-gray-500">
                          <div className="animate-spin inline-block mr-2">⏳</div>
                          加载扫描结果...
                        </div>
                      ) : results.length > 0 ? (
                        <div className="p-2 space-y-1">
                          {results.map((result) => {
                            const sevConfig = getSeverityConfig(result.severity);
                            return (
                              <div 
                                key={result.id}
                                className={`p-2.5 rounded-lg ${sevConfig.bg} border ${sevConfig.border} flex items-center justify-between`}
                              >
                                <div className="flex items-center gap-2 flex-1 min-w-0">
                                  <span>{sevConfig.icon}</span>
                                  <span className={`text-xs px-1.5 py-0.5 rounded ${sevConfig.color} ${sevConfig.bg}`}>
                                    {result.severity?.toUpperCase() || 'INFO'}
                                  </span>
                                  <span className="text-sm text-white truncate flex-1">
                                    {result.title || result.data?.substring(0, 60) || '扫描发现'}
                                  </span>
                                </div>
                                <div className="flex items-center gap-2">
                                  <span className="text-xs text-gray-500 font-mono truncate max-w-[150px]" title={result.url || result.target}>
                                    {result.url || result.target}
                                  </span>
                                  <button
                                    onClick={() => navigator.clipboard.writeText(result.url || result.target || '')}
                                    className="text-xs px-1.5 py-0.5 bg-gray-700/50 text-gray-400 rounded hover:text-white"
                                  >
                                    复制
                                  </button>
                                </div>
                              </div>
                            );
                          })}
                        </div>
                      ) : (
                        <div className="p-4 text-center text-gray-500 text-sm">
                          暂无扫描结果，点击 🚀 开始扫描
                        </div>
                      )}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>

      {showAddModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold text-white">添加目标资产</h3>
              <button onClick={() => setShowAddModal(false)} className="text-gray-400 hover:text-white text-xl">&times;</button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">名称 *</label>
                <input
                  type="text"
                  value={newTarget.name}
                  onChange={(e) => setNewTarget({ ...newTarget, name: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-cyan-500 outline-none"
                  placeholder="例如：VIPC6资源网"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">值 *</label>
                <input
                  type="text"
                  value={newTarget.value}
                  onChange={(e) => setNewTarget({ ...newTarget, value: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-cyan-500 outline-none"
                  placeholder="例如：vipc6.com 或 192.168.1.1"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">类型</label>
                <select
                  value={newTarget.type}
                  onChange={(e) => setNewTarget({ ...newTarget, type: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-cyan-500 outline-none"
                >
                  <option value="domain">域名</option>
                  <option value="ip">IP地址</option>
                  <option value="url">URL地址</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">描述</label>
                <input
                  type="text"
                  value={newTarget.description}
                  onChange={(e) => setNewTarget({ ...newTarget, description: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-cyan-500 outline-none"
                  placeholder="可选描述"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">标签</label>
                <input
                  type="text"
                  value={newTarget.tags}
                  onChange={(e) => setNewTarget({ ...newTarget, tags: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-cyan-500 outline-none"
                  placeholder="例如：重要,测试 (逗号分隔)"
                />
              </div>
            </div>
            <div className="flex justify-end space-x-3 mt-6">
              <button onClick={() => setShowAddModal(false)} className="px-4 py-2 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 transition">取消</button>
              <button onClick={handleAddTarget} disabled={loading} className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-500 transition disabled:opacity-50">{loading ? '添加中...' : '添加'}</button>
            </div>
          </div>
        </div>
      )}

      {showImportModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 w-full max-w-2xl">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold text-white">批量导入目标</h3>
              <button onClick={() => setShowImportModal(false)} className="text-gray-400 hover:text-white text-xl">&times;</button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">目标类型</label>
                <select
                  value={importType}
                  onChange={(e) => setImportType(e.target.value as any)}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-cyan-500 outline-none"
                >
                  <option value="domain">域名</option>
                  <option value="ip">IP地址</option>
                  <option value="url">URL地址</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">导入内容 (每行一个)</label>
                <textarea
                  value={importText}
                  onChange={(e) => setImportText(e.target.value)}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white h-48 font-mono text-sm focus:border-cyan-500 outline-none"
                  placeholder={importType === 'domain' ? `vipc6.com\nexample.com\nbaidu.com` : importType === 'ip' ? `192.168.1.1\n10.0.0.1\n8.8.8.8` : `https://example.com\nhttps://baidu.com`}
                />
              </div>
              <div className="bg-gray-800/50 rounded-lg p-3 flex items-center justify-between">
                <span className="text-sm text-gray-400">💡 每行一个目标</span>
                <span className="text-sm text-cyan-400">{importText.split('\n').filter(l => l.trim()).length} 行</span>
              </div>
            </div>
            <div className="flex justify-end space-x-3 mt-6">
              <button onClick={() => setShowImportModal(false)} className="px-4 py-2 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 transition">取消</button>
              <button onClick={handleImport} disabled={loading} className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500 transition disabled:opacity-50">{loading ? '导入中...' : '开始导入'}</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AssetOverview;
