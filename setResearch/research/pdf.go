package research

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func ProcessPDF(fileURL string) (string, error) {
	pdfData, err := downloadPDF(fileURL)
	if err != nil {
		return "", fmt.Errorf("error downloading PDF: %v", err)
	}

	tempFile, err := os.CreateTemp("", "pdf-*.pdf")
	if err != nil {
		return "", fmt.Errorf("error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(pdfData); err != nil {
		return "", fmt.Errorf("error writing to temporary file: %v", err)
	}
	tempFile.Close()

	pdfContent, err := extractTextFromPDF(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("error extracting text from PDF: %v", err)
	}

	return pdfContent, nil
}

func downloadPDF(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	pdfData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return pdfData, nil
}

func extractTextFromPDF(pdfPath string) (string, error) {
	cmd := exec.Command("java", "-jar", "pdfbox-app-3.0.2.jar", "export:text",
		"-encoding", "UTF-8", "-i", pdfPath, "-console")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("error running PDFBox CLI to extract text: %v", err)
	}

	return string(output), nil
}
