import React, { useState } from 'react';

const AssetDetail: React.FC = () => {
  const [asset, setAsset] = useState<any>(null);

  return (
    <div className="h-full overflow-auto bg-gray-950 p-6">
      {asset ? (
        <div className="max-w-4xl mx-auto">
          <h2 className="text-2xl font-bold text-white mb-6">{asset.name}</h2>
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-6">
            <pre className="text-gray-300 whitespace-pre-wrap">{JSON.stringify(asset, null, 2)}</pre>
          </div>
        </div>
      ) : (
        <div className="text-center py-12 text-gray-500">
          <div className="text-4xl mb-4">👆</div>
          <p>选择一个资产查看详情</p>
        </div>
      )}
    </div>
  );
};

export default AssetDetail;
