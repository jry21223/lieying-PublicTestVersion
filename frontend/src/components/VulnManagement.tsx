import React, { useState, useEffect, useCallback } from 'react';
import { vulnerabilitiesApi } from '../services/api';

interface Vulnerability {
  id: string;
  task_id?: string;
  type: string;
  target: string;
  url: string;
  title: string;
  description: string;
  data: string;
  severity: string;
  confirmed: boolean;
  cve_id?: string;
  cnvd_id?: string;
  payload?: string;
  request?: string;
  response?: string;
  reproduce_steps?: string;
  fix_suggestion?: string;
  references?: string;
  tags?: string;
  proof?: string;
  risk_level: number;
  affected_version?: string;
  created_at: string;
  updated_at: string;
}

const SEVERITY_CONFIG: Record<string, { color: string; bg: string; icon: string; label: string }> = {
  CRITICAL: { color: 'text-red-400', bg: 'bg-red-500/10 border-red-500/30', icon: '☠️', label: '严重' },
  HIGH: { color: 'text-orange-400', bg: 'bg-orange-500/10 border-orange-500/30', icon: '🔴', label: '高危' },
  MEDIUM: { color: 'text-yellow-400', bg: 'bg-yellow-500/10 border-yellow-500/30', icon: '🟡', label: '中危' },
  LOW: { color: 'text-blue-400', bg: 'bg-blue-500/10 border-blue-500/30', icon: '🔵', label: '低危' },
  INFO: { color: 'text-gray-400', bg: 'bg-gray-500/10 border-gray-500/30', icon: 'ℹ️', label: '信息' },
};

const TYPE_LABELS: Record<string, string> = {
  sql_injection: 'SQL注入',
  xss: 'XSS跨站',
  directory_traversal: '目录遍历',
  xxe: 'XXE注入',
  ssrf: 'SSRF服务端请求伪造',
  security_headers: '安全头缺失',
  sensitive_info: '敏感信息泄露',
  cors_misconfig: 'CORS配置错误',
  cookie_security: 'Cookie安全问题',
  tech_stack: '技术栈指纹',
  form_security: '表单安全问题',
  csrf: 'CSRF跨站请求伪造',
  open_redirect: '开放重定向',
  clickjacking: '点击劫持',
  idor: '越权访问(IDOR)',
  file_upload: '文件上传漏洞',
  command_injection: '命令注入',
  ssti: '服务端模板注入(SSTI)',
  deserialization: '反序列化漏洞',
  weak_password: '弱口令',
  default_config: '默认配置',
};

