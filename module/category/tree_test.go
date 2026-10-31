package category

import (
	"testing"

	"expense.com/m/entity"
)

func category(id, name string, parent *string, sortOrder int) *entity.CategoryDTO {
	return &entity.CategoryDTO{
		CommonDTOFields: entity.CommonDTOFields{ID: id},
		Name:            name,
		ParentId:        parent,
		SortOrder:       sortOrder,
	}
}

func ptr(s string) *string {
	return &s
}

func TestBuildTree(t *testing.T) {
	flat := []*entity.CategoryDTO{
		category("kopi", "Kopi", ptr("makan"), 1),
		category("makan", "Makan & Minum", nil, 1),
		category("transport", "Transportasi", nil, 0),
		category("bensin", "Bensin", ptr("transport"), 0),
		category("makanluar", "Makan di Luar", ptr("makan"), 0),
		category("orphan", "Yatim", ptr("tidak-ada"), 0),
	}

	roots := buildTree(flat)

	if len(roots) != 3 {
		t.Fatalf("jumlah root = %d, want 3", len(roots))
	}

	// Root diurutkan berdasarkan sortOrder lalu nama.
	if roots[0].ID != "transport" {
		t.Fatalf("root pertama harus transport, got %s", roots[0].ID)
	}
	if roots[2].ID != "makan" {
		t.Fatalf("root terakhir harus makan, got %s", roots[2].ID)
	}

	var makan *entity.CategoryDTO
	for _, root := range roots {
		if root.ID == "makan" {
			makan = root
		}
	}
	if makan == nil {
		t.Fatal("root makan tidak ditemukan")
	}
	if len(makan.Children) != 2 {
		t.Fatalf("jumlah child Makan = %d, want 2", len(makan.Children))
	}
	if makan.Children[0].ID != "makanluar" || makan.Children[1].ID != "kopi" {
		t.Fatalf("urutan child salah: %s, %s", makan.Children[0].ID, makan.Children[1].ID)
	}

	found := false
	for _, root := range roots {
		if root.ID == "orphan" {
			found = true
		}
	}
	if !found {
		t.Fatal("kategori tanpa induk valid harus menjadi root")
	}
}

func TestBuildTreeDepthTwo(t *testing.T) {
	flat := []*entity.CategoryDTO{
		category("a", "A", nil, 0),
		category("b", "B", ptr("a"), 0),
	}

	roots := buildTree(flat)
	if len(roots) != 1 || len(roots[0].Children) != 1 {
		t.Fatalf("tree tidak sesuai: %+v", roots)
	}
	if roots[0].Children[0].Children != nil {
		t.Fatal("anak tidak boleh memiliki anak pada depth 2")
	}
}
