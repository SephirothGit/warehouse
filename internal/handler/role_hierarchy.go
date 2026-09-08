package handler

var roleLevels = map[string]int{
	"warehouse_worker": 1,
	"manager":          2,
	"admin":            3,
}

func hasRequiredLevel(userRoles []string, requiredRole string) bool {
	requiredLevel, ok := roleLevels[requiredRole]
	if !ok {
		return false
	}

	for _, userRole := range userRoles {
		if level, ok := roleLevels[userRole]; ok && level >= requiredLevel {
			return true
		}
	}
	return false
}
