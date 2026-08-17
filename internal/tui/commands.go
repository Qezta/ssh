package tui

import (
	"fmt"
	"sort"
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

// multiWordCommands are matched against the full trimmed input line before
// falling back to first-token dispatch, so easter eggs like "sudo rm -rf /"
// actually fire instead of being shadowed by tokenization.
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

var commands map[string]commandDef

func init() {
	commands = map[string]commandDef{
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
		"clear": {
			description: "Clear terminal",
			execute: func(_ []string) commandResult {
				return commandResult{clear: true}
			},
		},
		"nix-build": {
			description: "Build profile",
			execute: func(_ []string) commandResult {
				return commandResult{output: `these 42 derivations will be built:
  /nix/store/...-divit-profile-0.1.0.drv
  /nix/store/...-expertise-nix-infrastructure.drv
  /nix/store/...-expertise-functional-programming.drv

building '/nix/store/...-divit-profile-0.1.0.drv'...
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
output: /nix/store/...-divit-profile-0.1.0

Success! Profile built successfully.`}
			},
		},
		"cowsay": {
			description: "Make the cow say something",
			execute: func(args []string) commandResult {
				text := strings.Join(args, " ")
				if text == "" {
					text = "Hello from Divit!"
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
}

// executeCommand ports executeCommand() from terminal-commands.ts.
func executeCommand(input string) commandResult {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return commandResult{}
	}

	if def, ok := multiWordCommands[trimmed]; ok {
		return def.execute(nil)
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

// sortedCommandNames is used by the help/completion UI.
func sortedCommandNames() []string {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
