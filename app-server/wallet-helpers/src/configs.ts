export default {
  app: {
    port: Number(process.env.APP_PORT) || 3000,
    authToken: process.env.APP_AUTH_TOKEN || "changemenow",
  },
  web3: {
    endpoint: process.env.WEB3_ENDPOINT || "https://bsc-dataseed.binance.org/",
    chainId: Number(process.env.WEB3_CHAIN_ID) || 56,
    networkId: Number(process.env.WEB3_NETWORK_ID) || 56,
  },
  adresss: {
    token:
      process.env.ADDRESS_TOKEN || "0xa1daa25016d83895f215ee724c78fcc1085accb2",
  },
};
