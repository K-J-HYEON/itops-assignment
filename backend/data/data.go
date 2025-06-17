package data

import (
	"fmt"
	"time"
	"itops-assignment/backend/models" // models 패키지 임포트
)

// IssueStore 인터페이스는 이슈 데이터 저장소에 대한 계약을 정의합니다.
type IssueStore interface {
	CreateIssue(issue models.Issue) (models.Issue, error)
	GetIssues(statusFilter string) ([]models.Issue, error)
	GetIssueByID(id uint) (models.Issue, error)
	UpdateIssue(issue models.Issue) (models.Issue, error)
	GetNextIssueID() uint
	GetUserByID(id uint) (*models.User, bool) // 사용자 조회 메서드 반환 타입을 *models.User로 변경
}

// MockStore는 IssueStore 인터페이스의 메모리 내 구현입니다.
type MockStore struct {
	issues      map[uint]models.Issue
	nextIssueID uint
	Users       map[uint]*models.User // 미리 정의된 사용자 맵을 *models.User 타입으로 변경
}

// NewMockStore는 MockStore의 새 인스턴스를 초기화하고 반환합니다.
// 여기에 초기 이슈 데이터를 추가합니다.
func NewMockStore() *MockStore {
	// users 맵에 User 구조체의 포인터를 저장하도록 수정
	users := map[uint]*models.User{
		1: {ID: 1, Name: "김개발"},
		2: {ID: 2, Name: "이디자인"},
		3: {ID: 3, Name: "박기획"},
	}

	store := &MockStore{
		issues:      make(map[uint]models.Issue),
		nextIssueID: 1,
		Users:       users,
	}

	// 초기 이슈 데이터 추가
	now := time.Now()
	issuesToSeed := []models.Issue{
		{
			Title:       "로그인 기능 버그 수정",
			Description: "사용자 로그인 시 간헐적으로 발생하는 오류 해결",
			Status:      "IN_PROGRESS",
			User:        users[1], // users[1]은 이미 포인터이므로 & 연산자 불필요
			CreatedAt:   now.Add(-time.Hour * 24 * 5),
			UpdatedAt:   now.Add(-time.Hour * 24 * 5),
		},
		{
			Title:       "메인 페이지 UI 개선",
			Description: "메인 페이지의 레이아웃을 더 직관적으로 개선하고 사용자 경험 향상",
			Status:      "PENDING",
			User:        nil, // 미할당
			CreatedAt:   now.Add(-time.Hour * 24 * 3),
			UpdatedAt:   now.Add(-time.Hour * 24 * 3),
		},
		{
			Title:       "데이터베이스 쿼리 최적화",
			Description: "대규모 데이터 조회 시 쿼리 성능 저하 문제 해결",
			Status:      "COMPLETED",
			User:        users[3], // users[3]은 이미 포인터이므로 & 연산자 불필요
			CreatedAt:   now.Add(-time.Hour * 24 * 7),
			UpdatedAt:   now.Add(-time.Hour * 24 * 2),
		},
		{
			Title:       "모바일 앱 푸시 알림 구현",
			Description: "모바일 앱에 푸시 알림 기능 추가 필요",
			Status:      "PENDING",
			User:        users[2], // users[2]은 이미 포인터이므로 & 연산자 불필요
			CreatedAt:   now.Add(-time.Hour * 24 * 1),
			UpdatedAt:   now.Add(-time.Hour * 24 * 1),
		},
		{
			Title:       "레거시 코드 리팩토링",
			Description: "오래된 코드 베이스를 현대적인 표준에 맞게 리팩토링",
			Status:      "IN_PROGRESS",
			User:        users[1], // users[1]은 이미 포인터이므로 & 연산자 불필요
			CreatedAt:   now.Add(-time.Hour * 24 * 10),
			UpdatedAt:   now.Add(-time.Hour * 24 * 8),
		},
	}

	for _, issue := range issuesToSeed {
		store.CreateIssue(issue) // 미리 정의된 이슈들을 스토어에 추가
	}

	return store
}

// CreateIssue는 MockStore에 새 이슈를 추가합니다.
// 초기 시딩 시 ID가 이미 부여된 경우를 처리하고, nextIssueID를 올바르게 관리합니다.
func (s *MockStore) CreateIssue(issue models.Issue) (models.Issue, error) {
	if issue.ID == 0 { // ID가 지정되지 않았다면 새로 할당
		issue.ID = s.nextIssueID
	}
	s.issues[issue.ID] = issue
	// 현재 할당된 ID가 nextIssueID보다 크거나 같으면 nextIssueID를 업데이트
	if issue.ID >= s.nextIssueID { 
		s.nextIssueID = issue.ID + 1
	}
	// CreatedAt, UpdatedAt은 NewIssueRequest 핸들러에서 또는 시딩 시 설정됨.
	return issue, nil
}


// GetIssues는 statusFilter에 따라 이슈 목록을 반환합니다.
func (s *MockStore) GetIssues(statusFilter string) ([]models.Issue, error) {
	var filteredIssues []models.Issue
	for _, issue := range s.issues {
		if statusFilter == "" || issue.Status == statusFilter {
			filteredIssues = append(filteredIssues, issue)
		}
	}
	return filteredIssues, nil
}

// GetIssueByID는 주어진 ID에 해당하는 이슈를 반환합니다.
func (s *MockStore) GetIssueByID(id uint) (models.Issue, error) {
	issue, ok := s.issues[id]
	if !ok {
		return models.Issue{}, fmt.Errorf("ID %d를 가진 이슈를 찾을 수 없습니다", id)
	}
	return issue, nil
}

// UpdateIssue는 MockStore에서 기존 이슈를 업데이트합니다.
func (s *MockStore) UpdateIssue(issue models.Issue) (models.Issue, error) {
	if _, ok := s.issues[issue.ID]; !ok {
		return models.Issue{}, fmt.Errorf("ID %d를 가진 이슈를 찾을 수 없습니다", issue.ID)
	}
	issue.UpdatedAt = time.Now()
	s.issues[issue.ID] = issue
	return issue, nil
}

// GetNextIssueID는 다음에 사용할 고유 이슈 ID를 반환합니다.
func (s *MockStore) GetNextIssueID() uint {
	return s.nextIssueID
}

// GetUserByID는 주어진 ID에 해당하는 사용자를 반환합니다.
func (s *MockStore) GetUserByID(id uint) (*models.User, bool) { // 반환 타입을 *models.User로 변경
	user, ok := s.Users[id]
	return user, ok
}
