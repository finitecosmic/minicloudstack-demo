# Project Name

This project provides a common intermediate representation (IR) for infrastructure resources.
Resources are represented in a provider-independent resource graph, then a planner determines
how those resources can be implemented using the capabilities of a specific provider.

## Overview

Infrastructure definitions are often tightly coupled to a specific provider, making them difficult to reuse or translate across AWS, Kubernetes, CloudStack, and other platforms. This project separates infrastructure intent from provider implementation, allowing the same resource model to be planned against different providers based on their capabilities. 

## Setup

## 