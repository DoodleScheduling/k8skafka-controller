#!/bin/bash
set -e

sh ./scripts/e2e/test_create_new_topic.sh
sh ./scripts/e2e/test_add_new_partitions.sh
sh ./scripts/e2e/test_change_cleanup_policy.sh
sh ./scripts/e2e/test_delete_topic_manifest.sh