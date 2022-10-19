#!/usr/bin/env bash

# Migration database (tables not deleted with this action)
/app database migration

# Start the server
/app server
