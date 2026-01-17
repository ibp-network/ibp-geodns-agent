# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial repository structure matching ibp-geodns-* repos patterns
- Agent bootstrap and configuration loading
- NATS client integration using ibp-geodns-libs
- Health check endpoints (/health, /ready, /live)
- Service monitoring framework
- Status reporting via NATS
- Systemd service installation support
- Structured logging with configurable levels
- Configuration hot-reload support
- Makefile for build and development tasks
