package main

import (
	"fmt"
	authutils "login_token/utils"
	"net/http"
	"time"
)

func main() {
	// Replace these with your actual cookies
	cookies := []*http.Cookie{
		{Name: "JSESSIONID", Value: "722B82A6B55C2B5E1516492B6BE1504D.localhost", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
		{Name: "_gid", Value: "GA1.2.45474629.1720413854", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
		{Name: "SET_COOKIE_POLICY", Value: "1.0.0", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
		{Name: "access_grant", Value: "eyJlbmMiOiJBMTI4R0NNIiwiYWxnIjoiUlNBLU9BRVAifQ.FgePM2Z9qFSQ0Q5JFN0jPOIx1iPqAW9tE6kC66YW765HX8tzNoOhC5kgOsvu4jdiwMUtWshGZfpv9vxzVbqzjUM_BCzLxjWxRi0Fhw6o8qp6lhtzuYLTtAPVtzXY-JlXWPtILyMJcL23KN-rM36leFqcCzJYp2WqVnOgKfOH6sSjjPiutQTiOgq3H4cKOqR-5beTk1UoYxt9Ki7g3F46-dIKq1-fmX9IlV7wDt9-IPZCz7acuqGZjbsaeC7fdZXn6np1Kqf1oRHd16lnZxvmW0Fp-VpTLUZ4wtCFuwb8GZTPGJmOCRZQMkfQyVkJJxpegw_-qweOw-SYLzoPO8jlVNeTyfRmtI1zV6v5TPdp48Txjubj6P_U6OERSoojp5EsmsZzN3Owsakea_S3DoTlqrtml5Oh2kTgUnzCYv39chsNjMsIA4uSlCM-U59nIpUSvDBTakxNncSUIg7LzdU8c0TzEVF4CaBAyXcoY34CXevoqQwTs-TVjrHUtfjTG6rIau5zDyjvQ_HPxe39OS6kMlfXitdqYTszyuM2nHMO7Qat3cgVi9GBS-2jr0ie0QpByqQCSv_Tk_BgHfGO2qHvIFrtlkRAHBoD-gt-lbEIt-rzNfiauEusayhH0KVNaOagJ87PMJlXaVWMkrLXr7OtnWePOItG2QRoFzqBbBWjAIs.hFeBhCsMHZSa2EUH.TQGZriQdDRXDmoPCEBSacMFFLcd2wkAYlawZA8iueHKaqHeipyAjjEjbzY8lwCPZexcz_fdgsZt0OkhQerzJWOT8NxPOmQ7eHk_YHJORotRoPspAv9urgvqvM6ApMoiFn-7OcDlZmsmdIiEgDxl_NOFW5uP0qvBFbn6L6MYbAF5p7hlfReiZYDGuIjpJCRI2NlBdisIsiNlzZQOhCpegxdZ76rijRjtG5mI0cbFP_iAy152xNirpV02fGpGGQK7xrdu42M4D9m9bz6PrUx1l9aVZEKZQilGd_frJ02KXonSkB6SJ7iO0BgehwD5iDC878hI93sbYw-Q8FpQ2s1dOjr69izq3AR5iknVZEgcNJ3z8OwtZ7Phjjr0HrYtDwgJn4PsK2To_DH_MX9SSnQEN4eyHHe1GgiV2GKklBMnijeL7z5i4KiYL0AVLxSnkrYq41FX6m430hZa1YqUGsNNpRn08CJNRiknbeqccLIs9WGqeh_WdnTH7k2IMyyby6xiPkTlQLWZOSWce8tVDFERTKkWhZyzeY-F9JMLRXrza2Nw2yFPoGa4o8iWMzMmxy-M1jdlx0fC6g6404DoF2LgHBAu37RAO8cXBk8aqtXTDrTrC89rC_YzmudTnZaHjIkx3OwBYOxlxbUCwFtg1GvuO8BsjiRqq9c7yDptbZWV9adQuQ4FG-P39rBojCeZtvwSvh9Jwqpb-hCaIIYmhZXFjb42frfm7FKMdTWOH4uGmviBviQmMGg-5rWyjtjCpemV0-TkR6MtfQ9IWKwh7Lszv6hUfKbtLVfPmcrp9UXfenTT1FKMDJdphwZjTEVzxusqDKZM3Hs8fpcbeEmR3UaAUt-GuN7h53qMebs5O50HRp0bev3jWircG63DidMU1UIe9MnCVuPDmu3vQ9vo5JiHrLJw2hvYpnbd6fcs1098OOWqDDB8SfTkEJZWSpFV3R6wZO_e3idXn-3garB2baSKv9SXKMgRfw7oJtykhmae4HFBMWg4iYE6E-Hot3A.b_AFgaCiV20MEsSCcuRTNQ", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
		{Name: "_ga", Value: "GA1.1.434521276.1720413853", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
		{Name: "route", Value: "ea93f49e827559a4637ccec88d06e835", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
	}

	client, err := authutils.SetupClientWithCookies(cookies)
	if err != nil {
		fmt.Println("Error setting up client with cookies:", err)
		return
	}

	fmt.Println("Cookies:")
	authutils.PrintCookies(cookies)

	// Headers from your cURL request
	headers := map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Accept-Language":           "en-US,en;q=0.9,th-TH;q=0.8,th;q=0.7",
		"Priority":                  "u=0, i",
		"Sec-CH-UA":                 `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`,
		"Sec-CH-UA-Mobile":          "?0",
		"Sec-CH-UA-Platform":        `"macOS"`,
		"Sec-Fetch-Dest":            "document",
		"Sec-Fetch-Mode":            "navigate",
		"Sec-Fetch-Site":            "cross-site",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
	}

	// Example request to the specified page
	requestURL := "https://www.setsmart.com/ism/companyprofile.html?locale=en_US"
	response, err := authutils.MakeAuthenticatedRequest(client, requestURL, headers, cookies)
	if err != nil {
		fmt.Println("Error making custom request:", err)
		return
	}

	fmt.Println("Custom Request Response:")
	fmt.Println(response)
}
