# Phi-3 Mini

<img src="https://github.com/jalonsogo/model-cards/blob/4c39899ef2d3eff3bfe28253b557283c8933c811/logos/microsoft.svg" width="120" />

Phi-3 is a family of lightweight, state-of-the-art open models developed by Microsoft, available in 3B (Mini) and 14B (Medium) parameter sizes.

## Characteristics

| Attribute             | Details       |
|---------------------- |-------------- |
| **Provider**          | Microsoft     |
| **Architecture**      | phi3          |
| **Cutoff Date**       | October 2023  |
| **Languages**         | English       |
| **Input Modalities**  | Text          |
| **Output Modalities** | Text          |
| **License**           | MIT           |

## Available Model Variants

| Model Variant  | Parameters | Size   | Quantization | Context Window | VRAM Usage | Download  |
|--------------- |----------- |------- |------------- |--------------- |----------- |---------- |
| `phi3:3B-FP16` | 2.3B       | 5GB | None         | 128K           | 5GB¹       | [Link]    |
| `phi3:3B-Q4`   | 2.3B       | 2.39GB | Q4           | 4K             | 1.3GB¹     | [Link]    |
¹: VRAM estimates based on model characteristics.

## Intended Uses

Phi-3 is designed for both commercial and research applications in English, specifically optimized for:

- Memory/compute-constrained environments
- Latency-sensitive scenarios
- Strong reasoning tasks (math & logic)
- Long-context processing

It facilitates research in language and multimodal models, serving as a foundation for generative AI applications.

### Considerations

The model is not optimized for all applications. Developers should:
- Assess limitations and potential biases
- Ensure accuracy, safety, and fairness
- Follow relevant laws and regulations

This Model Card does not modify the model’s license.

### How to Run This AI Model

You can pull the AI model from Docker Hub using its name:

```sh
docker model pull phi3
```

AI models are OCI Artifacts. If no tag is specified, the `:default` model will be downloaded (usually optimized for a balance between performance and hardware requirements). To specify a particular variant, use:

```sh
docker model pull phi3:2.3B-Q4_K_M
```

Once downloaded, you can run it with the Docker Model Runner™:

```sh
docker model run phi3
```

### Benchmark Performance

| Category                         | Benchmark              | Phi-3-Mini-4K-In | Phi-3-Mini-128K-In |
|--------------------------------- |----------------------- |----------------- |------------------- |
| **Popular Aggregate Benchmarks** | AGI Eval (0-shot)      | 37.5             | 36.9               |
|                                  | MMLU (5-shot)          | 68.8             | 68.1               |
|                                  | BigBench Hard (0-shot) | 71.7             | 71.5               |
| **Language Understanding**       | ANLI (7-shot)          | 52.8             | 52.2               |
|                                  | HellaSwag (5-shot)     | 76.7             | 74.5               |
| **Reasoning**                    | ARC Challenge (10-shot)| 84.9             | 84                 |
|                                  | ARC Easy (10-shot)     | 94.6             | 95.2               |
|                                  | BoolQ (5-shot)         | 77.6             | 78.7               |
|                                  | CommonsenseQA (10-shot)| 80.2             | 78                 |
|                                  | MedQA (2-shot)         | 53.8             | 55.3               |
|                                  | OpenBookQA (10-shot)   | 83.2             | 80.6               |
|                                  | PIQA (5-shot)          | 84.2             | 83.6               |
|                                  | Social IQA (5-shot)    | 76.6             | 76.1               |
|                                  | TruthfulQA (MC2)       | 65               | 63.2               |
|                                  | WinoGrande (5-shot)    | 70.8             | 72.5               |
| **Factual Knowledge**            | TriviaQA (5-shot)      | 64               | 57.1               |
| **Math**                         | GSM8K CoT (0-shot)     | 82.5             | 83.6               |
| **Code Generation**              | HumanEval (0-shot)     | 59.1             | 57.9               |
|                                  | MBPP (0-shot)          | 53.8             | 62.5               |

### Links
- [Microsoft Announcement](https://azure.microsoft.com/en-us/blog/introducing-phi-3-redefining-whats-possible-with-slms/)
- [Phi-3 Technical Report](https://arxiv.org/abs/2404.14219)
