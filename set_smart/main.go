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
		{Name: "JSESSIONID", Value: "A814B1C212C2DE78CE4BDC02FACB480E.localhost", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
		{Name: "_gid", Value: "GA1.2.45474629.1720413854", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
		{Name: "SET_COOKIE_POLICY", Value: "1.0.0", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
		{Name: "access_grant", Value: "eyJlbmMiOiJBMTI4R0NNIiwiYWxnIjoiUlNBLU9BRVAifQ.nOqhIOREr4RqM-IabCt6Ke3vfVZMWkBmMlKADd6nUmwpFDDdEGX1AWQuNp157ieTrJw-e0r7itBi83nwE2pVlX42bVh9dqdhVdxeEdKkyg9tOiTyB7mAW8yUC-ZLTzox98R-aj0XRwPM00AqKppFFZKagaPJKpP5IiYIeTlXfQ3zjD-xC59Vkvcr00Hdb_e2Qpca7KDx19OkWu5HGiBMNTWStY5Owv9qJKPJp37KEokEYsqD1luKdxrZ86oitNlyp8z7Pht2a-yjgIncR_CiCdmhpDRiHGKqogTAINfAYlECVSYImdIDA1w2CAIvqNXgmxfr7H77XG1BFK6hpgSb-kPivOJDfq9Gsy6xpQ1VbnMhW3pHIfYsT0MZonxF6lslpnTJ3KkwRrIqCAKoKRBuq_s6EidZXnc20JwhxHrRHnvM9eF3pcPvFDXFscm2NYL_53LauV56NhsssQ6IR6bEOpuygUu-tZQniRr_0Nuj-mr4L_bCl5fqKp3eKpew0KETxKKEp-sIg23xY23ZWuM-9-qoc9IgImzd1bPZys9ZZkmz4nObl7B0rgykpo5fieU9BndM_4CUyzygNNp25x7rWUNo96VcFaum1if8wADbWOTGH6JK_pHIJfgvrRZ4Tdcysur1Au4ZEOYx7wEWlomr0RWu4Vxsihm99HAcw5LGVdw.RQ8tHNJJqMRir9d_.0IoO_e_MYsWsJ8lUptGdw7rCV6HZgq-WN5nq_Kb33Y6S_CZ1Z1KOftmHOMLfxKKZqCmaQrC2YfleH0gV9aXlChPik2RrSI5OrRSEJ96tBGDhTc712KHH2fPFFD5U9eGsVX5N7qlVnp4G-2VcrUkgnu0yFa2OlXLrj-ioQgNroNUspoRitEszY4M6j2lK8pyuNNS3pQkxO0VMN-zcQkx1E2aAMxH5g4eShRNOmKfMxjWUJNY7RRAvnlvwi8Pj1ociVnEU2nvL6uPrrfqqe-4pj8a7S9NBYKWgkAVLLJ1ELINj85Lpx4eyz4Q5QW1gj_Hlwwol8eY8XkKHntvoOIalz_sWlFCVwdyucj1TyqSvQLm7yKSs9rdTz-wpwIKlreGFm9KxqVhHnAo6PN5cbZk2OdGD4jdGs0bQqzi6a9BtN8PdaPslqmTq4CKWmvJLB5XJVdhfYY6p8P7fR8fJOUULNwcumbHjwwdxVs6MQkK3ZKYXUHtTVl4KoKrP9dpfpoWa4jKSBLRNSf6jfaaTjav0xUY519xSYcJfLtrFXkBsdiz5aeF-oV2PyBspYIW-jpGw3rvrXZWpGW_QUGhcl5rVgDxZ0GdIJp1Ncz8BGsGEdrFbDFQLx0F0zP4m3gtIVDjfenCnK1Iy-cyDeoh_19WN0tUcXccYOEOACvG8M6Wua5hPmFCwzz5BTgM-mZY1nc3j4OE6H-UfmlcvNR2UPUbJ3vUn0ZOHEuGfKpAOTfbu4zo4PY-t9kSq5On3liCVRsFy4_-9SjQMgK5MVmhIiCkq1Xc3nG2t-YQ6b0oSd80Y_wLcOkcQoqVstgm8b8gvfryQqpw11fUmjVfPKFl5VagnHYW3yIPH1hbvgDFcJCOp-Njhm6JNG4Fq8I_YiCyplx8wg0B-fr-hc12T5MZt7dNa20Znhd5WRBfyTpVLUMmJO7sIJVtDe-OYxewBzh1_yM8-7GIGEmtFDU8bcfxUOiyWvzMLcGfxN9kvCSN_hUz1FzIFxKu8hVJwdi0LEw", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
		{Name: "_ga", Value: "GA1.1.434521276.1720413853", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
		{Name: "route", Value: "458c7c73f43669bbf1b7afd17ec1dfdc", Domain: "www.setsmart.com", Path: "/", Expires: time.Now().Add(24 * time.Hour)},
	}

	client, err := authutils.SetupClientWithToken(cookies)
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
