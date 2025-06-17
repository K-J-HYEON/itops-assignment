package models

import "time"

// User 구조체는 사용자의 ID와 이름을 정의합니다.
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Issue 구조체는 이슈의 모든 속성을 정의합니다.
type Issue struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // PENDING, IN_PROGRESS, COMPLETED, CANCELLED
	User        *User     `json:"user,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// NewIssueRequest 구조체는 이슈 생성 요청 시 클라이언트로부터 받는 데이터를 정의합니다.
type NewIssueRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      *uint  `json:"userId"` // 담당자 ID, nil 가능 (선택 사항)
}

// UpdateIssueRequest 구조체는 이슈 수정 요청 시 클라이언트로부터 받는 데이터를 정의합니다.
type UpdateIssueRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	UserID      *uint   `json:"userId"` // 담당자 ID, nil 가능 (담당자 해제)
}

// APIError 구조체는 API 응답에서 에러 정보를 전달하는 표준 형식을 정의합니다.
type APIError struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
