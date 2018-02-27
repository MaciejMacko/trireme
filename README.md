# Trireme Project

<img src="docs/trireme.png" width="400">

This is the meta-repo for the Zero-Trust Trireme project.

## Trireme specs

Trireme defines a protocol and specifications on how to make remote ressources such as containers, Kubernetes pods, Linux processes etc communicate seamlessly and securely.
Those specification are based on the Zero-Trust concept, which requires strong authentication accross all communications.

## Trireme Library

The Trireme specifications are implemented by [Trireme-lib](https://github.com/aporeto-inc/trireme-lib)

## Implementations using Trireme

The following projects use Trireme-lib:

- [Aporeto Enterprise](https://console.aporeto.com) is implemented with Trireme at its core. This implementation uses a full API and vizualization in order to help you design security policies accross multiple cloud instances.
- [Trireme-Kubernetes](https://github.com/aporeto-inc/trireme-kubernetes) is a minimalist implementation specifically aimed to implement Kubernetes network policies.
- [Trireme-Example](https://github.com/aporeto-inc/trireme-example) is a very simple example implementation wrapped around Trireme-Lib. It can be used to deno Trireme application for Docker Containers or Linux Processes.
- [Trireme-bare-metal](https://github.com/aporeto-inc/trireme-bare-metal) is an implementation of Trireme that can be used with a routed backbone.

## Other libraries//utilities maintained as part of the Trireme umbrella

- [Kubepox](https://github.com/aporeto-inc/kubepox). The Kubernetes Policy eXploration tool is a simple library//executable that can be used in order to decide which NetworkPolicy applies to which set of pods or vcice-versa.
- [TG](https://github.com/aporeto-inc/tg)
- [Apobeer](https://github.com/aporeto-inc/apobeer) A simple policed demo application that can be deployed on multiple orchestration platforms.
