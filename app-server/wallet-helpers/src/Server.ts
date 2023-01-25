import pino from "pino";
import express from "express";
import cors from "cors";
import helmet from "helmet";
import { json } from "body-parser";
import { Transaction } from "ethereumjs-tx";
import EthCommon from "ethereumjs-common";

import * as Web3 from "./Web3";
import configs from "./configs";
const ABI = require("./ABI.json");

const web3 = Web3.create();

export function create(logger: pino.Logger) {
  const server = express();

  server.use(cors());
  server.use(json());
  server.use(helmet());

  server.use(async function auth(req, res, next) {
    const authToken = (req.query.auth_token as string) || req.get("AuthToken");

    if (authToken !== configs.app.authToken) {
      // IGNORE AUTH
      // return res.status(401).json({ error: "Permission defined!" });
    }

    return next();
  });

  server.get("/", (req, res) => {
    return res.json({ name: "wallet-helpers", version: "1.0.0" });
  });

  server.get("/check", async (req, res) => {
    const source = String(req.query.source);
    const walletId = req.query.wallet_id as string;
    const fromWalletId = req.query.from_wallet_id as string;
    const transactionId = req.query.transaction_id as string;
    const checkAmount = Number(req.query.amount as string);
    const paymentToken = String(req.query.payment_token || "").toLowerCase();

    logger.info(
      `[${source}] expected:${checkAmount} expected_token:${paymentToken}`
    );

    if (!walletId || !fromWalletId) {
      return res
        .status(400)
        .json({ error: "Both FROM and TO Wallet ID cannot be blank!" });
    }

    if (!transactionId) {
      return res.status(400).json({ error: "Transaction ID cannot be blank!" });
    }

    if (!Number.isFinite(checkAmount)) {
      return res
        .status(400)
        .json({ error: "Transaction amount cannot be blank!" });
    }

    const { from, to, amount, token } = await getTxnDetails(transactionId);
    let ok =
      from === fromWalletId.toLowerCase() &&
      to === walletId.toLowerCase() &&
      amount >= checkAmount;
    if (!!paymentToken) {
      ok = ok && paymentToken == token;
    }

    logger.info(
      `[${source}] ${ok} expected:${checkAmount} acutal:${amount} expected_token:${paymentToken} actual_token:${token}`
    );
    return res.json({ ok });
  });

  server.post("/transfer", async (req, res) => {
    const body = JSON.stringify(req.body);
    const recipientWallet = req.body.recipient_wallet as string;
    const amount = Number(req.body.amount as string) * 1e18;
    const currency = (req.query.currency as string) || "ECR20Token";

    logger.info(
      `body:${body} recipient_wallet:${recipientWallet}, amount:${amount}, currency:${currency}`
    );

    if (!recipientWallet) {
      return res
        .status(400)
        .json({ error: "Recipient wallet cannot be blank!" });
    }

    if (!Number.isInteger(amount)) {
      return res
        .status(400)
        .json({ error: "Transfer amount cannot be blank!" });
    }

    let result;

    if (currency === "bnb") {
      result = await sendBnB(recipientWallet, amount);
    } else {
      result = await sendECR20Token(recipientWallet, amount).catch(
        console.error
      );
    }

    if (!result) {
      return res
        .status(400)
        .json({ error: "Transfer amount is greater than current balance!" });
    }

    console.log(`Result ${JSON.stringify(result)}`);
    return res.json(result);
  });

  return server;
}

async function getTxnDetails(transactionId: string) {
  const receipt = await web3.eth.getTransactionReceipt(transactionId);

  // ECR20 token
  if (receipt.logs.length) {
    const fromHash = receipt.logs[0].topics[1];
    const toHash = receipt.logs[0].topics[2];
    const to = "0x" + toHash.slice(toHash.length - 40).toLowerCase();
    const from = "0x" + fromHash.slice(toHash.length - 40).toLowerCase();

    const amount =
      Number(web3.utils.hexToNumberString(receipt.logs[0].data)) / 1e18;

    const token = String(receipt.logs[0].address).toLowerCase();
    return { from, to, amount, token };
  }

  // BNB
  const transaction = await web3.eth.getTransaction(transactionId);
  const from = transaction.from?.toLocaleLowerCase();
  const to = transaction.to?.toLocaleLowerCase();
  const amount = Number(transaction.value) / 1e18;
  return { from, to, amount };
}

async function sendECR20Token(recipientWallet: string, amount: number) {
  const contract = new web3.eth.Contract(ABI, configs.adresss.token);
  // Check balance of first
  const balance = await contract.methods
    .balanceOf(web3.eth.defaultAccount)
    .call()
    .then(Number);
  if (balance < amount) return null;

  const gasPrice = await web3.eth.getGasPrice();

  const gas = await contract.methods
    .transfer(recipientWallet, amount.toString())
    .estimateGas({ from: web3.eth.defaultAccount });

  return contract.methods
    .transfer(recipientWallet, amount.toString())
    .send({ from: web3.eth.defaultAccount, gasPrice, gas });
}

async function sendBnB(recipientWallet: string, amount: number) {
  // Check balance of first
  const balance = await web3.eth
    .getBalance(web3.eth.defaultAccount as string)
    .then(Number);
  if (balance < amount) return null;

  const gasPrice = await web3.eth.getGasPrice().then(Number);
  const nonce = await web3.eth.getTransactionCount(
    web3.eth.defaultAccount as string
  );
  const rawTx = {
    from: web3.eth.defaultAccount,
    to: recipientWallet,
    nonce: web3.utils.toHex(nonce),
    gasPrice: web3.utils.toHex(gasPrice),
    gasLimit: web3.utils.toHex(3000000),
    value: amount,
    data: "0x",
  };
  const common = EthCommon.forCustomChain(
    "mainnet",
    {
      name: "bnb",
      networkId: configs.web3.networkId,
      chainId: configs.web3.chainId,
    },
    "petersburg"
  );
  var tx = new Transaction(rawTx, { common });

  const pk = Buffer.from(process.env.APP_PRIVATE_KEY as string, "hex");
  tx.sign(pk);
  const serializedTx = tx.serialize();

  return web3.eth.sendSignedTransaction("0x" + serializedTx.toString("hex"));
}
