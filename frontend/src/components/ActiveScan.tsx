import React, { useState, useEffect, useCallback } from 'react';
import { activeScanApi, vulnerabilitiesApi } from '../services/api';

interface ActiveScanTask {
  id: string;
  name: string;
  target: string;
  type: 'web' | 'port' | 'vuln' | 'dir' | 'fuzz';
  status: 'pending' | 'running' | 'completed' | 'failed' | 'cancelled';
  progress: number;
  options: any;
  total_items: number;
  processed_items: number;
  error: string | null;
  started_at: string | null;
  finished_at: string | null;
  created_at: string;
  updated_at: string;
}

interface ScanResult {
  id: string;
  task_id: string;
  target: string;
  type: string;
  severity: 'critical' | 'high' | 'medium' | 'low' | 'info';
  data: string;
  confirmed: boolean;
  created_at: string;
}

const scanTypes = [
  { value: 'web', label: 'Web扫描', desc: '基础Web信息收集，检测敏感路径' },
  { value: 'port', label: '端口扫描', desc: '检测开放端口和服务识别' },
  { value: 'vuln', label: '漏洞扫描', desc: 'OWASP Top 10漏洞检测' },
  { value: 'dir', label: '目录扫描', desc: '敏感目录和文件探测' },
  { value: 'fuzz', label: '模糊测试', desc: '参数模糊测试和边界检测' },
];

const portPresets = [
  { label: '常见端口', ports: '21,22,23,25,53,80,110,143,443,445,993,995,3306,3389,5432,6379,8080,8443,27017' },
  { label: '高危端口', ports: '21,22,23,445,1433,3306,3389,5432,6379,9200,27017' },
  { label: 'Web端口', ports: '80,443,8080,8443,8888,9000,9090,3000,5000' },
  { label: '全端口', ports: '1-65535' },
];

