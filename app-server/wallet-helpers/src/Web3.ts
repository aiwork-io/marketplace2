import Web3 from "web3";
import configs from "./configs";

export function create() {
  const web3 = new Web3(configs.web3.endpoint);

  const privateKey = process.env.APP_PRIVATE_KEY as string;
  if (!privateKey) throw new Error("Private key could not be empty!");

  const account = web3.eth.accounts.privateKeyToAccount(privateKey);
  web3.eth.accounts.wallet.add(account);
  web3.eth.defaultAccount = account.address;
  
  return web3;
}
