# MergePlease

🚀 **MergePlease** is a Go CLI tool to pick your latest GitLab Merge Request and send it to Slack with style!

---

## ✨ Features

- 🔍 Fetches your open, non-draft MRs (created by you) from multiple GitLab instances
- 🖱️ Interactive picker for your most recent MRs
- 🎲 Sends a random message, MR title, and MR URL to Slack
- 🛠️ Configurable via YAML in your home directory
- 💬 Customizable Slack messages

---

## 📦 Installation

1. **Build the app:**
   ```sh
   make build
   ```
2. **Copy the binary (optional):**
   ```sh
   cp mergeplease ~/bin/mergeplease
   ```
3. **Ensure `~/.merge/config.yaml` exists and is configured (see below).**

---

## ⚙️ Configuration

Create a config file at `~/.merge/config.yaml`:

```yaml
gitlab:
  - api: https://gitlab.com/api/v4
    token: your_gitlab1_private_token
  - api: https://your.other.gitlab.instance/api/v4
    token: your_gitlab2_private_token
slack_webhook: "https://hooks.slack.com/triggers/..."
slack_messages:
  - Merge please!
  - Ready for review!
  - Can you take a look?
  - PTAL
  - Time to merge!
mr_picker_count: 5
```

- `gitlab`: List of GitLab API endpoints and tokens
- `slack_webhook`: Your Slack webhook URL
- `slack_messages`: List of possible random messages
- `mr_picker_count`: How many MRs to show in the picker

---

## 🧑‍💻 Usage

```sh
./mergeplease
```

- Pick a merge request from the interactive list
- Confirm to send to Slack

---

## 📝 Slack Webhook Requirements

- The Slack webhook must accept a JSON payload with a `content` key, for example:
  ```json
  { "content": "Your message here" }
  ```
- The webhook URL should be set in the `slack_webhook` field of your config.

---

## 📄 Example Slack Message

```
:gitlab: Merge please!
Add new login page
https://gitlab.com/yourproject/merge_requests/123
```

---

## 🛡️ License

MIT

---

## ❤️ Contributing

Pull requests and issues are welcome! If you have ideas for more features or improvements, open an issue or PR.

---

## 🤝 Acknowledgements

- Inspired by modern CLI tools and the need for fast, beautiful merge request reminders!
