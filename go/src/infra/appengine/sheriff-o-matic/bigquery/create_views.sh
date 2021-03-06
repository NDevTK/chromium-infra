#!/bin/bash
# Copyright 2019 The Chromium Authors. All rights reserved.
# Use of this source code is governed by a BSD-style license that can be
# found in the LICENSE file.

# Change these values as required to set up new views.
APP_ID=sheriff-o-matic-staging
project_names_without_test_results=("chromeos" "fuchsia")

resultdb_dataset="chrome-luci-data.chromium_staging"
if [ "$APP_ID" == "sheriff-o-matic" ]; then
  resultdb_dataset="chrome-luci-data.chromium"
fi

project_name="chrome"
echo "creating data set and views for project chrome"
bq --project_id $APP_ID mk -d "$project_name"
sed -e s/APP_ID/$APP_ID/g step_status_transitions_chrome.sql | bq --project_id $APP_ID query --use_legacy_sql=false
sed -e s/APP_ID/$APP_ID/g -e s/RESULTDB_DATASET/"$resultdb_dataset"/g failing_steps_chrome.sql | bq query --project_id $APP_ID --use_legacy_sql=false
sed -e s/APP_ID/$APP_ID/g -e s/PROJECT_NAME/"$project_name"/g sheriffable_failures.sql | bq --project_id $APP_ID query --use_legacy_sql=false

for project_name in "${project_names_without_test_results[@]}"
do
    echo "creating data set and views for project: $project_name"
    bq --project_id $APP_ID mk -d "$project_name"
    sed -e s/APP_ID/$APP_ID/g -e s/PROJECT_NAME/"$project_name"/g step_status_transitions.sql | bq --project_id $APP_ID query --use_legacy_sql=false
    sed -e s/APP_ID/$APP_ID/g -e s/PROJECT_NAME/"$project_name"/g failing_steps_without_test_results.sql | bq query --project_id $APP_ID --use_legacy_sql=false
    sed -e s/APP_ID/$APP_ID/g -e s/PROJECT_NAME/"$project_name"/g sheriffable_failures.sql | bq --project_id $APP_ID query --use_legacy_sql=false
done