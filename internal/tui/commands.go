package tui

import (
	"fmt"
	"strings"
)

// commandResult is the outcome of executing one command line.
type commandResult struct {
	output  string
	isError bool
	clear   bool
}

// commandDef mirrors CommandDefinition from divit.qezta.com's
// src/lib/data/terminal-commands.ts, ported so `ssh divit@qezta.com` gives
// the same terminal experience as the browser widget on the website.
type commandDef struct {
	description string
	execute     func(args []string) commandResult
}

// identity is which persona a session sees, chosen by the username the
// client connects with. There's only one SSH listener on qezta.com, so
// routing happens on username rather than running separate servers.
type identity int

const (
	// identityDivit is Divit's personal portfolio — `ssh divit@qezta.com`.
	identityDivit identity = iota
	// identityQezta is the generic org/company landing — anything else,
	// e.g. `ssh qezta@qezta.com` or `ssh hi@qezta.com`.
	identityQezta
)

// identityFor maps a connecting username to a persona. Only "divit" (and
// the shorthand "div") get the personal portfolio; every other username
// falls back to the org identity, so callers don't need to guess a single
// "correct" alias like qezta/hi/hello/guest.
func identityFor(username string) identity {
	switch strings.ToLower(strings.TrimSpace(username)) {
	case "divit", "div":
		return identityDivit
	default:
		return identityQezta
	}
}

type project struct {
	name, category, stars, desc string
}

var projects = []project{
	{"HybridTransformer-MFIF", "ml", "*", "Multi-focus image fusion transformer"},
	{"OS-nixCfg", "infrastructure", "*** 14", "Multi-platform NixOS configuration"},
	{"hs-faust", "functional", "* 5", "Haskell DSL for audio processing"},
	{"CARLA-Autonomous-Driving", "ml", "*** 18", "Autonomous vehicle simulation"},
	{"DocAssist-LLM", "ml", "* 8", "RAG-enhanced LLM from scratch"},
	{"kanata-service", "infrastructure", "* 11", "macOS launchctl service manager"},
}

var skillCategories = []string{
	"Functional Programming",
	"Infrastructure",
	"AI/ML",
	"Data Engineering",
	"DevOps",
	"Languages",
}

var skills = map[string][]string{
	"Functional Programming": {"Haskell", "Nix", "FFI", "DSL Design", "Type Systems"},
	"Infrastructure":         {"NixOS", "Docker", "Kubernetes", "AWS", "CI/CD", "Terraform"},
	"AI/ML":                  {"PyTorch", "Transformers", "Computer Vision", "RAG", "HuggingFace"},
	"Data Engineering":       {"Spark", "AWS Glue", "Hadoop", "Kafka", "PowerBI"},
	"DevOps":                 {"GitHub Actions", "Jenkins", "Hydra", "Nix Binary Caching"},
	"Languages":              {"Python", "TypeScript", "Haskell", "Go", "Swift", "Bash"},
}

// orgRepo mirrors one entry of terranix/orgs/Qezta.nix in the ghOrg-terraform
// repo, which is the actual source of truth for what's public under
// github.com/Qezta.
type orgRepo struct {
	name, desc string
}

var orgRepos = []orgRepo{
	{"divit", "divit.qezta.com — personal portfolio w/ Gemini-powered RAG chatbot"},
	{"qezta-website", "Personal site + Qezta company site development base"},
	{"site", "qezta.com — marketing site for Qezta (under development)"},
	{"ui", "Shared Svelte component library & design system"},
	{"infra", "Cloudflare DNS/domain management (Terranix)"},
	{"ssh", "This terminal — ssh divit@qezta.com or ssh qezta@qezta.com"},
}

// multiWordCommands are matched against the full trimmed input line before
// falling back to first-token dispatch, so easter eggs like "sudo rm -rf /"
// actually fire instead of being shadowed by tokenization. Shared across
// identities.
var multiWordCommands = map[string]commandDef{
	"sudo rm -rf /": {
		description: "DO NOT RUN THIS",
		execute: func(_ []string) commandResult {
			return commandResult{output: `Nice try!

But this is an SSH terminal, not a real shell.
Your system is safe... for now.

Fun fact: On NixOS, even if you did run this, the system could
be rebuilt from configuration. That's the power of declarative
infrastructure!`}
		},
	},
}

