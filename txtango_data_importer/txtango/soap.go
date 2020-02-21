package txtango

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

//soapCall generate a request given a request and a template and sends it
func soapCall(params interface{}, tmplName, tmplRaw string) ([]byte, error) {
	//construct the request using a template
	tmpl := template.Must(template.New(tmplName).Parse(tmplRaw))
	tmpl = template.Must(tmpl.Parse(loginTemplate))

	doc := &bytes.Buffer{}
	//fill in template values with actual values
	err := tmpl.Execute(doc, params)
	if err != nil {
		return nil, errors.Wrap(err, "Error while filling template")
	}

	if err := xml.Unmarshal([]byte(doc.String()), new(interface{})); err != nil {
		return nil, errors.Wrap(err, "There is an error in xml request. Please dig in the code for")
	}

	//build request
	httpRequest, err := http.NewRequest(
		http.MethodPost,
		os.Getenv("TX_HOST"),
		bytes.NewBuffer([]byte(doc.String())))
	if err != nil {
		return nil, errors.Wrap(err, "Error while generating request")
	}
	//add request header
	httpRequest.Header.Add("Content-Type", "text/xml; charset=utf-8")

	//send request
	client := &http.Client{}
	response, err := client.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	//read response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	//return response
	return body, nil
}
