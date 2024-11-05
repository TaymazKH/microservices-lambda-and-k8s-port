#!/bin/bash -eu

npm install
zip -r deployment.zip node_modules genproto server.js currency_service.js
