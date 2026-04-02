import React, { useState, useEffect } from 'react';
import { eduApi, activeScanApi } from '../services/api';

interface University {
  id: string;
  name: string;
  short_name: string;
  province: string;
  city: string;
  type: string;
  level: string;
  domains: string[];
  ips: string[];
  has_src: boolean;
  src_url: string;
  enabled: boolean;
  last_scan_at: string | null;
  created_at: string;
  updated_at: string;
}

const DEFAULT_UNIVERSITIES: Omit<University, 'id' | 'created_at' | 'updated_at' | 'last_scan_at'>[] = [
  { name: '清华大学', short_name: '清华', province: '北京', city: '北京', type: '综合', level: '985', domains: ['tsinghua.edu.cn', 'www.tsinghua.edu.cn', 'mail.tsinghua.edu.cn', 'lib.tsinghua.edu.cn'], ips: [], has_src: true, src_url: 'https://security.tsinghua.edu.cn', enabled: true },
  { name: '北京大学', short_name: '北大', province: '北京', city: '北京', type: '综合', level: '985', domains: ['pku.edu.cn', 'www.pku.edu.cn', 'mail.pku.edu.cn', 'library.pku.edu.cn'], ips: [], has_src: true, src_url: 'https://security.pku.edu.cn', enabled: true },
  { name: '浙江大学', short_name: '浙大', province: '浙江', city: '杭州', type: '综合', level: '985', domains: ['zju.edu.cn', 'www.zju.edu.cn', 'mail.zju.edu.cn', 'cz.zju.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '复旦大学', short_name: '复旦', province: '上海', city: '上海', type: '综合', level: '985', domains: ['fudan.edu.cn', 'www.fudan.edu.cn', 'mail.fudan.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '上海交通大学', short_name: '上交', province: '上海', city: '上海', type: '综合', level: '985', domains: ['sjtu.edu.cn', 'www.sjtu.edu.cn', 'mail.sjtu.edu.cn', 'pan.sjtu.edu.cn'], ips: [], has_src: true, src_url: 'https://security.sjtu.edu.cn', enabled: true },
  { name: '南京大学', short_name: '南大', province: '江苏', city: '南京', type: '综合', level: '985', domains: ['nju.edu.cn', 'www.nju.edu.cn', 'mail.nju.edu.cn', 'lib.nju.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '中国科学技术大学', short_name: '中科大', province: '安徽', city: '合肥', type: '理工', level: '985', domains: ['ustc.edu.cn', 'www.ustc.edu.cn', 'mail.ustc.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '哈尔滨工业大学', short_name: '哈工大', province: '黑龙江', city: '哈尔滨', type: '理工', level: '985', domains: ['hit.edu.cn', 'www.hit.edu.cn', 'mail.hit.edu.cn', 'today.hit.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '西安交通大学', short_name: '西交', province: '陕西', city: '西安', type: '综合', level: '985', domains: ['xjtu.edu.cn', 'www.xjtu.edu.cn', 'mail.xjtu.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '北京航空航天大学', short_name: '北航', province: '北京', city: '北京', type: '理工', level: '985', domains: ['buaa.edu.cn', 'www.buaa.edu.cn', 'mail.buaa.edu.cn', 'lib.buaa.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '同济大学', short_name: '同济', province: '上海', city: '上海', type: '理工', level: '985', domains: ['tongji.edu.cn', 'www.tongji.edu.cn', 'mail.tongji.edu.cn', 'lib.tongji.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '华中科技大学', short_name: '华科', province: '湖北', city: '武汉', type: '综合', level: '985', domains: ['hust.edu.cn', 'www.hust.edu.cn', 'mail.hust.edu.cn', 'lib.hust.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '武汉大学', short_name: '武大', province: '湖北', city: '武汉', type: '综合', level: '985', domains: ['whu.edu.cn', 'www.whu.edu.cn', 'mail.whu.edu.cn', 'lib.whu.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '中山大学', short_name: '中大', province: '广东', city: '广州', type: '综合', level: '985', domains: ['sysu.edu.cn', 'www.sysu.edu.cn', 'mail.sysu.edu.cn', 'library.sysu.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '四川大学', short_name: '川大', province: '四川', city: '成都', type: '综合', level: '985', domains: ['scu.edu.cn', 'www.scu.edu.cn', 'mail.scu.edu.cn', 'lib.scu.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '电子科技大学', short_name: '成电', province: '四川', city: '成都', type: '理工', level: '985', domains: ['uestc.edu.cn', 'www.uestc.edu.cn', 'mail.uestc.edu.cn', 'lib.uestc.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '东南大学', short_name: '东大', province: '江苏', city: '南京', type: '综合', level: '985', domains: ['seu.edu.cn', 'www.seu.edu.cn', 'mail.seu.edu.cn', 'lib.seu.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '山东大学', short_name: '山大', province: '山东', city: '济南', type: '综合', level: '985', domains: ['sdu.edu.cn', 'www.sdu.edu.cn', 'mail.sdu.edu.cn', 'lib.sdu.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '厦门大学', short_name: '厦大', province: '福建', city: '厦门', type: '综合', level: '985', domains: ['xmu.edu.cn', 'www.xmu.edu.cn', 'mail.xmu.edu.cn', 'lib.xmu.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '中国人民大学', short_name: '人大', province: '北京', city: '北京', type: '综合', level: '985', domains: ['ruc.edu.cn', 'www.ruc.edu.cn', 'mail.ruc.edu.cn', 'lib.ruc.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '北京理工大学', short_name: '北理', province: '北京', city: '北京', type: '理工', level: '985', domains: ['bit.edu.cn', 'www.bit.edu.cn', 'mail.bit.edu.cn', 'lib.bit.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '北京邮电大学', short_name: '北邮', province: '北京', city: '北京', type: '理工', level: '211', domains: ['bupt.edu.cn', 'www.bupt.edu.cn', 'mail.bupt.edu.cn', 'lib.bupt.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '华南理工大学', short_name: '华工', province: '广东', city: '广州', type: '理工', level: '985', domains: ['scut.edu.cn', 'www.scut.edu.cn', 'mail.scut.edu.cn', 'lib.scut.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '天津大学', short_name: '天大', province: '天津', city: '天津', type: '理工', level: '985', domains: ['tju.edu.cn', 'www.tju.edu.cn', 'mail.tju.edu.cn', 'lib.tju.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '南开大学', short_name: '南开', province: '天津', city: '天津', type: '综合', level: '985', domains: ['nankai.edu.cn', 'www.nankai.edu.cn', 'mail.nankai.edu.cn', 'lib.nankai.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '重庆大学', short_name: '重大', province: '重庆', city: '重庆', type: '综合', level: '985', domains: ['cqu.edu.cn', 'www.cqu.edu.cn', 'mail.cqu.edu.cn', 'lib.cqu.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '中南大学', short_name: '中南', province: '湖南', city: '长沙', type: '综合', level: '985', domains: ['csu.edu.cn', 'www.csu.edu.cn', 'mail.csu.edu.cn', 'lib.csu.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
  { name: '大连理工大学', short_name: '大工', province: '辽宁', city: '大连', type: '理工', level: '985', domains: ['dlut.edu.cn', 'www.dlut.edu.cn', 'mail.dlut.edu.cn', 'lib.dlut.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '吉林大学', short_name: '吉大', province: '吉林', city: '长春', type: '综合', level: '985', domains: ['jlu.edu.cn', 'www.jlu.edu.cn', 'mail.jlu.edu.cn', 'lib.jlu.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '湖南大学', short_name: '湖大', province: '湖南', city: '长沙', type: '综合', level: '985', domains: ['hnu.edu.cn', 'www.hnu.edu.cn', 'mail.hnu.edu.cn', 'lib.hnu.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '华东师范大学', short_name: '华师大', province: '上海', city: '上海', type: '师范', level: '985', domains: ['ecnu.edu.cn', 'www.ecnu.edu.cn', 'mail.ecnu.edu.cn', 'lib.ecnu.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '兰州大学', short_name: '兰大', province: '甘肃', city: '兰州', type: '综合', level: '985', domains: ['lzu.edu.cn', 'www.lzu.edu.cn', 'mail.lzu.edu.cn', 'lib.lzu.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '北京师范大学', short_name: '北师大', province: '北京', city: '北京', type: '师范', level: '985', domains: ['bnu.edu.cn', 'www.bnu.edu.cn', 'mail.bnu.edu.cn', 'lib.bnu.edu.cn'], ips: [], has_src: false, src_url: '', enabled: true },
  { name: '西北工业大学', short_name: '西工大', province: '陕西', city: '西安', type: '理工', level: '985', domains: ['nwpu.edu.cn', 'www.nwpu.edu.cn', 'mail.nwpu.edu.cn', 'lib.nwpu.edu.cn'], ips: [], has_src: true, src_url: '', enabled: true },
];

const EduUniversityManager: React.FC = () => {
  const [universities, setUniversities] = useState<University[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedProvince, setSelectedProvince] = useState('');
  const [selectedLevel, setSelectedLevel] = useState('');
  const [searchKeyword, setSearchKeyword] = useState('');
  const [showAddModal, setShowAddModal] = useState(false);
  const [selectedUni, setSelectedUni] = useState<University | null>(null);
  const [scanningId, setScanningId] = useState<string | null>(null);
  const [newUniversity, setNewUniversity] = useState({
    name: '',
    short_name: '',
    province: '',
    city: '',
    type: '综合',
    level: '本科',
    domains: '',
    has_src: false,
    src_url: '',
    enabled: true,
  });

  useEffect(() => {
    loadUniversities();
  }, []);

  const loadUniversities = async () => {
    setLoading(true);
    try {
      let uniList: University[] = [];
      try {
        const data: any = await eduApi.getUniversities();
        uniList = (data?.data || data || []).map((u: any) => ({
          ...u,
          domains: Array.isArray(u.domains) ? u.domains : (typeof u.domains === 'string' ? JSON.parse(u.domains || '[]') : []),
          ips: Array.isArray(u.ips) ? u.ips : (typeof u.ips === 'string' ? JSON.parse(u.ips || '[]') : []),
        }));
      } catch (e) {
        console.warn('后端eduApi暂不可用，使用本地数据');
      }

      if (uniList.length === 0) {
        const savedData = localStorage.getItem('eduUniversities');
        if (savedData) {
          uniList = JSON.parse(savedData);
        }
      }

      if (uniList.length === 0) {
        uniList = DEFAULT_UNIVERSITIES.map((u, idx) => ({
          ...u,
          id: `default-${idx}`,
          last_scan_at: null,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        }));
        localStorage.setItem('eduUniversities', JSON.stringify(uniList));
      }
      setUniversities(uniList);
    } catch (error) {
      console.error('加载高校列表失败:', error);
      const savedData = localStorage.getItem('eduUniversities');
      if (savedData) setUniversities(JSON.parse(savedData));
    } finally {
      setLoading(false);
    }
  };

  const handleStartScan = async (uni: University) => {
    if (!confirm(`确定要扫描 ${uni.name} 的域名吗？\n将扫描以下域名:\n${uni.domains.join('\n')}`)) return;

    setScanningId(uni.id);
    try {
      const primaryDomain = uni.domains[0] || '';
      const result = await activeScanApi.createTask({
        name: `高校扫描 - ${uni.name}`,
        target: primaryDomain.startsWith('http') ? primaryDomain : `https://${primaryDomain}`,
        type: 'vuln',
        options: {
          severity: 'all',
          tags: ['education', 'university'],
          timeout: 60,
          concurrency: 15,
        },
      });
      alert(`✅ 扫描任务已创建！\n任务ID: ${result?.id || '未知'}\n目标: ${primaryDomain}\n请到【主动扫描】模块查看进度`);
    } catch (error) {
      console.error('启动扫描失败:', error);
      alert('❌ 启动扫描失败，请检查后端服务是否运行');
    } finally {
      setScanningId(null);
    }
  };

  const handleBatchScan = async () => {
    const enabledUnis = universities.filter(u => u.enabled && u.domains.length > 0);
    if (!confirm(`确定要批量扫描 ${enabledUnis.length} 所已启用的高校吗？`)) return;

    alert(`🚀 开始批量扫描 ${enabledUnis.length} 所高校...\n\n每个域名将创建一个漏洞扫描任务\n请到【主动扫描】模块查看所有任务`);

    for (const uni of enabledUnis.slice(0, 10)) {
      try {
        const domain = uni.domains[0];
        await activeScanApi.createTask({
          name: `批量扫描 - ${uni.name}`,
          target: domain.startsWith('http') ? domain : `https://${domain}`,
          type: 'web',
          options: { timeout: 30, concurrency: 10 },
        });
      } catch (e) {
        console.warn(`${uni.name} 扫描创建失败:`, e);
      }
      await new Promise(r => setTimeout(r, 500));
    }
    alert(`✅ 已提交 ${Math.min(enabledUnis.length, 10)} 个扫描任务！`);
  };

  const handleAddUniversity = () => {
    if (!newUniversity.name || !newUniversity.province) {
      alert('请填写学校名称和省份');
      return;
    }
    const parsedDomains = (() => { try { return JSON.parse(newUniversity.domains || '[]'); } catch { return newUniversity.domains.split(',').map((d: string) => d.trim()).filter(Boolean); } })();

    const newUni: University = {
      id: `custom-${Date.now()}`,
      name: newUniversity.name,
      short_name: newUniversity.short_name,
      province: newUniversity.province,
      city: newUniversity.city,
      type: newUniversity.type,
      level: newUniversity.level,
      domains: parsedDomains,
      ips: [],
      has_src: newUniversity.has_src,
      src_url: newUniversity.src_url,
      enabled: newUniversity.enabled !== false,
      last_scan_at: null,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
    };
    const updated = [...universities, newUni];
    setUniversities(updated);
    localStorage.setItem('eduUniversities', JSON.stringify(updated));
    setShowAddModal(false);
    setNewUniversity({ name: '', short_name: '', province: '', city: '', type: '综合', level: '本科', domains: '', has_src: false, src_url: '', enabled: true });
    alert(`✅ 已添加 ${newUni.name}`);
  };

  const handleDeleteUniversity = (id: string) => {
    if (!confirm('确定要删除这所高校吗？')) return;
    const updated = universities.filter(u => u.id !== id);
    setUniversities(updated);
    localStorage.setItem('eduUniversities', JSON.stringify(updated));
  };

  const handleExport = () => {
    const allDomains = universities.flatMap(u => u.domains.map(d => ({ university: u.name, domain: d, province: u.province, level: u.level, has_src: u.has_src })));
    const headers = '学校名称\t域名\t省份\t级别\t有SRC\n';
    const rows = allDomains.map(d => `${d.university}\t${d.domain}\t${d.province}\t${d.level}\t${d.has_src}`).join('\n');
    const blob = new Blob([headers + rows], { type: 'text/plain;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `高校域名库_${new Date().toISOString().slice(0, 10)}.tsv`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const filteredUniversities = universities.filter(uni => {
    if (selectedProvince && uni.province !== selectedProvince) return false;
    if (selectedLevel && uni.level !== selectedLevel) return false;
    if (searchKeyword) {
      const kw = searchKeyword.toLowerCase();
      return uni.name.toLowerCase().includes(kw) || (uni.short_name?.toLowerCase().includes(kw)) || uni.domains.some(d => d.toLowerCase().includes(kw));
    }
    return true;
  });

  const provinces = [...new Set(universities.map(u => u.province).filter(Boolean))];
  const levels = [...new Set(universities.map(u => u.level).filter(Boolean))];
  const stats = { total: universities.length, enabled: universities.filter(u => u.enabled).length, withSrc: universities.filter(u => u.has_src).length, totalDomains: universities.reduce((s, u) => s + u.domains.length, 0) };

  const getLevelBadge = (level: string) => {
    const colors: Record<string, string> = {
      '985': 'bg-red-600/20 text-red-400 border-red-500/30',
      '211': 'bg-purple-600/20 text-purple-400 border-purple-500/30',
      '本科': 'bg-green-600/20 text-green-400 border-green-500/30',
      '专科': 'bg-blue-600/20 text-blue-400 border-blue-500/30',
    };
    return <span className={`px-2 py-0.5 rounded-full text-xs font-medium border ${colors[level] || 'bg-gray-600/20 text-gray-400 border-gray-500/30'}`}>{level}</span>;
  };

  return (
    <div className="h-full flex flex-col bg-gray-950">
      <div className="p-4 bg-gray-900 border-b border-gray-800">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-lg font-semibold text-white">🎓 高校域名库</h2>
            <p className="text-xs text-gray-400 mt-0.5">共 {stats.total} 所高校 · {stats.totalDomains} 个域名 · {stats.withSrc} 所有SRC</p>
          </div>
          <div className="flex space-x-2">
            <button onClick={handleBatchScan} className="px-3 py-1.5 bg-orange-600 hover:bg-orange-500 text-white rounded-lg text-sm transition" title="批量扫描所有已启用高校">
              🚀 批量扫描
            </button>
            <button onClick={handleExport} className="px-3 py-1.5 bg-gray-700 text-gray-200 rounded-lg text-sm hover:bg-gray-600 transition">
              📤 导出
            </button>
            <button onClick={() => setShowAddModal(true)} className="px-3 py-1.5 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-500 transition">
              + 添加
            </button>
          </div>
        </div>

        {/* 统计条 */}
        <div className="grid grid-cols-4 gap-2 mt-3">
          {[{ label: '总计', value: stats.total, icon: '🎓' }, { label: '已启用', value: stats.enabled, icon: '✅' }, { label: '有SRC', value: stats.withSrc, icon: '🏆' }, { label: '总域名', value: stats.totalDomains, icon: '🌐' }].map(s => (
            <div key={s.label} className="bg-gray-800/50 rounded px-3 py-1.5 text-center">
              <span className="text-lg font-bold text-white">{s.value}</span>
              <span className="text-xs text-gray-500 ml-1">{s.icon} {s.label}</span>
            </div>
          ))}
        </div>
      </div>

      <div className="p-3 bg-gray-900/50 border-b border-gray-800">
        <div className="flex gap-3">
          <input type="text" value={searchKeyword} onChange={(e) => setSearchKeyword(e.target.value)} placeholder="🔍 搜索高校名称/简称/域名..." className="flex-1 min-w-[180px] px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-sm text-white placeholder-gray-500 focus:border-cyan-500 outline-none" />
          <select value={selectedProvince} onChange={(e) => setSelectedProvince(e.target.value)} className="px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-sm text-white outline-none">
            <option value="">全部省份 ({provinces.length})</option>
            {provinces.map(p => <option key={p} value={p}>{p}</option>)}
          </select>
          <select value={selectedLevel} onChange={(e) => setSelectedLevel(e.target.value)} className="px-3 py-1.5 bg-gray-800 border border-gray-700 rounded-lg text-sm text-white outline-none">
            <option value="">全部级别</option>
            {levels.map(l => <option key={l} value={l}>{l}</option>)}
          </select>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-4">
        {loading ? (
          <div className="flex items-center justify-center h-full"><div className="animate-spin w-8 h-8 border-4 border-cyan-500 border-t-transparent rounded-full"></div></div>
        ) : filteredUniversities.length === 0 ? (
          <div className="text-center py-12 text-gray-500">
            <div className="text-4xl mb-4">🎓</div>
            <p>无匹配结果</p>
            <button onClick={() => { setSelectedProvince(''); setSelectedLevel(''); setSearchKeyword(''); }} className="mt-2 text-cyan-400 text-sm underline">清除筛选</button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-3">
            {filteredUniversities.map((uni) => (
              <div key={uni.id} className={`bg-gray-900 rounded-lg border p-3 hover:border-gray-600 transition group ${!uni.enabled ? 'opacity-50' : ''}`}>
                <div className="flex items-start justify-between mb-2">
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-white text-sm truncate">{uni.name}</h3>
                    {uni.short_name && <p className="text-xs text-gray-500">{uni.short_name}</p>}
                  </div>
                  <div className="flex items-center gap-1 ml-2 shrink-0">
                    {getLevelBadge(uni.level)}
                    {uni.has_src && <span className="px-1.5 py-0.5 bg-green-600/20 text-green-400 rounded text-[10px] font-bold">SRC</span>}
                  </div>
                </div>

                <div className="text-xs text-gray-500 mb-2">{uni.province}{uni.city ? ` · ${uni.city}` : ''} · {uni.type}</div>

                {uni.domains.length > 0 && (
                  <div className="mb-2">
                    <p className="text-[10px] text-gray-600 mb-1">域名 ({uni.domains.length})</p>
                    <div className="flex flex-wrap gap-1">
                      {uni.domains.slice(0, 3).map((domain, i) => (
                        <span key={i} className="px-1.5 py-0.5 bg-gray-800 text-cyan-300 rounded text-[10px] font-mono cursor-pointer hover:bg-cyan-900/30" onClick={() => navigator.clipboard.writeText(domain)} title="点击复制">{domain}</span>
                      ))}
                      {uni.domains.length > 3 && <span className="px-1.5 py-0.5 text-gray-600 text-[10px]">+{uni.domains.length - 3}</span>}
                    </div>
                  </div>
                )}

                <div className="flex gap-1.5 mt-2 pt-2 border-t border-gray-800">
                  <button onClick={() => setSelectedUni(uni)} className="flex-1 px-2 py-1 text-[11px] bg-blue-600/20 text-blue-400 rounded hover:bg-blue-600/30 transition">详情</button>
                  <button onClick={() => handleStartScan(uni)} disabled={scanningId === uni.id || !uni.enabled} className={`flex-1 px-2 py-1 text-[11px] rounded transition ${scanningId === uni.id ? 'bg-yellow-600/30 text-yellow-400 animate-pulse' : 'bg-red-600/20 text-red-400 hover:bg-red-600/30'} disabled:opacity-30`}>
                    {scanningId === uni.id ? '⏳ 扫描中...' : '🔍 扫描'}
                  </button>
                  <button onClick={() => handleDeleteUniversity(uni.id)} className="px-2 py-1 text-[11px] bg-gray-700/50 text-gray-500 rounded hover:bg-red-600/20 hover:text-red-400 transition">✕</button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {showAddModal && (
        <div className="fixed inset-0 bg-black/60 flex items-center justify-center z-50 p-4" onClick={() => setShowAddModal(false)}>
          <div className="bg-gray-900 rounded-xl border border-gray-800 p-5 w-full max-w-lg max-h-[90vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-white">添加高校</h3>
              <button onClick={() => setShowAddModal(false)} className="text-gray-400 hover:text-white text-xl">&times;</button>
            </div>
            <div className="space-y-3">
              <div><label className="block text-xs text-gray-400 mb-1">学校名称 *</label><input type="text" value={newUniversity.name} onChange={(e) => setNewUniversity({ ...newUniversity, name: e.target.value })} className="w-full px-3 py-1.5 bg-gray-800 border border-gray-700 rounded text-sm text-white outline-none focus:border-cyan-500" placeholder="例如：XX大学" /></div>
              <div className="grid grid-cols-2 gap-3">
                <div><label className="block text-xs text-gray-400 mb-1">简称</label><input type="text" value={newUniversity.short_name} onChange={(e) => setNewUniversity({ ...newUniversity, short_name: e.target.value })} className="w-full px-3 py-1.5 bg-gray-800 border border-gray-700 rounded text-sm text-white outline-none" /></div>
                <div><label className="block text-xs text-gray-400 mb-1">省份 *</label><select value={newUniversity.province} onChange={(e) => setNewUniversity({ ...newUniversity, province: e.target.value })} className="w-full px-3 py-1.5 bg-gray-800 border border-gray-700 rounded text-sm text-white outline-none"><option value="">选择</option>{['北京','上海','天津','重庆','河北','山西','辽宁','吉林','黑龙江','江苏','浙江','安徽','福建','江西','山东','河南','湖北','湖南','广东','广西','海南','四川','贵州','云南','陕西','甘肃','内蒙古','新疆','西藏','宁夏'].map(p => <option key={p} value={p}>{p}</option>)}</select></div>
              </div>
              <div className="grid grid-cols-2 gap-3">
                <div><label className="block text-xs text-gray-400 mb-1">类型</label><select value={newUniversity.type} onChange={(e) => setNewUniversity({ ...newUniversity, type: e.target.value })} className="w-full px-3 py-1.5 bg-gray-800 border border-gray-700 rounded text-sm text-white outline-none">{['综合','理工','农林','医药','师范','财经','政法','体育','艺术','民族','军事'].map(t => <option key={t} value={t}>{t}</option>)}</select></div>
                <div><label className="block text-xs text-gray-400 mb-1">级别</label><select value={newUniversity.level} onChange={(e) => setNewUniversity({ ...newUniversity, level: e.target.value })} className="w-full px-3 py-1.5 bg-gray-800 border border-gray-700 rounded text-sm text-white outline-none">{['985','211','本科','专科'].map(l => <option key={l} value={l}>{l}</option>)}</select></div>
              </div>
              <div><label className="block text-xs text-gray-400 mb-1">域名列表 (逗号或JSON)</label><textarea value={newUniversity.domains} onChange={(e) => setNewUniversity({ ...newUniversity, domains: e.target.value })} className="w-full px-3 py-1.5 bg-gray-800 border border-gray-700 rounded text-sm text-white font-mono outline-none rows-2" placeholder="a.edu.cn, b.edu.cn 或 JSON数组" /></div>
              <div className="flex gap-4">
                <label className="flex items-center gap-1.5 text-sm text-gray-300 cursor-pointer"><input type="checkbox" checked={newUniversity.has_src} onChange={(e) => setNewUniversity({ ...newUniversity, has_src: e.target.checked })} className="rounded bg-gray-700" />有SRC平台</label>
                <label className="flex items-center gap-1.5 text-sm text-gray-300 cursor-pointer"><input type="checkbox" checked={newUniversity.enabled} onChange={(e) => setNewUniversity({ ...newUniversity, enabled: e.target.checked })} className="rounded bg-gray-700" />默认启用</label>
              </div>
            </div>
            <div className="flex justify-end gap-2 mt-4 pt-3 border-t border-gray-800">
              <button onClick={() => setShowAddModal(false)} className="px-4 py-1.5 bg-gray-700 text-gray-200 rounded-lg text-sm hover:bg-gray-600">取消</button>
              <button onClick={handleAddUniversity} disabled={!newUniversity.name || !newUniversity.province} className="px-4 py-1.5 bg-blue-600 text-white rounded-lg text-sm hover:bg-blue-500 disabled:opacity-50">保存</button>
            </div>
          </div>
        </div>
      )}

      {selectedUni && (
        <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50 p-4" onClick={() => setSelectedUni(null)}>
          <div className="bg-gray-900 rounded-xl border border-gray-700 w-full max-w-lg max-h-[85vh] overflow-hidden" onClick={(e) => e.stopPropagation()}>
            <div className="p-4 border-b border-gray-800 flex items-center justify-between">
              <div>
                <h3 className="text-lg font-bold text-white">{selectedUni.name}</h3>
                <p className="text-xs text-gray-400 mt-0.5">{selectedUni.short_name} · {selectedUni.province} · {getLevelBadge(selectedUni.level)}</p>
              </div>
              <button onClick={() => setSelectedUni(null)} className="text-gray-400 hover:text-white text-xl">&times;</button>
            </div>
            <div className="p-4 space-y-3 overflow-y-auto max-h-[calc(85vh-140px)]">
              <InfoRow label="学校类型" value={selectedUni.type} />
              <InfoRow label="所在城市" value={`${selectedUni.province} ${selectedUni.city || ''}`} />
              <InfoRow label="SRC平台" value={selectedUni.has_src ? (selectedUni.src_url || '是') : '无'} />
              <InfoRow label="状态" value={selectedUni.enabled ? '✅ 启用' : '❌ 禁用'} />
              <div>
                <p className="text-xs text-gray-500 mb-1">域名列表 ({selectedUni.domains.length})</p>
                <div className="space-y-1 max-h-40 overflow-y-auto">
                  {selectedUni.domains.map((d, i) => (
                    <div key={i} className="flex items-center justify-between bg-gray-800 rounded px-2 py-1 group">
                      <code className="text-xs text-cyan-300 font-mono">{d}</code>
                      <button onClick={() => navigator.clipboard.writeText(d)} className="text-[10px] text-gray-600 group-hover:text-cyan-400 opacity-0 group-hover:opacity-100 transition">复制</button>
                    </div>
                  ))}
                </div>
              </div>
              {selectedUni.src_url && (
                <div>
                  <p className="text-xs text-gray-500 mb-1">SRC地址</p>
                  <a href={selectedUni.src_url} target="_blank" rel="noopener noreferrer" className="text-sm text-blue-400 hover:underline break-all">{selectedUni.src_url}</a>
                </div>
              )}
              <div className="pt-2 border-t border-gray-800 flex gap-2">
                <button onClick={() => { handleStartScan(selectedUni); setSelectedUni(null); }} className="flex-1 py-2 bg-red-600 hover:bg-red-500 text-white rounded-lg text-sm font-medium transition">🔍 开始扫描此高校</button>
                <button onClick={() => { navigator.clipboard.writeText(selectedUni.domains.join('\n')); alert('域名列表已复制！'); }} className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg text-sm transition">📋 复制全部域名</button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

const InfoRow: React.FC<{ label: string; value: string }> = ({ label, value }) => (
  <div className="bg-gray-800/50 rounded px-3 py-2">
    <div className="text-xs text-gray-500">{label}</div>
    <div className="text-sm text-gray-200">{value || '-'}</div>
  </div>
);

export default EduUniversityManager;
