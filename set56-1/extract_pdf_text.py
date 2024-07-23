import sys
import fitz

def extract_text_from_pdf(file_path):
    try:
        pdf_document = fitz.open(file_path)
        text = ""
        for page_num in range(len(pdf_document)):
            page = pdf_document.load_page(page_num)
            text += page.get_text()
        return text.encode('utf-8')  # Encode the text in UTF-8
    except Exception as e:
        return f"Error extracting text: {e}".encode('utf-8')

def main():
    if len(sys.argv) < 2:
        print("Usage: extract_pdf_text.exe <path_to_pdf>")
        sys.exit(1)

    file_path = sys.argv[1]
    text = extract_text_from_pdf(file_path)
    sys.stdout.buffer.write(text)  

if __name__ == "__main__":
    main()
