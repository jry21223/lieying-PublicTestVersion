import React, { useState, useEffect } from 'react';
import { activeScanApi, vulnerabilitiesApi } from '../services/api';

interface ScanResult {
  id: string;
  task_id: string;
  target: string;
  url: string;
  type: string;
  severity: 'critical' | 'high' | 'medium' | 'low' | 'info';
  title: string;
  description: string;
  data: string;
  payload: string;
  request: string;
  response: string;
  reproduce_steps: string;
  fix_suggestion: string;
  references: string;
  cve_id: string;
  cnvd_id: string;
  tags: string;
  proof: string;
  risk_level: number;
  confirmed: boolean;
  created_at: string;
}

const VulnResultsPanel: React.FC = () => {
  const [results, setResults] = useState<ScanResult[]>([]);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState('all');
  const [autoRefresh, setAutoRefresh] = useState(false);
  const [selectedResult, setSelectedResult] = useState<ScanResult | null>(null);

  useEffect(() => {
    loadResults();
  }, []);

  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(loadResults, 5000);
      return () => clearInterval(interval);
    }
  }, [autoRefresh]);

  const loadResults = async () => {
    try {
      const tasks = await activeScanApi.listTasks();
      const allResults: ScanResult[] = [];
      
      if (Array.isArray(tasks)) {
        for (const task of tasks) {
          try {
            const taskResults = await activeScanApi.getResults(task.id);
            if (Array.isArray(taskResults)) {
              allResults.push(...taskResults);
            }
          } catch (e) {
            // 忽略单个任务的错误
          }
        }
      }
      
      allResults.sort((a, b) => {
        const timeA = a.created_at ? new Date(a.created_at).getTime() : 0;
        const timeB = b.created_at ? new Date(b.created_at).getTime() : 0;
        return timeB - timeA;
      });
      
      setResults(allResults);
    } catch (error) {
      console.error('加载扫描结果失败:', error);
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  const addToVulnerabilities = async (result: ScanResult) => {
    try {
      await vulnerabilitiesApi.create({
        title: result.title || result.data,
        severity: result.severity,
        url: result.url || result.target,
        status: 'pending',
      });
      alert('已添加到漏洞列表！');
    } catch (error) {
      console.error('添加漏洞失败:', error);
      alert('添加失败！');
    }
  };

  const getSeverityColor = (severity: string) => {
    const colors: Record<string, string> = {
      critical: 'bg-red-600/20 text-red-400 border-red-600/30',
      high: 'bg-orange-600/20 text-orange-400 border-orange-600/30',
      medium: 'bg-yellow-600/20 text-yellow-400 border-yellow-600/30',
      low: 'bg-blue-600/20 text-blue-400 border-blue-600/30',
      info: 'bg-gray-600/20 text-gray-400 border-gray-600/30',
    };
    return colors[severity] || colors.info;
  };

  const getSeverityIcon = (severity: string) => {
    const icons: Record<string, string> = {
      critical: '🔴',
      high: '🟠',
      medium: '🟡',
      low: '🔵',
      info: '⚪',
    };
    return icons[severity] || '⚪';
  };

  const getTypeLabel = (type: string) => {
    const labels: Record<string, string> = {
      'web_scan': 'Web扫描',
      'missing_header': '安全头',
      'tech_detection': '技术栈',
      'sensitive_info': '敏感信息',
      'sensitive_page': '敏感页面',
      'cors_misconfig': 'CORS',
      'cors_config': 'CORS',
      'cookie_security': 'Cookie',
      'csrf_token': 'CSRF',
      'insecure_form': '表单安全',
      'mixed_content': '混合内容',
      'vuln_detected': '漏洞',
      'sql_injection': 'SQL注入',
      'xss': 'XSS',
      'directory_traversal': '目录遍历',
      'xxe': 'XXE',
      'ssrf': 'SSRF',
      'auth_required': '认证',
      'server_error': '服务器错误',
    };
    return labels[type] || type;
  };

  const getTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      'sql_injection': 'bg-red-600/20 text-red-400 border border-red-600/30',
      'xss': 'bg-red-600/20 text-red-400 border border-red-600/30',
      'directory_traversal': 'bg-red-600/20 text-red-400 border border-red-600/30',
      'xxe': 'bg-red-600/20 text-red-400 border border-red-600/30',
      'ssrf': 'bg-red-600/20 text-red-400 border border-red-600/30',
      'vuln_detected': 'bg-red-600/20 text-red-400 border border-red-600/30',
      'missing_header': 'bg-yellow-600/20 text-yellow-400 border border-yellow-600/30',
      'cors_misconfig': 'bg-orange-600/20 text-orange-400 border border-orange-600/30',
      'cookie_security': 'bg-yellow-600/20 text-yellow-400 border border-yellow-600/30',
      'csrf_token': 'bg-orange-600/20 text-orange-400 border border-orange-600/30',
      'sensitive_info': 'bg-red-600/20 text-red-400 border border-red-600/30',
      'sensitive_page': 'bg-yellow-600/20 text-yellow-400 border border-yellow-600/30',
      'tech_detection': 'bg-cyan-600/20 text-cyan-400 border border-cyan-600/30',
    };
    return colors[type] || 'bg-gray-600/20 text-gray-400 border border-gray-600/30';
  };

  const filteredResults = filter === 'all' 
    ? results 
    : results.filter(r => r.severity === filter);

  const stats = {
    critical: results.filter(r => r.severity === 'critical').length,
    high: results.filter(r => r.severity === 'high').length,
    medium: results.filter(r => r.severity === 'medium').length,
    low: results.filter(r => r.severity === 'low').length,
    info: results.filter(r => r.severity === 'info').length,
  };

  const extractVulnInfo = (result: ScanResult): { title: string; description: string } => {
    if (result.title && result.title !== result.data) {
      return { title: result.title, description: result.description || '' };
    }
    
    const data = result.data || '';
    
    if (data.includes('发现登录页面')) {
      return { title: '发现登录页面', description: '目标站点存在登录入口，可用于后续的暴力破解、弱口令测试、SQL注入、XSS等漏洞测试' };
    }
    if (data.includes('发现管理员关键字')) {
      return { title: '发现管理后台关键字', description: '页面中包含admin关键字，可能存在管理后台入口' };
    }
    if (data.includes('发现密码字段')) {
      return { title: '发现密码字段', description: '页面中包含password关键字，可能存在密码输入表单' };
    }
    if (data.includes('发现管理后台关键字')) {
      return { title: '发现管理后台', description: '页面中包含dashboard或console关键字，可能存在管理后台' };
    }
    if (data.includes('发现默认首页')) {
      return { title: '发现默认首页', description: '目标使用默认欢迎页面，可能泄露服务器信息' };
    }
    
    const bracketMatch = data.match(/\[.*?\]\s*(.+)/);
    if (bracketMatch) {
      return { title: bracketMatch[1].split(' - ')[0] || bracketMatch[1], description: result.description || '' };
    }
    
    if (data.length > 100) {
      return { title: getTypeLabel(result.type), description: data.substring(0, 100) + '...' };
    }
    
    return { title: data || getTypeLabel(result.type), description: result.description || '' };
  };

  if (loading) {
    return (
      <div className="h-full flex items-center justify-center bg-gray-950">
        <div className="text-gray-400">
          <div className="animate-spin text-4xl mb-4">⏳</div>
          <p>加载扫描结果...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-6xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-white">漏洞扫描结果</h2>
          <div className="flex items-center gap-3">
            <label className="flex items-center gap-2 px-3 py-1.5 bg-gray-800 rounded-lg text-sm cursor-pointer">
              <input
                type="checkbox"
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
                className="w-4 h-4 rounded bg-gray-700 border-gray-600"
              />
              <span className="text-gray-300">自动刷新</span>
            </label>
            <button
              onClick={loadResults}
              className="px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition"
            >
              刷新
            </button>
          </div>
        </div>

        {results.length > 0 && (
          <>
            <div className="grid grid-cols-5 gap-4 mb-6">
              <div className="bg-red-900/20 border border-red-800/30 rounded-xl p-4">
                <p className="text-sm text-gray-400">严重</p>
                <p className="text-2xl font-bold text-red-500">{stats.critical}</p>
              </div>
              <div className="bg-orange-900/20 border border-orange-800/30 rounded-xl p-4">
                <p className="text-sm text-gray-400">高危</p>
                <p className="text-2xl font-bold text-orange-500">{stats.high}</p>
              </div>
              <div className="bg-yellow-900/20 border border-yellow-800/30 rounded-xl p-4">
                <p className="text-sm text-gray-400">中危</p>
                <p className="text-2xl font-bold text-yellow-500">{stats.medium}</p>
              </div>
              <div className="bg-blue-900/20 border border-blue-800/30 rounded-xl p-4">
                <p className="text-sm text-gray-400">低危</p>
                <p className="text-2xl font-bold text-blue-500">{stats.low}</p>
              </div>
              <div className="bg-gray-800/50 border border-gray-700/30 rounded-xl p-4">
                <p className="text-sm text-gray-400">信息</p>
                <p className="text-2xl font-bold text-gray-400">{stats.info}</p>
              </div>
            </div>

            <div className="flex items-center gap-2 mb-4">
              <span className="text-sm text-gray-400">筛选：</span>
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
                  className={`px-3 py-1.5 rounded-lg text-sm transition-colors ${
                    filter === key
                      ? 'bg-blue-600 text-white'
                      : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
                  }`}
                >
                  {label}
                </button>
              ))}
            </div>
          </>
        )}

        {filteredResults.length > 0 ? (
          <div className="space-y-3">
            {filteredResults.map((result) => {
              const vulnInfo = extractVulnInfo(result);
              return (
                <div 
                  key={result.id} 
                  className="bg-gray-900 rounded-lg border border-gray-800 p-4 hover:border-gray-700 transition cursor-pointer"
                  onClick={() => setSelectedResult(result)}
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <span className="text-lg">{getSeverityIcon(result.severity)}</span>
                      <span className={`px-2 py-0.5 rounded text-xs border ${getSeverityColor(result.severity)}`}>
                        {result.severity.toUpperCase()}
                      </span>
                      <span className={`px-2 py-0.5 rounded text-xs ${getTypeColor(result.type)}`}>
                        {getTypeLabel(result.type)}
                      </span>
                      {result.cve_id && (
                        <span className="px-2 py-0.5 rounded text-xs bg-red-600/20 text-red-400 border border-red-600/30">
                          {result.cve_id}
                        </span>
                      )}
                    </div>
                    <div className="flex items-center gap-2">
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          navigator.clipboard.writeText(result.url || result.target || '');
                        }}
                        className="text-xs px-2 py-1 bg-gray-600/20 text-gray-400 rounded hover:bg-gray-600/30"
                      >
                        复制URL
                      </button>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          addToVulnerabilities(result);
                        }}
                        className="text-xs px-3 py-1 bg-green-600/20 text-green-400 rounded hover:bg-green-600/30"
                      >
                        添加到漏洞
                      </button>
                      <span className="text-xs text-gray-500">
                        {result.created_at ? new Date(result.created_at).toLocaleString() : ''}
                      </span>
                    </div>
                  </div>
                  <p className="text-sm text-blue-400 font-mono mb-1 truncate hover:underline">
                    🔗 {result.url || result.target}
                  </p>
                  <p className="text-white font-medium">{vulnInfo.title}</p>
                  {vulnInfo.description && (
                    <p className="text-sm text-gray-400 mt-1">{vulnInfo.description}</p>
                  )}
                  <div className="mt-2 text-xs text-gray-500 flex items-center gap-2">
                    <span>点击查看详情 →</span>
                  </div>
                </div>
              );
            })}
          </div>
        ) : (
          <div className="text-center py-12 text-gray-500">
            <div className="text-4xl mb-4">📋</div>
            <p>暂无扫描结果</p>
            <p className="text-sm mt-2">请在"主动扫描"模块创建扫描任务</p>
          </div>
        )}

        {selectedResult && (
          <div 
            className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4"
            onClick={() => setSelectedResult(null)}
          >
            <div 
              className="bg-gray-900 rounded-xl border border-gray-700 w-full max-w-4xl max-h-[90vh] overflow-hidden"
              onClick={(e) => e.stopPropagation()}
            >
              <div className="p-4 border-b border-gray-800 flex items-center justify-between bg-gray-800/50">
                <div className="flex items-center gap-3">
                  <span className="text-2xl">{getSeverityIcon(selectedResult.severity)}</span>
                  <div>
                    <h3 className="text-lg font-semibold text-white">
                      {selectedResult.title || extractVulnInfo(selectedResult).title}
                    </h3>
                    <div className="flex items-center gap-2 mt-1">
                      <span className={`px-2 py-0.5 rounded text-xs ${getSeverityColor(selectedResult.severity)}`}>
                        {selectedResult.severity.toUpperCase()}
                      </span>
                      <span className={`px-2 py-0.5 rounded text-xs ${getTypeColor(selectedResult.type)}`}>
                        {getTypeLabel(selectedResult.type)}
                      </span>
                      {selectedResult.cve_id && (
                        <span className="px-2 py-0.5 rounded text-xs bg-red-600/20 text-red-400 border border-red-600/30">
                          {selectedResult.cve_id}
                        </span>
                      )}
                      {selectedResult.cnvd_id && (
                        <span className="px-2 py-0.5 rounded text-xs bg-orange-600/20 text-orange-400 border border-orange-600/30">
                          {selectedResult.cnvd_id}
                        </span>
                      )}
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={() => {
                      navigator.clipboard.writeText(selectedResult.url || selectedResult.target || '');
                    }}
                    className="px-3 py-1.5 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 text-sm"
                  >
                    复制URL
                  </button>
                  <button
                    onClick={() => setSelectedResult(null)}
                    className="text-gray-400 hover:text-white text-xl"
                  >
                    ✕
                  </button>
                </div>
              </div>

              <div className="overflow-auto max-h-[calc(90vh-80px)]">
                <div className="p-4 space-y-4">
                  <div className="bg-gray-800/50 rounded-lg p-4">
                    <h4 className="text-sm font-medium text-gray-400 mb-2 flex items-center gap-2">
                      <span>🎯</span> 漏洞地址
                    </h4>
                    <a 
                      href={selectedResult.url || selectedResult.target} 
                      target="_blank" 
                      rel="noopener noreferrer"
                      className="text-blue-400 hover:text-blue-300 font-mono text-sm break-all hover:underline"
                    >
                      {selectedResult.url || selectedResult.target}
                    </a>
                  </div>

                  {selectedResult.description && (
                    <div className="bg-gray-800/50 rounded-lg p-4">
                      <h4 className="text-sm font-medium text-gray-400 mb-2 flex items-center gap-2">
                        <span>📝</span> 漏洞描述
                      </h4>
                      <p className="text-gray-300 text-sm">{selectedResult.description}</p>
                    </div>
                  )}

                  {selectedResult.reproduce_steps && (
                    <div className="bg-gray-800/50 rounded-lg p-4">
                      <h4 className="text-sm font-medium text-gray-400 mb-2 flex items-center gap-2">
                        <span>🔬</span> 复现步骤
                      </h4>
                      <pre className="text-gray-300 text-sm whitespace-pre-wrap bg-gray-900 p-3 rounded overflow-x-auto">
                        {selectedResult.reproduce_steps}
                      </pre>
                    </div>
                  )}

                  {selectedResult.payload && (
                    <div className="bg-gray-800/50 rounded-lg p-4">
                      <h4 className="text-sm font-medium text-gray-400 mb-2 flex items-center gap-2">
                        <span>💉</span> Payload
                      </h4>
                      <pre className="text-green-400 text-sm font-mono bg-gray-900 p-3 rounded overflow-x-auto">
                        {selectedResult.payload}
                      </pre>
                    </div>
                  )}

                  {selectedResult.request && (
                    <div className="bg-gray-800/50 rounded-lg p-4">
                      <h4 className="text-sm font-medium text-gray-400 mb-2 flex items-center gap-2">
                        <span>📤</span> 请求包
                      </h4>
                      <pre className="text-cyan-400 text-xs font-mono bg-gray-900 p-3 rounded overflow-x-auto max-h-48">
                        {selectedResult.request}
                      </pre>
                    </div>
                  )}

                  {selectedResult.response && (
                    <div className="bg-gray-800/50 rounded-lg p-4">
                      <h4 className="text-sm font-medium text-gray-400 mb-2 flex items-center gap-2">
                        <span>📥</span> 响应包
                      </h4>
                      <pre className="text-yellow-400 text-xs font-mono bg-gray-900 p-3 rounded overflow-x-auto max-h-48">
                        {selectedResult.response.length > 2000 ? selectedResult.response.substring(0, 2000) + '\n...(已截断)' : selectedResult.response}
                      </pre>
                    </div>
                  )}

                  {selectedResult.fix_suggestion && (
                    <div className="bg-green-900/20 rounded-lg p-4 border border-green-600/30">
                      <h4 className="text-sm font-medium text-green-400 mb-2 flex items-center gap-2">
                        <span>🛡️</span> 修复建议
                      </h4>
                      <pre className="text-green-300 text-sm whitespace-pre-wrap">
                        {selectedResult.fix_suggestion}
                      </pre>
                    </div>
                  )}

                  {selectedResult.references && (
                    <div className="bg-gray-800/50 rounded-lg p-4">
                      <h4 className="text-sm font-medium text-gray-400 mb-2 flex items-center gap-2">
                        <span>📚</span> 参考资料
                      </h4>
                      <div className="space-y-1">
                        {selectedResult.references.split('\n').map((ref, idx) => (
                          <a
                            key={idx}
                            href={ref}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-blue-400 hover:text-blue-300 text-sm block hover:underline"
                          >
                            {ref}
                          </a>
                        ))}
                      </div>
                    </div>
                  )}

                  {selectedResult.proof && (
                    <div className="bg-gray-800/50 rounded-lg p-4">
                      <h4 className="text-sm font-medium text-gray-400 mb-2 flex items-center gap-2">
                        <span>📷</span> 漏洞证明
                      </h4>
                      <pre className="text-gray-300 text-sm whitespace-pre-wrap bg-gray-900 p-3 rounded">
                        {selectedResult.proof}
                      </pre>
                    </div>
                  )}

                  <div className="grid grid-cols-2 gap-4">
                    {selectedResult.tags && (
                      <div className="bg-gray-800/50 rounded-lg p-4">
                        <h4 className="text-sm font-medium text-gray-400 mb-1">标签</h4>
                        <div className="flex flex-wrap gap-1">
                          {selectedResult.tags.split(',').map((tag, idx) => (
                            <span key={idx} className="px-2 py-0.5 bg-gray-700 text-gray-300 rounded text-xs">
                              {tag.trim()}
                            </span>
                          ))}
                        </div>
                      </div>
                    )}
                    {selectedResult.risk_level && (
                      <div className="bg-gray-800/50 rounded-lg p-4">
                        <h4 className="text-sm font-medium text-gray-400 mb-1">风险评分</h4>
                        <div className="flex items-center gap-2">
                          <div className="flex-1 bg-gray-700 rounded-full h-2">
                            <div
                              className={`h-2 rounded-full ${
                                selectedResult.risk_level >= 8 ? 'bg-red-500' :
                                selectedResult.risk_level >= 6 ? 'bg-orange-500' :
                                selectedResult.risk_level >= 4 ? 'bg-yellow-500' : 'bg-blue-500'
                              }`}
                              style={{ width: `${(selectedResult.risk_level / 10) * 100}%` }}
                            />
                          </div>
                          <span className="text-white text-sm font-medium">{selectedResult.risk_level}/10</span>
                        </div>
                      </div>
                    )}
                  </div>

                  <div className="flex justify-end gap-3 pt-4 border-t border-gray-800">
                    <a
                      href={selectedResult.url || selectedResult.target}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500 transition"
                    >
                      在浏览器中打开
                    </a>
                    <button
                      onClick={() => {
                        addToVulnerabilities(selectedResult);
                        setSelectedResult(null);
                      }}
                      className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-500 transition"
                    >
                      添加到漏洞列表
                    </button>
                    <button
                      onClick={() => setSelectedResult(null)}
                      className="px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition"
                    >
                      关闭
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default VulnResultsPanel;
