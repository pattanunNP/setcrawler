package utils

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func SetupCookieAndToken() (*http.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	// Set the cookies from environment variables
	baseURL, _ := url.Parse("https://www.setsmart.com")
	cookies := []*http.Cookie{
		{Name: "JSESSIONID", Value: os.Getenv("JSESSIONID")},
		{Name: "_gid", Value: os.Getenv("_gid")},
		{Name: "SET_COOKIE_POLICY", Value: os.Getenv("SET_COOKIE_POLICY")},
		{Name: "access_grant", Value: os.Getenv("access_grant")},
		{Name: "route", Value: os.Getenv("route")},
		{Name: "_ga", Value: os.Getenv("_ga")},
		{Name: "_ga_W4KEXZN4YX", Value: os.Getenv("_ga_W4KEXZN4YX")},
		{Name: "_ga_6WS2P0P25V", Value: os.Getenv("_ga_6WS2P0P25V")},
		{Name: "JSESSIONID2", Value: os.Getenv("JSESSIONID2")},
	}
	jar.SetCookies(baseURL, cookies)

	return client, nil
}
