﻿﻿﻿﻿﻿﻿﻿﻿﻿﻿﻿//﻿// API 服务层 - 连接后端
const API_BASE_URL = 'http://localhost:8081/api/v1';

// 通用请求函数
async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;

  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Unknown error' }));
      throw new Error(error.error || `HTTP ${response.status}`);
    }

    return response.json();
  } catch (error) {
    console.error('API 请求失败:', error);
    throw error;
  }
}

// 目标管理
export const targetsApi = {
  list: () => request('/targets'),
  create: (data: { name: string; type: string; value: string; description?: string; tags?: string }) =>
    request('/targets', { method: 'POST', body: JSON.stringify(data) }),
  delete: (id: string) =>
    request(`/targets/${id}`, { method: 'DELETE' }),
  getAssetTree: (id: string) => request(`/targets/${id}/assets/tree`),
};

// 漏洞管理
export const vulnerabilitiesApi = {
  list: () => request('/vulnerabilities'),
  get: (id: string) => request(`/vulnerabilities/${id}`),
  confirm: (id: string) =>
    request(`/vulnerabilities/${id}/confirm`, { method: 'POST' }),
  patch: (id: string, data: any) =>
    request(`/vulnerabilities/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  bulkAction: (data: { ids: string[]; action: string; note?: string }) =>
    request('/vulnerabilities/bulk', { method: 'POST', body: JSON.stringify(data) }),
  addComment: (id: string, content: string) =>
    request(`/vulnerabilities/${id}/comments`, { method: 'POST', body: JSON.stringify({ content }) }),
  getHistory: (id: string) => request(`/vulnerabilities/${id}/history`),
};

// 漏洞看板统计
export const dashboardApi = {
  getVulnStats: () => request('/dashboard/stats'),
  getTrend: () => request('/dashboard/trend'),
  getTopAssets: () => request('/dashboard/top-assets'),
};

// 信息收集
export const reconApi = {
  list: () => request('/scans'),
  start: (data: { target_id: string; modules?: string[]; options?: any }) =>
    request('/scans', { method: 'POST', body: JSON.stringify(data) }),
  startRecon: (data: { target: string; type: string; options?: any }) =>
    request('/recon', { method: 'POST', body: JSON.stringify(data) }),
  get: (scanId: string) => request(`/scans/${scanId}`),
  getProgress: (scanId: string) => request(`/scans/${scanId}/progress`),
  stop: (scanId: string) => request(`/scans/${scanId}`, { method: 'DELETE' }),
};

// 漏洞扫描
export const scanApi = {
  list: () => request('/vuln-scans'),
  start: (data: { target_id: string; assets?: string[]; scanners?: string[]; poc_tags?: string[] }) =>
    request('/vuln-scans', { method: 'POST', body: JSON.stringify(data) }),
  get: (scanId: string) => request(`/vuln-scans/${scanId}`),
  pause: (scanId: string) => request(`/vuln-scans/${scanId}/pause`, { method: 'POST' }),
  resume: (scanId: string) => request(`/vuln-scans/${scanId}/resume`, { method: 'POST' }),
  stop: (scanId: string) => request(`/vuln-scans/${scanId}`, { method: 'DELETE' }),
};

// AI 助手
export const aiApi = {
  status: () => request('/ai/status'),
  chat: (message: string, context?: any) =>
    request('/ai/chat', { method: 'POST', body: JSON.stringify({ message, context }) }),
  getConfig: () => request('/ai/config'),
  updateConfig: (config: any) =>
    request('/ai/config', { method: 'PUT', body: JSON.stringify(config) }),
  testConnection: () => request('/ai/test', { method: 'POST' }),
  analyze: (data: any) =>
    request('/ai/analyze', { method: 'POST', body: JSON.stringify(data) }),
  generatePOC: (data: any) =>
    request('/ai/generate-poc', { method: 'POST', body: JSON.stringify(data) }),
  bypassWAF: (data: any) =>
    request('/ai/bypass-waf', { method: 'POST', body: JSON.stringify(data) }),
};

// SRC 平台
export const srcApi = {
  listPlatforms: () => request('/src/platforms'),
};

// 报告生成
export const reportApi = {
  getTemplates: () => request('/reports/templates'),
  generate: (data: any) =>
    request('/reports/generate', { method: 'POST', body: JSON.stringify(data) }),
  uploadTemplate: (data: any) =>
    request('/reports/templates/upload', { method: 'POST', body: JSON.stringify(data) }),
  download: (id: string) => request(`/reports/download/${id}`),
};

// 资产管理
export const assetsApi = {
  list: () => request('/assets'),
  get: (id: string) => request(`/assets/${id}`),
  create: (data: any) =>
    request('/assets', { method: 'POST', body: JSON.stringify(data) }),
  update: (id: string, data: any) =>
    request(`/assets/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  patch: (id: string, data: { tags?: string[]; notes?: string }) =>
    request(`/assets/${id}`, { method: 'PATCH', body: JSON.stringify(data) }),
  delete: (id: string) =>
    request(`/assets/${id}`, { method: 'DELETE' }),
  getStats: (targetId?: string) => {
    const query = targetId ? `?target_id=${targetId}` : '';
    return request(`/assets/stats${query}`);
  },
  search: (targetId: string, keyword: string) =>
    request(`/assets/search?target_id=${targetId}&keyword=${encodeURIComponent(keyword)}`),
};

// 信息收集工具管理
export const reconToolsApi = {
  listTools: (category?: string) => {
    const query = category ? `?category=${encodeURIComponent(category)}` : '';
    return request(`/recon/tools${query}`);
  },
  getTool: (id: string) => request(`/recon/tools/${id}`),
  checkToolInstall: (id: string) =>
    request(`/recon/tools/${id}/check-install`, { method: 'POST' }),

  listWorkflows: () => request('/recon/workflows'),
  listWorkflowTemplates: () => request('/recon/workflows/templates'),
  createWorkflow: (data: any) =>
    request('/recon/workflows', { method: 'POST', body: JSON.stringify(data) }),
  getWorkflow: (id: string) => request(`/recon/workflows/${id}`),
  updateWorkflow: (id: string, data: any) =>
    request(`/recon/workflows/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteWorkflow: (id: string) =>
    request(`/recon/workflows/${id}`, { method: 'DELETE' }),
  executeWorkflow: (id: string, data: any) =>
    request(`/recon/workflows/${id}/execute`, { method: 'POST', body: JSON.stringify(data) }),

  getWorkflowExecution: (id: string) => request(`/recon/executions/workflow/${id}`),
  listToolExecutions: (workflowExecId: string) =>
    request(`/recon/executions/workflow/${workflowExecId}/tools`),

  listDataSources: (enabledOnly?: boolean) => {
    const query = enabledOnly ? '?enabled=true' : '';
    return request(`/recon/datasources${query}`);
  },
  createDataSource: (data: any) =>
    request('/recon/datasources', { method: 'POST', body: JSON.stringify(data) }),
  updateDataSource: (id: string, data: any) =>
    request(`/recon/datasources/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
  deleteDataSource: (id: string) =>
    request(`/recon/datasources/${id}`, { method: 'DELETE' }),
};

// 教育 SRC
export const eduApi = {
  getUniversities: (params?: { level?: string; province?: string }) => {
    const query = new URLSearchParams(params as Record<string, string>).toString();
    return request(`/edu/domains${query ? `?${query}` : ''}`);
  },
  getStats: () => request('/edu/stats'),
  scan: (data: { domain_ids?: string[]; target?: string }) =>
    request('/edu/scan', { method: 'POST', body: JSON.stringify(data) }),
  getScanResults: (scanId: string) => request(`/edu/scan/${scanId}/results`),
  batchScan: (data: { level?: string; province?: string; concurrency?: number }) =>
    request('/edu/batch', { method: 'POST', body: JSON.stringify(data) }),
  generateReport: (data: { scan_id: string; src: string }) =>
    request('/edu/report', { method: 'POST', body: JSON.stringify(data) }),
  getSystems: () => request('/edu/systems'),
  importDomains: (data: any) =>
    request('/edu/domains/import', { method: 'POST', body: JSON.stringify(data) }),
  deleteDomain: (id: string) =>
    request(`/edu/domains/${id}`, { method: 'DELETE' }),
};

// 漏洞扫描模块
export const vulnScanApi = {
  createScanTask: (data: any) => request('/vuln-scans', { method: 'POST', body: JSON.stringify(data) }),
  listScanTasks: (targetId: string) => request(`/vuln-scans?target_id=${targetId}`),
  getScanTask: (id: string) => request(`/vuln-scans/${id}`),
  getScanResults: (taskId: string) => request(`/vuln-scan/results?task_id=${taskId}`),
  getScanResult: (id: string) => request(`/vuln-scan/results/${id}`),
  updateResultStatus: (id: string, data: any) =>
    request(`/vuln-scan/results/${id}/status`, { method: 'PUT', body: JSON.stringify(data) }),
  aiVerifyResult: (id: string) =>
    request(`/vuln-scan/results/${id}/ai-verify`, { method: 'POST' }),
  listScanConfigs: () => request('/vuln-scan/configs'),
  getDefaultConfig: () => request('/vuln-scan/configs/default'),
  listPOCs: () => request('/vuln-scan/pocs'),
  analyzeVulnChains: (targetId: string) => request(`/vuln-scan/chains/target/${targetId}`),
  generateWAFBypassPayloads: (data: { category: string; payload: string }) =>
    request('/vuln-scan/waf-bypass', { method: 'POST', body: JSON.stringify(data) }),
};

// 主动扫描
export const activeScanApi = {
  listTasks: () => request('/active-scans'),
  getTask: (id: string) => request(`/active-scans/${id}`),
  createTask: (data: { name?: string; target: string; type: string; options?: any }) =>
    request('/active-scans', { method: 'POST', body: JSON.stringify(data) }),
  stopTask: (id: string) => request(`/active-scans/${id}`, { method: 'DELETE' }),
  getResults: (taskId: string) => request(`/active-scans/${taskId}/results`),
};

// 网络设置
export const networkApi = {
  getProxyConfig: () => request('/proxy'),
  updateProxyConfig: (data: { enabled: boolean; type: string; host: string; port: number; username?: string; password?: string }) =>
    request('/proxy', { method: 'PUT', body: JSON.stringify(data) }),
  getNetworkConfig: () => request('/network/config'),
  updateNetworkConfig: (data: any) =>
    request('/network/config', { method: 'PUT', body: JSON.stringify(data) }),
  testConnection: (data: { url: string; proxy?: string }) =>
    request('/network/test', { method: 'POST', body: JSON.stringify(data) }),
};

// 漏洞技巧模块
export const vulnTechniquesApi = {
  list: (params?: { category?: string; severity?: string; difficulty?: string; keyword?: string }) => {
    const query = new URLSearchParams(params as Record<string, string>).toString();
    return request(`/vuln-techniques${query ? `?${query}` : ''}`);
  },
  get: (id: string) => request(`/vuln-techniques/${id}`),
  getStats: () => request('/vuln-techniques/stats'),
  practice: (data: { technique_id: string; target?: string; url?: string; payload?: string; mode?: string }) =>
    request('/vuln-techniques/practice', { method: 'POST', body: JSON.stringify(data) }),
  getPracticeResults: (techniqueId?: string) => {
    const query = techniqueId ? `?technique_id=${techniqueId}` : '';
    return request(`/vuln-techniques/practice/results${query}`);
  },
  getPracticeResult: (id: string) => request(`/vuln-techniques/practice/${id}`),
};

// 请求历史
export const requestHistoryApi = {
  list: (params?: { module?: string; target?: string; method?: string }) => {
    const query = new URLSearchParams(params as Record<string, string>).toString();
    return request(`/request-history${query ? `?${query}` : ''}`);
  },
  get: (id: string) => request(`/request-history/${id}`),
  delete: (id: string) => request(`/request-history/${id}`, { method: 'DELETE' }),
  clear: () => request('/request-history', { method: 'DELETE' }),
  export: () => {
    window.open(`${API_BASE_URL}/request-history/export`, '_blank');
  },
};

export default {
  targets: targetsApi,
  vulnerabilities: vulnerabilitiesApi,
  dashboard: dashboardApi,
  assets: assetsApi,
  recon: reconApi,
  reconTools: reconToolsApi,
  scan: scanApi,
  ai: aiApi,
  src: srcApi,
  report: reportApi,
  edu: eduApi,
  vulnScan: vulnScanApi,
  activeScan: activeScanApi,
  network: networkApi,
  vulnTechniques: vulnTechniquesApi,
  requestHistory: requestHistoryApi,
};
