/*
	export ADYEN_CLIENT_TOKEN="YOUR_ADYEN_ENCRYPTED_URL"
	export ADYEN_USERNAME="YOUR_ADYEN_API_USERNAME"
	export ADYEN_PASSWORD="YOUR_API_PASSWORD"
	export ADYEN_ACCOUNT="YOUR_MERCHANT_ACCOUNT"
*/

package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	adyen "github.com/zhutik/adyen-api-go"
)

// TempateConfig for HTML template
type TempateConfig struct {
	EncURL string
	Time   string
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

/**
 * Show Adyen Payment form
 */
func showForm(w http.ResponseWriter, r *http.Request) {
	adyen := adyen.New(
		os.Getenv("ADYEN_USERNAME"),
		os.Getenv("ADYEN_PASSWORD"),
		os.Getenv("ADYEN_CLIENT_TOKEN"),
		os.Getenv("ADYEN_ACCOUNT"),
	)
	now := time.Now()

	config := TempateConfig{
		EncURL: adyen.ClientURL(),
		Time:   now.Format(time.RFC3339),
	}

	t := template.Must(template.ParseFiles("index.html"))
	t.Execute(w, config)
}

/**
 * Handle post request and perform payment authosization
 */
func performPayment(w http.ResponseWriter, r *http.Request) {
	adyen := adyen.New(
		os.Getenv("ADYEN_USERNAME"),
		os.Getenv("ADYEN_PASSWORD"),
		os.Getenv("ADYEN_CLIENT_TOKEN"),
		os.Getenv("ADYEN_ACCOUNT"),
	)

	r.ParseForm()

	rand.Seed(time.Now().UTC().UnixNano())

	g, err := adyen.Authorise().Payment(
		r.Form.Get("adyen-encrypted-data"),
		"DE-100"+randomString(6),
		1000,
	)

	if err == nil {
		fmt.Fprintf(w, "<h1>Success!</h1><code><pre>"+g.AuthCode+" "+g.PspReference+"</pre></code>")
	} else {
		fmt.Fprintf(w, "<h1>Something went wrong: "+err.Error()+"</h1>")
	}
}

func performCapture(w http.ResponseWriter, r *http.Request) {
	adyen := adyen.New(
		os.Getenv("ADYEN_USERNAME"),
		os.Getenv("ADYEN_PASSWORD"),
		os.Getenv("ADYEN_CLIENT_TOKEN"),
		os.Getenv("ADYEN_ACCOUNT"),
	)

	r.ParseForm()

	amount, err := strconv.ParseFloat(r.Form.Get("amount"), 32)

	if err != nil {
		fmt.Fprintf(w, "<h1>Failed! Can not convert amount to float</h1>")
		return
	}

	g, err := adyen.Capture().Payment(
		r.Form.Get("original-reference"),
		r.Form.Get("reference"),
		float32(amount),
	)

	if err == nil {
		fmt.Fprintf(w, "<h1>Success!</h1><code><pre>"+g.PspReference+" "+g.Response+"</pre></code>")
	} else {
		fmt.Fprintf(w, "<h1>Something went wrong: "+err.Error()+"</h1>")
	}
}

func main() {
	http.HandleFunc("/", showForm)
	http.HandleFunc("/perform_payment", performPayment)
	http.HandleFunc("/perform_capture", performCapture)
	http.ListenAndServe(":8080", nil)
}
