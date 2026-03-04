package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"lottery-app/internal/storage"
)

const baseURL = "https://datachart.500.com/ssq/history/newinc/history.php"

var (
	rowRe     = regexp.MustCompile(`<tr class="t_tr1">.*?</tr>`)
	tdValRe   = regexp.MustCompile(`<td[^>]*>([^<]*)</td>`)
	commentRe = regexp.MustCompile(`<!--.*?-->`)
)

// Fetch 从500彩票网抓取最近N期双色球数据
func Fetch(count int) ([]storage.LotteryRecord, error) {
	body, err := httpGet(baseURL)
	if err != nil {
		return nil, err
	}

	records, err := parseHTML(body)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("未找到数据，请检查网络连接")
	}

	if count <= len(records) {
		return records[:count], nil
	}

	// 需要更多数据，按期号范围重新请求
	latestIssue, _ := strconv.Atoi(records[0].Issue)
	startIssue := calcStartIssue(latestIssue, count)

	url := fmt.Sprintf("%s?start=%d&end=%d", baseURL, startIssue, latestIssue)
	body, err = httpGet(url)
	if err != nil {
		return nil, err
	}

	records, err = parseHTML(body)
	if err != nil {
		return nil, err
	}
	if len(records) > count {
		records = records[:count]
	}
	return records, nil
}

func calcStartIssue(latest, count int) int {
	year := latest / 1000
	seq := latest % 1000
	if seq >= count {
		return year*1000 + seq - count + 1
	}
	return (year-1)*1000 + 1
}

func httpGet(url string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://datachart.500.com/ssq/")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func parseHTML(body string) ([]storage.LotteryRecord, error) {
	rows := rowRe.FindAllString(body, -1)
	records := make([]storage.LotteryRecord, 0, len(rows))

	for _, row := range rows {
		row = commentRe.ReplaceAllString(row, "")
		matches := tdValRe.FindAllStringSubmatch(row, -1)
		if len(matches) < 9 {
			continue
		}

		issue := strings.TrimSpace(matches[0][1])

		red := make([]int, 6)
		for i := 0; i < 6; i++ {
			n, err := strconv.Atoi(strings.TrimSpace(matches[i+1][1]))
			if err != nil {
				return nil, fmt.Errorf("解析红球失败 [%s]: %w", issue, err)
			}
			red[i] = n
		}

		blue, err := strconv.Atoi(strings.TrimSpace(matches[7][1]))
		if err != nil {
			return nil, fmt.Errorf("解析蓝球失败 [%s]: %w", issue, err)
		}

		date := strings.TrimSpace(matches[len(matches)-1][1])

		records = append(records, storage.LotteryRecord{
			Issue: issue,
			Date:  date,
			Red:   red,
			Blue:  blue,
		})
	}
	return records, nil
}
