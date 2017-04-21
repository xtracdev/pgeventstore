# pgeventstore

Event store using Postgres as the storage engine, ported (back) from
oraeventstore.

## Database Set Up

For the database set up, we assume the availability of a database user who can create schema
objects.

To install the schema, use [flyway](https://flywaydb.org/) to install 
the schema. Installation involves downloading the schema and dropping
the postgres JDBC jar into the flyway drivers directory.

Edit the flyway.conf in the db directory with your particulars, then from
the db directory run:

<pre>
flyway -user=esdbo -password=password -locations=filesystem:migration migrate
</pre>

## A Note on the Publish Table

The publish table simply writes the aggregate IDs of recently stored
aggregates, which picks up creation and updates. Another process will need
to read from the table to pick up the published aggregate, read the
actual data from the event store table, do something with it (publish it
to a queue, write out CQRS query views, etc), then delete the record from the
publish table.


## Dependencies

<pre>
go get github.com/xtracdev/goes
go get github.com/Sirupsen/logrus
go get github.com/gucumber/gucumber
go get github.com/stretchr/testify/assert
go get github.com/lib/pq
</pre>

## Contributing

To contribute, you must certify you agree with the [Developer Certificate of Origin](http://developercertificate.org/)
by signing your commits via `git -s`. To create a signature, configure your user name and email address in git.
Sign with your real name, do not use pseudonyms or submit anonymous commits.


In terms of workflow:

0. For significant changes or improvement, create an issue before commencing work.
1. Fork the respository, and create a branch for your edits.
2. Add tests that cover your changes, unit tests for smaller changes, acceptance test
for more significant functionality.
3. Run gofmt on each file you change before committing your changes.
4. Run golint on each file you change before committing your changes.
5. Make sure all the tests pass before committing your changes.
6. Commit your changes and issue a pull request.
## License

(c) 2017 Fidelity Investments
Licensed under the Apache License, Version 2.0
