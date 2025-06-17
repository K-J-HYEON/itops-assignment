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
	GetUserByID(id uint) (models.User, bool) // 사용자 조회 메서드 추가
}

// MockStore는 IssueStore 인터페이스의 메모리 내 구현입니다.
type MockStore struct {
	issues     map[uint]models.Issue
	nextIssueID uint
	Users      map[uint]models.User // 미리 정의된 사용자 맵
}

// NewMockStore는 MockStore의 새 인스턴스를 초기화하고 반환합니다.
func NewMockStore() *MockStore {
	users := map[uint]models.User{
		1: {ID: 1, Name: "김개발"},
		2: {ID: 2, Name: "이디자인"},
		3: {ID: 3, Name: "박기획"},
	}
	return &MockStore{
		issues:      make(map[uint]models.Issue),
		nextIssueID: 1,
		Users:       users,
	}
}

// CreateIssue는 MockStore에 새 이슈를 추가합니다.
func (s *MockStore) CreateIssue(issue models.Issue) (models.Issue, error) {
	issue.ID = s.GetNextIssueID()
	issue.CreatedAt = time.Now()
	issue.UpdatedAt = time.Now()
	s.issues[issue.ID] = issue
	s.nextIssueID++
	return issue, nil
}

// GetIssues는 statusFilter에 따라 이슈 목록을 반환합니다.
func (s *MockStore) GetIssues(statusFilter string) ([]models.Issue, error) {
	var filteredIssues []models.Issue
	// 맵은 순서가 없으므로, 정렬이 필요하면 여기서 수행하거나 클라이언트에서 처리해야 합니다.
	// 간단한 예제이므로 순서 없는 상태로 반환합니다.
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
func (s *MockStore) GetUserByID(id uint) (models.User, bool) {
	user, ok := s.Users[id]
	return user, ok
}
