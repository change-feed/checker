# Checker

## Overview

Checker is a component of the change-feed system that detects changes in web pages, changelogs, and documentation by periodically fetching snapshots and comparing them against previous versions. It integrates with [Dapr](https://docs.dapr.io/concepts/building-blocks-concept/) for state management and runs efficiently within a Kubernetes environment.

## Features

- Fetches and stores web page snapshots (configured via the `CHECK_URL` env).
- Compares snapshots using a FIFO-like storage mechanism.
- Uses Dapr state management for backend storage (MinIO in homelab, AWS S3 in production).
- Runs periodic checks via Kubernetes CronJobs.

## Architecture

- Dapr-powered storage

    The system interacts with [Dapr’s state management API](https://docs.dapr.io/developing-applications/building-blocks/state-management/state-management-overview/) to store and retrieve snapshots efficiently.

- Kubernetes-native orchestration

    Uses CronJobs to schedule snapshot checks without requiring long-running processes.

- Configurable backends

    Designed to support multiple storage backends, with current focus on S3-compatible object storage.
