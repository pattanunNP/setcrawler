import sys
import pdfplumber

def extract_text_from_pdf(file_path):
    try:
        with pdfplumber.open(file_path) as pdf:
            text = ""
            for page in pdf.pages:
                text += page.extract_text()
        return text
    except Exception as e:
        return f"Error extracting text: {e}"
    
def main():
    if len(sys.argv) < 2:
        print("Usage: extract_pdf_text.exe <path_to_path>")
        sys.exit(1)

    file_path = sys.argv[1]
    text = extract_text_from_pdf(file_path)
    print(text.encode('utf-8', errors='replace').decode('utf-8',errors='replace'))

if __name__ == "__main__":
    main()