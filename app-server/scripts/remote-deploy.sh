#!/bin/bash

set -e

ssh ubuntu@18.138.103.117 -i ~/.ssh/tuannguyen930708 'source ~/.profile; /home/ubuntu/aiworkmarketplace/scripts/redeploy.sh' || true
ssh ubuntu@18.138.103.117 -i ~/.ssh/id_kubeplusplus_x300 'source ~/.profile; /home/ubuntu/aiworkmarketplace/scripts/redeploy.sh' || true