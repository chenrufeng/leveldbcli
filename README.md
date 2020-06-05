LevelDB-CLI: a simple utility for debugging LevelDB
===========
Command-line utility for working with levelDB.
This utility is useful for debugging applications using the database LevelDB


![Demo GIF](https://raw.githubusercontent.com/liderman/leveldb-cli/master/docs/live-demo.gif)

Installation and build
----------------------

```
go get github.com/liderman/leveldb-cli
go install
```

Requirements
------------
 * `go1.5` or newer.

Usage
-----

```
# ./leveldb-cli
```

```
» open testdb
Database not exist! Create new database.
Success
testdb» set key100 value100
Success
testdb» set key200 value200
Success
testdb» set key300 value300
Success
testdb» set "key \"123" value
Success
testdb» show prefix key
Key	      | Value
key100	  | value100
key200	  | value200
key300	  | value300
key \"123 | value

testdb» show range key2 key3
Key	| Value
key200	| value200

testdb» close
Success
» exit
```

Commands
--------

### open
> open `DATABASE_NAME`

Opens database.
If the database does not exist, it is created.
You can use this method to create a new database.
 * `DATABASE_NAME` - The database name or path

### close
> close

It closes a previously opened database.

### set
> set `KEY` `VALUE`

Set the value of for a key.
 * `KEY` - The key
 * `VALUE` - The value

### delete
> delete `KEY`

Delete the record by key.
 * `KEY` - The key

### get
> get `KEY` [`FORMAT`] [`KEYFILTER`]

Display value by key.
 * `KEY` - The key
 * `FORMAT` - Data Display Format (Optional)

### show
> show -filter=[`filter`] -keyformat=[`format`] -format=[`format`] -start=[`start`] -end=[`end`] -prefix=[`prefix`] -contain=[`contain`]

for example,
> show -filter=all -contain=a


#### The list of filters available to use
 * `all` default filter, show all key-values
 * `prefix` display all values the keys that begin with thee prefix
 * `range` disokat all values, the keys of which are in the range between "start" and "end"


#### The list of formats available to display
 used by all filter
 * `ignore` - Not display
 * `raw` - Raw data without processing (default)
 * `base64` - Attempts to convert the data to base64 string
 * `bson` - Attempts to convert the data to be displayed from `bson` to `json`
 * `geohash` - Attempts to convert the data format of the `geohash` in the coordinates of the center (lat, lng)
 * `int64` - Attempts to display the data as an integer 64-bit number
 * `float64` - Attempts to display the data as a 64-bit number c with a floating point

#### start and end of range
 used only by range-filter
 * `START` - The key or key prefix indicating the beginning of the range
 * `LIMIT` - The key or key prefix indicating the end of the range

#### prefix
 used only by prefix-filter
 * display all value ,the key that begins with prefix

#### constain
 used by all filters
  * display all values ,that key contains specific string

### help
> help

Displays short usage software

### version
> version

Displays the current version of software and operating systems on which it runs

LICENSE
-------
Project distributes with standard MIT license
