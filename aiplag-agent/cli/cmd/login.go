package cmd

import (
	"aiplag-agent/common/config"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	//BackendBaseURL = "http://localhost:8080"
	BackendBaseURL       = "https://plaggy.xyz"
	MagicRequestEndpoint = BackendBaseURL + "/api/v1/auth/magic-request"
	MagicStatusEndpoint  = BackendBaseURL + "/api/v1/auth/magic-status"
)

func requestMagicLink(email string) (string, error) {
	jsonData, err := json.Marshal(map[string]string{"email": email})
	if err != nil {
		return "", err
	}

	resp, err := http.Post(MagicRequestEndpoint, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status %s", resp.Status)
	}

	var result struct {
		MagicId string `json:"magicId"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.MagicId, nil
}

func checkLoginStatus(magicID string, email string) (bool, string, error) {
	u, err := url.Parse(MagicStatusEndpoint)
	if err != nil {
		return false, "", err
	}

	q := u.Query()
	q.Set("magic", magicID)
	u.RawQuery = q.Encode()

	bodyData := map[string]string{"email": email}
	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return false, "", err
	}

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(bodyBytes))
	if err != nil {
		return false, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("server returned status %s", resp.Status)
	}

	var status struct {
		Authenticated bool   `json:"authenticated"`
		Token         string `json:"token,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return false, "", err
	}

	return status.Authenticated, status.Token, nil
}

var loginCmd = &cobra.Command{
	Use:   "login [email]",
	Short: "Manages user login",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var email string

		if len(args) == 0 || args[0] == "" {
			fmt.Print("Enter your email: ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			email = strings.TrimSpace(input)
		} else {
			email = args[0]
		}

		if email == "" {
			fmt.Println("Email field is empty! Aborting.")
			return
		}

		magicId, err := requestMagicLink(email)
		if err != nil {
			fmt.Println("Failed to request magic link:", err)
			return
		}

		viper.Set("session.email", email)
		viper.Set("session.token", "")

		cfgPath := config.ConfigPath()
		err = viper.WriteConfigAs(cfgPath)
		if err != nil {
			err = viper.SafeWriteConfigAs(cfgPath)
			if err != nil {
				fmt.Println("Failed to save config:", err)
			}
		}

		fmt.Println("Magic link sent! Please check your email and click the link to complete login.")

		waitForLogin(magicId, email)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func waitForLogin(magicId, email string) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		authenticated, token, err := checkLoginStatus(magicId, email)
		if err != nil {
			log.Println("Error checking login status:", err)
			return
		}

		if authenticated {
			fmt.Println("Login successful! You can now run commands.")

			viper.Set("session.token", token)
			cfgPath := config.ConfigPath()
			_ = viper.WriteConfigAs(cfgPath)

			break
		}
	}
}
