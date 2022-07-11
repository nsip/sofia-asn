#!/bin/bash

set -e

rm -rf ./partition/partition
rm -rf ./partition/out
rm -rf ./partition/*.log

rm -rf ./asn-json/asn-json
rm -rf ./asn-json/out
rm -rf ./asn-json/*.log

rm -rf ./asn-json-ld/asn-json-ld
rm -rf ./asn-json-ld/out*
rm -rf ./asn-json-ld/*.log

rm -rf ./precheck/*.log