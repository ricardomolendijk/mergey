package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v3"
)

type MergeRequest struct {
	Title     string    `json:"title"`
	WebURL    string    `json:"web_url"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GitLabConfig struct {
	APIURL string
	Token  string
}

type Config struct {
	Gitlab []struct {
		API   string `yaml:"api"`
		Token string `yaml:"token"`
	} `yaml:"gitlab"`
	Slack struct {
		Webhook  string   `yaml:"webhook"`
		Messages []string `yaml:"messages"`
	} `yaml:"slack"`
	MRPickerCount int `yaml:"mr_picker_count"`
}

func fetchLastMR(cfg GitLabConfig) (*MergeRequest, error) {
	req, err := http.NewRequest("GET", cfg.APIURL+"/merge_requests?scope=created_by_me&order_by=updated_at&sort=desc&per_page=1", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", cfg.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var mrs []MergeRequest
	if err := json.NewDecoder(resp.Body).Decode(&mrs); err != nil {
		return nil, err
	}
	if len(mrs) == 0 {
		return nil, fmt.Errorf("no merge requests found")
	}
	return &mrs[0], nil
}

func fetchRecentMRs(cfg GitLabConfig, count int) ([]MergeRequest, error) {
	url := cfg.APIURL + fmt.Sprintf("/merge_requests?scope=created_by_me&state=opened&order_by=updated_at&sort=desc&per_page=%d", count*2) // fetch more to allow for draft filtering
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", cfg.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var mrs []MergeRequest
	if err := json.NewDecoder(resp.Body).Decode(&mrs); err != nil {
		return nil, err
	}
	// Filter out drafts (title starts with 'Draft:' or 'WIP:' or 'Draft ' or 'WIP ')
	var filtered []MergeRequest
	for _, mr := range mrs {
		title := strings.ToLower(mr.Title)
		if strings.HasPrefix(title, "draft:") || strings.HasPrefix(title, "wip:") || strings.HasPrefix(title, "draft ") || strings.HasPrefix(title, "wip ") {
			continue
		}
		filtered = append(filtered, mr)
		if len(filtered) >= count {
			break
		}
	}
	return filtered, nil
}

func askConfirmation(mr *MergeRequest) bool {
	fmt.Printf("Latest MR: %s\nURL: %s\nUpdated: %s\nConfirm send to Slack? (y/n): ", mr.Title, mr.WebURL, mr.UpdatedAt.Format(time.RFC822))
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return strings.ToLower(scanner.Text()) == "y"
}

func loadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func getRandomMessageYAML(cfg *Config) string {
	if len(cfg.Slack.Messages) == 0 {
		return "Merge please!"
	}
	rand.Seed(time.Now().UnixNano())
	return strings.TrimSpace(cfg.Slack.Messages[rand.Intn(len(cfg.Slack.Messages))])
}

func sendToSlack(webhookURL string, mr *MergeRequest, cfg *Config) error {
	msg := map[string]string{
		"content": fmt.Sprintf(":gitlab: %s\n%s\n%s", getRandomMessageYAML(cfg), mr.Title, mr.WebURL),
	}
	b, _ := json.Marshal(msg)
	resp, err := http.Post(webhookURL, "application/json", strings.NewReader(string(b)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("slack webhook error: %s", resp.Status)
	}
	return nil
}

func pickMR(mrs []MergeRequest) *MergeRequest {
	items := make([]string, len(mrs))
	for i, mr := range mrs {
		items[i] = fmt.Sprintf("%s (Updated: %s)", mr.Title, mr.UpdatedAt.Format(time.RFC822))
	}
	prompt := promptui.Select{
		Label: "Select a Merge Request",
		Items: items,
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return nil
	}
	return &mrs[idx]
}

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Could not determine current user:", err)
		return
	}
	configPath := filepath.Join(usr.HomeDir, ".merge", "config.yaml")
	cfg, err := loadConfig(configPath)
	if err != nil {
		fmt.Println("Failed to load config at", configPath, ":", err)
		return
	}
	pickerCount := cfg.MRPickerCount
	if pickerCount <= 0 {
		pickerCount = 5
	}
	var allMRs []MergeRequest
	for _, g := range cfg.Gitlab {
		mrs, err := fetchRecentMRs(GitLabConfig{APIURL: g.API, Token: g.Token}, 1000)
		if err == nil {
			allMRs = append(allMRs, mrs...)
		}
	}
	if len(allMRs) == 0 {
		fmt.Println("No merge requests found on either GitLab instance. Check your config.yaml and network access.")
		return
	}
	// Sort allMRs by UpdatedAt desc
	for i := 0; i < len(allMRs)-1; i++ {
		for j := i + 1; j < len(allMRs); j++ {
			if allMRs[j].UpdatedAt.After(allMRs[i].UpdatedAt) {
				allMRs[i], allMRs[j] = allMRs[j], allMRs[i]
			}
		}
	}
	max := pickerCount
	if len(allMRs) < pickerCount {
		max = len(allMRs)
	}
	picked := pickMR(allMRs[:max])
	if picked == nil {
		fmt.Println("Aborted.")
		return
	}
	if askConfirmation(picked) {
		if err := sendToSlack(cfg.Slack.Webhook, picked, cfg); err != nil {
			fmt.Println("Failed to send to Slack:", err)
		} else {
			fmt.Println("Sent to Slack!")
		}
	} else {
		fmt.Println("Aborted.")
	}
}
