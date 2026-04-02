import React, { useState, useEffect } from 'react';
import { activeScanApi, targetsApi } from '../services/api';

interface Target {
  id: string;
  name: string;
  value: string;
  type: string;
}

const ScanPanel: React.FC = () => {
  const [targets, setTargets] = useState<Target[]>([]);
  const [selectedTargets, setSelectedTargets] = useState<string[]>([]);
  const [scanType, setScanType] = useState<'web' | 'port' | 'vuln' | 'dir' | 'fuzz'>('vuln');
  const [scanning, setScanning] = useState(false);
  const [progress, setProgress] = useState(0);
  const [currentTarget, setCurrentTarget] = useState('');
  const [scanResults, setScanResults] = useState<any[]>([]);
  const [ports, setPorts] = useState('21,22,23,25,53,80,110,143,443,445,993,995,3306,3389,5432,6379,8080,8443,27017');
  const [concurrency, setConcurrency] = useState(10);
  const [timeout, setTimeout] = useState(30);

  useEffect(() => {
    loadTargets();
  }, []);

  const loadTargets = async () => {
    try {
      const data = await targetsApi.list();
      setTargets(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('加载目标失败:', error);
    }
  };

  const toggleTarget = (id: string) => {
    setSelectedTargets(prev => 
      prev.includes(id) 
        ? prev.filter(t => t !== id)
        : [...prev, id]
    );
  };

  const selectAll = () => {
    if (selectedTargets.length === targets.length) {
      setSelectedTargets([]);
    } else {
      setSelectedTargets(targets.map(t => t.id));
    }
  };

  const startScan = async () => {
    if (selectedTargets.length === 0) {
      alert('请选择至少一个目标');
      return;
    }

    setScanning(true);
    setProgress(0);
    setScanResults([]);

    const total = selectedTargets.length;
    let completed = 0;

    for (const targetId of selectedTargets) {
      const target = targets.find(t => t.id === targetId);
      if (!target) continue;

      setCurrentTarget(target.value);
      
      try {
        const result = await activeScanApi.createTask({
          name: `批量扫描 - ${target.name}`,
          target: target.value,
          type: scanType,
          options: {
            ports: scanType === 'port' ? ports : undefined,
            concurrency,
            timeout,
          },
        });

        setScanResults(prev => [...prev, {
          target: target.value,
          success: true,
          taskId: result.id,
        }]);
      } catch (error) {
        console.error(`扫描 ${target.value} 失败:`, error);
        setScanResults(prev => [...prev, {
          target: target.value,
          success: false,
          error: String(error),
        }]);
      }

      completed++;
      setProgress(Math.round((completed / total) * 100));
    }

    setScanning(false);
    setCurrentTarget('');
    alert(`扫描完成！成功: ${scanResults.filter(r => r.success).length}, 失败: ${scanResults.filter(r => !r.success).length}`);
  };

  const scanTypes = [
    { value: 'web', label: 'Web扫描', desc: '基础Web信息收集' },
    { value: 'port', label: '端口扫描', desc: '检测开放端口' },
    { value: 'vuln', label: '漏洞扫描', desc: 'OWASP Top 10检测' },
    { value: 'dir', label: '目录扫描', desc: '敏感目录探测' },
    { value: 'fuzz', label: '模糊测试', desc: '参数模糊测试' },
  ];

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-4xl mx-auto">
        <h2 className="text-2xl font-bold text-white mb-6">漏洞扫描</h2>

        <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 mb-6">
          <h3 className="text-lg font-semibold text-white mb-4">扫描类型</h3>
          <div className="grid grid-cols-5 gap-3">
            {scanTypes.map(type => (
              <button
                key={type.value}
                onClick={() => setScanType(type.value as any)}
                className={`p-3 rounded-lg text-left transition ${
                  scanType === type.value
                    ? 'bg-blue-600/20 border border-blue-600/30'
                    : 'bg-gray-800 border border-gray-700 hover:border-gray-600'
                }`}
              >
                <div className="text-white font-medium">{type.label}</div>
                <div className="text-xs text-gray-400 mt-1">{type.desc}</div>
              </button>
            ))}
          </div>
        </div>

        {scanType === 'port' && (
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 mb-6">
            <h3 className="text-lg font-semibold text-white mb-4">端口配置</h3>
            <input
              type="text"
              value={ports}
              onChange={(e) => setPorts(e.target.value)}
              className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
              placeholder="端口列表，如: 80,443,8080 或 1-1000"
            />
          </div>
        )}

        <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 mb-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-white">选择目标 ({selectedTargets.length}/{targets.length})</h3>
            <button
              onClick={selectAll}
              className="text-sm text-blue-400 hover:text-blue-300"
            >
              {selectedTargets.length === targets.length ? '取消全选' : '全选'}
            </button>
          </div>

          {targets.length > 0 ? (
            <div className="max-h-60 overflow-y-auto space-y-2">
              {targets.map(target => (
                <label
                  key={target.id}
                  className="flex items-center gap-3 p-3 bg-gray-800 rounded-lg cursor-pointer hover:bg-gray-700 transition"
                >
                  <input
                    type="checkbox"
                    checked={selectedTargets.includes(target.id)}
                    onChange={() => toggleTarget(target.id)}
                    className="w-4 h-4 rounded bg-gray-700 border-gray-600"
                  />
                  <div className="flex-1">
                    <p className="text-white font-medium">{target.name}</p>
                    <p className="text-sm text-gray-400">{target.value}</p>
                  </div>
                  <span className="text-xs px-2 py-0.5 bg-gray-700 text-gray-300 rounded">
                    {target.type}
                  </span>
                </label>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">
              <p>暂无目标</p>
              <p className="text-sm mt-2">请先在左侧添加目标</p>
            </div>
          )}
        </div>

        <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 mb-6">
          <h3 className="text-lg font-semibold text-white mb-4">扫描选项</h3>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm text-gray-400 mb-1">并发数</label>
              <input
                type="number"
                value={concurrency}
                onChange={(e) => setConcurrency(parseInt(e.target.value) || 10)}
                className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
              />
            </div>
            <div>
              <label className="block text-sm text-gray-400 mb-1">超时时间(秒)</label>
              <input
                type="number"
                value={timeout}
                onChange={(e) => setTimeout(parseInt(e.target.value) || 30)}
                className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white focus:border-blue-500 focus:outline-none"
              />
            </div>
          </div>
        </div>

        {scanning && (
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 mb-6">
            <div className="flex items-center justify-between mb-2">
              <span className="text-white">扫描进度</span>
              <span className="text-gray-400">{progress}%</span>
            </div>
            <div className="w-full bg-gray-800 rounded-full h-2 mb-3">
              <div
                className="bg-blue-600 h-2 rounded-full transition-all"
                style={{ width: `${progress}%` }}
              />
            </div>
            {currentTarget && (
              <p className="text-sm text-gray-400">
                正在扫描: <span className="text-white">{currentTarget}</span>
              </p>
            )}
          </div>
        )}

        {scanResults.length > 0 && (
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 mb-6">
            <h3 className="text-lg font-semibold text-white mb-4">扫描结果</h3>
            <div className="space-y-2">
              {scanResults.map((result, index) => (
                <div
                  key={index}
                  className={`p-3 rounded-lg ${
                    result.success ? 'bg-green-600/10 border border-green-600/20' : 'bg-red-600/10 border border-red-600/20'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="text-white">{result.target}</span>
                    <span className={result.success ? 'text-green-400' : 'text-red-400'}>
                      {result.success ? '✓ 成功' : '✗ 失败'}
                    </span>
                  </div>
                  {result.taskId && (
                    <p className="text-xs text-gray-400 mt-1">任务ID: {result.taskId}</p>
                  )}
                  {result.error && (
                    <p className="text-xs text-red-400 mt-1">{result.error}</p>
                  )}
                </div>
              ))}
            </div>
          </div>
        )}

        <div className="flex gap-3">
          <button
            onClick={startScan}
            disabled={scanning || selectedTargets.length === 0}
            className="flex-1 px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-500 transition disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {scanning ? '扫描中...' : `开始扫描 (${selectedTargets.length} 个目标)`}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ScanPanel;
