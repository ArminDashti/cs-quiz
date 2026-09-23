INSERT INTO quizzes (name, slug, description, category)
VALUES
    ('AI & Agents', 'ai-agents', 'Artificial intelligence, LLMs, and AI agent concepts')
ON CONFLICT (slug) DO NOTHING;

UPDATE quizzes SET category = 'AI' WHERE slug = 'ai-agents' AND category = '';

-- Intermediate (20 questions)
INSERT INTO questions (quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index, sort_order, question_type, difficulty)
SELECT q.id, v.prompt, v.a, v.b, v.c, v.d, v.correct, v.sort_order, 'multiple_choice', 'medium'
FROM quizzes q
CROSS JOIN (VALUES
    ('What does LLM stand for?', 'Learned Layer Matrix', 'Linear Logic Machine', 'Linguistic Lookup Module', 'Large Language Model', 3, 1),
    ('Which component is typically used to store embeddings for similarity search in RAG systems?', 'Vector database', 'Key-value cache', 'Relational table with indexes', 'Message queue', 0, 2),
    ('In an LLM, what does the “temperature” parameter control?', 'Context window size', 'Model training speed', 'Randomness of token sampling', 'Number of layers', 2, 3),
    ('What is a “token” in the context of LLM APIs?', 'A GPU core', 'A unit of text the model processes', 'An authentication key', 'A training epoch', 1, 4),
    ('RLHF stands for...', 'Recursive Loss Heuristic Function', 'Reinforced Latent Hidden Fields', 'Rapid Labelled Hybrid Fine-tuning', 'Reinforcement Learning from Human Feedback', 3, 5),
    ('Fine-tuning an LLM primarily involves...', 'Increasing the context window', 'Continuing training on task-specific data', 'Deleting unused layers', 'Quantizing weights only', 1, 6),
    ('What is prompt engineering?', 'Designing inputs to steer model output', 'Tuning hyperparameters at train time', 'Writing compiler passes', 'Building database schemas', 0, 7),
    ('In agentic AI, what is tool use / function calling?', 'Compressing prompts', 'The model invokes external functions to act on the world', 'Using a GPU for inference', 'Swapping models mid-request', 1, 8),
    ('What does a system prompt define?', 'The tokenizer vocabulary', 'The API rate limit', 'The training dataset', 'Baseline instructions and persona for the model', 3, 9),
    ('Chain-of-thought prompting asks the model to...', 'Produce intermediate reasoning steps before answering', 'Repeat the prompt', 'Use multiple models', 'Answer as briefly as possible', 0, 10),
    ('Embeddings are best described as...', 'Encrypted files', 'Hashed passwords', 'Dense numeric vectors representing text meaning', 'Sorted word lists', 2, 11),
    ('What is the main purpose of RAG (Retrieval-Augmented Generation)?', 'Ground answers in retrieved external documents', 'Encrypt user prompts', 'Speed up tokenization', 'Reduce model parameters', 0, 12),
    ('Hallucination in LLMs means...', 'Slow response times', 'Fluent but factually incorrect output', 'Refusing to answer', 'Duplicate tokens', 1, 13),
    ('Which metric is commonly used to evaluate classification accuracy?', 'TFLOPS', 'p99 latency', 'F1 score', 'Frames per second', 2, 14),
    ('Supervised learning differs from unsupervised learning because it uses...', 'Labeled training examples', 'No data', 'Human play', 'Only raw text', 0, 15),
    ('Gradient descent is used to...', 'Compress models', 'Sort datasets', 'Tokenize input', 'Minimize a loss function by adjusting weights', 3, 16),
    ('Overfitting occurs when a model...', 'Trains too slowly', 'Has too few parameters', 'Memorizes training data and generalizes poorly', 'Skips validation', 2, 17),
    ('A convolutional neural network (CNN) is most often used for...', 'Text tokenization', 'Network routing', 'Image recognition tasks', 'Database indexing', 2, 18),
    ('What is quantization in model deployment?', 'Counting parameters', 'Reducing weight precision to shrink and speed up models', 'Splitting a model into shards across regions', 'Encrypting model weights', 1, 19),
    ('Which protocol lets an AI agent expose reusable capabilities to other tools?', 'Model Context Protocol (MCP)', 'SMTP', 'FTP', 'SNMP', 0, 20)
) AS v(prompt, a, b, c, d, correct, sort_order)
WHERE q.slug = 'ai-agents'
  AND NOT EXISTS (SELECT 1 FROM questions x WHERE x.quiz_id = q.id AND x.difficulty = 'medium');

