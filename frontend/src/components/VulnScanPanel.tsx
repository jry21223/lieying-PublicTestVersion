import React, { useState } from 'react';

const VulnScanPanel: React.FC = () => {
  const [tasks, setTasks] = useState<any[]>([]);

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      <div className="max-w-4xl mx-auto">
        <h2 className="text-2xl font-bold text-white mb-6">漏洞扫描任务</h2>
        {tasks.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            <div className="text-4xl mb-4">🔍</div>
            <p>暂无扫描任务</p>
          </div>
        ) : (
          <div className="space-y-3">
            {tasks.map(t => (
              <div key={t.id} className="bg-gray-900 rounded-lg border border-gray-800 p-4">
                <p className="text-white">{t.name}</p>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default VulnScanPanel;
