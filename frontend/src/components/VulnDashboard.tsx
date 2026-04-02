import React, { useState, useEffect } from 'react';
import { vulnerabilitiesApi } from '../services/api';

interface VulnStats {
  critical: number;
  high: number;
  medium: number;
  low: number;
  info: number;
  total: number;
  confirmed: number;
}

const VulnDashboard: React.FC = () => {
  const [stats, setStats] = useState<VulnStats>({
    critical: 0,
    high: 0,
    medium: 0,
    low: 0,
    info: 0,
    total: 0,
    confirmed: 0,
  });
  const [loading, setLoading] = useState(true);
  const [recentVulns, setRecentVulns] = useState<any[]>([]);
  const [autoRefresh, setAutoRefresh] = useState(false);

  useEffect(() => {
    loadStats();
  }, []);

  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(loadStats, 5000);
      return () => clearInterval(interval);
    }
  }, [autoRefresh]);

  const loadStats = async () => {
    try {
      const data = await vulnerabilitiesApi.list();
      const vulns = Array.isArray(data) ? data : [];
      
      const newStats: VulnStats = {
        critical: vulns.filter((v: any) => v.severity === 'critical').length,
        high: vulns.filter((v: any) => v.severity === 'high').length,
        medium: vulns.filter((v: any) => v.severity === 'medium').length,
        low: vulns.filter((v: any) => v.severity === 'low').length,
        info: vulns.filter((v: any) => v.severity === 'info').length,
        total: vulns.length,
        confirmed: vulns.filter((v: any) => v.confirmed).length,
      };
      
      setStats(newStats);
      setRecentVulns(vulns.slice(0, 5));
    } catch (error) {
      console.error('加载漏洞统计失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const getSeverityColor = (severity: string) => {
    const colors: Record<string, string> = {
      critical: 'text-red-500 bg-red-500/10 border-red-500/20',
      high: 'text-orange-500 bg-orange-500/10 border-orange-500/20',
      medium: 'text-yellow-500 bg-yellow-500/10 border-yellow-500/20',
      low: 'text-blue-500 bg-blue-500/10 border-blue-500/20',
      info: 'text-gray-500 bg-gray-500/10 border-gray-500/20',
    };
    return colors[severity] || colors.info;
  };

  const getSeverityLabel = (severity: string) => {
    const labels: Record<string, string> = {
      critical: '严重',
      high: '高危',
      medium: '中危',
      low: '低危',
      info: '信息',
    };
    return labels[severity] || '未知';
  };

  if (loading) {
    return (
      <div className="h-full flex items-center justify-center bg-gray-950">
        <div className="text-gray-400">
          <div className="animate-spin text-4xl mb-4">⏳</div>
          <p>加载统计数据...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-6xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-white">漏洞看板</h2>
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
              onClick={loadStats}
              className="px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition"
            >
              刷新
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-5 gap-4 mb-6">
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <p className="text-sm text-gray-400">严重</p>
            <p className="text-3xl font-bold text-red-500 mt-1">{stats.critical}</p>
          </div>
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <p className="text-sm text-gray-400">高危</p>
            <p className="text-3xl font-bold text-orange-500 mt-1">{stats.high}</p>
          </div>
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <p className="text-sm text-gray-400">中危</p>
            <p className="text-3xl font-bold text-yellow-500 mt-1">{stats.medium}</p>
          </div>
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <p className="text-sm text-gray-400">低危</p>
            <p className="text-3xl font-bold text-blue-500 mt-1">{stats.low}</p>
          </div>
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <p className="text-sm text-gray-400">信息</p>
            <p className="text-3xl font-bold text-gray-500 mt-1">{stats.info}</p>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <div className="bg-gradient-to-br from-blue-600/20 to-purple-600/20 rounded-xl border border-blue-500/20 p-5">
            <p className="text-sm text-gray-400">漏洞总数</p>
            <p className="text-4xl font-bold text-white mt-1">{stats.total}</p>
          </div>
          <div className="bg-gradient-to-br from-green-600/20 to-cyan-600/20 rounded-xl border border-green-500/20 p-5">
            <p className="text-sm text-gray-400">已确认</p>
            <p className="text-4xl font-bold text-green-400 mt-1">{stats.confirmed}</p>
          </div>
          <div className="bg-gradient-to-br from-orange-600/20 to-red-600/20 rounded-xl border border-orange-500/20 p-5">
            <p className="text-sm text-gray-400">待处理</p>
            <p className="text-4xl font-bold text-orange-400 mt-1">{stats.total - stats.confirmed}</p>
          </div>
        </div>

        {stats.total > 0 && (
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6">
            <h3 className="text-lg font-semibold text-white mb-4">漏洞分布</h3>
            <div className="space-y-3">
              {[
                { key: 'critical', label: '严重', count: stats.critical, total: stats.total, color: 'bg-red-500' },
                { key: 'high', label: '高危', count: stats.high, total: stats.total, color: 'bg-orange-500' },
                { key: 'medium', label: '中危', count: stats.medium, total: stats.total, color: 'bg-yellow-500' },
                { key: 'low', label: '低危', count: stats.low, total: stats.total, color: 'bg-blue-500' },
                { key: 'info', label: '信息', count: stats.info, total: stats.total, color: 'bg-gray-500' },
              ].map((item) => (
                <div key={item.key} className="flex items-center gap-3">
                  <span className="text-sm text-gray-400 w-12">{item.label}</span>
                  <div className="flex-1 bg-gray-800 rounded-full h-3 overflow-hidden">
                    <div
                      className={`h-full ${item.color} transition-all duration-500`}
                      style={{ width: `${item.total > 0 ? (item.count / item.total) * 100 : 0}%` }}
                    />
                  </div>
                  <span className="text-sm text-gray-400 w-16 text-right">
                    {item.count} ({item.total > 0 ? ((item.count / item.total) * 100).toFixed(1) : 0}%)
                  </span>
                </div>
              ))}
            </div>
          </div>
        )}

        {recentVulns.length > 0 && (
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 mt-6">
            <h3 className="text-lg font-semibold text-white mb-4">最近发现</h3>
            <div className="space-y-3">
              {recentVulns.map((vuln: any) => (
                <div key={vuln.id} className="flex items-center justify-between p-3 bg-gray-800/50 rounded-lg">
                  <div className="flex items-center gap-3">
                    <span className={`px-2 py-0.5 rounded text-xs ${getSeverityColor(vuln.severity)}`}>
                      {getSeverityLabel(vuln.severity)}
                    </span>
                    <div>
                      <p className="text-white font-medium">{vuln.title}</p>
                      <p className="text-xs text-gray-400 font-mono truncate max-w-md">{vuln.url}</p>
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    {vuln.confirmed && (
                      <span className="text-xs text-green-400">✓ 已确认</span>
                    )}
                    <span className="text-xs text-gray-500">
                      {vuln.created_at ? new Date(vuln.created_at).toLocaleDateString() : ''}
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {stats.total === 0 && (
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-12 text-center">
            <div className="text-4xl mb-4">🔍</div>
            <p className="text-gray-400">暂无漏洞数据</p>
            <p className="text-sm text-gray-500 mt-2">开始扫描以发现漏洞</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default VulnDashboard;
