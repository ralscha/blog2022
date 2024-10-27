import torch
from peft import PeftModel
from transformers import AutoModelForCausalLM, AutoTokenizer, pipeline

from finetune.main import PROMPT_TEMPLATE

base_model = "unsloth/Llama-3.2-1B-Instruct"
new_model = "Llama-3.2-1B-Instruct-Country-SQL"

tokenizer = AutoTokenizer.from_pretrained(base_model)

base_model_reload = AutoModelForCausalLM.from_pretrained(
  base_model,
  return_dict=True,
  low_cpu_mem_usage=True,
  torch_dtype=torch.float16,
  device_map="auto",
  trust_remote_code=True,
)

# Merge adapter with base model
model = PeftModel.from_pretrained(base_model_reload, new_model)
model = model.merge_and_unload()

# Quick test
messages = [{"role": "user", "content": PROMPT_TEMPLATE.format(question="Which is the largest country?")}]
prompt = tokenizer.apply_chat_template(messages, tokenize=False, add_generation_prompt=True)
pipe = pipeline(
  "text-generation",
  model=model,
  tokenizer=tokenizer,
  torch_dtype=torch.float16,
  device_map="auto",
)

outputs = pipe(prompt, max_new_tokens=120, do_sample=True, temperature=0.7, top_k=50, top_p=0.95)
print(outputs[0]["generated_text"])

# Save and push to hub
model.save_pretrained(new_model)
tokenizer.save_pretrained(new_model)
model.push_to_hub(new_model, use_temp_dir=False, token="hf_")
tokenizer.push_to_hub(new_model, use_temp_dir=False, token="hf_")
