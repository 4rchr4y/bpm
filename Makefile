SHELL := /bin/bash

export BPM_PATH=$(HOME)/.bpm
export BPM_VERSION=v0.0.1


.PHONY: count
count:
	@ find . -name tests -prune -o -type f -name '*.go' | xargs wc -l
