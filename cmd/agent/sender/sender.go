package sender

import (
	"fmt"
	"io"
	"net/http"
)

func SendPostRequest(url string) {

	// Create a new POST request with no body
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Optionally add headers
	req.Header.Set("Content-Type", "text/plain")

	// Use http.DefaultClient to send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body:", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}
	fmt.Println("Response:", string(body))
}
