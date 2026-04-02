import React, { useState } from 'react';

interface Proxy {
  id: string;
  type: 'http' | 'https' | 'socks5';
  host: string;
  port: number;
  username?: string;
  password?: string;
  enabled: boolean;
  latency?: number;
  lastCheck?: string;
}

const ProxyPool: React.FC = () => {
  const [proxies, setProxies] = useState<Proxy[]>(() => {
    const saved = localStorage.getItem('proxyPool');
    return saved ? JSON.parse(saved) : [];
  });
  const [showAddModal, setShowAddModal] = useState(false);
  const [editingProxy, setEditingProxy] = useState<Proxy | null>(null);
  const [newProxy, setNewProxy] = useState<Partial<Proxy>>({
    type: 'http',
    host: '',
    port: 8080,
    enabled: true,
  });

  const saveToStorage = (newProxies: Proxy[]) => {
    setProxies(newProxies);
    localStorage.setItem('proxyPool', JSON.stringify(newProxies));
  };

  const addProxy = () => {
    if (!newProxy.host || !newProxy.port) {
      alert('请填写代理主机和端口');
      return;
    }

    const proxy: Proxy = {
      id: Date.now().toString(),
      type: newProxy.type as 'http' | 'https' | 'socks5',
      host: newProxy.host!,
      port: newProxy.port!,
      username: newProxy.username,
      password: newProxy.password,
      enabled: newProxy.enabled !== false,
    };

    saveToStorage([...proxies, proxy]);
    setShowAddModal(false);
    setNewProxy({ type: 'http', host: '', port: 8080, enabled: true });
  };

  const updateProxy = () => {
    if (!editingProxy) return;

    const updated = proxies.map(p =>
      p.id === editingProxy.id ? editingProxy : p
    );
    saveToStorage(updated);
    setEditingProxy(null);
  };

  const deleteProxy = (id: string) => {
    if (!confirm('确定要删除这个代理吗？')) return;
    saveToStorage(proxies.filter(p => p.id !== id));
  };

  const toggleProxy = (id: string) => {
    const updated = proxies.map(p =>
      p.id === id ? { ...p, enabled: !p.enabled } : p
    );
    saveToStorage(updated);
  };

  const testProxy = async (proxy: Proxy) => {
    const startTime = Date.now();
    try {
      const response = await fetch(`http://${proxy.host}:${proxy.port}`, {
        method: 'HEAD',
        mode: 'no-cors',
      });
      const latency = Date.now() - startTime;

      const updated = proxies.map(p =>
        p.id === proxy.id ? { ...p, latency, lastCheck: new Date().toISOString() } : p
      );
      saveToStorage(updated);
      alert(`代理可用！延迟: ${latency}ms`);
    } catch (error) {
      const updated = proxies.map(p =>
        p.id === proxy.id ? { ...p, latency: -1, lastCheck: new Date().toISOString() } : p
      );
      saveToStorage(updated);
      alert('代理不可用');
    }
  };

  const importProxies = () => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = '.txt,.json';
    input.onchange = async (e: any) => {
      const file = e.target.files[0];
      if (!file) return;

      try {
        const content = await file.text();
        let importedProxies: Proxy[] = [];

        if (file.name.endsWith('.json')) {
          importedProxies = JSON.parse(content);
        } else {
          const lines = content.split('\n').filter(line => line.trim());
          importedProxies = lines.map(line => {
            const [host, port] = line.split(':');
            return {
              id: Date.now().toString() + Math.random(),
              type: 'http' as const,
              host: host.trim(),
              port: parseInt(port.trim()) || 8080,
              enabled: true,
            };
          });
        }

        saveToStorage([...proxies, ...importedProxies]);
        alert(`成功导入 ${importedProxies.length} 个代理`);
      } catch (error) {
        alert('导入失败，请检查文件格式');
      }
    };
    input.click();
  };

  const exportProxies = () => {
    const dataStr = JSON.stringify(proxies, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(dataBlob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `proxies-${new Date().toISOString().split('T')[0]}.json`;
    link.click();
  };

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-6xl mx-auto">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-2xl font-bold text-white">代理池管理</h2>
            <p className="text-sm text-gray-400 mt-1">管理扫描代理列表</p>
          </div>
          <div className="flex gap-2">
            <button
              onClick={importProxies}
              className="px-4 py-2 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 transition"
            >
              导入
            </button>
            <button
              onClick={exportProxies}
              className="px-4 py-2 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600 transition"
            >
              导出
            </button>
            <button
              onClick={() => setShowAddModal(true)}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500 transition"
            >
              添加代理
            </button>
          </div>
        </div>

        {proxies.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            <div className="text-4xl mb-4">🌐</div>
            <p>暂无代理</p>
            <p className="text-sm mt-2">点击"添加代理"开始添加</p>
          </div>
        ) : (
          <div className="bg-gray-900 rounded-xl border border-gray-800 overflow-hidden">
            <table className="w-full">
              <thead className="bg-gray-800">
                <tr className="text-left text-xs text-gray-400">
                  <th className="px-4 py-3 font-medium">状态</th>
                  <th className="px-4 py-3 font-medium">代理</th>
                  <th className="px-4 py-3 font-medium">类型</th>
                  <th className="px-4 py-3 font-medium">延迟</th>
                  <th className="px-4 py-3 font-medium">最后检查</th>
                  <th className="px-4 py-3 font-medium">操作</th>
                </tr>
              </thead>
              <tbody>
                {proxies.map(proxy => (
                  <tr key={proxy.id} className="border-t border-gray-800 hover:bg-gray-800/50">
                    <td className="px-4 py-3">
                      <button
                        onClick={() => toggleProxy(proxy.id)}
                        className={`w-12 h-6 rounded-full transition relative ${
                          proxy.enabled ? 'bg-green-600' : 'bg-gray-600'
                        }`}
                      >
                        <span
                          className={`absolute top-1 w-4 h-4 bg-white rounded-full transition ${
                            proxy.enabled ? 'right-1' : 'left-1'
                          }`}
                        />
                      </button>
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-2">
                        <span className={`text-sm ${proxy.enabled ? 'text-white' : 'text-gray-500'}`}>
                          {proxy.host}:{proxy.port}
                        </span>
                        {proxy.username && (
                          <span className="text-xs text-gray-500">({proxy.username})</span>
                        )}
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-1 rounded text-xs font-medium ${
                        proxy.type === 'http' ? 'bg-blue-600/20 text-blue-400' :
                        proxy.type === 'https' ? 'bg-green-600/20 text-green-400' :
                        'bg-purple-600/20 text-purple-400'
                      }`}>
                        {proxy.type.toUpperCase()}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      {proxy.latency !== undefined ? (
                        <span className={`text-sm ${
                          proxy.latency < 0 ? 'text-red-400' :
                          proxy.latency < 500 ? 'text-green-400' :
                          proxy.latency < 1000 ? 'text-yellow-400' : 'text-red-400'
                        }`}>
                          {proxy.latency < 0 ? '超时' : `${proxy.latency}ms`}
                        </span>
                      ) : (
                        <span className="text-gray-500 text-sm">-</span>
                      )}
                    </td>
                    <td className="px-4 py-3 text-sm text-gray-400">
                      {proxy.lastCheck ? new Date(proxy.lastCheck).toLocaleString() : '-'}
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex gap-2">
                        <button
                          onClick={() => testProxy(proxy)}
                          className="px-2 py-1 text-xs bg-blue-600/20 text-blue-400 rounded hover:bg-blue-600/30"
                        >
                          测试
                        </button>
                        <button
                          onClick={() => setEditingProxy(proxy)}
                          className="px-2 py-1 text-xs bg-gray-600/20 text-gray-400 rounded hover:bg-gray-600/30"
                        >
                          编辑
                        </button>
                        <button
                          onClick={() => deleteProxy(proxy.id)}
                          className="px-2 py-1 text-xs bg-red-600/20 text-red-400 rounded hover:bg-red-600/30"
                        >
                          删除
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      {showAddModal && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold text-white">添加代理</h3>
              <button onClick={() => setShowAddModal(false)} className="text-gray-400 hover:text-white">
                ✕
              </button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">类型</label>
                <select
                  value={newProxy.type}
                  onChange={(e) => setNewProxy({ ...newProxy, type: e.target.value as any })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
                >
                  <option value="http">HTTP</option>
                  <option value="https">HTTPS</option>
                  <option value="socks5">SOCKS5</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">主机 *</label>
                <input
                  type="text"
                  value={newProxy.host}
                  onChange={(e) => setNewProxy({ ...newProxy, host: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
                  placeholder="127.0.0.1"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">端口 *</label>
                <input
                  type="number"
                  value={newProxy.port}
                  onChange={(e) => setNewProxy({ ...newProxy, port: parseInt(e.target.value) || 0 })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
                  placeholder="8080"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">用户名 (可选)</label>
                <input
                  type="text"
                  value={newProxy.username || ''}
                  onChange={(e) => setNewProxy({ ...newProxy, username: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">密码 (可选)</label>
                <input
                  type="password"
                  value={newProxy.password || ''}
                  onChange={(e) => setNewProxy({ ...newProxy, password: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
                />
              </div>
            </div>
            <div className="flex justify-end gap-3 mt-6">
              <button
                onClick={() => setShowAddModal(false)}
                className="px-4 py-2 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600"
              >
                取消
              </button>
              <button
                onClick={addProxy}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500"
              >
                添加
              </button>
            </div>
          </div>
        </div>
      )}

      {editingProxy && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 w-full max-w-md">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-xl font-semibold text-white">编辑代理</h3>
              <button onClick={() => setEditingProxy(null)} className="text-gray-400 hover:text-white">
                ✕
              </button>
            </div>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">类型</label>
                <select
                  value={editingProxy.type}
                  onChange={(e) => setEditingProxy({ ...editingProxy, type: e.target.value as any })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
                >
                  <option value="http">HTTP</option>
                  <option value="https">HTTPS</option>
                  <option value="socks5">SOCKS5</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">主机</label>
                <input
                  type="text"
                  value={editingProxy.host}
                  onChange={(e) => setEditingProxy({ ...editingProxy, host: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">端口</label>
                <input
                  type="number"
                  value={editingProxy.port}
                  onChange={(e) => setEditingProxy({ ...editingProxy, port: parseInt(e.target.value) || 0 })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">用户名</label>
                <input
                  type="text"
                  value={editingProxy.username || ''}
                  onChange={(e) => setEditingProxy({ ...editingProxy, username: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1">密码</label>
                <input
                  type="password"
                  value={editingProxy.password || ''}
                  onChange={(e) => setEditingProxy({ ...editingProxy, password: e.target.value })}
                  className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
                />
              </div>
            </div>
            <div className="flex justify-end gap-3 mt-6">
              <button
                onClick={() => setEditingProxy(null)}
                className="px-4 py-2 bg-gray-700 text-gray-200 rounded-lg hover:bg-gray-600"
              >
                取消
              </button>
              <button
                onClick={updateProxy}
                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500"
              >
                保存
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ProxyPool;
