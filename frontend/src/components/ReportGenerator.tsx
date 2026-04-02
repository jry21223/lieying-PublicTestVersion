import React, { useState, useEffect } from 'react';
import { targetsApi, vulnerabilitiesApi, activeScanApi } from '../services/api';

interface Vulnerability {
  id: string;
  title: string;
  severity: string;
  status: string;
  target?: string;
  description?: string;
  created_at: string;
}

interface ScanResult {
  id: string;
  type: string;
  target: string;
  data: string;
  severity: string;
  created_at: string;
}

interface ReportTemplate {
  id: string;
  name: string;
  description: string;
  type: 'full' | 'executive' | 'technical' | 'edu-src';
}

const ReportGenerator: React.FC = () => {
  const [generating, setGenerating] = useState(false);
  const [reportType, setReportType] = useState<'full' | 'executive' | 'technical' | 'edu-src'>('full');
  const [includeStats, setIncludeStats] = useState(true);
  const [includeVulns, setIncludeVulns] = useState(true);
  const [includeAssets, setIncludeAssets] = useState(true);
  const [includeRecommendations, setIncludeRecommendations] = useState(true);
  const [vulnerabilities, setVulnerabilities] = useState<Vulnerability[]>([]);
  const [targets, setTargets] = useState<any[]>([]);
  const [scanResults, setScanResults] = useState<ScanResult[]>([]);
  const [generatedReport, setGeneratedReport] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [vulnData, targetData] = await Promise.all([
        vulnerabilitiesApi.list(),
        targetsApi.list(),
      ]);
      setVulnerabilities(Array.isArray(vulnData) ? vulnData : []);
      setTargets(Array.isArray(targetData) ? targetData : []);
    } catch (error) {
      console.error('加载数据失败:', error);
    }
  };

  const templates: ReportTemplate[] = [
    { id: 'full', name: '📋 完整报告', description: '包含所有内容的详细渗透测试报告', type: 'full' },
    { id: 'executive', name: '📊 执行摘要', description: '面向管理层的简洁报告，突出关键发现', type: 'executive' },
    { id: 'technical', name: '🔧 技术报告', description: '面向技术团队的详细技术分析', type: 'technical' },
    { id: 'edu-src', name: '🎓 教育SRC报告', description: '高校漏洞挖掘专项报告', type: 'edu-src' },
  ];

  const generateReport = async () => {
    setGenerating(true);
    setGeneratedReport(null);

    try {
      const stats = {
        totalTargets: targets.length,
        totalVulns: vulnerabilities.length,
        criticalVulns: vulnerabilities.filter(v => v.severity === 'critical').length,
        highVulns: vulnerabilities.filter(v => v.severity === 'high').length,
        mediumVulns: vulnerabilities.filter(v => v.severity === 'medium').length,
        lowVulns: vulnerabilities.filter(v => v.severity === 'low').length,
      };

      const html = generateHTMLReport(reportType, stats, vulnerabilities, targets, {
        includeStats,
        includeVulns,
        includeAssets,
        includeRecommendations,
      });

      setTimeout(() => {
        setGeneratedReport(html);
        setGenerating(false);
      }, 1500);
    } catch (error) {
      alert('报告生成失败');
      setGenerating(false);
    }
  };

  const generateHTMLReport = (
    type: string,
    stats: any,
    vulns: Vulnerability[],
    targets: any[],
    options: any
  ): string => {
    const now = new Date().toLocaleString('zh-CN');

    let content = '';
    if (type === 'executive') {
      content = generateExecutiveSummary(stats, vulns);
    } else if (type === 'technical') {
      content = generateTechnicalReport(stats, vulns);
    } else if (type === 'edu-src') {
      content = generateEduSRCReport(stats, vulns);
    } else {
      content = generateFullReport(stats, vulns, targets, options);
    }

    return `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>猎影渗透测试报告 - ${now}</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { font-family: 'Microsoft YaHei', Arial, sans-serif; background: #0f172a; color: #e2e8f0; line-height: 1.6; padding: 40px; }
    .container { max-width: 1000px; margin: 0 auto; }
    .header { text-align: center; margin-bottom: 40px; padding: 30px; background: linear-gradient(135deg, #1e293b, #334155); border-radius: 16px; }
    .header h1 { font-size: 32px; margin-bottom: 10px; color: #60a5fa; }
    .header .subtitle { color: #94a3b8; font-size: 14px; }
    .section { background: #1e293b; border-radius: 12px; padding: 24px; margin-bottom: 24px; border: 1px solid #334155; }
    .section h2 { font-size: 20px; color: #60a5fa; margin-bottom: 16px; padding-bottom: 8px; border-bottom: 2px solid #3b82f6; }
    .stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 16px; }
    .stat-card { background: #0f172a; padding: 20px; border-radius: 8px; text-align: center; }
    .stat-card .value { font-size: 36px; font-weight: bold; margin-bottom: 4px; }
    .stat-card .label { font-size: 12px; color: #94a3b8; text-transform: uppercase; }
    .stat-critical .value { color: #dc2626; }
    .stat-high .value { color: #f97316; }
    .stat-medium .value { color: #eab308; }
    .stat-low .value { color: #3b82f6; }
    .vuln-list { list-style: none; }
    .vuln-item { background: #0f172a; padding: 16px; border-radius: 8px; margin-bottom: 12px; border-left: 4px solid; }
    .vuln-item.critical { border-color: #dc2626; }
    .vuln-item.high { border-color: #f97316; }
    .vuln-item.medium { border-color: #eab308; }
    .vuln-item.low { border-color: #3b82f6; }
    .vuln-item h3 { font-size: 16px; margin-bottom: 8px; }
    .vuln-item .meta { font-size: 12px; color: #94a3b8; }
    .severity-badge { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; font-weight: bold; margin-left: 8px; }
    .severity-critical { background: #dc2626; color: white; }
    .severity-high { background: #f97316; color: white; }
    .severity-medium { background: #eab308; color: black; }
    .severity-low { background: #3b82f6; color: white; }
    .footer { text-align: center; margin-top: 40px; padding: 20px; color: #64748b; font-size: 12px; }
    .recommendation { background: #0f172a; padding: 16px; border-radius: 8px; margin-bottom: 12px; }
    .recommendation h4 { color: #22c55e; margin-bottom: 8px; }
    table { width: 100%; border-collapse: collapse; background: #0f172a; border-radius: 8px; overflow: hidden; }
    th, td { padding: 12px 16px; text-align: left; border-bottom: 1px solid #334155; }
    th { background: #334155; color: #e2e8f0; font-weight: 600; }
    tr:hover { background: #1e293b; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🔍 猎影渗透测试报告</h1>
      <p class="subtitle">Lieying Penetration Testing Report</p>
      <p class="subtitle">生成时间: ${now}</p>
    </div>

    ${content}

    <div class="footer">
      <p>本报告由猎影渗透测试平台自动生成</p>
      <p>昆仑安全实验室 | KunLun Security Lab</p>
    </div>
  </div>
</body>
</html>`;
  };

  const generateExecutiveSummary = (stats: any, vulns: Vulnerability[]): string => `
    <div class="section">
      <h2>📊 执行摘要</h2>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="value" style="color: #60a5fa">${stats.totalTargets}</div>
          <div class="label">测试目标</div>
        </div>
        <div class="stat-card stat-critical">
          <div class="value">${stats.criticalVulns}</div>
          <div class="label">严重漏洞</div>
        </div>
        <div class="stat-card stat-high">
          <div class="value">${stats.highVulns}</div>
          <div class="label">高危漏洞</div>
        </div>
        <div class="stat-card">
          <div class="value" style="color: #22c55e">${stats.totalVulns === 0 ? '无' : stats.totalVulns > 0 ? '有' : '待处理'}</div>
          <div class="label">风险状态</div>
        </div>
      </div>
    </div>

    <div class="section">
      <h2>⚠️ 关键发现</h2>
      ${vulns.filter(v => v.severity === 'critical' || v.severity === 'high').length === 0
        ? '<p>未发现严重或高危漏洞</p>'
        : `<ul class="vuln-list">
            ${vulns.filter(v => v.severity === 'critical' || v.severity === 'high').slice(0, 5).map(v => `
              <li class="vuln-item ${v.severity}">
                <h3>${v.title}<span class="severity-badge severity-${v.severity}">${v.severity.toUpperCase()}</span></h3>
                <p class="meta">${v.target || '未知目标'} | ${new Date(v.created_at).toLocaleDateString('zh-CN')}</p>
              </li>
            `).join('')}
          </ul>`
      }
    </div>

    <div class="section">
      <h2>💡 建议</h2>
      <div class="recommendation">
        <h4>🎯 立即处理</h4>
        <p>建议立即修复所有严重和高危漏洞，特别是可能导致数据泄露或服务器沦陷的漏洞</p>
      </div>
      <div class="recommendation">
        <h4>📅 定期扫描</h4>
        <p>建议每周进行一次自动化漏洞扫描，每月进行一次完整渗透测试</p>
      </div>
    </div>
  `;

  const generateTechnicalReport = (stats: any, vulns: Vulnerability[]): string => `
    <div class="section">
      <h2>🔍 漏洞详情</h2>
      <table>
        <thead>
          <tr>
            <th>漏洞名称</th>
            <th>严重性</th>
            <th>状态</th>
            <th>发现时间</th>
          </tr>
        </thead>
        <tbody>
          ${vulns.map(v => `
            <tr>
              <td>${v.title}</td>
              <td><span class="severity-badge severity-${v.severity}">${v.severity.toUpperCase()}</span></td>
              <td>${v.status}</td>
              <td>${new Date(v.created_at).toLocaleDateString('zh-CN')}</td>
            </tr>
          `).join('')}
          ${vulns.length === 0 ? '<tr><td colspan="4" style="text-align:center">未发现漏洞</td></tr>' : ''}
        </tbody>
      </table>
    </div>

    <div class="section">
      <h2>🛡️ 加固建议</h2>
      <div class="recommendation">
        <h4>1. 输入验证与过滤</h4>
        <p>对所有用户输入进行严格的验证和过滤，防止SQL注入、XSS等注入类攻击</p>
      </div>
      <div class="recommendation">
        <h4>2. 身份认证与授权</h4>
        <p>实施强密码策略，多因素认证，最小权限原则</p>
      </div>
      <div class="recommendation">
        <h4>3. 敏感数据保护</h4>
        <p>对敏感数据进行加密存储，传输过程使用HTTPS</p>
      </div>
      <div class="recommendation">
        <h4>4. 安全配置</h4>
        <p>及时更新补丁，关闭不必要的端口和服务，配置WAF</p>
      </div>
    </div>
  `;

  const generateEduSRCReport = (stats: any, vulns: Vulnerability[]): string => `
    <div class="section">
      <h2>🎓 教育SRC漏洞统计</h2>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="value" style="color: #60a5fa">${stats.totalTargets}</div>
          <div class="label">测试域名</div>
        </div>
        <div class="stat-card stat-critical">
          <div class="value">${stats.criticalVulns}</div>
          <div class="label">严重漏洞</div>
        </div>
        <div class="stat-card stat-high">
          <div class="value">${stats.highVulns}</div>
          <div class="label">高危漏洞</div>
        </div>
        <div class="stat-card stat-medium">
          <div class="value">${stats.mediumVulns}</div>
          <div class="label">中危漏洞</div>
        </div>
      </div>
    </div>

    <div class="section">
      <h2>🏛️ 漏洞列表</h2>
      ${vulns.length === 0
        ? '<p>暂未发现需要提交的漏洞，建议继续深入测试教务系统、VPN、邮件系统等</p>'
        : `<ul class="vuln-list">
            ${vulns.map(v => `
              <li class="vuln-item ${v.severity}">
                <h3>${v.title}<span class="severity-badge severity-${v.severity}">${v.severity.toUpperCase()}</span></h3>
                <p class="meta">${v.description || '无详细描述'} | ${v.target || '未知目标'}</p>
              </li>
            `).join('')}
          </ul>`
      }
    </div>

    <div class="section">
      <h2>💡 SRC挖掘建议</h2>
      <div class="recommendation">
        <h4>🎯 重点目标</h4>
        <ul style="margin-left: 20px; margin-top: 8px;">
          <li>教务管理系统、选课系统</li>
          <li>统一身份认证(CAS)、VPN系统</li>
          <li>邮件系统、OA办公系统</li>
          <li>图书馆系统、一卡通系统</li>
        </ul>
      </div>
      <div class="recommendation">
        <h4>📝 漏洞类型优先级</h4>
        <ul style="margin-left: 20px; margin-top: 8px;">
          <li>严重：Getshell、RCE、SQL注入（大量数据）</li>
          <li>高危：敏感信息泄露、任意文件读取、SSRF</li>
          <li>中危：存储型XSS、CSRF、弱口令</li>
        </ul>
      </div>
    </div>
  `;

  const generateFullReport = (stats: any, vulns: Vulnerability[], targets: any[], options: any): string => {
    let sections = '';

    if (options.includeStats) {
      sections += `
        <div class="section">
          <h2>📊 测试概览</h2>
          <div class="stats-grid">
            <div class="stat-card">
              <div class="value" style="color: #60a5fa">${stats.totalTargets}</div>
              <div class="label">测试目标</div>
            </div>
            <div class="stat-card stat-critical">
              <div class="value">${stats.criticalVulns}</div>
              <div class="label">严重漏洞</div>
            </div>
            <div class="stat-card stat-high">
              <div class="value">${stats.highVulns}</div>
              <div class="label">高危漏洞</div>
            </div>
            <div class="stat-card stat-medium">
              <div class="value">${stats.mediumVulns}</div>
              <div class="label">中危漏洞</div>
            </div>
          </div>
        </div>
      `;
    }

    if (options.includeAssets && targets.length > 0) {
      sections += `
        <div class="section">
          <h2>🎯 测试资产</h2>
          <table>
            <thead>
              <tr>
                <th>名称</th>
                <th>类型</th>
                <th>值</th>
              </tr>
            </thead>
            <tbody>
              ${targets.slice(0, 10).map(t => `
                <tr>
                  <td>${t.name}</td>
                  <td>${t.type}</td>
                  <td>${t.value}</td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      `;
    }

    if (options.includeVulns) {
      sections += `
        <div class="section">
          <h2>🔍 漏洞详情</h2>
          ${vulns.length === 0
            ? '<p>未发现漏洞</p>'
            : `<ul class="vuln-list">
                ${vulns.map(v => `
                  <li class="vuln-item ${v.severity}">
                    <h3>${v.title}<span class="severity-badge severity-${v.severity}">${v.severity.toUpperCase()}</span></h3>
                    <p class="meta">${v.description || '无详细描述'} | ${v.target || '未知目标'}</p>
                  </li>
                `).join('')}
              </ul>`
          }
        </div>
      `;
    }

    if (options.includeRecommendations) {
      sections += `
        <div class="section">
          <h2>💡 安全建议</h2>
          <div class="recommendation">
            <h4>🔒 紧急修复</h4>
            <p>立即修复所有严重和高危漏洞</p>
          </div>
          <div class="recommendation">
            <h4>📅 定期扫描</h4>
            <p>建议每周自动化扫描，每月完整渗透测试</p>
          </div>
          <div class="recommendation">
            <h4>🛡️ 纵深防御</h4>
            <p>实施WAF、入侵检测、最小权限原则</p>
          </div>
        </div>
      `;
    }

    return sections;
  };

  const downloadReport = () => {
    if (!generatedReport) return;

    const blob = new Blob([generatedReport], { type: 'text/html' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `猎影渗透测试报告_${new Date().toISOString().slice(0, 10)}.html`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  };

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-4xl mx-auto">
        <h2 className="text-2xl font-bold text-white mb-6">报告生成</h2>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <h3 className="text-sm font-medium text-gray-400 mb-3">报告模板</h3>
            <div className="space-y-2">
              {templates.map((t) => (
                <button
                  key={t.id}
                  onClick={() => setReportType(t.type)}
                  className={`w-full p-3 rounded-lg text-left transition ${
                    reportType === t.type
                      ? 'bg-blue-600/20 border border-blue-500'
                      : 'bg-gray-800 hover:bg-gray-750'
                  }`}
                >
                  <p className="text-white font-medium">{t.name}</p>
                  <p className="text-xs text-gray-400">{t.description}</p>
                </button>
              ))}
            </div>
          </div>

          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5">
            <h3 className="text-sm font-medium text-gray-400 mb-3">报告选项</h3>
            <div className="space-y-3">
              {[
                { checked: includeStats, setChecked: setIncludeStats, label: '📊 测试统计' },
                { checked: includeAssets, setChecked: setIncludeAssets, label: '🎯 资产清单' },
                { checked: includeVulns, setChecked: setIncludeVulns, label: '🔍 漏洞详情' },
                { checked: includeRecommendations, setChecked: setIncludeRecommendations, label: '💡 安全建议' },
              ].map((option, i) => (
                <label key={i} className="flex items-center gap-3 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={option.checked}
                    onChange={(e) => option.setChecked(e.target.checked)}
                    className="w-4 h-4 rounded bg-gray-800 border-gray-700 text-blue-600 focus:ring-blue-600"
                  />
                  <span className="text-white">{option.label}</span>
                </label>
              ))}
            </div>
          </div>
        </div>

        <div className="flex gap-3">
          <button
            onClick={generateReport}
            disabled={generating}
            className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-500 disabled:opacity-50 transition"
          >
            {generating ? '生成中...' : '🔍 生成报告'}
          </button>

          {generatedReport && (
            <button
              onClick={downloadReport}
              className="px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-500 transition"
            >
              📥 下载HTML报告
            </button>
          )}
        </div>

        {generatedReport && (
          <div className="mt-6 bg-gray-900 rounded-xl border border-gray-800 p-5">
            <h3 className="text-sm font-medium text-gray-400 mb-3">📋 报告预览</h3>
            <div
              className="bg-slate-900 rounded-lg p-4 max-h-96 overflow-auto"
              dangerouslySetInnerHTML={{ __html: generatedReport }}
            />
          </div>
        )}
      </div>
    </div>
  );
};

export default ReportGenerator;
