#!/bin/bash
set -e

sh ./test_create_new_topic.sh
sh ./test_add_new_partitions.sh
sh ./test_change_cleanup_policy.sh
sh ./test_delete_topic_manifest.sh