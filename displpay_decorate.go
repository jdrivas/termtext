package termtext

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jdrivas/vconfig"
	"github.com/juju/ansiterm"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
)

// We want to decorate List and Describe with some context dependent
// display of the HTTP response and any errors.
//
// List, Describe are created as methods on objects that are mirror to the
// jupyterhub objects. We want to use method displatch to deal with the
// different kinds of objects and sicne we can't add functions outside of
// a package, we'll create mirror types.
// e.g.   type Group jh.Group
//
// To use these you cast an object from jh to the mirror type, then
// call the function List() or Describe() with the result from the JH function:
//
//   groups, resp, err := jh.GetGroups()
//   List(Groups(groups), resp, err)
//
// The method's are not called directly to alllow decorattion of the output with resp and
// error display dependent  on: error condition, verbose vs. debug etc.
// Also, we want the same display output in those cases where the is no object returned
// from the jh functions (e.g. jh.GetGroups).
//
// Display() is a method that is called when there is only an http.Resoionse and error returned.
// This is typical, for example, on Delete calls.

//
// To do this, List(...) and Describe(...) , Display(...) all call render() which sets up
// a decorotor pipeline as needed.

// There are two basic object display functions: List and Describe.
// Not every data object supports both. Generally, they
// all spport List, and some also support Describe.

// Listable suppots List()
type Listable interface {
	List()
}

// Describable supports Describe()
type Describable interface {
	Describe()
}

// List and Describe display their objects by calling render, but
// first checking that an object is there. If not they send along
// an empty function for the descoroator to call.
func List(d Listable, resp *http.Response, err error) {
	renderer := func() {}
	if d != nil {
		renderer = d.List
	}
	render(renderer, resp, err)
}

// Describe provides detailed output on the object.
func Describe(d Describable, resp *http.Response, err error) {
	renderer := func() {}
	if d != nil {
		renderer = d.Describe
	}
	render(renderer, resp, err)
}

/*
// Display dispolays only the resp and error through the normal pipeline
func Display(resp *http.Response, err error) {
	render(func() {}, resp, err)
}

// DisplayF calls the displayRenderer function as part of the standard render pipeline.
// This is useful for printing out status information bracketed by the usual
// verbose/debugt etc. influenced response and error output from the normal pipeline.
func DisplayF(displayRenderer func(), resp *http.Response, err error) {
	render(displayRenderer, resp, err)

}


func displayServerStartedF(started bool, resp *http.Response, err error) func() {
	return (func() {
		result := Success("started")
		if started == false {
			result = Success("requested")
			if err != nil {
				result = Fail("probably not started")
			}
		}
		fmt.Printf("%s %s\n", Title("Server"), result)
	})
}

func displpayServerStopedF(stopped bool, resp *http.Response, err error) func() {
	return (func() {
		result := Success("stopped")
		if stopped == false {
			result = Success("requested")
			if err != nil {
				result = Fail("probably not stopped")
			}
		}
		fmt.Printf("%s %s\n", Title("Server"), result)
	})
}
*/

// private API
func render(renderer func(), resp *http.Response, err error) {

	switch {
	case vconfig.Debug():
		httpDecorate((errorDecorate(renderer, err)), resp)()
	case vconfig.Verbose():
		shortHTTPDecorate((errorDecorate(renderer, err)), resp)()
	case viper.GetBool(JSONDisplayKey):
		if resp != nil {
			jsonDisplay(resp, err)
			return
		}
		fallthrough
	default:
		if err == nil {
			errorDecorate(renderer, err)()
		} else {
			errorHTTPDecorate((errorDecorate(renderer, err)), resp)()
		}
	}
}

// The Decorators are built as pre function call. That is print, then
// call your argument. So this goes frist to last, with the list.
// Thus httpDecorate(errorDecorate(d.List)) will first print
// the http Response, then the error message, then the List().

// HTTPDisplay - This is for the HTTP direct commands.
func HTTPDisplay(resp *http.Response, err error) {
	if viper.GetBool(JSONDisplayKey) {
		jsonDisplay(resp, err)
	} else {
		httpDecorate(errorDecorate(func() {}, err), resp)()
	}
}

