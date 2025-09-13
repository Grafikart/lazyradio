package radio

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type RadioItem struct {
	Name   string
	Artist string
}

func FetchLastNovaTracks() ([]RadioItem, error) {
	// f, err := os.Open("radio/nova.html")
	f, err := fetchNovaPage()
	if err != nil {
		return nil, fmt.Errorf("failed to open radio playlist: %w", err)
	}
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse radio document: %w", err)
	}
	var items []RadioItem
	fmt.Printf("%+v", doc.Find(".wwtt_content .wwtt_right").Length())
	doc.Find(".wwtt_content .wwtt_right").Each(func(i int, s *goquery.Selection) {
		items = append(items, RadioItem{
			Name:   s.Find("h2").First().Text(),
			Artist: s.Find("p").Last().Text(),
		})
	})

	return items, nil
}

func fetchNovaPage() (io.Reader, error) {
	// Get current time and format it as H:M
	now := time.Now()
	currentTime := fmt.Sprintf("%d:%d", now.Hour(), now.Minute())

	// Prepare the form data
	data := url.Values{}
	data.Set("action", "loadmore_programs")
	data.Set("date", "")
	data.Set("time", currentTime)
	data.Set("page", "1")
	data.Set("radio", "910")

	// Create the request
	req, err := http.NewRequest("POST", "https://www.nova.fr/wp-admin/admin-ajax.php", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "fr-FR,fr;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Origin", "https://www.nova.fr")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Referer", "https://www.nova.fr/c-etait-quoi-ce-titre/?radio=radio-nova")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func main() {
	reader, err := fetchNovaPage()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer reader.(io.ReadCloser).Close()

	// Read and display response
	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Response: %s\n", string(body))
}
