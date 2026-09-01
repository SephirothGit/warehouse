package handler

import (
	"net/http"
	"strconv"

	"github.com/SephirothGit/warehouse/internal/repository"
)

func RequireRole(role string, userRepo repository.UserRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value("user_id").(string)
			if !ok {
				http.Error(w, "user id is not found in context", http.StatusUnauthorized)
				return
			}

			userIDInt, err := strconv.Atoi(userID)
			if err != nil {
				http.Error(w, "invalid user id", http.StatusInternalServerError)
				return 
			}

			userRoles, err := userRepo.GetRoles(userIDInt)
			if err != nil {
				http.Error(w, "unable to fetch user roles", http.StatusInternalServerError)
				return
			}

			hasRole := false
			for _, roleName := range userRoles {
				if roleName == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}