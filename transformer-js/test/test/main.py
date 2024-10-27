from transformers import AutoTokenizer, AutoModelForCausalLM

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
    tokenizer = AutoTokenizer.from_pretrained("ralscha/Llama-3.2-1B-Instruct-Country-SQL")
    model = AutoModelForCausalLM.from_pretrained("ralscha/Llama-3.2-1B-Instruct-Country-SQL")
    messages = [
        {
            "role": "user",
            "content": (
                PROMPT_TEMPLATE.format(question="Show me countries where the population is greater than 10 million."))
        }
    ]
    prompt = tokenizer.apply_chat_template(messages, tokenize=False,
                                           add_generation_prompt=True)
    inputs = tokenizer(prompt, return_tensors="pt", padding=True,
                       truncation=True)

    outputs = model.generate(**inputs, max_new_tokens=150,
                             num_return_sequences=1)

    text = tokenizer.decode(outputs[0], skip_special_tokens=True)
    print(text.split("assistant")[1])
