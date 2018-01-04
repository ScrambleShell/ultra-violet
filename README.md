# Ultra Violet(uv): "Friend Computer" is watching your activities.

__This repo bases clone of "github.com/sourcegraph/thyme".__

Automatically track which applications you use and for how long.

- Simple CLI to track and analyze your application usage
- Detailed charts that let you profile how you spend your time
- Stores data locally, giving you full control and privacy
- Open-source and easily extensible

uv is a work in progress, so please report bugs! Want to see how it works? 

## Citizen, are you happy?

### Simple CLI

1. Record which applications you use every 30 seconds:
   ```
   $ while true; do uv track -o infraRed.json; sleep 30s; done;
   ```

2. Create charts showing application usage over time. In a new window:
   ```
   $ uv show -i infraRed.json -w stats > infraRed.html
   ```

3. Open `infraRed.html` in your browser of choice to see the charts
   below.

### Application usage timeline

![Application usage timeline](/assets/images/app_coarse.png)

### Detailed application window timeline

![Application usage timeline](/assets/images/app_fine.png)

### Aggregate time usage by app

![Application usage timeline](/assets/images/agg.png)


## Dependencies

uv's dependencies vary by system. See `uv dep` (mentioned in the installation instructions below).

## Install

1. [Install Go](https://golang.org/dl/) and run
   ```
   $ go get -u github.com/aimof/ultraRed/cmd/uv
   ```
1. Follow the instructions printed by `uv dep`.
   ```
   $ vu dep
   ```
   __ToDo__: `uv dep` should check dependency.

1. Verify `uv` works with
   ```
   $ uv track
   ```
   This should display JSON describing which applications are currently active, visible, and present on your system.

UV currently supports Linux(GNOME) and darwin.

## Use cases

UV was designed for developers who want to investigate their
application usage to make decisions that boost their day-to-day
productivity.

It can also be for other purposes such as:

- Tracking billable hours and constructing timesheets
- Studying application usage behavior in a given population

## How is UV different from other time trackers?

There are many time tracking products and services on the market.
UV differs from available offerings in the following ways:

- UV does not intend to be a fully featured time management product
  or service. UV adopts the Unix philosophy of a command-line tool
  that does one thing well and plays nicely with other command-line
  tools.

- UV does not require you to manually signal when you start or stop
  an activity. It automatically records which applications you use.

- UV is open source and free of charge.

- UV does not send data over the network. It stores the data it
  collects on local disk. It's up to you whether you want to share it
  or not.

## LICENSE

* original code is Licensed by Sourcegraph in MIT.
* Additional part is Licensed by Aito Shiroshita in MIT.
