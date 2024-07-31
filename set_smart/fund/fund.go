package fund

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IPOYear struct {
	Year string `json:"year"`
}

type Fund struct {
	FundName string `json:"fundname"`
}

func FetchIPOYears(cookieStr string) ([]IPOYear, error) {
	url := "https://www.setsmart.com/api/fund/ipo-year/list"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("authorization", "Bearer eyJlbmMiOiJBMTI4R0NNIiwiYWxnIjoiUlNBLU9BRVAifQ.0iLYEVb3OZMdUl0bVNmltA8D48AT7Sxsz2JFeyZUz7uCF5f2MRxaRsVkxFAiMfCuKaAO6bWnU0rbA5E-2tYT4-pTF0PsU4AIXz2E5z-0isf8HE4JADt8D0fNJVUD8H2tSu9SKRNpIAWS8CkNiP_rTa1eEmoprK6PtATcTUsq6x-t7LvpH7V1DvPnJa8bsLgD68Q8Sn2tVE5ULvckrfHknLbPlep-G9k65jKd5RUR9W9OPEWJLEyg8r-DSsHCpRERrh9Mq0SBmREEbsV8zxZy1uvVOvGSeV5HpUrAxtgwFJhpWUyNZ8YOyXsZIlvHb17YFAYNqoDLYVcthXjeVgE81jQ76Qv1vP8DtVJ4m4U5bGsy_5l2QBlcD2Tz7H3SoWzRKxFZEDKLBvO_oyOw9WJ9KibBK9D4JVOWydAVWh3YprmzU4XpKuGxyiKzl97Mb21EG1pCNhYp5FxHyUvme9xw6vLhBoKZit8ZjimW-ShpZOiOdWrP9YfMVyAEOzrCVxBJm9g4xgumLUu_q_gxweL6gFZSSruf1yFZR04ymPchmNQ4yv8t-HtgGcFnlk9F0bFTpmkf_SdmbO6Ajgk8wXLPjgJp69J1vIOA29J2KdEm9OJ9UEbhx8sF5I24LYkXyjGorKrxoh8aYvtdpWLeACc09XISehKQ4TTOVWuI0ZZ-O5E.G7mEo3C3STWyYyRs.v4SEx2VRQTkVQRQUQ9NmuwRpMd19O5WNDF39ILdZEJ5PiDu-HmYAN4j8HbIkISPl-EGYfjJ5XUw19arQkqwSoPeuzyUAGiDWeQFMuoUJJu2f7CLkWZbeJfbqwwAEgl7vVm8gnW7Du63iihQl_GIhSgbMQvteve7yK1KJKzqaqsehWmX_ZyY07zJmHhqP73fj87FoyBUndV193Gwhoq21rkJAps52Iz14JY0tatxzL4UOuebWadmWdn3zOW0AYk-txH0M0MW-HOF9nETVURhsC6_ppPMTJUvviYwCm29OiwMN9PRxRHMCwwkl7y0Uo6Mb9fRtgeiQtO8RqFHFle02vieniQTJqiVRNaWgZNHKDrUix7TLXcF9rE1FsPWLso8ctEv8BaROaUgzEVztU2JI1BbQvpcW7KQcAoligO7pZsR63PrqErtLkaFGiwCHIZendA7nYizgeg8qKU7xW1Dla6qGi9EqYJ-KFI-DuiLAKi36P3xE-f7KjkX5x4qMg5X87eXYC8sbnLkH-L1NeHaXg4rNPtXjwOrh6p6IJ9PnTZM0j_d2cGphwTK7FyIYB0ou1uBfqFwf1vE4VuI_ust0aEzMBd4NDWFr9sOf7_ZEi6kRTw2RVd9CeOHDvXplVOtKYTuCPYsBsq459lnukDmd1xhhOZDd4gNjcu1SWb1vt-zkVwsGn6bvzROVbusBzcUsv3Y83zN69EyuUo_H8ZuLGy6QN9-HkvH_UmKjpWflqRoMCfNXncCtqqQOEepcSZfblRQ7lorU6ScMLfHc-213wB_F_AtKFNyKjYIDjCbNCIrCmH4SsVSpawxURMQvDUrkEWtogQfuAjiTHzQp65_-X8YPvahFlx5UrP4YbIdqs5lvglIJ_uAgUJjtZiLSHrdsl7aWK8xZSntNIZkIX8O6NcqE6jZmP7d3ds4scSZtnS1BQof6Z7auMXBOkKfTGp8_d8-HOMC26etbxdR3gLPLNufleC_dqyYfwfn1QuirVpYu9LrIr23yJK0tvg.uyLoG_xvzzN0PTGa1Vr5Mg")
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("receive non-200 response code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var years []IPOYear
	err = json.Unmarshal(body, &years)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return years, nil
}

func FetchFundsByYear(cookieStr, year string) ([]Fund, error) {
	url := fmt.Sprintf("https://www.setsmart.com/api/fund/ipo/list?year=%s", year)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Add("accept", "application/json, text/plain, */*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Add("authorization", "Bearer eyJlbmMiOiJBMTI4R0NNIiwiYWxnIjoiUlNBLU9BRVAifQ.0iLYEVb3OZMdUl0bVNmltA8D48AT7Sxsz2JFeyZUz7uCF5f2MRxaRsVkxFAiMfCuKaAO6bWnU0rbA5E-2tYT4-pTF0PsU4AIXz2E5z-0isf8HE4JADt8D0fNJVUD8H2tSu9SKRNpIAWS8CkNiP_rTa1eEmoprK6PtATcTUsq6x-t7LvpH7V1DvPnJa8bsLgD68Q8Sn2tVE5ULvckrfHknLbPlep-G9k65jKd5RUR9W9OPEWJLEyg8r-DSsHCpRERrh9Mq0SBmREEbsV8zxZy1uvVOvGSeV5HpUrAxtgwFJhpWUyNZ8YOyXsZIlvHb17YFAYNqoDLYVcthXjeVgE81jQ76Qv1vP8DtVJ4m4U5bGsy_5l2QBlcD2Tz7H3SoWzRKxFZEDKLBvO_oyOw9WJ9KibBK9D4JVOWydAVWh3YprmzU4XpKuGxyiKzl97Mb21EG1pCNhYp5FxHyUvme9xw6vLhBoKZit8ZjimW-ShpZOiOdWrP9YfMVyAEOzrCVxBJm9g4xgumLUu_q_gxweL6gFZSSruf1yFZR04ymPchmNQ4yv8t-HtgGcFnlk9F0bFTpmkf_SdmbO6Ajgk8wXLPjgJp69J1vIOA29J2KdEm9OJ9UEbhx8sF5I24LYkXyjGorKrxoh8aYvtdpWLeACc09XISehKQ4TTOVWuI0ZZ-O5E.G7mEo3C3STWyYyRs.v4SEx2VRQTkVQRQUQ9NmuwRpMd19O5WNDF39ILdZEJ5PiDu-HmYAN4j8HbIkISPl-EGYfjJ5XUw19arQkqwSoPeuzyUAGiDWeQFMuoUJJu2f7CLkWZbeJfbqwwAEgl7vVm8gnW7Du63iihQl_GIhSgbMQvteve7yK1KJKzqaqsehWmX_ZyY07zJmHhqP73fj87FoyBUndV193Gwhoq21rkJAps52Iz14JY0tatxzL4UOuebWadmWdn3zOW0AYk-txH0M0MW-HOF9nETVURhsC6_ppPMTJUvviYwCm29OiwMN9PRxRHMCwwkl7y0Uo6Mb9fRtgeiQtO8RqFHFle02vieniQTJqiVRNaWgZNHKDrUix7TLXcF9rE1FsPWLso8ctEv8BaROaUgzEVztU2JI1BbQvpcW7KQcAoligO7pZsR63PrqErtLkaFGiwCHIZendA7nYizgeg8qKU7xW1Dla6qGi9EqYJ-KFI-DuiLAKi36P3xE-f7KjkX5x4qMg5X87eXYC8sbnLkH-L1NeHaXg4rNPtXjwOrh6p6IJ9PnTZM0j_d2cGphwTK7FyIYB0ou1uBfqFwf1vE4VuI_ust0aEzMBd4NDWFr9sOf7_ZEi6kRTw2RVd9CeOHDvXplVOtKYTuCPYsBsq459lnukDmd1xhhOZDd4gNjcu1SWb1vt-zkVwsGn6bvzROVbusBzcUsv3Y83zN69EyuUo_H8ZuLGy6QN9-HkvH_UmKjpWflqRoMCfNXncCtqqQOEepcSZfblRQ7lorU6ScMLfHc-213wB_F_AtKFNyKjYIDjCbNCIrCmH4SsVSpawxURMQvDUrkEWtogQfuAjiTHzQp65_-X8YPvahFlx5UrP4YbIdqs5lvglIJ_uAgUJjtZiLSHrdsl7aWK8xZSntNIZkIX8O6NcqE6jZmP7d3ds4scSZtnS1BQof6Z7auMXBOkKfTGp8_d8-HOMC26etbxdR3gLPLNufleC_dqyYfwfn1QuirVpYu9LrIr23yJK0tvg.uyLoG_xvzzN0PTGa1Vr5Mg")
	req.Header.Add("cookie", cookieStr)
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Read body for debugging purposes
		return nil, fmt.Errorf("received non-200 response code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var funds []Fund
	err = json.Unmarshal(body, &funds)
	if err != nil {
		return nil, fmt.Errorf("error unmarshling JSON: %w", err)
	}

	return funds, nil
}