// Raw JSON display, no decoration.
func jsonDisplay(resp *http.Response, err error) {

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err == nil && resp.StatusCode != http.StatusNoContent {
		prettyJSON := bytes.Buffer{}
		err := json.Indent(&prettyJSON, body, "", "  ")
		if err == nil {
			prettyJSON.WriteTo(os.Stdout)
			fmt.Println()
		} else {
			fmt.Printf("%s", string(body))
		}
	} else {
		fmt.Printf("Body read error: %v\n", err)
	}
}

func errorDecorate(f func(), err error) func() {
	return (func() {
		if err != nil {
			fmt.Printf(fmt.Sprintf("%s\n", Error(err)))
		}

		f()
	})
}

// What to say if there is no response.
func nilResp() {
	fmt.Printf("Nil HTTP Response.\n")
}

// One liner update on the response.
func shortHTTPDecorate(f func(), resp *http.Response) func() {
	return (func() {
		if resp == nil {
			nilResp()
		} else {
			fmt.Printf("%s %s\n", Title("HTTP Response: "), httpStatusFunc(resp.StatusCode)("%s", resp.Status))
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil && resp.StatusCode != http.StatusNoContent {
				prettyPrintBody(body)
			}
		}

		f()
	})
}

func errorHTTPDecorate(f func(), resp *http.Response) func() {
	return (func() {
		if resp == nil {
			nilResp()
		} else {
			fmt.Printf("%s %s\n", Title("HTTP Response: "), httpStatusFunc(resp.StatusCode)("%s", resp.Status))
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil && resp.StatusCode != http.StatusNoContent {
				m, err := getMessage(body)
				if err == nil && m.Message != "" {
					fmt.Printf("%s %s\n", Title("Message:"), Alert(m.Message))
				}
			}
		}

		f()
	})
}

// Message is a simple struct to pull out JSON that is often
// embedded in error returns.
type Message struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Tabled based HTTP reseponse with headers.
func httpDecorate(f func(), resp *http.Response) func() {
	return (func() {
		if resp != nil {
			w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
			fmt.Fprintf(w, "%s\n", Title("Status\tLength\tEncoding\tUncompressed"))
			fmt.Fprintf(w, "%s\t%s\n",
				httpStatusFunc(resp.StatusCode)("%s", resp.Status),
				Text("%d\t%#v\t%t", resp.ContentLength, resp.TransferEncoding, resp.Uncompressed))
			w.Flush()

			// Headers
			w = ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
			fmt.Fprintf(w, "%s\n", Title("Header\tValue"))
			for k, v := range resp.Header {
				fmt.Fprintf(w, "%s\n", Text("%s\t%s", k, v))
			}
			w.Flush()

			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil && resp.StatusCode != http.StatusNoContent {
				prettyPrintBody(body)
			} else {
				fmt.Printf("%s %s\n", Title("Body Read Error:"), Text("%v", err))
			}
		} else {
			nilResp()
		}

		f()

	})
}

func getMessage(jsonString []byte) (m Message, err error) {
	err = json.Unmarshal(jsonString, &m)
	return m, err
}

// Assume it's json and try to pretty print.
func prettyPrintBody(body []byte) {
	prettyJSON := bytes.Buffer{}
	err := json.Indent(&prettyJSON, body, "", "  ")
	if err == nil {

		// Yes, this is totally gratuitous.
		m, err := getMessage(body)
		if err == nil && m.Message != "" {
			fmt.Printf("%s %s\n", Title("Message:"), Alert(m.Message))
		}

		// Print out the pretty body
		fmt.Printf("%s\n%s\n", Title("RESP JSON Body:"), Text("%s", string(prettyJSON.Bytes())))
	} else {
		// Much of the time, we've just read past the body in properly handling the response.
		if len(body) > 0 {
			fmt.Printf("%s %s \n", Title("JSON indenting error:"), Fail("%v", err))
			fmt.Printf("%s\n%s\n", Title("So here is the RESP body:"), Text("%s", string(body)))
		}
	}

}

func httpStatusFunc(httpStatus int) (f ColorSprintfFunc) {
	switch {
	case httpStatus < 300:
		f = Success
	case httpStatus < 400:
		f = Warn
	default:
		f = Fail
	}
	return f
}

func cmdError(e error) {
	fmt.Printf("Error: %s\n", Fail(e.Error()))
}