// sharedCommands work the same regardless of identity.
var sharedCommands = map[string]commandDef{
	"clear": {
		description: "Clear terminal",
		execute: func(_ []string) commandResult {
			return commandResult{clear: true}
		},
	},
	"nix-build": {
		description: "Build profile (Easter egg)",
		execute: func(_ []string) commandResult {
			return commandResult{output: `these 42 derivations will be built:
  /nix/store/...-qezta-profile-0.1.0.drv
  /nix/store/...-expertise-nix-infrastructure.drv
  /nix/store/...-expertise-functional-programming.drv

building '/nix/store/...-qezta-profile-0.1.0.drv'...
unpacking sources
patching sources
configuring
building
  [1/42] Compiling Haskell skills... OK
  [2/42] Building Nix expertise... OK
  [3/42] Training ML models... OK
  [4/42] Optimizing infrastructure... OK
  ...
  [42/42] Final assembly... OK

build time: 0.5s
output: /nix/store/...-qezta-profile-0.1.0

Success! Profile built successfully.`}
		},
	},
	"cowsay": {
		description: "Make the cow say something",
		execute: func(args []string) commandResult {
			text := strings.Join(args, " ")
			if text == "" {
				text = "Hello from Qezta!"
			}
			line := strings.Repeat("_", len(text)+2)
			return commandResult{output: fmt.Sprintf(` %s
< %s >
 %s
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||`, line, text, line)}
		},
	},
}

var personalCommands map[string]commandDef

var orgCommands map[string]commandDef

