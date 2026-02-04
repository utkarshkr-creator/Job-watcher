# Job Watcher

Job Watcher is a powerful, automated tool designed to help developers find job opportunities efficiently. It aggregates job listings from multiple sources, filters them based on your preferences, scores them using AI (optional), and delivers real-time notifications directly to your Telegram.

## Features

-   **Multi-Source Aggregation**: Fetches jobs from:
    -   Indeed
    -   LinkedIn
    -   Naukri
    -   Instahyre
    -   RemoteOK
    -   Razorpay
    -   Wellfound
    -   YCombinator (Work at a Startup)
    -   Hacker News ("Who's Hiring")
    -   Reddit (r/cscareerquestions, r/forhire)
    -   Triplebyte
    -   Company Career Pages (90+ tech companies)
    -   Shared Lists (Google Sheets / GitHub Tables)
-   **High Performance**: Fetches from all sources in parallel for maximum speed.
-   **Smart Filtering**:
    -   Keyword matching (titles, skills)
    -   Experience level filtering
    -   Location filtering (Remote/India specific)
    -   Exclusion keywords (avoid senior/manager roles if desired)
-   **AI Scoring**: Integrates with local LLMs (via Ollama) to score jobs based on your resume and preferences.
-   **Real-time Alerts**: Sends instant notifications to Telegram with job details and links.
-   **Deduplication**: Tracks seen jobs to prevent duplicate alerts.
-   **Job Enrichment**: Can fetch additional recruiter details for high-scoring jobs.

## Prerequisites

-   **Go**: Version 1.21 or higher.
-   **Telegram Bot**: You need a Bot Token and Chat ID to receive alerts.
-   **Ollama** (Optional): Only required if you want AI-powered job scoring and filtering.

## Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/yourusername/job-watcher.git
    cd job-watcher
    ```

2.  **Install dependencies**:
    ```bash
    go mod download
    ```

3.  **Setup Configuration**:
    -   The project uses `config.yaml` to define search parameters. Modify it to match your skills and preferences.
    -   Create a `.env` file for your secrets (see configuration below).

### AI Configuration (Optional)

The project uses Ollama to run a local LLM (like Llama 3) to score jobs based on your resume. This helps prioritize the most relevant listings.

**If you don't want to use AI:**
1.  Open `config.yaml`.
2.  Set `enabled: false` under the `ai` section:
    ```yaml
    ai:
      enabled: false
    ```
3.  You can skip installing Ollama.

**If you DO want to use AI:**
1.  Install [Ollama](https://ollama.com/).
2.  Pull the model specified in `config.yaml` (default `llama3.1:latest`):
    ```bash
    ollama pull llama3.1:latest
    ```
3.  Ensure `ai.enabled` is set to `true` in `config.yaml`.

### Configuration

### 1. Telegram Bot Setup

To receive notifications, you need to set up a Telegram bot:

1.  **Create a Bot**:
    -   Open Telegram and search for **@BotFather**.
    -   Send the command `/newbot`.
    -   Follow the instructions to name your bot.
    -   **BotFather** will give you a **HTTP API Token**. This is your `TG_TOKEN`.

2.  **Get Chat ID**:
    -   Start a chat with your new bot and send a message (e.g., "Hello").
    -   Visit this URL in your browser: `https://api.telegram.org/bot<YOUR_TG_TOKEN>/getUpdates`
    -   Look for the `"chat":{"id":12345678,...}` part in the JSON response.
    -   The number (e.g., `12345678`) is your `TG_CHAT`.

### 2. Environment Variables

Create a `.env` file in the root directory (you can use `.env.example` as a reference):

```bash
TG_TOKEN=your_telegram_bot_token
TG_CHAT=your_telegram_chat_id
```

### 3. config.yaml

The `config.yaml` file controls the behavior of the scraper. For a detailed guide on how to personalize this for your specific needs (keywords, experience level, AI matching), please read **[CUSTOMIZATION.md](CUSTOMIZATION.md)**.

Key sections include:

-   `keywords`: List of technologies and roles you are interested in (e.g., "python", "react", "SDE").
-   `locations`: Target locations (e.g., "remote", "bangalore").
-   `exclude_keywords`: keywords to filter out (e.g., "senior", "lead").
-   `max_experience_years`: Maximum experience required for the job.
-   `sources`: Enable/disable specific job sources.
-   `ai`: Configuration for AI scoring (enabled/disabled, model name, threshold).

## Usage

1.  **Run the watcher**:
    ```bash
    go run .
    ```
    Or build and run:
    ```bash
    go build -o job-watcher
    ./job-watcher
    ```

2.  **AI Feature**:
    If AI is enabled in `config.yaml`, make sure Ollama is running:
    ```bash
    ollama serve
    ```
    And ensure you have the configured model pulled (default is `llama3.1:latest`):
    ```bash
    ollama pull llama3.1:latest
    ```

## Contributing

We welcome contributions! Please see **[CONTRIBUTING.md](CONTRIBUTING.md)** for details on how to submit pull requests, report issues, and follows our code style.

## License

[MIT](LICENSE)
