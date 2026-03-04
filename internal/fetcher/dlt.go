package fetcher

import (
	"fmt"
	"strconv"
	"strings"

	"lottery-app/internal/storage"
)

const dltBaseURL = "https://datachart.500.com/dlt/history/newinc/history.php"

// FetchDLT 从500彩票网抓取最近N期大乐透数据
func FetchDLT(count int) ([]storage.DLTRecord, error) {
	body, err := httpGet(dltBaseURL)
	if err != nil {
		return nil, err
	}

	records, err := parseDLTHTML(body)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("未找到数据，请检查网络连接")
	}

	if count <= len(records) {
		return records[:count], nil
	}

	latestIssue, _ := strconv.Atoi(records[0].Issue)
	startIssue := calcStartIssue(latestIssue, count)

	url := fmt.Sprintf("%s?start=%d&end=%d", dltBaseURL, startIssue, latestIssue)
	body, err = httpGet(url)
	if err != nil {
		return nil, err
	}

	records, err = parseDLTHTML(body)
	if err != nil {
		return nil, err
	}
	if len(records) > count {
		records = records[:count]
	}
	return records, nil
}

func parseDLTHTML(body string) ([]storage.DLTRecord, error) {
	rows := rowRe.FindAllString(body, -1)
	records := make([]storage.DLTRecord, 0, len(rows))

	for _, row := range rows {
		row = commentRe.ReplaceAllString(row, "")
		matches := tdValRe.FindAllStringSubmatch(row, -1)
		if len(matches) < 9 {
			continue
		}

		issue := strings.TrimSpace(matches[0][1])

		front := make([]int, 5)
		for i := 0; i < 5; i++ {
			n, err := strconv.Atoi(strings.TrimSpace(matches[i+1][1]))
			if err != nil {
				return nil, fmt.Errorf("解析前区失败 [%s]: %w", issue, err)
			}
			front[i] = n
		}

		back := make([]int, 2)
		for i := 0; i < 2; i++ {
			n, err := strconv.Atoi(strings.TrimSpace(matches[i+6][1]))
			if err != nil {
				return nil, fmt.Errorf("解析后区失败 [%s]: %w", issue, err)
			}
			back[i] = n
		}

		date := strings.TrimSpace(matches[len(matches)-1][1])

		records = append(records, storage.DLTRecord{
			Issue: issue,
			Date:  date,
			Front: front,
			Back:  back,
		})
	}
	return records, nil
}
