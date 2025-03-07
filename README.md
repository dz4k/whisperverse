# Whisperverse 👻

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://pkg.go.dev/github.com/whisperverse/whisperverse)
[![Build Status](https://img.shields.io/github/workflow/status/whisperverse/whisperverse/Go/master)](https://github.com/whisperverse/whisperverse/actions/workflows/go.yml)
[![Codecov](https://img.shields.io/codecov/c/github/whisperverse/whisperverse.svg?style=flat-square)](https://codecov.io/gh/whisperverse/whisperverse)
[![Go Report Card](https://goreportcard.com/badge/github.com/whisperverse/whisperverse?style=flat-square)](https://goreportcard.com/report/github.com/whisperverse/whisperverse)
[![Version](https://img.shields.io/github/v/release/whisperverse/whisperverse?include_prereleases&style=flat-square&color=brightgreen)](https://github.com/whisperverse/whisperverse/releases)

## This is a work-in-progress that is *guaranteed to NOT work* for anyone, in any capacity, at any time.  Do not use this code in production, in development, or even in theory, because it's all wrong and may never be corrected, maintained, or supported

## Why Whisperverse?

In the early 2000's, social media got many things right: it's easy to publish, easy to subscribe, and share.
But things went off the rails, and our current social media landscape is far from perfect.  Now, social media is synonymous with disinformation, invasive tracking, and lack of control.  

Malicious algorithms with global reach churn human beings into ad revenue.  It didn't have to be this way.  Ghost is a social CMS with a small reach, allowing you to stay connected to the most important people -- the friends and family in your inner circle.

## What is It?

Ghost is a new kind of decentralized, private media server that will connect people instead of driving them apart, and will return power and privacy to users and content creators.

### Decentralized

When completed, this will be a new kind of personal media server, meant to be an open, [federated](https://en.wikipedia.org/wiki/Fediverse) replacement for many of the closed, centralized services that we all use today.  

### Private

Ghost belongs to the users, not the service providers.  There is no tracking built in to whisper, and we will work to keep it that way.  Strong access controls make your content easy to share, and easy to manage.

### Social

It will work with customizable templates that will replicate many of the social media services out there: posts, comments, images, videos, real time communications and more.

### Real-Time

Ghost will support several real-time messaging interfaces, pushing live content to your community instantly.  

## Technology Goals

Ghost must be extremely service-provider-friendly: easy to virtualize, provision, and deploy. To make this easy , it should be self-contained, with as few dependencies as possible.  Here are a few of the interfaces that I'd like to implement:

* [RSS](https://en.wikipedia.org/wiki/RSS) / [JSON Feed](https://jsonfeed.org)
* [ActivityPub](https://activitypub.rocks) / [OpenSocial](https://www.getopensocial.com) / [Diaspora](https://diasporafoundation.org)
* [OAuth](https://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=&ved=2ahUKEwjByq6-_K3wAhVeIDQIHdMuCmsQFjAQegQIBBAD&url=https%3A%2F%2Foauth.net%2F&usg=AOvVaw3GDFM0pkIJMe4FATEf5VSd)
* [WebRTC](https://webrtc.org)
* [Server Sent Events](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events)
* [oEmbed](https://oembed.com)

### Toolkit

| Tool | Info|
|---|---|
| [Go](https://golang.org) | Single file executable server, compiled for every OS and architecture |
| [Mongodb](https://mongodb.org) | Database server (possibly swappable) |
| [htmx](https://htmx.org) & [hyperscript](https://hyperscript.org)  | Interactive HTML content
| [ckEditor](https://ckeditor.com/ckeditor-5/) | Rich content editing
| ??? | Various local and online file storage systems

## Pitch In

There's a lot to do, and I'd love to have your help.  If you're interested in building the federated web, please get in touch or submit a pull request. 👻
