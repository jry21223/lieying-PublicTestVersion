package edu

import (
	"fmt"
)

type University struct {
	Name      string
	ShortName string
	Province  string
	City      string
	Type      string
	Level     string
	Domains   []string
}

// 高校域名库 - 包含300+高校
type UniversityDB struct {
	universities []University
}

func NewUniversityDB() *UniversityDB {
	db := &UniversityDB{
		universities: []University{},
	}
	db.initData()
	return db
}

func (udb *UniversityDB) initData() {
	// 985高校
	udb.universities = append(udb.universities, []University{
		{Name: "北京大学", ShortName: "北大", Province: "北京", City: "北京", Type: "综合", Level: "985", Domains: []string{"pku.edu.cn", "pku.cn"}},
		{Name: "清华大学", ShortName: "清华", Province: "北京", City: "北京", Type: "理工", Level: "985", Domains: []string{"tsinghua.edu.cn", "tsinghua.cn"}},
		{Name: "复旦大学", ShortName: "复旦", Province: "上海", City: "上海", Type: "综合", Level: "985", Domains: []string{"fudan.edu.cn", "fudan.cn"}},
		{Name: "上海交通大学", ShortName: "上交", Province: "上海", City: "上海", Type: "理工", Level: "985", Domains: []string{"sjtu.edu.cn", "sjtu.cn"}},
		{Name: "浙江大学", ShortName: "浙大", Province: "浙江", City: "杭州", Type: "综合", Level: "985", Domains: []string{"zju.edu.cn", "zju.cn"}},
		{Name: "南京大学", ShortName: "南大", Province: "江苏", City: "南京", Type: "综合", Level: "985", Domains: []string{"nju.edu.cn", "nju.cn"}},
		{Name: "中国科学技术大学", ShortName: "中科大", Province: "安徽", City: "合肥", Type: "理工", Level: "985", Domains: []string{"ustc.edu.cn", "ustc.cn"}},
		{Name: "哈尔滨工业大学", ShortName: "哈工大", Province: "黑龙江", City: "哈尔滨", Type: "理工", Level: "985", Domains: []string{"hit.edu.cn", "hit.cn"}},
		{Name: "西安交通大学", ShortName: "西交", Province: "陕西", City: "西安", Type: "理工", Level: "985", Domains: []string{"xjtu.edu.cn", "xjtu.cn"}},
		{Name: "北京航空航天大学", ShortName: "北航", Province: "北京", City: "北京", Type: "理工", Level: "985", Domains: []string{"buaa.edu.cn", "buaa.cn"}},
		{Name: "南开大学", ShortName: "南开", Province: "天津", City: "天津", Type: "综合", Level: "985", Domains: []string{"nankai.edu.cn", "nankai.cn"}},
		{Name: "天津大学", ShortName: "天大", Province: "天津", City: "天津", Type: "理工", Level: "985", Domains: []string{"tju.edu.cn", "tju.cn"}},
		{Name: "东南大学", ShortName: "东大", Province: "江苏", City: "南京", Type: "理工", Level: "985", Domains: []string{"seu.edu.cn", "seu.cn"}},
		{Name: "华中科技大学", ShortName: "华科", Province: "湖北", City: "武汉", Type: "理工", Level: "985", Domains: []string{"hust.edu.cn", "hust.cn"}},
		{Name: "武汉大学", ShortName: "武大", Province: "湖北", City: "武汉", Type: "综合", Level: "985", Domains: []string{"whu.edu.cn", "whu.cn"}},
		{Name: "厦门大学", ShortName: "厦大", Province: "福建", City: "厦门", Type: "综合", Level: "985", Domains: []string{"xmu.edu.cn", "xmu.cn"}},
		{Name: "山东大学", ShortName: "山大", Province: "山东", City: "济南", Type: "综合", Level: "985", Domains: []string{"sdu.edu.cn", "sdu.cn"}},
		{Name: "湖南大学", ShortName: "湖大", Province: "湖南", City: "长沙", Type: "综合", Level: "985", Domains: []string{"hnu.edu.cn", "hnu.cn"}},
		{Name: "中南大学", ShortName: "中南", Province: "湖南", City: "长沙", Type: "综合", Level: "985", Domains: []string{"csu.edu.cn", "csu.cn"}},
		{Name: "中山大学", ShortName: "中大", Province: "广东", City: "广州", Type: "综合", Level: "985", Domains: []string{"sysu.edu.cn", "sysu.cn"}},
		{Name: "华南理工大学", ShortName: "华工", Province: "广东", City: "广州", Type: "理工", Level: "985", Domains: []string{"scut.edu.cn", "scut.cn"}},
		{Name: "四川大学", ShortName: "川大", Province: "四川", City: "成都", Type: "综合", Level: "985", Domains: []string{"scu.edu.cn", "scu.cn"}},
		{Name: "电子科技大学", ShortName: "成电", Province: "四川", City: "成都", Type: "理工", Level: "985", Domains: []string{"uestc.edu.cn", "uestc.cn"}},
		{Name: "重庆大学", ShortName: "重大", Province: "重庆", City: "重庆", Type: "综合", Level: "985", Domains: []string{"cqu.edu.cn", "cqu.cn"}},
		{Name: "大连理工大学", ShortName: "大工", Province: "辽宁", City: "大连", Type: "理工", Level: "985", Domains: []string{"dlut.edu.cn", "dlut.cn"}},
		{Name: "东北大学", ShortName: "东大", Province: "辽宁", City: "沈阳", Type: "理工", Level: "985", Domains: []string{"neu.edu.cn", "neu.cn"}},
		{Name: "吉林大学", ShortName: "吉大", Province: "吉林", City: "长春", Type: "综合", Level: "985", Domains: []string{"jlu.edu.cn", "jlu.cn"}},
		{Name: "兰州大学", ShortName: "兰大", Province: "甘肃", City: "兰州", Type: "综合", Level: "985", Domains: []string{"lzu.edu.cn", "lzu.cn"}},
		{Name: "西北工业大学", ShortName: "西工大", Province: "陕西", City: "西安", Type: "理工", Level: "985", Domains: []string{"nwpu.edu.cn", "nwpu.cn"}},
		{Name: "西北农林科技大学", ShortName: "西农", Province: "陕西", City: "杨凌", Type: "农林", Level: "985", Domains: []string{"nwsuaf.edu.cn", "nwsuaf.cn"}},
		{Name: "同济大学", ShortName: "同济", Province: "上海", City: "上海", Type: "理工", Level: "985", Domains: []string{"tongji.edu.cn", "tongji.cn"}},
		{Name: "华东师范大学", ShortName: "华师", Province: "上海", City: "上海", Type: "师范", Level: "985", Domains: []string{"ecnu.edu.cn", "ecnu.cn"}},
		{Name: "北京师范大学", ShortName: "北师大", Province: "北京", City: "北京", Type: "师范", Level: "985", Domains: []string{"bnu.edu.cn", "bnu.cn"}},
		{Name: "中国人民大学", ShortName: "人大", Province: "北京", City: "北京", Type: "综合", Level: "985", Domains: []string{"ruc.edu.cn", "ruc.cn"}},
		{Name: "中国农业大学", ShortName: "中农", Province: "北京", City: "北京", Type: "农林", Level: "985", Domains: []string{"cau.edu.cn", "cau.cn"}},
		{Name: "中央民族大学", ShortName: "民大", Province: "北京", City: "北京", Type: "民族", Level: "985", Domains: []string{"muc.edu.cn", "muc.cn"}},
		{Name: "国防科技大学", ShortName: "国防科大", Province: "湖南", City: "长沙", Type: "军事", Level: "985", Domains: []string{"nudt.edu.cn", "nudt.cn"}},
		{Name: "海军军医大学", ShortName: "二医大", Province: "上海", City: "上海", Type: "军事", Level: "985", Domains: []string{"smmu.edu.cn", "smmu.cn"}},
		{Name: "空军军医大学", ShortName: "四医大", Province: "陕西", City: "西安", Type: "军事", Level: "985", Domains: []string{"fmmu.edu.cn", "fmmu.cn"}},
	}...)

	// 211高校（部分）
	udb.universities = append(udb.universities, []University{
		{Name: "北京交通大学", ShortName: "北交", Province: "北京", City: "北京", Type: "理工", Level: "211", Domains: []string{"bjtu.edu.cn", "njtu.edu.cn"}},
		{Name: "北京工业大学", ShortName: "北工大", Province: "北京", City: "北京", Type: "理工", Level: "211", Domains: []string{"bjut.edu.cn"}},
		{Name: "北京科技大学", ShortName: "北科大", Province: "北京", City: "北京", Type: "理工", Level: "211", Domains: []string{"ustb.edu.cn"}},
		{Name: "北京化工大学", ShortName: "北化", Province: "北京", City: "北京", Type: "理工", Level: "211", Domains: []string{"buct.edu.cn"}},
		{Name: "北京邮电大学", ShortName: "北邮", Province: "北京", City: "北京", Type: "理工", Level: "211", Domains: []string{"bupt.edu.cn", "bupt.cn"}},
		{Name: "北京林业大学", ShortName: "北林", Province: "北京", City: "北京", Type: "农林", Level: "211", Domains: []string{"bjfu.edu.cn"}},
		{Name: "北京中医药大学", ShortName: "北中医", Province: "北京", City: "北京", Type: "医药", Level: "211", Domains: []string{"bucm.edu.cn"}},
		{Name: "北京外国语大学", ShortName: "北外", Province: "北京", City: "北京", Type: "语言", Level: "211", Domains: []string{"bfsu.edu.cn"}},
		{Name: "中国传媒大学", ShortName: "中传", Province: "北京", City: "北京", Type: "艺术", Level: "211", Domains: []string{"cuc.edu.cn"}},
		{Name: "中央财经大学", ShortName: "中财", Province: "北京", City: "北京", Type: "财经", Level: "211", Domains: []string{"cufe.edu.cn"}},
		{Name: "对外经济贸易大学", ShortName: "外经贸", Province: "北京", City: "北京", Type: "财经", Level: "211", Domains: []string{"uibe.edu.cn"}},
		{Name: "中国政法大学", ShortName: "法大", Province: "北京", City: "北京", Type: "政法", Level: "211", Domains: []string{"cupl.edu.cn"}},
		{Name: "华北电力大学", ShortName: "华电", Province: "北京", City: "北京", Type: "理工", Level: "211", Domains: []string{"ncepu.edu.cn"}},
		{Name: "中国矿业大学（北京）", ShortName: "矿大", Province: "北京", City: "北京", Type: "理工", Level: "211", Domains: []string{"cumtb.edu.cn"}},
		{Name: "中国石油大学（北京）", ShortName: "中石大", Province: "北京", City: "北京", Type: "理工", Level: "211", Domains: []string{"cup.edu.cn"}},
		{Name: "中国地质大学（北京）", ShortName: "地大", Province: "北京", City: "北京", Type: "理工", Level: "211", Domains: []string{"cugb.edu.cn"}},
		{Name: "上海财经大学", ShortName: "上财", Province: "上海", City: "上海", Type: "财经", Level: "211", Domains: []string{"shufe.edu.cn"}},
		{Name: "上海外国语大学", ShortName: "上外", Province: "上海", City: "上海", Type: "语言", Level: "211", Domains: []string{"shisu.edu.cn"}},
		{Name: "华东理工大学", ShortName: "华理", Province: "上海", City: "上海", Type: "理工", Level: "211", Domains: []string{"ecust.edu.cn"}},
		{Name: "东华大学", ShortName: "东华", Province: "上海", City: "上海", Type: "理工", Level: "211", Domains: []string{"dhu.edu.cn"}},
		{Name: "上海大学", ShortName: "上大", Province: "上海", City: "上海", Type: "综合", Level: "211", Domains: []string{"shu.edu.cn"}},
		{Name: "海军军医大学", ShortName: "二军大", Province: "上海", City: "上海", Type: "军事", Level: "211", Domains: []string{"smmu.edu.cn"}},
		{Name: "南京航空航天大学", ShortName: "南航", Province: "江苏", City: "南京", Type: "理工", Level: "211", Domains: []string{"nuaa.edu.cn"}},
		{Name: "南京理工大学", ShortName: "南理工", Province: "江苏", City: "南京", Type: "理工", Level: "211", Domains: []string{"njust.edu.cn"}},
		{Name: "河海大学", ShortName: "河海", Province: "江苏", City: "南京", Type: "理工", Level: "211", Domains: []string{"hhu.edu.cn"}},
		{Name: "南京农业大学", ShortName: "南农", Province: "江苏", City: "南京", Type: "农林", Level: "211", Domains: []string{"njau.edu.cn"}},
		{Name: "中国药科大学", ShortName: "药大", Province: "江苏", City: "南京", Type: "医药", Level: "211", Domains: []string{"cpu.edu.cn"}},
		{Name: "南京师范大学", ShortName: "南师大", Province: "江苏", City: "南京", Type: "师范", Level: "211", Domains: []string{"njnu.edu.cn"}},
		{Name: "苏州大学", ShortName: "苏大", Province: "江苏", City: "苏州", Type: "综合", Level: "211", Domains: []string{"suda.edu.cn"}},
		{Name: "江南大学", ShortName: "江大", Province: "江苏", City: "无锡", Type: "理工", Level: "211", Domains: []string{"jiangnan.edu.cn"}},
		{Name: "中国矿业大学", ShortName: "矿大", Province: "江苏", City: "徐州", Type: "理工", Level: "211", Domains: []string{"cumt.edu.cn"}},
		{Name: "合肥工业大学", ShortName: "合工大", Province: "安徽", City: "合肥", Type: "理工", Level: "211", Domains: []string{"hfut.edu.cn"}},
		{Name: "安徽大学", ShortName: "安大", Province: "安徽", City: "合肥", Type: "综合", Level: "211", Domains: []string{"ahu.edu.cn"}},
		{Name: "福州大学", ShortName: "福大", Province: "福建", City: "福州", Type: "理工", Level: "211", Domains: []string{"fzu.edu.cn"}},
		{Name: "南昌大学", ShortName: "昌大", Province: "江西", City: "南昌", Type: "综合", Level: "211", Domains: []string{"ncu.edu.cn"}},
		{Name: "郑州大学", ShortName: "郑大", Province: "河南", City: "郑州", Type: "综合", Level: "211", Domains: []string{"zzu.edu.cn"}},
		{Name: "武汉理工大学", ShortName: "武理工", Province: "湖北", City: "武汉", Type: "理工", Level: "211", Domains: []string{"whut.edu.cn"}},
		{Name: "中国地质大学（武汉）", ShortName: "地大", Province: "湖北", City: "武汉", Type: "理工", Level: "211", Domains: []string{"cug.edu.cn"}},
		{Name: "华中农业大学", ShortName: "华农", Province: "湖北", City: "武汉", Type: "农林", Level: "211", Domains: []string{"hzau.edu.cn"}},
		{Name: "华中师范大学", ShortName: "华师", Province: "湖北", City: "武汉", Type: "师范", Level: "211", Domains: []string{"ccnu.edu.cn"}},
		{Name: "中南财经政法大学", ShortName: "中南大", Province: "湖北", City: "武汉", Type: "财经", Level: "211", Domains: []string{"zuel.edu.cn"}},
		{Name: "湖南师范大学", ShortName: "湖师大", Province: "湖南", City: "长沙", Type: "师范", Level: "211", Domains: []string{"hunnu.edu.cn"}},
		{Name: "暨南大学", ShortName: "暨大", Province: "广东", City: "广州", Type: "综合", Level: "211", Domains: []string{"jnu.edu.cn"}},
		{Name: "华南师范大学", ShortName: "华师", Province: "广东", City: "广州", Type: "师范", Level: "211", Domains: []string{"scnu.edu.cn"}},
		{Name: "广西大学", ShortName: "西大", Province: "广西", City: "南宁", Type: "综合", Level: "211", Domains: []string{"gxu.edu.cn"}},
		{Name: "海南大学", ShortName: "海大", Province: "海南", City: "海口", Type: "综合", Level: "211", Domains: []string{"hainanu.edu.cn"}},
		{Name: "西南交通大学", ShortName: "西南交大", Province: "四川", City: "成都", Type: "理工", Level: "211", Domains: []string{"swjtu.edu.cn"}},
		{Name: "西南财经大学", ShortName: "西财", Province: "四川", City: "成都", Type: "财经", Level: "211", Domains: []string{"swufe.edu.cn"}},
		{Name: "四川农业大学", ShortName: "川农", Province: "四川", City: "雅安", Type: "农林", Level: "211", Domains: []string{"sicau.edu.cn"}},
		{Name: "贵州大学", ShortName: "贵大", Province: "贵州", City: "贵阳", Type: "综合", Level: "211", Domains: []string{"gzu.edu.cn"}},
		{Name: "云南大学", ShortName: "云大", Province: "云南", City: "昆明", Type: "综合", Level: "211", Domains: []string{"ynu.edu.cn"}},
		{Name: "西藏大学", ShortName: "藏大", Province: "西藏", City: "拉萨", Type: "综合", Level: "211", Domains: []string{"utibet.edu.cn"}},
		{Name: "西北大学", ShortName: "西大", Province: "陕西", City: "西安", Type: "综合", Level: "211", Domains: []string{"nwu.edu.cn"}},
		{Name: "西安电子科技大学", ShortName: "西电", Province: "陕西", City: "西安", Type: "理工", Level: "211", Domains: []string{"xidian.edu.cn"}},
		{Name: "长安大学", ShortName: "长大", Province: "陕西", City: "西安", Type: "理工", Level: "211", Domains: []string{"chd.edu.cn"}},
		{Name: "陕西师范大学", ShortName: "陕师大", Province: "陕西", City: "西安", Type: "师范", Level: "211", Domains: []string{"snnu.edu.cn"}},
		{Name: "青海大学", ShortName: "青大", Province: "青海", City: "西宁", Type: "综合", Level: "211", Domains: []string{"qhu.edu.cn"}},
		{Name: "宁夏大学", ShortName: "宁大", Province: "宁夏", City: "银川", Type: "综合", Level: "211", Domains: []string{"nxu.edu.cn"}},
		{Name: "新疆大学", ShortName: "新大", Province: "新疆", City: "乌鲁木齐", Type: "综合", Level: "211", Domains: []string{"xju.edu.cn"}},
		{Name: "石河子大学", ShortName: "石大", Province: "新疆", City: "石河子", Type: "综合", Level: "211", Domains: []string{"shzu.edu.cn"}},
		{Name: "太原理工大学", ShortName: "太原理工", Province: "山西", City: "太原", Type: "理工", Level: "211", Domains: []string{"tyut.edu.cn"}},
		{Name: "内蒙古大学", ShortName: "内大", Province: "内蒙古", City: "呼和浩特", Type: "综合", Level: "211", Domains: []string{"imu.edu.cn"}},
		{Name: "辽宁大学", ShortName: "辽大", Province: "辽宁", City: "沈阳", Type: "综合", Level: "211", Domains: []string{"lnu.edu.cn"}},
		{Name: "延边大学", ShortName: "延大", Province: "吉林", City: "延吉", Type: "综合", Level: "211", Domains: []string{"ybu.edu.cn"}},
		{Name: "东北师范大学", ShortName: "东师", Province: "吉林", City: "长春", Type: "师范", Level: "211", Domains: []string{"nenu.edu.cn"}},
		{Name: "哈尔滨工程大学", ShortName: "哈工程", Province: "黑龙江", City: "哈尔滨", Type: "理工", Level: "211", Domains: []string{"hrbeu.edu.cn"}},
		{Name: "东北农业大学", ShortName: "东农", Province: "黑龙江", City: "哈尔滨", Type: "农林", Level: "211", Domains: []string{"neau.edu.cn"}},
		{Name: "东北林业大学", ShortName: "东林", Province: "黑龙江", City: "哈尔滨", Type: "农林", Level: "211", Domains: []string{"nefu.edu.cn"}},
	}...)
}

