#!/bin/bash
set -e

peg -switch -inline unit_parser.peg
peg -switch -inline value_parser.peg
