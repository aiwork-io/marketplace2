import path from "path";
import pino from "pino";

export function create(name = "wallet-helpers") {
  return pino(
    { name },
    pino.destination(path.resolve(__dirname, "../logs/default.log"))
  );
}