func init() {
	personalCommands = merge(sharedCommands, map[string]commandDef{
		"help": {
			description: "Show available commands",
			execute: func(_ []string) commandResult {
				return commandResult{output: `Available commands:

  help              Show this help message
  about             Display information about Divit
  projects          List all projects
  projects --ml     Filter projects by category (ml|infrastructure|functional)
  skills            Show technical expertise
  skills --expert   Show only expert-level skills
  contact           Display contact information
  experience        Show work experience
  resume            View CV
  clear             Clear terminal
  nix-build         Build profile (Easter egg)
  sudo rm -rf /     Don't even try this
  cowsay <text>     Make the cow say something

Type any command to get started!`}
			},
		},
		"about": {
			description: "Display information about Divit",
			execute: func(_ []string) commandResult {
				return commandResult{output: `Divit Mittal - Infrastructure Engineer & ML Researcher

B.Tech Data Science @ Manipal University Jaipur

Specializing in:
• Functional Programming (Haskell, Nix)
• Infrastructure as Code (NixOS, 100+ modules)
• Computer Vision Transformers
• Cloud Data Engineering (AWS, Spark)

Focus areas: declarative infrastructure, type-safe systems,
and production ML deployment.`}
			},
		},
		"projects": {
			description: "List projects",
			execute: func(args []string) commandResult {
				filter := ""
				if len(args) > 0 {
					filter = strings.TrimPrefix(args[0], "--")
				}

				filtered := projects
				if filter == "ml" || filter == "infrastructure" || filter == "functional" {
					filtered = nil
					for _, p := range projects {
						if p.category == filter {
							filtered = append(filtered, p)
						}
					}
				} else if filter != "" {
					return commandResult{
						output:  fmt.Sprintf("No projects found for category: %s\nValid categories: ml, infrastructure, functional", filter),
						isError: true,
					}
				}

				var b strings.Builder
				for i, p := range filtered {
					if i > 0 {
						b.WriteString("\n")
					}
					b.WriteString(fmt.Sprintf("%-12s %-30s %s", p.stars, p.name, p.desc))
				}

				header := "Projects"
				if filter != "" {
					header = fmt.Sprintf("Projects (%s)", filter)
				}
				return commandResult{output: fmt.Sprintf("%s:\n\n%s\n\nTip: Try 'projects --ml' or 'projects --infrastructure'", header, b.String())}
			},
		},
		"skills": {
			description: "Show technical expertise",
			execute: func(args []string) commandResult {
				expertOnly := false
				for _, a := range args {
					if a == "--expert" {
						expertOnly = true
					}
				}

				if expertOnly {
					return commandResult{output: `Expert-Level Skills:

Functional Programming & Type Systems
  • Haskell (FFI, Cabal, type-safe DSL design)
  • Nix (Advanced derivations, flakes, overlays)
  • Type Systems, Pure Functions, Immutability

Infrastructure & System Administration
  • NixOS (100+ modules, 4+ platforms)
  • Linux/macOS SysAdmin (systemd, launchd)
  • Infrastructure as Code (Declarative configs)
  • CI/CD Automation (Hydra, GitHub Actions)

Level: ███████████░░░ Expert`}
				}

				var b strings.Builder
				for i, cat := range skillCategories {
					if i > 0 {
						b.WriteString("\n\n")
					}
					b.WriteString(fmt.Sprintf("%s:\n  • %s", cat, strings.Join(skills[cat], "\n  • ")))
				}

				return commandResult{output: fmt.Sprintf("Technical Skills:\n\n%s\n\nTip: Try 'skills --expert' for expert-level skills only", b.String())}
			},
		},
		"contact": {
			description: "Display contact information",
			execute: func(_ []string) commandResult {
				return commandResult{output: `Contact Information:

Email:        divitmittal@outlook.in
GitHub:       github.com/DivitMittal
LinkedIn:     linkedin.com/in/divit-mittal
HuggingFace:  huggingface.co/divitmittal
Kaggle:       kaggle.com/divitmittal
Twitter:      x.com/Divit_Mittal

Feel free to reach out!`}
			},
		},
		"experience": {
			description: "Show work experience",
			execute: func(_ []string) commandResult {
				return commandResult{output: `Work Experience:

Computer Vision Researcher @ Manipal University (May-Aug 2025)
   • Hybrid Transformer MFIF research
   • Model on HuggingFace with live demo
   • Benchmark evaluation on LLVIP, TNO, MSRS datasets

DevOps & SysAdmin @ UbiOps, Yes!Delft (Dec 2021 - Feb 2022)
   • Containerized ML workloads (Docker/Kubernetes)
   • CI/CD pipelines, AWS/GCP infrastructure`}
			},
		},
		"resume": {
			description: "View CV",
			execute: func(_ []string) commandResult {
				return commandResult{output: `Europass CV:

https://europa.eu/europass/eportfolio/screen/share/33683761-9b95-4b23-8184-fa65f7a14e06?lang=en`}
			},
		},
	})

	orgCommands = merge(sharedCommands, map[string]commandDef{
		"help": {
			description: "Show available commands",
			execute: func(_ []string) commandResult {
				return commandResult{output: `Available commands:

  help              Show this help message
  about             What is Qezta?
  repos             List public github.com/Qezta repositories
  contact           Display contact information
  clear             Clear terminal
  nix-build         Build profile (Easter egg)
  sudo rm -rf /     Don't even try this
  cowsay <text>     Make the cow say something

Type any command to get started!`}
			},
		},
		"about": {
			description: "What is Qezta?",
			execute: func(_ []string) commandResult {
				return commandResult{output: `Qezta

qezta.com is currently under active development as a proper company/
product site. In the meantime this org hosts Divit Mittal's personal
portfolio and a handful of shared sub-products.

Looking for Divit specifically? Try: ssh divit@qezta.com
Type 'repos' to see what's public under github.com/Qezta.`}
			},
		},
		"repos": {
			description: "List public github.com/Qezta repositories",
			execute: func(_ []string) commandResult {
				var b strings.Builder
				for i, r := range orgRepos {
					if i > 0 {
						b.WriteString("\n")
					}
					b.WriteString(fmt.Sprintf("%-16s %s", r.name, r.desc))
				}
				return commandResult{output: fmt.Sprintf("Qezta Repositories:\n\n%s\n\nAll public: github.com/Qezta", b.String())}
			},
		},
		"contact": {
			description: "Display contact information",
			execute: func(_ []string) commandResult {
				return commandResult{output: `Contact:

GitHub Org:  github.com/Qezta
Security:    security@qezta.com
Personal:    ssh divit@qezta.com (or divit.qezta.com)`}
			},
		},
	})
}

func merge(base, overlay map[string]commandDef) map[string]commandDef {
	out := make(map[string]commandDef, len(base)+len(overlay))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overlay {
		out[k] = v
	}
	return out
}

// executeCommand ports executeCommand() from terminal-commands.ts, extended
// to dispatch against whichever persona the session's username resolved to.
func executeCommand(id identity, input string) commandResult {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return commandResult{}
	}

	if def, ok := multiWordCommands[trimmed]; ok {
		return def.execute(nil)
	}

	commands := personalCommands
	if id == identityQezta {
		commands = orgCommands
	}

	fields := strings.Fields(trimmed)
	cmd, args := fields[0], fields[1:]

	def, ok := commands[cmd]
	if !ok {
		return commandResult{
			output:  fmt.Sprintf("Command not found: %s\nType 'help' for available commands.", cmd),
			isError: true,
		}
	}

	return def.execute(args)
}
