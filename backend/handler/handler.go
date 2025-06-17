package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"itops-assignment/backend/data"
	"itops-assignment/backend/models"
)

// handleError는 주어진 에러 메시지와 HTTP 상태 코드로 JSON 에러 응답을 보냅니다.
func handleError(w http.ResponseWriter, message string, statusCode int) {
	log.Printf("Error: %s (Status: %d)", message, statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.APIError{Error: message, Code: statusCode})
}

// isValidStatus는 주어진 문자열이 유효한 이슈 상태인지 확인합니다.
func isValidStatus(status string) bool {
	switch status {
	case "PENDING", "IN_PROGRESS", "COMPLETED", "CANCELLED":
		return true
	default:
		return false
	}
}

// CreateIssueHandler는 POST /issue 요청을 처리합니다.
func CreateIssueHandler(store data.IssueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.NewIssueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleError(w, "유효하지 않은 요청 바디", http.StatusBadRequest)
			return
		}

		if req.Title == "" || req.Description == "" {
			handleError(w, "제목과 설명은 필수 입력 사항입니다", http.StatusBadRequest)
			return
		}

		var assignedUser *models.User
		initialStatus := "PENDING"

		if req.UserID != nil && *req.UserID != 0 { // 프론트에서 0을 null 대신 보낼 경우 처리
			user, ok := store.GetUserByID(*req.UserID)
			if !ok {
				handleError(w, fmt.Sprintf("존재하지 않는 사용자 ID입니다: %d", *req.UserID), http.StatusBadRequest)
				return
			}
			assignedUser = &user
			initialStatus = "IN_PROGRESS"
		}

		issue := models.Issue{
			Title:       req.Title,
			Description: req.Description,
			Status:      initialStatus,
			User:        assignedUser,
		}

		createdIssue, err := store.CreateIssue(issue)
		if err != nil {
			handleError(w, "이슈 생성에 실패했습니다: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdIssue)
	}
}

// GetIssuesHandler는 GET /issues 요청을 처리합니다.
func GetIssuesHandler(store data.IssueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusFilter := r.URL.Query().Get("status")

		if statusFilter != "" && !isValidStatus(statusFilter) {
			handleError(w, fmt.Sprintf("유효하지 않은 상태 필터입니다: %s", statusFilter), http.StatusBadRequest)
			return
		}

		issues, err := store.GetIssues(statusFilter)
		if err != nil {
			handleError(w, "이슈 목록 조회에 실패했습니다: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string][]models.Issue{"issues": issues})
	}
}

// GetIssueByIDHandler는 GET /issue/:id 요청을 처리합니다.
func GetIssueByIDHandler(store data.IssueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			handleError(w, "유효하지 않은 이슈 ID 형식입니다", http.StatusBadRequest)
			return
		}

		issue, err := store.GetIssueByID(uint(id))
		if err != nil {
			handleError(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(issue)
	}
}

// UpdateIssueHandler는 PATCH /issue/:id 요청을 처리합니다.
func UpdateIssueHandler(store data.IssueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			handleError(w, "유효하지 않은 이슈 ID 형식입니다", http.StatusBadRequest)
			return
		}

		existingIssue, err := store.GetIssueByID(uint(id))
		if err != nil {
			handleError(w, err.Error(), http.StatusNotFound)
			return
		}

		if existingIssue.Status == "COMPLETED" || existingIssue.Status == "CANCELLED" {
			handleError(w, "완료되거나 취소된 이슈는 수정할 수 없습니다", http.StatusForbidden)
			return
		}

		var req models.UpdateIssueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			handleError(w, "유효하지 않은 요청 바디", http.StatusBadRequest)
			return
		}

		// 부분 업데이트 로직
		if req.Title != nil {
			existingIssue.Title = *req.Title
		}
		if req.Description != nil {
			existingIssue.Description = *req.Description
		}

		// 상태 변경 로직
		if req.Status != nil {
			newStatus := *req.Status
			if !isValidStatus(newStatus) {
				handleError(w, fmt.Sprintf("유효하지 않은 상태 값입니다: %s", newStatus), http.StatusBadRequest)
				return
			}
			// 백엔드 비즈니스 규칙: 담당자가 없는 상태에서 IN_PROGRESS 또는 COMPLETED로 변경 불가
			// 프론트엔드에서 이미 UI 제약을 걸지만, 백엔드에서도 유효성 검증 필요
			if existingIssue.User == nil && (newStatus == "IN_PROGRESS" || newStatus == "COMPLETED") {
				handleError(w, "담당자가 없는 이슈는 'IN_PROGRESS' 또는 'COMPLETED' 상태로 변경할 수 없습니다", http.StatusBadRequest)
				return
			}
			existingIssue.Status = newStatus
		}

		// 담당자 변경 로직
		if req.UserID != nil { // 요청에 userId 필드가 포함된 경우 (null 또는 실제 ID)
			if *req.UserID == 0 { // 프론트에서 0을 'null' 대신 보낼 경우 담당자 해제로 간주
				existingIssue.User = nil
				// 백엔드 비즈니스 규칙: 담당자 제거 시 상태는 PENDING으로 변경
				existingIssue.Status = "PENDING"
			} else {
				user, ok := store.GetUserByID(*req.UserID)
				if !ok {
					handleError(w, fmt.Sprintf("존재하지 않는 사용자 ID입니다: %d", *req.UserID), http.StatusBadRequest)
					return
				}
				existingIssue.User = &user

				// 백엔드 비즈니스 규칙: PENDING 상태에서 담당자 할당 시, 상태가 명시적으로 지정되지 않았다면 IN_PROGRESS로 변경
				if existingIssue.Status == "PENDING" && req.Status == nil {
					existingIssue.Status = "IN_PROGRESS"
				}
			}
		}

		updatedIssue, err := store.UpdateIssue(existingIssue)
		if err != nil {
			handleError(w, "이슈 업데이트에 실패했습니다: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedIssue)
	}
}
