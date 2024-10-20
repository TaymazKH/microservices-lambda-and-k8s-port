#!/bin/bash -eu

mkdir -p site_packages
pip install -r requirements.txt -t site_packages

cd site_packages
zip -r ../deployment.zip .
cd ..

zip -r deployment.zip genproto recommendation_service.py product_catalog_stub.py logger.py common.py server.py