package kong

import (
	"fmt"
	"strconv"
)

var computedPluginProperties = []string{"created_at", "id", "consumer_id"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

// For the plugin config blobs, since there is no schema,
// we need to convert the values to a string to match config and avoid a diff.
func convertInterfaceToString(val interface{}) string {
	if s, ok := val.(string); ok {
		return s
	} else if b, ok := val.(bool); ok {
		return strconv.FormatBool(b)
	} else if f, ok := val.(float64); ok {
		return fmt.Sprintf("%f", f)
	} else if i, ok := val.(int); ok {
		return strconv.Itoa(i)
	}

	return ""
}
