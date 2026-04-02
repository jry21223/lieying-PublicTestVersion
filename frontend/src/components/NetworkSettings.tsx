import React, { useState, useEffect } from 'react';
import { networkApi } from '../services/api';

interface ProxyConfig {
  id?: string;
  enabled: boolean;
  type: string;
  host: string;
  port: number;
  username: string;
  password: string;
}

interface NetworkConfig {
  id?: string;
  concurrency: number;
  rate_limit: number;
  timeout: number;
  retries: number;
  user_agent: string;
  follow_redirects: boolean;
  skip_tls_verify: boolean;
}

const NetworkSettings: React.FC = () => {
  const [proxyConfig, setProxyConfig] = useState<ProxyConfig>({
    enabled: false,
    type: 'http',
    host: '127.0.0.1',
    port: 8080,
    username: '',
    password: '',
  });

  const [networkConfig, setNetworkConfig] = useState<NetworkConfig>({
    concurrency: 10,
    rate_limit: 50,
    timeout: 30,
    retries: 3,
    user_agent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
    follow_redirects: true,
    skip_tls_verify: false,
  });

  const [testUrl, setTestUrl] = useState('https://example.com');
  const [testResult, setTestResult] = useState<{ success: boolean; message: string; status?: number; duration?: number } | null>(null);
  const [testing, setTesting] = useState(false);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    loadConfigs();
  }, []);

  const loadConfigs = async () => {
    setLoading(true);
    try {
      const [proxyData, networkData] = await Promise.all([
        networkApi.getProxyConfig(),
        networkApi.getNetworkConfig(),
      ]);
      if (proxyData) {
        setProxyConfig(proxyData as ProxyConfig);
      }
      if (networkData) {
        setNetworkConfig(networkData as NetworkConfig);
      }
    } catch (error) {
      console.error('加载配置失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const saveProxyConfig = async () => {
    setSaving(true);
    try {
      await networkApi.updateProxyConfig(proxyConfig);
      alert('代理配置已保存到后端！');
    } catch (error) {
      console.error('保存代理配置失败:', error);
      alert('保存失败！');
    } finally {
      setSaving(false);
    }
  };

  const saveNetworkConfig = async () => {
    setSaving(true);
    try {
      await networkApi.updateNetworkConfig(networkConfig);
      alert('网络配置已保存到后端！');
    } catch (error) {
      console.error('保存网络配置失败:', error);
      alert('保存失败！');
    } finally {
      setSaving(false);
    }
  };

  const testConnection = async () => {
    if (!testUrl) {
      setTestResult({ success: false, message: '请输入测试URL' });
      return;
    }

    setTesting(true);
    setTestResult(null);

    try {
      const proxyStr = proxyConfig.enabled
        ? `${proxyConfig.type}://${proxyConfig.host}:${proxyConfig.port}`
        : undefined;

      const startTime = Date.now();
      const result = await networkApi.testConnection({
        url: testUrl,
        proxy: proxyStr,
      });
      const duration = Date.now() - startTime;

      setTestResult({
        success: true,
        message: `连接成功`,
        status: result.status || 200,
        duration: duration,
      });
    } catch (error: any) {
      setTestResult({
        success: false,
        message: `连接失败: ${error.message}`,
      });
    } finally {
      setTesting(false);
    }
  };

  const userAgents = [
    { label: 'Chrome 120 (Windows)', value: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36' },
    { label: 'Chrome 120 (Mac)', value: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36' },
    { label: 'Firefox 121 (Windows)', value: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0' },
    { label: 'Safari 17 (Mac)', value: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15' },
    { label: 'Edge 120 (Windows)', value: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0' },
    { label: 'Googlebot', value: 'Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)' },
    { label: 'Bingbot', value: 'Mozilla/5.0 (compatible; bingbot/2.0; +http://www.bing.com/bingbot.htm)' },
    { label: 'curl/7.68.0', value: 'curl/7.68.0' },
  ];

  if (loading) {
    return (
      <div className="h-full flex items-center justify-center bg-gray-950">
        <div className="text-gray-400">
          <div className="animate-spin text-4xl mb-4 text-center">⏳</div>
          <p>加载配置中...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-4xl mx-auto space-y-6">
        <div className="bg-gray-900 rounded-xl border border-gray-800 p-6">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h3 className="text-lg font-semibold text-white">代理设置</h3>
              <p className="text-sm text-gray-400 mt-1">配置HTTP/HTTPS/SOCKS5代理，用于扫描请求</p>
            </div>
            <div className="flex items-center gap-3">
              <span className={`text-sm ${proxyConfig.enabled ? 'text-green-400' : 'text-gray-500'}`}>
                {proxyConfig.enabled ? '✓ 已启用' : '○ 已禁用'}
              </span>
              <button
                onClick={saveProxyConfig}
                disabled={saving}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50"
              >
                {saving ? '保存中...' : '保存'}
              </button>
            </div>
          </div>

          <div className="space-y-4">
            <label className="flex items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                checked={proxyConfig.enabled}
                onChange={(e) => setProxyConfig({ ...proxyConfig, enabled: e.target.checked })}
                className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-blue-600 focus:ring-blue-500"
              />
              <span className="text-gray-300">启用代理</span>
            </label>

            {proxyConfig.enabled && (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pl-8">
                <div>
                  <label className="block text-sm font-medium text-gray-400 mb-1">代理类型</label>
                  <select
                    value={proxyConfig.type}
                    onChange={(e) => setProxyConfig({ ...proxyConfig, type: e.target.value })}
                    className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                  >
                    <option value="http">HTTP</option>
                    <option value="https">HTTPS</option>
                    <option value="socks5">SOCKS5</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-400 mb-1">主机地址</label>
                  <input
                    type="text"
                    value={proxyConfig.host}
                    onChange={(e) => setProxyConfig({ ...proxyConfig, host: e.target.value })}
                    className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                    placeholder="127.0.0.1"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-400 mb-1">端口</label>
                  <input
                    type="number"
                    value={proxyConfig.port}
                    onChange={(e) => setProxyConfig({ ...proxyConfig, port: parseInt(e.target.value) || 0 })}
                    className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                    placeholder="8080"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-400 mb-1">用户名 (可选)</label>
                  <input
                    type="text"
                    value={proxyConfig.username}
                    onChange={(e) => setProxyConfig({ ...proxyConfig, username: e.target.value })}
                    className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                    placeholder="username"
                  />
                </div>

                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-gray-400 mb-1">密码 (可选)</label>
                  <input
                    type="password"
                    value={proxyConfig.password}
                    onChange={(e) => setProxyConfig({ ...proxyConfig, password: e.target.value })}
                    className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                    placeholder="password"
                  />
                </div>
              </div>
            )}
          </div>
        </div>

        <div className="bg-gray-900 rounded-xl border border-gray-800 p-6">
          <div className="flex items-center justify-between mb-6">
            <div>
              <h3 className="text-lg font-semibold text-white">扫描策略</h3>
              <p className="text-sm text-gray-400 mt-1">配置并发、超时、重试等扫描参数</p>
            </div>
            <button
              onClick={saveNetworkConfig}
              disabled={saving}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition disabled:opacity-50"
            >
              {saving ? '保存中...' : '保存'}
            </button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-400 mb-1">并发数 (1-100)</label>
              <input
                type="number"
                min="1"
                max="100"
                value={networkConfig.concurrency}
                onChange={(e) => setNetworkConfig({ ...networkConfig, concurrency: parseInt(e.target.value) || 1 })}
                className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
              />
              <p className="text-xs text-gray-500 mt-1">同时进行的请求数量</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-400 mb-1">请求限制 (QPS)</label>
              <input
                type="number"
                min="1"
                max="1000"
                value={networkConfig.rate_limit}
                onChange={(e) => setNetworkConfig({ ...networkConfig, rate_limit: parseInt(e.target.value) || 1 })}
                className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
              />
              <p className="text-xs text-gray-500 mt-1">每秒最大请求数</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-400 mb-1">超时时间 (秒)</label>
              <input
                type="number"
                min="1"
                max="300"
                value={networkConfig.timeout}
                onChange={(e) => setNetworkConfig({ ...networkConfig, timeout: parseInt(e.target.value) || 30 })}
                className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
              />
              <p className="text-xs text-gray-500 mt-1">单个请求超时时间</p>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-400 mb-1">重试次数</label>
              <input
                type="number"
                min="0"
                max="10"
                value={networkConfig.retries}
                onChange={(e) => setNetworkConfig({ ...networkConfig, retries: parseInt(e.target.value) || 0 })}
                className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
              />
              <p className="text-xs text-gray-500 mt-1">失败后重试次数</p>
            </div>

            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-400 mb-1">User-Agent</label>
              <select
                value={networkConfig.user_agent}
                onChange={(e) => setNetworkConfig({ ...networkConfig, user_agent: e.target.value })}
                className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none mb-2"
              >
                {userAgents.map((ua) => (
                  <option key={ua.value} value={ua.value}>{ua.label}</option>
                ))}
                <option value="custom">自定义...</option>
              </select>
              {networkConfig.user_agent === 'custom' && (
                <input
                  type="text"
                  placeholder="输入自定义 User-Agent"
                  value={networkConfig.user_agent}
                  onChange={(e) => setNetworkConfig({ ...networkConfig, user_agent: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
                />
              )}
            </div>

            <div className="md:col-span-2 space-y-3">
              <label className="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  checked={networkConfig.follow_redirects}
                  onChange={(e) => setNetworkConfig({ ...networkConfig, follow_redirects: e.target.checked })}
                  className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-blue-600 focus:ring-blue-500"
                />
                <div>
                  <span className="text-gray-300">跟随重定向</span>
                  <p className="text-xs text-gray-500">自动跟随HTTP 3xx重定向</p>
                </div>
              </label>

              <label className="flex items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  checked={networkConfig.skip_tls_verify}
                  onChange={(e) => setNetworkConfig({ ...networkConfig, skip_tls_verify: e.target.checked })}
                  className="w-5 h-5 rounded bg-gray-700 border-gray-600 text-blue-600 focus:ring-blue-500"
                />
                <div>
                  <span className="text-gray-300">跳过TLS验证</span>
                  <p className="text-xs text-yellow-500">⚠️ 不安全：跳过SSL证书验证</p>
                </div>
              </label>
            </div>
          </div>
        </div>

        <div className="bg-gray-900 rounded-xl border border-gray-800 p-6">
          <div className="flex items-center justify-between mb-4">
            <div>
              <h3 className="text-lg font-semibold text-white">连接测试</h3>
              <p className="text-sm text-gray-400 mt-1">测试当前代理和网络配置是否正常工作</p>
            </div>
          </div>
          <div className="flex gap-3 mb-4">
            <input
              type="text"
              className="flex-1 px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
              placeholder="https://example.com"
              value={testUrl}
              onChange={(e) => setTestUrl(e.target.value)}
            />
            <button
              onClick={testConnection}
              disabled={testing}
              className="px-6 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition disabled:opacity-50"
            >
              {testing ? '测试中...' : '测试'}
            </button>
          </div>
          {testResult && (
            <div className={`p-4 rounded-lg ${testResult.success ? 'bg-green-600/20 text-green-400 border border-green-600/30' : 'bg-red-600/20 text-red-400 border border-red-600/30'}`}>
              <div className="flex items-center gap-2">
                {testResult.success ? '✅' : '❌'}
                <span className="font-medium">{testResult.message}</span>
              </div>
              {testResult.status && (
                <p className="text-sm mt-1">状态码: {testResult.status}</p>
              )}
              {testResult.duration && (
                <p className="text-sm mt-1">耗时: {testResult.duration}ms</p>
              )}
            </div>
          )}
        </div>

        <div className="bg-gray-900 rounded-xl border border-gray-800 p-6">
          <h3 className="text-lg font-semibold text-white mb-4">快速配置</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            <button
              onClick={() => {
                setNetworkConfig({
                  ...networkConfig,
                  concurrency: 5,
                  rate_limit: 10,
                  timeout: 60,
                  retries: 5,
                });
              }}
              className="p-3 bg-gray-800 hover:bg-gray-700 rounded-lg text-left transition"
            >
              <div className="text-white font-medium">🛡️ 安全模式</div>
              <div className="text-xs text-gray-400 mt-1">低并发、慢速扫描</div>
            </button>
            <button
              onClick={() => {
                setNetworkConfig({
                  ...networkConfig,
                  concurrency: 20,
                  rate_limit: 100,
                  timeout: 30,
                  retries: 3,
                });
              }}
              className="p-3 bg-gray-800 hover:bg-gray-700 rounded-lg text-left transition"
            >
              <div className="text-white font-medium">⚡ 标准模式</div>
              <div className="text-xs text-gray-400 mt-1">平衡速度与安全</div>
            </button>
            <button
              onClick={() => {
                setNetworkConfig({
                  ...networkConfig,
                  concurrency: 50,
                  rate_limit: 500,
                  timeout: 15,
                  retries: 1,
                });
              }}
              className="p-3 bg-gray-800 hover:bg-gray-700 rounded-lg text-left transition"
            >
              <div className="text-white font-medium">🚀 高速模式</div>
              <div className="text-xs text-gray-400 mt-1">高并发快速扫描</div>
            </button>
            <button
              onClick={() => {
                setNetworkConfig({
                  ...networkConfig,
                  concurrency: 100,
                  rate_limit: 1000,
                  timeout: 10,
                  retries: 0,
                });
              }}
              className="p-3 bg-gray-800 hover:bg-gray-700 rounded-lg text-left transition"
            >
              <div className="text-white font-medium">🔥 极速模式</div>
              <div className="text-xs text-gray-400 mt-1">最大速度（可能触发WAF）</div>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NetworkSettings;
