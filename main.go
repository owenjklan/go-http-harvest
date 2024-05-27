package main

import (
	"bytes"
	"flag"
	"fmt"
	nw "go-http-harvest/netobjects"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"sync"
	"time"
)

type ScanConfig struct {
	port           uint16
	connectTimeout time.Duration
}

type ConnectedResponse struct {
	Url         string            `json:"url"`
	Status      int               `json:"status"`
	Server      string            `json:"server"`
	Host        string            `json:"host"`
	Extras      map[string]string `json:"extras"`
	extrasStack []string          // names of extras functions to dispatch to
}

func (cr ConnectedResponse) String(showExtras bool) string {
	var headerString = fmt.Sprintf("%-16s|[%d] <Server: %-30s> | %s", cr.Host, cr.Status, cr.Server, cr.Url)
	if !showExtras || cr.Extras == nil {
		return headerString
	}

	if len(cr.Extras) == 0 {
		return headerString
	}

	// Iterate through "Extras" map
	const boldText = "\033[1m"
	const normalText = "\033[0m"
	const headerIndent = 2
	const contentIndent = 4
	var buffer bytes.Buffer
	buffer.WriteString(headerString + "\n") // The always-present header summary
	for header, content := range cr.Extras {
		buffer.WriteString(strings.Repeat(" ", headerIndent) + boldText + header + normalText + ":\n")

		// Iterate through lines of content, adding indent appropriately
		for _, contentLine := range strings.Split(content, "\n") {
			buffer.WriteString(strings.Repeat(" ", contentIndent) + contentLine)
		}
	}
	return buffer.String()
}

func main() {
	// REMEMBER: flag package returns *pointers*
	port := flag.Int(`port`, 80, "Port to attempt connections on")
	connectTimeout := flag.Int("connect-timeout", 5, "Connection timeout in seconds")
	var successList []ConnectedResponse

	// parse flags
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: go-http-harvest BASE-IPV4-ADDRESS [OPTIONS]")
		fmt.Printf("\nScan a 'Class C' netobjects block, based off the supplied IP address\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Create Scan config struct for scanner goroutines to use with their clients
	scanConfig := ScanConfig{
		port:           uint16(*port),
		connectTimeout: time.Duration(*connectTimeout * int(time.Second)),
	}

	// Create TargetNetwork for iterating randomly through class C block
	baseAddressString := flag.Arg(0)
	targetNetwork, err := nw.NewTargetNetwork(baseAddressString)

	if err != nil {
		fmt.Printf("Error trying to determine base address: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Target Network: ", targetNetwork)
	fmt.Printf("%s. Will connect on port %d\n", baseAddressString, *port)

	wg := sync.WaitGroup{}
	startTime := time.Now()

	// Start goroutines for the success and error handlers
	successChannel := make(chan *http.Response)
	// Create channel to feed addresses into scanning goroutines
	addressChannel := make(chan netip.Addr)

	// consumes *http.Response from successChannel, puts processed results in successList
	go processSuccess(successChannel, &successList)

	// start multiple scanner goroutines
	const scanWorkers = 5 // TODO: Make this configurable via flags
	for i := 0; i < scanWorkers; i++ {
		go scanHost(scanConfig, addressChannel, successChannel, &wg)
		wg.Add(1)
	}

	// Feed all the addresses into the scan workers
	for {
		nextIPAddr, ok := targetNetwork.NextHostAddress()

		if !ok { // Did we use up all the addresses?
			break
		}
		addressChannel <- nextIPAddr
	}
	close(addressChannel)

	wg.Wait()

	// All scanning finished. Note duration and output results
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	// output results summary here
	fmt.Println("Finished scanning in ", duration.String())
	fmt.Println("Successful Hosts:")
	for _, success := range successList {
		fmt.Println(success.String(true))
	}
}

func scanHost(
	scanConfig ScanConfig,
	addressChannel <-chan netip.Addr,
	successChannel chan<- *http.Response,
	wg *sync.WaitGroup,
) {
	// Create a re-usable HTTP client, one for each function invocation
	client := http.Client{
		Timeout: scanConfig.connectTimeout,
	}

	for targetAddress := range addressChannel {
		fmt.Printf("Scanning %s...\n", targetAddress)
		urlString := fmt.Sprintf("http://%s:%d", targetAddress.String(), scanConfig.port)
		resp, err := client.Get(urlString)

		if err != nil {
			//fmt.Printf("Error trying to connect to %s: %v\n", targetAddress, err)
			continue
		}

		successChannel <- resp
	}
	wg.Done()
}

// processSuccess adds the processes result struct to a slice. As slices are not
// thread safe, this should be a single consumer of the "successful connection" channel.
// Also, we need to remember to call response.Body.Close() when we're done.
func processSuccess(
	responseChannel <-chan *http.Response,
	successList *[]ConnectedResponse,
) {
	for response := range responseChannel {
		var serverString string
		status := response.StatusCode
		serverHeader, ok := response.Header["Server"]
		if ok {
			serverString = serverHeader[0]
		} else {
			// Essentially zero value. Do it like this to allow more specific handling potential
			// in future.
			serverString = ""
		}
		url := response.Request.URL

		// Add any extras we might want
		// TODO: Break this out
		// 401 response? Add WWW-Authorization header
		var extras = make(map[string]string)
		var extrasStack []string
		if status == http.StatusUnauthorized {
			extrasStack = append(extrasStack, "authHeader")
			//extras["WWW-Authenticate"] = response.Header.Get("WWW-Authenticate")
		}

		var finalisedResponse = ConnectedResponse{
			Status:      status,
			Server:      serverString,
			Url:         url.String(),
			Host:        url.Host,
			Extras:      extras,
			extrasStack: extrasStack,
		}

		// If the "extrasStack" is not empty, send it off to the Extras Dispatcher
		// It will consume entries from the stack, and send it back on the processedExtras
		// channel where it will *then* be added to the "success list", with a populated
		// map of "extra" details.
		if len(finalisedResponse.extrasStack) > 0 {
			finalisedResponse = ProcessExtras(finalisedResponse, response)
		}

		*successList = append(
			*successList,
			finalisedResponse,
		)
		// Be sure to close the body
		err := response.Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}
}

func processError(errorChannel <-chan error) {}

func addAuthorizationDetail(request *http.Request, response ConnectedResponse) error { return nil }
