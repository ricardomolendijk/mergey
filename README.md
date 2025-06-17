# MergePlease

A Go CLI tool to pick your latest GitLab Merge Request and send it to Slack via webhook.

## Features

- Fetches your open, non-draft MRs (created by you) from multiple GitLab instances
- Presents a picker for the most recent MRs
- Sends a random message, MR title, and MR URL to Slack
- Configurable via YAML

## Configuration

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

## Usage

1. Build the app:
   ```sh
   make build
   ```
2. Run the app:
   ```sh
   ./mergeplease
   ```

## License

MIT
