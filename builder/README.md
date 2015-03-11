# Hazel : Builder

## What is hapenning here?

Well we are aiming at building an AWS ami using Packer all that for different
use cases.

So, you will want to deploy diffrent services that have different needs, Golang
apis, nodejs frontends, ruby webapps, ... Hazel's way of enabling all that
diversity is by making the DevOps or Infra folks define "stacks" that dev can
in term use.

**Stacks** are composed of a `config.json` and a `packer.json` file at the root
of it's folder. This will be used in the UI to show a name and description for
each stack and to know which files in the stack folder are templates and should
be preprocessed.

**Templates?** Yes, stacks can define files that the builder should pass through
Go's `text/template` with info about the **environment**, the **current build**,
and the **application** being built.

The **builder's** job ends when an ami was created named `<app_slug>-<build_id>`,
packer logs are uploaded to S3 and the build is updated.
