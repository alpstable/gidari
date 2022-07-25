package morningstar

import "path"

// * This is a generated file, do not edit

type rawPath uint8

const (
	_ rawPath = iota
	HouseholdsPath
)

// Households will get all household IDs for an advisor. The household ID can be used with the Portfolios API to return
// information about household accounts.
func getHouseholdsPath(params map[string]string) string {
	return path.Join("/households")
}

// Get takes an rawPath const and rawPath arguments to parse the URL rawPath path.
func (p rawPath) Path(params map[string]string) string {
	return map[rawPath]func(map[string]string) string{
		HouseholdsPath: getHouseholdsPath,
	}[p](params)
}

func (p rawPath) Scope() string {
	return map[rawPath]string{}[p]
}
