# marketplace2
revamped marketplace

## Overview
There are currently 4 packages
app-docker
app-frontend-admin
app-frontend-user
app-server

App-docker contains the code for the docker client
App-frontend-admin & app-frontend-user contain code for the front end associated with the admin and normal user.  There is some difference in what the admin and normal user will see via the marketplace

App-server contains the backend code for the marketplace.

## Deployment
### Set-up infrastructure
1.	Setup a EC2 instance on AWS with Ubuntu OS: https://medium.com/nerd-for-tech/how-to-create-a-ubuntu-20-04-server-on-aws-ec2-elastic-cloud-computing-5b423b5bf635. 
You should select at least 100GB storage and t3.medium instance class
2.	Follow this link to get your AWS access key ID and secret: https://docs.aws.amazon.com/powershell/latest/userguide/pstools-appendix-sign-up.html
We will need these for step 3.
3.	Follow this link to install aws cli https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
Finalize the setup by running these commands:
 	export AWS_ACCESS_KEY_ID={access key id from step 2}
export AWS_SECRET_ACCESS_KEY={access secret from step 2} 
4.	Follow this guide to install docker to that ec2 instance: https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-20-04
NOTE: in production we will not setup docker for ec2, we will straight up use docker service provided by aws. But to save cost for this POC we have to do this hack.
5.	This is tricky, you need to get the source code to the ec2 instance. There r many ways to do this. 
Step 1: SSH into the ec2 instance, using this guide: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/AccessingInstancesLinux.html. Then u can either:
Step 2a: Install git => git login => git checkout the source code. 
Step 2b: use scp, following this guide https://docs.aws.amazon.com/managedservices/latest/appguide/qs-file-transfer.html
NOTE: For production we will not do this, we will be setting up a CICD pipeline instead.

### Run the code
1.	From the ec2 instance terminal, CD into the code folder. Start all of the components by command docker compose up -d
Note: Make sure this command curl localhost return value {"version":"2022.8.17"}. Let’s say your EC2 has an IP 18.138.103.117, make sure the command curl 18.138.103.117 return value {"version":"2022.8.17"} as well