-- Advanced (30 questions)
INSERT INTO questions (quiz_id, prompt, option_a, option_b, option_c, option_d, correct_index, sort_order, question_type, difficulty)
SELECT q.id, v.prompt, v.a, v.b, v.c, v.d, v.correct, v.sort_order, 'multiple_choice', 'hard'
FROM quizzes q
CROSS JOIN (VALUES
    ('In the transformer architecture, self-attention computes similarity between...', 'Adjacent pixels only', 'Convolution kernels', 'Hash buckets', 'Queries, keys, and values per token', 3, 1),
    ('Multi-head attention allows the model to...', 'Use multiple GPUs', 'Skip residual connections', 'Attend to different representation subspaces in parallel', 'Stack encoders', 2, 2),
    ('KV caching during autoregressive decoding speeds inference by...', 'Storing past key/value tensors so they need not be recomputed', 'Pruning heads', 'Caching user requests', 'Batching tokens alphabetically', 0, 3),
    ('RoPE (rotary position embeddings) encodes positional information by...', 'Adding sinusoidal scalars to weights', 'Rotating query/key vectors by position-dependent angles', 'Sorting tokens by index', 'Hashing positions', 1, 4),
    ('FlashAttention achieves efficiency primarily through...', 'Weight sharing across layers', 'Lower-rank approximation', 'Removing softmax', 'IO-aware tiled computation avoiding materializing the full attention matrix', 3, 5),
    ('Speculative decoding accelerates generation by...', 'Drafting tokens with a small model and verifying them with the large one', 'Halving the vocabulary', 'Skipping stop tokens', 'Greedy sampling only', 0, 6),
    ('LoRA fine-tuning reduces cost by...', 'Freezing the tokenizer', 'Deleting FFN layers', 'Training low-rank adapter matrices instead of all weights', 'Using fp16 storage', 2, 7),
    ('QLoRA combines LoRA with...', 'Reinforcement learning', 'Sparse MoE routing', 'Beam search', '4-bit quantized base-model weights', 3, 8),
    ('MoE (mixture-of-experts) architectures route tokens via...', 'A learned gating network selecting a subset of expert FFNs', 'Round-robin over layers', 'Static hash of the token ID', 'Random dropout', 0, 9),
    ('PPO, used in RLHF, constrains policy updates using...', 'Dropout noise', 'L2 weight decay', 'Early stopping only', 'A clipped surrogate objective', 3, 10),
    ('DPO differs from PPO-based RLHF because it...', 'Optimizes directly on preference pairs without a separate reward model', 'Trains two actors and a critic', 'Requires human labels at inference', 'Uses Monte Carlo rollouts', 0, 11),
    ('In RAG, chunk overlap mainly helps to...', 'Increase embedding dimensionality', 'Reduce vector DB size', 'Speed up retrieval', 'Preserve context that would be split at boundaries', 3, 12),
    ('Hybrid retrieval in RAG typically combines...', 'Two vector databases', 'Zipf and Bayes estimators', 'BM25 lexical search with dense vector search', 'SQL joins and views', 2, 13),
    ('Cross-encoder reranking improves RAG by...', 'Caching LLM outputs', 'Shrinking the corpus', 'Scoring query-document pairs jointly after initial retrieval', 'Replacing embeddings with hashes', 2, 14),
    ('HNSW indexes accelerate approximate nearest-neighbor search using...', 'Inverted file lists only', 'Radix tries', 'B+ trees', 'A layered navigable small-world graph', 3, 15),
    ('ANN search trades exactness for speed by returning results that are...', 'Randomly sampled', 'Lexicographically first', 'Near-nearest with high probability', 'Guaranteed optimal', 2, 16),
    ('ReAct-style agents interleave...', 'Tokens and sentences', 'Encoders and decoders', 'Training and inference', 'Reasoning traces with environment actions', 3, 17),
    ('In multi-agent systems, a supervisor/orchestrator pattern assigns work by...', 'Routing subtasks to specialized worker agents and aggregating results', 'Broadcasting every message to all agents', 'Shared memory only', 'Lockstep round-robin', 0, 18),
    ('To prevent prompt injection in tool-using agents, a key mitigation is...', 'Shorter prompts', 'Using larger models', 'Disabling temperature', 'Treating retrieved/user content as untrusted and sandboxing tool permissions', 3, 19),
    ('Guardrails for LLM apps typically enforce...', 'GPU quotas', 'Input/output validation policies such as topic or format constraints', 'Tokenizer versions', 'Model checkpoints', 1, 20),
    ('Constitutional AI aligns models by...', 'Deleting unsafe tokens', 'Fixed rule filters only', 'Manual redaction of the corpus', 'Self-critiquing and revising outputs against explicit principles', 3, 21),
    ('Perplexity measures language-model quality as...', 'Accuracy on multiple choice', 'The exponentiated average negative log-likelihood per token', 'Inference latency', 'BLEU against references', 1, 22),
    ('BERT is distinguished from GPT in that BERT...', 'Has no attention layers', 'Uses bidirectional encoder-only attention', 'Is trained only on code', 'Generates text autoregressively', 1, 23),
    ('Teacher forcing during sequence training means...', 'Forcing max-length outputs', 'Feeding ground-truth previous tokens as decoder inputs', 'Disabling dropout', 'Cloning the teacher model', 1, 24),
    ('Knowledge distillation trains a student to match...', 'The teacher''s hardware config', 'The teacher''s soft output distributions', 'The teacher''s random seeds', 'The teacher''s dataset exactly', 1, 25),
    ('Weight decay acts as which regularization term during optimization?', 'L2 penalty on the parameters', 'L1 penalty on activations', 'Label smoothing', 'Dropout probability', 0, 26),
    ('In beam search, increasing beam width generally...', 'Makes decoding faster', 'Eliminates repetition', 'Raises the chance of higher-probability sequences at more compute cost', 'Always lowers perplexity of the true output', 2, 27),
    ('Top-p (nucleus) sampling selects tokens from...', 'The top k logits always', 'The smallest set whose cumulative probability exceeds p', 'Uniform distribution', 'Only repeated tokens', 1, 28),
    ('Continuous batching in LLM serving improves throughput by...', 'Disabling KV cache', 'Concatenating all prompts offline', 'Pinning one request per GPU', 'Admitting new requests into a batch as soon as slots free up', 3, 29),
    ('A common failure mode of naive ReAct loops requiring mitigation is...', 'Too-fast convergence', 'Excessive tokenization cost', 'Underfitting', 'Infinite action loops without progress detection', 3, 30)
) AS v(prompt, a, b, c, d, correct, sort_order)
WHERE q.slug = 'ai-agents'
  AND NOT EXISTS (SELECT 1 FROM questions x WHERE x.quiz_id = q.id AND x.difficulty = 'hard');
