package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"itops-assignment/backend/data"
	"itops-assignment/backend/handler"
)

func main() {
	store := data.NewMockStore() // 메모리 기반 MockStore 초기화

	router := mux.NewRouter()

	// CORS 설정:
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 모든 Origin 허용
			w.Header().Set("Access-Control-Allow-Origin", "*")
			// 허용할 HTTP 메서드 (OPTIONS 포함)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
			// 허용할 요청 헤더 (특히 "Content-Type"이 중요)
			// "Content-Type"과 "X-Requested-With"는 웹 브라우저가 AJAX 요청 시 자동으로 추가하는 경우가 많습니다.
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With")
			// 클라이언트가 자격 증명(쿠키, HTTP 인증)을 포함하여 요청을 보낼 수 있도록 허용 (필요시)
			// w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Preflight 요청 (OPTIONS 메서드) 처리: 실제 요청 전 브라우저가 보내는 사전 확인 요청
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK) // 200 OK 응답으로 Preflight 요청 성공을 알림
				return
			}
			// 실제 요청을 다음 핸들러로 전달
			next.ServeHTTP(w, r)
		})
	})

	// API 엔드포인트 정의 및 핸들러 함수 연결
	router.HandleFunc("/issue", handler.CreateIssueHandler(store)).Methods("POST")
	router.HandleFunc("/issues", handler.GetIssuesHandler(store)).Methods("GET")
	router.HandleFunc("/issue/{id}", handler.GetIssueByIDHandler(store)).Methods("GET")
	router.HandleFunc("/issue/{id}", handler.UpdateIssueHandler(store)).Methods("PATCH")

	port := ":8080"
	log.Printf("서버가 %s에서 시작되었습니다...", port)
	log.Fatal(http.ListenAndServe(port, router))
}
