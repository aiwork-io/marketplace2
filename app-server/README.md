# AIWork Marketplace

## Get started

1. [Setup a EC2 instance on AWS with Ubuntu OS](https://medium.com/nerd-for-tech/how-to-create-a-ubuntu-20-04-server-on-aws-ec2-elastic-cloud-computing-5b423b5bf635) You should select at least 100GB storage and t3.medium instance class
2. Export both AWS access key and private key to environment variables. For example

   ```
   export AWS_ACCESS_KEY_ID=AKIA4HKIEP2QVOQ2YVVX
   export AWS_SECRET_ACCESS_KEY=8+FywJQTki/bGB7Pzyp8BXmpcpT8FWG4/UWBszhw
   ```

3. Clone or extract source code to folder `aiworkmarketplace` then move working directory to that folder by command `cd aiworkmarketplace`
4. [Install Docker Ubuntu OS](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-20-04)
5. Start all of the components by command `docker compose up -d`
6. Make sure this command `curl localhost` return value `{"version":"2022.8.17"}`
7. Let's say your EC2 has an IP `18.138.103.117`, make sure the command `curl 18.138.103.117` return value `{"version":"2022.8.17"}` as well
