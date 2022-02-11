#!/bin/bash
set -e

# This is a test suite. Add tests that need to run here.
bash ./scripts/tests/e2e/test_create_new_topic.sh
bash ./scripts/tests/e2e/test_add_new_partitions.sh
bash ./scripts/tests/e2e/test_change_cleanup_policy.sh
bash ./scripts/tests/e2e/test_delete_topic_manifest.sh