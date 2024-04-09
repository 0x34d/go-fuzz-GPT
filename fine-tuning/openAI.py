import os
import openai

openai.api_key = os.getenv("OPENAI_API_KEY")

response = openai.File.create(
  file=open("data.jsonl", "rb"),
  purpose='fine-tune'
)

response = openai.FineTuningJob.create(training_file="file-VXDnyJn14ejnKVQyEdUvYgvD", model="ft:gpt-3.5-turbo-0613:0x34d::7yaoKJRt")

response = openai.FineTuningJob.retrieve("ftjob-RL1Yjc2M8mDrdzGhEbN2dlCW")

response = openai.Model.delete("ft:gpt-3.5-turbo-0613:0x34d::7yZLjRMZ")

print(response)
