package model

import "fmt"

type Pagination struct {
	PageIndex  int   `json:"page_index"`
	PageSize   int   `json:"page_size"`
	TotalCount int64 `json:"total_count"`
	PageCount  int   `json:"page_count"`
}

func (p *Pagination) GetPageCount() int {
	if p.PageSize == 0 {
		return 0
	}
	if int(p.TotalCount)%p.PageSize == 0 {
		p.PageCount = int(p.TotalCount) / p.PageSize
	} else {
		p.PageCount = int(p.TotalCount)/p.PageSize + 1
	}
	return p.PageCount
}

func (p *Pagination) GetPageOffset() int {
	if p.PageIndex <= 0 {
		return 1
	}
	return p.PageSize * (p.PageIndex - 1)
}

func (p *Pagination) String() string {
	return fmt.Sprintf("Page size: %d, total count: %d, page index: %d, page count: %d, page offset: %d\n", p.PageSize, p.TotalCount, p.PageIndex, p.GetPageCount(), p.GetPageOffset())
}
