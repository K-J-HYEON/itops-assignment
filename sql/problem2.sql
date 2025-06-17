-- 요구사항:
-- 1. ETD 기간: 2025년 5월 1일 ~ 2025년 5월 11일
-- 2. POL(출발항) + 컨테이너 타입(CNTR_TYPE) 별로 집계
-- 3. 컨테이너 수량(CNTR 개수) 및 총 중량 합계(CNTR_WGT) 조회
-- 4. 결과 컬럼: POL_CD, CNTR_TYPE, CNTR_COUNT, TOTAL_WGT
-- 5. 정렬: POL_CD, CNTR_TYPE

-- ETD 기간 내 POL별 컨테이너 타입별 컨테이너 수량 및 중량 합계 조회
SELECT
    MST.POL_CD,          -- 출발항 코드 (요구사항 4: 결과 컬럼에 포함)
    CNTR.CNTR_TYPE,      -- 컨테이너 타입 (요구사항 4: 결과 컬럼에 포함)
    COUNT(CNTR.CNTR_NO) AS CNTR_COUNT, -- 컨테이너 수량 (요구사항 3: CNTR 개수 조회)
    SUM(CNTR.CNTR_WGT) AS TOTAL_WGT  -- 총 중량 합계 (요구사항 3: CNTR_WGT 합계 조회)
FROM
    FMS_HBL_MST MST      -- HOUSE B/L MASTER 테이블
JOIN
    FMS_HBL_CNTR CNTR    -- 컨테이너 테이블
ON
    MST.HBL_NO = CNTR.HBL_NO -- HBL_NO를 기준으로 두 테이블 조인
WHERE
    MST.ETD BETWEEN '20250501' AND '20250511' -- 지정한 ETD 기간 내 필터링 (요구사항 1: ETD 기간 적용)
GROUP BY
    MST.POL_CD,          -- POL(출발항) 별 그룹화 (요구사항 2: POL별 집계)
    CNTR.CNTR_TYPE       -- 컨테이너 타입 별 그룹화 (요구사항 2: 컨테이너 타입별 집계)
ORDER BY
    MST.POL_CD ASC,      -- POL_CD 오름차순 정렬 (요구사항 5: 정렬 조건)
    CNTR.CNTR_TYPE ASC;  -- CNTR_TYPE 오름차순 정렬 (요구사항 5: 정렬 조건)