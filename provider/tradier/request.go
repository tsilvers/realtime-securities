package tradier

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func (t *Tradier) request(endPoint string) ([]byte, error) {
	t.Lock()
	defer func() {
		// Update the last request timestamp.
		t.lastMillis = time.Now().UnixNano() / nanoToMilli

		t.Unlock()
	}()

	// Check request rate limit.
	nowMillis := time.Now().UnixNano() / nanoToMilli
	elapsedMillis := nowMillis - t.lastMillis
	millisLeft := tradierRequestDelayMillis - elapsedMillis
	if millisLeft > 0 {
		time.Sleep(time.Duration(millisLeft) * time.Millisecond)
	}

	// Create the request.
	requestURL := tradierURL + endPoint
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, errors.New("invalid request: " + err.Error())
	}

	req.Header.Set("Authorization", t.auth)
	req.Header.Set("Accept", tradierFormat)

	// Execute the request.
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, errors.New("error encountered during request: " + err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error encountered while reading response: " + err.Error())
	}

	if len(body) == 0 {
		return nil, nil
	}

	// Start: Show formatted json response.
	//fmt.Println("Request:", requestURL)
	//rawMap := make(map[string]interface{})
	//if err := json.Unmarshal(body, &rawMap); err != nil {
	//	fmt.Println("*** Error parsing response as json:", err)
	//	fmt.Println("Original response:")
	//	fmt.Println(string(body))
	//} else {
	//	if pretty, err := json.MarshalIndent(rawMap, "", "  "); err != nil {
	//		fmt.Println("*** Error formatting json response:", err)
	//		fmt.Println("Original response:")
	//		fmt.Println(string(body))
	//	} else {
	//		fmt.Println("JSON formatted response:")
	//		fmt.Println(string(pretty))
	//	}
	//}
	//fmt.Println()
	// End: Show formatted json response.

	return body, nil
}

func showRequestError(url string, err error) {
	log.Printf("Error encountered processing request: %s\nError: %v\n", url, err)
}

func showResponseErrorFunc(url string) func(error, ...string) {
	return func(err error, contextInfo ...string) {
		log.Printf("Error encountered processing response for request: %s\n", url)
		log.Printf("  Error: %v\n", err)
		for _, info := range contextInfo {
			log.Printf("  %s\n", info)
		}
		log.Println()
	}
}
