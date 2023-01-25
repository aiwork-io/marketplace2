require("dotenv").config();

import * as Server from "./Server";
import * as Logger from "./Logger";
import configs from "./configs";

const logger = Logger.create();
const server = Server.create(logger);
server.listen(configs.app.port, () =>
  logger.info(`RUN http://0.0.0.0:${configs.app.port}`)
);

process
  .on("unhandledRejection", (reason, p) => {
    console.error(reason, "Unhandled Rejection at Promise", p);
    process.exit(1);
  })
  .on("uncaughtException", (err) => {
    console.error(err, "Uncaught Exception thrown");
    process.exit(1);
  });
