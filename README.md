# Trireme Project

<img src="https://github.com/aporeto-inc/trireme-lib/blob/master/docs/trireme.png" width="400">

This is the meta-repo for the Zero-Trust Trireme project.

## Trireme specs

Trireme defines a protocol and specifications on how to make remote ressources such as containers, Kubernetes pods, Linux processes etc communicate seamlessly and securely.
Those specification are based on the Zero-Trust concept, which requires strong authentication accross all communications.

## Trireme Library

The Trireme specifications are implemented by [Trireme-lib](https://github.com/aporeto-inc/trireme-lib)

## Implementations using Trireme

The following projects use Trireme-lib:

- [Aporeto Enterprise](https://console.aporeto.com) is implemented with Trireme at its core. This implementation uses a full API and visualization in order to help you design security policies accross multiple cloud instances.
- [Trireme-Kubernetes](https://github.com/aporeto-inc/trireme-kubernetes) is a minimalist implementation specifically aimed to implement Kubernetes network policies.
- [Trireme-Example](https://github.com/aporeto-inc/trireme-example) is a very simple example implementation wrapped around Trireme-Lib. It can be used to demo Trireme application for Docker Containers or Linux Processes.
- [Trireme-bare-metal](https://github.com/aporeto-inc/trireme-bare-metal) is an implementation of Trireme that can be used with a routed backbone.

## Other libraries//utilities maintained as part of the Trireme umbrella

- [Trireme-CSR](https://github.com/aporeto-inc/trireme-csr) is a library that generates certificates for Trireme based on a CSR. This library got an implementation that is ready to run on Kubernetes as a controller.
- [Trireme-Statistics](https://github.com/aporeto-inc/trireme-statistics) is an implementation of the `collector` interface for Trireme that sends all statistics to InfluxDB. It also contains some example implementation of visualization based on InfluxB, Grafana and Chronograf.
- [Kubepox](https://github.com/aporeto-inc/kubepox). The Kubernetes Policy eXploration tool is a simple library//executable that can be used in order to decide which NetworkPolicy applies to which set of pods or vcice-versa.
- [TG](https://github.com/aporeto-inc/tg) makes issuing certificates easy. It wraps around the standard golang crypto lib.
- [Apobeer](https://github.com/aporeto-inc/apobeer) A simple policed demo application that can be deployed on multiple orchestration platforms.
