package kong

var computedPluginProperties = []string{"created_at", "id", "consumer", "service", "route"}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
