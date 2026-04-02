import React, { useState } from 'react';

const ToolManager: React.FC = () => {
  const [tools, setTools] = useState<any[]>([
    { id: '1', name: 'Nuclei', type: 'scanner', status: 'active' },
    { id: '2', name: 'Xray', type: 'scanner', status: 'active' },
    { id: '3', name: 'Subfinder', type: 'recon', status: 'active' },
  ]);

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-6xl mx-auto">
        <h2 className="text-2xl font-bold text-white mb-6">工具管理</h2>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {tools.map(tool => (
            <div key={tool.id} className="bg-gray-900 rounded-lg border border-gray-800 p-4">
              <div className="flex items-center justify-between mb-2">
                <h3 className="text-white font-medium">{tool.name}</h3>
                <span className={`px-2 py-1 rounded text-xs ${
                  tool.status === 'active' ? 'bg-green-600/20 text-green-400' : 'bg-gray-600/20 text-gray-400'
                }`}>{tool.status}</span>
              </div>
              <p className="text-sm text-gray-400">{tool.type}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ToolManager;
