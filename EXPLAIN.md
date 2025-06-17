## 추가 고려사항

- API 테스트 코드
- 백엔드 서버 구조(아키텍쳐) 기술
- 프론트엔드 컴포넌트 구조 기술
- 쿼리 최적화

### API 테스트 코드

현재 백엔드 Go 애플리케이션은 인메모리 MockStore를 사용하여 데이터를 관리하고 있습니다.

data 패키지의 IssueStore 인터페이스 구현체(MockStore) 메서드들이 예상대로 동작하는지 확인하는 테스트를 작성합니다. (예: CreateIssue, GetIssueByID, UpdateIssue 등)

handler 패키지의 각 핸들러 함수(CreateIssueHandler, GetIssuesHandler, UpdateIssueHandler 등)가 유효한 요청과 유효하지 않은 요청에 대해 올바른 HTTP 응답(상태 코드, JSON 바디)을 반환하는지 테스트합니다. Go의 net/http/httptest 패키지를 활용하여 실제 HTTP 요청-응답 시나리오를 시뮬레이션할 수 있습니다.

handlers/handlers_test.go

```
// package handler_test (또는 package handler)
// import (
//     "net/http"
//     "net/http/httptest"
//     "strings"
//     "testing"
//     "itops-assignment/backend/data"
//     "itops-assignment/backend/models"
//     "itops-assignment/backend/handler"
// )

// func TestCreateIssueHandler(t *testing.T) {
//     store := data.NewMockStore()
//     // 테스트 데이터 초기화 (필요하다면)
//     store.CreateIssue(models.Issue{Title: "Existing Issue", Description: "Desc", Status: "PENDING", User: nil})

//     // 유효한 요청
//     payload := `{"title": "새로운 이슈", "description": "이슈 설명입니다", "userId": 1}`
//     req := httptest.NewRequest("POST", "/issue", strings.NewReader(payload))
//     req.Header.Set("Content-Type", "application/json")
//     rr := httptest.NewRecorder()

//     handler := handler.CreateIssueHandler(store)
//     handler.ServeHTTP(rr, req)

//     if status := rr.Code; status != http.StatusCreated {
//         t.Errorf("핸들러가 %v 상태 코드를 반환했습니다. 예상: %v", status, http.StatusCreated)
//     }

//     // 응답 바디 검증 (옵션)
//     // var issue models.Issue
//     // json.NewDecoder(rr.Body).Decode(&issue)
//     // if issue.Title != "새로운 이슈" {
//     //     t.Errorf("잘못된 이슈 제목: %v", issue.Title)
//     // }
// }

```

### 백엔드 서버 구조 (아키텍처) 기술

현재 백엔드 서버는 간단한 3계층 아키텍처를 따르고 있습니다.

프레젠테이션 계층 (Presentation Layer):

main.go: HTTP 서버를 시작하고, 라우터를 설정하며, CORS 미들웨어를 적용합니다.

handlers 패키지: HTTP 요청을 처리하고 응답을 생성하는 핸들러 함수들을 포함합니다. 요청 파싱, 비즈니스 로직 호출, 응답 포맷팅을 담당합니다.

비즈니스/도메인 계층 (Business/Domain Layer):

(현재는 핸들러 내에 비즈니스 로직이 일부 포함되어 있지만) 별도의 서비스 또는 유즈케이스 패키지를 도입하여 핸들러에서 비즈니스 규칙을 분리할 수 있습니다. 이는 코드의 재사용성과 테스트 용이성을 높입니다.

(예: service 패키지 내 issue_service.go에서 이슈 생성/수정/조회 비즈니스 로직 처리)

데이터 접근 계층 (Data Access Layer):

data 패키지: 데이터 저장소(IssueStore 인터페이스)와 그 구현체(MockStore)를 정의합니다. 데이터의 CRUD(Create, Read, Update, Delete) 작업을 캡슐화합니다. 실제 데이터베이스(RDBMS, NoSQL 등)를 사용할 경우 이 계층에서 데이터베이스 드라이버와 상호작용합니다.

models 패키지: 데이터 구조(이슈, 사용자)와 API 요청/응답에 사용되는 모델들을 정의합니다.

추가 고려사항:

로깅 (Logging): 현재 log.Printf를 사용하고 있지만, 더 정교한 로깅 시스템(예: Logrus, Zap)을 도입하여 로그 레벨 관리, 구조화된 로그 출력 등을 구현할 수 있습니다.

에러 처리 (Error Handling): 현재 handleError 함수를 통해 일관된 에러 응답을 제공하고 있으나, Go의 에러 래핑(fmt.Errorf %w)과 사용자 정의 에러 타입을 활용하여 더 체계적인 에러 전파 및 처리가 가능합니다.

환경 설정 (Configuration): 현재 포트 번호 등이 코드에 하드코딩되어 있지만, 환경 변수(os.Getenv)나 설정 파일(YAML, JSON)을 통해 외부에서 관리하도록 개선할 수 있습니다.

데이터베이스 연동: 실제 프로덕션 환경에서는 MySQL, PostgreSQL, MongoDB 등 영구적인 데이터 저장소를 연결해야 합니다. MockStore 인터페이스를 구현하는 새로운 데이터베이스 어댑터를 작성하여 쉽게 교체할 수 있도록 설계되었습니다.

### 프론트엔드 컴포넌트 구조 기술

현재 프론트엔드 Vue.js 애플리케이션은 다음과 같은 컴포넌트 구조를 가집니다.

App.vue (루트 컴포넌트):

애플리케이션의 최상위 구조를 정의합니다.

공통 헤더(네비게이션 포함)를 포함하고, router-view를 통해 현재 라우트에 해당하는 컴포넌트를 렌더링합니다.

