package models

import "time"

// User 구조체는 사용자의 ID와 이름을 정의합니다.
// `json:"..."` 태그는 JSON 마샬링/언마샬링 시 필드 이름을 지정합니다.
type User struct {
	ID   uint   `json:"id"`   // 사용자의 고유 ID
	Name string `json:"name"` // 사용자의 이름
}

// Issue 구조체는 이슈의 모든 속성을 정의합니다.
// User 필드는 이슈의 담당자를 나타내며, `omitempty`는 nil일 경우 JSON에서 생략됩니다.
type Issue struct {
	ID          uint      `json:"id"`          // 이슈의 고유 ID
	Title       string    `json:"title"`       // 이슈의 제목
	Description string    `json:"description"` // 이슈의 상세 설명
	Status      string    `json:"status"`      // 이슈의 현재 상태 (PENDING, IN_PROGRESS, COMPLETED, CANCELLED)
	User        *User     `json:"user,omitempty"` // 이슈의 담당자 (User 구조체 포인터), 담당자가 없으면 JSON에서 생략
	CreatedAt   time.Time `json:"createdAt"`   // 이슈가 생성된 시간
	UpdatedAt   time.Time `json:"updated updatedAt"`   // 이슈가 마지막으로 업데이트된 시간
}

// NewIssueRequest 구조체는 이슈 생성 요청 시 클라이언트로부터 받는 데이터를 정의합니다.
type NewIssueRequest struct {
	Title       string `json:"title"`       // 이슈 제목 (필수)
	Description string `json:"description"` // 이슈 설명 (필수)
	UserID      *uint  `json:"userId"`      // 담당자 ID (선택 사항), nil 가능
}

// UpdateIssueRequest 구조체는 이슈 수정 요청 시 클라이언트로부터 받는 데이터를 정의합니다.
// 모든 필드가 포인터로 선언되어 있어 부분 업데이트(PATCH)를 처리할 수 있습니다.
type UpdateIssueRequest struct {
	Title       *string `json:"title"`       // 수정할 이슈 제목 (선택 사항)
	Description *string `json:"description"` // 수정할 이슈 설명 (선택 사항)
	Status      *string `json:"status"`      // 수정할 이슈 상태 (선택 사항)
	UserID      *uint   `json:"userId"`      // 담당자 ID (선택 사항), nil 가능 (담당자 해제)
}

// APIError 구조체는 API 응답에서 에러 정보를 전달하는 표준 형식을 정의합니다.
type APIError struct {
	Error string `json:"error"` // 에러 메시지
	Code  int    `json:"code"`  // HTTP 상태 코드 (예: 400, 404, 500)
}
