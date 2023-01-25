# AIWork Client

## Get started

1. Register an account. Let's say you will have a username`niceclient` with password `alonglongpassword`
2. Replace values of `AUTH_USERNAME` and `AUTH_PASSWORD` wiht the username and password you have. For example

   ```
   AUTH_USERNAME=niceclient
   AUTH_PASSWORD=alonglongpassword
   ```

3. Download the Docker image `capx-object` then import it by command `docker load < inference-engine.tar`
4. Start the client by command `docker compose up -d`
