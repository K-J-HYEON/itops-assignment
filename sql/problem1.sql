-- 요구사항:
-- 1. HOUSE B/L (HBL_NO) 를 기준으로 각 B/L 별 컨테이너 수량(CNTR 개수) 을 집계
-- 2. 가장 많은 수량을 가진 B/L 1건을 조회
-- 3. 컨테이너 수량은 CNTR_NO 기준으로 COUNT
-- 4. 동일 수량 시 ETD 빠른 순으로 우선 선택
-- 5. 결과 컬럼: HBL_NO, CNTR_COUNT, ETD
-- 6. 정렬: CNTR_COUNT DESC, ETD ASC

-- 컨테이너 수량이 가장 많은 B/L 1건을 조회하는 SQL문
SELECT
    MST.HBL_NO,          -- HOUSE B/L 번호 (요구사항 5: HBL_NO 포함)
    COUNT(CNTR.CNTR_NO) AS CNTR_COUNT, -- 각 B/L 별 컨테이너 수량 (요구사항 1, 3: CNTR_NO 기준으로 COUNT)
    MST.ETD              -- 출항일자 (요구사항 5: ETD 포함)
FROM
    FMS_HBL_MST MST      -- HOUSE B/L MASTER 테이블
JOIN
    FMS_HBL_CNTR CNTR    -- 컨테이너 테이블
ON
    MST.HBL_NO = CNTR.HBL_NO -- HBL_NO를 기준으로 두 테이블 조인 (요구사항 1: HBL_NO 기준 집계)
GROUP BY
    MST.HBL_NO,          -- HBL_NO 별로 그룹화
    MST.ETD              -- ETD도 함께 그룹화하여 집계에 포함 (ETD가 SELECT 절에 있으므로 GROUP BY에 포함)
ORDER BY
    CNTR_COUNT DESC,     -- 컨테이너 수량(CNTR_COUNT)이 많은 순서로 정렬 (내림차순) (요구사항 6)
    MST.ETD ASC          -- 컨테이너 수량이 동일할 경우, ETD(출항일자)가 빠른 순서로 정렬 (오름차순) (요구사항 4, 6)
FETCH FIRST 1 ROW ONLY;  -- 가장 많은 수량을 가진 B/L 중 첫 번째 1건만 조회 (Oracle 12c 이상) (요구사항 2)
