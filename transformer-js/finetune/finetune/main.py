import json

import torch
from datasets import Dataset
from peft import (
  LoraConfig,
  get_peft_model,
)
from transformers import (
  AutoModelForCausalLM,
  AutoTokenizer,
  BitsAndBytesConfig,
)
from trl import SFTTrainer, SFTConfig

PROMPT_TEMPLATE = """You are given a database schema and a question.
Based on the schema, generate SQL SELECT statement that answers the question.

Schema:
CREATE TABLE countries (
  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  name TEXT NOT NULL,
  area       INTEGER,
  area_land  INTEGER,
  area_water INTEGER,
  population        INTEGER,
  population_growth REAL,
  birth_rate        REAL,
  death_rate        REAL,
  migration_rate    REAL,
  flag_description TEXT
)

Question:
{question}
"""

if __name__ == '__main__':

  base_model = "unsloth/Llama-3.2-1B-Instruct"
  dataset_name = "countries-training-data.json"
  new_model = "Llama-3.2-1B-Instruct-Country-SQL"

  torch_dtype = torch.float16
  attn_implementation = "eager"

  # QLoRA config
  bnb_config = BitsAndBytesConfig(
    load_in_4bit=False,
    bnb_4bit_quant_type="nf4",
    bnb_4bit_compute_dtype=torch_dtype,
    bnb_4bit_use_double_quant=True,
  )

  # Load model
  model = AutoModelForCausalLM.from_pretrained(
    base_model,
    quantization_config=bnb_config,
    device_map="auto",
    attn_implementation=attn_implementation
  )

  tokenizer = AutoTokenizer.from_pretrained(base_model)

  # LoRA config
  peft_config = LoraConfig(
    r=16,
    lora_alpha=32,
    lora_dropout=0.05,
    bias="none",
    task_type="CAUSAL_LM",
    target_modules=['up_proj', 'down_proj', 'gate_proj', 'k_proj', 'q_proj', 'v_proj', 'o_proj']
  )
  model = get_peft_model(model, peft_config)

  messages = []
  with open(dataset_name, 'r') as f:
    for line in f:
      conversation = json.loads(line)
      user_message = next(filter(lambda x: x["role"] == "user", conversation))
      user_message["content"] = PROMPT_TEMPLATE.format(question=user_message["content"])
      obj = {"text": (tokenizer.apply_chat_template(conversation, tokenize=False, add_generation_prompt=False))}
      messages.append(obj)

  dataset = Dataset.from_list(messages)

  dataset = dataset.train_test_split(test_size=0.1)

  training_arguments = SFTConfig(
    output_dir=new_model,
    per_device_train_batch_size=1,
    per_device_eval_batch_size=1,
    gradient_accumulation_steps=2,
    optim="paged_adamw_32bit",
    num_train_epochs=5,
    evaluation_strategy="steps",
    eval_steps=0.2,
    logging_steps=1,
    warmup_steps=10,
    logging_strategy="steps",
    learning_rate=2e-4,
    fp16=True,
    bf16=False,
    group_by_length=True,
  )

  trainer = SFTTrainer(
    model=model,
    train_dataset=dataset["train"],
    eval_dataset=dataset["test"],
    peft_config=peft_config,
    max_seq_length=512,
    dataset_text_field="text",
    tokenizer=tokenizer,
    args=training_arguments,
    packing=False,
  )

  # Train model
  trainer.train()

  # Quick test
  model.config.use_cache = True
  user_prompt = PROMPT_TEMPLATE.format(question="Which is the smallest country?")
  messages = [
    {
      "role": "user",
      "content": user_prompt
    }
  ]
  prompt = tokenizer.apply_chat_template(messages, tokenize=False,
                                         add_generation_prompt=True)
  inputs = tokenizer(prompt, return_tensors='pt', padding=True,
                     truncation=True).to("cuda")
  outputs = model.generate(**inputs, max_length=150,
                           num_return_sequences=1)
  text = tokenizer.decode(outputs[0], skip_special_tokens=True)
  print(text.split("assistant")[1])

  # Save model
  trainer.model.save_pretrained(new_model)
  trainer.model.push_to_hub(new_model, use_temp_dir=False, token="hf_")
