# Docker AI Model Cards

Official model cards and assets for AI models on Docker Hub (`ai` namespace).

## Repository layout

- `ai/` — Model card markdown files (one file per model).
- `template.md` — Model card template to use for new models.
- `images/` — Supplemental images referenced by model cards.
- `logos/` — Provider/model logos used in cards.
- `tools/` — Utilities for maintaining cards and validating models.

## Add or update a model card

1. Copy `template.md` into `ai/` and name it after the model (e.g., `ai/llama3.3.md`).
2. Fill in all sections and tables.
3. If you add images or logos, place them in `images/` or `logos/` and reference them from the card.

## Tools

- Model Cards CLI (`tools/model-cards-cli`): updates the “Available model variants” table by inspecting model registries.
