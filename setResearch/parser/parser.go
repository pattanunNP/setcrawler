package parser

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

func ExtractFileURLFromHTML(htmlContent string) (string, error) {
	scriptStart := strings.Index(htmlContent, "window.__NUXT__")
	if scriptStart == -1 {
		return "", fmt.Errorf("NUXT script not found in HTML content")
	}
	scriptEnd := strings.Index(htmlContent[scriptStart:], "</script>")
	if scriptEnd == -1 {
		return "", fmt.Errorf("end of script tag not found")
	}
	scriptContent := htmlContent[scriptStart : scriptStart+scriptEnd]

	re := regexp.MustCompile(`"fileUrl"\s*:\s*"(https:\\u002F\\u002F[^"]+\.pdf)"`)
	matches := re.FindStringSubmatch(scriptContent)
	if len(matches) > 0 {
		cleanURL := strings.ReplaceAll(matches[1], `\\u00F`, `/`)
		log.Printf("Extracted fileUrl: %s\n", cleanURL)
		return cleanURL, nil
	}

	truncateContent := TruncateString(scriptContent, 2000)
	log.Printf("script Content (truncated): %s\n", truncateContent)

	return "", fmt.Errorf("fileUrl not found in script content")
}

func TruncateString(str string, num int) string {
	if len(str) > num {
		return str[:num] + "..."
	}
	return str
}
