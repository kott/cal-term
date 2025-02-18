package auth

import (
  "context"
  "fmt"
  "net/http"
  "os/exec"
  "runtime"
  
  "golang.org/x/oauth2"
)

func redirectServer(port int, requestState string, config *oauth2.Config, ch chan<- *oauth2.Token) *http.Server {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")

		if code == "" {
			fmt.Fprintf(w, "No code found")
			ch <- nil
			return
		}

		token, err := config.Exchange(context.Background(), code)
		if err != nil {
			fmt.Fprintf(w, "Error exchanging code for token: %v", err)
			ch <- nil
			return
		}

		if state != requestState {
			fmt.Fprintf(w, "State does not match original request: %s", state)
			ch <- nil
			return
		}

		fmt.Fprintf(w, "Success! You can close this window.")
		ch <- token
	})

	return server
}

func openUrl(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}

func waitForKeyPress(msg string) {
	fmt.Println(msg)
	var input string
	fmt.Scanln(&input)
}

