# Hazel!

#### Deployment framework/app

Hazel aim's at helping you implement __Continuous Integration__ and __Continuous
Delivery__ by provinding you with an platform that will keep track of your apps,
environments, versions, and builds. It has a nice web UI so that anybody
can trigger a deploy of the latest commit on master or a specific branch, but
also support webhooks so that you can deploy right after the CI runs.

It is similar in functionality to [Deployinator](https://github.com/etsy/deployinator)
by Etsy or [Dreadnot](https://github.com/racker/dreadnot) by Rackspace but adds
a little more flexibility to fit in more use cases. Plus intead is being a monolith
like those last two it is composes of multiple small apps working together all
doing one thing well and providing functionalities. That way you can easily
replace components you would have wish where built in an other way.

## The big plan

The end goal of Hazel is enabling DevOps and shipping continuously and efficiently
a SOA. Microservices makes you code easy to write in those decoupled but what's
in between of them and the infrastructure realted to them gets more complex compared
to a monolith.

With the end goal in mind, the secondary goal is to be able to use **Hazel** yesterday
in production, so, we are shooting for MVP. Meaning, stuff like service discovery,
dynamic DNS updating, rolling updates, supporting more than 5-10 apps and other
stacks than Node.js is not needed for basic usage and will come in time.

### Phase 1 : Sketching things out, POC

**Api**

Holds the logic for interacting with models in _Hazel_, namely, **Apps**, **Builds**,
and **Stacks**. It will be used by the UI to fetch all info needed. Also used by
webhooks to create new builds on demand.

**Builder**

The builder is the meat of Hazel. It's where a build goes through all the necessary
steps to become deployable using terraform. Here's what it does:

- `status=created` Waits for new builds
- `status=waiting` Polls S3 waiting for a tar to be uploaded with the application
   code to be uploaded
- `status=building` Fetches up to date env vars, parses and outputs all templates
  defined by the stack to have dynamic configuration, runs **packer** (from
  Hashicorp) and saves the logs to S3
- `status=built` This is the status after `building`, at this point the **builder**
  saved any build artifacts it needed to save back to S3 and updated the build
  with the newly build **ami** name, plus, updated the app metadata to deploy 10%
  of the normal amount of servers using the the image, this way updates roll foward
  slowly and complete propagation to all instances can be done via the UI or will
  be done automagicly after a timeout.

**Scheduler**

This part of Hazel is like an overseer that ensures the right amount of instances
are deployed for each app based on what is specified in the current app config.

See **Phase 2** for the desired role of this service. But, for **Phase 1** this
app will simply run terraform with the newly build **ami**'s when it sees undeployed
builds.

**Environment (vars) API**

Environment variables being strored centrally is a huge plus as you can have a web
UI to easily change them and trigger deploys, this service is a simple JSON over
HTTP api that only holds a few endpoints for apps and their config vars.

This service is used by the builder when creating new **ami**'s and by the
**Environment UI** that allows editing of those environment variables.

**Environment (vars) UI**

This UI allow editing environment variables for one app. It's not directly in
the **Web UI** as many more feature needing web interfaces are comming and if
everything is added to one web project it will become an unmaintainable pile
of javascript.

**Web UI**

This would be the interface to every services composing Hazel, it would only query
the API service as it's the one that holds all the config, metadata and state and
would be able to:

- Trigger builds for an app
- See a list of builds for an app and their state
- See a specific build's logs
- Promote a canary deploy (all instances become most recent build)
- Reject a canary deploy (roll back)
- Add a new app

**Datastore**

Let's be relly hacky and use only the filesystem, ok S3. After all _Deployinator_
from Etsy only keep log files in a folder structure, that's all. So that will do
for an MVP. No message queue and no MongoDB.

### Phase 2 : Road to production ready

**Canary pusher**

Simple, this small service polls app metadata and check from **canary deploys**
that hapenned more than X minutes ago, of for degrading metrics past a certain
threshold, then makes all instances use the canary ami (finish rolling update)
or removes and canary instances and marks the build as `rejected` (roll back)

**Scheduler**

This part of Hazel is like an overseer that ensures the right amount of instances
are deployed for each app based on what is specified in the current app config.

This would allow for nice use cases where changing the number of v23 instances
to 0 and v24 instances to 7 in the web UI would simply update the **Datastore**
and soon after the scheduler will pickup the difference with the current state
and adjust. Same thing if, let's say, _AWS kill one of your instances!_.

**Web UI**

More features, more control, more polish

**Datastore**

Probably move to a S3 for config and artifacts (e.g. build logs), MongoDB for
state and metadata and some MQ for Hazel components communication.

### Phase 3 : Better that Heroku

**Deployment**

This is not a service, but still is a feature. Deployment at this point should
become automated but mostly robust and secure for all **Hazel** components as
they will become the backbone of companies using it in production. This includes
double checking:

- Configuration management
- Alerts for unhealth core services
- Dogfooding (deploying Hazel non-core services using hazel)
- Not dogfooding, Hazel should be able to operate event when parts of it are
  also it should be monitored and alerts should go off even it, say the monitoring
  service of Hazel is not up
- Security double checked, firewalls up, unused ports closed, login needed everywhere,
  VPN needed for internal only and maintenance, SSL everywhere

**Log Aggregator**

Util now it was OK to use logentries directly on machines, now we might want to
start using `syslog` everywhere and push that to an **ElasticSearch** cluster all
while saving it to **Glacier** just in case. This service is the piece in the middle,
it can also be a fleet of instances running **Logstash** and a `syslog` reciever.

**CLI**

This might seem like long overdue, but I feel web UIs do a better job a selling
a product, even to developpers, than CLIs do. But, we still need one for efficiency
and scripting's sake.

Start with:

- Trigger build
- Set env vars
- Tail logs
- List live versions and server count

**Monitoring Collector**

The pull model, at scale, works better than push (reads Soundcloud's experience)
this would pull and stash in a datastore info from logs, loads balancers /stats
endpoints, nsq /stats endpoints, different datastores used by apps but also
`expvar` (the Golang package) style metrics.

**Monitoring Producers**

Deploy a fleet of **statsd** instances.

**Monitoring web UI**

Most likely **Graphana**, depending on the datastore, maybee something custom
built.

**Monitoring Datastore**

Hazel will use **statsd** for collecting metrics but the actual _drain_ for it
is not decided yet. Maybee Graphite, maybee InfluxDB, maybee something else. But
one thing is sure, it wont be written specificaly for Hazel. It would still imply
we need to be able to deploy it in a scalable way.

### Phase 4 : Embrace microservices

**Service Discovery**

Up till now DNS was fine, now Consul or else is needed

**Request Tracing**

Think Google Dapper or Twitter Zipkin, gokit is cooking an implementation,
can simply be a new service writing to Kafka what it sees in logs.

**Better Overseer Version Management**

With microservices comming you might want to keep multiple versions of services
up and running for longer than till next build, this should be simple but requires
thinking.

**Better Monitoring Datastore**

The scale of events will grow really fater with microservices interacting with
each other that with one monolith

**Better Monitoring Collector (more integrations)**

**More complex overseer, builder, log agregator, monitoring**

Support deploying container to a cluster of beefy **CoreOS** hosts.

## Getting started

```bash
fig up
```

## Plan

__This part (the design) is constantly changing, reading the code is more acurate)

![Plan](https://raw.githubusercontent.com/kiasaki/hazel/master/images/plan.png)

![Screenshot](https://raw.githubusercontent.com/kiasaki/hazel/master/images/screenshot.png)
