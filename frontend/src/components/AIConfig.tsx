import React, { useState } from 'react';

const AIConfig: React.FC = () => {
  const [config, setConfig] = useState({
    provider: 'openai',
    apiKey: '',
    model: 'gpt-4',
    temperature: 0.7,
  });

  const saveConfig = () => {
    localStorage.setItem('aiConfig', JSON.stringify(config));
    alert('AI配置已保存！');
  };

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-2xl mx-auto">
        <h2 className="text-2xl font-bold text-white mb-6">AI配置</h2>
        <div className="bg-gray-900 rounded-xl border border-gray-800 p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">AI提供商</label>
            <select
              value={config.provider}
              onChange={(e) => setConfig({ ...config, provider: e.target.value })}
              className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
            >
              <option value="openai">OpenAI</option>
              <option value="claude">Claude</option>
              <option value="local">本地模型</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">API Key</label>
            <input
              type="password"
              value={config.apiKey}
              onChange={(e) => setConfig({ ...config, apiKey: e.target.value })}
              className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
              placeholder="sk-..."
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">模型</label>
            <input
              type="text"
              value={config.model}
              onChange={(e) => setConfig({ ...config, model: e.target.value })}
              className="w-full px-3 py-2 bg-gray-800 border border-gray-700 rounded-lg text-white"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-400 mb-1">Temperature: {config.temperature}</label>
            <input
              type="range"
              min="0"
              max="1"
              step="0.1"
              value={config.temperature}
              onChange={(e) => setConfig({ ...config, temperature: parseFloat(e.target.value) })}
              className="w-full"
            />
          </div>
          <button
            onClick={saveConfig}
            className="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-500"
          >
            保存配置
          </button>
        </div>
      </div>
    </div>
  );
};

export default AIConfig;