const VulnManagement: React.FC = () => {
  const [vulns, setVulns] = useState<Vulnerability[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedVuln, setSelectedVuln] = useState<Vulnerability | null>(null);
  const [filterSeverity, setFilterSeverity] = useState<string>('all');
  const [filterType, setFilterType] = useState<string>('all');
  const [searchKeyword, setSearchKeyword] = useState('');
  const [showConfirmedOnly, setShowConfirmedOnly] = useState(false);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const [actionLoading, setActionLoading] = useState(false);

  const loadVulnerabilities = useCallback(async () => {
    try {
      setLoading(true);
      const data = await vulnerabilitiesApi.list();
      setVulns(data.data || data || []);
    } catch (error) {
      console.error('加载漏洞失败:', error);
      setVulns([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadVulnerabilities();
  }, [loadVulnerabilities]);

  const filteredVulns = vulns.filter(v => {
    if (filterSeverity !== 'all' && v.severity !== filterSeverity) return false;
    if (filterType !== 'all' && v.type !== filterType) return false;
    if (showConfirmedOnly && !v.confirmed) return false;
    if (searchKeyword) {
      const keyword = searchKeyword.toLowerCase();
      return (
        v.title?.toLowerCase().includes(keyword) ||
        v.target?.toLowerCase().includes(keyword) ||
        v.url?.toLowerCase().includes(keyword) ||
        v.type?.toLowerCase().includes(keyword) ||
        v.cve_id?.toLowerCase().includes(keyword)
      );
    }
    return true;
  });

  const stats = {
    total: vulns.length,
    critical: vulns.filter(v => v.severity === 'CRITICAL').length,
    high: vulns.filter(v => v.severity === 'HIGH').length,
    medium: vulns.filter(v => v.severity === 'MEDIUM').length,
    low: vulns.filter(v => v.severity === 'LOW').length,
    info: vulns.filter(v => v.severity === 'INFO').length,
    confirmed: vulns.filter(v => v.confirmed).length,
  };

  const uniqueTypes = [...new Set(vulns.map(v => v.type))];

  const handleConfirm = async (id: string) => {
    try {
      await vulnerabilitiesApi.confirm(id);
      setVulns(prev => prev.map(v => v.id === id ? { ...v, confirmed: true } : v));
    } catch (error) {
      console.error('确认漏洞失败:', error);
    }
  };

  const handleBulkAction = async (action: string) => {
    if (selectedIds.size === 0) return;
    try {
      setActionLoading(true);
      await vulnerabilitiesApi.bulkAction({ ids: Array.from(selectedIds), action });
      if (action === 'confirm') {
        setVulns(prev => prev.map(v => selectedIds.has(v.id) ? { ...v, confirmed: true } : v));
      }
      setSelectedIds(new Set());
      await loadVulnerabilities();
    } catch (error) {
      console.error('批量操作失败:', error);
    } finally {
      setActionLoading(false);
    }
  };

  const toggleSelect = (id: string) => {
    setSelectedIds(prev => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id); else next.add(id);
      return next;
    });
  };

  const selectAll = () => {
    if (selectedIds.size === filteredVulns.length) {
      setSelectedIds(new Set());
    } else {
      setSelectedIds(new Set(filteredVulns.map(v => v.id)));
    }
  };

  const exportReport = () => {
    const reportData = filteredVulns.map(v => ({
      '漏洞标题': v.title || '-',
      '严重等级': SEVERITY_CONFIG[v.severity]?.label || v.severity,
      '漏洞类型': TYPE_LABELS[v.type] || v.type,
      '目标地址': v.target || v.url || '-',
      'CVE编号': v.cve_id || '-',
      'CNVD编号': v.cnvd_id || '-',
      '是否确认': v.confirmed ? '已确认' : '未确认',
      '发现时间': v.created_at,
      '风险评分': v.risk_level || '-',
      'Payload': v.payload?.substring(0, 100) || '-',
      '修复建议': v.fix_suggestion?.substring(0, 100) || '-',
    }));
    const headers = Object.keys(reportData[0] || {}).join('\t');
    const rows = reportData.map(row => Object.values(row).join('\t'));
    const blob = new Blob([headers + '\n' + rows.join('\n')], { type: 'text/plain;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `漏洞报告_${new Date().toISOString().slice(0, 10)}.tsv`;
    a.click();
    URL.revokeObjectURL(url);
  };

  if (loading) {
    return (
      <div className="h-full overflow-auto bg-gray-950 p-6 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin w-12 h-12 border-4 border-cyan-500 border-t-transparent rounded-full mx-auto mb-4"></div>
          <p className="text-gray-400">正在加载漏洞数据...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-7xl mx-auto space-y-6">

        {/* 统计面板 */}
        <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-7 gap-3">
          {[
            { label: '总计', value: stats.total, color: 'text-white', bg: 'bg-gray-800', icon: '📊' },
            { label: '严重', value: stats.critical, color: 'text-red-400', bg: 'bg-red-950/50', icon: '☠️' },
            { label: '高危', value: stats.high, color: 'text-orange-400', bg: 'bg-orange-950/50', icon: '🔴' },
            { label: '中危', value: stats.medium, color: 'text-yellow-400', bg: 'bg-yellow-950/50', icon: '🟡' },
            { label: '低危', value: stats.low, color: 'text-blue-400', bg: 'bg-blue-950/50', icon: '🔵' },
            { label: '信息', value: stats.info, color: 'text-gray-400', bg: 'bg-gray-800', icon: 'ℹ️' },
            { label: '已确认', value: stats.confirmed, color: 'text-green-400', bg: 'bg-green-950/50', icon: '✅' },
          ].map(s => (
            <div key={s.label} className={`${s.bg} rounded-lg border border-gray-800 p-3 text-center`}>
              <div className={`text-2xl font-bold ${s.color}`}>{s.value}</div>
              <div className="text-xs text-gray-500 mt-1">{s.icon} {s.label}</div>
            </div>
          ))}
        </div>

        {/* 工具栏 */}
        <div className="flex flex-wrap gap-3 items-center bg-gray-900 rounded-lg border border-gray-800 p-4">
          <input
            type="text"
            placeholder="🔍 搜索漏洞(标题/地址/CVE/类型)..."
            value={searchKeyword}
            onChange={(e) => setSearchKeyword(e.target.value)}
            className="flex-1 min-w-[200px] bg-gray-800 border border-gray-700 rounded-lg px-4 py-2 text-sm text-white placeholder-gray-500 focus:border-cyan-500 focus:outline-none"
          />
          <select
            value={filterSeverity}
            onChange={(e) => setFilterSeverity(e.target.value)}
            className="bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:border-cyan-500 focus:outline-none"
          >
            <option value="all">全部等级</option>
            <option value="CRITICAL">严重</option>
            <option value="HIGH">高危</option>
            <option value="MEDIUM">中危</option>
            <option value="LOW">低危</option>
            <option value="INFO">信息</option>
          </select>
          <select
            value={filterType}
            onChange={(e) => setFilterType(e.target.value)}
            className="bg-gray-800 border border-gray-700 rounded-lg px-3 py-2 text-sm text-white focus:border-cyan-500 focus:outline-none"
          >
            <option value="all">全部类型</option>
            {uniqueTypes.map(t => (
              <option key={t} value={t}>{TYPE_LABELS[t] || t}</option>
            ))}
          </select>
          <label className="flex items-center gap-2 text-sm text-gray-400 cursor-pointer">
            <input
              type="checkbox"
              checked={showConfirmedOnly}
              onChange={(e) => setShowConfirmedOnly(e.target.checked)}
              className="rounded bg-gray-700 border-gray-600"
            />
            仅已确认
          </label>
          <button onClick={exportReport} className="bg-cyan-600 hover:bg-cyan-500 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2">
            📥 导出报告
          </button>
          <button onClick={loadVulnerabilities} className="bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded-lg text-sm transition-colors">
            🔄 刷新
          </button>
        </div>

        {/* 批量操作栏 */}
        {selectedIds.size > 0 && (
          <div className="flex gap-3 items-center bg-cyan-950/30 border border-cyan-500/30 rounded-lg px-4 py-3">
            <span className="text-cyan-400 text-sm font-medium">已选 {selectedIds.size} 项</span>
            <button
              onClick={() => handleBulkAction('confirm')}
              disabled={actionLoading}
              className="bg-green-600 hover:bg-green-500 disabled:opacity-50 text-white px-3 py-1.5 rounded text-sm transition-colors"
            >
              ✅ 批量确认
            </button>
            <button
              onClick={() => handleBulkAction('ignore')}
              disabled={actionLoading}
              className="bg-gray-600 hover:bg-gray-500 disabled:opacity-50 text-white px-3 py-1.5 rounded text-sm transition-colors"
            >
              🚫 批量忽略
            </button>
            <button onClick={() => setSelectedIds(new Set())} className="text-gray-400 hover:text-white text-sm ml-auto">
              取消选择
            </button>
          </div>
        )}

        {/* 漏洞列表 */}
        {filteredVulns.length === 0 ? (
          <div className="text-center py-16 text-gray-500">
            <div className="text-6xl mb-4">🛡️</div>
            <p className="text-xl mb-2">暂无漏洞数据</p>
            <p className="text-sm">请先执行扫描任务，发现的安全漏洞将在此显示</p>
            {vulns.length > 0 && filteredVulns.length === 0 && (
              <button onClick={() => { setFilterSeverity('all'); setFilterType('all'); setSearchKeyword(''); setShowConfirmedOnly(false); }} className="mt-4 text-cyan-400 hover:text-cyan-300 text-sm underline">
                清除筛选条件
              </button>
            )}
          </div>
        ) : (
          <div className="space-y-2">
            {/* 表头 */}
            <div className="flex items-center gap-3 bg-gray-900 rounded-lg border border-gray-800 px-4 py-2 text-xs text-gray-500 uppercase tracking-wider">
              <input type="checkbox" checked={selectedIds.size === filteredVulns.length && filteredVulns.length > 0} onChange={selectAll} className="rounded" />
              <span className="w-20">等级</span>
              <span className="flex-1">漏洞标题</span>
              <span className="w-48 hidden md:block">目标地址</span>
              <span className="w-28 hidden lg:block">类型</span>
              <span className="w-20 text-center">状态</span>
              <span className="w-24 text-right">操作</span>
            </div>

            {/* 漏洞条目 */}
            {filteredVulns.map((vuln) => {
              const severityCfg = SEVERITY_CONFIG[vuln.severity] || SEVERITY_CONFIG.INFO;
              return (
                <div
                  key={vuln.id}
                  onClick={() => setSelectedVuln(vuln)}
                  className={`group flex items-center gap-3 bg-gray-900 rounded-lg border ${severityCfg.border} px-4 py-3 cursor-pointer hover:bg-gray-800/80 transition-all ${selectedIds.has(vuln.id) ? 'ring-2 ring-cyan-500' : ''}`}
                >
                  <input
                    type="checkbox"
                    checked={selectedIds.has(vuln.id)}
                    onChange={(e) => { e.stopPropagation(); toggleSelect(vuln.id); }}
                    className="rounded"
                    onClick={(e) => e.stopPropagation()}
                  />

                  {/* 等级标签 */}
                  <div className={`w-20 flex items-center gap-1 ${severityCfg.color} text-xs font-bold`}>
                    <span>{severityCfg.icon}</span>
                    <span>{severityCfg.label}</span>
                  </div>

                  {/* 标题 + CVE */}
                  <div className="flex-1 min-w-0">
                    <div className="text-white text-sm font-medium truncate group-hover:text-cyan-400 transition-colors">
                      {vuln.title || vuln.type || '未知漏洞'}
                    </div>
                    {(vuln.cve_id || vuln.cnvd_id) && (
                      <div className="flex gap-2 mt-0.5">
                        {vuln.cve_id && <span className="text-xs text-red-400 font-mono">{vuln.cve_id}</span>}
                        {vuln.cnvd_id && <span className="text-xs text-orange-400 font-mono">{vuln.cnvd_id}</span>}
                      </div>
                    )}
                  </div>

                  {/* 目标地址 */}
                  <div className="w-48 hidden md:block min-w-0">
                    <div className="text-xs text-gray-400 truncate font-mono" title={vuln.url || vuln.target}>
                      {vuln.url || vuln.target || '-'}
                    </div>
                  </div>

                  {/* 类型 */}
                  <div className="w-28 hidden lg:block">
                    <span className="text-xs text-gray-400 bg-gray-800 px-2 py-0.5 rounded">
                      {TYPE_LABELS[vuln.type] || vuln.type}
                    </span>
                  </div>

                  {/* 状态 */}
                  <div className="w-20 text-center">
                    {vuln.confirmed ? (
                      <span className="text-xs text-green-400 bg-green-500/10 px-2 py-0.5 rounded border border-green-500/30">✅ 已确认</span>
                    ) : (
                      <span className="text-xs text-gray-500 bg-gray-800 px-2 py-0.5 rounded">待验证</span>
                    )}
                  </div>

                  {/* 操作按钮 */}
                  <div className="w-24 flex justify-end gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
                    {!vuln.confirmed && (
                      <button
                        onClick={(e) => { e.stopPropagation(); handleConfirm(vuln.id); }}
                        className="text-xs bg-green-600 hover:bg-green-500 text-white px-2 py-1 rounded transition-colors"
                      >
                        确认
                      </button>
                    )}
                    <span className="text-xs text-gray-500">查看 →</span>
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {/* 漏洞详情弹窗 */}
        {selectedVuln && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm" onClick={() => setSelectedVuln(null)}>
            <div className="bg-gray-900 border border-gray-700 rounded-xl max-w-4xl w-full max-h-[90vh] overflow-hidden shadow-2xl" onClick={(e) => e.stopPropagation()}>
              {/* 弹窗头部 */}
              <div className="flex items-center justify-between p-5 border-b border-gray-800">
                <div className="flex items-center gap-3">
                  <span className="text-2xl">{SEVERITY_CONFIG[selectedVuln.severity]?.icon || '⚠️'}</span>
                  <div>
                    <h3 className="text-lg font-bold text-white">{selectedVuln.title || selectedVuln.type || '漏洞详情'}</h3>
                    <div className="flex gap-2 mt-1">
                      <span className={`text-xs px-2 py-0.5 rounded ${SEVERITY_CONFIG[selectedVuln.severity]?.bg} ${SEVERITY_CONFIG[selectedVuln.severity]?.color}`}>
                        {SEVERITY_CONFIG[selectedVuln.severity]?.label || selectedVuln.severity}
                      </span>
                      <span className="text-xs text-gray-400 bg-gray-800 px-2 py-0.5 rounded">
                        {TYPE_LABELS[selectedVuln.type] || selectedVuln.type}
                      </span>
                      {selectedVuln.risk_level > 0 && (
                        <span className="text-xs text-red-400 bg-red-500/10 px-2 py-0.5 rounded">
                          风险分: {selectedVuln.risk_level}/10
                        </span>
                      )}
                    </div>
                  </div>
                </div>
                <button onClick={() => setSelectedVuln(null)} className="text-gray-400 hover:text-white text-2xl leading-none">&times;</button>
              </div>

              {/* 弹窗内容 - 可滚动 */}
              <div className="overflow-y-auto max-h-[calc(90vh-140px)] p-5 space-y-4">
                {/* 基本信息 */}
                <section className="space-y-3">
                  <h4 className="text-sm font-semibold text-cyan-400 flex items-center gap-2">
                    🎯 基本信息
                  </h4>
                  <div className="grid grid-cols-2 gap-3">
                    <InfoRow label="目标地址" value={selectedVuln.url || selectedVuln.target || '-'} mono />
                    <InfoRow label="CVE编号" value={selectedVuln.cve_id || '-'} mono />
                    <InfoRow label="CNVD编号" value={selectedVuln.cnvd_id || '-'} mono />
                    <InfoRow label="发现时间" value={selectedVuln.created_at || '-'} />
                    <InfoRow label="更新时间" value={selectedVuln.updated_at || '-'} />
                    <InfoRow label="受影响版本" value={selectedVuln.affected_version || '-'} />
                    <InfoRow label="状态" value={selectedVuln.confirmed ? '✅ 已确认' : '⏳ 待验证'} />
                    {selectedVuln.tags && <InfoRow label="标签" value={selectedVuln.tags} />}
                  </div>
                </section>

                {/* 漏洞描述 */}
                {selectedVuln.description && (
                  <section className="space-y-2">
                    <h4 className="text-sm font-semibold text-cyan-400 flex items-center gap-2">📝 漏洞描述</h4>
                    <div className="bg-gray-800/50 rounded-lg p-4 text-sm text-gray-300 leading-relaxed border border-gray-700/50">
                      {selectedVuln.description}
                    </div>
                  </section>
                )}

                {/* 复现步骤 */}
                {selectedVuln.reproduce_steps && (
                  <section className="space-y-2">
                    <h4 className="text-sm font-semibold text-red-400 flex items-center gap-2">🔬 复现步骤</h4>
                    <div className="bg-red-950/20 rounded-lg p-4 text-sm text-gray-300 leading-relaxed border border-red-500/20 whitespace-pre-line">
                      {selectedVuln.reproduce_steps}
                    </div>
                  </section>
                )}

                {/* Payload */}
                {selectedVuln.payload && (
                  <section className="space-y-2">
                    <h4 className="text-sm font-semibold text-orange-400 flex items-center gap-2">💉 测试Payload</h4>
                    <div className="bg-gray-950 rounded-lg p-4 font-mono text-xs text-orange-300 border border-gray-700 overflow-x-auto whitespace-pre-wrap break-all">
                      {selectedVuln.payload}
                    </div>
                    <button
                      onClick={() => navigator.clipboard.writeText(selectedVuln.payload!)}
                      className="text-xs text-gray-500 hover:text-cyan-400 transition-colors"
                    >
                      📋 复制Payload
                    </button>
                  </section>
                )}

                {/* 请求包 */}
                {selectedVuln.request && (
                  <section className="space-y-2">
                    <h4 className="text-sm font-semibold text-blue-400 flex items-center gap-2">📤 HTTP请求包</h4>
                    <div className="bg-gray-950 rounded-lg p-4 font-mono text-xs text-blue-300 border border-gray-700 max-h-60 overflow-y-auto whitespace-pre-wrap break-all">
                      {selectedVuln.request}
                    </div>
                  </section>
                )}

                {/* 响应包 */}
                {selectedVuln.response && (
                  <section className="space-y-2">
                    <h4 className="text-sm font-semibold text-green-400 flex items-center gap-2">📥 HTTP响应包</h4>
                    <div className="bg-gray-950 rounded-lg p-4 font-mono text-xs text-green-300 border border-gray-700 max-h-60 overflow-y-auto whitespace-pre-wrap break-all">
                      {selectedVuln.response.substring(0, 5000)}
                      {selectedVuln.response.length > 5000 && <span className="text-gray-500 ml-2">...(响应过长，已截断)</span>}
                    </div>
                  </section>
                )}

                {/* 验证证据 */}
                {selectedVuln.proof && (
                  <section className="space-y-2">
                    <h4 className="text-sm font-semibold text-purple-400 flex items-center gap-2">🔍 验证证据</h4>
                    <div className="bg-gray-800/50 rounded-lg p-4 text-sm text-gray-300 border border-gray-700/50 whitespace-pre-wrap">
                      {selectedVuln.proof}
                    </div>
                  </section>
                )}

                {/* 修复建议 */}
                {selectedVuln.fix_suggestion && (
                  <section className="space-y-2">
                    <h4 className="text-sm font-semibold text-green-400 flex items-center gap-2">🛡️ 修复建议</h4>
                    <div className="bg-green-950/20 rounded-lg p-4 text-sm text-gray-300 leading-relaxed border border-green-500/20 whitespace-pre-line">
                      {selectedVuln.fix_suggestion}
                    </div>
                  </section>
                )}

                {/* 参考资料 */}
                {selectedVuln.references && (
                  <section className="space-y-2">
                    <h4 className="text-sm font-semibold text-yellow-400 flex items-center gap-2">📚 参考资料</h4>
                    <div className="bg-gray-800/50 rounded-lg p-4 text-sm text-gray-300 border border-gray-700/50 whitespace-pre-line">
                      {selectedVuln.references}
                    </div>
                  </section>
                )}

                {/* 原始数据 */}
                {selectedVuln.data && (
                  <section className="space-y-2">
                    <h4 className="text-sm font-semibold text-gray-400 flex items-center gap-2">📦 原始扫描数据</h4>
                    <details className="bg-gray-950 rounded-lg border border-gray-800">
                      <summary className="px-4 py-2 text-xs text-gray-500 cursor-pointer hover:text-gray-300">展开查看原始数据</summary>
                      <pre className="p-4 text-xs text-gray-400 overflow-x-auto whitespace-pre-wrap break-all max-h-60 overflow-y-auto">
                        {JSON.stringify(JSON.parse(selectedVuln.data), null, 2)}
                      </pre>
                    </details>
                  </section>
                )}
              </div>

              {/* 弹窗底部操作栏 */}
              <div className="flex justify-end gap-3 p-4 border-t border-gray-800 bg-gray-900/50">
                {!selectedVuln.confirmed && (
                  <button
                    onClick={() => handleConfirm(selectedVuln.id)}
                    className="bg-green-600 hover:bg-green-500 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors"
                  >
                    ✅ 确认此漏洞
                  </button>
                )}
                <button
                  onClick={() => setSelectedVuln(null)}
                  className="bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded-lg text-sm transition-colors"
                >
                  关闭
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

const InfoRow: React.FC<{ label: string; value: string; mono?: boolean }> = ({ label, value, mono }) => (
  <div className="bg-gray-800/50 rounded px-3 py-2 border border-gray-700/30">
    <div className="text-xs text-gray-500 mb-0.5">{label}</div>
    <div className={`text-sm text-gray-200 break-all ${mono ? 'font-mono' : ''}`}>{value}</div>
  </div>
);

export default VulnManagement;
