import React, { useState } from 'react';

interface TreeNode {
  id: string;
  name: string;
  type: 'folder' | 'domain' | 'host' | 'url';
  children?: TreeNode[];
}

const AssetTree: React.FC = () => {
  const [tree, setTree] = useState<TreeNode[]>([]);
  const [selected, setSelected] = useState<string | null>(null);

  const loadTree = () => {
    const saved = localStorage.getItem('assetTree');
    if (saved) setTree(JSON.parse(saved));
  };

  React.useEffect(() => { loadTree(); }, []);

  const renderNode = (node: TreeNode) => (
    <div key={node.id} className="ml-4 my-1">
      <div
        onClick={() => setSelected(node.id)}
        className={`flex items-center gap-2 p-2 rounded cursor-pointer ${
          selected === node.id ? 'bg-blue-600/20' : 'hover:bg-gray-800'
        }`}
      >
        <span>{node.type === 'folder' ? '📁' : node.type === 'domain' ? '🌐' : node.type === 'host' ? '💻' : '🔗'}</span>
        <span className="text-white text-sm">{node.name}</span>
      </div>
      {node.children?.map(renderNode)}
    </div>
  );

  return (
    <div className="h-full overflow-auto bg-gray-950 p-4">
      <h3 className="text-lg font-semibold text-white mb-4">资产树</h3>
      {tree.length === 0 ? (
        <div className="text-center py-8 text-gray-500">
          <div className="text-4xl mb-4">🌲</div>
          <p>暂无资产</p>
        </div>
      ) : (
        <div>{tree.map(renderNode)}</div>
      )}
    </div>
  );
};

export default AssetTree;
