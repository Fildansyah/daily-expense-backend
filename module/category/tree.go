package category

import (
	"sort"

	"expense.com/m/entity"
)

// buildTree mengubah daftar kategori datar menjadi hierarki dua level.
// Kategori yang induknya tidak ditemukan diperlakukan sebagai root.
func buildTree(flat []*entity.CategoryDTO) []*entity.CategoryDTO {
	nodes := make(map[string]*entity.CategoryDTO, len(flat))
	ordered := make([]*entity.CategoryDTO, 0, len(flat))

	for _, category := range flat {
		if category == nil {
			continue
		}
		node := *category
		node.Children = nil
		nodes[category.ID] = &node
		ordered = append(ordered, &node)
	}

	roots := make([]*entity.CategoryDTO, 0)
	for _, node := range ordered {
		if node.ParentId != nil {
			if parent, ok := nodes[*node.ParentId]; ok && parent.ID != node.ID {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}

	sortNodes(roots)
	for _, node := range ordered {
		sortNodes(node.Children)
	}

	return roots
}

func sortNodes(nodes []*entity.CategoryDTO) {
	sort.SliceStable(nodes, func(i, j int) bool {
		if nodes[i].SortOrder != nodes[j].SortOrder {
			return nodes[i].SortOrder < nodes[j].SortOrder
		}
		return nodes[i].Name < nodes[j].Name
	})
}
