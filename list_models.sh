#!/bin/bash

# initialize output files
OUTPUT_FILE="assets/models.yaml"
TMP_FILE="assets/models.json"

> $OUTPUT_FILE
> $TMP_FILE

# regions to check
REGIONS=(
  "us-east-1"
  "us-east-2"
  "us-west-1"
  "us-west-2"
  "af-south-1"
  "ap-east-1"
  "ap-south-1"
  "ap-south-2"
  "ap-northeast-1"
  "ap-northeast-2"
  "ap-northeast-3"
  "ap-southeast-1"
  "ap-southeast-2"
  "ap-southeast-3"
  "ap-southeast-4"
  "ca-central-1"
  "ca-west-1"
  "eu-central-1"
  "eu-central-2"
  "eu-west-1"
  "eu-west-2"
  "eu-west-3"
  "eu-south-1"
  "eu-south-2"
  "eu-north-1"
  "il-central-1"
  "me-central-1"
  "me-south-1"
  "sa-east-1"
)

echo "Fetching Bedrock models from all available regions..."

# loop through regions
for region in "${REGIONS[@]}"; do
  echo "Checking region: $region"
  if models=$(aws bedrock list-foundation-models --region $region --query "modelSummaries[*].{modelId:modelId,modelName:modelName}" --output json 2>/dev/null); then
    if [ -n "$models" ] && [ "$models" != "[]" ]; then
      echo "$models" | jq -c '.[]' >> $TMP_FILE
    fi
  fi
  if models=$(aws bedrock list-inference-profiles --region $region --query "inferenceProfileSummaries[*].{modelId:inferenceProfileId,modelName:inferenceProfileName}" --output json 2>/dev/null); then
    if [ -n "$models" ] && [ "$models" != "[]" ]; then
      echo "$models" | jq -c '.[]' >> $TMP_FILE
    fi
  fi
done

# remove duplicates
cat $TMP_FILE | jq -s 'unique_by(.modelId) | .[]' | \
  jq -r 'select((.modelId | split(":") | length) < 3) | "- modelId: \"" + .modelId + "\"\n  modelName: \"" + .modelName + "\""' >> $OUTPUT_FILE
rm $TMP_FILE
