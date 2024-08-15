package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// URL and payload
	url := "https://app.bot.or.th/1213/MCPD/ProductApp/Credit/CompareProduct"
	payload := "forgeryToken=AWNUxhW3op7dB8gp4w1Xmx-5AR7AbnpZwmA-ekOqgighc9QymI3Mjb0catTJ5WJBvpeM-A7pDNLz0ExtthM99dGVFwxu2WHG4pIaEuXlROo1%2CCe-0gPmf4KwT3cjvZyw1IGJO4I2R84Yk0gwdnOag1uwhgF-57B_C_-mJ9lIED4s2FHwQl9GrYz0nCgy5A2bIzbLYW1zO8Q8VdwmplCDW2sI1&ProductIdList=3629%2C3622%2C3630%2C4633%2C4632%2C4634%2C4655%2C4656%2C1604%2C1607%2C2259%2C5573%2C5161%2C2251%2C5534%2C5537%2C5177%2C5491%2C5522%2C5523%2C3627%2C3631%2C3600%2C3664%2C3624%2C3601%2C3625%2C3603%2C3633%2C3626%2C3604%2C3661%2C3620%2C3634%2C3636%2C3635%2C3606%2C3638%2C3637%2C3640%2C3639%2C3605%2C3672%2C3668%2C3628%2C3662%2C3621%2C3648%2C3607%2C3642%2C3649%2C3643%2C3651%2C3613%2C3644%2C3652%2C3609%2C3645%2C3653%2C3646%2C3647%2C3615%2C3656%2C3610%2C3650%2C3671%2C3658%2C3655%2C3616%2C3670%2C3612%2C4653%2C4657%2C4658%2C4659%2C4660%2C5568%2C5528%2C2256%2C5570%2C2245%2C5531%2C4471%2C4472%2C4467%2C5479%2C5494%2C5484%2C5497%2C5483%2C5503%2C1603%2C4760%2C4765%2C5473%2C3804%2C3800%2C4445%2C4482%2C4483%2C4479%2C4480%2C4452%2C4453%2C4476%2C4477%2C4454%2C4456%2C4457%2C4458%2C4460%2C4461%2C4462%2C4463%2C4464%2C4465%2C4444%2C4446%2C4447%2C4449%2C4450%2C4473%2C4470%2C5486%2C5498%2C5501%2C5540%2C5555%2C2242%2C5560%2C5496%2C1605%2C1600%2C1601%2C1602%2C4475%2C2246%2C5563%2C5565%2C5539%2C5567%2C5525%2C5556%2C2244%2C5562%2C2255%2C5536%2C2260%2C5574%2C2252%2C5546%2C5548%2C5538%2C5577%2C5518%2C5514%2C5521%2C5517%2C5527%2C3802%2C5529%2C5542%2C2258%2C5572%2C2249%2C5533%2C5516,4459%2C4474%2C4469%2C5499%2C5492%2C5481%2C5148%2C5435%2C3608%2C3641%2C3669%2C3611%2C3617%2C3660%2C5114%2C5287%2C5282%2C5285%2C5193%2C5137%2C5211%2C5173%2C5162%2C5156%2C5222%2C5126%2C5202%2C5158%2C5147%2C5240%2C5294%2C5296%2C5323%2C5306%2C5300%2C5329%2C749%2C753%2C752%2C751%2C750%2C754%2C5489%2C5500%2C5480%2C5495%2C4484%2C4481%2C4478%2C4455%2C4466%2C4451%2C3808%2C3797%2C3798%2C5502%2C5482%2C5566%2C2257%2C5571%2C2247%2C5532%2C5575%2C5576%2C5541%2C5554%2C2248%2C5544%2C2240%2C5558%2C5429%2C3996%2C3993%2C5520%2C3994%2C3995%2C3602%2C3632%2C3654%2C3665%2C3663%2C3618%2C3657%2C3623%2C1606%2C5526%2C5553%2C5552%2C3806%2C3799%2C3801%2C2253%2C5535%2C5564%2C3614%2C3666%2C3667%2C3619%2C3659%2C2250%2C5545%2C2241%2C5559%2C5549%2C2254%2C5547%2C2243%2C5561%2C5550%2C4448%2C4468%2C5519%2C5493%2C5266%2C5305%2C5304%2C5246%2C4654%2C5462%2C5213%2C5233%2C5478%2C5467%2C3997%2C5180%2C5557%2C5543%2C5551%2C5524%2C5335"

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", "verify=test; verify=test; verify=test; mycookie=!IsS8/5OcS8Q0oZFt1qgTjHeotiaw/uf3hv2XArUSbliquroZB6FlPVm+jBUqejYL+P5bGsj9RsNitRreh2XsxLuXC1jpnlz1ck93CzQKDHU=; _cbclose=1; _cbclose6672=1; _ga=GA1.1.412122074.1723533133; AMCVS_F915091E62ED182D0A495F95%40AdobeOrg=1; AMCV_F915091E62ED182D0A495F95%40AdobeOrg=179643557%7CMCIDTS%7C19951%7CMCMID%7C53550622918316951353729640026118558196%7CMCAAMLH-1724305541%7C3%7CMCAAMB-1724305541%7CRKhpRz8krg2tLO6pguXWp5olkAcUniQYPHaMWWgdJ3xzPWQmdj0y%7CMCOPTOUT-1723707941s%7CNONE%7CvVersion%7C5.5.0; _ga_R8HGFHEVB7=GS1.1.1723700741.1.0.1723700743.58.0.0; _uid6672=16B5DEBD.4; _ctout6672=1; _ga_NLQFGWVNXN=GS1.1.1723704905.5.1.1723705211.53.0.0")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	// Load the HTML document using goquery
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	// Extract and clean all text content from the document
	allText := doc.Text()

	// Use regular expressions to remove JavaScript, CSS, and other unwanted content
	re := regexp.MustCompile(`(?s)<script.*?>.*?</script>|<style.*?>.*?</style>|<.*?>`)
	cleanedText := re.ReplaceAllString(allText, "")

	// Remove excessive whitespace
	cleanedText = strings.TrimSpace(cleanedText)
	cleanedText = strings.ReplaceAll(cleanedText, "\n", " ")
	cleanedText = strings.ReplaceAll(cleanedText, "\t", " ")
	cleanedText = strings.Join(strings.Fields(cleanedText), " ")

	// Optionally, save the cleaned text to a file
	err = os.WriteFile("cleaned_text.txt", []byte(cleanedText), 0644)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	fmt.Println("Cleaned text content extracted and saved to cleaned_text.txt")
}