const ActiveScan: React.FC = () => {
  const [tasks, setTasks] = useState<ActiveScanTask[]>([]);
  const [results, setResults] = useState<ScanResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [selectedTask, setSelectedTask] = useState<string | null>(null);
  const [selectedResult, setSelectedResult] = useState<ScanResult | null>(null);
  const [autoRefresh, setAutoRefresh] = useState(false);
  const [activeTab, setActiveTab] = useState<'tasks' | 'results'>('tasks');
  const [newTask, setNewTask] = useState({
    name: '',
    target: '',
    type: 'web' as const,
    ports: '21,22,23,25,53,80,110,143,443,445,993,995,3306,3389,5432,6379,8080,8443,27017',
    severity: 'all',
    template_tags: '',
    custom_headers: '',
    cookies: '',
    user_agent: '',
    timeout: 30,
    concurrency: 10,
    follow_redirects: true,
    skip_tls_verify: false,
  });

  const loadTasks = useCallback(async () => {
    try {
      const data: any = await activeScanApi.listTasks();
      setTasks(data?.data || data?.tasks || data || []);
    } catch (error) {
      console.error('加载任务失败:', error);
    }
  }, []);

  const loadTaskResults = async (taskId: string) => {
    try {
      const data: any = await activeScanApi.getResults(taskId);
      setResults(data?.data || data?.results || data || []);
    } catch (error) {
      console.error('加载结果失败:', error);
      setResults([]);
    }
  };

  useEffect(() => {
    loadTasks();
  }, [loadTasks]);

  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(() => {
        loadTasks();
        if (selectedTask) {
          const selectedTaskData = tasks.find(t => t.id === selectedTask);
          if (selectedTaskData && (selectedTaskData.status === 'running' || selectedTaskData.status === 'pending')) {
            loadTaskResults(selectedTask);
          }
        }
      }, 3000);
      return () => clearInterval(interval);
    }
  }, [autoRefresh, loadTasks, selectedTask, tasks]);

  useEffect(() => {
    if (selectedTask) {
      loadTaskResults(selectedTask);
    }
  }, [selectedTask]);

  const createTask = async () => {
    if (!newTask.target) {
      alert('请填写目标地址');
      return;
    }

    setLoading(true);
    try {
      const options: any = {
        timeout: newTask.timeout,
        concurrency: newTask.concurrency,
        follow_redirects: newTask.follow_redirects,
        skip_tls_verify: newTask.skip_tls_verify,
      };

      if (newTask.type === 'port') {
        options.ports = newTask.ports;
      }

      if (newTask.type === 'vuln' || newTask.type === 'fuzz') {
        options.severity = newTask.severity;
        if (newTask.template_tags) {
          options.tags = newTask.template_tags.split(',').map(t => t.trim());
        }
      }

      if (newTask.custom_headers) {
        try {
          options.headers = JSON.parse(newTask.custom_headers);
        } catch (e) {
          options.headers = newTask.custom_headers;
        }
      }

      if (newTask.cookies) {
        options.cookies = newTask.cookies;
      }

      if (newTask.user_agent) {
        options.user_agent = newTask.user_agent;
      }

      const result = await activeScanApi.createTask({
        name: newTask.name || `扫描任务 - ${newTask.target}`,
        target: newTask.target,
        type: newTask.type,
        options,
      });

      alert('任务创建成功！任务ID: ' + result.id);
      setShowCreateModal(false);
      setNewTask({
        name: '',
        target: '',
        type: 'web',
        ports: '21,22,23,25,53,80,110,143,443,445,993,995,3306,3389,5432,6379,8080,8443,27017',
        severity: 'all',
        template_tags: '',
        custom_headers: '',
        cookies: '',
        user_agent: '',
        timeout: 30,
        concurrency: 10,
        follow_redirects: true,
        skip_tls_verify: false,
      });
      loadTasks();
    } catch (error) {
      console.error('创建任务失败:', error);
      alert('创建任务失败！');
    } finally {
      setLoading(false);
    }
  };

  const stopTask = async (taskId: string) => {
    if (!confirm('确定要停止这个任务吗？')) return;

    try {
      await activeScanApi.stopTask(taskId);
      loadTasks();
    } catch (error) {
      console.error('停止任务失败:', error);
    }
  };

  const deleteTask = async (taskId: string) => {
    if (!confirm('确定要删除这个任务吗？')) return;

    try {
      await activeScanApi.stopTask(taskId);
      loadTasks();
      if (selectedTask === taskId) {
        setSelectedTask(null);
        setResults([]);
      }
    } catch (error) {
      console.error('删除任务失败:', error);
    }
  };

  const createVulnFromResult = async (result: ScanResult) => {
    const vulnInfo = {
      title: result.title || result.data || '扫描发现',
      severity: result.severity?.toUpperCase() || 'MEDIUM',
      target: result.target || result.url || '',
      type: result.type || 'web_scan',
      description: result.description || '',
      payload: result.payload || '',
      reproduce_steps: result.reproduce_steps || '',
      fix_suggestion: '请根据漏洞类型进行相应修复',
    };
    try {
      const response = await fetch(`http://localhost:8081/api/v1/vulnerabilities`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(vulnInfo),
      });
      if (response.ok) {
        alert('✅ 已成功添加到漏洞管理！');
      } else {
        const text = await response.text();
        console.warn('后端返回:', text);
        alert('⚠️ 漏洞信息已记录（后端可能不支持直接创建，请手动确认）');
      }
    } catch (error) {
      console.error('添加漏洞失败:', error);
      const blob = new Blob([JSON.stringify(vulnInfo, null, 2)], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `vuln_${result.id}.json`;
      a.click();
      URL.revokeObjectURL(url);
      alert('已导出漏洞信息为JSON文件（后端未响应）');
    }
  };

  const getStatusBadge = (status: string) => {
    const badges: Record<string, { class: string; icon: string }> = {
      pending: { class: 'bg-gray-600/20 text-gray-400 border border-gray-600/30', icon: '⏳' },
      running: { class: 'bg-blue-600/20 text-blue-400 border border-blue-600/30', icon: '🔄' },
      completed: { class: 'bg-green-600/20 text-green-400 border border-green-600/30', icon: '✅' },
      failed: { class: 'bg-red-600/20 text-red-400 border border-red-600/30', icon: '❌' },
      cancelled: { class: 'bg-yellow-600/20 text-yellow-400 border border-yellow-600/30', icon: '⏹️' },
    };
    return badges[status] || badges.pending;
  };

  const getTypeBadge = (type: string) => {
    const badges: Record<string, { class: string; label: string }> = {
      web: { class: 'bg-blue-600/20 text-blue-400 border border-blue-600/30', label: 'Web' },
      port: { class: 'bg-green-600/20 text-green-400 border border-green-600/30', label: '端口' },
      vuln: { class: 'bg-red-600/20 text-red-400 border border-red-600/30', label: '漏洞' },
      dir: { class: 'bg-purple-600/20 text-purple-400 border border-purple-600/30', label: '目录' },
      fuzz: { class: 'bg-orange-600/20 text-orange-400 border border-orange-600/30', label: '模糊' },
    };
    return badges[type] || badges.web;
  };

  const getSeverityBadge = (severity: string) => {
    const colors: Record<string, string> = {
      critical: 'bg-red-600/20 text-red-400 border border-red-600/30',
      high: 'bg-orange-600/20 text-orange-400 border border-orange-600/30',
      medium: 'bg-yellow-600/20 text-yellow-400 border border-yellow-600/30',
      low: 'bg-blue-600/20 text-blue-400 border border-blue-600/30',
      info: 'bg-gray-600/20 text-gray-400 border border-gray-600/30',
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

  const runningTasks = tasks.filter(t => t.status === 'running');
  const completedTasks = tasks.filter(t => t.status === 'completed');
  const failedTasks = tasks.filter(t => t.status === 'failed');

  return (
    <div className="h-full flex flex-col bg-gray-950">
      <div className="p-4 bg-gray-900 border-b border-gray-800">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-white">主动扫描</h2>
            <p className="text-sm text-gray-400">
              运行中: {runningTasks.length} | 已完成: {completedTasks.length} | 失败: {failedTasks.length}
            </p>
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
              onClick={loadTasks}
              className="px-4 py-2 bg-gray-700 text-white rounded-lg hover:bg-gray-600 transition"
            >
              刷新
            </button>
            <button
              onClick={() => setShowCreateModal(true)}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500 transition"
            >
              新建任务
            </button>
          </div>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden">
        <div className="w-1/2 border-r border-gray-800 overflow-auto flex flex-col">
          <div className="p-3 bg-gray-900 border-b border-gray-800 flex items-center justify-between">
            <h3 className="text-sm font-medium text-gray-400">扫描任务 ({tasks.length})</h3>
            <div className="flex gap-2">
              {runningTasks.length > 0 && (
                <span className="px-2 py-0.5 bg-blue-600/20 text-blue-400 rounded text-xs">
                  {runningTasks.length} 运行中
                </span>
              )}
            </div>
          </div>

          <div className="flex-1 overflow-auto p-4">
            {loading && tasks.length === 0 ? (
              <div className="text-center py-12 text-gray-500">
                <div className="animate-spin text-4xl mb-4">⏳</div>
                <p>加载中...</p>
              </div>
            ) : tasks.length === 0 ? (
              <div className="text-center py-12 text-gray-500">
                <div className="text-4xl mb-4">🔍</div>
                <p>暂无扫描任务</p>
                <p className="text-sm mt-2">点击"新建任务"开始扫描</p>
              </div>
            ) : (
              <div className="space-y-2">
                {tasks.map((task: any) => {
                  const statusBadge = getStatusBadge(task.status || 'pending');
                  const typeBadge = getTypeBadge(task.type || 'web');
                  return (
                    <div
                      key={task.id}
                      onClick={() => {
                        setSelectedTask(task.id);
                        loadTaskResults(task.id);
                      }}
                      className={`p-4 rounded-lg border cursor-pointer transition ${
                        selectedTask === task.id
                          ? 'bg-blue-600/20 border-blue-600/30'
                          : 'bg-gray-900 border-gray-800 hover:border-gray-700'
                      }`}
                    >
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center gap-2">
                          <span className="text-white font-medium truncate max-w-[200px]">
                            {task.name || task.target}
                          </span>
                          <span className={`px-2 py-0.5 rounded-full text-xs border ${typeBadge.class}`}>
                            {typeBadge.label}
                          </span>
                        </div>
                        <span className={`px-2 py-0.5 rounded-full text-xs border ${statusBadge.class}`}>
                          {statusBadge.icon} {task.status === 'pending' ? '等待中' :
                           task.status === 'running' ? '运行中' :
                           task.status === 'completed' ? '已完成' :
                           task.status === 'failed' ? '失败' :
                           task.status === 'cancelled' ? '已取消' : task.status}
                        </span>
                      </div>
                      <p className="text-sm text-gray-400 mb-2 truncate">{task.target}</p>
                      {task.status === 'running' && (
                        <div className="mb-2">
                          <div className="flex justify-between text-xs text-gray-500 mb-1">
                            <span>进度</span>
                            <span>{task.progress || 0}%</span>
                          </div>
                          <div className="w-full bg-gray-800 rounded-full h-1.5">
                            <div
                              className="bg-blue-600 h-1.5 rounded-full transition-all"
                              style={{ width: `${task.progress || 0}%` }}
                            />
                          </div>
                        </div>
                      )}
                      <div className="flex items-center justify-between">
                        <span className="text-xs text-gray-500">
                          {task.created_at ? new Date(task.created_at).toLocaleString() : ''}
                        </span>
                        <div className="flex gap-2">
                          {task.status === 'running' && (
                            <button
                              onClick={(e) => {
                                e.stopPropagation();
                                stopTask(task.id);
                              }}
                              className="text-xs px-2 py-1 bg-red-600/20 text-red-400 rounded hover:bg-red-600/30"
                            >
                              停止
                            </button>
                          )}
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              deleteTask(task.id);
                            }}
                            className="text-xs px-2 py-1 bg-gray-600/20 text-gray-400 rounded hover:bg-gray-600/30"
                          >
                            删除
                          </button>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>

        <div className="w-1/2 overflow-auto flex flex-col">
          <div className="p-3 bg-gray-900 border-b border-gray-800 flex items-center justify-between">
            <h3 className="text-sm font-medium text-gray-400">
              扫描结果 {selectedTask ? `(${results.length})` : ''}
            </h3>
            {results.length > 0 && (
              <div className="flex gap-2">
                <span className="text-xs px-2 py-0.5 bg-red-600/20 text-red-400 rounded">
                  严重: {results.filter((r: any) => r.severity === 'critical').length}
                </span>
                <span className="text-xs px-2 py-0.5 bg-orange-600/20 text-orange-400 rounded">
                  高危: {results.filter((r: any) => r.severity === 'high').length}
                </span>
                <span className="text-xs px-2 py-0.5 bg-yellow-600/20 text-yellow-400 rounded">
                  中危: {results.filter((r: any) => r.severity === 'medium').length}
                </span>
              </div>
            )}
          </div>

          <div className="flex-1 overflow-auto p-4">
            {selectedTask === null ? (
              <div className="text-center py-12 text-gray-500">
                <div className="text-4xl mb-4">📋</div>
                <p>选择任务查看结果</p>
              </div>
            ) : results.length === 0 ? (
              <div className="text-center py-12 text-gray-500">
                <div className="text-4xl mb-4">✅</div>
                <p>暂无扫描结果</p>
                <p className="text-sm mt-2">任务完成后将显示结果</p>
              </div>
            ) : (
              <div className="space-y-3">
                {results.map((result: any) => {
                  const resultType = result.type || 'info';
                  const typeLabels: Record<string, { label: string; color: string }> = {
                    'web_scan': { label: 'Web扫描', color: 'bg-blue-600/20 text-blue-400 border border-blue-600/30' },
                    'missing_header': { label: '安全头', color: 'bg-yellow-600/20 text-yellow-400 border border-yellow-600/30' },
                    'tech_detection': { label: '技术栈', color: 'bg-cyan-600/20 text-cyan-400 border border-cyan-600/30' },
                    'sensitive_info': { label: '敏感信息', color: 'bg-red-600/20 text-red-400 border border-red-600/30' },
                    'cors_misconfig': { label: 'CORS', color: 'bg-orange-600/20 text-orange-400 border border-orange-600/30' },
                    'cors_config': { label: 'CORS', color: 'bg-gray-600/20 text-gray-400 border border-gray-600/30' },
                    'cookie_security': { label: 'Cookie', color: 'bg-yellow-600/20 text-yellow-400 border border-yellow-600/30' },
                    'csrf_token': { label: 'CSRF', color: 'bg-orange-600/20 text-orange-400 border border-orange-600/30' },
                    'insecure_form': { label: '表单安全', color: 'bg-yellow-600/20 text-yellow-400 border border-yellow-600/30' },
                    'mixed_content': { label: '混合内容', color: 'bg-yellow-600/20 text-yellow-400 border border-yellow-600/30' },
                    'vuln_detected': { label: '漏洞', color: 'bg-red-600/20 text-red-400 border border-red-600/30' },
                    'sql_injection': { label: 'SQL注入', color: 'bg-red-600/20 text-red-400 border border-red-600/30' },
                    'xss': { label: 'XSS', color: 'bg-red-600/20 text-red-400 border border-red-600/30' },
                    'directory_traversal': { label: '目录遍历', color: 'bg-red-600/20 text-red-400 border border-red-600/30' },
                    'xxe': { label: 'XXE', color: 'bg-red-600/20 text-red-400 border border-red-600/30' },
                    'ssrf': { label: 'SSRF', color: 'bg-red-600/20 text-red-400 border border-red-600/30' },
                    'auth_required': { label: '认证', color: 'bg-purple-600/20 text-purple-400 border border-purple-600/30' },
                    'server_error': { label: '服务器错误', color: 'bg-orange-600/20 text-orange-400 border border-orange-600/30' },
                  };
                  const typeInfo = typeLabels[resultType] || { label: '其他', color: 'bg-gray-600/20 text-gray-400 border border-gray-600/30' };
                  
                  const getRecommendation = (type: string, severity: string): string => {
                    const recommendations: Record<string, string> = {
                      'missing_header': '建议配置安全响应头，如CSP、HSTS、X-Frame-Options等',
                      'cors_misconfig': '限制Access-Control-Allow-Origin为可信域名，避免使用通配符*',
                      'cookie_security': '为Cookie设置HttpOnly、Secure、SameSite属性',
                      'csrf_token': '为所有状态改变的操作添加CSRF Token验证',
                      'sql_injection': '使用参数化查询或ORM框架，对用户输入进行严格过滤',
                      'xss': '对用户输入进行HTML实体编码，配置CSP策略',
                      'directory_traversal': '验证用户输入的文件路径，使用白名单限制访问范围',
                      'xxe': '禁用XML外部实体处理，使用JSON替代XML',
                      'ssrf': '验证和限制用户提供的URL，禁止访问内网地址',
                      'sensitive_info': '移除页面中的敏感信息，检查备份文件和配置文件',
                      'tech_detection': '隐藏服务器技术栈信息，配置安全响应头',
                    };
                    return recommendations[type] || '';
                  };
                  
                  const recommendation = getRecommendation(resultType, result.severity);
                  
                  return (
                    <div 
                      key={result.id} 
                      className="p-4 bg-gray-900 rounded-lg border border-gray-800 hover:border-gray-700 transition cursor-pointer"
                      onClick={() => setSelectedResult(result)}
                    >
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center gap-2 flex-wrap">
                          <span className="text-lg">{getSeverityIcon(result.severity || 'info')}</span>
                          <span className={`px-2 py-0.5 rounded-full text-xs border ${getSeverityBadge(result.severity || 'info')}`}>
                            {(result.severity || 'info').toUpperCase()}
                          </span>
                          <span className={`px-2 py-0.5 rounded text-xs ${typeInfo.color}`}>
                            {typeInfo.label}
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
                              navigator.clipboard.writeText(result.target || '');
                            }}
                            className="text-xs px-2 py-1 bg-gray-600/20 text-gray-400 rounded hover:bg-gray-600/30"
                            title="复制URL"
                          >
                            复制
                          </button>
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              createVulnFromResult(result);
                            }}
                            className="text-xs px-2 py-1 bg-green-600/20 text-green-400 rounded hover:bg-green-600/30"
                          >
                            添加到漏洞
                          </button>
                          <span className="text-xs text-gray-500">
                            {result.created_at ? new Date(result.created_at).toLocaleString() : ''}
                          </span>
                        </div>
                      </div>
                      <p className="text-sm text-blue-400 mb-1 break-all font-mono" title={result.target}>
                        {result.target}
                      </p>
                      <p className="text-white font-medium mb-2">{result.title || result.data || '扫描结果'}</p>
                      {result.description && (
                        <p className="text-sm text-gray-400 mb-2 line-clamp-2">{result.description}</p>
                      )}
                      {recommendation && (
                        <div className="mt-2 p-2 bg-gray-800/50 rounded text-xs">
                          <span className="text-gray-400">💡 修复建议: </span>
                          <span className="text-gray-300">{recommendation}</span>
                        </div>
                      )}
                      {result.confirmed && (
                        <span className="text-xs px-2 py-0.5 bg-green-600/20 text-green-400 rounded mt-2 inline-block">
                          ✓ 已确认
                        </span>
                      )}
                      <div className="mt-2 text-xs text-gray-500 flex items-center gap-2">
                        <span>点击查看详情 →</span>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      </div>

      {showCreateModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 w-full max-w-2xl max-h-[90vh] overflow-auto">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold text-white">新建扫描任务</h3>
              <button
                onClick={() => setShowCreateModal(false)}
                className="text-gray-400 hover:text-white text-xl"
              >
                ✕
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">任务名称</label>
                <input
                  type="text"
                  value={newTask.name}
                  onChange={(e) => setNewTask({ ...newTask, name: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                  placeholder="例如：测试example.com"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">目标地址 *</label>
                <input
                  type="text"
                  value={newTask.target}
                  onChange={(e) => setNewTask({ ...newTask, target: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                  placeholder="例如：https://example.com 或 192.168.1.1"
                />
                <p className="text-xs text-gray-500 mt-1">支持域名、IP、URL格式</p>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-400 mb-2">扫描类型</label>
                <div className="grid grid-cols-5 gap-2">
                  {scanTypes.map((type) => (
                    <button
                      key={type.value}
                      onClick={() => setNewTask({ ...newTask, type: type.value as any })}
                      className={`p-3 rounded-lg text-left transition ${
                        newTask.type === type.value
                          ? 'bg-blue-600/20 border border-blue-600/30'
                          : 'bg-gray-800 border border-gray-700 hover:border-gray-600'
                      }`}
                    >
                      <div className="text-white text-sm font-medium">{type.label}</div>
                      <div className="text-xs text-gray-500 mt-1">{type.desc}</div>
                    </button>
                  ))}
                </div>
              </div>

              {newTask.type === 'port' && (
                <div>
                  <label className="block text-sm font-medium text-gray-400 mb-1">端口配置</label>
                  <div className="flex gap-2 mb-2">
                    {portPresets.map((preset) => (
                      <button
                        key={preset.label}
                        onClick={() => setNewTask({ ...newTask, ports: preset.ports })}
                        className="px-2 py-1 bg-gray-800 text-gray-300 rounded text-xs hover:bg-gray-700"
                      >
                        {preset.label}
                      </button>
                    ))}
                  </div>
                  <input
                    type="text"
                    value={newTask.ports}
                    onChange={(e) => setNewTask({ ...newTask, ports: e.target.value })}
                    className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                    placeholder="例如：1-1000, 80, 443"
                  />
                </div>
              )}

              {(newTask.type === 'vuln' || newTask.type === 'fuzz') && (
                <>
                  <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1">严重程度过滤</label>
                    <select
                      value={newTask.severity}
                      onChange={(e) => setNewTask({ ...newTask, severity: e.target.value })}
                      className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                    >
                      <option value="all">全部</option>
                      <option value="critical">严重</option>
                      <option value="high">高危</option>
                      <option value="medium">中危</option>
                      <option value="low">低危</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1">模板标签</label>
                    <input
                      type="text"
                      value={newTask.template_tags}
                      onChange={(e) => setNewTask({ ...newTask, template_tags: e.target.value })}
                      className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                      placeholder="例如：cve, sql-injection, xss (逗号分隔)"
                    />
                  </div>
                </>
              )}

              <div className="border-t border-gray-800 pt-4">
                <h4 className="text-sm font-medium text-gray-300 mb-3">高级选项</h4>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1">超时时间 (秒)</label>
                    <input
                      type="number"
                      value={newTask.timeout}
                      onChange={(e) => setNewTask({ ...newTask, timeout: parseInt(e.target.value) || 30 })}
                      className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-400 mb-1">并发数</label>
                    <input
                      type="number"
                      value={newTask.concurrency}
                      onChange={(e) => setNewTask({ ...newTask, concurrency: parseInt(e.target.value) || 10 })}
                      className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                    />
                  </div>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">自定义Headers (JSON格式)</label>
                <textarea
                  value={newTask.custom_headers}
                  onChange={(e) => setNewTask({ ...newTask, custom_headers: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none font-mono text-sm"
                  rows={2}
                  placeholder='{"Authorization": "Bearer xxx", "X-Custom": "value"}'
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">Cookies</label>
                <input
                  type="text"
                  value={newTask.cookies}
                  onChange={(e) => setNewTask({ ...newTask, cookies: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                  placeholder="session=xxx; token=yyy"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">自定义User-Agent</label>
                <input
                  type="text"
                  value={newTask.user_agent}
                  onChange={(e) => setNewTask({ ...newTask, user_agent: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                  placeholder="留空使用默认UA"
                />
              </div>

              <div className="flex gap-4">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={newTask.follow_redirects}
                    onChange={(e) => setNewTask({ ...newTask, follow_redirects: e.target.checked })}
                    className="w-4 h-4 rounded bg-gray-700 border-gray-600"
                  />
                  <span className="text-gray-300 text-sm">跟随重定向</span>
                </label>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={newTask.skip_tls_verify}
                    onChange={(e) => setNewTask({ ...newTask, skip_tls_verify: e.target.checked })}
                    className="w-4 h-4 rounded bg-gray-700 border-gray-600"
                  />
                  <span className="text-gray-300 text-sm">跳过TLS验证</span>
                </label>
              </div>
            </div>

            <div className="flex justify-end space-x-3 mt-6 pt-4 border-t border-gray-800">
              <button
                onClick={() => setShowCreateModal(false)}
                className="px-4 py-2 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 transition"
              >
                取消
              </button>
              <button
                onClick={createTask}
                disabled={loading || !newTask.target}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500 transition disabled:opacity-50"
              >
                {loading ? '创建中...' : '创建任务'}
              </button>
            </div>
          </div>
        </div>
      )}

      {selectedResult && (
        <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4" onClick={() => setSelectedResult(null)}>
          <div className="bg-gray-900 rounded-xl border border-gray-700 w-full max-w-4xl max-h-[90vh] overflow-hidden" onClick={(e) => e.stopPropagation()}>
            <div className="p-4 border-b border-gray-800 flex items-center justify-between bg-gray-800/50">
              <div className="flex items-center gap-3">
                <span className="text-2xl">{getSeverityIcon(selectedResult.severity || 'info')}</span>
                <div>
                  <h3 className="text-lg font-semibold text-white">{selectedResult.title || '漏洞详情'}</h3>
                  <div className="flex items-center gap-2 mt-1">
                    <span className={`px-2 py-0.5 rounded text-xs ${getSeverityBadge(selectedResult.severity || 'info')}`}>
                      {(selectedResult.severity || 'info').toUpperCase()}
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
                    {selectedResult.tags && (
                      <span className="px-2 py-0.5 rounded text-xs bg-gray-600/20 text-gray-400">
                        {selectedResult.tags}
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
                  <p className="text-blue-400 font-mono text-sm break-all bg-gray-900 p-3 rounded">
                    {selectedResult.url || selectedResult.target}
                  </p>
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
                          className="text-blue-400 hover:text-blue-300 text-sm block"
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
                  {selectedResult.affected_version && (
                    <div className="bg-gray-800/50 rounded-lg p-4">
                      <h4 className="text-sm font-medium text-gray-400 mb-1">影响版本</h4>
                      <p className="text-white text-sm">{selectedResult.affected_version}</p>
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
                  <button
                    onClick={() => createVulnFromResult(selectedResult)}
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
  );
};

export default ActiveScan;
