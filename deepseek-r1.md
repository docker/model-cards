# Deepseek-R1-Distill-Llama

![logo](https://github.com/docker/model-cards/raw/refs/heads/main/logos/deepseek-280x184-overview@2x.svg)

DeepSeek introduced its first-generation reasoning models, DeepSeek-R1-Zero and DeepSeek-R1, leveraging reinforcement learning to enhance reasoning performance, with DeepSeek-R1 achieving state-of-the-art results and open-sourcing multiple distilled models.

The models provided here are the distill-llama variants, which are llama based models that have been fine-tuned on the responses and reasoning output of the full DeepSeek-R1 model.

## Characteristics

| Attribute             | Details          |
|---------------------- |----------------- |
| **Provider**          | Deepseek         |
| **Architecture**      | llama            |
| **Cutoff Date**       | May 2024ⁱ        |
| **Languages**         | English, Chinese |
| **Tool Calling**      | ✅               |
| **Input Modalities**  | Text             |
| **Output Modalities** | Text             |
| **License**           | [MIT](https://github.com/deepseek-ai/DeepSeek-R1/blob/main/LICENSE)           |

i: Estimated

## Available Model Variants

| Model Variant                                                                      | Parameters | Quantization   | Context Window  | VRAM     | Size  |
|------------------------------------------------------------------------------------|----------- |----------------|---------------- |--------- |-------|
| `ai/deepseek-r1-distill-llama:70B-Q4_K_M`                                          | 70B        | IQ2_XXS/Q4_K_M | 128K tokens     | 42GB¹    | 42GB  |
| `ai/deepseek-r1-distill-llama:8B-F16`                                              | 8B         | F16            | 128K tokens     | 19.2GB¹  | 16GB  |
| `ai/deepseek-r1-distill-llama:latest`<br><br>`ai/deepseek-r1-distill-llama:8B-Q4_K_M`                                           | 8B         | IQ2_XXS/Q4_K_M | 128K tokens     | 4.5GB¹   | 5GB   |

¹: VRAM estimated based on model characteristics.

> `:latest` → `70B-Q4_K_M`

## Intended Uses

Deepseek-R1-Distill-Llama can help with:
- **Software Development:** Generates code, debugs, and explains complex concepts.
- **Mathematics:** Solves and explains complex problems for research and education.
- **Content Creation & Editing:** Writes, edits, and summarizes content for various industries.
- **Customer Service:** Powers chatbots to engage users and answer queries.
- **Data Analysis:** Extracts insights and generates reports from large datasets.
- **Education:** Acts as a digital tutor, providing clear explanations and personalized lessons.

## Considerations

- Set the **temperature between 0.5 and 0.7 (recommended: 0.6)** to avoid repetition or incoherence.
- **Do not use a system prompt**; include all instructions within the user prompt.
- For math problems, add a directive like: "Please reason step by step and enclose the final answer in \boxed{}."

This model is sensitive to prompts. Few-shot prompting consistently degrades its performance. Therefore, we
recommend users directly describe the problem and specify the output format using a
zero-shot setting for optimal results.


## How to Run This AI Model

You can pull the model using:
```
docker model pull ai/deepseek-r1-distill-llama
```

To run the model:
```
docker model run ai/deepseek-r1-distill-llama
```


## Benchmark Performance

| Category    | Benchmark                   | DeepSeek R1  |
|-------------|-----------------------------|------------- |
| **English** |                             |              |
|             | MMLU (Pass@1)               | 90.8         |
|             | MMLU-Redux (EM)             | 92.9         |
|             | MMLU-Pro (EM) |             | 84.0         |
|             | DROP (3-shot F1) |          | 92.2         |
|             | IF-Eval (Prompt Strict) |   | 83.3         |
|             | GPQA-Diamond (Pass@1) |     | 71.5         |
|             | SimpleQA (Correct) |        | 30.1         |
|             | FRAMES (Acc.) |             | 82.5         |
|             | AlpacaEval2.0 (LC-winrate)  | 87.6         |
|             | ArenaHard (GPT-4-1106)      | 92.3         |
| **Code**    |                             |              |
|             | LiveCodeBench (Pass@1-COT)  | 65.9         |
|             | Codeforces (Percentile)     | 96.3         |
|             | Codeforces (Rating)         | 2029         |
|             | SWE Verified (Resolved)     | 49.2         |
|             | Aider-Polyglot (Acc.)       | 53.3         |
| **Math**    |                             |              |
|             | AIME 2024 (Pass@1)          | 79 .8        |
|             | MATH-500 (Pass@1)           | 97.3         |
|             | CNMO 2024 (Pass@1)          | 78.8         |
| **Chinese** |                             |              |
|             | CLUEWSC (EM)                | 92.8         |
|             | C-Eval (EM)                 | 91.8         |
|             | C-SimpleQA (Correct)        | 63.7         |


## Links
- [DeepSeek-R1: Incentivizing Reasoning Capability in LLMs via Reinforcement Learning](https://arxiv.org/abs/2501.12948)
