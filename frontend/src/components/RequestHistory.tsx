import React, { useState, useEffect } from 'react';
import { requestHistoryApi } from '../services/api';

interface RequestLog {
  id: string;
  module: string;
  target: string;
  method: string;
  url: string;
  request_headers: string;
  request_body: string;
  response_status: number;
  response_headers: string;
  response_body: string;
  duration: number;
  error: string;
  timestamp: string;
}

const RequestHistory: React.FC = () => {
  const [logs, setLogs] = useState<RequestLog[]>([]);
  const [loading, setLoading] = useState(false);
  const [selectedLog, setSelectedLog] = useState<RequestLog | null>(null);
  const [filter, setFilter] = useState({
    module: '',
    target: '',
    method: '',
  });
  const [page, setPage] = useState(1);
  const pageSize = 50;
  const [autoRefresh, setAutoRefresh] = useState(false);

  const loadLogs = async () => {
    setLoading(true);
    try {
      const data = await requestHistoryApi.list({
        module: filter.module || undefined,
        target: filter.target || undefined,
        method: filter.method || undefined,
      } as any);
      setLogs(data || []);
    } catch (error) {
      console.error('加载日志失败:', error);
      setLogs([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadLogs();
  }, []);

  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(loadLogs, 3000);
      return () => clearInterval(interval);
    }
  }, [autoRefresh]);

  const handleFilter = () => {
    loadLogs();
    setPage(1);
  };

  const exportLogs = (format: 'json' | 'curl') => {
    if (format === 'json') {
      requestHistoryApi.export();
    } else if (format === 'curl' && selectedLog) {
      const curlCmd = generateCurlCommand(selectedLog);
      navigator.clipboard.writeText(curlCmd);
      alert('curl命令已复制到剪贴板！');
    }
  };

  const generateCurlCommand = (log: RequestLog): string => {
    let cmd = `curl -X ${log.method} '${log.url}'`;
    if (log.request_headers) {
      try {
        const headers = JSON.parse(log.request_headers);
        Object.entries(headers).forEach(([key, value]) => {
          cmd += ` -H '${key}: ${value}'`;
        });
      } catch (e) {}
    }
    if (log.request_body) {
      cmd += ` -d '${log.request_body}'`;
    }
    return cmd;
  };

  const clearLogs = async () => {
    if (!confirm('确定要清空所有请求日志吗？此操作不可恢复！')) return;

    try {
      await requestHistoryApi.clear();
      setLogs([]);
      setSelectedLog(null);
    } catch (error) {
      console.error('清空日志失败:', error);
      alert('清空失败！');
    }
  };

  const deleteLog = async (id: string) => {
    try {
      await requestHistoryApi.delete(id);
      setLogs(logs.filter(l => l.id !== id));
      if (selectedLog?.id === id) {
        setSelectedLog(null);
      }
    } catch (error) {
      console.error('删除日志失败:', error);
    }
  };

  const filteredLogs = logs.filter(log => {
    if (filter.module && log.module !== filter.module) return false;
    if (filter.target && !log.target.toLowerCase().includes(filter.target.toLowerCase())) return false;
    if (filter.method && log.method !== filter.method) return false;
    return true;
  });

  const paginatedLogs = filteredLogs.slice((page - 1) * pageSize, page * pageSize);
  const totalPages = Math.ceil(filteredLogs.length / pageSize);

  const getStatusColor = (status: number) => {
    if (status >= 200 && status < 300) return 'text-green-400';
    if (status >= 300 && status < 400) return 'text-blue-400';
    if (status >= 400 && status < 500) return 'text-yellow-400';
    if (status >= 500) return 'text-red-400';
    return 'text-gray-400';
  };

  const getMethodColor = (method: string) => {
    const colors: Record<string, string> = {
      GET: 'bg-green-600/20 text-green-400 border border-green-600/30',
      POST: 'bg-blue-600/20 text-blue-400 border border-blue-600/30',
      PUT: 'bg-yellow-600/20 text-yellow-400 border border-yellow-600/30',
      DELETE: 'bg-red-600/20 text-red-400 border border-red-600/30',
      PATCH: 'bg-purple-600/20 text-purple-400 border border-purple-600/30',
      OPTIONS: 'bg-gray-600/20 text-gray-400 border border-gray-600/30',
      HEAD: 'bg-gray-600/20 text-gray-400 border border-gray-600/30',
    };
    return colors[method] || 'bg-gray-600/20 text-gray-400 border border-gray-600/30';
  };

  const getModuleColor = (module: string) => {
    const colors: Record<string, string> = {
      recon: 'bg-cyan-600/20 text-cyan-400',
      scan: 'bg-orange-600/20 text-orange-400',
      'active-scan': 'bg-red-600/20 text-red-400',
      ai: 'bg-purple-600/20 text-purple-400',
      edu: 'bg-green-600/20 text-green-400',
      general: 'bg-gray-600/20 text-gray-400',
    };
    return colors[module] || 'bg-gray-600/20 text-gray-400';
  };

  const formatDuration = (ms: number) => {
    if (!ms) return '-';
    if (ms < 1000) return `${ms}ms`;
    return `${(ms / 1000).toFixed(2)}s`;
  };

  const formatBody = (body: string, maxLength: number = 500): string => {
    if (!body) return '(空)';
    try {
      const parsed = JSON.parse(body);
      const formatted = JSON.stringify(parsed, null, 2);
      return formatted.length > maxLength ? formatted.substring(0, maxLength) + '...\n(已截断)' : formatted;
    } catch {
      return body.length > maxLength ? body.substring(0, maxLength) + '...\n(已截断)' : body;
    }
  };

  return (
    <div className="h-full flex flex-col bg-gray-950">
      <div className="p-4 bg-gray-900 border-b border-gray-800">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-white">请求历史</h2>
            <p className="text-sm text-gray-400">所有网络请求记录 · 共 {logs.length} 条</p>
          </div>
          <div className="flex gap-2">
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
              onClick={loadLogs}
              className="px-3 py-1.5 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 text-sm"
            >
              刷新
            </button>
            <button
              onClick={() => exportLogs('json')}
              className="px-3 py-1.5 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 text-sm"
            >
              导出JSON
            </button>
            <button
              onClick={clearLogs}
              className="px-3 py-1.5 bg-red-600/20 text-red-400 rounded-lg hover:bg-red-600/30 text-sm"
            >
              清空
            </button>
          </div>
        </div>
      </div>

      <div className="p-4 bg-gray-900 border-b border-gray-800">
        <div className="flex gap-4">
          <div className="flex-1">
            <label className="block text-xs text-gray-400 mb-1">模块</label>
            <select
              value={filter.module}
              onChange={(e) => setFilter({ ...filter, module: e.target.value })}
              className="w-full px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm focus:border-blue-500 focus:outline-none"
            >
              <option value="">全部模块</option>
              <option value="recon">信息收集</option>
              <option value="scan">漏洞扫描</option>
              <option value="active-scan">主动扫描</option>
              <option value="ai">AI助手</option>
              <option value="edu">教育SRC</option>
              <option value="general">通用</option>
            </select>
          </div>
          <div className="flex-1">
            <label className="block text-xs text-gray-400 mb-1">目标搜索</label>
            <input
              type="text"
              value={filter.target}
              onChange={(e) => setFilter({ ...filter, target: e.target.value })}
              onKeyDown={(e) => e.key === 'Enter' && handleFilter()}
              placeholder="搜索目标..."
              className="w-full px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm focus:border-blue-500 focus:outline-none"
            />
          </div>
          <div className="flex-1">
            <label className="block text-xs text-gray-400 mb-1">请求方法</label>
            <select
              value={filter.method}
              onChange={(e) => setFilter({ ...filter, method: e.target.value })}
              className="w-full px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-white text-sm focus:border-blue-500 focus:outline-none"
            >
              <option value="">全部</option>
              <option value="GET">GET</option>
              <option value="POST">POST</option>
              <option value="PUT">PUT</option>
              <option value="DELETE">DELETE</option>
              <option value="PATCH">PATCH</option>
            </select>
          </div>
          <div className="flex items-end">
            <button
              onClick={handleFilter}
              className="px-4 py-1.5 bg-blue-600 text-white rounded-lg hover:bg-blue-500 text-sm"
            >
              筛选
            </button>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-auto">
        {loading ? (
          <div className="flex items-center justify-center h-full">
            <div className="text-gray-400">
              <div className="animate-spin text-4xl mb-4">⏳</div>
              <p>加载中...</p>
            </div>
          </div>
        ) : filteredLogs.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-full text-gray-500">
            <div className="text-4xl mb-4">📋</div>
            <p>暂无请求记录</p>
            <p className="text-sm mt-2">发起网络请求后将在此显示</p>
          </div>
        ) : (
          <table className="w-full">
            <thead className="bg-gray-900 sticky top-0 z-10">
              <tr className="text-left text-xs text-gray-400 border-b border-gray-800">
                <th className="px-4 py-2 font-medium w-32">时间</th>
                <th className="px-4 py-2 font-medium w-20">方法</th>
                <th className="px-4 py-2 font-medium">URL</th>
                <th className="px-4 py-2 font-medium w-20">状态</th>
                <th className="px-4 py-2 font-medium w-20">耗时</th>
                <th className="px-4 py-2 font-medium w-24">模块</th>
                <th className="px-4 py-2 font-medium w-20">操作</th>
              </tr>
            </thead>
            <tbody>
              {paginatedLogs.map((log) => (
                <tr
                  key={log.id}
                  onClick={() => setSelectedLog(log)}
                  className={`border-b border-gray-800 cursor-pointer hover:bg-gray-900/50 transition ${
                    selectedLog?.id === log.id ? 'bg-blue-600/10' : ''
                  }`}
                >
                  <td className="px-4 py-2 text-sm text-gray-400">
                    {log.timestamp ? new Date(log.timestamp).toLocaleTimeString() : '-'}
                  </td>
                  <td className="px-4 py-2">
                    <span className={`px-2 py-0.5 rounded text-xs font-medium ${getMethodColor(log.method)}`}>
                      {log.method}
                    </span>
                  </td>
                  <td className="px-4 py-2 text-sm text-white max-w-md truncate" title={log.url}>
                    {log.url}
                  </td>
                  <td className={`px-4 py-2 text-sm font-medium ${getStatusColor(log.response_status)}`}>
                    {log.response_status || (log.error ? 'Error' : '-')}
                  </td>
                  <td className="px-4 py-2 text-sm text-gray-400">
                    {formatDuration(log.duration)}
                  </td>
                  <td className="px-4 py-2">
                    <span className={`px-2 py-0.5 rounded text-xs ${getModuleColor(log.module)}`}>
                      {log.module || 'general'}
                    </span>
                  </td>
                  <td className="px-4 py-2">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        deleteLog(log.id);
                      }}
                      className="text-xs text-red-400 hover:text-red-300"
                    >
                      删除
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {totalPages > 1 && (
        <div className="p-4 bg-gray-900 border-t border-gray-800 flex items-center justify-between">
          <button
            onClick={() => setPage(Math.max(1, page - 1))}
            disabled={page === 1}
            className="px-3 py-1.5 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 disabled:opacity-50 text-sm"
          >
            上一页
          </button>
          <span className="text-sm text-gray-400">
            第 {page} / {totalPages} 页，共 {filteredLogs.length} 条
          </span>
          <button
            onClick={() => setPage(Math.min(totalPages, page + 1))}
            disabled={page === totalPages}
            className="px-3 py-1.5 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 disabled:opacity-50 text-sm"
          >
            下一页
          </button>
        </div>
      )}

      {selectedLog && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50" onClick={() => setSelectedLog(null)}>
          <div className="bg-gray-900 rounded-xl border border-gray-800 w-full max-w-5xl max-h-[85vh] overflow-hidden" onClick={(e) => e.stopPropagation()}>
            <div className="p-4 border-b border-gray-800 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <h3 className="text-lg font-semibold text-white">请求详情</h3>
                <span className={`px-2 py-0.5 rounded text-xs font-medium ${getMethodColor(selectedLog.method)}`}>
                  {selectedLog.method}
                </span>
                <span className={`text-sm font-medium ${getStatusColor(selectedLog.response_status)}`}>
                  {selectedLog.response_status || 'Error'}
                </span>
              </div>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => exportLogs('curl')}
                  className="px-3 py-1.5 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 text-sm"
                >
                  复制为curl
                </button>
                <button onClick={() => setSelectedLog(null)} className="text-gray-400 hover:text-white text-xl">
                  ✕
                </button>
              </div>
            </div>
            <div className="p-4 overflow-auto max-h-[70vh]">
              <div className="grid grid-cols-3 gap-4 mb-4">
                <div>
                  <p className="text-xs text-gray-400 mb-1">URL</p>
                  <p className="text-sm text-white break-all bg-gray-800 p-2 rounded">{selectedLog.url}</p>
                </div>
                <div>
                  <p className="text-xs text-gray-400 mb-1">目标</p>
                  <p className="text-sm text-white bg-gray-800 p-2 rounded">{selectedLog.target || '-'}</p>
                </div>
                <div>
                  <p className="text-xs text-gray-400 mb-1">模块 / 耗时</p>
                  <p className="text-sm text-white bg-gray-800 p-2 rounded">
                    <span className={`px-1.5 py-0.5 rounded text-xs ${getModuleColor(selectedLog.module)}`}>
                      {selectedLog.module || 'general'}
                    </span>
                    <span className="ml-2">{formatDuration(selectedLog.duration)}</span>
                  </p>
                </div>
              </div>

              <div className="mb-4">
                <p className="text-xs text-gray-400 mb-1">请求头</p>
                <pre className="text-xs text-gray-300 bg-gray-800 p-3 rounded-lg overflow-x-auto max-h-32">
                  {formatBody(selectedLog.request_headers, 1000)}
                </pre>
              </div>

              {selectedLog.request_body && (
                <div className="mb-4">
                  <p className="text-xs text-gray-400 mb-1">请求体</p>
                  <pre className="text-xs text-gray-300 bg-gray-800 p-3 rounded-lg overflow-x-auto max-h-48">
                    {formatBody(selectedLog.request_body, 2000)}
                  </pre>
                </div>
              )}

              <div className="mb-4">
                <p className="text-xs text-gray-400 mb-1">响应头</p>
                <pre className="text-xs text-gray-300 bg-gray-800 p-3 rounded-lg overflow-x-auto max-h-32">
                  {formatBody(selectedLog.response_headers, 1000)}
                </pre>
              </div>

              <div className="mb-4">
                <p className="text-xs text-gray-400 mb-1">响应体</p>
                <pre className="text-xs text-gray-300 bg-gray-800 p-3 rounded-lg overflow-x-auto max-h-64">
                  {formatBody(selectedLog.response_body, 5000)}
                </pre>
              </div>

              {selectedLog.error && (
                <div className="p-3 bg-red-600/20 rounded-lg border border-red-600/30">
                  <p className="text-xs text-red-400 font-medium mb-1">错误信息</p>
                  <p className="text-sm text-red-300">{selectedLog.error}</p>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default RequestHistory;
