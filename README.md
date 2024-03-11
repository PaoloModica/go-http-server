# Go HTTP server

## Overview
Go HTTP server application built following the indications provided by *Build an application* section of [Learn Go with tests](https://quii.gitbook.io/learn-go-with-tests/).


## Assignment

You have been asked to create a web server where users can track how many games players have won.
- `GET /players/{name}` should return a number indicating the total number of wins;
- `POST /players/{name}` should record a win for that name, incrementing for every subsequent `POST`.
- `GET /league` should return the list of all players, in `JSON` format.
