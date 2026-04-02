package edu

import (
	"testing"
)

func TestNewUniversityDB(t *testing.T) {
	db := NewUniversityDB()
	if db == nil {
		t.Fatal("NewUniversityDB() returned nil")
	}

	universities := db.GetAll()
	if len(universities) == 0 {
		t.Fatal("University database is empty")
	}

	t.Logf("Loaded %d universities", len(universities))
}

func TestGetByLevel(t *testing.T) {
	db := NewUniversityDB()

	// 测试985高校
	985Unis := db.GetByLevel("985")
	if len(985Unis) == 0 {
		t.Error("No 985 universities found")
	}
	t.Logf("Found %d 985 universities", len(985Unis))

	// 测试211高校
	211Unis := db.GetByLevel("211")
	if len(211Unis) == 0 {
		t.Error("No 211 universities found")
	}
	t.Logf("Found %d 211 universities", len(211Unis))
}

func TestGetByProvince(t *testing.T) {
	db := NewUniversityDB()

	// 测试北京高校
	beijingUnis := db.GetByProvince("北京")
	if len(beijingUnis) == 0 {
		t.Error("No Beijing universities found")
	}
	t.Logf("Found %d Beijing universities", len(beijingUnis))
}

func TestSearch(t *testing.T) {
	db := NewUniversityDB()

	// 搜索清华
	results := db.Search("清华")
	if len(results) == 0 {
		t.Error("Search for '清华' returned no results")
	}

	found := false
	for _, u := range results {
		if u.Name == "清华大学" {
			found = true
			break
		}
	}

	if !found {
		t.Error("清华大学 not found in search results")
	}
}

func TestGetUniversities(t *testing.T) {
	universities := GetUniversities()
	if len(universities) == 0 {
		t.Fatal("GetUniversities() returned empty slice")
	}

	// 验证每个大学都有域名
	for _, u := range universities {
		if len(u.Domains) == 0 {
			t.Errorf("University %s has no domains", u.Name)
		}
	}
}

func BenchmarkGetByLevel(b *testing.B) {
	db := NewUniversityDB()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.GetByLevel("985")
	}
}

func BenchmarkSearch(b *testing.B) {
	db := NewUniversityDB()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db.Search("大学")
	}
}
