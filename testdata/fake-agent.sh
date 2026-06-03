#!/bin/sh
# Fake agent for integration testing.
# Creates a marker file so tests can verify the agent ran in the correct directory.
touch agent-was-here.txt
