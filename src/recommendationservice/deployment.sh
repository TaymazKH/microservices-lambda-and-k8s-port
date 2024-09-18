#!/bin/bash -eu

mkdir -p site_packages
pip install -r requirements.txt -t site_packages

cd site_packages
zip -r ../deployment.zip .
cd ..

zip -r deployment.zip genproto dummy_email_service.py logger.py common.py server.py
