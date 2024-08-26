package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// FetchData makes an HTTP request to the given URL and returns the response body
func FetchData(payload string) ([]byte, error) {
	const url = "https://app.bot.or.th/1213/MCPD/ProductApp/NanoFinance/CompareProductList"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Setting headers (same as before)
	req.Header.Set("accept", "text/plain, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("cookie", `verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; RT="z=1&dm=app.bot.or.th&si=e09a0124-640e-4621-b914-91629b46456e&ss=m07v2kro&sl=0&tt=0"; _uid6672=16B5DEBD.27; _ga_NLQFGWVNXN=GS1.1.1724641558.32.0.1724641558.60.0.0; visit_time=2677`)
	req.Header.Set("origin", "https://app.bot.or.th")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://app.bot.or.th/1213/MCPD/ProductApp/NanoFinance/CompareProduct")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="99", "Google Chrome";v="127", "Chromium";v="127"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("verificationtoken", `r-0vhfITJl4VEk_yQxi3boYllKaWXnfbYYQsTgPF1ym1U3Axe4vh4mlOeoFL1rbFyrqftcgVni_AuWQymKri9L-caTpQenyeKLStAiJc5Cg1,dRC9vDGwBYzteur_1nlNmZA_kdP-YokSaPiUSD4hiNIp_okYbMV6WhLjbGdQB3GAmQv4ShTVrHGf1ALwTliKJ4YA2MSe8-ZPO7P9zkFVTYo1`)
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("Status Code: %d\n", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}
