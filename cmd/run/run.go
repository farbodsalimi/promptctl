package run

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/farbodsalimi/promptctl/internal/db"
	"github.com/farbodsalimi/promptctl/internal/providers"
	"github.com/farbodsalimi/promptctl/internal/templates"
)

var promptRunCmd = &cobra.Command{
	Use:   "prompt <vault> <name>",
	Short: "Run a prompt with an LLM provider",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		vaultName := args[0]
		promptName := args[1]

		provider, _ := cmd.Flags().GetString("provider")
		model, _ := cmd.Flags().GetString("model")
		vars, _ := cmd.Flags().GetString("vars")
		version, _ := cmd.Flags().GetInt("version")
		temperature, _ := cmd.Flags().GetFloat32("temperature")

		if provider == "" {
			log.Fatal("provider is required (use --provider)")
		}
		if model == "" {
			log.Fatal("model is required (use --model)")
		}

		// Parse variables
		varsMap, err := templates.ParseVars(vars)
		if err != nil {
			log.Fatalf("failed to parse variables: %v", err)
		}

		// Get prompt content
		var content string
		var promptVersionID int

		if version > 0 {
			content, err = db.GetPromptVersionContentByVersion(vaultName, promptName, version)
		} else {
			prompt, err2 := db.GetPromptByName(vaultName, promptName)
			if err2 != nil {
				log.Fatalf("prompt not found: %s/%s", vaultName, promptName)
			}
			promptVersionID, content, err = db.GetPromptVersionContent(prompt.ID)
		}

		if err != nil {
			log.Fatalf("failed to get prompt content: %v", err)
		}

		// Render template
		renderedPrompt, err := templates.RenderTemplate(content, varsMap)
		if err != nil {
			log.Fatalf("failed to render template: %v", err)
		}

		fmt.Printf("Rendered prompt:\n---\n%s\n---\n\n", renderedPrompt)

		// Get LLM provider
		ctx := context.Background()
		router, err := providers.GetProvider(ctx)
		if err != nil {
			log.Fatalf("failed to get provider: %v", err)
		}

		// Get specific LLM provider
		llm, ok := router.Get(provider)
		if !ok {
			log.Fatalf("provider not found: %s", provider)
		}

		// Execute request
		response, err := llm.Complete(renderedPrompt)
		if err != nil {
			log.Fatalf("failed to generate response: %v", err)
		}

		fmt.Printf("Response:\n%s\n", response)

		// Store run in database
		paramsData := map[string]any{
			"provider":    provider,
			"model":       model,
			"temperature": temperature,
			"vars":        varsMap,
		}
		paramsJSON, _ := json.Marshal(paramsData)

		err = db.CreateRun(promptVersionID, provider, string(paramsJSON), response)
		if err != nil {
			log.Printf("warning: failed to store run in database: %v", err)
		}
	},
}

var runListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent runs",
	Run: func(cmd *cobra.Command, args []string) {
		promptName, _ := cmd.Flags().GetString("prompt")
		vaultName, _ := cmd.Flags().GetString("vault")

		runs, err := db.GetRuns(vaultName, promptName)
		if err != nil {
			log.Fatalf("failed to list runs: %v", err)
		}

		if len(runs) == 0 {
			fmt.Println("No runs found")
			return
		}

		fmt.Println("Recent runs:")
		for _, run := range runs {
			fmt.Printf("  Run %d: %s with %s (created: %s)\n",
				run.ID, run.PromptName, run.Provider, run.Created)
		}
	},
}

var runShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show details of a specific run",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("invalid run ID: %s", idStr)
		}

		run, err := db.GetRunByID(id)
		if err != nil {
			log.Fatalf("run not found: %d", id)
		}

		fmt.Printf("Run %d:\n", run.ID)
		fmt.Printf("  Prompt: %s\n", run.PromptName)
		fmt.Printf("  Provider: %s\n", run.Provider)
		fmt.Printf("  Created: %s\n", run.Created)
		fmt.Printf("  Parameters:\n%s\n", run.Params)
		fmt.Printf("  Response:\n---\n%s\n---\n", run.Response)
	},
}

func NewRootCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run prompts and manage run history",
		Long:  `Execute prompts with LLM providers and view run history.`,
	}

	runCmd.AddCommand(promptRunCmd)
	runCmd.AddCommand(runListCmd)
	runCmd.AddCommand(runShowCmd)

	promptRunCmd.Flags().StringP("provider", "p", "", "LLM provider to use (openai, anthropic, google)")
	promptRunCmd.Flags().StringP("model", "m", "", "Model name (e.g., gpt-4, claude-3-sonnet, gemini-pro)")
	promptRunCmd.Flags().StringP("vars", "", "", "Template variables as JSON object or key=value pairs (e.g., '{\"name\":\"John\"}' or 'name=John,age=30')")
	promptRunCmd.Flags().IntP("version", "v", 0, "Specific prompt version to use (default: latest version)")
	promptRunCmd.Flags().Float32P("temperature", "t", 0.7, "Sampling temperature for response generation (0.0-2.0, higher = more creative)")

	runListCmd.Flags().StringP("prompt", "p", "", "Filter results by prompt name")
	runListCmd.Flags().StringP("vault", "v", "", "Filter results by vault name")

	return runCmd
}
