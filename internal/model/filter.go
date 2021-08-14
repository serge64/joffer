package model

import (
	"fmt"
	"strings"
)

type Filter struct {
	UserID        int      `json:"-"`
	ExludeSites   []string `json:"excluded_sites"`
	ExcludeGroups []string `json:"excluded_groups"`
	Positions     []string `json:"positions"`
	Salary        string   `json:"salary"`
	Companies     []string `json:"companies"`
	Areas         []string `json:"areas"`
	OlderThan     string   `json:"older_than"`
}

func (f *Filter) ToString() string {
	var filter []string

	if len(f.ExludeSites) > 0 {
		filter = append(filter, fmt.Sprintf("pl.full_name NOT IN ('%s')", strings.Join(f.ExludeSites, "', '")))
	}

	if len(f.ExcludeGroups) > 0 {
		filter = append(filter, fmt.Sprintf("g.name NOT IN ('%s')", strings.Join(f.ExcludeGroups, "', '")))
	}

	if len(f.Positions) > 0 {
		filter = append(filter, fmt.Sprintf("LOWER(v.name) IN ('%s')", strings.Join(f.Positions, "', '")))
	}

	if len(f.Companies) > 0 {
		filter = append(filter, fmt.Sprintf("LOWER(company) IN ('%s')", strings.Join(f.Companies, "', '")))
	}

	if len(f.Areas) > 0 {
		filter = append(filter, fmt.Sprintf("LOWER(area) IN ('%s')", strings.Join(f.Areas, "', '")))
	}

	if len(f.Salary) > 0 {
		filter = append(filter, fmt.Sprintf("CASE salary_to WHEN 0 THEN salary_from >= %s ELSE salary_from BETWEEN %s AND salary_to END", f.Salary, f.Salary))
	}

	if len(f.OlderThan) > 0 {
		filter = append(filter, fmt.Sprintf("'%s'::timestamp > at_published", f.OlderThan))
	}

	filter = append(filter, fmt.Sprintf("r.user_id = %d", f.UserID))
	filter = append(filter, "responsed = false")

	return fmt.Sprintf("WHERE %s", strings.Join(filter, " AND "))
}
