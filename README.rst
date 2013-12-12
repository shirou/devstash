/dev/stash
===================

This tiny tool can store the information from stdin on the remote server.


usage
-------


Store
+++++++++

::

   % ifconfig eth0 | devstash
   % iostat | devstash
   % iostat | devstash -t stat,io  # store with 'stat' and 'io' tag

or specify the file itself.

::

   % devstash -f /path/to/something.xls


List
++++++

::

   % devstash -l

Search
++++++++

devstash can search information in the stash by using FM-Index text
search algoritm.

At first, you need to create index

::

   % devstash -make-index

Then, you can just type a search word with -s option

::

   % devstash -s <searchword>

   ex:
   % devstash -s MB
   Keyword "MB" found 1 entries. (0.000024secs)
          tty             ad0              ad1             cpu
    tin  tout  KB/t tps  MB/s   KB/t tps  MB/s  us ni sy in id
      1    90 13.29   7  0.09  50.15   6  0.28   6  0  2  0 92
    2 Hit(s)



How to use
-----------------


1. get
+++++++++++++++

::

   % go get github.com/shirou/devstash


TODO: make binary file to each architecture.


2. config
++++++++++++

Write config to ~/.devstash.cfg .

::

  [default]
  # choose one of those

  uri = https://example.com:9000   # http server
  # uri = ssh://example@example.com/path/to/storedir
  # uri = file:///path/to/storedir  # localhost

  [server]
  # used when act as a HTTP/HTTPS server
  port = 9000
  directory = /path/to/storedir


FM-Index
------------

Golang FM-Index implementation is comes from https://bitbucket.org/oov/go-shellinford/

Thank you!


Inspired by
------------

https://opslet.com/


nameing
-------------

At first, I want to name '/dev/null' but it too difficult to googling.


License
---------

2 clause BSD license


