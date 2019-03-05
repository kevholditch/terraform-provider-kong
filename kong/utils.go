package kong

import "regexp"

var computedPluginProperties = []string{"created_at", "id", "consumer", "service", "route"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func getRegex(r *regexp.Regexp, err error) *regexp.Regexp {
	return r
}

func readIntArrayFromInterface(in interface{}) []int {
	if arr := in.([]interface{}); arr != nil {
		array := make([]int, len(arr))
		for i, x := range arr {
			item := x.(int)
			array[i] = item
		}

		return array
	}

	return []int{}
}