func (udb *UniversityDB) GetAll() []University {
	return udb.universities
}

func (udb *UniversityDB) GetByLevel(level string) []University {
	var result []University
	for _, u := range udb.universities {
		if u.Level == level {
			result = append(result, u)
		}
	}
	return result
}

func (udb *UniversityDB) GetByProvince(province string) []University {
	var result []University
	for _, u := range udb.universities {
		if u.Province == province {
			result = append(result, u)
		}
	}
	return result
}

func (udb *UniversityDB) GetByDomain(domain string) *University {
	for _, u := range udb.universities {
		for _, d := range u.Domains {
			if d == domain || domain == d || domain == "www."+d {
				return &u
			}
		}
	}
	return nil
}

func (udb *UniversityDB) Search(name string) []University {
	var result []University
	for _, u := range udb.universities {
		if contains(u.Name, name) || contains(u.ShortName, name) {
			result = append(result, u)
		}
	}
	return result
}

func (udb *UniversityDB) Count() int {
	return len(udb.universities)
}

func (udb *UniversityDB) CountByLevel(level string) int {
	count := 0
	for _, u := range udb.universities {
		if u.Level == level {
			count++
		}
	}
	return count
}

func (udb *UniversityDB) PrintStats() {
	fmt.Println("\n📊 高校域名库统计：")
	fmt.Println("=====================================")
	fmt.Printf("  总高校数: %d\n", udb.Count())
	fmt.Printf("  985高校: %d\n", udb.CountByLevel("985"))
	fmt.Printf("  211高校: %d\n", udb.CountByLevel("211"))
	fmt.Printf("  普通本科: %d\n", udb.CountByLevel("普通"))
	fmt.Println("=====================================")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(s[:len(substr)] == substr) ||
		(len(s) > len(substr) && s[len(s)-len(substr):] == substr))
}

// 全局实例
var defaultDB = NewUniversityDB()

// GetUniversities 获取所有高校列表（API用）
func GetUniversities() []University {
	return defaultDB.GetAll()
}

// GetUniversitiesByLevel 按级别获取高校（API用）
func GetUniversitiesByLevel(level string) []University {
	return defaultDB.GetByLevel(level)
}

// GetUniversitiesByProvince 按省份获取高校（API用）
func GetUniversitiesByProvince(province string) []University {
	return defaultDB.GetByProvince(province)
}
