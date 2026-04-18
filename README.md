# Aggregator

Aggregator is a CLI tool written in Go that allows the user to aggregate RSS feeds on a personally-selected timeframe and stores them in a Postgre SQL database. Users can follow feeds and browse through posts in the feeds.

## Installation

###Requirements
- Go (1.20+ recommended)
- PostgreSQL

###Steps

Clone the repository:

```bash
git clone https://github.com/scottw0173/aggregator.git
cd aggregator
go build
```

mv aggregator /usr/local/bin/   #This is an optional move to add aggregator to PATH

##Database Setup

Install goose (if not already installed):

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```
```
goose postgres "postgres://username:password@localhost:5432/dbname?sslmode=disable" up
```

#you will need to run this up migration 5 times to create the tables/relationships for the database

### Configuration
Create a `.gatorconfig.json` file in your home directory (not in the project directory):
It needs to contain the following information: 
    ```json
{
  "db_url": "postgres://username:password@localhost:5432/dbname?sslmode=disable",
  "user_name": ""
}  #the db_url requires username and password that are created by the user. 
   #here is the official documentation for PostgreSQL: https://www.postgresql.org/docs/8.0/sql-createuser.html
   #TLDR: sudo -u postgres createuser --interactive --pwprompt

#### Usage
if not added to path, then terminal command is ./aggregator
if added to PATH then terminal command is just aggregator (we will assume it is in PATH)

aggregator register <username> #creates user

aggregator login <username> #login as an existing user

aggregator users #lists all registered users (and denotes who is logged in currently)

aggregator agg <duration> #begins aggregation cycle based on duration given. Duration should be provided in the form of a string (e.g. 30s, 1d, 30m, etc)

aggregator addfeed <name> <url> #adds feed to be aggregated. Name is colloquial and url should be the RSS that the user wants added

aggregator feeds #lists all feeds that have been added and who added them

aggregator follow <url> #used to follow a feed that was already added by another user

aggregator following #lists all feeds added/followed by current user

aggregator unfollow <url> #unfollow RSS feed at given url for current user

aggregator browse <integer> #displays information for number of posts provided by user. Defaults to 2 if no number given

aggregator reset #deletes all registered users and feeds. clears all tables in the database. 