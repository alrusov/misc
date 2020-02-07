package misc

import (
	"regexp"
)

//----------------------------------------------------------------------------------------------------------------------------//

// Replace --
type Replace []replaceDef
type replaceDef struct {
	exp       *regexp.Regexp
	replaceTo string
}

//----------------------------------------------------------------------------------------------------------------------------//

// NewReplace --
func NewReplace() *Replace {
	return &Replace{}
}

//----------------------------------------------------------------------------------------------------------------------------//

// Add --
func (r *Replace) Add(re string, replaceTo string) error {
	exp, err := regexp.Compile(re)
	if err != nil {
		return err
	}

	*r = append(*r,
		replaceDef{
			exp:       exp,
			replaceTo: replaceTo,
		},
	)
	return nil
}

// AddMulti --
func (r *Replace) AddMulti(list map[string]string) error {
	for re, replaceTo := range list {
		err := r.Add(re, replaceTo)
		if err != nil {
			return err
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------------//

// Do --
func (r *Replace) Do(s string) string {
	for _, rr := range *r {
		s = rr.exp.ReplaceAllString(s, rr.replaceTo)
	}
	return s
}

//----------------------------------------------------------------------------------------------------------------------------//
