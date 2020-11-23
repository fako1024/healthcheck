# Health check tool to monitor services (e.g. for use with Docker health checks)

[![Github Release](https://img.shields.io/github/release/fako1024/healthcheck.svg)](https://github.com/fako1024/healthcheck/releases)
[![GoDoc](https://godoc.org/github.com/fako1024/healthcheck?status.svg)](https://godoc.org/github.com/fako1024/healthcheck/)
[![Go Report Card](https://goreportcard.com/badge/github.com/fako1024/healthcheck)](https://goreportcard.com/report/github.com/fako1024/healthcheck)
[![Build/Test Status](https://github.com/fako1024/healthcheck/workflows/Go/badge.svg)](https://github.com/fako1024/healthcheck/actions?query=workflow%3AGo)

This package provides a single binary capable of performing a wide set of health calls to different services, e.g. HTTP(S) & SQL, as well as low-level protocol endpoints, e.g. TCP. It is entended to be used in Docker (Compose) deployments, but can be applied to different scenarios just as well.
