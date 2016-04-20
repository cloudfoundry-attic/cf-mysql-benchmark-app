# CF MySQL Benchmarking App

The purpose of this app is to benchmark a deployment of a MySQL cluster. Benchmarking allows the user to collect useful metrics about the service, such as the number of read/write requests received per second, the number of DB transactions performed per second, the number of deadlock conditions hit, etc.

### How it works

The app works by invoking the [`sysbench` utility](https://launchpad.net/sysbench), which is a benchmark tool for evaluating a system running a database under intensive load. It creates a test database in the MySQL deployment, and then writes and reads back a given number of rows from it, while measuring the performance of these transactions.

`sysbench` can be run in isolation from the command line, but the app provides a user-friendly wrapper around it which works in tandem with the MySQL deployment. By deploying the app to the same subnet/AZ as your MySQL cluster, the app is able to avoid having to jump through firewalls.

### How to build and run the app using Cloud Foundry


##### Building the Dockerfile
The app runs on a Docker image which is described in `Dockerfile`. Among other things, the image has `sysbench` installed on it. The image can be built and pushed to Docker Hub by running the following command from the top-level directory of the app <br><br>
`docker build -t pivotalcf/cf-mysql-benchmark-app . && docker push pivotalcf/cf-mysql-benchmark-app`

##### Running the application

The application can be run on Diego, assuming that the Docker image has already been pushed with a simple<br>

`cf push -f <PATH_TO_APP_MANIFEST> -o pivotalcf/cf-mysql-benchmark-app`

A sample manifest file for the app is structured as follows:

```
---
applications:
  - name: sysbench
    env:
      MYSQL_HOSTS: "mysql_z1=<IP_OF_MYSQL_NODE_1>,mysql_z2=<IP_OF_MYSQL_NODE_2>,mysql_z3=<IP_OF_MYSQL_NODE_3>"
      MYSQL_PORT: <MYSQL_DATABASE_PORT>
      MYSQL_USER: <MYSQL_USERNAME>
      MYSQL_PASSWORD: <MYSQL_PASSWORD>
      TEST_DB: <NAME_OF_TEST_DATABASE_CREATED_BY_SYSBENCH>
      NUMBER_TEST_ROWS: <NUMBER OF ROWS >
      MAX_TIME: <NUMBER OF SECONDS TO RUN TEST>
      NUM_THREADS: <NUMBER OF CORES TO USE FOR TEST>
```

- `MYSQL_HOSTS` contains the IP addresses of the nodes in the MySQL cluster, MySQL Proxies, or MySQL cluster ELBs – essentially any IP address that the app can use to talk to the cluster.
- `MYSQL_PORT` the port on which to connect to the MySQL cluster
- `NUMBER_TEST_ROWS` specifies the number of rows we want sysbench to write to the test database it creates, the name of which is specified in `TEST_DB`.
- `TEST_DB` is the name of the database that should be used to run the tests against
- `MAX_TIME` is the length to run the `sysbench run` command, in SECONDS
- `NUM_THREADS` is the number of cores to use for the test

### Usage

`sysbench` works in two steps, called `prepare` and `run`. The `prepare` step logs in to the MySQL host, and sets up the database. The `run` step logs in to the MySQL host and runs the test against the specified host.

Analogously, the application has two endpoints that call through to the respective sysbench command.
- `-X POST /prepare/:node_index` where node_index is the 0-based position of the node you which to target as specified in your CF application manifest. This will call through to `sysbench prepare`
- `-X POST /start/:node_index` where node_index is the 0-based position of the node you which to target as specified in your CF application manifest. This will call through to `sysbench run`
