package models

import (
	"strings"
)

func GetPagingOption(pageReq int, limitReq int, sortReq string) (ModelOption, int, int) {
	var optionsQuery ModelOption
	page := 1
	limit := 10
	if pageReq >= 0 {
		page = pageReq
	}

	if limitReq > 0 {
		limit = limitReq
	}

	optionsQuery = ModelOption{
		Limit: int64(limit),
		Skip:  (int64(page) - 1) * int64(limit),
	}
	if sortReq != "" {
		if strings.HasPrefix(sortReq, "-") {
			runes := []rune(sortReq)
			sortBy := string(runes[1:len(sortReq)])
			optionsQuery.SortBy = sortBy
			optionsQuery.SortDir = -1
		} else {
			optionsQuery.SortBy = sortReq
			optionsQuery.SortDir = 1
		}
	}
	return optionsQuery, page, limit
}
