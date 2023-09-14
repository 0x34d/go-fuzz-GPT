import os
import openai

openai.api_key = os.getenv("OPENAI_API_KEY")

#response = openai.File.create(
#  file=open("ldap.jsonl", "rb"),
#  purpose='fine-tune'
#)

#response = openai.FineTuningJob.create(training_file="file-yaQpWsd5DDzT9esj6dPq7AXs", model="gpt-3.5-turbo")


# Retrieve the state of a fine-tune
#response = openai.FineTuningJob.retrieve("ftjob-4ZFxFOjEMJhpFqKyetHh3F4O")

response = openai.FineTuningJob.list_events(id="ftjob-4ZFxFOjEMJhpFqKyetHh3F4O", limit=10)

print(response)

#ft:gpt-3.5-turbo-0613:0x34d::7yZLjRMZ