전역 CSS(index.css) 및 Font Awesome을 임포트합니다.

pages/ 디렉토리:

각 라우트(페이지)에 매핑되는 주요 뷰 컴포넌트들을 포함합니다.

IssueList.vue: 이슈 목록을 표시하고, 필터링 및 새 이슈 생성/상세 페이지로의 이동을 담당합니다. 백엔드 API와 직접 통신하여 데이터를 가져옵니다.

IssueForm.vue: 이슈 생성 및 상세/수정을 위한 폼을 제공합니다. 백엔드 API와 통신하여 이슈를 생성하거나 업데이트합니다.

router/index.js:

Vue Router의 인스턴스를 생성하고, URL 경로와 컴포넌트 간의 매핑(routes)을 정의합니다.

data/mockData.js:

개발 단계에서 백엔드 없이 프론트엔드를 독립적으로 테스트하거나, 담당자 드롭다운과 같이 정적인 데이터를 제공하는 데 사용되는 목(Mock) 데이터를 포함합니다.

추가 고려사항:

상태 관리 (State Management):

현재는 각 컴포넌트 내의 ref나 data를 통해 상태를 관리하지만, 애플리케이션 규모가 커지면 Vuex 또는 Pinia와 같은 중앙 집중식 상태 관리 라이브러리를 도입하여 데이터 흐름을 명확하게 하고 복잡한 상태를 효율적으로 관리할 수 있습니다.

재사용 가능한 컴포넌트 분리:

components/ 디렉토리를 생성하여 버튼, 입력 필드, 모달, 로딩 스피너 등 여러 페이지에서 재사용될 수 있는 UI 요소들을 컴포넌트로 분리하여 관리할 수 있습니다. 이는 코드 중복을 줄이고 유지보수성을 높입니다.

유효성 검증 라이브러리:

폼 유효성 검증 로직이 현재 컴포넌트 내에 직접 구현되어 있지만, VeeValidate나 FormKit과 같은 라이브러리를 사용하여 더 강력하고 유연한 폼 유효성 검증을 구현할 수 있습니다.

테마/스타일 관리:

Tailwind CSS를 사용하여 스타일링 일관성을 유지하고 있지만, 디자인 시스템이 복잡해질 경우 Storybook과 같은 도구를 사용하여 컴포넌트 라이브러리를 구축하고 디자인 토큰을 관리할 수 있습니다.

### 쿼리 최적화 (SQL)

현재 제공된 SQL 쿼리는 주어진 요구사항에 맞춰 기능적으로 올바르게 동작합니다. 하지만 실제 운영 환경에서는 데이터 양과 쿼리 빈도를 고려한 성능 최적화가 중요합니다.

고려사항:

인덱스 (Indexes) 생성:

FMS_HBL_MST 테이블:

HBL_NO: PK 역할을 하므로 이미 인덱스가 있을 가능성이 높습니다.

ETD: WHERE 절의 필터링 조건(BETWEEN) 및 ORDER BY 조건으로 자주 사용되므로, ONBD_YMD와 함께 인덱스를 고려할 수 있습니다.

POL_CD: GROUP BY 및 WHERE 절의 필터링 조건으로 사용되므로 인덱스를 고려할 수 있습니다.

FMS_HBL_CNTR 테이블:

HBL_NO: FMS_HBL_MST와의 조인 조건(ON)으로 사용되므로 인덱스가 중요합니다. FK이므로 인덱스를 고려할 수 있습니다.

CNTR_NO: COUNT 대상이므로 유니크 인덱스가 있다면 카운트에 더 효율적일 수 있습니다.

CNTR_TYPE: GROUP BY 및 WHERE 절의 필터링 조건으로 사용되므로 인덱스를 고려할 수 있습니다.

CNTR_WGT: SUM 대상이지만, WHERE 절에 포함된다면 인덱스 고려 가능.

예시 (Oracle 인덱스 생성):

```
CREATE INDEX IX_FMS_HBL_MST_ETD_POL ON FMS_HBL_MST (ETD, POL_CD);
CREATE INDEX IX_FMS_HBL_CNTR_HBL_NO_TYPE ON FMS_HBL_CNTR (HBL_NO, CNTR_TYPE);
```

### 실행 계획 (Execution Plan) 분석:

쿼리를 실행하기 전에 데이터베이스의 실행 계획을 분석하여 쿼리가 어떤 방식으로 데이터를 읽고 처리하는지 이해합니다. EXPLAIN PLAN FOR 명령을 사용하여 쿼리 성능 병목 현상을 식별할 수 있습니다.

```
EXPLAIN PLAN FOR
SELECT ... FROM ... WHERE ...;
SELECT PLAN_TABLE_OUTPUT FROM TABLE(DBMS_XPLAN.DISPLAY());
```

### 파티셔닝 (Partitioning):

데이터 양이 매우 커질 경우, FMS_HBL_MST 테이블을 ETD 또는 ONBD_YMD와 같은 날짜 컬럼을 기준으로 파티셔닝하여 특정 기간의 데이터를 더 빠르게 조회할 수 있습니다.

### 물리적 설계 최적화:

데이터 타입의 적절성, 정규화/비정규화 수준 조정 등 테이블의 물리적 설계 자체를 최적화하여 쿼리 성능에 영향을 줄 수 있습니다.

### 쿼리 힌트 (Hints):

특정 상황에서 데이터베이스 옵티마이저의 기본 동작을 오버라이드하여 더 나은 실행 계획을 유도할 수 있습니다. (초보자에게는 권장되지 않으며, 데이터베이스 버전에 따라 다를 수 있습니다.)
